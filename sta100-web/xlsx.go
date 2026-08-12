package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func writeXLSX(w http.ResponseWriter, filename, sheetName string, rows [][]string) {
	var output bytes.Buffer
	archive := zip.NewWriter(&output)
	files := map[string]string{
		"[Content_Types].xml":        `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types"><Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/><Default Extension="xml" ContentType="application/xml"/><Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/><Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/></Types>`,
		"_rels/.rels":                `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/></Relationships>`,
		"xl/_rels/workbook.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"><Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/></Relationships>`,
		"xl/workbook.xml":            `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"><sheets><sheet name="` + xmlText(sheetName) + `" sheetId="1" r:id="rId1"/></sheets></workbook>`,
		"xl/worksheets/sheet1.xml":   worksheetXML(rows),
	}
	for _, name := range []string{"[Content_Types].xml", "_rels/.rels", "xl/workbook.xml", "xl/_rels/workbook.xml.rels", "xl/worksheets/sheet1.xml"} {
		entry, err := archive.Create(name)
		if err != nil {
			writeAPIError(w, http.StatusInternalServerError, "XLSX_CREATE_FAILED", "Excel 文件生成失败")
			return
		}
		if _, err := entry.Write([]byte(files[name])); err != nil {
			writeAPIError(w, http.StatusInternalServerError, "XLSX_CREATE_FAILED", "Excel 文件生成失败")
			return
		}
	}
	if err := archive.Close(); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "XLSX_CREATE_FAILED", "Excel 文件生成失败")
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(filename, `"`, "")+`"`)
	w.Header().Set("Content-Length", strconv.Itoa(output.Len()))
	_, _ = w.Write(output.Bytes())
}

func worksheetXML(rows [][]string) string {
	var output strings.Builder
	output.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)
	for rowIndex, row := range rows {
		output.WriteString(fmt.Sprintf(`<row r="%d">`, rowIndex+1))
		for columnIndex, value := range row {
			output.WriteString(fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t xml:space="preserve">%s</t></is></c>`, xlsxColumn(columnIndex+1), rowIndex+1, xmlText(value)))
		}
		output.WriteString(`</row>`)
	}
	output.WriteString(`</sheetData></worksheet>`)
	return output.String()
}

func xlsxColumn(column int) string {
	var result string
	for column > 0 {
		column--
		result = string(rune('A'+column%26)) + result
		column /= 26
	}
	return result
}

func xmlText(value string) string {
	value = strings.ReplaceAll(value, "&", "&amp;")
	value = strings.ReplaceAll(value, "<", "&lt;")
	value = strings.ReplaceAll(value, ">", "&gt;")
	value = strings.ReplaceAll(value, `"`, "&quot;")
	value = strings.ReplaceAll(value, "'", "&apos;")
	return value
}
