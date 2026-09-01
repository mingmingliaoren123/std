package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
)

const (
	quotePDFWidth          = 595.0
	quotePDFHeight         = 842.0
	quotePDFLeft           = 57.0
	quotePDFRight          = 538.0
	quotePDFPink           = "#E6417D"
	quotePDFLightPink      = "#FBE8F0"
	quotePDFText           = "#172033"
	quotePDFMuted          = "#5E5E5E"
	quotePDFLine           = "#E5E7EB"
	quotePDFTableTop       = 380.0
	quotePDFContinueTop    = 570.0
	quotePDFTableHeaderH   = 39.0
	quotePDFRowH           = 18.0
	quotePDFFirstPageFinal = 8
	quotePDFFirstPageRows  = 14
	quotePDFNextPageFinal  = 18
	quotePDFNextPageRows   = 22
	quotePDFLastPageMin    = 4
)

func renderDynamicQuotePDF(item Quote, values map[string]string, title string) ([]byte, string, error) {
	if err := ensurePDFUserFonts(); err != nil {
		return nil, "", fmt.Errorf("PDF 字体初始化失败: %w", err)
	}
	lines := item.Lines
	if len(lines) == 0 {
		lines = []BusinessLine{{ProductID: "", ProductName: "待补充", Quantity: 0, UnitPrice: 0, Amount: 0}}
	}
	pages := splitQuoteLines(lines)
	model, err := dynamicQuotePDFCreateModel(item, values, pages)
	if err != nil {
		return nil, "", err
	}
	payload, err := json.Marshal(model)
	if err != nil {
		return nil, "", err
	}
	var output bytes.Buffer
	if err := api.Create(nil, bytes.NewReader(payload), &output, nil); err != nil {
		return nil, "", fmt.Errorf("报价单 PDF 生成失败: %w", err)
	}
	return output.Bytes(), safeDownloadName(title) + ".pdf", nil
}

func splitQuoteLines(lines []BusinessLine) [][]BusinessLine {
	if len(lines) <= quotePDFFirstPageFinal {
		return [][]BusinessLine{lines}
	}
	firstEnd := min(quotePDFFirstPageRows, len(lines))
	out := [][]BusinessLine{append([]BusinessLine{}, lines[:firstEnd]...)}
	offset := firstEnd
	for len(lines)-offset > quotePDFNextPageFinal {
		take := min(quotePDFNextPageRows, len(lines)-offset)
		remainingAfterTake := len(lines) - offset - take
		if remainingAfterTake > 0 && remainingAfterTake < quotePDFLastPageMin {
			take = max(1, len(lines)-offset-quotePDFLastPageMin)
		}
		end := offset + take
		out = append(out, append([]BusinessLine{}, lines[offset:end]...))
		offset = end
	}
	if offset < len(lines) {
		out = append(out, append([]BusinessLine{}, lines[offset:]...))
	}
	return out
}

func dynamicQuotePDFCreateModel(item Quote, values map[string]string, pages [][]BusinessLine) (map[string]any, error) {
	font := pdfTemplateFontName()
	pdfPages := map[string]any{}
	lineOffset := 0
	for i, pageLines := range pages {
		pageNumber := i + 1
		isLast := pageNumber == len(pages)
		tableTop := quotePDFTableY(pageNumber == 1)
		pdfPages[fmt.Sprintf("%d", pageNumber)] = map[string]any{
			"content": map[string]any{
				"box":  quotePDFPageBoxes(pageNumber == 1, isLast, len(pageLines), tableTop),
				"text": quotePDFPageText(item, values, pageLines, pageNumber, len(pages), isLast, tableTop, lineOffset),
			},
		}
		lineOffset += len(pageLines)
	}
	return map[string]any{
		"paper":  "A4P",
		"origin": "LowerLeft",
		"guides": false,
		"fonts": map[string]any{
			"body":   map[string]any{"name": font, "size": 8, "col": quotePDFText},
			"strong": map[string]any{"name": font, "size": 10, "col": quotePDFMuted},
			"title":  map[string]any{"name": font, "size": 18, "col": quotePDFPink},
		},
		"pages": pdfPages,
	}, nil
}

func quotePDFTableY(first bool) float64 {
	if first {
		return quotePDFTableTop
	}
	return quotePDFContinueTop
}

func quotePDFPageBoxes(first, last bool, rowCount int, tableTop float64) []map[string]any {
	boxes := []map[string]any{
		quotePDFBox(quotePDFLeft, 700, quotePDFRight-quotePDFLeft, 1.4, quotePDFPink, quotePDFPink, 0),
		quotePDFBox(quotePDFLeft, tableTop, quotePDFRight-quotePDFLeft, quotePDFTableHeaderH, quotePDFPink, quotePDFPink, 0),
	}
	if first {
		boxes = append(boxes,
			quotePDFBox(quotePDFLeft, 427, quotePDFRight-quotePDFLeft, 112, quotePDFLightPink, quotePDFPink, 1),
			quotePDFBox(quotePDFLeft+(quotePDFRight-quotePDFLeft)/2, 427, 0.8, 112, quotePDFPink, quotePDFPink, 0),
			quotePDFBox(66, 510, 228, 0.8, quotePDFPink, quotePDFPink, 0),
			quotePDFBox(307, 510, 228, 0.8, quotePDFPink, quotePDFPink, 0),
		)
	}
	tableBottom := tableTop - quotePDFRowH*float64(rowCount)
	for i := 0; i <= rowCount; i++ {
		y := tableTop - quotePDFRowH*float64(i)
		boxes = append(boxes, quotePDFBox(quotePDFLeft, y, quotePDFRight-quotePDFLeft, 0.45, quotePDFLine, quotePDFLine, 0))
	}
	if last {
		totalTop := math.Max(230, tableBottom-28)
		termsTop := quotePDFTermsTop(totalTop)
		boxes = append(boxes,
			quotePDFBox(quotePDFLeft+205, totalTop-36, 214, 1.2, quotePDFPink, quotePDFPink, 0),
			quotePDFBox(quotePDFLeft, termsTop-2, 6, 14, quotePDFPink, quotePDFPink, 0),
			quotePDFBox(quotePDFLeft, 31, 230, 1.1, "#333333", "#333333", 0),
			quotePDFBox(quotePDFRight-230, 31, 230, 1.1, "#333333", "#333333", 0),
		)
	}
	return boxes
}

func quotePDFBox(x, y, w, h float64, fill, border string, borderWidth int) map[string]any {
	box := map[string]any{
		"pos":     []float64{x, y},
		"width":   w,
		"height":  h,
		"fillCol": fill,
	}
	if borderWidth > 0 {
		box["border"] = map[string]any{"width": borderWidth, "col": border}
	}
	return box
}

func quotePDFPageText(item Quote, values map[string]string, lines []BusinessLine, page, pageCount int, last bool, tableTop float64, lineOffset int) []map[string]any {
	text := quotePDFCommonText(values, page, pageCount)
	if page == 1 {
		text = append(text, quotePDFFirstPageText(values)...)
	}
	text = append(text, quotePDFTableHeaderText(tableTop, values["quote.currency"])...)
	y := tableTop - 16
	for i, line := range lines {
		rowY := y - quotePDFRowH*float64(i)
		text = append(text, quotePDFLineText(lineOffset+i+1, line, values["quote.currency"], rowY)...)
	}
	if last {
		tableBottom := tableTop - quotePDFRowH*float64(len(lines))
		totalTop := math.Max(230, tableBottom-28)
		termsTop := quotePDFTermsTop(totalTop)
		text = append(text, quotePDFTotalsText(values, totalTop)...)
		text = append(text, quotePDFTermsText(values, termsTop)...)
	}
	return text
}

func quotePDFTermsTop(totalTop float64) float64 {
	return math.Min(216, math.Max(150, totalTop-100))
}

func quotePDFCommonText(values map[string]string, page, pageCount int) []map[string]any {
	text := []map[string]any{
		quotePDFTextBoxWithFont("STRATRONIX", 63, 781, 24, quotePDFPink, "l", 190, "$title"),
		quotePDFTextBoxWithFont("STRATRONIX / 鼎图", 449, 754, 10, quotePDFPink, "l", 100, "$strong"),
		quotePDFTextBox("Room 1203D, Building C6, Hengfeng Industrial City, Shenzhen", 76, 735, 8, quotePDFMuted, "l", 360),
		quotePDFTextBox("Tel: 86-755-23086689 | Web: www.stratronix.ai", 332, 722, 8, quotePDFMuted, "l", 206),
		quotePDFTextBox("Email: info@stratronix.ai", 443, 709, 8, quotePDFMuted, "l", 95),
		quotePDFTextBoxWithFont("Q U O T A T I O N", 216, 663, 22, quotePDFPink, "l", 210, "$title"),
		quotePDFTextBoxWithFont("报 价 单", 279, 645, 10, quotePDFMuted, "l", 70, "$strong"),
	}
	if pageCount > 1 {
		text = append(text, quotePDFTextBox(fmt.Sprintf("Page %d / %d", page, pageCount), 483, 681, 8, quotePDFMuted, "r", 55))
	}
	return text
}

func quotePDFFirstPageText(values map[string]string) []map[string]any {
	return []map[string]any{
		quotePDFTextBoxWithFont("报价单号 / Quote\nNo.", 65, 606, 10, quotePDFMuted, "l", 105, "$strong"),
		quotePDFTextBox(values["pdf.record.id"], 171, 606, 8.5, quotePDFText, "l", 120),
		quotePDFTextBoxWithFont("报价日期 / Date", 282, 606, 10, quotePDFMuted, "l", 105, "$strong"),
		quotePDFTextBox(values["record.date"], 388, 606, 8.5, quotePDFText, "l", 80),
		quotePDFTextBoxWithFont("有效期 / Valid Until", 65, 571, 10, quotePDFMuted, "l", 115, "$strong"),
		quotePDFTextBox(values["quote.validity"], 171, 571, 8.5, quotePDFText, "l", 100),
		quotePDFTextBoxWithFont("付款条件 / Terms", 282, 571, 10, quotePDFMuted, "l", 105, "$strong"),
		quotePDFTextBox(values["pdf.quote.paymentTerms"], 388, 571, 8.2, quotePDFText, "l", 150),
		quotePDFTextBoxWithFont("客户 / CUSTOMER (TO)", 66, 522, 11, quotePDFPink, "l", 220, "$strong"),
		quotePDFTextBoxWithFont("供方 / SUPPLIER (FROM)", 307, 522, 11, quotePDFPink, "l", 220, "$strong"),
		quotePDFTextBoxWithFont("公司:", 66, 498, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(compactPDFText(values["customer.name"], 24), 128, 498, 7.0, quotePDFText, "l", 160),
		quotePDFTextBoxWithFont("国家:", 66, 483, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(values["customer.country"], 128, 483, 8.2, quotePDFText, "l", 188),
		quotePDFTextBoxWithFont("联系人:", 66, 468, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(values["pdf.customer.contact"], 128, 468, 8.2, quotePDFText, "l", 188),
		quotePDFTextBoxWithFont("邮箱:", 66, 453, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(values["customer.email"], 128, 453, 8.2, quotePDFText, "l", 188),
		quotePDFTextBoxWithFont("地址:", 66, 438, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(compactPDFText(firstNonEmpty(values["customer.address"], values["pdf.customer.address"]), 52), 128, 438, 7.8, quotePDFText, "l", 180),
		quotePDFTextBoxWithFont("公司:", 307, 498, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(values["supplier.company"], 372, 498, 8.2, quotePDFText, "l", 150),
		quotePDFTextBoxWithFont("联系人:", 307, 483, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(values["supplier.contact"], 372, 483, 8.2, quotePDFText, "l", 150),
		quotePDFTextBoxWithFont("电话:", 307, 468, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(values["supplier.phone"], 372, 468, 8.2, quotePDFText, "l", 150),
		quotePDFTextBoxWithFont("邮箱:", 307, 453, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(values["supplier.email"], 372, 453, 8.2, quotePDFText, "l", 150),
		quotePDFTextBoxWithFont("网址:", 307, 438, 9, quotePDFMuted, "l", 60, "$strong"),
		quotePDFTextBox(values["supplier.website"], 372, 438, 8.2, quotePDFText, "l", 150),
	}
}

func quotePDFTableHeaderText(tableTop float64, currency string) []map[string]any {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		currency = "USD"
	}
	return []map[string]any{
		quotePDFTextBoxWithFont("#", 67, tableTop+18, 9, "#FFFFFF", "c", 20, "$strong"),
		quotePDFTextBoxWithFont("型号 /", 91, tableTop+23, 10, "#FFFFFF", "l", 55, "$strong"),
		quotePDFTextBoxWithFont("Model", 91, tableTop+10, 10, "#FFFFFF", "l", 55, "$strong"),
		quotePDFTextBoxWithFont("描述 / Description", 159, tableTop+18, 10, "#FFFFFF", "l", 175, "$strong"),
		quotePDFTextBoxWithFont("数量 /", 349, tableTop+23, 10, "#FFFFFF", "c", 45, "$strong"),
		quotePDFTextBoxWithFont("Qty", 349, tableTop+10, 10, "#FFFFFF", "c", 45, "$strong"),
		quotePDFTextBoxWithFont("单价 ("+currency+")", 430, tableTop+18, 10, "#FFFFFF", "r", 70, "$strong"),
		quotePDFTextBoxWithFont("小计 ("+currency+")", 514, tableTop+18, 10, "#FFFFFF", "r", 65, "$strong"),
	}
}

func quotePDFLineText(no int, line BusinessLine, currency string, y float64) []map[string]any {
	return []map[string]any{
		quotePDFTextBox(fmt.Sprintf("%d", no), 67, y, 7, quotePDFText, "c", 20),
		quotePDFTextBoxWithFont(compactPDFText(firstNonEmpty(line.ProductID, line.ProductName), 22), 91, y, 7.2, quotePDFText, "l", 60, "$strong"),
		quotePDFTextBox(compactPDFText(firstNonEmpty(line.ProductName, line.ProductID), 46), 159, y, 7.2, quotePDFText, "l", 175),
		quotePDFTextBox(cleanQuantity(line.Quantity), 349, y, 7.2, quotePDFText, "c", 45),
		quotePDFTextBox(formatPDFMoney(line.UnitPrice), 430, y, 7.2, quotePDFText, "r", 70),
		quotePDFTextBox(formatPDFMoney(line.Amount), 514, y, 7.2, quotePDFText, "r", 65),
	}
}

func quotePDFTotalsText(values map[string]string, top float64) []map[string]any {
	y := top
	destination := firstNonEmpty(values["quote.destination"], "Rotterdam")
	return []map[string]any{
		quotePDFTextBoxWithFont("产品小计 / Subtotal (FOB Shenzhen)", 268, y, 10, quotePDFMuted, "l", 230, "$strong"),
		quotePDFTextBox(values["pdf.record.subtotal"], 526, y, 8, quotePDFText, "r", 65),
		quotePDFTextBoxWithFont("运费 / Freight (→ "+destination+")", 268, y-28, 10, quotePDFMuted, "l", 230, "$strong"),
		quotePDFTextBox(values["pdf.record.freight"], 526, y-28, 8, quotePDFText, "r", 65),
		quotePDFTextBoxWithFont("总计 / TOTAL (FOB "+destination+")", 268, y-66, 12, quotePDFMuted, "l", 245, "$strong"),
		quotePDFTextBox(values["pdf.record.total"], 526, y-66, 10, quotePDFText, "r", 65),
	}
}

func quotePDFTermsText(values map[string]string, top float64) []map[string]any {
	return []map[string]any{
		quotePDFTextBoxWithFont("条款与条件 / Terms & Conditions", 69, top, 11, quotePDFPink, "l", 260, "$strong"),
		quotePDFTextBox("价格条件 / Price Terms: "+values["quote.priceTerms"], 61, top-18, 6.8, quotePDFText, "l", 440),
		quotePDFTextBox("付款条件 / Payment: "+values["quote.paymentTerms"], 61, top-31, 6.8, quotePDFText, "l", 440),
		quotePDFTextBox("交期 / Lead Time: "+values["quote.leadTime"], 61, top-44, 6.8, quotePDFText, "l", 440),
		quotePDFTextBox("质保 / Warranty: "+values["quote.warranty"], 61, top-57, 6.8, quotePDFText, "l", 440),
		quotePDFTextBox("认证 / Certification: "+values["quote.certification"], 61, top-70, 6.8, quotePDFText, "l", 440),
		quotePDFTextBox("客户签字 / Customer Signature & Date", 57, 18, 8, quotePDFMuted, "l", 230),
		quotePDFTextBox("供方签字 / Supplier Signature & Date\nSTRATRONIX / 鼎图", 397, 12, 8, quotePDFMuted, "l", 142),
		quotePDFTextBox("STRATRONIX / 鼎图  |  86-755-23086689  |  www.stratronix.ai  |  info@stratronix.ai", 142, 4, 7.5, quotePDFMuted, "l", 330),
	}
}

func quotePDFTextBox(value string, x, y, fontSize float64, color, align string, width float64) map[string]any {
	return quotePDFTextBoxWithFont(value, x, y, fontSize, color, align, width, "$body")
}

func quotePDFTextBoxWithFont(value string, x, y, fontSize float64, color, align string, width float64, fontName string) map[string]any {
	box := map[string]any{
		"value": strings.TrimSpace(value),
		"pos":   []float64{x, y},
		"font":  map[string]any{"name": fontName, "size": fontSize, "col": color},
		"align": align,
	}
	if width > 0 {
		box["width"] = width
	}
	return box
}
