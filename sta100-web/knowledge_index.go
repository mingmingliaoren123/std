package main

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

const (
	knowledgeVectorDimensions = 384
	knowledgeChunkRuneLimit   = 1200
	knowledgeChunkOverlap     = 160
	knowledgeMaxFileBytes     = 25 << 20
)

type knowledgeDocument struct {
	ID              string
	Visibility      string
	SourcePath      string
	Name            string
	Category        string
	Language        string
	SHA256          string
	Bytes           int64
	SourceUpdatedAt string
	IndexedAt       string
	Status          string
	ErrorText       string
}

type knowledgeHit struct {
	DocumentID string
	ChunkID    string
	Title      string
	SourcePath string
	Category   string
	Visibility string
	Content    string
	UpdatedAt  string
	Score      float64
}

func knowledgeDataRoot() (string, error) {
	root := strings.TrimSpace(os.Getenv("STA100_KNOWLEDGE_DATA_DIR"))
	if root == "" {
		privateRoot, err := privateDataRoot()
		if err != nil {
			return "", err
		}
		root = filepath.Join(filepath.Dir(privateRoot), "knowledge")
	}
	if !filepath.IsAbs(root) {
		return "", errors.New("知识库数据目录必须使用绝对路径")
	}
	if err := os.MkdirAll(root, 0700); err != nil {
		return "", err
	}
	return filepath.Clean(root), nil
}

func sharedKnowledgeRoot() (string, error) {
	if root := strings.TrimSpace(os.Getenv("STA100_SHARED_KNOWLEDGE_DIR")); root != "" {
		if !filepath.IsAbs(root) {
			return "", errors.New("共有知识库目录必须使用绝对路径")
		}
		return filepath.Clean(root), nil
	}
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(executable), "..", "knowledge", "shared"), nil
}

func knowledgePrivateFileRoot() (string, error) {
	root, err := privateDataRoot()
	if err != nil {
		return "", err
	}
	return filepath.Clean(root), nil
}

func allowedKnowledgeExtension(extension string) bool {
	switch strings.ToLower(strings.TrimSpace(extension)) {
	case ".md", ".txt", ".csv", ".tsv", ".json", ".xml", ".html", ".htm", ".docx", ".xlsx", ".rtf", ".pdf":
		return true
	default:
		return false
	}
}

func knowledgeUploadValidation(name string, size int64) error {
	extension := strings.ToLower(filepath.Ext(name))
	if !allowedKnowledgeExtension(extension) {
		return fmt.Errorf("不支持 %s；可上传 Markdown、TXT、CSV/TSV、JSON、XML、HTML、DOCX、XLSX、RTF 和 PDF", extension)
	}
	if size < 1 {
		return errors.New("文件为空，无法建立知识库索引")
	}
	if size > knowledgeMaxFileBytes {
		return fmt.Errorf("文件不能超过 %d MB", knowledgeMaxFileBytes>>20)
	}
	return nil
}

func (a *businessAPI) syncKnowledgeSources(ctx context.Context) (int, error) {
	count := 0
	sharedRoot, err := sharedKnowledgeRoot()
	if err != nil {
		return 0, err
	}
	var syncErrors []string
	if info, statErr := os.Stat(sharedRoot); statErr == nil && info.IsDir() {
		seen := map[string]bool{}
		if err := filepath.WalkDir(sharedRoot, func(path string, entry os.DirEntry, walkErr error) error {
			if walkErr != nil {
				syncErrors = append(syncErrors, path+": "+walkErr.Error())
				return nil
			}
			if entry == nil || entry.IsDir() || !entry.Type().IsRegular() || !allowedKnowledgeExtension(filepath.Ext(entry.Name())) {
				return nil
			}
			relative, relErr := filepath.Rel(sharedRoot, path)
			if relErr == nil {
				seen[filepath.ToSlash(relative)] = true
			}
			if indexErr := a.indexKnowledgeFile(ctx, "shared", sharedRoot, path, ""); indexErr != nil {
				syncErrors = append(syncErrors, path+": "+indexErr.Error())
			}
			count++
			return nil
		}); err != nil {
			return count, err
		}
		if cleanupErr := a.removeMissingKnowledgeDocuments(ctx, "shared", seen); cleanupErr != nil {
			syncErrors = append(syncErrors, cleanupErr.Error())
		}
	}
	privateRoot, err := knowledgePrivateFileRoot()
	if err != nil {
		return count, err
	}
	seen := map[string]bool{}
	if err := filepath.WalkDir(privateRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			syncErrors = append(syncErrors, path+": "+walkErr.Error())
			return nil
		}
		if entry == nil || entry.IsDir() || !entry.Type().IsRegular() || !allowedKnowledgeExtension(filepath.Ext(entry.Name())) {
			return nil
		}
		relative, relErr := filepath.Rel(privateRoot, path)
		if relErr == nil {
			seen[filepath.ToSlash(relative)] = true
		}
		if indexErr := a.indexKnowledgeFile(ctx, "private", privateRoot, path, ""); indexErr != nil {
			syncErrors = append(syncErrors, path+": "+indexErr.Error())
		}
		count++
		return nil
	}); err != nil {
		return count, err
	}
	if cleanupErr := a.removeMissingKnowledgeDocuments(ctx, "private", seen); cleanupErr != nil {
		syncErrors = append(syncErrors, cleanupErr.Error())
	}
	if len(syncErrors) > 0 {
		return count, fmt.Errorf("知识库同步完成但有 %d 个文件未处理：%s", len(syncErrors), strings.Join(syncErrors, "；"))
	}
	return count, nil
}

func (a *businessAPI) indexPrivateKnowledgeFile(ctx context.Context, item PrivateFile) error {
	path, err := privateFileStoragePath(item)
	if err != nil {
		return err
	}
	root, err := knowledgePrivateFileRoot()
	if err != nil {
		return err
	}
	return a.indexKnowledgeFile(ctx, "private", root, path, item.Category)
}

func (a *businessAPI) indexKnowledgeFile(ctx context.Context, visibility, root, path, category string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(root, path)
	if err != nil || strings.HasPrefix(relative, "..") {
		return errors.New("知识文件路径不在指定知识库目录内")
	}
	digest, err := fileSHA256(path)
	if err != nil {
		return err
	}
	if category == "" {
		category = knowledgeCategoryFromPath(relative)
	}
	updatedAt := info.ModTime().UTC().Format(time.RFC3339)
	var existing knowledgeDocument
	err = a.store.db.QueryRowContext(ctx, `SELECT id,visibility,source_path,name,category,language,sha256,bytes,source_updated_at,indexed_at,status,error_text FROM knowledge_documents WHERE visibility=? AND source_path=?`, visibility, filepath.ToSlash(relative)).Scan(
		&existing.ID, &existing.Visibility, &existing.SourcePath, &existing.Name, &existing.Category, &existing.Language, &existing.SHA256, &existing.Bytes, &existing.SourceUpdatedAt, &existing.IndexedAt, &existing.Status, &existing.ErrorText)
	if err == nil && existing.SHA256 == digest && existing.Status == "indexed" && existing.Category == category {
		return nil
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	content, parseErr := "", error(nil)
	status, errorText := "indexed", ""
	if validationErr := knowledgeUploadValidation(filepath.Base(path), info.Size()); validationErr != nil {
		status, errorText = "needs_processing", validationErr.Error()
	} else {
		content, parseErr = parseKnowledgeFile(path)
	}
	if parseErr != nil {
		status, errorText = "needs_processing", parseErr.Error()
		content = ""
	}
	documentID := "KD-" + shortKnowledgeHash(visibility+"\x00"+filepath.ToSlash(relative))
	now := time.Now().UTC().Format(time.RFC3339)
	tx, err := a.store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `INSERT INTO knowledge_documents(id,visibility,source_path,name,category,language,sha256,bytes,source_updated_at,indexed_at,status,error_text) VALUES(?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET category=excluded.category,language=excluded.language,sha256=excluded.sha256,bytes=excluded.bytes,source_updated_at=excluded.source_updated_at,indexed_at=excluded.indexed_at,status=excluded.status,error_text=excluded.error_text`, documentID, visibility, filepath.ToSlash(relative), filepath.Base(path), category, knowledgeLanguage(relative), digest, info.Size(), updatedAt, now, status, errorText); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM knowledge_chunk_agents WHERE chunk_id IN (SELECT id FROM knowledge_chunks WHERE document_id=?)`, documentID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM knowledge_chunks WHERE document_id=?`, documentID); err != nil {
		return err
	}
	if status == "indexed" {
		agents := knowledgeAgentsForCategory(category, relative)
		for number, chunk := range splitKnowledgeChunks(content) {
			chunkID := fmt.Sprintf("%s-%04d", documentID, number+1)
			if _, err := tx.ExecContext(ctx, `INSERT INTO knowledge_chunks(id,document_id,chunk_no,content,vector,updated_at) VALUES(?,?,?,?,?,?)`, chunkID, documentID, number+1, chunk, encodeKnowledgeVector(knowledgeVector(chunk)), now); err != nil {
				return err
			}
			for _, agentID := range agents {
				if _, err := tx.ExecContext(ctx, `INSERT INTO knowledge_chunk_agents(chunk_id,agent_id) VALUES(?,?)`, chunkID, agentID); err != nil {
					return err
				}
			}
		}
	}
	return tx.Commit()
}

func (a *businessAPI) removeMissingKnowledgeDocuments(ctx context.Context, visibility string, seen map[string]bool) error {
	rows, err := a.store.db.QueryContext(ctx, `SELECT id,source_path FROM knowledge_documents WHERE visibility=?`, visibility)
	if err != nil {
		return err
	}
	defer rows.Close()
	type documentRef struct{ id, path string }
	missing := make([]documentRef, 0)
	for rows.Next() {
		var item documentRef
		if err := rows.Scan(&item.id, &item.path); err != nil {
			return err
		}
		if !seen[item.path] {
			missing = append(missing, item)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(missing) == 0 {
		return nil
	}
	tx, err := a.store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, item := range missing {
		if _, err := tx.ExecContext(ctx, `DELETE FROM knowledge_chunk_agents WHERE chunk_id IN (SELECT id FROM knowledge_chunks WHERE document_id=?)`, item.id); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM knowledge_chunks WHERE document_id=?`, item.id); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM knowledge_documents WHERE id=?`, item.id); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (a *businessAPI) removePrivateKnowledgeDocument(ctx context.Context, item PrivateFile) error {
	root, err := knowledgePrivateFileRoot()
	if err != nil {
		return err
	}
	path, err := privateFileStoragePath(item)
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(root, path)
	if err != nil || strings.HasPrefix(relative, "..") {
		return errors.New("私有文件路径不在指定知识库目录内")
	}
	var documentID string
	err = a.store.db.QueryRowContext(ctx, `SELECT id FROM knowledge_documents WHERE visibility='private' AND source_path=?`, filepath.ToSlash(relative)).Scan(&documentID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return err
	}
	tx, err := a.store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM knowledge_chunk_agents WHERE chunk_id IN (SELECT id FROM knowledge_chunks WHERE document_id=?)`, documentID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM knowledge_chunks WHERE document_id=?`, documentID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM knowledge_documents WHERE id=?`, documentID); err != nil {
		return err
	}
	return tx.Commit()
}

func (a *businessAPI) privateKnowledgeStatus(ctx context.Context, item PrivateFile, indexErr error) string {
	if indexErr != nil {
		return "NeedsProcessing"
	}
	root, err := knowledgePrivateFileRoot()
	if err != nil {
		return "NeedsProcessing"
	}
	path, err := privateFileStoragePath(item)
	if err != nil {
		return "NeedsProcessing"
	}
	relative, err := filepath.Rel(root, path)
	if err != nil {
		return "NeedsProcessing"
	}
	var status string
	if err := a.store.db.QueryRowContext(ctx, `SELECT status FROM knowledge_documents WHERE visibility='private' AND source_path=?`, filepath.ToSlash(relative)).Scan(&status); err != nil {
		return "NeedsProcessing"
	}
	if status == "indexed" {
		return "Indexed"
	}
	return "NeedsProcessing"
}

func parseKnowledgeFile(path string) (string, error) {
	extension := strings.ToLower(filepath.Ext(path))
	switch extension {
	case ".md", ".txt", ".csv", ".tsv", ".json", ".xml", ".html", ".htm":
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		text := string(data)
		if extension == ".html" || extension == ".htm" || extension == ".xml" {
			text = stripMarkup(text)
		}
		return normalizeKnowledgeText(text), nil
	case ".docx", ".xlsx":
		return parseKnowledgeOfficeZip(path, extension)
	case ".rtf":
		data, err := os.ReadFile(path)
		return normalizeKnowledgeText(stripRTF(string(data))), err
	case ".pdf":
		return "", errors.New("PDF 已保存；当前部署未内置 OCR/文本提取器，需先转换为可检索文本或补充 OCR 组件")
	default:
		return "", errors.New("不支持的知识文件格式")
	}
}

func parseKnowledgeOfficeZip(path, extension string) (string, error) {
	reader, err := zip.OpenReader(path)
	if err != nil {
		return "", err
	}
	defer reader.Close()
	parts := make([]string, 0)
	for _, file := range reader.File {
		name := strings.ToLower(file.Name)
		wanted := extension == ".docx" && name == "word/document.xml" || extension == ".xlsx" && (name == "xl/sharedstrings.xml" || strings.HasPrefix(name, "xl/worksheets/"))
		if !wanted || file.UncompressedSize64 > knowledgeMaxFileBytes {
			continue
		}
		content, openErr := file.Open()
		if openErr != nil {
			return "", openErr
		}
		data, readErr := io.ReadAll(io.LimitReader(content, knowledgeMaxFileBytes+1))
		_ = content.Close()
		if readErr != nil {
			return "", readErr
		}
		parts = append(parts, stripMarkup(string(data)))
	}
	text := normalizeKnowledgeText(strings.Join(parts, "\n"))
	if text == "" {
		return "", errors.New("未从 Office 文件中提取到可检索文本")
	}
	return text, nil
}

var markupPattern = regexp.MustCompile(`(?s)<[^>]*>`)
var rtfControlPattern = regexp.MustCompile(`\\[a-zA-Z]+-?\d* ?|[{}]`)

func stripMarkup(value string) string {
	return html.UnescapeString(markupPattern.ReplaceAllString(value, " "))
}
func stripRTF(value string) string { return rtfControlPattern.ReplaceAllString(value, " ") }
func normalizeKnowledgeText(value string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(value), " "))
}

func truncateKnowledgeText(value string, limit int) string {
	runes := []rune(strings.TrimSpace(value))
	if limit < 1 || len(runes) <= limit {
		return string(runes)
	}
	return string(runes[:limit]) + "..."
}

func splitKnowledgeChunks(value string) []string {
	value = normalizeKnowledgeText(value)
	runes := []rune(value)
	chunks := make([]string, 0, len(runes)/knowledgeChunkRuneLimit+1)
	for start := 0; start < len(runes); {
		end := start + knowledgeChunkRuneLimit
		if end > len(runes) {
			end = len(runes)
		}
		chunks = append(chunks, string(runes[start:end]))
		if end == len(runes) {
			break
		}
		start = end - knowledgeChunkOverlap
	}
	return chunks
}

func knowledgeVector(value string) []float32 {
	vector := make([]float32, knowledgeVectorDimensions)
	for _, token := range knowledgeTokens(value) {
		hash := uint32(2166136261)
		for _, r := range token {
			hash ^= uint32(r)
			hash *= 16777619
		}
		index := hash % knowledgeVectorDimensions
		if hash&1 == 0 {
			vector[index] += 1
		} else {
			vector[index] -= 1
		}
	}
	norm := float32(0)
	for _, item := range vector {
		norm += item * item
	}
	if norm > 0 {
		norm = float32(math.Sqrt(float64(norm)))
		for index := range vector {
			vector[index] /= norm
		}
	}
	return vector
}

func knowledgeTokens(value string) []string {
	words := strings.FieldsFunc(strings.ToLower(value), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) })
	runes := []rune(strings.ToLower(value))
	for index := 0; index+1 < len(runes); index++ {
		if unicode.Is(unicode.Han, runes[index]) && unicode.Is(unicode.Han, runes[index+1]) {
			words = append(words, string(runes[index:index+2]))
		}
	}
	return words
}

func encodeKnowledgeVector(vector []float32) []byte {
	data := make([]byte, len(vector)*4)
	for index, value := range vector {
		binary.LittleEndian.PutUint32(data[index*4:], math.Float32bits(value))
	}
	return data
}

func decodeKnowledgeVector(data []byte) []float32 {
	if len(data) != knowledgeVectorDimensions*4 {
		return nil
	}
	vector := make([]float32, knowledgeVectorDimensions)
	for index := range vector {
		vector[index] = math.Float32frombits(binary.LittleEndian.Uint32(data[index*4:]))
	}
	return vector
}

func knowledgeCategoryFromPath(path string) string {
	parts := strings.Split(filepath.ToSlash(path), "/")
	if len(parts) > 1 {
		return parts[len(parts)-2]
	}
	return "未分类"
}

func knowledgeLanguage(path string) string {
	path = filepath.ToSlash(path)
	if strings.HasPrefix(path, "zh/") {
		return "zh"
	}
	if strings.HasPrefix(path, "en/") {
		return "en"
	}
	return "und"
}

func knowledgeAgentsForCategory(category, path string) []string {
	file := strings.ToLower(category + " " + path)
	agents := map[string]bool{"rag-agent": true}
	add := func(ids ...string) {
		for _, id := range ids {
			if assistantDomainAgentIDs[id] {
				agents[id] = true
			}
		}
	}
	switch {
	case containsAnyKnowledgeTerm(file, "客户", "customer", "渠道", "dealer", "distribution", "经销"):
		add("customer-measurement-agent", "b2b-marketplace-agent", "market-analyzer", "export-agent")
	case containsAnyKnowledgeTerm(file, "供应商", "supplier", "oem", "odm", "工厂"):
		add("supplier-aggregator", "price-tracker", "brand-value-crawler")
	case containsAnyKnowledgeTerm(file, "合同", "agreement", "贸易", "出口", "报关"):
		add("export-agent", "payment-advisor", "invoice-agent", "sinosure-advisor")
	case containsAnyKnowledgeTerm(file, "维修", "repair", "产品", "compatibility"):
		add("repair-qa", "compatibility-agent", "design-advisor", "inventory-agent")
	case containsAnyKnowledgeTerm(file, "品牌", "brand", "竞品", "价格", "price", "策略", "market"):
		add("brand-value-crawler", "price-tracker", "market-analyzer")
	case containsAnyKnowledgeTerm(file, "法规", "regulation", "cbam", "知识", "faq"):
		add("country-advisor", "cbam-calculator", "market-analyzer")
	default:
		add("market-analyzer")
	}
	result := make([]string, 0, len(agents))
	for id := range agents {
		result = append(result, id)
	}
	sort.Strings(result)
	return result
}

func fileSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, io.LimitReader(file, knowledgeMaxFileBytes+1)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}
func shortKnowledgeHash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])[:20]
}

func (a *businessAPI) searchKnowledge(ctx context.Context, agentIDs []string, query string, limit int) ([]knowledgeHit, error) {
	if len(agentIDs) == 0 || strings.TrimSpace(query) == "" {
		return nil, nil
	}
	if limit < 1 {
		limit = 8
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(agentIDs)), ",")
	args := make([]any, len(agentIDs))
	for i, id := range agentIDs {
		args[i] = id
	}
	rows, err := a.store.db.QueryContext(ctx, `SELECT DISTINCT c.id,c.document_id,d.name,d.source_path,d.category,d.visibility,c.content,c.vector,d.source_updated_at FROM knowledge_chunks c JOIN knowledge_chunk_agents ca ON ca.chunk_id=c.id JOIN knowledge_documents d ON d.id=c.document_id WHERE ca.agent_id IN (`+placeholders+`) AND d.status='indexed'`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	queryVector := knowledgeVector(query)
	hits := make([]knowledgeHit, 0)
	for rows.Next() {
		var hit knowledgeHit
		var vector []byte
		if err := rows.Scan(&hit.ChunkID, &hit.DocumentID, &hit.Title, &hit.SourcePath, &hit.Category, &hit.Visibility, &hit.Content, &vector, &hit.UpdatedAt); err != nil {
			return nil, err
		}
		chunkVector := decodeKnowledgeVector(vector)
		if chunkVector == nil {
			continue
		}
		for i := range queryVector {
			hit.Score += float64(queryVector[i] * chunkVector[i])
		}
		hit.Score += knowledgeLexicalScore(query, hit.Content)
		if hit.Score > 0 {
			hits = append(hits, hit)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Slice(hits, func(i, j int) bool { return hits[i].Score > hits[j].Score })
	if len(hits) > limit {
		hits = hits[:limit]
	}
	return hits, nil
}

func knowledgeLexicalScore(query, content string) float64 {
	score := 0.0
	lower := strings.ToLower(content)
	for _, token := range knowledgeTokens(query) {
		if len([]rune(token)) > 1 && strings.Contains(lower, token) {
			score += 0.08
		}
	}
	return math.Min(score, 0.35)
}

func (a *businessAPI) knowledgeEvidence(ctx context.Context, agentIDs []string, query string) []assistantEvidence {
	hits, err := a.searchKnowledge(ctx, agentIDs, query, 8)
	if err != nil {
		return nil
	}
	evidence := make([]assistantEvidence, 0, len(hits))
	for _, hit := range hits {
		label := map[string]string{"shared": "共有知识库", "private": "私有知识库"}[hit.Visibility]
		evidence = append(evidence, assistantEvidence{ID: hit.ChunkID, Entity: "knowledge", Title: hit.Title, Content: hit.Content, UpdatedAt: hit.UpdatedAt, Freshness: knowledgeFreshness(hit.UpdatedAt), Source: label + " / " + hit.Category})
	}
	return evidence
}

func knowledgeFreshness(updatedAt string) string {
	updated, err := time.Parse(time.RFC3339, strings.TrimSpace(updatedAt))
	if err != nil || updated.IsZero() {
		return "更新时间待核实"
	}
	age := time.Since(updated)
	if age < 0 {
		return "时间异常，需核实"
	}
	switch {
	case age > 365*24*time.Hour:
		return "历史资料，回答前需核实"
	case age > 90*24*time.Hour:
		return "较早资料，建议核实"
	default:
		return "近期资料"
	}
}

func (a *businessAPI) knowledgeIndexStatus(ctx context.Context) map[string]any {
	var documents, chunks, pending int
	_ = a.store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM knowledge_documents WHERE status='indexed'`).Scan(&documents)
	_ = a.store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM knowledge_chunks`).Scan(&chunks)
	_ = a.store.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM knowledge_documents WHERE status<>'indexed'`).Scan(&pending)
	failed := make([]map[string]any, 0, 5)
	rows, _ := a.store.db.QueryContext(ctx, `SELECT name,status,error_text FROM knowledge_documents WHERE status<>'indexed' ORDER BY name LIMIT 5`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var name, status, reason string
			if rows.Scan(&name, &status, &reason) == nil {
				failed = append(failed, map[string]any{"name": name, "status": status, "reason": reason})
			}
		}
	}
	root, _ := knowledgeDataRoot()
	shared, _ := sharedKnowledgeRoot()
	return map[string]any{"documents": documents, "chunks": chunks, "pending": pending, "pendingFiles": failed, "dataRoot": root, "sharedRoot": shared, "vectorVersion": "local-hybrid-zh-en-v1"}
}

// knowledgeCategories is intentionally derived from indexed documents. The
// OEM selector must describe the data that can actually be retrieved, rather
// than a catalogue maintained separately in the UI.
func (a *businessAPI) knowledgeCategories(ctx context.Context) []map[string]any {
	return a.knowledgeCategoriesForAgents(ctx, nil)
}

func (a *businessAPI) oemKnowledgeCategories(ctx context.Context) []map[string]any {
	// OEM 页面必须展示可直接用于需求匹配的产品品类，不能把“品牌、竞品、策略”
	// 这类文档标签误当作工厂生产分类。先从 OEM 专题知识库中识别骑行产品类目。
	rows, err := a.store.db.QueryContext(ctx, `SELECT d.name, d.category, c.content FROM knowledge_documents d
		JOIN knowledge_chunks c ON c.document_id=d.id
		JOIN knowledge_chunk_agents ca ON ca.chunk_id=c.id
		WHERE ca.agent_id=? AND d.status='indexed'`, "oem-match-agent")
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	counts := map[string]int{}
	for rows.Next() {
		var name, category, content string
		if rows.Scan(&name, &category, &content) != nil {
			continue
		}
		for _, item := range oemProductCategories(name + " " + category + " " + content) {
			counts[item]++
		}
	}
	order := []string{"公路整车", "山地整车", "E-bike 整车", "E-bike 电池", "中置电机", "轮组", "车架", "头盔", "骑行服", "自行车配件"}
	result := make([]map[string]any, 0, len(order))
	for _, category := range order {
		if counts[category] > 0 {
			result = append(result, map[string]any{"value": category, "label": category, "documents": counts[category]})
		}
	}
	if len(result) == 0 {
		// No OEM-specific source has been indexed yet. Keep the selector useful
		// with the supported OEM business taxonomy, without claiming data exists.
		for _, category := range []string{"公路整车", "E-bike 电池", "中置电机", "头盔"} {
			result = append(result, map[string]any{"value": category, "label": category, "documents": 0})
		}
	}
	return result
}

func oemProductCategories(value string) []string {
	value = strings.ToLower(value)
	match := func(terms ...string) bool { return containsAnyKnowledgeTerm(value, terms...) }
	items := []string{}
	switch {
	case match("road bike", "公路", "road frame"):
		items = append(items, "公路整车")
	case match("mountain bike", "mtb", "山地"):
		items = append(items, "山地整车")
	case match("e-bike", "ebike", "电助力", "电动自行车"):
		items = append(items, "E-bike 整车")
	}
	if match("battery", "电池", "battery pack") {
		items = append(items, "E-bike 电池")
	}
	if match("mid-drive", "中置电机", "motor", "电机") {
		items = append(items, "中置电机")
	}
	if match("wheelset", "wheel", "轮组") {
		items = append(items, "轮组")
	}
	if match("frame", "车架") {
		items = append(items, "车架")
	}
	if match("helmet", "头盔", "mips") {
		items = append(items, "头盔")
	}
	if match("apparel", "jersey", "骑行服") {
		items = append(items, "骑行服")
	}
	if match("accessor", "配件", "pedal", "车把") {
		items = append(items, "自行车配件")
	}
	return items
}

func (a *businessAPI) knowledgeCategoriesForAgents(ctx context.Context, agentIDs []string) []map[string]any {
	query := `SELECT d.category, COUNT(DISTINCT d.id) FROM knowledge_documents d`
	args := make([]any, 0, len(agentIDs))
	if len(agentIDs) > 0 {
		placeholders := strings.TrimRight(strings.Repeat("?,", len(agentIDs)), ",")
		query += ` JOIN knowledge_chunks c ON c.document_id=d.id JOIN knowledge_chunk_agents ca ON ca.chunk_id=c.id WHERE ca.agent_id IN (` + placeholders + `) AND d.status='indexed' AND TRIM(d.category)<>'' GROUP BY d.category ORDER BY d.category`
		for _, agentID := range agentIDs {
			args = append(args, agentID)
		}
	} else {
		query += ` WHERE d.status='indexed' AND TRIM(d.category)<>'' GROUP BY d.category ORDER BY d.category`
	}
	rows, err := a.store.db.QueryContext(ctx, query, args...)
	if err != nil {
		return []map[string]any{}
	}
	defer rows.Close()
	categories := make([]map[string]any, 0)
	for rows.Next() {
		var category string
		var count int
		if rows.Scan(&category, &count) != nil {
			continue
		}
		category = strings.TrimSpace(category)
		if category == "" {
			continue
		}
		categories = append(categories, map[string]any{"value": category, "label": category, "documents": count})
	}
	return categories
}

// Keep json imported by this file for future manifest compatibility and to
// ensure the index format is deliberately local and serializable.
var _ = json.Valid
