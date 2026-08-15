package main

import (
	"context"
	"fmt"
)

func (s *businessStore) seed(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM records`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	customers := []Customer{
		{ID: "ACC-0001", Name: "VeloTrade GmbH", Type: "Distributor", Country: "德国", Contact: "Anna Keller", Phone: "+49 30 555 0188", Email: "anna@velotrade.example", Website: "https://velotrade.example", Owner: "Donald", Rating: "Active", Source: "展会", Updated: "2026-08-10 09:42"},
		{ID: "ACC-0002", Name: "Nordic Cycle AB", Type: "Importer", Country: "瑞典", Contact: "Erik Lund", Phone: "+46 8 410 2250", Email: "erik@nordiccycle.example", Website: "https://nordiccycle.example", Owner: "Donald", Rating: "Prospect", Source: "朋友介绍", Updated: "2026-08-09 16:20"},
		{ID: "ACC-0003", Name: "SIM Sp. z o.o.", Type: "Customer", Country: "波兰", Contact: "Marek Nowak", Phone: "+48 34 310 2260", Email: "sales@sim.example", Website: "https://sim.example", Owner: "Donald", Rating: "Active", Source: "电话", Updated: "2026-08-08 11:15"},
		{ID: "ACC-0004", Name: "Ciclo Iberia S.L.", Type: "Reseller", Country: "西班牙", Contact: "Lucia Martin", Phone: "+34 91 778 2301", Email: "lucia@cicloiberia.example", Website: "https://cicloiberia.example", Owner: "Donald", Rating: "Prospect", Source: "拜访", Updated: "2026-08-07 14:02"},
		{ID: "ACC-0005", Name: "Alpine Motion SAS", Type: "Integrator", Country: "法国", Contact: "Louis Bernard", Phone: "+33 1 8420 1189", Email: "louis@alpine.example", Website: "https://alpine.example", Owner: "Donald", Rating: "Active", Source: "互联网线索", Updated: "2026-08-06 10:30"},
	}
	products := []Product{
		{ID: "STA-100-EU", Name: "STA-100 私有智能体设备", Category: "智能设备", Manufacturer: "STRATRONIX", Price: "EUR 159.00", Stock: 426, HS: "8471504090", Status: "Active", Description: "4G+32G，Linux OS，预装 OpenClaw Framework，支持客户私有数据。", Updated: "2026-08-10 08:20"},
		{ID: "PM-DUAL-01", Name: "双边功率计套装", Category: "智能骑行", Manufacturer: "STRATRONIX", Price: "EUR 438.00", Stock: 68, HS: "9029209000", Status: "Active", Description: "支持 ANT+ / BLE，提供曲柄兼容信息和出厂校准记录。", Updated: "2026-08-09 08:20"},
		{ID: "GPS-PRO-02", Name: "骑行码表 Pro", Category: "智能骑行", Manufacturer: "STRATRONIX", Price: "EUR 219.00", Stock: 112, HS: "8526919090", Status: "Active", Description: "多卫星定位、路线导航与训练数据同步。", Updated: "2026-08-08 08:20"},
		{ID: "EBK-CITY-03", Name: "城市 E-bike 方案", Category: "整车方案", Manufacturer: "STRATRONIX", Price: "EUR 1,280.00", Stock: 19, HS: "8711601000", Status: "Review", Description: "面向欧洲城市通勤渠道的电助力整车组合方案。", Updated: "2026-08-07 08:20"},
	}
	quotes := []Quote{
		{ID: "QUO-2026-0188", Subject: "STA-100 首批设备报价", Customer: "VeloTrade GmbH", Value: "EUR 34,800.00", Valid: "2026-09-08", Status: "Delivered", Products: "STA-100 私有智能体设备 x 200", Owner: "Donald", Currency: "EUR", Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "STA-100 私有智能体设备", Quantity: 200, UnitPrice: 174, Amount: 34800}}, Updated: "2026-08-10 09:40"},
		{ID: "QUO-2026-0187", Subject: "智能骑行组件组合", Customer: "Nordic Cycle AB", Value: "EUR 18,450.00", Valid: "2026-09-02", Status: "Draft", Products: "双边功率计套装 x 50", Owner: "Donald", Currency: "EUR", Lines: []BusinessLine{{ProductID: "PM-DUAL-01", ProductName: "双边功率计套装", Quantity: 50, UnitPrice: 369, Amount: 18450}}, Updated: "2026-08-09 09:40"},
		{ID: "QUO-2026-0185", Subject: "波兰渠道补货报价", Customer: "SIM Sp. z o.o.", Value: "EUR 22,100.00", Valid: "2026-08-26", Status: "Accepted", Products: "STA-100 私有智能体设备 x 120", Owner: "Donald", Currency: "EUR", Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "STA-100 私有智能体设备", Quantity: 120, UnitPrice: 184.1666667, Amount: 22100}}, Updated: "2026-08-08 09:40"},
		{ID: "QUO-2026-0179", Subject: "西班牙测试订单报价", Customer: "Ciclo Iberia S.L.", Value: "EUR 5,760.00", Valid: "2026-08-15", Status: "Rejected", Products: "STA-100 私有智能体设备 x 32", Owner: "Donald", Currency: "EUR", Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "STA-100 私有智能体设备", Quantity: 32, UnitPrice: 180, Amount: 5760}}, Updated: "2026-08-07 09:40"},
	}
	orders := []Order{
		{ID: "SO-2026-0106", Customer: "SIM Sp. z o.o.", Quote: "QUO-2026-0185", Products: "STA-100 私有智能体设备 x 120", Value: "EUR 22,100.00", Currency: "EUR", Status: "Production", Delivery: "2026-09-10", Progress: 55, Lines: quotes[2].Lines, Updated: "2026-08-10 10:20"},
		{ID: "SO-2026-0105", Customer: "VeloTrade GmbH", Quote: "QUO-2026-0174", Products: "STA-100 私有智能体设备 x 80", Value: "EUR 14,900.00", Currency: "EUR", Status: "Confirmed", Delivery: "2026-08-30", Progress: 30, Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "STA-100 私有智能体设备", Quantity: 80, UnitPrice: 186.25, Amount: 14900}}, Updated: "2026-08-09 10:20"},
		{ID: "SO-2026-0102", Customer: "Alpine Motion SAS", Quote: "QUO-2026-0162", Products: "双边功率计套装 x 40", Value: "EUR 31,600.00", Currency: "EUR", Status: "Shipped", Delivery: "2026-08-18", Progress: 82, Lines: []BusinessLine{{ProductID: "PM-DUAL-01", ProductName: "双边功率计套装", Quantity: 40, UnitPrice: 790, Amount: 31600}}, Updated: "2026-08-08 10:20"},
		{ID: "SO-2026-0098", Customer: "Nordic Cycle AB", Quote: "QUO-2026-0151", Products: "STA-100 私有智能体设备 x 25", Value: "EUR 4,650.00", Currency: "EUR", Status: "Completed", Delivery: "2026-08-03", Progress: 100, Lines: []BusinessLine{{ProductID: "STA-100-EU", ProductName: "STA-100 私有智能体设备", Quantity: 25, UnitPrice: 186, Amount: 4650}}, Updated: "2026-08-07 10:20"},
	}
	documents := []Document{
		{ID: "PI-20260810-003", Type: "PI", Customer: "SIM Sp. z o.o.", Order: "SO-2026-0106", Template: "STRATRONIX 标准 PI v3", Language: "英文", Status: "Draft", Value: orders[0].Value, Lines: orders[0].Lines, Updated: "2026-08-10 10:22"},
		{ID: "CI-20260809-012", Type: "CI", Customer: "Alpine Motion SAS", Order: "SO-2026-0102", Template: "欧盟商业发票 v2", Language: "英文", Status: "Confirmed", Value: orders[2].Value, Lines: orders[2].Lines, Updated: "2026-08-09 17:45"},
		{ID: "PL-20260809-009", Type: "PL", Customer: "Alpine Motion SAS", Order: "SO-2026-0102", Template: "标准装箱单 v2", Language: "英文", Status: "Confirmed", Value: orders[2].Value, Lines: orders[2].Lines, Updated: "2026-08-09 17:41"},
		{ID: "CD-20260808-006", Type: "报关单", Customer: "VeloTrade GmbH", Order: "SO-2026-0105", Template: "出口报关单 v1", Language: "中文 / 英文双语", Status: "Review", Value: orders[1].Value, Lines: orders[1].Lines, Updated: "2026-08-08 15:12"},
	}
	suppliers := []Supplier{
		{ID: "SUP-0001", Company: "苏州骑行动力科技", Phone: "0512-6688 2100", Contact: "周敏", Email: "sales@suzhou-cycle.example", Product: "E-bike 电池", Specification: "36V / 48V，定制 BMS", Quote: "EUR 82 / 组", Notes: "可提供 CE、UN38.3 资料", Source: "展会", Updated: "2026-08-10 08:50"},
		{ID: "SUP-0002", Company: "东莞骑迹智能装备", Phone: "0769-8822 1177", Contact: "陈杰", Email: "business@qiji-cycle.example", Product: "码表、功率计", Specification: "ANT+ / BLE，OEM 包装", Quote: "EUR 38 / 件起", Notes: "支持小批量 ODM", Source: "朋友介绍", Updated: "2026-08-09 15:20"},
		{ID: "SUP-0003", Company: "厦门轻量运动用品", Phone: "0592-6200 3318", Contact: "林晓", Email: "export@lightmotion.example", Product: "头盔、骑行装备", Specification: "MIPS，EN1078", Quote: "EUR 19 / 件起", Notes: "认证文件待复核", Source: "电话", Updated: "2026-08-08 10:05"},
	}
	privateFiles := []PrivateFile{
		{ID: "FILE-0001", Name: "STA100_Product_Manual_EN.pdf", Category: "产品手册", Tags: []string{"STA-100", "英文"}, Size: "6.8 MB", Source: "出厂通用资料", Status: "Indexed", Updated: "2026-08-10 08:20"},
		{ID: "FILE-0002", Name: "VeloTrade_Distribution_Agreement.docx", Category: "合同", Tags: []string{"客户私有", "德国"}, Size: "1.2 MB", Source: "客户上传", Status: "Indexed", Updated: "2026-08-09 16:10"},
		{ID: "FILE-0003", Name: "EU_Battery_Regulation_2026.pdf", Category: "法规", Tags: []string{"欧盟", "E-bike"}, Size: "3.1 MB", Source: "通用知识库", Status: "Indexed", Updated: "2026-08-09 12:45"},
		{ID: "FILE-0004", Name: "Shimano_Compatibility_List.xlsx", Category: "产品资料", Tags: []string{"兼容", "组件"}, Size: "846 KB", Source: "客户上传", Status: "Review", Updated: "2026-08-08 18:33"},
	}
	news := []NewsItem{
		{ID: "NEWS-0001", Category: "欧洲市场", Title: "欧洲自行车产业进入补库存周期，渠道更关注小批量和快速交付", Summary: "多家欧洲经销商在 2026 年下半年调整采购节奏，订单结构从大批量预采转向小批量、多批次。", Source: "Bike Europe", Time: "2026-08-10 09:10", Relevance: "96%"},
		{ID: "NEWS-0002", Category: "法规", Title: "欧盟更新电池尽职调查实施指引，E-bike 供应链资料需同步准备", Summary: "新指引强化材料来源、碳足迹和供应链证明要求。", Source: "EUR-Lex", Time: "2026-08-10 08:35", Relevance: "93%"},
		{ID: "NEWS-0003", Category: "智能骑行", Title: "无线电子变速与功率数据融合成为高端整车配置趋势", Summary: "整车厂正在把兼容性和训练数据完整度作为高端产品卖点。", Source: "Cycling Industry News", Time: "2026-08-09 17:40", Relevance: "89%"},
		{ID: "NEWS-0004", Category: "渠道", Title: "北欧经销商加快建设线上线下一体的维修服务网络", Summary: "服务能力和备件响应速度正在影响品牌进入门槛。", Source: "Nordic Cycling", Time: "2026-08-09 15:20", Relevance: "84%"},
		{ID: "NEWS-0005", Category: "产品", Title: "欧洲城市通勤市场对轻量化 E-bike 的关注持续上升", Summary: "重量、可维护性和电池合规成为渠道选品主要指标。", Source: "E-bike News", Time: "2026-08-09 11:05", Relevance: "82%"},
	}
	recommendations := []Recommendation{
		{ID: "REC-0001", Title: "Eurobike 2026 展商名录新增 86 家欧洲采购商", Desc: "与欧洲经销商、整车进口商关注条件匹配，可进一步生成客户候选清单。", Source: "Eurobike", Type: "展会情报", Time: "12 分钟前"},
		{ID: "REC-0002", Title: "欧盟电池法规尽职调查条款进入新执行阶段", Desc: "可能影响 E-bike 电池产品资料与供应商声明，建议同步检查现有模板。", Source: "EUR-Lex", Type: "法规", Time: "36 分钟前"},
		{ID: "REC-0003", Title: "德国通勤 E-bike 价格带出现结构性变化", Desc: "中端价格段热度上升，适合评估产品组合和报价策略。", Source: "Bike Europe", Type: "市场趋势", Time: "1 小时前"},
	}
	jobs := []Job{
		{ID: "JOB-RECOMMEND", Name: "每日推荐更新", Kind: "recommendations", Description: "定时汇总本地推荐缓存，后续接入客户确认的数据源。", AgentID: "sta100-coordinator", Prompt: "根据当前关注国家、主题和本地业务数据生成推荐摘要。", Schedule: "每 60 分钟", Enabled: true, BuiltIn: true, Status: "WaitingSource", UpdatedAt: "2026-08-10 08:00"},
		{ID: "JOB-WEEKLY", Name: "智能体周报", Kind: "weekly_report", Description: "读取本机 Agent 会话、Token 使用和业务审计日志，生成周报草稿。", AgentID: "sta100-coordinator", Prompt: "汇总最近 7 天 STA-100 智能体使用情况、关键业务操作和待跟进事项。", Schedule: "每周", Enabled: true, BuiltIn: true, Status: "Ready", UpdatedAt: "2026-08-10 08:00"},
		{ID: "JOB-NEWS", Name: "行业新闻更新", Kind: "news", Description: "按客户确认的新闻来源和频率抓取行业新闻，当前来源规则待确认。", AgentID: "market-analyzer", Prompt: "围绕骑行行业、欧洲渠道、法规和产品趋势整理新闻候选。", Schedule: "每 60 分钟", Enabled: true, BuiltIn: true, Status: "WaitingSource", UpdatedAt: "2026-08-10 08:00"},
		{ID: "JOB-INDEX", Name: "数据索引维护", Kind: "index", Description: "扫描本机私有文件元数据，正式正文解析和向量索引等待原始数据格式。", AgentID: "sta100-knowledge", Prompt: "检查本地私有文件是否需要解析、分类、去重和索引。", Schedule: "每天", Enabled: true, BuiltIn: true, Status: "WaitingData", UpdatedAt: "2026-08-10 08:00"},
	}
	plugins := []Plugin{
		{ID: "wechat", Name: "微信", Enabled: false, Status: "Unbound", Capabilities: []string{}, UpdatedAt: "2026-08-10 08:00"},
		{ID: "feishu", Name: "飞书", Enabled: false, Status: "Unbound", Capabilities: []string{}, UpdatedAt: "2026-08-10 08:00"},
	}

	for _, entry := range []struct {
		kind  string
		items any
	}{
		{"accounts", customers}, {"products", products}, {"quotes", quotes}, {"orders", orders},
		{"documents", documents}, {"suppliers", suppliers}, {"private_files", privateFiles},
		{"news", news}, {"recommendations", recommendations}, {"jobs", jobs}, {"plugins", plugins},
	} {
		if err := s.seedSlice(ctx, entry.kind, entry.items); err != nil {
			return fmt.Errorf("seed %s: %w", entry.kind, err)
		}
	}
	preferences := UserPreferences{
		RecommendationEnabled: true,
		NewsShowLimit:         20,
		NewsFrequency:         "1小时",
		NewsCountries:         "德国、法国、波兰、瑞典",
		NewsTopics:            "E-bike、智能骑行、经销商、欧盟法规",
		NewsSources:           "EUR-Lex\nBike Europe\nCycling Industry News\nEurobike",
		AgentAllowlists:       map[string][]string{},
		AgentModelOverrides:   map[string]string{},
	}
	return s.putSetting(ctx, "preferences", preferences)
}

func (s *businessStore) seedSlice(ctx context.Context, kind string, input any) error {
	items, err := sliceRecords(input)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := s.create(ctx, kind, item.id, item.value); err != nil {
			return err
		}
	}
	return nil
}

type seedRecord struct {
	id    string
	value any
}

func sliceRecords(input any) ([]seedRecord, error) {
	result := make([]seedRecord, 0)
	switch values := input.(type) {
	case []Customer:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Product:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Quote:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Order:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Document:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Supplier:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []PrivateFile:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []NewsItem:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Recommendation:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Job:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	case []Plugin:
		for _, value := range values {
			result = append(result, seedRecord{value.ID, value})
		}
	default:
		return nil, fmt.Errorf("unsupported seed type %T", input)
	}
	return result, nil
}
