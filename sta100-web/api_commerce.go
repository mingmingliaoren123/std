package main

import (
	"context"
	"errors"
	"net/http"
	"sort"
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
	if item.Status != "Draft" {
		writeAPIError(w, http.StatusConflict, "QUOTE_STATUS_INVALID", "只有草稿报价可以发送")
		return
	}
	writeTODO(w, "TODO_QUOTE_DELIVERY", "报价发送需要正式模板、PDF 规则、邮件样式和发件通道，当前不会仅修改状态冒充发送成功", []string{"报价模板", "PDF 版式", "邮件主题和正文", "发件账户", "收件人规则", "发送失败重试规则"})
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
	writeTODO(w, "TODO_QUOTE_PDF", "报价 PDF 需在客户确认正式模板、纸张、语言和签章规则后生成", []string{"报价模板", "纸张尺寸", "语言", "Logo 和签章", "文件命名规则"})
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
	total += item.Freight + item.Tax
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
	a.orderItem(w, r, parts[0])
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
		request.Types = []string{"PI", "CI", "PL", "报关单"}
	}
	for _, documentType := range request.Types {
		if !validDocumentType(documentType) {
			writeAPIError(w, http.StatusBadRequest, "INVALID_DOCUMENT_TYPE", "单据类型仅支持 PI、CI、PL、报关单")
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
			writeAPIError(w, http.StatusBadRequest, "INVALID_DOCUMENT_TYPE", "单据类型仅支持 PI、CI、PL、报关单")
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
		return Document{}, errors.New("单据类型仅支持 PI、CI、PL、报关单")
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
	writeTODO(w, "TODO_DOCUMENT_EXPORT", "正式单据文件需在对应模板、字段映射和输出格式确认后生成", []string{"PI/CI/PL/报关单模板", "字段映射", "输出格式", "语言", "签章规则", "文件命名规则"})
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
	return value == "PI" || value == "CI" || value == "PL" || value == "报关单"
}

func validQuoteStatus(value string) bool {
	return value == "Draft" || value == "Delivered" || value == "Accepted" || value == "Rejected" || value == "Archived"
}

func validOrderStatus(value string) bool {
	return value == "Draft" || value == "Confirmed" || value == "Production" || value == "Shipped" || value == "Completed" || value == "Cancelled"
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
