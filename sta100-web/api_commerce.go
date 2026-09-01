package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (a *businessAPI) quotesRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.quotesCollection(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "deliver" {
		a.quoteDeliver(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "convert-order" {
		a.quoteConvertOrder(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "download" {
		a.quoteDownload(w, r, parts[0])
		return
	}
	a.quoteItem(w, r, parts[0])
}

func (a *businessAPI) quotesCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := listRecords[Quote](r.Context(), a.store, "quotes")
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		status := r.URL.Query().Get("status")
		validFrom, validTo := r.URL.Query().Get("valid_from"), r.URL.Query().Get("valid_to")
		filtered := make([]Quote, 0, len(items))
		for _, item := range items {
			if query != "" && !containsFold(item.ID+" "+item.Subject+" "+item.Customer+" "+item.Products, query) {
				continue
			}
			if status != "" && status != "all" && item.Status != status {
				continue
			}
			if validFrom != "" && item.Valid < validFrom {
				continue
			}
			if validTo != "" && item.Valid > validTo {
				continue
			}
			filtered = append(filtered, item)
		}
		sortCommerce(filtered, r.URL.Query().Get("sort"), r.URL.Query().Get("direction"), func(item Quote) string { return item.Value }, func(item Quote) string { return item.Updated })
		listResponse(w, filtered)
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var item Quote
		if err := decodeJSONBody(w, r, &item); err != nil {
			return
		}
		if message := a.prepareQuote(r.Context(), &item); message != "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_QUOTE", message)
			return
		}
		id, err := a.store.nextSequence(r.Context(), "quotes", "QUO-"+time.Now().Format("2006"), 4)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		item.ID, item.Status, item.Updated = id, "Draft", currentText()
		if item.Owner == "" {
			item.Owner = requestOperator(r)
		}
		if err := a.store.create(r.Context(), "quotes", id, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "create", "quote", id, requestOperator(r), map[string]any{"total": item.Value})
		writeJSON(w, http.StatusCreated, item)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) quoteItem(w http.ResponseWriter, r *http.Request, id string) {
	var item Quote
	if err := a.store.get(r.Context(), "quotes", id, &item); err != nil {
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
		if item.Status == "Accepted" || item.Status == "Rejected" || item.Status == "Archived" {
			writeAPIError(w, http.StatusConflict, "QUOTE_LOCKED", "报价单已进入锁定状态，不能修改明细")
			return
		}
		var request Quote
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.ID = id
		if request.Status == "" {
			request.Status = item.Status
		}
		if !validQuoteStatus(request.Status) {
			writeAPIError(w, http.StatusBadRequest, "INVALID_QUOTE_STATUS", "报价状态不受支持")
			return
		}
		if message := a.prepareQuote(r.Context(), &request); message != "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_QUOTE", message)
			return
		}
		request.Updated = currentText()
		if err := a.store.put(r.Context(), "quotes", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "quote", id, requestOperator(r), map[string]any{"total": request.Value})
		writeJSON(w, http.StatusOK, request)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if a.quoteReferenced(r.Context(), id) {
			writeAPIError(w, http.StatusConflict, "QUOTE_IN_USE", "报价单已生成订单，不能删除")
			return
		}
		if item.Status != "Draft" && item.Status != "Rejected" {
			writeAPIError(w, http.StatusConflict, "QUOTE_DELETE_FORBIDDEN", "仅草稿或已拒绝报价可以归档")
			return
		}
		if err := a.store.softDelete(r.Context(), "quotes", id); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "archive", "quote", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, map[string]any{"archived": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) quoteDeliver(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var item Quote
	if err := a.store.get(r.Context(), "quotes", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	var request EmailSendRequest
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	if request.Subject == "" || request.Body == "" {
		writeAPIError(w, http.StatusBadRequest, "EMAIL_CONTENT_REQUIRED", "请填写邮件主题和正文")
		return
	}
	if err := a.sendQuoteEmail(r.Context(), item, request); err != nil {
		writeAPIError(w, http.StatusBadGateway, "QUOTE_EMAIL_FAILED", err.Error())
		return
	}
	if item.Status == "Draft" {
		item.Status = "Delivered"
		item.Updated = currentText()
		_ = a.store.put(r.Context(), "quotes", item.ID, item)
	}
	a.store.audit(r.Context(), "email", "quote", item.ID, requestOperator(r), map[string]any{"to": request.To, "attachment": request.AttachRecord})
	writeJSON(w, http.StatusOK, map[string]any{"sent": true, "quote": item, "message": "报价单邮件已通过 SMTP 发送"})
}

func (a *businessAPI) quoteConvertOrder(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var quote Quote
	if err := a.store.get(r.Context(), "quotes", id, &quote); err != nil {
		writeBusinessError(w, err)
		return
	}
	if quote.Status != "Accepted" {
		writeAPIError(w, http.StatusConflict, "QUOTE_NOT_ACCEPTED", "只有已接受报价可以转换为订单")
		return
	}
	if a.quoteReferenced(r.Context(), id) {
		writeAPIError(w, http.StatusConflict, "ORDER_ALREADY_EXISTS", "该报价已生成订单")
		return
	}
	lines := make([]BusinessLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		line.UnitPrice = line.UnitPrice * (1 - line.Discount/100)
		line.Discount = 0
		line.Amount = line.Quantity * line.UnitPrice
		lines = append(lines, line)
	}
	order := Order{Customer: quote.Customer, Quote: quote.ID, Currency: quote.Currency, Status: "Confirmed", Delivery: quote.Valid, Progress: 0, Lines: lines, Updated: currentText()}
	if message := a.prepareOrder(r.Context(), &order); message != "" {
		writeAPIError(w, http.StatusBadRequest, "INVALID_ORDER", message)
		return
	}
	orderID, err := a.store.nextSequence(r.Context(), "orders", "SO-"+time.Now().Format("2006"), 4)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	order.ID = orderID
	if err := a.store.create(r.Context(), "orders", orderID, order); err != nil {
		writeBusinessError(w, err)
		return
	}
	a.store.audit(r.Context(), "convert_order", "quote", id, requestOperator(r), map[string]any{"orderId": orderID})
	writeJSON(w, http.StatusCreated, order)
}

func (a *businessAPI) quoteDownload(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	var item Quote
	if err := a.store.get(r.Context(), "quotes", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	tpl, err := a.templateForQuote(r.Context(), item)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	if isStratronixQuotePDFTemplate(tpl) {
		values := a.quoteTemplateValues(r.Context(), item, tpl)
		content, filename, err := renderDynamicQuotePDF(item, values, "报价单-"+item.ID)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		writeAttachment(w, filename, "application/pdf", content)
		return
	}
	content, filename, contentType, err := a.renderBusinessTemplate(r.Context(), tpl, a.quoteTemplateValues(r.Context(), item, tpl), "报价单-"+item.ID)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	writeAttachment(w, filename, contentType, content)
}

func (a *businessAPI) prepareQuote(ctx context.Context, item *Quote) string {
	if strings.TrimSpace(item.Subject) == "" {
		return "报价主题不能为空"
	}
	if !a.activeCustomerExists(ctx, item.Customer) {
		return "关联客户不存在或已归档"
	}
	if item.Valid == "" {
		return "有效期不能为空"
	}
	if _, err := time.Parse("2006-01-02", item.Valid); err != nil {
		return "有效期格式必须为 YYYY-MM-DD"
	}
	if item.Currency == "" {
		item.Currency = "EUR"
	}
	lines, total, message := a.normalizeLines(ctx, item.Lines, false)
	if message != "" {
		return message
	}
	item.Lines = lines
	total += item.Freight
	item.Value, item.Products = formatMoney(total, item.Currency), lineSummary(lines)
	return ""
}

func (a *businessAPI) ordersRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.ordersCollection(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "documents" {
		a.orderDocuments(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "deliver" {
		a.orderDeliver(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "download" {
		a.orderDownload(w, r, parts[0])
		return
	}
	a.orderItem(w, r, parts[0])
}

func (a *businessAPI) orderDeliver(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var item Order
	if err := a.store.get(r.Context(), "orders", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	var request EmailSendRequest
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	if request.Subject == "" || request.Body == "" {
		writeAPIError(w, http.StatusBadRequest, "EMAIL_CONTENT_REQUIRED", "请填写邮件主题和正文")
		return
	}
	if err := a.sendOrderEmail(r.Context(), item, request); err != nil {
		writeAPIError(w, http.StatusBadGateway, "ORDER_EMAIL_FAILED", err.Error())
		return
	}
	item.Updated = currentText()
	_ = a.store.put(r.Context(), "orders", item.ID, item)
	a.store.audit(r.Context(), "email", "order", item.ID, requestOperator(r), map[string]any{"to": request.To, "documents": request.DocumentIDs})
	writeJSON(w, http.StatusOK, map[string]any{"sent": true, "order": item, "message": "订单邮件已通过 SMTP 发送"})
}

func (a *businessAPI) ordersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := listRecords[Order](r.Context(), a.store, "orders")
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		status := r.URL.Query().Get("status")
		deliveryFrom, deliveryTo := r.URL.Query().Get("delivery_from"), r.URL.Query().Get("delivery_to")
		filtered := make([]Order, 0, len(items))
		for _, item := range items {
			if query != "" && !containsFold(item.ID+" "+item.Customer+" "+item.Products+" "+item.Quote, query) {
				continue
			}
			if status != "" && status != "all" && item.Status != status {
				continue
			}
			if deliveryFrom != "" && item.Delivery < deliveryFrom {
				continue
			}
			if deliveryTo != "" && item.Delivery > deliveryTo {
				continue
			}
			filtered = append(filtered, item)
		}
		sortCommerce(filtered, r.URL.Query().Get("sort"), r.URL.Query().Get("direction"), func(item Order) string { return item.Value }, func(item Order) string { return item.Updated })
		listResponse(w, filtered)
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var item Order
		if err := decodeJSONBody(w, r, &item); err != nil {
			return
		}
		if message := a.prepareOrder(r.Context(), &item); message != "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_ORDER", message)
			return
		}
		id, err := a.store.nextSequence(r.Context(), "orders", "SO-"+time.Now().Format("2006"), 4)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		item.ID, item.Updated = id, currentText()
		if item.Status == "" {
			item.Status = "Confirmed"
		}
		if !validOrderStatus(item.Status) {
			writeAPIError(w, http.StatusBadRequest, "INVALID_ORDER_STATUS", "订单状态不受支持")
			return
		}
		if err := a.store.create(r.Context(), "orders", id, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "create", "order", id, requestOperator(r), map[string]any{"total": item.Value})
		writeJSON(w, http.StatusCreated, item)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) orderItem(w http.ResponseWriter, r *http.Request, id string) {
	var item Order
	if err := a.store.get(r.Context(), "orders", id, &item); err != nil {
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
		if item.Status == "Production" || item.Status == "Shipped" || item.Status == "Completed" || item.Status == "Cancelled" {
			writeAPIError(w, http.StatusConflict, "ORDER_LOCKED", "订单进入生产后关键字段已锁定")
			return
		}
		var request Order
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.ID = id
		if request.Status == "" {
			request.Status = item.Status
		}
		if !validOrderStatus(request.Status) {
			writeAPIError(w, http.StatusBadRequest, "INVALID_ORDER_STATUS", "订单状态不受支持")
			return
		}
		if message := a.prepareOrder(r.Context(), &request); message != "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_ORDER", message)
			return
		}
		request.Updated = currentText()
		if err := a.store.put(r.Context(), "orders", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "order", id, requestOperator(r), map[string]any{"total": request.Value})
		writeJSON(w, http.StatusOK, request)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if a.orderReferenced(r.Context(), id) {
			writeAPIError(w, http.StatusConflict, "ORDER_IN_USE", "订单已有单据，不能删除")
			return
		}
		if item.Status != "Draft" && item.Status != "Confirmed" && item.Status != "Cancelled" {
			writeAPIError(w, http.StatusConflict, "ORDER_DELETE_FORBIDDEN", "生产或发运订单不能删除")
			return
		}
		if err := a.store.softDelete(r.Context(), "orders", id); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "archive", "order", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, map[string]any{"archived": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) prepareOrder(ctx context.Context, item *Order) string {
	if !a.activeCustomerExists(ctx, item.Customer) {
		return "关联客户不存在或已归档"
	}
	if item.Delivery == "" {
		return "预计交付日期不能为空"
	}
	if _, err := time.Parse("2006-01-02", item.Delivery); err != nil {
		return "预计交付日期格式必须为 YYYY-MM-DD"
	}
	if item.Quote != "" {
		var quote Quote
		if err := a.store.get(ctx, "quotes", item.Quote, &quote); err != nil {
			return "关联报价不存在"
		}
		if quote.Customer != item.Customer {
			return "关联报价与客户不一致"
		}
	}
	if item.Currency == "" {
		item.Currency = "EUR"
	}
	lines, total, message := a.normalizeLines(ctx, item.Lines, true)
	if message != "" {
		return message
	}
	item.Lines, item.Value, item.Products = lines, formatMoney(total, item.Currency), lineSummary(lines)
	if item.Progress < 0 || item.Progress > 100 {
		return "订单进度必须在 0-100 之间"
	}
	return ""
}

func (a *businessAPI) orderDocuments(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMutation(w, r, http.MethodPost) {
		return
	}
	var request struct {
		Types    []string `json:"types"`
		Template string   `json:"template"`
		Language string   `json:"language"`
	}
	if err := decodeJSONBody(w, r, &request); err != nil {
		return
	}
	if len(request.Types) == 0 {
		request.Types = []string{"PI", "CI", "PL", "报关单", "合同"}
	}
	for _, documentType := range request.Types {
		if !validDocumentType(documentType) {
			writeAPIError(w, http.StatusBadRequest, "INVALID_DOCUMENT_TYPE", "单据类型仅支持 PI、CI、PL、报关单、合同")
			return
		}
	}
	var sourceOrder Order
	if err := a.store.get(r.Context(), "orders", id, &sourceOrder); err != nil {
		writeBusinessError(w, err)
		return
	}
	created := make([]Document, 0, len(request.Types))
	for _, documentType := range request.Types {
		document, err := a.createDocumentFromOrder(r.Context(), id, documentType, request.Template, request.Language)
		if err != nil {
			if errors.Is(err, errRecordNotFound) {
				writeBusinessError(w, err)
			} else {
				writeAPIError(w, http.StatusBadRequest, "INVALID_DOCUMENT", err.Error())
			}
			return
		}
		created = append(created, document)
	}
	a.store.audit(r.Context(), "generate_documents", "order", id, requestOperator(r), map[string]any{"count": len(created)})
	writeJSON(w, http.StatusCreated, map[string]any{"items": created, "total": len(created)})
}

func (a *businessAPI) orderDownload(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	var item Order
	if err := a.store.get(r.Context(), "orders", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	tpl := builtInTemplate("order")
	content, filename, contentType, err := a.renderBusinessTemplate(r.Context(), tpl, orderTemplateValues(item), "订单-"+item.ID)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	writeAttachment(w, filename, contentType, content)
}

func (a *businessAPI) documentsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.documentsCollection(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "download" {
		a.documentDownload(w, r, parts[0])
		return
	}
	a.documentItem(w, r, parts[0])
}

func (a *businessAPI) documentsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := listRecords[Document](r.Context(), a.store, "documents")
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		documentType, status := r.URL.Query().Get("type"), r.URL.Query().Get("status")
		filtered := make([]Document, 0, len(items))
		for _, item := range items {
			if query != "" && !containsFold(item.ID+" "+item.Customer+" "+item.Order+" "+item.Template, query) {
				continue
			}
			if documentType != "" && documentType != "all" && item.Type != documentType {
				continue
			}
			if status != "" && status != "all" && item.Status != status {
				continue
			}
			filtered = append(filtered, item)
		}
		listResponse(w, filtered)
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var request Document
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		item, err := a.createDocumentFromOrder(r.Context(), request.Order, request.Type, request.Template, request.Language)
		if err != nil {
			if errors.Is(err, errRecordNotFound) {
				writeBusinessError(w, err)
			} else {
				writeAPIError(w, http.StatusBadRequest, "INVALID_DOCUMENT", err.Error())
			}
			return
		}
		a.store.audit(r.Context(), "create", "document", item.ID, requestOperator(r), map[string]any{"orderId": item.Order})
		writeJSON(w, http.StatusCreated, item)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) documentItem(w http.ResponseWriter, r *http.Request, id string) {
	var item Document
	if err := a.store.get(r.Context(), "documents", id, &item); err != nil {
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
		if item.Status == "Confirmed" || item.Status == "Sent" || item.Status == "Cancelled" {
			writeAPIError(w, http.StatusConflict, "DOCUMENT_LOCKED", "已确认、已发送或已取消单据不能编辑")
			return
		}
		var request Document
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.Order = firstNonEmpty(request.Order, item.Order)
		if request.Order != item.Order {
			writeAPIError(w, http.StatusConflict, "DOCUMENT_ORDER_IMMUTABLE", "单据关联订单不能修改")
			return
		}
		request.Type = firstNonEmpty(request.Type, item.Type)
		request.Template = firstNonEmpty(request.Template, item.Template)
		request.Language = firstNonEmpty(request.Language, item.Language)
		request.Status = firstNonEmpty(request.Status, item.Status)
		request.ID, request.Customer, request.Value, request.Lines, request.Updated = id, item.Customer, item.Value, item.Lines, currentText()
		if !validDocumentType(request.Type) {
			writeAPIError(w, http.StatusBadRequest, "INVALID_DOCUMENT_TYPE", "单据类型仅支持 PI、CI、PL、报关单、合同")
			return
		}
		if !validDocumentStatus(request.Status) {
			writeAPIError(w, http.StatusBadRequest, "INVALID_DOCUMENT_STATUS", "单据状态不受支持")
			return
		}
		if err := a.store.put(r.Context(), "documents", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "document", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, request)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if item.Status != "Draft" && item.Status != "Review" {
			writeAPIError(w, http.StatusConflict, "DOCUMENT_DELETE_FORBIDDEN", "仅草稿或待复核单据可以归档")
			return
		}
		if err := a.store.softDelete(r.Context(), "documents", id); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "archive", "document", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, map[string]any{"archived": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) createDocumentFromOrder(ctx context.Context, orderID, documentType, template, language string) (Document, error) {
	if !validDocumentType(documentType) {
		return Document{}, errors.New("单据类型仅支持 PI、CI、PL、报关单、合同")
	}
	var order Order
	if err := a.store.get(ctx, "orders", orderID, &order); err != nil {
		return Document{}, err
	}
	if template == "" {
		template = "各类型默认模板"
	}
	if language == "" {
		language = "英文"
	}
	prefix := documentType
	if documentType == "报关单" {
		prefix = "CD"
	} else if documentType == "合同" {
		prefix = "CT"
	}
	id, err := a.store.nextSequence(ctx, "documents", prefix+"-"+time.Now().Format("20060102"), 3)
	if err != nil {
		return Document{}, err
	}
	document := Document{ID: id, Type: documentType, Customer: order.Customer, Order: order.ID, Template: template, Language: language, Status: "Review", Value: order.Value, Lines: append([]BusinessLine{}, order.Lines...), Updated: currentText()}
	if err := a.store.create(ctx, "documents", id, document); err != nil {
		return Document{}, err
	}
	return document, nil
}

func (a *businessAPI) documentDownload(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	var item Document
	if err := a.store.get(r.Context(), "documents", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	tpl, err := a.templateForDocument(r.Context(), item)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	format, err := normalizeDocumentDownloadFormat(r.URL.Query().Get("format"), item.Type)
	if err != nil {
		writeAPIError(w, http.StatusBadRequest, "INVALID_DOCUMENT_FORMAT", err.Error())
		return
	}
	values := a.documentTemplateValues(r.Context(), item, tpl)
	content, filename, contentType, err := a.renderDocumentDownload(r.Context(), item, tpl, values, format)
	if err != nil {
		if errors.Is(err, errOfficeConverterMissing) {
			writeAPIError(w, http.StatusServiceUnavailable, "DOCUMENT_CONVERTER_MISSING", "当前环境未安装 LibreOffice/soffice，无法按原模板导出该格式")
			return
		}
		writeBusinessError(w, err)
		return
	}
	writeAttachment(w, filename, contentType, content)
}

func (a *businessAPI) quoteTemplateValues(ctx context.Context, item Quote, tpl BusinessTemplate) map[string]string {
	subtotal := lineAmountTotal(item.Lines)
	total := quoteTemplateTotal(item)
	values := map[string]string{
		"company.name":  "STRATRONIX",
		"customer.name": item.Customer,
		"record.id":     item.ID,
		"record.date":   time.Now().Format("2006-01-02"),
		"record.total":  item.Value,
		"record.lines":  businessLinesHTML(item.Lines),
		"quote.subject": item.Subject,
		"quote.valid":   item.Valid,
		"quote.terms":   item.Terms,
	}
	if item.Currency == "" {
		values["quote.currency"] = currencyFromMoney(item.Value)
	} else {
		values["quote.currency"] = item.Currency
	}
	values["record.total"] = formatMoney(total, values["quote.currency"])
	for key, value := range mergeTemplateDefaultValues(tpl.DefaultValues, defaultQuoteTemplateValues()) {
		values[key] = value
	}
	values["record.subtotal"] = formatMoney(subtotal, values["quote.currency"])
	values["record.freight"] = formatMoney(item.Freight, values["quote.currency"])
	values["record.tax"] = formatMoney(item.Tax, values["quote.currency"])
	for key, value := range item.TemplateFields {
		if strings.TrimSpace(key) != "" {
			values[strings.TrimSpace(key)] = strings.TrimSpace(value)
		}
	}
	values["quote.paymentTerms"] = firstNonEmpty(values["quote.paymentTerms"], item.Terms)
	if customer, ok := a.customerByName(ctx, item.Customer); ok {
		values["customer.country"] = firstNonEmpty(values["customer.country"], customer.Country)
		values["customer.city"] = firstNonEmpty(values["customer.city"], customer.City)
		values["customer.contact"] = firstNonEmpty(values["customer.contact"], customer.Contact)
		values["customer.phone"] = firstNonEmpty(values["customer.phone"], customer.Phone)
		values["customer.email"] = firstNonEmpty(values["customer.email"], customer.Email)
		values["customer.website"] = firstNonEmpty(values["customer.website"], customer.Website)
		values["customer.address"] = firstNonEmpty(values["customer.address"], customer.Description)
	}
	if strings.TrimSpace(item.Terms) != "" && strings.TrimSpace(item.TemplateFields["quote.paymentTerms"]) == "" {
		values["quote.paymentTerms"] = item.Terms
	}
	values["quote.validity"] = firstNonEmpty(values["quote.validity"], quoteValidityText(item.Valid))
	values["quote.validityDetail"] = quoteValidityDetail(values["quote.validityDays"], item.Valid)
	values["pdf.record.id"] = stratronixQuoteNumber(item.ID)
	values["pdf.quote.paymentTerms"] = compactPDFText(values["quote.paymentTerms"], 42)
	values["pdf.record.subtotal"] = formatPDFMoney(subtotal)
	values["pdf.record.freight"] = formatPDFMoney(item.Freight)
	values["pdf.record.total"] = formatPDFMoney(total)
	values["pdf.quote.priceTermsLine"] = "价格条件 / Price Terms: " + values["quote.priceTerms"]
	values["pdf.quote.paymentTermsLine"] = "付款条件 / Payment: " + values["quote.paymentTerms"]
	values["pdf.quote.leadTimeLine"] = "交期 / Lead Time: " + values["quote.leadTime"]
	values["pdf.quote.warrantyLine"] = "质保 / Warranty: " + values["quote.warranty"]
	values["pdf.quote.certificationLine"] = "认证 / Certification: " + values["quote.certification"]
	values["pdf.quote.packagingLine"] = "包装 / Packaging: " + values["quote.packaging"]
	values["pdf.quote.validityDetailLine"] = "报价有效期 / Validity: " + values["quote.validityDetail"]
	values["pdf.quote.noteLine"] = "备注 / Note: " + values["quote.note"]
	values["pdf.customer.contact"] = compactPDFContact(values["customer.contact"])
	values["pdf.customer.address"] = compactPDFText(values["customer.address"], 34)
	for i, line := range item.Lines {
		index := i + 1
		prefix := fmt.Sprintf("line.%d.", index)
		values[prefix+"no"] = strconv.Itoa(index)
		values[prefix+"model"] = firstNonEmpty(line.ProductID, line.ProductName)
		values[prefix+"description"] = firstNonEmpty(line.ProductName, line.ProductID)
		values[prefix+"quantity"] = cleanQuantity(line.Quantity)
		values[prefix+"unitPrice"] = formatMoney(line.UnitPrice, values["quote.currency"])
		values[prefix+"amount"] = formatMoney(line.Amount, values["quote.currency"])
		values["pdf."+prefix+"unitPrice"] = formatPDFMoney(line.UnitPrice)
		values["pdf."+prefix+"amount"] = formatPDFMoney(line.Amount)
	}
	return values
}

func quoteTemplateTotal(item Quote) float64 {
	if len(item.Lines) == 0 && item.Freight == 0 {
		return moneyNumber(item.Value)
	}
	return lineAmountTotal(item.Lines) + item.Freight
}

func (a *businessAPI) templateForQuote(ctx context.Context, _ Quote) (BusinessTemplate, error) {
	var tpl BusinessTemplate
	if err := a.store.get(ctx, businessTemplateKind, builtInQuotePDFTemplateID, &tpl); err == nil && tpl.Kind == "quote" && !strings.EqualFold(tpl.Status, "Archived") {
		return tpl, nil
	}
	return a.defaultTemplate(ctx, "quote")
}

func orderTemplateValues(item Order) map[string]string {
	return map[string]string{
		"company.name":   "STRATRONIX",
		"customer.name":  item.Customer,
		"record.id":      item.ID,
		"record.date":    time.Now().Format("2006-01-02"),
		"record.total":   item.Value,
		"record.lines":   businessLinesHTML(item.Lines),
		"order.po":       item.PO,
		"order.delivery": item.Delivery,
		"order.status":   item.Status,
	}
}

func (a *businessAPI) templateForDocument(ctx context.Context, item Document) (BusinessTemplate, error) {
	templateID := builtInDocumentInvoiceTemplateID
	switch item.Type {
	case "PI":
		templateID = builtInPIPDFTemplateID
	case "CI", "发票":
		templateID = builtInDocumentInvoiceTemplateID
	case "PL", "装箱单":
		templateID = builtInDocumentPackingTemplateID
	case "报关单":
		templateID = builtInDocumentCustomsTemplateID
	case "合同":
		templateID = builtInDocumentContractTemplateID
	}
	var tpl BusinessTemplate
	if err := a.store.get(ctx, businessTemplateKind, templateID, &tpl); err == nil && tpl.Kind == "document" && !strings.EqualFold(tpl.Status, "Archived") {
		return tpl, nil
	}
	return a.defaultTemplate(ctx, "document")
}

func (a *businessAPI) customerByName(ctx context.Context, name string) (Customer, bool) {
	items, err := listRecords[Customer](ctx, a.store, "accounts")
	if err != nil {
		return Customer{}, false
	}
	name = strings.TrimSpace(name)
	for _, item := range items {
		if !item.Archived && strings.EqualFold(strings.TrimSpace(item.Name), name) {
			return item, true
		}
	}
	return Customer{}, false
}

func currencyFromMoney(value string) string {
	upper := strings.ToUpper(strings.TrimSpace(value))
	for _, currency := range []string{"USD", "EUR", "CNY", "GBP"} {
		if strings.HasPrefix(upper, currency) || strings.Contains(upper, " "+currency+" ") {
			return currency
		}
	}
	return "EUR"
}

func lineAmountTotal(lines []BusinessLine) float64 {
	total := 0.0
	for _, line := range lines {
		total += line.Amount
	}
	return total
}

func cleanQuantity(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func formatPDFMoney(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func quoteValidityDetail(validityDays, valid string) string {
	validityDays = strings.TrimSpace(validityDays)
	valid = strings.TrimSpace(valid)
	if validityDays == "" {
		validityDays = "30"
	}
	if valid == "" {
		return validityDays + " days"
	}
	return validityDays + " days (until " + valid + ")"
}

func quotePDFSummaryLine(lines []BusinessLine) BusinessLine {
	if len(lines) == 0 {
		return BusinessLine{}
	}
	summary := lines[0]
	if len(lines) == 1 {
		return summary
	}
	sameProduct := true
	totalQty := 0.0
	totalAmount := 0.0
	for _, line := range lines {
		if line.ProductID != summary.ProductID || line.ProductName != summary.ProductName {
			sameProduct = false
		}
		totalQty += line.Quantity
		totalAmount += line.Amount
	}
	if sameProduct {
		summary.Quantity = totalQty
		summary.Amount = totalAmount
		if totalQty > 0 {
			summary.UnitPrice = totalAmount / totalQty
		}
		return summary
	}
	return BusinessLine{
		ProductID:   "MIXED",
		ProductName: lineSummary(lines),
		Quantity:    totalQty,
		UnitPrice:   0,
		Amount:      totalAmount,
	}
}

func stratronixQuoteNumber(id string) string {
	id = strings.TrimSpace(id)
	if id == "" {
		return ""
	}
	id = strings.TrimPrefix(strings.TrimPrefix(id, "QUO-"), "Q-")
	if strings.HasPrefix(id, "Q") {
		return "STRAT-" + id
	}
	return "STRAT-Q" + id
}

func compactPDFContact(value string) string {
	value = strings.TrimSpace(value)
	for _, marker := range []string{"（", "(", "，", ",", ";", "；"} {
		if index := strings.Index(value, marker); index > 0 {
			value = strings.TrimSpace(value[:index])
			break
		}
	}
	return compactPDFText(value, 34)
}

func compactPDFText(value string, limit int) string {
	value = strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
	if limit <= 0 || len([]rune(value)) <= limit {
		return value
	}
	runes := []rune(value)
	return strings.TrimSpace(string(runes[:limit])) + "..."
}

func quoteValidityText(validUntil string) string {
	validUntil = strings.TrimSpace(validUntil)
	if validUntil == "" {
		return ""
	}
	return validUntil
}

func (a *businessAPI) normalizeLines(ctx context.Context, input []BusinessLine, enforceStock bool) ([]BusinessLine, float64, string) {
	if len(input) == 0 {
		return nil, 0, "至少需要一条产品明细"
	}
	result := make([]BusinessLine, 0, len(input))
	total := 0.0
	for _, line := range input {
		var product Product
		if err := a.store.get(ctx, "products", strings.TrimSpace(line.ProductID), &product); err != nil {
			return nil, 0, "产品 " + line.ProductID + " 不存在"
		}
		if product.Status == "Inactive" {
			return nil, 0, "产品 " + product.Name + " 已停用"
		}
		if line.Quantity <= 0 {
			return nil, 0, "产品数量必须大于 0"
		}
		if line.UnitPrice < 0 {
			return nil, 0, "产品单价不能为负数"
		}
		if line.Discount < 0 || line.Discount > 100 {
			return nil, 0, "折扣必须在 0-100 之间"
		}
		if enforceStock && line.Quantity > float64(product.Stock) {
			return nil, 0, product.Name + " 数量超过当前库存"
		}
		line.ProductName = product.Name
		line.Amount = line.Quantity * line.UnitPrice * (1 - line.Discount/100)
		total += line.Amount
		result = append(result, line)
	}
	return result, total, ""
}

func (a *businessAPI) activeCustomerExists(ctx context.Context, name string) bool {
	items, err := listRecords[Customer](ctx, a.store, "accounts")
	if err != nil {
		return false
	}
	for _, item := range items {
		if !item.Archived && item.Name == strings.TrimSpace(name) {
			return true
		}
	}
	return false
}

func (a *businessAPI) quoteReferenced(ctx context.Context, id string) bool {
	items, _ := listRecords[Order](ctx, a.store, "orders")
	for _, item := range items {
		if item.Quote == id {
			return true
		}
	}
	return false
}

func (a *businessAPI) orderReferenced(ctx context.Context, id string) bool {
	items, _ := listRecords[Document](ctx, a.store, "documents")
	for _, item := range items {
		if item.Order == id {
			return true
		}
	}
	return false
}

func validDocumentType(value string) bool {
	return value == "PI" || value == "CI" || value == "PL" || value == "报关单" || value == "合同"
}

func validQuoteStatus(value string) bool {
	return value == "Draft" || value == "Delivered" || value == "Accepted" || value == "Rejected" || value == "Archived"
}

func validOrderStatus(value string) bool {
	return value == "Draft" || value == "Confirmed" || value == "Paid" || value == "Production" || value == "Shipped" || value == "Completed" || value == "Cancelled"
}

func validDocumentStatus(value string) bool {
	return value == "Draft" || value == "Review" || value == "Confirmed" || value == "Sent" || value == "Cancelled"
}

func sortCommerce[T any](items []T, field, direction string, money func(T) string, updated func(T) string) {
	ascending := direction == "asc"
	less := func(i, j int) bool { return updated(items[i]) < updated(items[j]) }
	if field == "value" {
		less = func(i, j int) bool { return moneyNumber(money(items[i])) < moneyNumber(money(items[j])) }
	}
	sort.SliceStable(items, func(i, j int) bool {
		if ascending {
			return less(i, j)
		}
		return less(j, i)
	})
}
