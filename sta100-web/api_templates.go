package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
)

const businessTemplateKind = "templates"
const builtInQuotePDFTemplateID = "QUOTE-TPL-STRAT-Q2026-0819"
const builtInQuotePDFOriginalAsset = "assets/STRAT-Q2026-0819-001-SD100-NL---198b86ea-260a-40f7-b3c6-faee0b9c912d.pdf"
const builtInQuotePDFTemplateAsset = "assets/STRAT-Q2026-0819-001-SD100-NL-standard-clean.pdf"
const builtInPIPDFTemplateID = "DOCUMENT-TPL-PI-20260625"
const builtInPIPDFTemplateAsset = "assets/PI-20260625-DW01.pdf"
const builtInDocumentCustomsTemplateID = "DOCUMENT-TPL-CUSTOMS-20260826"
const builtInDocumentContractTemplateID = "DOCUMENT-TPL-CONTRACT-20260826"
const builtInDocumentInvoiceTemplateID = "DOCUMENT-TPL-INVOICE-20260826"
const builtInDocumentPackingTemplateID = "DOCUMENT-TPL-PACKING-20260826"
const builtInDocumentWorkbookTemplateAsset = "assets/STA100-document-templates.xlsm"
const stratronixQuoteOriginalTemplateSHA256 = "8094e3d29424e2287f846a2cc1f9558c31ffc1a972de3227b8dc4bb1ae4aced4"
const stratronixQuoteCleanTemplateSHA256 = "b4cbc23eb8beca9055ed1eb2f52cacdda8cc5b97a1ce1d80e90ea9714b9f0cf3"

var templatePlaceholderPattern = regexp.MustCompile(`\{\{\s*([a-zA-Z0-9_.-]+)\s*\}\}`)

var pdfFontInstallMu sync.Mutex

func isStratronixQuotePDFTemplate(item BusinessTemplate) bool {
	if item.Kind != "quote" {
		return false
	}
	hash := strings.TrimSpace(item.SHA256)
	return item.ID == builtInQuotePDFTemplateID ||
		strings.EqualFold(hash, stratronixQuoteOriginalTemplateSHA256) ||
		strings.EqualFold(hash, stratronixQuoteCleanTemplateSHA256)
}

func readEmbeddedUIFile(assets embed.FS, name string) ([]byte, error) {
	data, err := assets.ReadFile(name)
	if err == nil {
		return data, nil
	}
	fileData, fileErr := os.ReadFile(filepath.Join("sta100-web", name))
	if fileErr == nil {
		return fileData, nil
	}
	fileData, fileErr = os.ReadFile(name)
	if fileErr == nil {
		return fileData, nil
	}
	return nil, fmt.Errorf("load embedded asset %s: embed: %v; file: %w", name, err, fileErr)
}

func builtInDocumentTemplateIDs() []string {
	return []string{
		builtInPIPDFTemplateID,
		builtInDocumentCustomsTemplateID,
		builtInDocumentContractTemplateID,
		builtInDocumentInvoiceTemplateID,
		builtInDocumentPackingTemplateID,
	}
}

func isBuiltInDocumentTemplateID(id string) bool {
	for _, candidate := range builtInDocumentTemplateIDs() {
		if candidate == id {
			return true
		}
	}
	return false
}

func (a *businessAPI) ensureBuiltInBusinessTemplates(ctx context.Context) error {
	if a == nil || a.store == nil {
		return nil
	}
	if err := a.ensureBuiltInQuotePDFTemplate(ctx, ui); err != nil {
		return err
	}
	if err := a.ensureBuiltInDocumentTemplates(ctx, ui); err != nil {
		return err
	}
	return a.ensurePDFBusinessTemplates(ctx)
}

func (a *businessAPI) ensureBuiltInQuotePDFTemplate(ctx context.Context, assets embed.FS) error {
	data, err := readEmbeddedUIFile(assets, builtInQuotePDFTemplateAsset)
	if err != nil {
		return err
	}
	root, err := templateStorageRoot()
	if err != nil {
		return err
	}
	storageName := builtInQuotePDFTemplateID + ".pdf"
	target := filepath.Join(root, storageName)
	if existing, err := os.ReadFile(target); err != nil || !bytes.Equal(existing, data) {
		if err := os.WriteFile(target, data, 0o600); err != nil {
			return err
		}
	}
	hash := sha256.Sum256(data)
	templates, err := listRecords[BusinessTemplate](ctx, a.store, businessTemplateKind)
	if err != nil {
		return err
	}
	hasQuoteDefault := false
	for _, item := range templates {
		if item.Kind == "quote" && item.Default && !strings.EqualFold(item.Status, "Archived") && item.ID != builtInQuotePDFTemplateID {
			hasQuoteDefault = true
			break
		}
	}
	now := currentText()
	item := BusinessTemplate{
		ID: builtInQuotePDFTemplateID, Kind: "quote", Name: "STRATRONIX 标准报价单 PDF 模板", Version: "2026-08-19",
		Source: "系统内置", Mode: "builtin-file", FileName: filepath.Base(builtInQuotePDFOriginalAsset), Mime: "application/pdf",
		Size: humanBytes(int64(len(data))), Bytes: int64(len(data)), SHA256: hex.EncodeToString(hash[:]),
		StorageFileName: storageName, Status: "Active", Default: !hasQuoteDefault, OutputFormat: "pdf",
		Placeholders: defaultTemplatePlaceholders("quote"), FieldMapping: defaultTemplateFieldMapping("quote", defaultTemplatePlaceholders("quote")),
		DefaultValues: defaultQuoteTemplateValues(), PDFFields: defaultPDFTemplateFields("quote"), CreatedAt: now, Updated: now,
		RecognitionNote: "已内置客户提供的 STRATRONIX 标准报价单 PDF；系统使用清理后的标准底版按坐标写入报价单数据，避免示例值重叠。",
	}
	var existing BusinessTemplate
	getErr := a.store.get(ctx, businessTemplateKind, builtInQuotePDFTemplateID, &existing)
	if getErr == nil {
		existing.Name = item.Name
		existing.Version = item.Version
		existing.Source = item.Source
		existing.Mode = item.Mode
		existing.FileName = item.FileName
		existing.Mime = item.Mime
		existing.Size = item.Size
		existing.Bytes = item.Bytes
		existing.SHA256 = item.SHA256
		existing.StorageFileName = item.StorageFileName
		existing.OutputFormat = item.OutputFormat
		existing.Placeholders = item.Placeholders
		existing.FieldMapping = item.FieldMapping
		existing.DefaultValues = mergeTemplateDefaultValues(existing.DefaultValues, item.DefaultValues)
		existing.PDFFields = item.PDFFields
		if existing.Status == "" || strings.EqualFold(existing.Status, "Draft") {
			existing.Status = "Active"
		}
		if !hasQuoteDefault && !existing.Default {
			existing.Default = true
		}
		existing.RecognitionNote = item.RecognitionNote
		existing.Updated = now
		return a.store.put(ctx, businessTemplateKind, builtInQuotePDFTemplateID, existing)
	}
	if errors.Is(getErr, errRecordNotFound) {
		return a.store.create(ctx, businessTemplateKind, builtInQuotePDFTemplateID, item)
	}
	return getErr
}

func (a *businessAPI) ensureBuiltInDocumentTemplates(ctx context.Context, assets embed.FS) error {
	defaultValues := defaultDocumentTemplateValues()
	if err := a.ensureBuiltInFileTemplate(ctx, assets, builtInPIPDFTemplateID, "document", "PI 标准 PDF 模板", "2026-06-25", builtInPIPDFTemplateAsset, "application/pdf", "pdf", false, defaultValues, defaultPDFTemplateFields("pi"), "已内置客户提供的 PI PDF 标准模板；系统按固定坐标叠加 PI、客户、订单和产品字段。", "PI"); err != nil {
		return err
	}
	if err := a.ensureBuiltInFileTemplate(ctx, assets, builtInDocumentCustomsTemplateID, "document", "报关单标准模板", "2026-08-26", builtInDocumentWorkbookTemplateAsset, "application/vnd.ms-excel.sheet.macroEnabled.12", "xlsm", true, defaultValues, nil, "已内置客户提供的单据模板 XLSM 报关单页；系统写入隐藏资料页并保留原公式和格式。", "报关单"); err != nil {
		return err
	}
	if err := a.ensureBuiltInFileTemplate(ctx, assets, builtInDocumentContractTemplateID, "document", "合同标准模板", "2026-08-26", builtInDocumentWorkbookTemplateAsset, "application/vnd.ms-excel.sheet.macroEnabled.12", "xlsm", false, defaultValues, nil, "已内置客户提供的单据模板 XLSM 合同页；系统写入隐藏资料页并保留原公式和格式。", "合同"); err != nil {
		return err
	}
	if err := a.ensureBuiltInFileTemplate(ctx, assets, builtInDocumentInvoiceTemplateID, "document", "发票 CI 标准模板", "2026-08-26", builtInDocumentWorkbookTemplateAsset, "application/vnd.ms-excel.sheet.macroEnabled.12", "xlsm", false, defaultValues, nil, "已内置客户提供的单据模板 XLSM 发票页；系统写入隐藏资料页并保留原公式和格式。", "CI"); err != nil {
		return err
	}
	return a.ensureBuiltInFileTemplate(ctx, assets, builtInDocumentPackingTemplateID, "document", "装箱单 PL 标准模板", "2026-08-26", builtInDocumentWorkbookTemplateAsset, "application/vnd.ms-excel.sheet.macroEnabled.12", "xlsm", false, defaultValues, nil, "已内置客户提供的单据模板 XLSM 装箱单页；系统写入隐藏资料页并保留原公式和格式。", "PL")
}

func (a *businessAPI) ensureBuiltInFileTemplate(ctx context.Context, assets embed.FS, id, kind, name, version, asset, mimeType, outputFormat string, defaultTemplate bool, defaultValues map[string]string, pdfFields []PDFTemplateField, recognitionNote string, documentType string) error {
	data, err := readEmbeddedUIFile(assets, asset)
	if err != nil {
		return err
	}
	root, err := templateStorageRoot()
	if err != nil {
		return err
	}
	ext := strings.ToLower(filepath.Ext(asset))
	storageName := id + ext
	target := filepath.Join(root, storageName)
	if existing, err := os.ReadFile(target); err != nil || !bytes.Equal(existing, data) {
		if err := os.WriteFile(target, data, 0o600); err != nil {
			return err
		}
	}
	hash := sha256.Sum256(data)
	now := currentText()
	item := BusinessTemplate{
		ID: id, Kind: kind, Name: name, Version: version, Source: "系统内置", Mode: "builtin-file",
		FileName: filepath.Base(asset), Mime: mimeType, Size: humanBytes(int64(len(data))), Bytes: int64(len(data)),
		SHA256: hex.EncodeToString(hash[:]), StorageFileName: storageName, Status: "Active", Default: defaultTemplate,
		OutputFormat: outputFormat, Placeholders: defaultTemplatePlaceholders(kind), FieldMapping: defaultTemplateFieldMapping(kind, defaultTemplatePlaceholders(kind)),
		DefaultValues: defaultValues, PDFFields: pdfFields, CreatedAt: now, Updated: now, RecognitionNote: recognitionNote, DocumentType: documentType,
	}
	var existing BusinessTemplate
	getErr := a.store.get(ctx, businessTemplateKind, id, &existing)
	if getErr == nil {
		existing.Name = item.Name
		existing.Version = item.Version
		existing.Source = item.Source
		existing.Mode = item.Mode
		existing.FileName = item.FileName
		existing.Mime = item.Mime
		existing.Size = item.Size
		existing.Bytes = item.Bytes
		existing.SHA256 = item.SHA256
		existing.StorageFileName = item.StorageFileName
		existing.Status = "Active"
		existing.OutputFormat = item.OutputFormat
		existing.Placeholders = item.Placeholders
		existing.FieldMapping = item.FieldMapping
		existing.DefaultValues = mergeTemplateDefaultValues(existing.DefaultValues, item.DefaultValues)
		existing.PDFFields = item.PDFFields
		existing.Default = item.Default
		existing.RecognitionNote = item.RecognitionNote
		existing.DocumentType = item.DocumentType
		existing.Updated = now
		return a.store.put(ctx, businessTemplateKind, id, existing)
	}
	if errors.Is(getErr, errRecordNotFound) {
		return a.store.create(ctx, businessTemplateKind, id, item)
	}
	return getErr
}

func (a *businessAPI) ensurePDFBusinessTemplates(ctx context.Context) error {
	items, err := listRecords[BusinessTemplate](ctx, a.store, businessTemplateKind)
	if err != nil {
		return err
	}
	now := currentText()
	for _, item := range items {
		if !isPDFTemplateFile(item.FileName) || strings.EqualFold(item.Status, "Archived") {
			continue
		}
		changed := false
		isStratronixQuote := isStratronixQuotePDFTemplate(item)
		if item.OutputFormat != "pdf" {
			item.OutputFormat = "pdf"
			changed = true
		}
		if isStratronixQuote || len(item.Placeholders) == 0 || len(item.Placeholders) < len(defaultTemplatePlaceholders(item.Kind)) {
			item.Placeholders = defaultTemplatePlaceholders(item.Kind)
			changed = true
		}
		if isStratronixQuote || len(item.FieldMapping) == 0 {
			item.FieldMapping = defaultTemplateFieldMapping(item.Kind, item.Placeholders)
			changed = true
		}
		if isStratronixQuote {
			merged := mergeTemplateDefaultValues(item.DefaultValues, defaultQuoteTemplateValues())
			if !stringMapEqual(item.DefaultValues, merged) {
				item.DefaultValues = merged
				changed = true
			}
		}
		if isStratronixQuote || len(item.PDFFields) == 0 {
			item.PDFFields = defaultPDFTemplateFields(item.Kind)
			changed = true
		} else {
			normalized := normalizePDFTemplateFields(item.PDFFields)
			if len(normalized) != len(item.PDFFields) {
				item.PDFFields = normalized
				changed = true
			}
		}
		if changed {
			item.Updated = now
			if item.RecognitionNote == "" || strings.Contains(item.RecognitionNote, "当前版本用于模板归档") {
				item.RecognitionNote = "已兼容升级为 PDF 模板；可按 PDF 表单字段或坐标字段生成原版式文件。"
			}
			if err := a.store.put(ctx, businessTemplateKind, item.ID, item); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *businessAPI) templatesRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.templatesCollection(w, r)
		return
	}
	switch parts[0] {
	case "upload":
		a.templateUpload(w, r, false)
	case "image-recognition":
		a.templateUpload(w, r, true)
	default:
		if len(parts) == 2 && parts[1] == "publish" {
			a.templatePublish(w, r, parts[0])
			return
		}
		if len(parts) == 2 && parts[1] == "render" {
			a.templateRender(w, r, parts[0])
			return
		}
		a.templateItem(w, r, parts[0])
	}
}

func (a *businessAPI) templatesCollection(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	items, err := listRecords[BusinessTemplate](r.Context(), a.store, businessTemplateKind)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	kind := normalizeTemplateKind(r.URL.Query().Get("kind"))
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	filtered := make([]BusinessTemplate, 0, len(items))
	for _, item := range items {
		if kind != "" && item.Kind != kind {
			continue
		}
		if kind == "quote" && item.ID != builtInQuotePDFTemplateID {
			continue
		}
		if kind == "document" && !isBuiltInDocumentTemplateID(item.ID) {
			continue
		}
		if query != "" && !containsFold(item.Name+" "+item.FileName+" "+item.Version+" "+item.Source, query) {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].Default != filtered[j].Default {
			return filtered[i].Default
		}
		return filtered[i].Updated > filtered[j].Updated
	})
	writeJSON(w, http.StatusOK, map[string]any{"items": filtered, "total": len(filtered), "status": "ready"})
}

func (a *businessAPI) templateUpload(w http.ResponseWriter, r *http.Request, imageMode bool) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	if err := r.ParseMultipartForm(24 << 20); err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_TEMPLATE_UPLOAD", "模板文件不能超过 20 MB")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "MISSING_TEMPLATE_FILE", "请选择模板文件")
		return
	}
	defer file.Close()
	kind := normalizeTemplateKind(r.FormValue("kind"))
	if kind == "" {
		writeAPIError(w, http.StatusBadRequest, "INVALID_TEMPLATE_KIND", "模板类型必须是 quote、order 或 document")
		return
	}
	if kind == "quote" {
		writeAPIError(w, http.StatusConflict, "STANDARD_QUOTE_TEMPLATE_LOCKED", "报价单固定使用系统内置 STRATRONIX 标准模板，不再支持上传或替换模板")
		return
	}
	if kind == "order" {
		writeAPIError(w, http.StatusConflict, "STANDARD_ORDER_TEMPLATE_LOCKED", "订单不单独维护模板，请在单据中按订单生成 PI、CI、PL、合同和报关单")
		return
	}
	if kind == "document" {
		writeAPIError(w, http.StatusConflict, "STANDARD_DOCUMENT_TEMPLATE_LOCKED", "单据固定使用系统内置 PI PDF、报关单、合同、发票和装箱单标准模板，不再支持上传或替换模板")
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedTemplateExtension(ext, imageMode) {
		if imageMode {
			writeAPIError(w, http.StatusBadRequest, "INVALID_TEMPLATE_TYPE", "图片模板仅支持 JPG、PNG、WebP")
		} else {
			writeAPIError(w, http.StatusBadRequest, "INVALID_TEMPLATE_TYPE", "文件模板仅支持 DOCX、XLSX、PDF、HTML、TXT、MD")
		}
		return
	}
	root, err := templateStorageRoot()
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	now := currentText()
	id, err := a.store.nextSequence(r.Context(), businessTemplateKind, strings.ToUpper(kind)+"-TPL", 4)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	var content bytes.Buffer
	hash := sha256.New()
	written, err := io.Copy(io.MultiWriter(&content, hash), io.LimitReader(file, 20*1024*1024+1))
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	if written > 20*1024*1024 {
		writeAPIError(w, http.StatusBadRequest, "TEMPLATE_TOO_LARGE", "模板文件不能超过 20 MB")
		return
	}
	digest := hex.EncodeToString(hash.Sum(nil))
	storageName := id + ext
	path := filepath.Join(root, storageName)
	if err := os.WriteFile(path, content.Bytes(), 0o600); err != nil {
		writeBusinessError(w, err)
		return
	}
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = mime.TypeByExtension(ext)
	}
	name := strings.TrimSuffix(filepath.Base(header.Filename), ext)
	if strings.TrimSpace(name) == "" {
		name = templateKindLabel(kind) + "模板"
	}
	item := BusinessTemplate{
		ID: id, Kind: kind, Name: name, Version: "v1", Source: map[bool]string{true: "图片识别上传", false: "文件上传"}[imageMode],
		Mode: map[bool]string{true: "image", false: "file"}[imageMode], FileName: filepath.Base(header.Filename), Mime: mimeType,
		Size: humanBytes(written), Bytes: written, SHA256: digest, Path: path, StorageFileName: storageName, Status: "Draft", OutputFormat: defaultTemplateOutputFormat(ext),
		Placeholders: templatePlaceholders(ext, content.Bytes()), RecognitionNote: templateRecognitionNote(imageMode, ext),
		CreatedAt: now, Updated: now,
	}
	item.FieldMapping = defaultTemplateFieldMapping(item.Kind, item.Placeholders)
	if ext == ".pdf" {
		item.Placeholders = defaultTemplatePlaceholders(item.Kind)
		item.FieldMapping = defaultTemplateFieldMapping(item.Kind, item.Placeholders)
		item.PDFFormFields = inspectPDFFormFieldNames(path)
		item.PDFFields = defaultPDFTemplateFields(item.Kind)
	}
	if err := a.store.create(r.Context(), businessTemplateKind, id, item); err != nil {
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "upload", "template", id, requestOperator(r), map[string]any{"kind": kind, "sha256": digest})
	writeJSON(w, http.StatusCreated, item)
}

func (a *businessAPI) templateItem(w http.ResponseWriter, r *http.Request, id string) {
	var item BusinessTemplate
	if err := a.store.get(r.Context(), businessTemplateKind, id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, item)
	case http.MethodPatch:
		if !allowMutation(w, r, http.MethodPatch) {
			return
		}
		var request BusinessTemplate
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		item.Name = firstNonEmpty(request.Name, item.Name)
		item.Version = firstNonEmpty(request.Version, item.Version)
		item.Status = firstNonEmpty(request.Status, item.Status)
		item.OutputFormat = firstNonEmpty(request.OutputFormat, item.OutputFormat)
		if request.FieldMapping != nil {
			item.FieldMapping = request.FieldMapping
		}
		if request.DefaultValues != nil {
			item.DefaultValues = sanitizeTemplateDefaultValues(item.Kind, request.DefaultValues)
		}
		if request.PDFFields != nil {
			item.PDFFields = normalizePDFTemplateFields(request.PDFFields)
		}
		if request.PDFFormFields != nil {
			item.PDFFormFields = normalizePDFFormFields(request.PDFFormFields)
		}
		if request.Placeholders != nil {
			item.Placeholders = request.Placeholders
		}
		item.Updated = currentText()
		if err := validateBusinessTemplate(item); err != nil {
			writeAPIError(w, http.StatusBadRequest, "INVALID_TEMPLATE", err.Error())
			return
		}
		if err := a.store.put(r.Context(), businessTemplateKind, id, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "template", id, requestOperator(r), map[string]any{"kind": item.Kind})
		writeJSON(w, http.StatusOK, item)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if item.Mode == "builtin-file" {
			writeAPIError(w, http.StatusConflict, "BUILTIN_TEMPLATE_DELETE_FORBIDDEN", "系统内置标准模板不能删除，只能维护默认填充值")
			return
		}
		if item.Default {
			writeAPIError(w, http.StatusConflict, "DEFAULT_TEMPLATE_DELETE_FORBIDDEN", "默认模板不能直接删除，请先发布其它模板为默认")
			return
		}
		if err := a.store.softDelete(r.Context(), businessTemplateKind, id); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "archive", "template", id, requestOperator(r), map[string]any{"kind": item.Kind})
		writeJSON(w, http.StatusOK, map[string]any{"archived": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) templatePublish(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	item, err := a.publishTemplate(r.Context(), id)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "publish", "template", id, requestOperator(r), map[string]any{"kind": item.Kind})
	writeJSON(w, http.StatusOK, item)
}

func (a *businessAPI) publishTemplate(ctx context.Context, id string) (BusinessTemplate, error) {
	var item BusinessTemplate
	if err := a.store.get(ctx, businessTemplateKind, id, &item); err != nil {
		return BusinessTemplate{}, err
	}
	items, err := listRecords[BusinessTemplate](ctx, a.store, businessTemplateKind)
	if err != nil {
		return BusinessTemplate{}, err
	}
	now := currentText()
	for _, existing := range items {
		if existing.Kind == item.Kind && existing.ID != item.ID && existing.Default {
			existing.Default = false
			existing.Updated = now
			if err := a.store.put(ctx, businessTemplateKind, existing.ID, existing); err != nil {
				return BusinessTemplate{}, err
			}
		}
	}
	item.Default, item.Status, item.Updated = true, "Active", now
	if err := a.store.put(ctx, businessTemplateKind, item.ID, item); err != nil {
		return BusinessTemplate{}, err
	}
	return item, nil
}

func (a *businessAPI) templateRender(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	var item BusinessTemplate
	if err := a.store.get(r.Context(), businessTemplateKind, id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	values := map[string]string{"template.name": item.Name, "template.version": item.Version, "template.kind": item.Kind}
	if isStratronixQuotePDFTemplate(item) {
		values = quoteTemplatePreviewValues(item.DefaultValues)
	} else if item.Kind == "document" {
		values = documentTemplatePreviewValues(item.DefaultValues, item.OutputFormat, item.DocumentType)
	}
	content, filename, contentType, err := a.renderBusinessTemplate(r.Context(), item, values, templateKindLabel(item.Kind)+"-模板预览")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	writeAttachment(w, filename, contentType, content)
}

func (a *businessAPI) renderBusinessTemplate(ctx context.Context, item BusinessTemplate, values map[string]string, title string) ([]byte, string, string, error) {
	if isPDFTemplateFile(item.FileName) {
		content, filename, err := a.renderPDFTemplateWithValues(ctx, item, values, title)
		if err != nil {
			return nil, "", "", err
		}
		return content, filename, "application/pdf", nil
	}
	if isWorkbookTemplateFile(item.FileName) {
		content, filename, err := a.renderWorkbookTemplateWithValues(ctx, item, values, title)
		if err != nil {
			return nil, "", "", err
		}
		return content, filename, workbookContentType(item.FileName), nil
	}
	content, filename, err := a.renderTemplateWithValues(ctx, item, values, title)
	if err != nil {
		return nil, "", "", err
	}
	return []byte(content), filename, "text/html; charset=utf-8", nil
}

func (a *businessAPI) defaultTemplate(ctx context.Context, kind string) (BusinessTemplate, error) {
	items, err := listRecords[BusinessTemplate](ctx, a.store, businessTemplateKind)
	if err != nil {
		return BusinessTemplate{}, err
	}
	var latest *BusinessTemplate
	for i := range items {
		item := items[i]
		if item.Kind != kind || strings.EqualFold(item.Status, "Archived") {
			continue
		}
		if item.Default {
			return item, nil
		}
		if latest == nil || item.Updated > latest.Updated {
			copy := item
			latest = &copy
		}
	}
	if latest != nil {
		return *latest, nil
	}
	return builtInTemplate(kind), nil
}

func builtInTemplate(kind string) BusinessTemplate {
	now := currentText()
	return BusinessTemplate{
		ID: "BUILTIN-" + strings.ToUpper(kind), Kind: kind, Name: templateKindLabel(kind) + "默认模板", Version: "v1",
		Source: "系统默认", Mode: "builtin", FileName: "", Status: "Active", Default: true, OutputFormat: "html",
		Placeholders: defaultTemplatePlaceholders(kind), FieldMapping: defaultTemplateFieldMapping(kind, defaultTemplatePlaceholders(kind)),
		CreatedAt: now, Updated: now,
	}
}

func defaultQuoteTemplateValues() map[string]string {
	return map[string]string{
		"supplier.company":      "STRATRONIX / 鼎图",
		"supplier.contact":      "Donald",
		"supplier.phone":        "86-755-23086689",
		"supplier.email":        "info@stratronix.ai",
		"supplier.website":      "www.stratronix.ai",
		"quote.priceTerms":      "FOB Shenzhen",
		"quote.paymentTerms":    "T/T 30% deposit, 70% before shipment",
		"quote.shipping":        "Shipping to be confirmed",
		"quote.destination":     "Rotterdam",
		"quote.leadTime":        "15-20 working days after deposit received",
		"quote.warranty":        "12 months under STRATRONIX standard terms",
		"quote.certification":   "CE / FCC / RoHS if applicable",
		"quote.packaging":       "Standard export carton",
		"quote.note":            "Quotation is subject to final quantity and specification confirmation.",
		"quote.validityDays":    "30",
		"quote.freightCurrency": "USD",
	}
}

func quoteTemplatePreviewValues(defaultValues map[string]string) map[string]string {
	values := mergeTemplateDefaultValues(defaultValues, defaultQuoteTemplateValues())
	for _, key := range quoteTemplateRealtimeFields() {
		values[key] = ""
	}
	values["pdf.quote.paymentTerms"] = compactPDFText(values["quote.paymentTerms"], 42)
	values["pdf.quote.priceTermsLine"] = "价格条件 / Price Terms: " + values["quote.priceTerms"]
	values["pdf.quote.paymentTermsLine"] = "付款条件 / Payment: " + values["quote.paymentTerms"]
	values["pdf.quote.leadTimeLine"] = "交期 / Lead Time: " + values["quote.leadTime"]
	values["pdf.quote.warrantyLine"] = "质保 / Warranty: " + values["quote.warranty"]
	values["pdf.quote.certificationLine"] = "认证 / Certification: " + values["quote.certification"]
	values["pdf.quote.packagingLine"] = "包装 / Packaging: " + values["quote.packaging"]
	values["pdf.quote.validityDetailLine"] = "报价有效期 / Validity: " + quoteValidityDetail(values["quote.validityDays"], "")
	values["pdf.quote.noteLine"] = "备注 / Note: " + values["quote.note"]
	return values
}

func quoteTemplateRealtimeFields() []string {
	return []string{
		"pdf.record.id",
		"record.id",
		"record.date",
		"quote.validity",
		"customer.name",
		"customer.country",
		"customer.contact",
		"pdf.customer.contact",
		"customer.email",
		"customer.address",
		"pdf.customer.address",
		"line.1.no",
		"line.1.model",
		"line.1.description",
		"line.1.quantity",
		"line.1.unitPrice",
		"line.1.amount",
		"pdf.line.1.unitPrice",
		"pdf.line.1.amount",
		"line.2.no",
		"line.2.model",
		"line.2.description",
		"line.2.quantity",
		"line.2.unitPrice",
		"line.2.amount",
		"pdf.line.2.unitPrice",
		"pdf.line.2.amount",
		"line.3.no",
		"line.3.model",
		"line.3.description",
		"line.3.quantity",
		"line.3.unitPrice",
		"line.3.amount",
		"pdf.line.3.unitPrice",
		"pdf.line.3.amount",
		"record.subtotal",
		"record.freight",
		"record.total",
		"pdf.record.subtotal",
		"pdf.record.freight",
		"pdf.record.total",
	}
}

func mergeTemplateDefaultValues(existing, defaults map[string]string) map[string]string {
	merged := make(map[string]string, len(defaults)+len(existing))
	for key, value := range defaults {
		merged[key] = value
	}
	for key, value := range existing {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		merged[key] = strings.TrimSpace(value)
	}
	return merged
}

func sanitizeTemplateDefaultValues(kind string, values map[string]string) map[string]string {
	defaults := defaultTemplateDefaultValues(kind)
	clean := make(map[string]string, len(defaults))
	for key := range defaults {
		if value, ok := values[key]; ok {
			clean[key] = strings.TrimSpace(value)
		}
	}
	return mergeTemplateDefaultValues(clean, defaults)
}

func defaultTemplateDefaultValues(kind string) map[string]string {
	if kind == "document" {
		return defaultDocumentTemplateValues()
	}
	return defaultQuoteTemplateValues()
}

func stringMapEqual(left, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}

func (a *businessAPI) renderTemplateWithValues(ctx context.Context, item BusinessTemplate, values map[string]string, title string) (string, string, error) {
	body := ""
	if isTextTemplateFile(item.FileName) {
		path, err := templateFilePath(item)
		if err != nil {
			return "", "", err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return "", "", err
		}
		body = string(data)
	}
	if strings.TrimSpace(body) == "" {
		body = builtInTemplateHTML(item.Kind)
	}
	for key, value := range values {
		replacement := html.EscapeString(value)
		if isTrustedTemplateHTMLValue(key) {
			replacement = value
		}
		body = strings.ReplaceAll(body, "{{"+key+"}}", replacement)
		body = templatePlaceholderPattern.ReplaceAllStringFunc(body, func(match string) string {
			parts := templatePlaceholderPattern.FindStringSubmatch(match)
			if len(parts) == 2 && parts[1] == key {
				return replacement
			}
			return match
		})
	}
	body = templatePlaceholderPattern.ReplaceAllString(body, `<span style="color:#9a3412">待补充:$1</span>`)
	if !strings.Contains(strings.ToLower(body), "<html") {
		body = "<!doctype html><html><head><meta charset=\"utf-8\"><title>" + html.EscapeString(title) + "</title>" + templatePrintStyle() + "</head><body>" + body + "</body></html>"
	}
	filename := safeDownloadName(title)
	return body, filename + ".html", nil
}

func (a *businessAPI) renderPDFTemplateWithValues(ctx context.Context, item BusinessTemplate, values map[string]string, title string) ([]byte, string, error) {
	path, err := templateFilePath(item)
	if err != nil {
		return nil, "", err
	}
	if err := api.ValidateFile(path, nil); err != nil {
		return nil, "", fmt.Errorf("PDF 模板无效或无法读取: %w", err)
	}
	out, err := os.CreateTemp("", "sta100-pdf-template-*.pdf")
	if err != nil {
		return nil, "", err
	}
	outPath := out.Name()
	_ = out.Close()
	defer os.Remove(outPath)

	if hasPDFFormFieldMappings(item) {
		if err := fillPDFFormTemplate(path, outPath, item, values); err != nil {
			return nil, "", err
		}
	} else {
		fields := normalizePDFTemplateFields(item.PDFFields)
		if len(fields) == 0 {
			return nil, "", errors.New("PDF 模板缺少字段坐标，请在模板管理中配置 pdfFields 后再生成")
		}
		if isStratronixQuotePDFTemplate(item) {
			cleanPath, cleanErr := ensureStratronixQuoteCleanTemplateFile()
			if cleanErr != nil {
				return nil, "", cleanErr
			}
			path = cleanPath
		}
		if err := stampPDFTemplate(path, outPath, fields, values); err != nil {
			return nil, "", err
		}
	}
	content, err := os.ReadFile(outPath)
	if err != nil {
		return nil, "", err
	}
	return content, safeDownloadName(title) + ".pdf", nil
}

func hasPDFFormFieldMappings(item BusinessTemplate) bool {
	for templateField := range item.FieldMapping {
		for _, formField := range item.PDFFormFields {
			if strings.EqualFold(strings.TrimSpace(templateField), strings.TrimSpace(formField)) {
				return true
			}
		}
	}
	return false
}

func fillPDFFormTemplate(inPath, outPath string, item BusinessTemplate, values map[string]string) error {
	fields := make([]map[string]any, 0)
	for formField, valueKey := range item.FieldMapping {
		formField = strings.TrimSpace(formField)
		if formField == "" {
			continue
		}
		if !pdfFormFieldExists(item.PDFFormFields, formField) {
			continue
		}
		value := strings.TrimSpace(values[valueKey])
		if value == "" {
			value = strings.TrimSpace(values[formField])
		}
		fields = append(fields, map[string]any{"name": formField, "value": value})
	}
	if len(fields) == 0 {
		return errors.New("PDF 表单模板没有可写入的字段映射，请检查 fieldMapping 是否使用 PDF 表单字段名")
	}
	payload := map[string]any{
		"header": map[string]string{"source": filepath.Base(inPath), "producer": "STA-100"},
		"forms":  []map[string]any{{"textfield": fields}},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	input, err := os.CreateTemp("", "sta100-pdf-form-*.json")
	if err != nil {
		return err
	}
	inputPath := input.Name()
	if _, err := input.Write(data); err != nil {
		_ = input.Close()
		_ = os.Remove(inputPath)
		return err
	}
	if err := input.Close(); err != nil {
		_ = os.Remove(inputPath)
		return err
	}
	defer os.Remove(inputPath)
	if err := api.FillFormFile(inPath, inputPath, outPath, nil); err != nil {
		return fmt.Errorf("PDF 表单字段填充失败: %w", err)
	}
	return nil
}

func stampPDFTemplate(inPath, outPath string, fields []PDFTemplateField, values map[string]string) error {
	if err := ensurePDFUserFonts(); err != nil {
		return fmt.Errorf("PDF 字体初始化失败: %w", err)
	}
	ctx, err := api.ReadContextFile(inPath)
	if err != nil {
		return fmt.Errorf("读取 PDF 模板失败: %w", err)
	}
	stamps := map[int][]*model.Watermark{}
	for _, field := range fields {
		key := strings.TrimSpace(field.Key)
		if key == "" {
			continue
		}
		value := strings.TrimSpace(values[key])
		if value == "" {
			continue
		}
		value = normalizePDFTemplateValue(key, value)
		page := field.Page
		if page <= 0 {
			page = 1
		}
		if page > ctx.PageCount {
			continue
		}
		fontSize := field.FontSize
		if fontSize <= 0 {
			fontSize = 10
		}
		align := normalizePDFAlign(field.Align)
		color := strings.TrimSpace(field.Color)
		if color == "" {
			color = "#172033"
		}
		desc := fmt.Sprintf("font:%s, points:%s, scale:1 abs, pos:bl, off:%s %s, align:%s, fillcol:%s, rot:0",
			pdfTemplateFontName(), formatPDFFontSize(fontSize), formatPDFNumber(field.X), formatPDFNumber(field.Y), align, color)
		if field.MaxWidth > 0 {
			desc += ", maxwidth:" + formatPDFNumber(field.MaxWidth)
		}
		if strings.TrimSpace(field.BackgroundColor) != "" {
			desc += ", bgcolor:" + strings.TrimSpace(field.BackgroundColor) + ", margins:" + pdfTemplateFieldMargins(key)
		}
		wm, err := api.TextWatermark(value, desc, true, false, types.POINTS)
		if err != nil {
			return fmt.Errorf("PDF 字段 %s 坐标配置无效: %w", key, err)
		}
		stamps[page] = append(stamps[page], wm)
	}
	if len(stamps) == 0 {
		return errors.New("PDF 模板没有可写入的字段值，请检查报价数据和 pdfFields 字段配置")
	}
	if err := api.AddWatermarksSliceMapFile(inPath, outPath, stamps, nil); err != nil {
		return fmt.Errorf("PDF 模板写入字段失败: %w", err)
	}
	return nil
}

func normalizePDFTemplateValue(key, value string) string {
	if key == "record.lines" {
		value = strings.NewReplacer("<br>", "\n", "<br/>", "\n", "<br />", "\n").Replace(value)
		value = regexp.MustCompile(`(?i)</tr>`).ReplaceAllString(value, "\n")
		value = regexp.MustCompile(`(?i)</t[dh]>`).ReplaceAllString(value, "  ")
		value = regexp.MustCompile(`(?s)<[^>]*>`).ReplaceAllString(value, "")
		value = html.UnescapeString(value)
		value = regexp.MustCompile(`[ \t]+`).ReplaceAllString(value, " ")
		lines := strings.Split(value, "\n")
		cleaned := make([]string, 0, len(lines))
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				cleaned = append(cleaned, line)
			}
		}
		return strings.Join(cleaned, "\n")
	}
	return value
}

func ensurePDFUserFonts() error {
	pdfFontInstallMu.Lock()
	defer pdfFontInstallMu.Unlock()

	configRoot := strings.TrimSpace(os.Getenv("STA100_PDFCPU_CONFIG_DIR"))
	if configRoot == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		configRoot = filepath.Join(home, ".local", "share", "sta100", "pdfcpu")
	}
	if err := api.EnsureDefaultConfigAt(configRoot); err != nil {
		return err
	}
	fontName := pdfTemplateFontName()
	if _, err := os.Stat(filepath.Join(configRoot, "pdfcpu", "fonts", fontName+".gob")); err == nil {
		return nil
	}
	candidates := []string{
		"/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
	}
	fonts := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			fonts = append(fonts, candidate)
		}
	}
	if len(fonts) == 0 {
		return errors.New("未找到可用字体文件")
	}
	return api.InstallFonts(fonts)
}

func pdfTemplateFontName() string {
	if _, err := os.Stat("/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc"); err == nil {
		return "WenQuanYiZenHei"
	}
	return "DejaVuSans"
}

func inspectPDFFormFieldNames(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	fields, err := api.FormFields(file, nil)
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(fields))
	seen := map[string]bool{}
	for _, field := range fields {
		name := strings.TrimSpace(firstNonEmpty(field.Name, field.ID))
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func pdfFormFieldExists(fields []string, name string) bool {
	for _, field := range fields {
		if strings.EqualFold(strings.TrimSpace(field), strings.TrimSpace(name)) {
			return true
		}
	}
	return false
}

func normalizePDFFormFields(fields []string) []string {
	out := make([]string, 0, len(fields))
	seen := map[string]bool{}
	for _, field := range fields {
		field = strings.TrimSpace(field)
		key := strings.ToLower(field)
		if field == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, field)
	}
	sort.Strings(out)
	return out
}

func templateFilePath(item BusinessTemplate) (string, error) {
	if strings.TrimSpace(item.Path) != "" {
		return item.Path, nil
	}
	storageName := strings.TrimSpace(item.StorageFileName)
	if storageName == "" && item.ID != "" && item.FileName != "" {
		storageName = item.ID + strings.ToLower(filepath.Ext(item.FileName))
	}
	if storageName == "" {
		return "", errors.New("模板文件路径缺失")
	}
	root, err := templateStorageRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, filepath.Base(storageName)), nil
}

func ensureStratronixQuoteCleanTemplateFile() (string, error) {
	data, err := readEmbeddedUIFile(ui, builtInQuotePDFTemplateAsset)
	if err != nil {
		return "", err
	}
	root, err := templateStorageRoot()
	if err != nil {
		return "", err
	}
	target := filepath.Join(root, builtInQuotePDFTemplateID+".pdf")
	if existing, err := os.ReadFile(target); err != nil || !bytes.Equal(existing, data) {
		if err := os.WriteFile(target, data, 0o600); err != nil {
			return "", err
		}
	}
	return target, nil
}

func isTrustedTemplateHTMLValue(key string) bool {
	switch key {
	case "record.lines":
		return true
	default:
		return false
	}
}

func writeAttachment(w http.ResponseWriter, filename, contentType string, content []byte) {
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(filename, `"`, "")+`"`)
	_, _ = w.Write(content)
}

func templateStorageRoot() (string, error) {
	root := strings.TrimSpace(os.Getenv("STA100_TEMPLATE_DIR"))
	if root == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		root = filepath.Join(home, ".local", "share", "sta100", "templates")
	}
	if err := os.MkdirAll(root, 0o700); err != nil {
		return "", err
	}
	return root, nil
}

func normalizeTemplateKind(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "quote", "quotes":
		return "quote"
	case "order", "orders":
		return "order"
	case "document", "documents", "doc":
		return "document"
	default:
		return ""
	}
}

func templateKindLabel(kind string) string {
	switch kind {
	case "quote":
		return "报价单"
	case "order":
		return "订单"
	case "document":
		return "单据"
	default:
		return "业务"
	}
}

func allowedTemplateExtension(ext string, imageMode bool) bool {
	if imageMode {
		return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
	}
	switch ext {
	case ".docx", ".xlsx", ".xlsm", ".pdf", ".html", ".htm", ".txt", ".md":
		return true
	default:
		return false
	}
}

func isTextTemplateFile(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".html", ".htm", ".txt", ".md":
		return true
	default:
		return false
	}
}

func templatePlaceholders(ext string, data []byte) []string {
	if ext != ".html" && ext != ".htm" && ext != ".txt" && ext != ".md" {
		return defaultTemplatePlaceholders("")
	}
	matches := templatePlaceholderPattern.FindAllStringSubmatch(string(data), -1)
	seen, out := map[string]bool{}, []string{}
	for _, match := range matches {
		if len(match) == 2 && !seen[match[1]] {
			seen[match[1]] = true
			out = append(out, match[1])
		}
	}
	if len(out) == 0 {
		return defaultTemplatePlaceholders("")
	}
	sort.Strings(out)
	return out
}

func pdfTemplateFieldMargins(key string) string {
	switch key {
	case "pdf.record.id":
		return "2 3"
	case "pdf.record.subtotal", "pdf.record.freight", "pdf.record.total", "pdf.line.1.unitPrice", "pdf.line.1.amount", "pdf.line.2.unitPrice", "pdf.line.2.amount", "pdf.line.3.unitPrice", "pdf.line.3.amount":
		return "2 3"
	case "customer.name", "customer.country", "pdf.customer.contact", "customer.email", "pdf.customer.address":
		return "2 3"
	case "quote.validity", "record.date", "pdf.quote.paymentTerms":
		return "2 3"
	default:
		return "2 3"
	}
}

func defaultTemplatePlaceholders(kind string) []string {
	common := []string{"company.name", "customer.name", "record.id", "record.date", "record.total", "record.lines"}
	switch kind {
	case "quote":
		return append(common,
			"quote.subject", "quote.valid", "quote.terms", "quote.currency", "quote.validity",
			"quote.paymentTerms", "quote.priceTerms", "quote.shipping", "quote.destination",
			"quote.leadTime", "quote.warranty", "quote.certification", "quote.packaging", "quote.note",
			"customer.country", "customer.city", "customer.contact", "customer.phone", "customer.email", "customer.website", "customer.address",
			"record.subtotal", "record.freight", "record.tax",
			"line.1.no", "line.1.model", "line.1.description", "line.1.quantity", "line.1.unitPrice", "line.1.amount",
			"line.2.no", "line.2.model", "line.2.description", "line.2.quantity", "line.2.unitPrice", "line.2.amount",
			"line.3.no", "line.3.model", "line.3.description", "line.3.quantity", "line.3.unitPrice", "line.3.amount",
			"pdf.record.id", "pdf.quote.paymentTerms", "pdf.record.subtotal", "pdf.record.freight", "pdf.record.total",
			"pdf.line.1.unitPrice", "pdf.line.1.amount", "pdf.line.2.unitPrice", "pdf.line.2.amount", "pdf.line.3.unitPrice", "pdf.line.3.amount",
			"pdf.customer.contact", "pdf.customer.address",
		)
	case "order":
		return append(common, "order.po", "order.delivery", "order.status")
	case "document":
		return append(common, "document.type", "document.language", "order.id")
	default:
		return common
	}
}

func defaultTemplateFieldMapping(kind string, placeholders []string) map[string]string {
	mapping := map[string]string{}
	for _, key := range placeholders {
		mapping[key] = key
	}
	if len(mapping) == 0 {
		for _, key := range defaultTemplatePlaceholders(kind) {
			mapping[key] = key
		}
	}
	return mapping
}

func templateRecognitionNote(imageMode bool, ext string) string {
	if imageMode {
		return "已保存图片模板；OCR/版式识别待接入 OpenClaw 或客户指定 OCR 引擎，当前可人工维护字段映射后发布。"
	}
	if ext == ".docx" || ext == ".xlsx" || ext == ".pdf" {
		return "已保存原始模板文件；PDF 可按表单字段或坐标字段生成原版式文件，DOCX/XLSX 当前用于模板归档和默认模板选择。"
	}
	return "已识别文本模板中的 {{placeholder}} 占位符，可直接按字段映射生成 HTML 文件。"
}

func defaultTemplateOutputFormat(ext string) string {
	if ext == ".pdf" {
		return "pdf"
	}
	if ext == ".xlsm" || ext == ".xlsx" {
		return strings.TrimPrefix(ext, ".")
	}
	return "html"
}

func isPDFTemplateFile(name string) bool {
	return strings.EqualFold(filepath.Ext(name), ".pdf")
}

func isWorkbookTemplateFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".xlsm" || ext == ".xlsx"
}

func workbookContentType(name string) string {
	if strings.EqualFold(filepath.Ext(name), ".xlsm") {
		return "application/vnd.ms-excel.sheet.macroEnabled.12"
	}
	return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
}

func normalizePDFAlign(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "left", "l":
		return "l"
	case "right", "r":
		return "r"
	case "center", "c":
		return "c"
	default:
		return "l"
	}
}

func formatPDFNumber(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatPDFFontSize(value float64) string {
	if value <= 0 {
		value = 10
	}
	return strconv.Itoa(int(value + 0.5))
}

func normalizePDFTemplateFields(fields []PDFTemplateField) []PDFTemplateField {
	out := make([]PDFTemplateField, 0, len(fields))
	for _, field := range fields {
		field.Key = strings.TrimSpace(field.Key)
		if field.Key == "" {
			continue
		}
		if field.FontSize <= 0 {
			field.FontSize = 10
		}
		if field.Page <= 0 {
			field.Page = 1
		}
		out = append(out, field)
	}
	return out
}

func defaultPDFTemplateFields(kind string) []PDFTemplateField {
	switch kind {
	case "quote":
		return stratronixQuotePDFTemplateFields()
	case "order":
		return append(genericPDFTemplateFields(),
			PDFTemplateField{Key: "order.po", Label: "PO 号", Page: 1, X: 42, Y: 744, FontSize: 11, Align: "l", BackgroundColor: "#ffffff"},
			PDFTemplateField{Key: "order.delivery", Label: "交付日期", Page: 1, X: 42, Y: 720, FontSize: 11, Align: "l", BackgroundColor: "#ffffff"},
			PDFTemplateField{Key: "order.status", Label: "状态", Page: 1, X: 42, Y: 696, FontSize: 10, Align: "l", BackgroundColor: "#ffffff"},
		)
	case "document":
		return append(genericPDFTemplateFields(),
			PDFTemplateField{Key: "document.type", Label: "单据类型", Page: 1, X: 42, Y: 744, FontSize: 11, Align: "l", BackgroundColor: "#ffffff"},
			PDFTemplateField{Key: "document.language", Label: "语言", Page: 1, X: 42, Y: 720, FontSize: 11, Align: "l", BackgroundColor: "#ffffff"},
			PDFTemplateField{Key: "order.id", Label: "订单号", Page: 1, X: 42, Y: 696, FontSize: 10, Align: "l", BackgroundColor: "#ffffff"},
		)
	case "pi":
		return stratronixPIPDFTemplateFields()
	}
	return genericPDFTemplateFields()
}

func stratronixPIPDFTemplateFields() []PDFTemplateField {
	return []PDFTemplateField{
		{Key: "supplier.company", Label: "供方公司", Page: 1, X: 44, Y: 697, FontSize: 9, MaxWidth: 210, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "supplier.address", Label: "供方地址", Page: 1, X: 44, Y: 677, FontSize: 8, MaxWidth: 200, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "supplier.country", Label: "供方国家", Page: 1, X: 44, Y: 663, FontSize: 8, MaxWidth: 160, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "supplier.phone", Label: "供方电话", Page: 1, X: 44, Y: 650, FontSize: 8, MaxWidth: 150, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "supplier.email", Label: "供方邮箱", Page: 1, X: 44, Y: 637, FontSize: 8, MaxWidth: 160, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "supplier.website", Label: "供方网址", Page: 1, X: 44, Y: 624, FontSize: 8, MaxWidth: 160, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "pdf.record.id", Label: "PI 编号", Page: 1, X: 459, Y: 736, FontSize: 9, MaxWidth: 100, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "pdf.record.balanceDue", Label: "待付金额", Page: 1, X: 493, Y: 692, FontSize: 11, MaxWidth: 80, Align: "r", BackgroundColor: "#ffffff"},
		{Key: "customer.name", Label: "客户公司", Page: 1, X: 45, Y: 579, FontSize: 8, MaxWidth: 210, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "pdf.customer.address", Label: "客户地址", Page: 1, X: 45, Y: 558, FontSize: 7, MaxWidth: 260, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "customer.country", Label: "客户国家", Page: 1, X: 45, Y: 542, FontSize: 8, MaxWidth: 150, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "pdf.record.date", Label: "PI 日期", Page: 1, X: 512, Y: 586, FontSize: 8, MaxWidth: 70, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "payment.terms", Label: "付款条款", Page: 1, X: 497, Y: 558, FontSize: 8, MaxWidth: 90, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "pdf.record.dueDate", Label: "到期日", Page: 1, X: 512, Y: 532, FontSize: 8, MaxWidth: 70, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "pdf.pi.lines", Label: "产品明细", Page: 1, X: 56, Y: 477, FontSize: 7, MaxWidth: 475, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "pdf.record.total", Label: "合计", Page: 1, X: 526, Y: 292, FontSize: 9, MaxWidth: 80, Align: "r", BackgroundColor: "#ffffff"},
		{Key: "amount.words", Label: "英文大写金额", Page: 1, X: 430, Y: 263, FontSize: 8, MaxWidth: 145, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "document.notes", Label: "备注", Page: 1, X: 45, Y: 176, FontSize: 7, MaxWidth: 480, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "bank.account", Label: "银行账户", Page: 1, X: 45, Y: 103, FontSize: 7, MaxWidth: 480, Align: "l", BackgroundColor: "#ffffff"},
	}
}

func genericPDFTemplateFields() []PDFTemplateField {
	return []PDFTemplateField{
		{Key: "company.name", Label: "公司名称", Page: 1, X: 42, Y: 800, FontSize: 12, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "customer.name", Label: "客户名称", Page: 1, X: 42, Y: 772, FontSize: 11, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "record.id", Label: "单号", Page: 1, X: 405, Y: 800, FontSize: 11, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "record.date", Label: "日期", Page: 1, X: 405, Y: 778, FontSize: 11, Align: "l", BackgroundColor: "#ffffff"},
		{Key: "record.total", Label: "总金额", Page: 1, X: 405, Y: 754, FontSize: 12, Align: "l", BackgroundColor: "#ffffff"},
	}
}

func stratronixQuotePDFTemplateFields() []PDFTemplateField {
	return []PDFTemplateField{
		{Key: "pdf.record.id", Label: "报价单号", Page: 1, X: 171, Y: 607, FontSize: 8, MaxWidth: 100, Align: "l"},
		{Key: "record.date", Label: "报价日期", Page: 1, X: 388, Y: 610, FontSize: 9, MaxWidth: 80, Align: "l"},
		{Key: "quote.validity", Label: "有效期", Page: 1, X: 171, Y: 574, FontSize: 8, MaxWidth: 100, Align: "l"},
		{Key: "pdf.quote.paymentTerms", Label: "付款条件", Page: 1, X: 388, Y: 574, FontSize: 8, MaxWidth: 150, Align: "l"},
		{Key: "customer.name", Label: "客户公司", Page: 1, X: 128, Y: 510, FontSize: 8, MaxWidth: 170, Align: "l"},
		{Key: "customer.country", Label: "客户国家", Page: 1, X: 128, Y: 495, FontSize: 8, MaxWidth: 170, Align: "l"},
		{Key: "pdf.customer.contact", Label: "客户联系人", Page: 1, X: 128, Y: 480, FontSize: 8, MaxWidth: 170, Align: "l"},
		{Key: "customer.email", Label: "客户邮箱", Page: 1, X: 128, Y: 465, FontSize: 8, MaxWidth: 170, Align: "l"},
		{Key: "pdf.customer.address", Label: "客户地址", Page: 1, X: 128, Y: 450, FontSize: 8, MaxWidth: 145, Align: "l"},
		{Key: "supplier.company", Label: "供方公司", Page: 1, X: 372, Y: 510, FontSize: 8, MaxWidth: 150, Align: "l"},
		{Key: "supplier.contact", Label: "供方联系人", Page: 1, X: 372, Y: 495, FontSize: 8, MaxWidth: 150, Align: "l"},
		{Key: "supplier.phone", Label: "供方电话", Page: 1, X: 372, Y: 480, FontSize: 8, MaxWidth: 150, Align: "l"},
		{Key: "supplier.email", Label: "供方邮箱", Page: 1, X: 372, Y: 465, FontSize: 8, MaxWidth: 150, Align: "l"},
		{Key: "supplier.website", Label: "供方网址", Page: 1, X: 372, Y: 450, FontSize: 8, MaxWidth: 150, Align: "l"},
		{Key: "line.1.no", Label: "序号", Page: 1, X: 67, Y: 368, FontSize: 7, MaxWidth: 20, Align: "c"},
		{Key: "line.1.model", Label: "型号", Page: 1, X: 91, Y: 368, FontSize: 7, MaxWidth: 52, Align: "l"},
		{Key: "line.1.description", Label: "描述", Page: 1, X: 159, Y: 368, FontSize: 7, MaxWidth: 170, Align: "l"},
		{Key: "line.1.quantity", Label: "数量", Page: 1, X: 349, Y: 368, FontSize: 7, MaxWidth: 36, Align: "l"},
		{Key: "pdf.line.1.unitPrice", Label: "单价", Page: 1, X: 430, Y: 368, FontSize: 7, MaxWidth: 52, Align: "r"},
		{Key: "pdf.line.1.amount", Label: "小计", Page: 1, X: 514, Y: 368, FontSize: 7, MaxWidth: 52, Align: "r"},
		{Key: "line.2.no", Label: "序号2", Page: 1, X: 67, Y: 356, FontSize: 7, MaxWidth: 20, Align: "c"},
		{Key: "line.2.model", Label: "型号2", Page: 1, X: 91, Y: 356, FontSize: 7, MaxWidth: 52, Align: "l"},
		{Key: "line.2.description", Label: "描述2", Page: 1, X: 159, Y: 356, FontSize: 7, MaxWidth: 170, Align: "l"},
		{Key: "line.2.quantity", Label: "数量2", Page: 1, X: 349, Y: 356, FontSize: 7, MaxWidth: 36, Align: "l"},
		{Key: "pdf.line.2.unitPrice", Label: "单价2", Page: 1, X: 430, Y: 356, FontSize: 7, MaxWidth: 52, Align: "r"},
		{Key: "pdf.line.2.amount", Label: "小计2", Page: 1, X: 514, Y: 356, FontSize: 7, MaxWidth: 52, Align: "r"},
		{Key: "line.3.no", Label: "序号3", Page: 1, X: 67, Y: 344, FontSize: 7, MaxWidth: 20, Align: "c"},
		{Key: "line.3.model", Label: "型号3", Page: 1, X: 91, Y: 344, FontSize: 7, MaxWidth: 52, Align: "l"},
		{Key: "line.3.description", Label: "描述3", Page: 1, X: 159, Y: 344, FontSize: 7, MaxWidth: 170, Align: "l"},
		{Key: "line.3.quantity", Label: "数量3", Page: 1, X: 349, Y: 344, FontSize: 7, MaxWidth: 36, Align: "l"},
		{Key: "pdf.line.3.unitPrice", Label: "单价3", Page: 1, X: 430, Y: 344, FontSize: 7, MaxWidth: 52, Align: "r"},
		{Key: "pdf.line.3.amount", Label: "小计3", Page: 1, X: 514, Y: 344, FontSize: 7, MaxWidth: 52, Align: "r"},
		{Key: "pdf.record.subtotal", Label: "产品小计", Page: 1, X: 526, Y: 322, FontSize: 8, MaxWidth: 82, Align: "r"},
		{Key: "pdf.record.freight", Label: "运费", Page: 1, X: 526, Y: 294, FontSize: 8, MaxWidth: 82, Align: "r"},
		{Key: "pdf.record.total", Label: "总计", Page: 1, X: 526, Y: 242, FontSize: 10, MaxWidth: 82, Align: "r"},
		{Key: "pdf.quote.priceTermsLine", Label: "价格条件", Page: 1, X: 61, Y: 200, FontSize: 7, MaxWidth: 440, Align: "l"},
		{Key: "pdf.quote.paymentTermsLine", Label: "付款条件条款", Page: 1, X: 61, Y: 187, FontSize: 7, MaxWidth: 440, Align: "l"},
		{Key: "pdf.quote.leadTimeLine", Label: "交期", Page: 1, X: 61, Y: 173, FontSize: 7, MaxWidth: 440, Align: "l"},
		{Key: "pdf.quote.warrantyLine", Label: "质保", Page: 1, X: 61, Y: 160, FontSize: 7, MaxWidth: 440, Align: "l"},
		{Key: "pdf.quote.certificationLine", Label: "认证", Page: 1, X: 61, Y: 147, FontSize: 7, MaxWidth: 440, Align: "l"},
		{Key: "pdf.quote.packagingLine", Label: "包装", Page: 1, X: 61, Y: 133, FontSize: 7, MaxWidth: 440, Align: "l"},
		{Key: "pdf.quote.validityDetailLine", Label: "报价有效期条款", Page: 1, X: 61, Y: 120, FontSize: 7, MaxWidth: 440, Align: "l"},
		{Key: "pdf.quote.noteLine", Label: "备注", Page: 1, X: 61, Y: 107, FontSize: 7, MaxWidth: 440, Align: "l"},
	}
}

func validateBusinessTemplate(item BusinessTemplate) error {
	if normalizeTemplateKind(item.Kind) == "" {
		return errors.New("模板类型必须是 quote、order 或 document")
	}
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("模板名称不能为空")
	}
	if item.Status == "" {
		item.Status = "Draft"
	}
	return nil
}

func safeDownloadName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		value = "STA100-template"
	}
	replacer := strings.NewReplacer("/", "-", "\\", "-", ":", "-", "*", "-", "?", "", "\"", "", "<", "", ">", "", "|", "-")
	return replacer.Replace(value)
}

func templatePrintStyle() string {
	return `<style>body{font-family:Arial,"Microsoft YaHei",sans-serif;margin:32px;color:#172033}table{width:100%;border-collapse:collapse;margin:18px 0}th,td{border:1px solid #cbd5e1;padding:8px;text-align:left}th{background:#eef2f7}.muted{color:#64748b}.header{display:flex;justify-content:space-between;border-bottom:2px solid #203047;padding-bottom:14px;margin-bottom:22px}.total{text-align:right;font-size:18px;font-weight:700}.section{margin-top:18px}</style>`
}

func builtInTemplateHTML(kind string) string {
	label := templateKindLabel(kind)
	return `<!doctype html><html><head><meta charset="utf-8"><title>` + label + `</title>` + templatePrintStyle() + `</head><body><div class="header"><h1>` + label + `</h1><div><strong>{{company.name}}</strong><br><span class="muted">{{record.date}}</span></div></div><section class="section"><h2>{{customer.name}}</h2><p>编号：{{record.id}}</p></section><section class="section">{{record.lines}}</section><p class="total">合计：{{record.total}}</p></body></html>`
}

func businessLinesHTML(lines []BusinessLine) string {
	var builder strings.Builder
	builder.WriteString("<table><thead><tr><th>产品</th><th>数量</th><th>单价</th><th>金额</th></tr></thead><tbody>")
	for _, line := range lines {
		builder.WriteString("<tr><td>")
		builder.WriteString(html.EscapeString(firstNonEmpty(line.ProductName, line.ProductID)))
		builder.WriteString("</td><td>")
		builder.WriteString(fmt.Sprintf("%.2f", line.Quantity))
		builder.WriteString("</td><td>")
		builder.WriteString(fmt.Sprintf("%.2f", line.UnitPrice))
		builder.WriteString("</td><td>")
		builder.WriteString(fmt.Sprintf("%.2f", line.Amount))
		builder.WriteString("</td></tr>")
	}
	builder.WriteString("</tbody></table>")
	return builder.String()
}
