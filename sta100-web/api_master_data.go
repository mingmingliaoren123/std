package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (a *businessAPI) accountsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.accountsCollection(w, r)
		return
	}
	if parts[0] == "ocr" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		writeTODO(w, "TODO_OCR_SPEC", "名片识别接口已保留，需客户确认 OCR 引擎、字段和置信度规则", []string{"OCR 引擎", "名片样例", "字段映射", "置信度阈值"})
		return
	}
	if parts[0] == "export" {
		a.accountsExport(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "relations" {
		a.accountRelations(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "communications" {
		a.accountCommunications(w, r, parts[0])
		return
	}
	a.accountItem(w, r, parts[0])
}

func (a *businessAPI) accountsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := a.listCustomers(r.Context())
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		customerType := strings.TrimSpace(r.URL.Query().Get("type"))
		country := strings.TrimSpace(r.URL.Query().Get("country"))
		includeArchived := r.URL.Query().Get("include_archived") == "true"
		filtered := make([]Customer, 0, len(items))
		for _, item := range items {
			if item.Archived && !includeArchived {
				continue
			}
			if query != "" && !containsFold(item.ID+" "+item.Name+" "+item.Contact+" "+item.Email+" "+item.Phone+" "+item.Country, query) {
				continue
			}
			if customerType != "" && customerType != "all" && item.Type != customerType {
				continue
			}
			if country != "" && country != "all" && item.Country != country {
				continue
			}
			filtered = append(filtered, item)
		}
		sortCustomers(filtered, r.URL.Query().Get("sort"), r.URL.Query().Get("direction"))
		total := len(filtered)
		writeJSON(w, http.StatusOK, map[string]any{"items": filterPage(filtered, r), "total": total})
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var request Customer
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		if message := validateCustomer(request); message != "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_ACCOUNT", message)
			return
		}
		if a.customerNameExists(r.Context(), request.Name, "") {
			writeAPIError(w, http.StatusConflict, "ACCOUNT_NAME_EXISTS", "客户名称已存在")
			return
		}
		id, err := a.store.nextSequence(r.Context(), "accounts", "ACC", 4)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		request.ID, request.Updated, request.Archived = id, currentText(), false
		if request.Type == "" {
			request.Type = "Customer"
		}
		if request.Rating == "" {
			request.Rating = "Prospect"
		}
		if request.Owner == "" {
			request.Owner = requestOperator(r)
		}
		request.Orders, request.Total = 0, "EUR 0.00"
		if err := a.store.create(r.Context(), "accounts", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "create", "account", id, requestOperator(r), map[string]any{"name": request.Name})
		writeJSON(w, http.StatusCreated, request)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) accountItem(w http.ResponseWriter, r *http.Request, id string) {
	var item Customer
	if err := a.store.get(r.Context(), "accounts", id, &item); err != nil {
		writeBusinessError(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		a.applyCustomerAggregate(r.Context(), &item)
		writeJSON(w, http.StatusOK, item)
	case http.MethodPatch:
		if !allowMutation(w, r, http.MethodPatch) {
			return
		}
		var request Customer
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.ID = id
		if message := validateCustomer(request); message != "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_ACCOUNT", message)
			return
		}
		if a.customerNameExists(r.Context(), request.Name, id) {
			writeAPIError(w, http.StatusConflict, "ACCOUNT_NAME_EXISTS", "客户名称已存在")
			return
		}
		request.Orders, request.Total, request.Archived = item.Orders, item.Total, item.Archived
		request.Updated = currentText()
		if err := a.updateCustomerWithReferences(r.Context(), item, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "account", id, requestOperator(r), map[string]any{"name": request.Name})
		a.applyCustomerAggregate(r.Context(), &request)
		writeJSON(w, http.StatusOK, request)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		var delReq struct {
			HardDelete bool `json:"hardDelete"`
		}
		if r.ContentLength > 0 || strings.TrimSpace(r.Header.Get("Content-Type")) != "" {
			if err := decodeJSONBody(w, r, &delReq); err != nil {
				return
			}
		}
		if delReq.HardDelete {
			if err := a.store.deleteRecord(r.Context(), "accounts", id); err != nil {
				writeBusinessError(w, err)
				return
			}
			a.store.audit(r.Context(), "delete", "account", id, requestOperator(r), nil)
			writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id})
			return
		}
		item.Archived, item.Updated = true, currentText()
		if err := a.store.put(r.Context(), "accounts", id, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "archive", "account", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, map[string]any{"archived": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) updateCustomerWithReferences(ctx context.Context, previous, updated Customer) error {
	quotes, err := listRecords[Quote](ctx, a.store, "quotes")
	if err != nil {
		return err
	}
	orders, err := listRecords[Order](ctx, a.store, "orders")
	if err != nil {
		return err
	}
	documents, err := listRecords[Document](ctx, a.store, "documents")
	if err != nil {
		return err
	}
	tx, err := a.store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	put := func(kind, id string, value any) error {
		data, marshalErr := json.Marshal(value)
		if marshalErr != nil {
			return marshalErr
		}
		result, execErr := tx.ExecContext(ctx, `UPDATE records SET data_json=?,updated_at=? WHERE kind=? AND id=? AND deleted_at IS NULL`, data, timeNowUTC(), kind, id)
		if execErr != nil {
			return execErr
		}
		changed, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return rowsErr
		}
		if changed == 0 {
			return errRecordNotFound
		}
		return nil
	}
	if err := put("accounts", updated.ID, updated); err != nil {
		return err
	}
	if previous.Name != updated.Name {
		for _, item := range quotes {
			if item.Customer == previous.Name {
				item.Customer = updated.Name
				if err := put("quotes", item.ID, item); err != nil {
					return err
				}
			}
		}
		for _, item := range orders {
			if item.Customer == previous.Name {
				item.Customer = updated.Name
				if err := put("orders", item.ID, item); err != nil {
					return err
				}
			}
		}
		for _, item := range documents {
			if item.Customer == previous.Name {
				item.Customer = updated.Name
				if err := put("documents", item.ID, item); err != nil {
					return err
				}
			}
		}
	}
	return tx.Commit()
}

func (a *businessAPI) accountRelations(w http.ResponseWriter, r *http.Request, id string) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	var customer Customer
	if err := a.store.get(r.Context(), "accounts", id, &customer); err != nil {
		writeBusinessError(w, err)
		return
	}
	quotes, _ := listRecords[Quote](r.Context(), a.store, "quotes")
	orders, _ := listRecords[Order](r.Context(), a.store, "orders")
	documents, _ := listRecords[Document](r.Context(), a.store, "documents")
	relatedQuotes := make([]Quote, 0)
	relatedOrders := make([]Order, 0)
	relatedDocuments := make([]Document, 0)
	for _, item := range quotes {
		if item.Customer == customer.Name {
			relatedQuotes = append(relatedQuotes, item)
		}
	}
	for _, item := range orders {
		if item.Customer == customer.Name {
			relatedOrders = append(relatedOrders, item)
		}
	}
	for _, item := range documents {
		if item.Customer == customer.Name {
			relatedDocuments = append(relatedDocuments, item)
		}
	}
	communications, _ := a.listCustomerCommunications(r.Context(), id)
	writeJSON(w, http.StatusOK, map[string]any{"customer": customer, "quotes": relatedQuotes, "orders": relatedOrders, "documents": relatedDocuments, "communications": communications})
}

func (a *businessAPI) accountCommunications(w http.ResponseWriter, r *http.Request, customerID string) {
	var customer Customer
	if err := a.store.get(r.Context(), "accounts", customerID, &customer); err != nil {
		writeBusinessError(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := a.listCustomerCommunications(r.Context(), customerID)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		listResponse(w, items)
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		if customer.Archived {
			writeAPIError(w, http.StatusConflict, "ACCOUNT_ARCHIVED", "客户已归档，历史沟通记录可查询但不能继续新增")
			return
		}
		var request CustomerCommunication
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.Type = strings.TrimSpace(request.Type)
		request.Subject = strings.TrimSpace(request.Subject)
		request.Content = strings.TrimSpace(request.Content)
		request.Contact = strings.TrimSpace(request.Contact)
		request.OccurredAt = strings.TrimSpace(request.OccurredAt)
		if request.Type == "" || request.Content == "" || request.OccurredAt == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_COMMUNICATION", "沟通方式、沟通时间和沟通内容不能为空")
			return
		}
		if len([]rune(request.Subject)) > 200 || len([]rune(request.Content)) > 10000 || len([]rune(request.Contact)) > 200 {
			writeAPIError(w, http.StatusBadRequest, "INVALID_COMMUNICATION", "沟通记录字段超过长度限制")
			return
		}
		id, err := a.store.nextSequence(r.Context(), "customer_communications", "COMM", 8)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		request.ID = id
		request.CustomerID = customerID
		request.CreatedAt = timeNowUTC()
		request.CreatedBy = requestOperator(r)
		if err := a.store.create(r.Context(), "customer_communications", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "append", "customer_communication", id, request.CreatedBy, map[string]any{"customerId": customerID, "type": request.Type, "occurredAt": request.OccurredAt})
		writeJSON(w, http.StatusCreated, request)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "COMMUNICATION_IMMUTABLE", "历史沟通记录只能追加，不能修改或删除")
	}
}

func (a *businessAPI) listCustomerCommunications(ctx context.Context, customerID string) ([]CustomerCommunication, error) {
	items, err := listRecords[CustomerCommunication](ctx, a.store, "customer_communications")
	if err != nil {
		return nil, err
	}
	filtered := make([]CustomerCommunication, 0)
	for _, item := range items {
		if item.CustomerID == customerID {
			filtered = append(filtered, item)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].OccurredAt == filtered[j].OccurredAt {
			return filtered[i].CreatedAt > filtered[j].CreatedAt
		}
		return filtered[i].OccurredAt > filtered[j].OccurredAt
	})
	return filtered, nil
}

func (a *businessAPI) accountsExport(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	items, err := a.listCustomers(r.Context())
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	rows := [][]string{{"客户编号", "客户名称", "类型", "国家", "城市", "联系人", "电话", "邮箱", "订单数", "累计金额", "评级", "来源", "更新时间"}}
	for _, item := range items {
		if item.Archived {
			continue
		}
		rows = append(rows, []string{item.ID, item.Name, item.Type, item.Country, item.City, item.Contact, item.Phone, item.Email, strconv.Itoa(item.Orders), item.Total, item.Rating, item.Source, item.Updated})
	}
	a.store.audit(r.Context(), "export", "account", "*", requestOperator(r), map[string]any{"count": len(rows) - 1, "format": "xlsx"})
	writeXLSX(w, "sta100-accounts.xlsx", "客户", rows)
}

func (a *businessAPI) listCustomers(ctx context.Context) ([]Customer, error) {
	items, err := listRecords[Customer](ctx, a.store, "accounts")
	if err != nil {
		return nil, err
	}
	for index := range items {
		a.applyCustomerAggregate(ctx, &items[index])
	}
	return items, nil
}

func (a *businessAPI) applyCustomerAggregate(ctx context.Context, customer *Customer) {
	orders, err := listRecords[Order](ctx, a.store, "orders")
	if err != nil {
		return
	}
	count, total := 0, 0.0
	for _, order := range orders {
		if order.Customer == customer.Name && order.Status != "Cancelled" {
			count++
			total += moneyNumber(order.Value)
		}
	}
	customer.Orders, customer.Total = count, formatMoney(total, "EUR")
}

func (a *businessAPI) customerNameExists(ctx context.Context, name, exceptID string) bool {
	items, err := listRecords[Customer](ctx, a.store, "accounts")
	if err != nil {
		return false
	}
	for _, item := range items {
		if item.ID != exceptID && !item.Archived && strings.EqualFold(strings.TrimSpace(item.Name), strings.TrimSpace(name)) {
			return true
		}
	}
	return false
}

func validateCustomer(item Customer) string {
	if strings.TrimSpace(item.Name) == "" {
		return "客户名称不能为空"
	}
	if strings.TrimSpace(item.Country) == "" {
		return "账单国家不能为空"
	}
	if len([]rune(item.Name)) > 200 {
		return "客户名称不能超过 200 个字符"
	}
	return ""
}

func sortCustomers(items []Customer, field, direction string) {
	ascending := direction == "asc"
	less := func(i, j int) bool { return items[i].Updated < items[j].Updated }
	switch field {
	case "orders":
		less = func(i, j int) bool { return items[i].Orders < items[j].Orders }
	case "total":
		less = func(i, j int) bool { return moneyNumber(items[i].Total) < moneyNumber(items[j].Total) }
	case "name":
		less = func(i, j int) bool { return items[i].Name < items[j].Name }
	}
	sort.SliceStable(items, func(i, j int) bool {
		if ascending {
			return less(i, j)
		}
		return less(j, i)
	})
}

func (a *businessAPI) productsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.productsCollection(w, r)
		return
	}
	if parts[0] == "import" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		writeTODO(w, "TODO_PRODUCT_IMPORT_FORMAT", "产品批量导入接口已保留，需确认导入模板、重复处理和错误回执格式", []string{"产品导入模板", "必填字段", "重复编码处理", "错误回执格式"})
		return
	}
	a.productItem(w, r, parts[0])
}

func (a *businessAPI) productsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := listRecords[Product](r.Context(), a.store, "products")
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		category := r.URL.Query().Get("category")
		filtered := make([]Product, 0, len(items))
		for _, item := range items {
			if query != "" && !containsFold(item.ID+" "+item.Name+" "+item.Category+" "+item.HS, query) {
				continue
			}
			if category != "" && category != "all" && item.Category != category {
				continue
			}
			filtered = append(filtered, item)
		}
		listResponse(w, filtered)
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var item Product
		if err := decodeJSONBody(w, r, &item); err != nil {
			return
		}
		if message := validateProduct(item); message != "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_PRODUCT", message)
			return
		}
		item.ID = strings.ToUpper(strings.TrimSpace(item.ID))
		item.Updated = currentText()
		if item.Status == "" {
			item.Status = "Active"
		}
		if err := a.store.create(r.Context(), "products", item.ID, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "create", "product", item.ID, requestOperator(r), map[string]any{"name": item.Name})
		writeJSON(w, http.StatusCreated, item)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) productItem(w http.ResponseWriter, r *http.Request, id string) {
	var item Product
	if err := a.store.get(r.Context(), "products", id, &item); err != nil {
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
		var request Product
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.ID = id
		if message := validateProduct(request); message != "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_PRODUCT", message)
			return
		}
		request.Updated = currentText()
		if err := a.store.put(r.Context(), "products", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "product", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, request)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if a.productReferenced(r.Context(), id) {
			writeAPIError(w, http.StatusConflict, "PRODUCT_IN_USE", "产品已被报价、订单或单据引用，只能停用")
			return
		}
		item.Status, item.Updated = "Inactive", currentText()
		if err := a.store.put(r.Context(), "products", id, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "deactivate", "product", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, map[string]any{"deactivated": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func validateProduct(item Product) string {
	if strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Name) == "" {
		return "产品编码和名称不能为空"
	}
	if strings.TrimSpace(item.HS) == "" {
		return "HS CODE 不能为空"
	}
	if item.Stock < 0 {
		return "库存不能为负数"
	}
	if moneyNumber(item.Price) < 0 {
		return "价格不能为负数"
	}
	return ""
}

func (a *businessAPI) productReferenced(ctx context.Context, id string) bool {
	quotes, _ := listRecords[Quote](ctx, a.store, "quotes")
	orders, _ := listRecords[Order](ctx, a.store, "orders")
	documents, _ := listRecords[Document](ctx, a.store, "documents")
	for _, quote := range quotes {
		for _, line := range quote.Lines {
			if line.ProductID == id {
				return true
			}
		}
	}
	for _, order := range orders {
		for _, line := range order.Lines {
			if line.ProductID == id {
				return true
			}
		}
	}
	for _, document := range documents {
		for _, line := range document.Lines {
			if line.ProductID == id {
				return true
			}
		}
	}
	return false
}

func (a *businessAPI) suppliersRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.suppliersCollection(w, r)
		return
	}
	if parts[0] == "export" {
		a.suppliersExport(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "communications" {
		a.supplierCommunications(w, r, parts[0])
		return
	}
	a.supplierItem(w, r, parts[0])
}

func (a *businessAPI) suppliersExport(w http.ResponseWriter, r *http.Request) {
	if !allowMethod(w, r, http.MethodGet) {
		return
	}
	items, err := listRecords[Supplier](r.Context(), a.store, "suppliers")
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	rows := [][]string{{"供应商编号", "公司", "电话", "联系人", "邮件", "产品", "规格", "报价", "来源", "备注", "更新时间"}}
	for _, item := range items {
		if !item.Archived {
			rows = append(rows, []string{item.ID, item.Company, item.Phone, item.Contact, item.Email, item.Product, item.Specification, item.Quote, item.Source, item.Notes, item.Updated})
		}
	}
	a.store.audit(r.Context(), "export", "supplier", "*", requestOperator(r), map[string]any{"count": len(rows) - 1, "format": "xlsx"})
	writeXLSX(w, "sta100-suppliers.xlsx", "供应商", rows)
}

func (a *businessAPI) suppliersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := listRecords[Supplier](r.Context(), a.store, "suppliers")
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		filtered := make([]Supplier, 0, len(items))
		for _, item := range items {
			if item.Archived {
				continue
			}
			if query != "" && !containsFold(item.ID+" "+item.Company+" "+item.Contact+" "+item.Product+" "+item.Source, query) {
				continue
			}
			filtered = append(filtered, item)
		}
		listResponse(w, filtered)
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var item Supplier
		if err := decodeJSONBody(w, r, &item); err != nil {
			return
		}
		if strings.TrimSpace(item.Company) == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_SUPPLIER", "供应商公司名称不能为空")
			return
		}
		id, err := a.store.nextSequence(r.Context(), "suppliers", "SUP", 4)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		item.ID, item.Updated = id, currentText()
		if err := a.store.create(r.Context(), "suppliers", id, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "create", "supplier", id, requestOperator(r), nil)
		writeJSON(w, http.StatusCreated, item)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) supplierItem(w http.ResponseWriter, r *http.Request, id string) {
	var item Supplier
	if err := a.store.get(r.Context(), "suppliers", id, &item); err != nil {
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
		var request Supplier
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		if strings.TrimSpace(request.Company) == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_SUPPLIER", "供应商公司名称不能为空")
			return
		}
		request.ID, request.Updated = id, currentText()
		if err := a.store.put(r.Context(), "suppliers", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "supplier", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, request)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		item.Archived, item.Updated = true, currentText()
		if err := a.store.put(r.Context(), "suppliers", id, item); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "archive", "supplier", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, map[string]any{"archived": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

// leadsRouter 线索路由
func (a *businessAPI) leadsRouter(w http.ResponseWriter, r *http.Request, parts []string) {
	if len(parts) == 0 || parts[0] == "" {
		a.leadsCollection(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "communications" {
		a.leadCommunications(w, r, parts[0])
		return
	}
	if len(parts) == 2 && parts[1] == "convert" {
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		a.leadConvert(w, r, parts[0])
		return
	}
	a.leadItem(w, r, parts[0])
}

func (a *businessAPI) leadsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := a.listLeads(r.Context())
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		query := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
		leadType := strings.TrimSpace(r.URL.Query().Get("type"))
		country := strings.TrimSpace(r.URL.Query().Get("country"))
		status := strings.TrimSpace(r.URL.Query().Get("status"))
		includeArchived := r.URL.Query().Get("include_archived") == "true"
		filtered := make([]Lead, 0, len(items))
		for _, item := range items {
			if item.Archived && !includeArchived {
				continue
			}
			if query != "" && !containsFold(item.ID+" "+item.Name+" "+item.Contact+" "+item.Email+" "+item.Phone+" "+item.Country, query) {
				continue
			}
			if leadType != "" && leadType != "all" && item.Type != leadType {
				continue
			}
			if country != "" && country != "all" && item.Country != country {
				continue
			}
			if status != "" && status != "all" {
				if status == "converted" && !item.Converted {
					continue
				}
				if (status == "pending" || status == "followup") && item.Converted {
					continue
				}
			}
			filtered = append(filtered, item)
		}
		sortLeads(filtered, r.URL.Query().Get("sort"), r.URL.Query().Get("direction"))
		total := len(filtered)
		writeJSON(w, http.StatusOK, map[string]any{"items": filterPage(filtered, r), "total": total})
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var request Lead
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.ID = request.ID
		if request.ID == "" {
			id, err := a.store.nextSequence(r.Context(), "leads", "LEAD", 8)
			if err != nil {
				writeBusinessError(w, err)
				return
			}
			request.ID = id
		}
		request.Name = strings.TrimSpace(request.Name)
		if request.Name == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_LEAD", "线索客户名称不能为空")
			return
		}
		request.Type = strings.TrimSpace(request.Type)
		request.Country = strings.TrimSpace(request.Country)
		request.City = strings.TrimSpace(request.City)
		request.Contact = strings.TrimSpace(request.Contact)
		request.Phone = strings.TrimSpace(request.Phone)
		request.Email = strings.TrimSpace(request.Email)
		request.Website = strings.TrimSpace(request.Website)
		request.Address = strings.TrimSpace(request.Address)
		request.Business = strings.TrimSpace(request.Business)
		request.Source = strings.TrimSpace(request.Source)
		request.SourceUrl = strings.TrimSpace(request.SourceUrl)
		request.Reason = strings.TrimSpace(request.Reason)
		request.CreatedAt = timeNowUTC()
		request.Updated = timeNowUTC()
		if err := a.store.create(r.Context(), "leads", request.ID, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "create", "lead", request.ID, requestOperator(r), map[string]any{"name": request.Name})
		writeJSON(w, http.StatusCreated, request)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) listLeads(ctx context.Context) ([]Lead, error) {
	return listRecords[Lead](ctx, a.store, "leads")
}

func sortLeads(items []Lead, field, direction string) {
	desc := strings.ToLower(strings.TrimSpace(direction)) != "asc"
	field = strings.ToLower(strings.TrimSpace(field))
	if field == "" {
		field = "updated"
	}
	sort.SliceStable(items, func(i, j int) bool {
		less := false
		switch field {
		case "score":
			less = items[i].Score < items[j].Score
		case "created", "createdat", "created_at":
			less = items[i].CreatedAt < items[j].CreatedAt
		default:
			less = items[i].Updated < items[j].Updated
		}
		if desc {
			return !less
		}
		return less
	})
}

func (a *businessAPI) leadItem(w http.ResponseWriter, r *http.Request, id string) {
	var item Lead
	if err := a.store.get(r.Context(), "leads", id, &item); err != nil {
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
		var request Lead
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.ID = id
		request.Name = strings.TrimSpace(request.Name)
		if request.Name == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_LEAD", "线索客户名称不能为空")
			return
		}
		request.Type = strings.TrimSpace(request.Type)
		request.Country = strings.TrimSpace(request.Country)
		request.City = strings.TrimSpace(request.City)
		request.Contact = strings.TrimSpace(request.Contact)
		request.Phone = strings.TrimSpace(request.Phone)
		request.Email = strings.TrimSpace(request.Email)
		request.Website = strings.TrimSpace(request.Website)
		request.Address = strings.TrimSpace(request.Address)
		request.Business = strings.TrimSpace(request.Business)
		request.Source = strings.TrimSpace(request.Source)
		request.SourceUrl = strings.TrimSpace(request.SourceUrl)
		request.Reason = strings.TrimSpace(request.Reason)
		request.Updated = timeNowUTC()
		request.Converted = item.Converted
		request.ConvertedAt = item.ConvertedAt
		request.ConvertedID = item.ConvertedID
		request.CreatedAt = item.CreatedAt
		if err := a.store.put(r.Context(), "leads", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "update", "lead", id, requestOperator(r), map[string]any{"name": request.Name})
		writeJSON(w, http.StatusOK, request)
	case http.MethodDelete:
		if !allowMutation(w, r, http.MethodDelete) {
			return
		}
		if err := a.store.deleteRecord(r.Context(), "leads", id); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "delete", "lead", id, requestOperator(r), nil)
		writeJSON(w, http.StatusOK, map[string]any{"deleted": true, "id": id})
	default:
		w.Header().Set("Allow", "GET, PATCH, DELETE")
		writeAPIError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "请求方法不受支持")
	}
}

func (a *businessAPI) leadConvert(w http.ResponseWriter, r *http.Request, leadID string) {
	var lead Lead
	if err := a.store.get(r.Context(), "leads", leadID, &lead); err != nil {
		writeBusinessError(w, err)
		return
	}
	if lead.Converted {
		writeAPIError(w, http.StatusConflict, "LEAD_ALREADY_CONVERTED", "该线索已转化为客户")
		return
	}
	if lead.Archived {
		writeAPIError(w, http.StatusConflict, "LEAD_ARCHIVED", "已归档的线索不能转化为客户")
		return
	}
	var req struct {
		KeepSource *bool `json:"keepSource"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if r.Body != nil {
		_ = r.Body.Close()
	}
	keepSource := true
	if req.KeepSource != nil {
		keepSource = *req.KeepSource
	}
	customerID, err := a.store.nextSequence(r.Context(), "accounts", "ACC", 8)
	if err != nil {
		writeBusinessError(w, err)
		return
	}
	customer := Customer{
		ID:          customerID,
		Name:        lead.Name,
		Type:        lead.Type,
		Country:     lead.Country,
		City:        lead.City,
		Contact:     lead.Contact,
		Phone:       lead.Phone,
		Email:       lead.Email,
		Website:     lead.Website,
		Source:      "应用发现线索",
		Description: lead.Business,
		Rating:      "Prospect",
		Orders:      0,
		Total:       "0",
		Updated:     timeNowUTC(),
		Archived:    false,
	}
	if err := a.store.create(r.Context(), "accounts", customerID, customer); err != nil {
		writeBusinessError(w, err)
		return
	}
	if err := a.copyLeadCommunicationsToCustomer(r.Context(), leadID, customerID); err != nil {
		writeBusinessError(w, err)
		return
	}
	if keepSource {
		lead.Converted = true
		lead.ConvertedAt = timeNowUTC()
		lead.ConvertedID = customerID
		lead.Updated = timeNowUTC()
		if err := a.store.put(r.Context(), "leads", leadID, lead); err != nil {
			writeBusinessError(w, err)
			return
		}
	} else {
		if err := a.store.deleteRecord(r.Context(), "leads", leadID); err != nil {
			writeBusinessError(w, err)
			return
		}
	}
	a.store.audit(r.Context(), "convert", "lead->customer", leadID, requestOperator(r), map[string]any{"leadId": leadID, "customerId": customerID, "keepSource": keepSource})
	writeJSON(w, http.StatusCreated, map[string]any{"lead": lead, "customer": customer, "keepSource": keepSource})
}

func (a *businessAPI) leadCommunications(w http.ResponseWriter, r *http.Request, leadID string) {
	var lead Lead
	if err := a.store.get(r.Context(), "leads", leadID, &lead); err != nil {
		writeBusinessError(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := a.listLeadCommunications(r.Context(), leadID)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		listResponse(w, items)
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var request LeadCommunication
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.Type = strings.TrimSpace(request.Type)
		request.Subject = strings.TrimSpace(request.Subject)
		request.Content = strings.TrimSpace(request.Content)
		request.Contact = strings.TrimSpace(request.Contact)
		request.OccurredAt = strings.TrimSpace(request.OccurredAt)
		if request.Type == "" || request.Content == "" || request.OccurredAt == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_COMMUNICATION", "沟通方式、沟通时间和沟通内容不能为空")
			return
		}
		if len([]rune(request.Subject)) > 200 || len([]rune(request.Content)) > 10000 || len([]rune(request.Contact)) > 200 {
			writeAPIError(w, http.StatusBadRequest, "INVALID_COMMUNICATION", "沟通记录字段超过长度限制")
			return
		}
		id, err := a.store.nextSequence(r.Context(), "lead_communications", "COMM", 8)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		request.ID = id
		request.LeadID = leadID
		request.CreatedAt = timeNowUTC()
		request.CreatedBy = requestOperator(r)
		if err := a.store.create(r.Context(), "lead_communications", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "append", "lead_communication", id, request.CreatedBy, map[string]any{"leadId": leadID, "type": request.Type, "occurredAt": request.OccurredAt})
		writeJSON(w, http.StatusCreated, request)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "COMMUNICATION_IMMUTABLE", "历史沟通记录只能追加，不能修改或删除")
	}
}

func (a *businessAPI) listLeadCommunications(ctx context.Context, leadID string) ([]LeadCommunication, error) {
	items, err := listRecords[LeadCommunication](ctx, a.store, "lead_communications")
	if err != nil {
		return nil, err
	}
	filtered := make([]LeadCommunication, 0)
	for _, item := range items {
		if item.LeadID == leadID {
			filtered = append(filtered, item)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].OccurredAt == filtered[j].OccurredAt {
			return filtered[i].CreatedAt > filtered[j].CreatedAt
		}
		return filtered[i].OccurredAt > filtered[j].OccurredAt
	})
	return filtered, nil
}

func (a *businessAPI) copyLeadCommunicationsToCustomer(ctx context.Context, leadID, customerID string) error {
	items, err := a.listLeadCommunications(ctx, leadID)
	if err != nil {
		return err
	}
	for _, item := range items {
		id, err := a.store.nextSequence(ctx, "customer_communications", "COMM", 8)
		if err != nil {
			return err
		}
		copy := CustomerCommunication{
			ID:         id,
			CustomerID: customerID,
			Type:       item.Type,
			Subject:    item.Subject,
			Content:    item.Content,
			Contact:    item.Contact,
			OccurredAt: item.OccurredAt,
			CreatedAt:  item.CreatedAt,
			CreatedBy:  item.CreatedBy,
		}
		if err := a.store.create(ctx, "customer_communications", id, copy); err != nil {
			return err
		}
	}
	return nil
}

// supplierCommunications 供应商跟进记录
func (a *businessAPI) supplierCommunications(w http.ResponseWriter, r *http.Request, supplierID string) {
	var supplier Supplier
	if err := a.store.get(r.Context(), "suppliers", supplierID, &supplier); err != nil {
		writeBusinessError(w, err)
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := a.listSupplierCommunications(r.Context(), supplierID)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		listResponse(w, items)
	case http.MethodPost:
		if !allowMutation(w, r, http.MethodPost) {
			return
		}
		var request SupplierCommunication
		if err := decodeJSONBody(w, r, &request); err != nil {
			return
		}
		request.Type = strings.TrimSpace(request.Type)
		request.Subject = strings.TrimSpace(request.Subject)
		request.Content = strings.TrimSpace(request.Content)
		request.Contact = strings.TrimSpace(request.Contact)
		request.OccurredAt = strings.TrimSpace(request.OccurredAt)
		if request.Type == "" || request.Content == "" || request.OccurredAt == "" {
			writeAPIError(w, http.StatusBadRequest, "INVALID_COMMUNICATION", "沟通方式、沟通时间和沟通内容不能为空")
			return
		}
		if len([]rune(request.Subject)) > 200 || len([]rune(request.Content)) > 10000 || len([]rune(request.Contact)) > 200 {
			writeAPIError(w, http.StatusBadRequest, "INVALID_COMMUNICATION", "沟通记录字段超过长度限制")
			return
		}
		id, err := a.store.nextSequence(r.Context(), "supplier_communications", "SPC", 8)
		if err != nil {
			writeBusinessError(w, err)
			return
		}
		request.ID = id
		request.SupplierID = supplierID
		request.CreatedAt = timeNowUTC()
		request.CreatedBy = requestOperator(r)
		if err := a.store.create(r.Context(), "supplier_communications", id, request); err != nil {
			writeBusinessError(w, err)
			return
		}
		a.store.audit(r.Context(), "append", "supplier_communication", id, request.CreatedBy, map[string]any{"supplierId": supplierID, "type": request.Type, "occurredAt": request.OccurredAt})
		writeJSON(w, http.StatusCreated, request)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeAPIError(w, http.StatusMethodNotAllowed, "COMMUNICATION_IMMUTABLE", "历史沟通记录只能追加，不能修改或删除")
	}
}

func (a *businessAPI) listSupplierCommunications(ctx context.Context, supplierID string) ([]SupplierCommunication, error) {
	items, err := listRecords[SupplierCommunication](ctx, a.store, "supplier_communications")
	if err != nil {
		return nil, err
	}
	filtered := make([]SupplierCommunication, 0)
	for _, item := range items {
		if item.SupplierID == supplierID {
			filtered = append(filtered, item)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].OccurredAt == filtered[j].OccurredAt {
			return filtered[i].CreatedAt > filtered[j].CreatedAt
		}
		return filtered[i].OccurredAt > filtered[j].OccurredAt
	})
	return filtered, nil
}
