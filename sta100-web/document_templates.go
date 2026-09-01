package main

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

var errOfficeConverterMissing = fmt.Errorf("模板需要转换，但当前环境未安装 LibreOffice/soffice")

func defaultDocumentTemplateValues() map[string]string {
	return map[string]string{
		"supplier.company":          "Stratronix Technology (Shenzhen) Company, Limited",
		"supplier.address":          "Shenzhen China 528201",
		"supplier.country":          "China",
		"supplier.phone":            "86-755-23086689",
		"supplier.email":            "info@stratronix.ai",
		"supplier.website":          "www.stratronix.ai",
		"domestic.shipper":          "Stratronix Technology (Shenzhen) Company, Limited",
		"domestic.shipper.address":  "Shenzhen China 528201",
		"customs.code":              "",
		"exit.customs":              "皇岗海关",
		"tax.exemption.nature":      "",
		"trade.mode":                "一般贸易",
		"contract.sign.place":       "Shenzhen",
		"port.departure":            "Shenzhen",
		"trade.term":                "FOB",
		"shipping.period":           "To be confirmed",
		"container.no":              "",
		"settlement.method":         "T/T",
		"destination.port":          "",
		"exemption.mode":            "",
		"origin.country":            "中国",
		"transport.mode":            "航空",
		"domestic.source.location":  "深圳",
		"freight":                   "",
		"packaging.type":            "纸箱",
		"package.count":             "1",
		"gross.weight":              "",
		"net.weight":                "",
		"quantity.unit":             "台",
		"product.spec":              "",
		"document.notes":            "1. STA-100, HDMI CABLE, MANUAL, GIFTBOX in total one box;\n2. Product details and customs declaration elements require final manual review.",
		"declaration.element":       "品牌类型、出口享惠情况、用途、材质、规格型号等申报要素待人工复核。",
		"payment.terms":             "Due on Receipt",
		"bank.account":              "BANK ACCOUNT:",
		"amount.words":              "",
		"document.missingTemplates": "订单标准模板暂未提供；单据页中的 PI、CI、PL、合同、报关单分别作为独立业务模板展示，底层共用客户提供的标准底板。",
	}
}

func documentTemplatePreviewValues(defaultValues map[string]string, outputFormat, documentType string) map[string]string {
	values := mergeTemplateDefaultValues(defaultValues, defaultDocumentTemplateValues())
	values["record.id"] = "PI-20260826-001"
	values["pdf.record.id"] = "PI-20260826-001"
	values["record.date"] = time.Now().Format("2006-01-02")
	values["pdf.record.date"] = formatPDFDate(values["record.date"])
	values["record.dueDate"] = values["record.date"]
	values["pdf.record.dueDate"] = formatPDFDate(values["record.dueDate"])
	values["record.total"] = "EUR 1,900.00"
	values["pdf.record.total"] = "1900.00"
	values["pdf.record.balanceDue"] = "0.00"
	values["customer.name"] = "Sample Customer GmbH"
	values["customer.country"] = "Germany"
	values["customer.city"] = "Berlin"
	values["customer.address"] = "Sample Street 1, Berlin"
	values["pdf.customer.address"] = "Sample Street 1, Berlin"
	values["order.id"] = "SO-2026-0001"
	values["order.po"] = "PO-SAMPLE-001"
	values["document.type"] = firstNonEmpty(documentType, "PI")
	values["document.language"] = "英文"
	values["invoice.no"] = "PI-20260826-001"
	values["contract.no"] = "SO-2026-0001"
	values["contract.sign.date"] = values["record.date"]
	values["currency"] = "EUR"
	values["line.1.no"] = "1"
	values["line.1.hs"] = "8471504090"
	values["line.1.model"] = "STA-100"
	values["line.1.description"] = "Smart cycling AI appliance"
	values["line.1.quantity"] = "10"
	values["line.1.unit"] = values["quantity.unit"]
	values["line.1.unitPrice"] = "190.00"
	values["line.1.amount"] = "1900.00"
	values["pdf.pi.lines"] = "1  STA-100\n   Smart cycling AI appliance\n   HS CODE 8471504090\n   Qty 10   Rate 190.00   Amount 1900.00"
	values["record.lines"] = businessLinesHTML([]BusinessLine{{ProductID: "STA-100", ProductName: "Smart cycling AI appliance", Quantity: 10, UnitPrice: 190, Amount: 1900}})
	if outputFormat == "xlsm" {
		values["document.type"] = firstNonEmpty(documentType, "CI")
	}
	return values
}

func (a *businessAPI) documentTemplateValues(ctx context.Context, item Document, tpl BusinessTemplate) map[string]string {
	values := mergeTemplateDefaultValues(tpl.DefaultValues, defaultDocumentTemplateValues())
	values["company.name"] = "STRATRONIX"
	values["customer.name"] = item.Customer
	values["record.id"] = item.ID
	values["pdf.record.id"] = item.ID
	values["record.date"] = time.Now().Format("2006-01-02")
	values["pdf.record.date"] = formatPDFDate(values["record.date"])
	values["record.dueDate"] = firstNonEmpty(values["record.dueDate"], values["record.date"])
	values["pdf.record.dueDate"] = formatPDFDate(values["record.dueDate"])
	values["record.total"] = item.Value
	values["record.lines"] = businessLinesHTML(item.Lines)
	values["document.type"] = item.Type
	values["document.language"] = item.Language
	values["order.id"] = item.Order
	values["invoice.no"] = item.ID
	values["contract.no"] = item.Order
	values["contract.sign.date"] = values["record.date"]
	values["currency"] = currencyFromMoney(item.Value)
	values["pdf.record.total"] = formatPDFMoney(moneyNumber(item.Value))
	values["pdf.record.balanceDue"] = "0.00"
	values["amount.words"] = firstNonEmpty(values["amount.words"], item.Value)

	if customer, ok := a.customerByName(ctx, item.Customer); ok {
		values["customer.country"] = firstNonEmpty(customer.Country, values["customer.country"])
		values["customer.city"] = firstNonEmpty(customer.City, values["customer.city"])
		values["customer.contact"] = firstNonEmpty(customer.Contact, values["customer.contact"])
		values["customer.phone"] = firstNonEmpty(customer.Phone, values["customer.phone"])
		values["customer.email"] = firstNonEmpty(customer.Email, values["customer.email"])
		values["customer.website"] = firstNonEmpty(customer.Website, values["customer.website"])
		values["customer.address"] = firstNonEmpty(customer.Description, customer.City, values["customer.address"])
	}
	values["pdf.customer.address"] = compactPDFText(values["customer.address"], 86)

	var order Order
	if item.Order != "" && a.store.get(ctx, "orders", item.Order, &order) == nil {
		values["order.po"] = order.PO
		values["order.delivery"] = order.Delivery
		values["order.status"] = order.Status
		values["contract.no"] = firstNonEmpty(order.PO, order.ID)
		values["shipping.period"] = firstNonEmpty(order.Delivery, values["shipping.period"])
		if order.Currency != "" {
			values["currency"] = order.Currency
		}
	}
	values["destination.port"] = firstNonEmpty(values["destination.port"], values["customer.country"])

	lines := item.Lines
	if len(lines) == 0 && len(order.Lines) > 0 {
		lines = order.Lines
	}
	pdfLines := make([]string, 0, len(lines))
	for i, line := range lines {
		index := i + 1
		prefix := fmt.Sprintf("line.%d.", index)
		values[prefix+"no"] = strconv.Itoa(index)
		values[prefix+"model"] = firstNonEmpty(line.ProductID, line.ProductName)
		values[prefix+"description"] = firstNonEmpty(line.ProductName, line.ProductID)
		values[prefix+"quantity"] = cleanQuantity(line.Quantity)
		values[prefix+"unit"] = values["quantity.unit"]
		values[prefix+"unitPrice"] = formatPDFMoney(line.UnitPrice)
		values[prefix+"amount"] = formatPDFMoney(line.Amount)
		values[prefix+"packageCount"] = values["package.count"]
		values[prefix+"grossWeight"] = values["gross.weight"]
		values[prefix+"netWeight"] = values["net.weight"]
		values[prefix+"spec"] = firstNonEmpty(values["product.spec"], line.ProductName)
		var product Product
		if line.ProductID != "" && a.store.get(ctx, "products", line.ProductID, &product) == nil {
			values[prefix+"hs"] = product.HS
			values[prefix+"description"] = firstNonEmpty(product.Name, values[prefix+"description"])
			values[prefix+"spec"] = firstNonEmpty(product.Description, values[prefix+"spec"])
		}
		if values[prefix+"hs"] == "" {
			values[prefix+"hs"] = "8471504090"
		}
		if index <= 5 {
			pdfLines = append(pdfLines, fmt.Sprintf("%d  %s\n   %s\n   HS CODE %s\n   Qty %s   Rate %s   Amount %s",
				index, values[prefix+"model"], compactPDFText(values[prefix+"description"], 84), values[prefix+"hs"], values[prefix+"quantity"], values[prefix+"unitPrice"], values[prefix+"amount"]))
		}
	}
	values["pdf.pi.lines"] = strings.Join(pdfLines, "\n")
	return values
}

func documentTemplateTypeLabel(value string) string {
	switch strings.TrimSpace(value) {
	case "PI":
		return "PI 标准 PDF 模板"
	case "CI", "发票":
		return "发票 CI 标准模板"
	case "PL", "装箱单":
		return "装箱单 PL 标准模板"
	case "报关单":
		return "报关单标准模板"
	case "合同":
		return "合同标准模板"
	default:
		return "单据模板"
	}
}

func (a *businessAPI) renderWorkbookTemplateWithValues(ctx context.Context, item BusinessTemplate, values map[string]string, title string) ([]byte, string, error) {
	path, err := templateFilePath(item)
	if err != nil {
		return nil, "", err
	}
	source, err := zip.OpenReader(path)
	if err != nil {
		return nil, "", fmt.Errorf("读取 XLSM 模板失败: %w", err)
	}
	defer source.Close()

	var output bytes.Buffer
	writer := zip.NewWriter(&output)
	for _, file := range source.File {
		header := file.FileHeader
		entry, err := writer.CreateHeader(&header)
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		reader, err := file.Open()
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		switch file.Name {
		case "xl/worksheets/sheet1.xml":
			data = fillDeclarationWorksheet(data, values)
		case "xl/workbook.xml":
			data = updateWorkbookActiveSheet(data, values["document.type"])
		}
		if _, err := entry.Write(data); err != nil {
			_ = writer.Close()
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return output.Bytes(), safeDownloadName(title) + ".xlsm", nil
}

// renderWorkbookAsXLSX keeps the supplied workbook's sheets, styles, images and
// print setup, but removes macro-only package metadata so the download is a real
// .xlsx file rather than an XLSM file with a renamed suffix.
func renderWorkbookAsXLSX(source []byte, title string) ([]byte, string, error) {
	archive, err := zip.NewReader(bytes.NewReader(source), int64(len(source)))
	if err != nil {
		return nil, "", fmt.Errorf("读取 Excel 模板失败: %w", err)
	}

	var output bytes.Buffer
	writer := zip.NewWriter(&output)
	for _, file := range archive.File {
		if isWorkbookMacroPart(file.Name) {
			continue
		}
		reader, err := file.Open()
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		switch file.Name {
		case "[Content_Types].xml":
			data = normalizeXLSXContentTypes(data)
		case "xl/_rels/workbook.xml.rels":
			data = removeWorkbookMacroRelationships(data)
		}
		header := file.FileHeader
		entry, err := writer.CreateHeader(&header)
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		if _, err := entry.Write(data); err != nil {
			_ = writer.Close()
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return output.Bytes(), safeDownloadName(title) + ".xlsx", nil
}

func isWorkbookMacroPart(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	return strings.Contains(name, "vbaproject") || strings.Contains(name, "vbasignature")
}

func normalizeXLSXContentTypes(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("application/vnd.ms-excel.sheet.macroEnabled.main+xml"), []byte("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"))
	macroOverride := regexp.MustCompile(`(?i)<Override PartName="/[^"]*vba[^"]*" ContentType="[^"]*"\s*/>`)
	return macroOverride.ReplaceAll(data, nil)
}

func removeWorkbookMacroRelationships(data []byte) []byte {
	macroRelationship := regexp.MustCompile(`(?i)<Relationship\s+[^>]*(?:vbaProject|vbaSignature)[^>]*/>`)
	return macroRelationship.ReplaceAll(data, nil)
}

func normalizeDocumentDownloadFormat(format, documentType string) (string, error) {
	format = strings.ToLower(strings.TrimSpace(format))
	format = strings.TrimPrefix(format, ".")
	switch format {
	case "", "default":
		return "pdf", nil
	case "xls", "xlsx", "xlsm":
		format = "excel"
	case "doc", "docx":
		format = "word"
	}
	for _, supported := range supportedDocumentFormats(documentType) {
		if format == supported {
			return format, nil
		}
	}
	return "", fmt.Errorf("%s 仅支持下载格式: %s", documentType, strings.Join(documentFormatLabels(supportedDocumentFormats(documentType)), "、"))
}

func supportedDocumentFormats(documentType string) []string {
	switch documentType {
	case "PI":
		return []string{"pdf"}
	case "CI", "发票", "PL", "装箱单", "合同":
		return []string{"pdf", "word", "excel"}
	case "报关单":
		return []string{"pdf", "xml"}
	default:
		return []string{"pdf"}
	}
}

func documentFormatLabels(formats []string) []string {
	labels := make([]string, 0, len(formats))
	for _, format := range formats {
		switch format {
		case "pdf":
			labels = append(labels, "PDF")
		case "word":
			labels = append(labels, "Word")
		case "excel":
			labels = append(labels, "Excel")
		case "xml":
			labels = append(labels, "XML")
		}
	}
	return labels
}

func (a *businessAPI) renderDocumentDownload(ctx context.Context, item Document, tpl BusinessTemplate, values map[string]string, format string) ([]byte, string, string, error) {
	title := item.Type + "-" + item.ID
	switch format {
	case "excel":
		content, filename, err := a.renderWorkbookTemplateWithValues(ctx, tpl, values, title)
		if err != nil {
			return nil, "", "", err
		}
		content, filename, err = renderWorkbookAsXLSX(content, strings.TrimSuffix(filename, ".xlsm"))
		if err != nil {
			return nil, "", "", err
		}
		return content, filename, workbookContentType(".xlsx"), nil
	case "word":
		// Calc does not provide a DOCX export filter. CI, PL and contract use
		// their own built-in DOCX layouts, populated from the same fields as XLSX.
		content, filename, err := renderDocumentWordDOCX(item, values, title)
		if err != nil {
			return nil, "", "", err
		}
		return content, filename, "application/vnd.openxmlformats-officedocument.wordprocessingml.document", nil
	case "xml":
		content := renderCustomsDeclarationXML(item, values)
		return content, safeDownloadName(title) + ".xml", "application/xml; charset=utf-8", nil
	case "pdf":
		if item.Type == "PI" && isPDFTemplateFile(tpl.FileName) {
			content, filename, err := a.renderPDFTemplateWithValues(ctx, tpl, values, title)
			if err != nil {
				return nil, "", "", err
			}
			return content, filename, "application/pdf", nil
		}
		if isWorkbookTemplateFile(tpl.FileName) {
			content, _, err := a.renderWorkbookTemplateWithValues(ctx, tpl, values, title)
			if err != nil {
				return nil, "", "", err
			}
			converted, filename, err := renderWorkbookDocumentPDF(ctx, content, item.Type, title)
			if err == nil {
				return converted, filename, "application/pdf", nil
			}
			if !errors.Is(err, errOfficeConverterMissing) {
				return nil, "", "", err
			}
		}
		content, filename, err := renderStandardDocumentPDF(item, values, title)
		if err != nil {
			return nil, "", "", err
		}
		return content, filename, "application/pdf", nil
	default:
		return nil, "", "", fmt.Errorf("单据下载格式不受支持: %s", format)
	}
}

// renderWorkbookDocumentPDF first lets Calc resolve the supplied workbook's
// formulas. It then keeps only the requested document sheet and turns the
// calculated cells into values before the final PDF export. LibreOffice prints
// hidden sheets as well, so merely hiding the data-entry sheet is insufficient.
func renderWorkbookDocumentPDF(ctx context.Context, content []byte, documentType, title string) ([]byte, string, error) {
	calculated, _, err := convertOfficeDocument(ctx, content, ".xlsm", "xlsx", title+"-calculated")
	if err != nil {
		return nil, "", err
	}
	singleSheet, err := keepWorkbookDocumentSheet(calculated, documentType)
	if err != nil {
		return nil, "", err
	}
	return convertOfficeDocument(ctx, singleSheet, ".xlsx", "pdf", title)
}

func workbookDocumentSheetName(documentType string) string {
	switch documentType {
	case "合同":
		return "合同"
	case "CI", "发票":
		return "发票"
	case "PL", "装箱单":
		return "装箱单"
	default:
		return "报关单"
	}
}

func workbookDocumentSheetPath(documentType string) string {
	switch documentType {
	case "合同":
		return "xl/worksheets/sheet3.xml"
	case "CI", "发票":
		return "xl/worksheets/sheet4.xml"
	case "PL", "装箱单":
		return "xl/worksheets/sheet5.xml"
	default:
		return "xl/worksheets/sheet2.xml"
	}
}

func workbookDocumentPrintArea(documentType string) string {
	switch documentType {
	case "合同":
		return "合同!$A$1:$G$46"
	case "CI", "发票":
		return "发票!$A$1:$G$27"
	case "PL", "装箱单":
		return "装箱单!$A$1:$J$41"
	default:
		return "报关单!$A$1:$AL$34"
	}
}

// keepWorkbookDocumentSheet produces a self-contained workbook containing the
// one printable standard sheet. The input worksheet is intentionally removed
// only after Calc has recalculated it, and remaining formulas are flattened to
// their cached values so exported documents do not retain an external source.
func keepWorkbookDocumentSheet(source []byte, documentType string) ([]byte, error) {
	archive, err := zip.NewReader(bytes.NewReader(source), int64(len(source)))
	if err != nil {
		return nil, fmt.Errorf("读取计算后的 Excel 模板失败: %w", err)
	}
	targetName := workbookDocumentSheetName(documentType)
	targetPath := workbookDocumentSheetPath(documentType)
	var output bytes.Buffer
	writer := zip.NewWriter(&output)
	for _, file := range archive.File {
		name := file.Name
		if shouldSkipWorkbookPart(name, targetPath) {
			continue
		}
		reader, err := file.Open()
		if err != nil {
			_ = writer.Close()
			return nil, err
		}
		data, err := io.ReadAll(reader)
		_ = reader.Close()
		if err != nil {
			_ = writer.Close()
			return nil, err
		}
		switch name {
		case "xl/workbook.xml":
			data = keepWorkbookSheetXML(data, targetName, workbookDocumentPrintArea(documentType))
		case "xl/_rels/workbook.xml.rels":
			data = keepWorkbookSheetRelationshipXML(data, targetPath)
		case "[Content_Types].xml":
			data = keepWorkbookContentTypesXML(data, targetPath)
		default:
			if name == targetPath {
				data = flattenWorksheetFormulas(data)
			}
		}
		header := file.FileHeader
		entry, err := writer.CreateHeader(&header)
		if err != nil {
			_ = writer.Close()
			return nil, err
		}
		if _, err := entry.Write(data); err != nil {
			_ = writer.Close()
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func shouldSkipWorkbookPart(name, targetPath string) bool {
	if strings.HasPrefix(name, "xl/externalLinks/") || name == "xl/calcChain.xml" {
		return true
	}
	if strings.HasPrefix(name, "xl/worksheets/") && strings.HasSuffix(name, ".xml") && !strings.HasPrefix(name, "xl/worksheets/_rels/") {
		return name != targetPath
	}
	if strings.HasPrefix(name, "xl/worksheets/_rels/") && strings.HasSuffix(name, ".rels") {
		return name != "xl/worksheets/_rels/"+filepath.Base(targetPath)+".rels"
	}
	return false
}

func keepWorkbookSheetXML(data []byte, targetName, printArea string) []byte {
	sheetTag := regexp.MustCompile(`<sheet name="([^"]+)"[^>]*/>`)
	data = sheetTag.ReplaceAllFunc(data, func(tag []byte) []byte {
		parts := sheetTag.FindSubmatch(tag)
		if len(parts) != 2 || string(parts[1]) != targetName {
			return nil
		}
		attrs := regexp.MustCompile(`\s+state="[^"]*"`).ReplaceAll(tag, nil)
		return attrs
	})
	data = regexp.MustCompile(`\s+(?:firstSheet|activeTab)="\d+"`).ReplaceAll(data, nil)
	data = bytes.Replace(data, []byte("<workbookView"), []byte("<workbookView firstSheet=\"0\" activeTab=\"0\""), 1)
	data = regexp.MustCompile(`<definedNames>.*?</definedNames>`).ReplaceAll(data, nil)
	data = regexp.MustCompile(`<externalReferences>.*?</externalReferences>`).ReplaceAll(data, nil)
	if strings.TrimSpace(printArea) != "" {
		definedName := []byte(`<definedNames><definedName name="_xlnm.Print_Area" localSheetId="0">` + xmlText(printArea) + `</definedName></definedNames>`)
		data = bytes.Replace(data, []byte("</sheets>"), append([]byte("</sheets>"), definedName...), 1)
	}
	return data
}

func keepWorkbookSheetRelationshipXML(data []byte, targetPath string) []byte {
	relationshipTag := regexp.MustCompile(`<Relationship\s+[^>]*/>`)
	target := strings.TrimPrefix(targetPath, "xl/")
	return relationshipTag.ReplaceAllFunc(data, func(tag []byte) []byte {
		text := string(tag)
		if strings.Contains(text, "worksheets/") && !strings.Contains(text, `Target="`+target+`"`) {
			return nil
		}
		if strings.Contains(text, "externalLink") || strings.Contains(text, "calcChain") {
			return nil
		}
		return tag
	})
}

func keepWorkbookContentTypesXML(data []byte, targetPath string) []byte {
	override := regexp.MustCompile(`<Override\s+[^>]*/>`)
	target := "/" + targetPath
	return override.ReplaceAllFunc(data, func(tag []byte) []byte {
		text := string(tag)
		if strings.Contains(text, "/xl/worksheets/") && !strings.Contains(text, `PartName="`+target+`"`) {
			return nil
		}
		if strings.Contains(text, "externalLink") || strings.Contains(text, "calcChain") {
			return nil
		}
		return tag
	})
}

func flattenWorksheetFormulas(data []byte) []byte {
	return regexp.MustCompile(`<f(?:\s[^>]*)?>.*?</f>`).ReplaceAll(data, nil)
}

func convertOfficeDocument(ctx context.Context, content []byte, sourceExt, targetFormat, title string) ([]byte, string, error) {
	bin := officeConverterBinary()
	if bin == "" {
		return nil, "", fmt.Errorf("%w: %s", errOfficeConverterMissing, strings.ToUpper(targetFormat))
	}
	tempDir, err := os.MkdirTemp("", "sta100-template-*")
	if err != nil {
		return nil, "", err
	}
	defer os.RemoveAll(tempDir)
	base := safeDownloadName(title)
	source := filepath.Join(tempDir, base+sourceExt)
	if err := os.WriteFile(source, content, 0o600); err != nil {
		return nil, "", err
	}
	profileDir := filepath.Join(tempDir, "office-profile")
	cmd := exec.CommandContext(ctx, bin, "--headless", "--convert-to", targetFormat, "--outdir", tempDir, source)
	cmd.Env = append(os.Environ(), "STA100_OFFICE_PROFILE_DIR="+profileDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", fmt.Errorf("模板转换为 %s 失败: %v: %s", strings.ToUpper(targetFormat), err, strings.TrimSpace(string(output)))
	}
	candidates, _ := filepath.Glob(filepath.Join(tempDir, base+".*"))
	targetExt := "." + strings.TrimPrefix(strings.ToLower(targetFormat), ".")
	for _, candidate := range candidates {
		if strings.EqualFold(filepath.Ext(candidate), targetExt) {
			data, readErr := os.ReadFile(candidate)
			if readErr != nil {
				return nil, "", readErr
			}
			return data, base + targetExt, nil
		}
	}
	return nil, "", fmt.Errorf("模板转换为 %s 后未找到输出文件", strings.ToUpper(targetFormat))
}

func officeConverterBinary() string {
	if configured := strings.TrimSpace(os.Getenv("STA100_OFFICE_BIN")); configured != "" {
		if info, err := os.Stat(configured); err == nil && !info.IsDir() {
			return configured
		}
	}
	if executable, err := os.Executable(); err == nil {
		bundled := filepath.Join(filepath.Dir(executable), "..", "office", "bin", "soffice")
		if info, statErr := os.Stat(bundled); statErr == nil && !info.IsDir() {
			return bundled
		}
	}
	for _, name := range []string{"soffice", "libreoffice"} {
		if path, err := exec.LookPath(name); err == nil && strings.TrimSpace(path) != "" {
			return path
		}
	}
	return ""
}

func renderDocumentWordHTML(item Document, values map[string]string) []byte {
	body := standardDocumentHTML(item, values, true)
	return []byte("<!doctype html><html><head><meta charset=\"utf-8\"><title>" + html.EscapeString(item.Type+"-"+item.ID) + "</title>" + documentExportStyle() + "</head><body>" + body + "</body></html>")
}

// renderDocumentWordDOCX produces a standards-compliant DOCX fallback for
// deployments without LibreOffice. It uses the same business fields as the PDF
// export, so a document can always be generated instead of failing at download.
func renderDocumentWordDOCX(item Document, values map[string]string, title string) ([]byte, string, error) {
	var output bytes.Buffer
	writer := zip.NewWriter(&output)
	entries := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/></Types>`,
		"_rels/.rels":         `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/></Relationships>`,
		"word/document.xml":   documentWordXML(item, values),
	}
	for _, name := range []string{"[Content_Types].xml", "_rels/.rels", "word/document.xml"} {
		entry, err := writer.Create(name)
		if err != nil {
			_ = writer.Close()
			return nil, "", err
		}
		if _, err := io.WriteString(entry, entries[name]); err != nil {
			_ = writer.Close()
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return output.Bytes(), safeDownloadName(title) + ".docx", nil
}

func documentWordXML(item Document, values map[string]string) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:body>`)
	builder.WriteString(wordParagraph(documentTitle(item.Type), true, "28"))
	builder.WriteString(wordParagraph(documentWordNumberLine(item, values), false, "20"))
	builder.WriteString(documentWordPartyTable(values))

	rows := documentLineRows(values, 100)
	switch item.Type {
	case "PL", "装箱单":
		builder.WriteString(wordParagraph("Packing Details", true, "22"))
		builder.WriteString(documentWordTable([]string{"No.", "Product", "Description", "Qty", "Packages", "Net Wt.", "Gross Wt."}, documentPackingWordRows(rows, values)))
		builder.WriteString(wordParagraph("Total Packages: "+firstNonEmpty(values["package.count"], "0")+"    Packaging: "+firstNonEmpty(values["packaging.type"], "N/A"), true, "20"))
	case "合同":
		builder.WriteString(wordParagraph("Goods and Pricing", true, "22"))
		builder.WriteString(documentWordTable([]string{"No.", "Product", "Description", "Qty", "Unit Price", "Amount"}, documentContractWordRows(rows)))
		builder.WriteString(wordParagraph("Contract Terms", true, "22"))
		builder.WriteString(wordParagraph("Trade Term: "+values["trade.term"]+"\nPayment: "+values["settlement.method"]+"\nDelivery: "+values["shipping.period"]+"\nPort of Destination: "+values["destination.port"]+"\nPlace of Signing: "+values["contract.sign.place"], false, "18"))
		builder.WriteString(wordParagraph("Total Contract Value: "+firstNonEmpty(values["record.total"], item.Value), true, "22"))
	default:
		builder.WriteString(wordParagraph("Commercial Goods", true, "22"))
		builder.WriteString(documentWordTable([]string{"No.", "HS Code", "Product", "Description", "Qty", "Unit Price", "Amount"}, rows))
		builder.WriteString(wordParagraph("Invoice Total: "+firstNonEmpty(values["record.total"], item.Value), true, "22"))
		builder.WriteString(wordParagraph("Payment: "+values["settlement.method"]+"    Trade Term: "+values["trade.term"], false, "18"))
	}
	builder.WriteString(`<w:sectPr><w:pgSz w:w="11906" w:h="16838"/><w:pgMar w:top="900" w:right="900" w:bottom="900" w:left="900"/></w:sectPr></w:body></w:document>`)
	return builder.String()
}

func documentWordNumberLine(item Document, values map[string]string) string {
	switch item.Type {
	case "合同":
		return "Contract No: " + firstNonEmpty(values["contract.no"], values["record.id"]) + "    Date: " + values["record.date"] + "    Order: " + values["order.id"]
	case "PL", "装箱单":
		return "Packing List No: " + values["record.id"] + "    Date: " + values["record.date"] + "    Order: " + values["order.id"]
	default:
		return "Invoice No: " + values["record.id"] + "    Date: " + values["record.date"] + "    Order: " + values["order.id"]
	}
}

func documentWordPartyTable(values map[string]string) string {
	return documentWordTable([]string{"Buyer / Consignee", "Seller / Shipper"}, [][]string{{
		firstNonEmpty(values["customer.name"], "N/A") + "\n" + firstNonEmpty(values["customer.address"], "N/A") + "\n" + values["customer.country"],
		firstNonEmpty(values["supplier.company"], "N/A") + "\n" + firstNonEmpty(values["supplier.address"], "N/A") + "\n" + values["supplier.country"],
	}})
}

func documentPackingWordRows(rows [][]string, values map[string]string) [][]string {
	output := make([][]string, 0, len(rows))
	for _, row := range rows {
		if len(row) < 7 {
			continue
		}
		output = append(output, []string{row[0], row[2], row[3], row[4], firstNonEmpty(values["package.count"], "1"), values["net.weight"], values["gross.weight"]})
	}
	return output
}

func documentContractWordRows(rows [][]string) [][]string {
	output := make([][]string, 0, len(rows))
	for _, row := range rows {
		if len(row) < 7 {
			continue
		}
		output = append(output, []string{row[0], row[2], row[3], row[4], row[5], row[6]})
	}
	return output
}

func wordParagraph(value string, bold bool, size string) string {
	properties := ""
	if bold {
		properties = "<w:b/>"
	}
	value = xmlText(value)
	value = strings.ReplaceAll(value, "\n", `</w:t><w:br/><w:t xml:space="preserve">`)
	return `<w:p><w:pPr><w:spacing w:after="120"/></w:pPr><w:r><w:rPr>` + properties + `<w:sz w:val="` + size + `"/><w:rFonts w:ascii="Arial" w:hAnsi="Arial" w:eastAsia="Microsoft YaHei"/></w:rPr><w:t xml:space="preserve">` + value + `</w:t></w:r></w:p>`
}

func documentWordTable(headings []string, rows [][]string) string {
	var builder strings.Builder
	builder.WriteString(`<w:tbl><w:tblPr><w:tblW w:w="0" w:type="auto"/><w:tblBorders><w:top w:val="single" w:sz="6"/><w:left w:val="single" w:sz="6"/><w:bottom w:val="single" w:sz="6"/><w:right w:val="single" w:sz="6"/><w:insideH w:val="single" w:sz="4"/><w:insideV w:val="single" w:sz="4"/></w:tblBorders></w:tblPr>`)
	builder.WriteString(wordTableRow(headings, true))
	for _, row := range rows {
		builder.WriteString(wordTableRow(row, false))
	}
	if len(rows) == 0 {
		builder.WriteString(wordTableRow([]string{"No item lines"}, false))
	}
	builder.WriteString(`</w:tbl>`)
	return builder.String()
}

func wordTableRow(cells []string, bold bool) string {
	var builder strings.Builder
	builder.WriteString("<w:tr>")
	for _, value := range cells {
		properties := ""
		if bold {
			properties = "<w:b/>"
		}
		builder.WriteString(`<w:tc><w:tcPr><w:tcW w:w="1300" w:type="dxa"/></w:tcPr><w:p><w:r><w:rPr>` + properties + `<w:sz w:val="16"/><w:rFonts w:ascii="Arial" w:hAnsi="Arial" w:eastAsia="Microsoft YaHei"/></w:rPr><w:t xml:space="preserve">` + xmlText(value) + `</w:t></w:r></w:p></w:tc>`)
	}
	builder.WriteString("</w:tr>")
	return builder.String()
}

func renderStandardDocumentPDF(item Document, values map[string]string, title string) ([]byte, string, error) {
	if err := ensurePDFUserFonts(); err != nil {
		return nil, "", fmt.Errorf("PDF 字体初始化失败: %w", err)
	}
	jsonPayload, err := json.Marshal(documentPDFCreateModel(item, values))
	if err != nil {
		return nil, "", err
	}
	var output bytes.Buffer
	if err := api.Create(nil, bytes.NewReader(jsonPayload), &output, nil); err != nil {
		return nil, "", fmt.Errorf("PDF 单据生成失败: %w", err)
	}
	return output.Bytes(), safeDownloadName(title) + ".pdf", nil
}

func documentPDFCreateModel(item Document, values map[string]string) map[string]any {
	font := pdfTemplateFontName()
	lines := documentLineRows(values, 10)
	text := []map[string]any{
		{"value": documentTitle(item.Type), "pos": []float64{44, 790}, "font": map[string]any{"name": "$title", "size": 19}},
		{"value": "STRATRONIX", "pos": []float64{430, 792}, "font": map[string]any{"name": "$strong", "size": 15}},
		{"value": "Document No: " + values["record.id"], "pos": []float64{44, 752}, "font": map[string]any{"name": "$body", "size": 10}},
		{"value": "Date: " + values["record.date"], "pos": []float64{360, 752}, "font": map[string]any{"name": "$body", "size": 10}},
		{"value": "Customer: " + values["customer.name"], "pos": []float64{44, 724}, "width": 230, "font": map[string]any{"name": "$body", "size": 10}},
		{"value": "Supplier: " + values["supplier.company"], "pos": []float64{310, 724}, "width": 240, "font": map[string]any{"name": "$body", "size": 10}},
		{"value": "Customer Address: " + values["customer.address"], "pos": []float64{44, 694}, "width": 500, "font": map[string]any{"name": "$body", "size": 9}},
		{"value": "Order: " + values["order.id"] + "    PO/Contract: " + values["contract.no"], "pos": []float64{44, 666}, "width": 500, "font": map[string]any{"name": "$body", "size": 9}},
		{"value": "Items", "pos": []float64{44, 630}, "font": map[string]any{"name": "$strong", "size": 12}},
	}
	y := 606.0
	for _, row := range lines {
		text = append(text, map[string]any{"value": strings.Join(row, "    "), "pos": []float64{44, y}, "width": 510, "font": map[string]any{"name": "$body", "size": 8}})
		y -= 28
	}
	text = append(text,
		map[string]any{"value": "Total: " + firstNonEmpty(values["record.total"], item.Value), "pos": []float64{380, y - 8}, "font": map[string]any{"name": "$strong", "size": 12}},
		map[string]any{"value": documentTermsBlock(item.Type, values), "pos": []float64{44, y - 56}, "width": 500, "font": map[string]any{"name": "$body", "size": 8}},
	)
	return map[string]any{
		"paper":  "A4P",
		"origin": "LowerLeft",
		"guides": false,
		"fonts": map[string]any{
			"body":   map[string]any{"name": font, "size": 9},
			"strong": map[string]any{"name": font, "size": 11},
			"title":  map[string]any{"name": font, "size": 18},
		},
		"pages": map[string]any{"1": map[string]any{"content": map[string]any{"text": text}}},
	}
}

func renderCustomsDeclarationXML(item Document, values map[string]string) []byte {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	builder.WriteString(`<CustomsDeclaration>` + "\n")
	writeXMLNode(&builder, 1, "DocumentNo", item.ID)
	writeXMLNode(&builder, 1, "OrderNo", values["order.id"])
	writeXMLNode(&builder, 1, "DeclarationDate", values["record.date"])
	writeXMLNode(&builder, 1, "Shipper", values["domestic.shipper"])
	writeXMLNode(&builder, 1, "Consignee", values["customer.name"])
	writeXMLNode(&builder, 1, "ConsigneeAddress", values["customer.address"])
	writeXMLNode(&builder, 1, "ExitCustoms", values["exit.customs"])
	writeXMLNode(&builder, 1, "TradeMode", values["trade.mode"])
	writeXMLNode(&builder, 1, "TradeTerm", values["trade.term"])
	writeXMLNode(&builder, 1, "Currency", values["currency"])
	writeXMLNode(&builder, 1, "DestinationPort", values["destination.port"])
	writeXMLNode(&builder, 1, "TransportMode", values["transport.mode"])
	writeXMLNode(&builder, 1, "Freight", values["freight"])
	writeXMLNode(&builder, 1, "PackageCount", values["package.count"])
	writeXMLNode(&builder, 1, "GrossWeight", values["gross.weight"])
	writeXMLNode(&builder, 1, "NetWeight", values["net.weight"])
	builder.WriteString("  <Items>\n")
	for i, row := range documentLineRows(values, 10) {
		if len(row) == 0 || strings.TrimSpace(strings.Join(row, "")) == "" {
			continue
		}
		prefix := fmt.Sprintf("line.%d.", i+1)
		builder.WriteString("    <Item>\n")
		writeXMLNode(&builder, 3, "LineNo", values[prefix+"no"])
		writeXMLNode(&builder, 3, "HSCode", values[prefix+"hs"])
		writeXMLNode(&builder, 3, "Model", values[prefix+"model"])
		writeXMLNode(&builder, 3, "Description", values[prefix+"description"])
		writeXMLNode(&builder, 3, "Spec", values[prefix+"spec"])
		writeXMLNode(&builder, 3, "Quantity", values[prefix+"quantity"])
		writeXMLNode(&builder, 3, "Unit", values[prefix+"unit"])
		writeXMLNode(&builder, 3, "UnitPrice", values[prefix+"unitPrice"])
		writeXMLNode(&builder, 3, "Amount", values[prefix+"amount"])
		builder.WriteString("    </Item>\n")
	}
	builder.WriteString("  </Items>\n")
	writeXMLNode(&builder, 1, "DeclarationElement", values["declaration.element"])
	builder.WriteString(`</CustomsDeclaration>` + "\n")
	return []byte(builder.String())
}

func writeXMLNode(builder *strings.Builder, indent int, name, value string) {
	builder.WriteString(strings.Repeat("  ", indent))
	builder.WriteString("<")
	builder.WriteString(name)
	builder.WriteString(">")
	builder.WriteString(xmlText(value))
	builder.WriteString("</")
	builder.WriteString(name)
	builder.WriteString(">\n")
}

func standardDocumentHTML(item Document, values map[string]string, forWord bool) string {
	return `<main class="document-export"><header><div><h1>` + html.EscapeString(documentTitle(item.Type)) + `</h1><p>` + html.EscapeString(values["record.id"]) + `</p></div><strong>STRATRONIX</strong></header><section class="meta-grid"><div><b>Customer</b><p>` + html.EscapeString(values["customer.name"]) + `<br>` + html.EscapeString(values["customer.address"]) + `</p></div><div><b>Supplier</b><p>` + html.EscapeString(values["supplier.company"]) + `<br>` + html.EscapeString(values["supplier.address"]) + `</p></div><div><b>Document</b><p>Date: ` + html.EscapeString(values["record.date"]) + `<br>Order: ` + html.EscapeString(values["order.id"]) + `</p></div><div><b>Terms</b><p>` + html.EscapeString(documentTermsBlock(item.Type, values)) + `</p></div></section>` + documentLinesTableHTML(values) + `<p class="total">Total: ` + html.EscapeString(firstNonEmpty(values["record.total"], item.Value)) + `</p></main>`
}

func documentExportStyle() string {
	return `<style>body{font-family:Arial,"Microsoft YaHei",sans-serif;margin:32px;color:#172033}.document-export header{display:flex;justify-content:space-between;align-items:flex-start;border-bottom:2px solid #172033;padding-bottom:14px;margin-bottom:18px}.document-export h1{margin:0;font-size:24px}.meta-grid{display:grid;grid-template-columns:1fr 1fr;gap:14px;margin:16px 0}.meta-grid div{border:1px solid #cbd5e1;padding:10px}.meta-grid p{margin:5px 0 0;line-height:1.55}.doc-lines{width:100%;border-collapse:collapse;margin-top:16px}.doc-lines th,.doc-lines td{border:1px solid #cbd5e1;padding:8px;text-align:left}.doc-lines th{background:#eef2f7}.total{text-align:right;font-size:18px;font-weight:700;margin-top:18px}</style>`
}

func documentLinesTableHTML(values map[string]string) string {
	var builder strings.Builder
	builder.WriteString(`<table class="doc-lines"><thead><tr><th>No.</th><th>HS Code</th><th>Product</th><th>Description</th><th>Qty</th><th>Unit Price</th><th>Amount</th></tr></thead><tbody>`)
	rows := documentLineRows(values, 10)
	if len(rows) == 0 {
		builder.WriteString(`<tr><td colspan="7">No item lines</td></tr>`)
	}
	for _, row := range rows {
		builder.WriteString("<tr>")
		for _, value := range row {
			builder.WriteString("<td>")
			builder.WriteString(html.EscapeString(value))
			builder.WriteString("</td>")
		}
		builder.WriteString("</tr>")
	}
	builder.WriteString(`</tbody></table>`)
	return builder.String()
}

func documentLineRows(values map[string]string, maxRows int) [][]string {
	rows := make([][]string, 0, maxRows)
	for i := 1; i <= maxRows; i++ {
		prefix := fmt.Sprintf("line.%d.", i)
		if strings.TrimSpace(values[prefix+"description"]+values[prefix+"model"]+values[prefix+"amount"]+values[prefix+"quantity"]) == "" {
			continue
		}
		rows = append(rows, []string{values[prefix+"no"], values[prefix+"hs"], values[prefix+"model"], values[prefix+"description"], values[prefix+"quantity"], values[prefix+"unitPrice"], values[prefix+"amount"]})
	}
	return rows
}

func documentTitle(documentType string) string {
	switch documentType {
	case "CI", "发票":
		return "Commercial Invoice"
	case "PL", "装箱单":
		return "Packing List"
	case "合同":
		return "Sales Contract"
	case "报关单":
		return "Customs Declaration"
	case "PI":
		return "Proforma Invoice"
	default:
		return documentType
	}
}

func documentTermsBlock(documentType string, values map[string]string) string {
	parts := []string{
		"Trade Term: " + values["trade.term"],
		"Settlement: " + values["settlement.method"],
		"Shipping Period: " + values["shipping.period"],
		"Destination: " + values["destination.port"],
	}
	if documentType == "报关单" {
		parts = append(parts, "Declaration Element: "+values["declaration.element"])
	}
	return strings.Join(parts, "\n")
}

func fillDeclarationWorksheet(data []byte, values map[string]string) []byte {
	fields := map[string]string{
		"E6":  "record.date",
		"E7":  "domestic.shipper",
		"O7":  "domestic.shipper.address",
		"E11": "customer.name",
		"O11": "customer.address",
		"E15": "customs.code",
		"M15": "exit.customs",
		"U15": "contract.no",
		"M17": "tax.exemption.nature",
		"U17": "contract.sign.date",
		"M19": "trade.mode",
		"U19": "contract.sign.place",
		"E21": "port.departure",
		"M21": "trade.term",
		"U21": "shipping.period",
		"E23": "container.no",
		"M23": "settlement.method",
		"U23": "invoice.no",
		"E25": "destination.port",
		"M25": "exemption.mode",
		"E27": "destination.port",
		"M27": "currency",
		"U27": "origin.country",
		"E29": "transport.mode",
		"U29": "domestic.source.location",
		"E31": "freight",
		"M31": "packaging.type",
		"L33": "declaration.element",
	}
	out := data
	for cell, key := range fields {
		out = setWorksheetCellInlineString(out, cell, values[key])
	}
	for i := 1; i <= 10; i++ {
		row := 43 + (i-1)*2
		prefix := fmt.Sprintf("line.%d.", i)
		out = setWorksheetCellInlineString(out, fmt.Sprintf("A%d", row), values[prefix+"no"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("B%d", row), values[prefix+"hs"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("E%d", row), values[prefix+"description"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("N%d", row), values[prefix+"packageCount"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("P%d", row), values[prefix+"grossWeight"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("R%d", row), values[prefix+"netWeight"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("T%d", row), values[prefix+"quantity"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("V%d", row), values[prefix+"unit"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("W%d", row), values[prefix+"unitPrice"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("E%d", row+1), values[prefix+"spec"])
		out = setWorksheetCellInlineString(out, fmt.Sprintf("T%d", row+1), "")
		out = setWorksheetCellInlineString(out, fmt.Sprintf("V%d", row+1), "")
	}
	out = setWorksheetCellInlineString(out, "N63", values["package.count"])
	out = setWorksheetCellInlineString(out, "P63", values["gross.weight"])
	out = setWorksheetCellInlineString(out, "R63", values["net.weight"])
	return out
}

func setWorksheetCellInlineString(data []byte, cell, value string) []byte {
	cell = strings.TrimSpace(cell)
	startToken := []byte(`<c r="` + cell + `"`)
	start := bytes.Index(data, startToken)
	if start < 0 {
		return data
	}
	tagEnd := bytes.IndexByte(data[start:], '>')
	if tagEnd < 0 {
		return data
	}
	tagEnd += start
	attrs := data[start+len(startToken) : tagEnd]
	end := tagEnd + 1
	if len(attrs) > 0 && attrs[len(attrs)-1] == '/' {
		attrs = bytes.TrimSpace(attrs[:len(attrs)-1])
	} else {
		closeToken := []byte(`</c>`)
		closeIndex := bytes.Index(data[tagEnd:], closeToken)
		if closeIndex < 0 {
			return data
		}
		end = tagEnd + closeIndex + len(closeToken)
	}
	attrs = regexp.MustCompile(`\s+t="[^"]*"`).ReplaceAll(attrs, nil)
	attrs = regexp.MustCompile(`\s+cm="[^"]*"`).ReplaceAll(attrs, nil)
	attrs = regexp.MustCompile(`\s+vm="[^"]*"`).ReplaceAll(attrs, nil)
	replacement := []byte(`<c r="` + cell + `"` + string(attrs) + ` t="inlineStr"><is><t xml:space="preserve">` + xmlText(value) + `</t></is></c>`)
	out := make([]byte, 0, len(data)-end+start+len(replacement))
	out = append(out, data[:start]...)
	out = append(out, replacement...)
	out = append(out, data[end:]...)
	return out
}

func updateWorkbookActiveSheet(data []byte, documentType string) []byte {
	activeTab := "1"
	targetSheet := "报关单"
	switch documentType {
	case "合同":
		activeTab = "2"
		targetSheet = "合同"
	case "CI", "发票":
		activeTab = "3"
		targetSheet = "发票"
	case "PL", "装箱单":
		activeTab = "4"
		targetSheet = "装箱单"
	case "报关单":
		activeTab = "1"
		targetSheet = "报关单"
	}
	re := regexp.MustCompile(`activeTab="\d+"`)
	if re.Match(data) {
		data = re.ReplaceAll(data, []byte(`activeTab="`+activeTab+`"`))
	}
	if bytes.Contains(data, []byte("<calcPr")) {
		data = regexp.MustCompile(`\s+fullCalcOnLoad="[^"]*"`).ReplaceAll(data, nil)
		data = regexp.MustCompile(`\s+forceFullCalc="[^"]*"`).ReplaceAll(data, nil)
		data = regexp.MustCompile(`<calcPr([^>]*)/>`).ReplaceAll(data, []byte(`<calcPr$1 fullCalcOnLoad="1" forceFullCalc="1"/>`))
	}
	// The original workbook contains every document sheet. The generated file
	// must expose only the requested document (plus the hidden input sheet), so
	// a PDF export produces that document rather than a five-page workbook.
	sheetTag := regexp.MustCompile(`<sheet name="([^"]+)"([^>]*)/>`)
	data = sheetTag.ReplaceAllFunc(data, func(tag []byte) []byte {
		parts := sheetTag.FindSubmatch(tag)
		if len(parts) != 3 {
			return tag
		}
		name := string(parts[1])
		attrs := regexp.MustCompile(`\s+state="[^"]*"`).ReplaceAll(parts[2], nil)
		if name != targetSheet {
			attrs = append(attrs, []byte(` state="hidden"`)...)
		}
		return append(append([]byte(`<sheet name="`+name+`"`), attrs...), []byte(`/>`)...)
	})
	return data
}

func formatPDFDate(value string) string {
	parsed, err := time.Parse("2006-01-02", strings.TrimSpace(value))
	if err != nil {
		return value
	}
	return parsed.Format("01.02.2006")
}
