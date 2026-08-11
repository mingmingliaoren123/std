const state = {
  page: 'overview',
  lang: 'zh',
  subscription: true,
  recExpanded: false,
  agentCategory: 'all',
  customerSearch: '',
  customerType: 'all',
  customerCountry: 'all',
  customerSort: { field: 'updated', direction: 'desc' },
  quoteSearch: '',
  quoteStatus: 'all',
  quoteView: 'table',
  quoteSort: { field: 'value', direction: 'desc' },
  orderSearch: '',
  orderStatus: 'all',
  orderSort: { field: 'value', direction: 'desc' },
  documentSearch: '',
  documentType: 'all',
  documentStatus: 'all',
  productView: 'grid',
  productSearch: '',
  productCategory: 'all',
  productSort: 'updated',
  supplierSearch: '',
  supplierSort: 'updated',
  fileSearch: '',
  newsExpanded: false,
  newsCategory: '全部',
  newsShowLimit: 20,
  newsFrequency: '1小时',
  newsCountries: '德国、法国、波兰、瑞典',
  newsTopics: 'E-bike、智能骑行、经销商、欧盟法规',
  newsSources: 'EUR-Lex\nBike Europe\nCycling Industry News\nEurobike',
  settingsTab: 'model',
  modelConfigured: true,
  oemSort: 'score',
  oemTop: 3,
  oemCategory: '全部骑行类目',
  oemQuery: 'E-bike 电池 OEM 100 组',
  customerSearchMode: 'local',
  customerSearchQuery: '德国 E-bike 经销商',
  customerHasContact: true,
  discoveryCountry: '德国',
  discoveryCity: '柏林',
  discoveryType: 'Distributor',
  openClawStatus: null,
  openClawStatusLoading: false,
  openClawModels: null,
  openClawModelsLoading: false,
  openClawAgents: null,
  openClawAgentsLoading: false,
  agentChats: {},
  agentInternetAllowlists: {},
  agentSourceSelections: {},
  showApiKey: false,
  formContext: null,
  selectedRows: {
    customers: new Set(),
    quotes: new Set(),
    orders: new Set(),
  },
  quoteDraftLines: [],
  orderDraftLines: [],
  templateUploads: [],
  templateKind: 'document',
};

const authState = {
  username: 'admin',
  mode: 'login',
  message: '',
  authenticated: false,
  masterPassword: '',
};

const pageMeta = {
  overview: ['概览', 'STA-100 / 工作台', '📊'],
  agents: ['智能体', 'STA-100 / 智能协作', '🤖'],
  customers: ['客户', 'STA-100 / 业务管理', '👥'],
  quotes: ['报价单', 'STA-100 / 业务管理', '📄'],
  orders: ['订单', 'STA-100 / 业务管理', '📦'],
  documents: ['单据', 'STA-100 / 外贸单据', '🧾'],
  products: ['产品库', 'STA-100 / 业务管理', '🚲'],
  suppliers: ['供应商', 'STA-100 / 业务管理', '🏭'],
  database: ['数据库', 'STA-100 / 私有知识', '📚'],
  news: ['行业新闻', 'STA-100 / 行业情报', '📰'],
  settings: ['设置', 'STA-100 / 系统管理', '⚙️'],
};

const translations = {
  zh: {
    productSubtitle: '骑行行业智能工作台', navOverview: '概览', navAgents: '智能体', navCustomers: '客户',
    navQuotes: '报价单', navOrders: '订单', navDocuments: '单据', navProducts: '产品库', navSuppliers: '供应商', navDatabase: '数据库',
    navNews: '行业新闻', navSettings: '设置', serverOnline: '服务在线', serverRegion: '欧洲节点', localAccount: '本机账户',
  },
  en: {
    productSubtitle: 'Cycling Industry Workspace', navOverview: 'Overview', navAgents: 'Agents', navCustomers: 'Customers',
    navQuotes: 'Quotes', navOrders: 'Orders', navDocuments: 'Documents', navProducts: 'Products', navSuppliers: 'Suppliers', navDatabase: 'Database',
    navNews: 'Industry News', navSettings: 'Settings', serverOnline: 'Service online', serverRegion: 'Europe node', localAccount: 'Local account',
  },
};

const metrics = [
  { key: 'tasks', label: '今日待办', value: 8, icon: 'list-checks', detail: '3 项高优先级，5 项普通任务' },
  { key: 'meetings', label: '今日会议', value: 3, icon: 'calendar-clock', detail: '下一场 14:30 欧洲渠道会议' },
  { key: 'documents', label: '今日处理文档', value: 12, icon: 'file-check-2', detail: '10 个已完成，2 个待校正' },
  { key: 'orders', label: '进行中订单', value: 24, icon: 'package-open', detail: '订单模块实时聚合' },
  { key: 'chats', label: '今日对话', value: 47, icon: 'messages-square', detail: '来自 9 个专业智能体' },
  { key: 'news', label: '行业资讯', value: 15, icon: 'radio-tower', detail: '3 条高相关，12 条待浏览' },
];

const recommendations = [
  { title: 'Eurobike 2026 展商名录新增 86 家欧洲采购商', desc: '与“欧洲经销商、整车进口商”关注条件匹配，可进一步生成客户候选清单。', source: 'Eurobike', type: '展会情报', time: '12 分钟前' },
  { title: '欧盟电池法规尽职调查条款进入新执行阶段', desc: '可能影响 E-bike 电池产品资料与供应商声明，建议同步检查现有模板。', source: 'EUR-Lex', type: '法规', time: '36 分钟前' },
  { title: '德国通勤 E-bike 价格带出现结构性变化', desc: '中端 2,000-3,500 EUR 价格段热度上升，适合评估产品组合和报价策略。', source: 'Bike Europe', type: '市场趋势', time: '1 小时前' },
  { title: 'Shimano 发布新一代无线组件兼容信息', desc: '与本地产品兼容数据库存在 18 个待复核组合，可交由产品兼容智能体处理。', source: 'Shimano', type: '产品更新', time: '2 小时前' },
  { title: '波兰两家区域经销商公开新增采购负责人', desc: '已发现公开商务联系方式，可加入候选客户并安排跟进。', source: '公开网站', type: '客户线索', time: '3 小时前' },
];

const agents = [
  ['出口业务助手', 'ExportAgent', 'trade', 'briefcase-business', '报价、贸易条款、出口流程与风险检查', ['生成欧洲报价', '检查贸易条款', '准备出口清单']],
  ['支付条款助手', 'PaymentAdvisor', 'trade', 'landmark', '分析 T/T、L/C、O/A 等收款条件', ['比较付款方式', '评估账期风险', '生成付款建议']],
  ['物流路线助手', 'ShippingETA', 'trade', 'ship', '估算航线、时效、节点和异常风险', ['查询欧洲航线', '估算到港时间', '检查物流风险']],
  ['外贸单据助手', 'InvoiceAgent', 'trade', 'file-spreadsheet', '基于订单生成 PI、CI、PL 和报关单', ['生成 PI', '核对 CI', '生成全套单据']],
  ['中信保助手', 'SinosureAdvisor', 'trade', 'shield-check', '客户信用、额度与保险建议', ['评估客户信用', '核对信用额度', '准备投保资料']],
  ['合规计算助手', 'CBAMCalculator', 'trade', 'scale', '欧盟法规、关税和合规检查', ['查询欧盟法规', '检查 HS CODE', '生成合规摘要']],
  ['商务邮件助手', 'EmailGenerator', 'trade', 'mail-plus', '生成多语言客户开发与跟进邮件', ['写首次开发信', '跟进未回复客户', '翻译商务邮件']],
  ['价格监测助手', 'PriceTracker', 'trade', 'chart-no-axes-combined', '跟踪产品、原料和市场价格变化', ['比较价格变化', '查看价格预警', '生成调价建议']],
  ['库存管理助手', 'InventoryAgent', 'retail', 'warehouse', '库存盘点、补货与积压预警', ['查看缺货风险', '生成补货建议', '分析库存周转']],
  ['库存出清助手', 'InventoryClearanceAgent', 'retail', 'package-x', '为滞销库存匹配处置渠道和节奏', ['制定 30 天计划', '匹配清仓渠道', '估算回款']],
  ['B2B 平台助手', 'B2BMarketplaceAgent', 'retail', 'store', '产品上架、平台订单和渠道建议', ['准备平台上架', '比较销售渠道', '检查产品资料']],
  ['二手交易助手', 'UsedBikeTradingAgent', 'retail', 'refresh-cw', '二手车辆估值、检测和交易流程', ['评估二手车', '生成检测清单', '准备交易说明']],
  ['维修诊断助手', 'RepairQA', 'retail', 'wrench', '基于手册和案例进行故障诊断', ['诊断变速异响', '查询维修手册', '生成维修步骤']],
  ['产品兼容助手', 'CompatibilityAgent', 'retail', 'combine', '整车、组件与智能骑行设备兼容匹配', ['检查 Di2 兼容', '匹配功率计', '推荐码表组合']],
  ['市场分析助手', 'MarketAnalyzer', 'market', 'chart-spline', '国家、品类、价格段和增长因素分析', ['分析德国市场', '比较欧洲国家', '生成市场报告']],
  ['国家进入助手', 'CountryAdvisor', 'market', 'map', '国家政策、渠道、客户与进入节奏建议', ['选择目标国家', '制定进入计划', '查看法规风险']],
  ['展会助手', 'ExhibitionAdvisor', 'market', 'calendar-search', '展会筛选、准备、联系人和跟进', ['查询欧洲展会', '准备参展清单', '生成跟进计划']],
  ['赛事营销助手', 'TeamRaceAdvisor', 'market', 'trophy', '赛事、车队、赞助与品牌曝光分析', ['筛选合作赛事', '评估赞助价值', '生成合作方案']],
  ['供应商聚合助手', 'SupplierAggregator', 'support', 'factory', '供应商发现、对比和风险核查', ['寻找配件供应商', '比较供应能力', '检查供应商风险']],
  ['品牌情报助手', 'BrandValueCrawler', 'support', 'badge-check', '品牌动态、产品线和渠道价值分析', ['分析竞品品牌', '跟踪新品', '生成品牌对比']],
  ['客户洞察助手', 'CustomerMeasurementAgent', 'support', 'scan-search', '客户公开信息、需求和匹配度分析', ['分析目标客户', '补充公开信息', '生成跟进重点']],
  ['知识检索助手', 'RAGAgent', 'support', 'book-open-check', '融合通用知识库和客户私有数据检索', ['搜索本地资料', '对比本地与互联网', '列出引用来源']],
  ['产品设计助手', 'DesignAdvisor', 'support', 'pen-tool', '产品配色、几何、人体工学和规格建议', ['生成配色建议', '检查几何参数', '整理产品规格']],
  ['骑行路线助手', 'RouteFetcher', 'support', 'route', '路线、天气、补给点和骑行计划', ['规划骑行路线', '检查天气风险', '生成补给计划']],
];

const agentEmojis = ['🛒','💱','🚢','🧾','🛡️','🌱','✉️','💰','📦','💸','🤝','🔄','🔧','⚙️','📈','🌍','🎪','🏆','🏭','💎','📏','🛍️','🎨','🛤️'];
const agentIDs = ['export-agent','payment-advisor','shipping-eta','invoice-agent','sinosure-advisor','cbam-calculator','email-generator','price-tracker','inventory-agent','inventory-clearance-agent','b2b-marketplace-agent','used-bike-trading-agent','repair-qa','compatibility-agent','market-analyzer','country-advisor','exhibition-advisor','team-race-advisor','supplier-aggregator','brand-value-crawler','customer-measurement-agent','rag-agent','design-advisor','route-fetcher'];
const defaultInternetAllowlist = ['eur-lex.europa.eu','bike-eu.com','cyclingindustry.news','eurobike.com'];

const customers = [
  { id: 'ACC-0001', name: 'VeloTrade GmbH', type: 'Distributor', country: '德国', contact: 'Anna Keller', phone: '+49 30 555 0188', email: 'anna@velotrade.example', website: 'https://velotrade.example', owner: 'Donald', rating: 'Active', source: '展会', orders: 12, total: 'EUR 286,400', updated: '2026-08-10 09:42' },
  { id: 'ACC-0002', name: 'Nordic Cycle AB', type: 'Importer', country: '瑞典', contact: 'Erik Lund', phone: '+46 8 410 2250', email: 'erik@nordiccycle.example', website: 'https://nordiccycle.example', owner: 'Donald', rating: 'Prospect', source: '朋友介绍', orders: 3, total: 'EUR 61,900', updated: '2026-08-09 16:20' },
  { id: 'ACC-0003', name: 'SIM Sp. z o.o.', type: 'Customer', country: '波兰', contact: 'Marek Nowak', phone: '+48 34 310 2260', email: 'sales@sim.example', website: 'https://sim.example', owner: 'Donald', rating: 'Active', source: '电话', orders: 8, total: 'EUR 119,700', updated: '2026-08-08 11:15' },
  { id: 'ACC-0004', name: 'Ciclo Iberia S.L.', type: 'Reseller', country: '西班牙', contact: 'Lucia Martin', phone: '+34 91 778 2301', email: 'lucia@cicloiberia.example', website: 'https://cicloiberia.example', owner: 'Donald', rating: 'Prospect', source: '拜访', orders: 1, total: 'EUR 18,250', updated: '2026-08-07 14:02' },
  { id: 'ACC-0005', name: 'Alpine Motion SAS', type: 'Integrator', country: '法国', contact: 'Louis Bernard', phone: '+33 1 8420 1189', email: 'louis@alpinemotion.example', website: 'https://alpinemotion.example', owner: 'Donald', rating: 'Active', source: '互联网线索', orders: 6, total: 'EUR 94,600', updated: '2026-08-06 10:30' },
];

const quotes = [
  { id: 'QUO-2026-0188', subject: 'STA-100 首批设备报价', customer: 'VeloTrade GmbH', value: 'EUR 34,800', valid: '2026-09-08', status: 'Delivered', products: 'STA-100 x 200', owner: 'Donald' },
  { id: 'QUO-2026-0187', subject: '智能骑行组件组合', customer: 'Nordic Cycle AB', value: 'EUR 18,450', valid: '2026-09-02', status: 'Draft', products: '功率计套装 x 50', owner: 'Donald' },
  { id: 'QUO-2026-0185', subject: '波兰渠道补货报价', customer: 'SIM Sp. z o.o.', value: 'EUR 22,100', valid: '2026-08-26', status: 'Accepted', products: 'STA-100 x 120', owner: 'Donald' },
  { id: 'QUO-2026-0179', subject: '西班牙测试订单报价', customer: 'Ciclo Iberia S.L.', value: 'EUR 5,760', valid: '2026-08-15', status: 'Rejected', products: 'STA-100 x 32', owner: 'Donald' },
];

const orders = [
  { id: 'SO-2026-0106', customer: 'SIM Sp. z o.o.', quote: 'QUO-2026-0185', products: 'STA-100 x 120', value: 'EUR 22,100', status: 'Production', delivery: '2026-09-10', progress: 55 },
  { id: 'SO-2026-0105', customer: 'VeloTrade GmbH', quote: 'QUO-2026-0174', products: 'STA-100 x 80', value: 'EUR 14,900', status: 'Confirmed', delivery: '2026-08-30', progress: 30 },
  { id: 'SO-2026-0102', customer: 'Alpine Motion SAS', quote: 'QUO-2026-0162', products: '智能组件套装 x 40', value: 'EUR 31,600', status: 'Shipped', delivery: '2026-08-18', progress: 82 },
  { id: 'SO-2026-0098', customer: 'Nordic Cycle AB', quote: 'QUO-2026-0151', products: 'STA-100 x 25', value: 'EUR 4,650', status: 'Completed', delivery: '2026-08-03', progress: 100 },
];

const documents = [
  { id: 'PI-20260810-003', type: 'PI', customer: 'SIM Sp. z o.o.', order: 'SO-2026-0106', template: 'STRATRONIX 标准 PI v3', status: 'Draft', updated: '2026-08-10 10:22' },
  { id: 'CI-20260809-012', type: 'CI', customer: 'Alpine Motion SAS', order: 'SO-2026-0102', template: '欧盟商业发票 v2', status: 'Confirmed', updated: '2026-08-09 17:45' },
  { id: 'PL-20260809-009', type: 'PL', customer: 'Alpine Motion SAS', order: 'SO-2026-0102', template: '标准装箱单 v2', status: 'Confirmed', updated: '2026-08-09 17:41' },
  { id: 'CD-20260808-006', type: '报关单', customer: 'VeloTrade GmbH', order: 'SO-2026-0105', template: '出口报关单 v1', status: 'Review', updated: '2026-08-08 15:12' },
];

const products = [
  { id: 'STA-100-EU', name: 'STA-100 私有智能体设备', category: '智能设备', price: 'EUR 159.00', stock: 426, hs: '8471504090', status: 'Active', desc: '4G+32G，Linux OS，预装 OpenClaw Framework，支持客户私有数据。' },
  { id: 'PM-DUAL-01', name: '双边功率计套装', category: '智能骑行', price: 'EUR 438.00', stock: 68, hs: '9029209000', status: 'Active', desc: '支持 ANT+ / BLE，提供曲柄兼容信息和出厂校准记录。' },
  { id: 'GPS-PRO-02', name: '骑行码表 Pro', category: '智能骑行', price: 'EUR 219.00', stock: 112, hs: '8526919090', status: 'Active', desc: '多卫星定位、路线导航与训练数据同步。' },
  { id: 'EBK-CITY-03', name: '城市 E-bike 方案', category: '整车方案', price: 'EUR 1,280.00', stock: 19, hs: '8711601000', status: 'Review', desc: '面向欧洲城市通勤渠道的电助力整车组合方案。' },
];

const suppliers = [
  { id: 'SUP-0001', company: '苏州骑行动力科技', phone: '0512-6688 2100', contact: '周敏', email: 'sales@suzhou-cycle.example', product: 'E-bike 电池', specification: '36V / 48V，定制 BMS', quote: 'EUR 82 / 组', notes: '可提供 CE、UN38.3 资料', source: '展会', updated: '2026-08-10 08:50' },
  { id: 'SUP-0002', company: '东莞骑迹智能装备', phone: '0769-8822 1177', contact: '陈杰', email: 'business@qiji.example', product: '码表、功率计', specification: 'ANT+ / BLE，OEM 包装', quote: 'EUR 38 / 件起', notes: '支持小批量 ODM', source: '朋友介绍', updated: '2026-08-09 15:20' },
  { id: 'SUP-0003', company: '厦门轻量运动用品', phone: '0592-6200 3318', contact: '林晓', email: 'export@lightmotion.example', product: '头盔、骑行装备', specification: 'MIPS，EN1078', quote: 'EUR 19 / 件起', notes: '认证文件待复核', source: '电话', updated: '2026-08-08 10:05' },
];

const files = [
  { name: 'STA100_Product_Manual_EN.pdf', category: '产品手册', tags: ['STA-100', '英文'], size: '6.8 MB', source: '出厂通用资料', status: 'Indexed', updated: '2026-08-10 08:20' },
  { name: 'VeloTrade_Distribution_Agreement.docx', category: '合同', tags: ['客户私有', '德国'], size: '1.2 MB', source: '客户上传', status: 'Indexed', updated: '2026-08-09 16:10' },
  { name: 'EU_Battery_Regulation_2026.pdf', category: '法规', tags: ['欧盟', 'E-bike'], size: '3.1 MB', source: '通用知识库', status: 'Indexed', updated: '2026-08-09 12:45' },
  { name: 'Shimano_Compatibility_List.xlsx', category: '产品资料', tags: ['兼容', '组件'], size: '846 KB', source: '客户上传', status: 'Review', updated: '2026-08-08 18:33' },
];

const news = [
  { category: '欧洲市场', title: '欧洲自行车产业进入补库存周期，渠道更关注小批量和快速交付', summary: '多家欧洲经销商在 2026 年下半年调整采购节奏，订单结构从大批量预采转向小批量、多批次。', source: 'Bike Europe', time: '2026-08-10 09:10', relevance: '96%' },
  { category: '法规', title: '欧盟更新电池尽职调查实施指引，E-bike 供应链资料需同步准备', summary: '新指引强化材料来源、碳足迹和供应链证明要求。', source: 'EUR-Lex', time: '2026-08-10 08:35', relevance: '93%' },
  { category: '智能骑行', title: '无线电子变速与功率数据融合成为高端整车配置趋势', summary: '整车厂正在把兼容性和训练数据完整度作为高端产品卖点。', source: 'Cycling Industry News', time: '2026-08-09 17:40', relevance: '89%' },
  { category: '渠道', title: '北欧经销商加快建设线上线下一体的维修服务网络', summary: '服务能力和备件响应速度正在影响品牌进入门槛。', source: 'Nordic Cycling', time: '2026-08-09 15:20', relevance: '84%' },
  { category: '产品', title: '欧洲城市通勤市场对轻量化 E-bike 的关注持续上升', summary: '重量、可维护性和电池合规成为渠道选品主要指标。', source: 'E-bike News', time: '2026-08-09 11:05', relevance: '82%' },
];

const oemFactories = [
  { name: '苏州骑行动力科技', category: 'E-bike 电池', city: '苏州', tier: 'S', score: 94, capacity: '30,000 组/月', moq: '100 组', source: '本地知识库', reason: '已有欧洲电池产品出货记录，MOQ 与当前需求匹配' },
  { name: '宁波远行整车制造', category: '公路/山地整车', city: '宁波', tier: 'A', score: 91, capacity: '12,000 台/月', moq: '300 台', source: '联网检索', reason: '联网公开资料经检索整理后显示具备整车研发和小批量 OEM 能力' },
  { name: '东莞骑迹智能装备', category: '码表/智能设备', city: '东莞', tier: 'A', score: 88, capacity: '8,000 台/月', moq: '200 台', source: '联网检索', reason: '联网公开资料显示支持 ANT+/BLE 产品和 ODM 服务' },
  { name: '厦门轻量运动用品', category: '头盔/骑行装备', city: '厦门', tier: 'B', score: 84, capacity: '50,000 件/月', moq: '500 件', source: '本地知识库 + 联网检索', reason: '本地档案结合联网资料后显示装备类目覆盖完整' },
  { name: '常州链动传动系统', category: '链条/传动', city: '常州', tier: 'B', score: 81, capacity: '80,000 件/月', moq: '1,000 件', source: '联网检索', reason: '联网行业资料较完整，认证和实际交付仍需复核' },
];

const unifiedSearchCustomers = [
  { name: 'VeloTrade GmbH', country: '德国', type: 'Distributor', contact: 'anna@velotrade.example', business: 'E-bike / 经销商', source: '本地知识库', score: 96 },
  { name: 'Nordic Cycle AB', country: '瑞典', type: 'Importer', contact: '+46 8 410 2250', business: '城市车 / 进口商', source: '本地知识库 + 联网检索', score: 91 },
  { name: 'Berlin Motion Handels', country: '德国', type: 'Prospect', contact: 'sales@berlinmotion.example', business: '智能骑行 / 零售渠道', source: '联网检索', score: 88 },
  { name: 'Ciclo Verde S.L.', country: '西班牙', type: 'Dealer', contact: 'contact@cicloverde.example', business: '公路车 / 车店', source: '联网检索', score: 83 },
];

const localDiscoveryLeads = [
  { name: 'Berlin Motion Handels', city: '柏林', country: '德国', type: 'Distributor', contact: 'sales@berlinmotion.example', score: 93, reason: '公开渠道显示具备 E-bike 区域经销网络' },
  { name: 'Rheinland Cycle Network', city: '科隆', country: '德国', type: 'Importer', contact: '+49 221 555 0142', score: 88, reason: '进口业务与产品资料需求和当前目标匹配' },
  { name: 'München Bike Lab', city: '慕尼黑', country: '德国', type: 'Dealer', contact: 'hello@munichbikelab.example', score: 82, reason: '高端智能骑行产品门店，适合小批测试' },
];

const discoveryCities = {
  中国: ['北京','上海','广州','深圳','杭州','成都','重庆','苏州','宁波','厦门','东莞','天津','青岛','武汉','西安'],
  阿尔巴尼亚: ['地拉那','都拉斯','发罗拉'], 安道尔: ['安道尔城','莱塞斯卡尔德'], 亚美尼亚: ['埃里温','久姆里'],
  奥地利: ['维也纳','格拉茨','林茨','萨尔茨堡','因斯布鲁克'], 阿塞拜疆: ['巴库','占贾'], 白俄罗斯: ['明斯克','戈梅利','布列斯特'],
  比利时: ['布鲁塞尔','安特卫普','根特','列日','布鲁日'], 波黑: ['萨拉热窝','巴尼亚卢卡','莫斯塔尔'], 保加利亚: ['索非亚','普罗夫迪夫','瓦尔纳','布尔加斯'],
  克罗地亚: ['萨格勒布','斯普利特','里耶卡','奥西耶克'], 塞浦路斯: ['尼科西亚','利马索尔','拉纳卡'], 捷克: ['布拉格','布尔诺','俄斯特拉发','比尔森'],
  丹麦: ['哥本哈根','奥胡斯','欧登塞','奥尔堡'], 爱沙尼亚: ['塔林','塔尔图','纳尔瓦'], 芬兰: ['赫尔辛基','埃斯波','坦佩雷','图尔库','奥卢'],
  法国: ['巴黎','里昂','马赛','图卢兹','波尔多','里尔','南特','斯特拉斯堡'], 格鲁吉亚: ['第比利斯','巴统','库塔伊西'],
  德国: ['柏林','汉堡','慕尼黑','科隆','法兰克福','斯图加特','杜塞尔多夫','莱比锡','纽伦堡'],
  希腊: ['雅典','塞萨洛尼基','帕特雷','伊拉克利翁'], 匈牙利: ['布达佩斯','德布勒森','塞格德','米什科尔茨'], 冰岛: ['雷克雅未克','科帕沃于尔','阿克雷里'],
  爱尔兰: ['都柏林','科克','利默里克','戈尔韦'], 意大利: ['罗马','米兰','都灵','佛罗伦萨','博洛尼亚','那不勒斯','威尼斯'], 科索沃: ['普里什蒂纳','普里兹伦'],
  哈萨克斯坦: ['阿斯塔纳','阿拉木图','阿克套'], 拉脱维亚: ['里加','陶格夫匹尔斯','利耶帕亚'], 列支敦士登: ['瓦杜兹','沙恩'], 立陶宛: ['维尔纽斯','考纳斯','克莱佩达'], 卢森堡: ['卢森堡市','埃施'],
  马耳他: ['瓦莱塔','斯利马','比尔基卡拉'], 摩尔多瓦: ['基希讷乌','伯尔兹'], 摩纳哥: ['摩纳哥城','蒙特卡洛'], 黑山: ['波德戈里察','布德瓦','科托尔'],
  荷兰: ['阿姆斯特丹','鹿特丹','乌得勒支','海牙','埃因霍温'], 北马其顿: ['斯科普里','比托拉','奥赫里德'], 挪威: ['奥斯陆','卑尔根','特隆赫姆','斯塔万格'],
  波兰: ['华沙','克拉科夫','弗罗茨瓦夫','波兹南','格但斯克','罗兹'], 葡萄牙: ['里斯本','波尔图','布拉加','科英布拉'],
  罗马尼亚: ['布加勒斯特','克卢日-纳波卡','蒂米什瓦拉','雅西'], 俄罗斯: ['莫斯科','圣彼得堡','喀山','加里宁格勒'], 圣马力诺: ['圣马力诺市','塞拉瓦莱'],
  塞尔维亚: ['贝尔格莱德','诺维萨德','尼什'], 斯洛伐克: ['布拉迪斯拉发','科希策','日利纳'], 斯洛文尼亚: ['卢布尔雅那','马里博尔','科佩尔'],
  西班牙: ['马德里','巴塞罗那','瓦伦西亚','塞维利亚','马拉加','毕尔巴鄂'], 瑞典: ['斯德哥尔摩','哥德堡','马尔默','乌普萨拉'],
  瑞士: ['苏黎世','日内瓦','巴塞尔','伯尔尼','洛桑'], 土耳其: ['伊斯坦布尔','安卡拉','伊兹密尔','布尔萨'], 乌克兰: ['基辅','利沃夫','敖德萨','第聂伯罗'],
  英国: ['伦敦','曼彻斯特','伯明翰','爱丁堡','布里斯托尔','利兹'], 梵蒂冈: ['梵蒂冈城'],
};

function icon(name) { return `<i data-lucide="${name}"></i>`; }
function escapeAttr(value) { return String(value).replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;'); }
function escapeHTML(value) { return String(value ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#39;'); }
function authMessage(message='') {
  authState.message = message;
  const node = document.getElementById('authMessage');
  if (node) node.textContent = message;
}
function clearToasts() {
  const region = document.getElementById('toastRegion');
  if (region) region.replaceChildren();
}
function renderAuthScreen() {
  const screen = document.getElementById('authScreen');
  if (!screen) return;
  const loginView = authState.mode === 'login';
  const masterView = authState.mode === 'master';
  const resetView = authState.mode === 'reset';
  const body = loginView ? `<form class="auth-form" id="loginForm"><div class="form-field"><label for="loginUsername">登录用户名</label><input class="input" id="loginUsername" value="${escapeAttr(authState.username)}" readonly></div><div class="form-field"><label for="loginPassword">密码</label><input class="input" id="loginPassword" type="password" autocomplete="current-password" autofocus required></div><div id="authMessage" class="auth-message">${escapeHTML(authState.message)}</div><button class="button primary" type="submit">${icon('log-in')}登录</button><div class="auth-actions"><span class="secondary-text">本机账户</span><button class="auth-link" type="button" id="forgotPassword">忘记密码？</button></div></form>`
    : masterView ? `<form class="auth-form" id="masterForm"><div class="auth-back"><button class="auth-link" type="button" id="backToLogin">${icon('arrow-left')}返回登录</button></div><div class="form-field"><label for="masterPassword">万能密码</label><input class="input" id="masterPassword" type="password" autocomplete="off" autofocus required></div><div id="authMessage" class="auth-message">${escapeHTML(authState.message)}</div><button class="button primary" type="submit">${icon('shield-check')}验证并重置</button><p class="auth-note">仅用于设备维护。验证通过后可重新设置当前本机账户的用户名和密码。</p></form>`
    : `<form class="auth-form" id="resetForm"><div class="auth-back"><button class="auth-link" type="button" id="backToMaster">${icon('arrow-left')}返回上一步</button></div><div class="form-field"><label for="resetUsername">新用户名</label><input class="input" id="resetUsername" value="${escapeAttr(authState.username)}" autocomplete="username" required></div><div class="form-field"><label for="resetPassword">新密码</label><input class="input" id="resetPassword" type="password" autocomplete="new-password" minlength="4" required></div><div class="form-field"><label for="resetPasswordConfirm">确认新密码</label><input class="input" id="resetPasswordConfirm" type="password" autocomplete="new-password" minlength="4" required></div><div id="authMessage" class="auth-message">${escapeHTML(authState.message)}</div><button class="button primary" type="submit">${icon('key-round')}保存并进入工作台</button></form>`;
  screen.innerHTML = `<section class="auth-card"><header class="auth-head"><div class="auth-brand"><img src="assets/cycling-agent-icon.jpg" alt="STA-100"><div><strong>STA-100</strong><span>骑行行业智能工作台</span></div></div><h1 id="authTitle">${loginView ? '登录工作台' : resetView ? '重置本机账户' : '账户恢复'}</h1><p>${loginView ? `请输入 ${escapeHTML(authState.username)} 的密码继续使用。` : resetView ? '请设置新的用户名和登录密码。' : '请输入设备维护万能密码进入重置流程。'}</p></header><div class="auth-body">${body}</div></section>`;
  applyIcons();
  document.getElementById('loginForm')?.addEventListener('submit', event => { event.preventDefault(); void loginUser(); });
  document.getElementById('masterForm')?.addEventListener('submit', event => { event.preventDefault(); void verifyMasterPassword(); });
  document.getElementById('resetForm')?.addEventListener('submit', event => { event.preventDefault(); void resetCredentials(); });
  document.getElementById('forgotPassword')?.addEventListener('click', () => { authState.mode='master'; authState.masterPassword=''; authMessage(''); renderAuthScreen(); });
  document.getElementById('backToLogin')?.addEventListener('click', () => { authState.mode='login'; authState.masterPassword=''; authMessage(''); renderAuthScreen(); });
  document.getElementById('backToMaster')?.addEventListener('click', () => { authState.mode='master'; authState.masterPassword=''; authMessage(''); renderAuthScreen(); });
}
async function loadAuthCredentials() {
  try {
    const status = await apiFetch('/api/v1/auth/status');
    authState.username = status.username || 'admin';
    if (status.authenticated) showAuthenticatedApp();
  } catch {
    authState.username = 'admin';
  }
  if (!authState.authenticated) renderAuthScreen();
}
function updateSidebarIdentity() {
  const username = authState.username || 'admin';
  const name = document.getElementById('sidebarUsername');
  const avatar = document.getElementById('sidebarAvatar');
  if (name) name.textContent = username;
  if (avatar) avatar.textContent = username.slice(0, 2).toUpperCase();
}
function showAuthenticatedApp() {
  authState.authenticated = true;
  document.getElementById('authScreen').hidden = true;
  document.getElementById('appShell').hidden = false;
  updateSidebarIdentity();
  setPage(state.page);
}
async function logoutUser() {
  try { await apiFetch('/api/v1/auth/logout', { method: 'POST', body: '{}' }); } catch { /* session may already be expired */ }
  closeModal(); closeDrawer(); closeCommand();
  authState.authenticated = false;
  authState.mode = 'login';
  authState.masterPassword = '';
  clearToasts();
  authMessage('');
  document.getElementById('appShell').hidden = true;
  document.getElementById('authScreen').hidden = false;
  renderAuthScreen();
}
async function loginUser() {
  const password = document.getElementById('loginPassword')?.value || '';
  try {
    const result = await apiFetch('/api/v1/auth/login', { method: 'POST', body: JSON.stringify({ username: authState.username, password }) });
    authState.username = result.username || authState.username;
    authMessage('');
    showAuthenticatedApp();
  } catch (error) {
    authMessage(error.message || '用户名或密码不正确。');
    document.getElementById('loginPassword')?.focus();
  }
}
async function verifyMasterPassword() {
  const password = document.getElementById('masterPassword')?.value || '';
  try {
    await apiFetch('/api/v1/auth/verify-master', { method: 'POST', body: JSON.stringify({ masterPassword: password }) });
    authState.masterPassword = password;
    authState.mode = 'reset'; authMessage(''); renderAuthScreen();
  } catch (error) {
    authMessage(error.message || '万能密码不正确，无法进入重置流程。');
  }
}
async function resetCredentials() {
  const username = document.getElementById('resetUsername')?.value.trim() || '';
  const password = document.getElementById('resetPassword')?.value || '';
  const confirmation = document.getElementById('resetPasswordConfirm')?.value || '';
  if (!username) { authMessage('用户名不能为空。'); return; }
  if (password.length < 4) { authMessage('新密码至少需要 4 个字符。'); return; }
  if (password !== confirmation) { authMessage('两次输入的新密码不一致。'); return; }
  try {
    const result = await apiFetch('/api/v1/auth/reset', { method: 'POST', body: JSON.stringify({ masterPassword: authState.masterPassword, username, password }) });
    authState.username = result.username || username; authState.masterPassword = ''; authState.mode='login'; authMessage('');
    showAuthenticatedApp();
    toast('账户已重置', `当前用户名已更新为 ${authState.username}。`, 'success');
  } catch (error) {
    authState.masterPassword = '';
    authState.mode = 'master';
    authMessage(error.message || '账户重置失败，请重新验证万能密码。');
    renderAuthScreen();
  }
}
function openAccountSettings() {
  openModal({ title: '账户与密码', eyebrow: '设置 / 本机账户', body: `<div class="form-grid"><div class="form-field full"><label>当前用户名</label><input class="input" value="${escapeAttr(authState.username)}" readonly></div>${inputField('当前密码','',true,false,'password','accountCurrentPassword')}${inputField('新用户名',authState.username,true,false,'text','accountNewUsername')}${inputField('新密码（留空则不修改）','',false,false,'password','accountNewPassword')}${inputField('确认新密码','',false,false,'password','accountNewPasswordConfirm')}<div class="auth-note form-field full"><span>${icon('shield-check')} 修改账户需要验证当前密码；忘记当前密码时请退出后使用登录页的万能密码重置流程。</span></div></div>`, footer: `${formFooter('保存账户','save-account-settings')}` });
}
async function saveAccountSettings() {
  const current = document.getElementById('accountCurrentPassword')?.value || '';
  const username = document.getElementById('accountNewUsername')?.value.trim() || '';
  const password = document.getElementById('accountNewPassword')?.value || '';
  const confirmation = document.getElementById('accountNewPasswordConfirm')?.value || '';
  if (!username) { toast('保存失败','用户名不能为空。','warning'); return; }
  if (password && password.length < 4) { toast('保存失败','新密码至少需要 4 个字符。','warning'); return; }
  if (password !== confirmation) { toast('保存失败','两次输入的新密码不一致。','warning'); return; }
  try {
    const result = await apiFetch('/api/v1/auth/account', { method: 'PATCH', body: JSON.stringify({ currentPassword: current, username, password: password || current }) });
    authState.username = result.username || username;
    updateSidebarIdentity(); closeModal(); renderPage(); toast('账户已更新', '下次登录将显示新的用户名。', 'success');
  } catch (error) {
    toast('保存失败', error.message || '账户保存失败。', 'warning');
  }
}
async function initAuth() {
  await loadAuthCredentials();
}
function badge(status) {
  const map = { Active: ['green', '活跃'], Prospect: ['blue', '潜在'], Customer: ['green', '客户'], Draft: ['neutral', '草稿'], Delivered: ['blue', '已发送'], Accepted: ['green', '已接受'], Rejected: ['red', '已拒绝'], Confirmed: ['green', '已确认'], Production: ['amber', '生产中'], Shipped: ['blue', '已发运'], Completed: ['green', '已完成'], Review: ['amber', '待复核'], Indexed: ['green', '已索引'] };
  const [cls, label] = map[status] || ['neutral', status];
  return `<span class="badge ${cls}">${label}</span>`;
}

function applyIcons() {
  if (window.lucide) window.lucide.createIcons({ attrs: { 'aria-hidden': 'true' } });
}

function applyTranslations() {
  document.querySelectorAll('[data-i18n]').forEach(el => {
    const value = translations[state.lang][el.dataset.i18n];
    if (value) el.textContent = value;
  });
}

async function apiFetch(path, options = {}) {
  const headers = { Accept: 'application/json', ...(options.headers || {}) };
  if (options.body) headers['Content-Type'] = 'application/json';
  if (options.method && options.method !== 'GET') headers['X-STA100-Request'] = '1';
  const response = await fetch(path, { ...options, headers });
  let data = null;
  try { data = await response.json(); } catch { data = null; }
  if (!response.ok) throw new Error(data?.error?.message || `请求失败 (${response.status})`);
  return data;
}

async function loadOpenClawStatus(force = false) {
  if (state.openClawStatusLoading || (state.openClawStatus && !force)) return state.openClawStatus;
  state.openClawStatusLoading = true;
  if (state.page === 'settings' && state.settingsTab === 'system') renderPage();
  try {
    state.openClawStatus = await apiFetch('/api/v1/openclaw/status');
  } catch (error) {
    state.openClawStatus = { available: false, error: error.message };
  } finally {
    state.openClawStatusLoading = false;
    if (state.page === 'settings' && state.settingsTab === 'system') renderPage();
  }
  return state.openClawStatus;
}

async function loadOpenClawModels(force = false) {
  if (state.openClawModelsLoading || (state.openClawModels && !force)) return state.openClawModels;
  state.openClawModelsLoading = true;
  if (state.page === 'settings' && state.settingsTab === 'model') renderPage();
  try {
    state.openClawModels = await apiFetch('/api/v1/openclaw/models');
    state.modelConfigured = Boolean(state.openClawModels.configured);
  } catch (error) {
    state.openClawModels = { configured: false, models: [], providers: [], error: error.message };
    state.modelConfigured = false;
  } finally {
    state.openClawModelsLoading = false;
    if (state.page === 'settings' && state.settingsTab === 'model') renderPage();
  }
  return state.openClawModels;
}

async function loadOpenClawAgents(force = false) {
  if (state.openClawAgentsLoading || (state.openClawAgents && !force)) return state.openClawAgents;
  state.openClawAgentsLoading = true;
  if (state.page === 'agents' || (state.page === 'settings' && state.settingsTab === 'system')) renderPage();
  try {
    const data = await apiFetch('/api/v1/openclaw/agents');
    state.openClawAgents = data.agents || [];
  } catch (error) {
    state.openClawAgents = [];
    toast('OpenClaw Agent 读取失败', error.message, 'warning');
  } finally {
    state.openClawAgentsLoading = false;
    if (state.page === 'agents' || (state.page === 'settings' && state.settingsTab === 'system')) renderPage();
  }
  return state.openClawAgents;
}

function loadPageOpenClawData() {
  if (state.page === 'agents') void loadOpenClawAgents();
  if (state.page !== 'settings') return;
  if (state.settingsTab === 'model') void loadOpenClawModels();
  if (state.settingsTab === 'system') {
    void loadOpenClawStatus();
    void loadOpenClawAgents();
  }
}

function setPage(page) {
  if (!pageMeta[page]) return;
  state.page = page;
  const [title, eyebrow, emoji] = pageMeta[page];
  document.getElementById('pageTitle').textContent = `${emoji} ${title}`;
  document.getElementById('pageEyebrow').textContent = eyebrow;
  document.querySelectorAll('.nav-item').forEach(item => item.classList.toggle('active', item.dataset.page === page));
  document.getElementById('sidebar').classList.remove('open');
  renderPage();
  loadPageOpenClawData();
  history.replaceState(null, '', `#${page}`);
  document.getElementById('pageRoot').focus({ preventScroll: true });
}

function renderPage() {
  const renderers = { overview: renderOverview, agents: renderAgents, customers: renderCustomers, quotes: renderQuotes, orders: renderOrders, documents: renderDocuments, products: renderProducts, suppliers: renderSuppliers, database: renderDatabase, news: renderNews, settings: renderSettings };
  document.getElementById('pageRoot').innerHTML = renderers[state.page]();
  applyIcons();
  wirePageSpecific();
}

function renderOverview() {
  const recs = state.recExpanded ? recommendations : recommendations.slice(0, 3);
  return `<div class="page-stack">
    <section class="hero-strip panel">
      <div>
        <span class="badge green">📊 日报订阅</span>
        <h2>行业信息和业务进展，按你的关注条件持续更新</h2>
        <p>系统每 ${state.newsFrequency} 由“为你推荐”智能体更新指定网站和平台，并与本地业务数据合并去重。</p>
      </div>
      <div class="subscribe-area">
        <div class="subscribe-meta"><strong>${state.subscription ? '订阅已开启' : '订阅已暂停'}</strong><span>${state.subscription ? '下次更新 10:45' : '不会自动更新'}</span></div>
        <label class="toggle"><input id="subscriptionToggle" type="checkbox" ${state.subscription ? 'checked' : ''}><span></span></label>
      </div>
    </section>

    <section class="metric-grid" aria-label="今日业务摘要">
      ${metrics.map(m => `<button class="metric-button" data-action="metric-detail" data-key="${m.key}"><span class="metric-icon">${icon(m.icon)}</span><span><strong class="metric-number">${m.value}</strong><span class="metric-label">${m.label}</span></span></button>`).join('')}
    </section>

    <section class="content-grid">
      <div class="panel">
        <header class="panel-head"><div><h3>为你推荐</h3><p>根据关注条件、用户操作和智能体记录生成</p></div><div class="inline-actions"><span class="badge blue">刚刚更新</span><button class="link-button" data-action="recommend-settings">推荐设置</button></div></header>
        <div class="recommendation-list">
          ${recs.map((r, i) => `<article class="recommendation-item"><span class="recommendation-rank">${String(i + 1).padStart(2, '0')}</span><div><h4>${r.title}</h4><p>${r.desc}</p><div class="source-line"><span class="mini-source">来源 <strong>${r.source}</strong></span><span class="mini-source">类型 <strong>${r.type}</strong></span><span class="mini-source">${r.time}</span></div></div><button class="table-icon" data-action="recommend-detail" data-index="${i}" aria-label="查看详情" title="查看详情">${icon('arrow-up-right')}</button></article>`).join('')}
        </div>
        <div class="panel-body" style="padding-top:10px;text-align:center"><button class="button ghost small" data-action="toggle-recommendations">${icon(state.recExpanded ? 'chevron-up' : 'chevron-down')}${state.recExpanded ? '收起推荐' : '查看更多推荐'}</button></div>
      </div>
      <aside class="panel">
        <header class="panel-head"><div><h3>行业新闻</h3><p>默认展示相关度最高的 3 条</p></div><button class="link-button" data-page="news">MORE</button></header>
        <div class="news-mini">${news.slice(0,3).map(n => `<button class="news-mini-item" data-action="news-detail" data-title="${n.title}"><span class="news-time">${n.time.slice(11)}</span><span><h4>${n.title}</h4><p>${n.source} · 相关度 ${n.relevance}</p></span></button>`).join('')}</div>
      </aside>
    </section>

    <div class="section-head tool-section-head"><div><h3>🛠️ 智能业务工具</h3><p>保留原应用的 OEM 匹配、客户统一搜索和本地客户发现能力。</p></div><span class="meta">本地知识库优先 · 来源分别展示</span></div>
    <section class="tool-grid">
      <article class="panel tool-panel tool-panel-wide">
        <header class="tool-header"><div><h3>🏭 OEM 工厂智能匹配</h3><p>按骑行类目融合本地知识库与联网检索；RAG 是检索实现方式，不单独作为数据来源。</p></div><span class="badge amber">规则 TODO</span></header>
        <div class="panel-body">
          <div class="filter-row tool-presets">${['公路整车 OEM 1000 台','E-bike 电池 OEM 100 组','中置电机 500 套','头盔 MIPS 500 个'].map(v=>`<button class="filter-chip" data-action="oem-preset" data-value="${v}">${v}</button>`).join('')}</div>
          <div class="tool-form">
            <label class="field-search tool-query">${icon('search')}<input id="oemQuery" value="${escapeAttr(state.oemQuery)}" placeholder="输入产品、数量、市场和要求"></label>
            <select class="select" id="oemCategory"><option ${state.oemCategory==='全部骑行类目'?'selected':''}>全部骑行类目</option><option ${state.oemCategory==='整车'?'selected':''}>整车</option><option ${state.oemCategory==='E-bike 电池'?'selected':''}>E-bike 电池</option><option ${state.oemCategory==='电机'?'selected':''}>电机</option><option ${state.oemCategory==='链条/传动'?'selected':''}>链条/传动</option><option ${state.oemCategory==='轮胎'?'selected':''}>轮胎</option><option ${state.oemCategory==='头盔'?'selected':''}>头盔</option><option ${state.oemCategory==='码表/智能设备'?'selected':''}>码表/智能设备</option><option ${state.oemCategory==='功率计'?'selected':''}>功率计</option></select>
            <select class="select" id="oemSort"><option value="score" ${state.oemSort==='score'?'selected':''}>按匹配度排序</option><option value="capacity" ${state.oemSort==='capacity'?'selected':''}>按产能排序</option><option value="moq" ${state.oemSort==='moq'?'selected':''}>按 MOQ 排序</option><option value="source" ${state.oemSort==='source'?'selected':''}>按数据来源排序</option></select>
            <select class="select" id="oemTop"><option value="3" ${state.oemTop===3?'selected':''}>Top 3</option><option value="5" ${state.oemTop===5?'selected':''}>Top 5</option><option value="10" ${state.oemTop===10?'selected':''}>Top 10</option></select>
            <button class="button primary" data-action="oem-run">${icon('scan-search')}开始匹配</button>
          </div>
          <div class="source-legend"><span>${icon('hard-drive')}本地知识库</span><span>${icon('globe-2')}联网检索（RAG 检索 + 公开网络采集）</span><small>正式分类和评分权重待数据提供后确认</small></div>
          <div class="tool-results" id="oemResults">${renderOEMCards()}</div>
        </div>
      </article>

      <article class="panel tool-panel">
        <header class="tool-header"><div><h3>🔍 客户统一搜索</h3><p>统一字段检索本地知识库和联网信息，默认使用本地知识库。</p></div></header>
        <div class="panel-body">
          <div class="tool-form compact">
            <label class="field-search tool-query">${icon('search')}<input id="unifiedCustomerQuery" value="${escapeAttr(state.customerSearchQuery)}" placeholder="国家、公司、业务、邮箱或电话"></label>
            <select class="select" id="unifiedSearchMode"><option value="local" ${state.customerSearchMode==='local'?'selected':''}>本地知识库（默认）</option><option value="rag" ${state.customerSearchMode==='rag'?'selected':''}>联网检索</option><option value="hybrid" ${state.customerSearchMode==='hybrid'?'selected':''}>本地知识库 + 联网检索</option></select>
            <label class="contact-check"><input class="checkbox" id="hasContactOnly" type="checkbox" ${state.customerHasContact?'checked':''}> 必有联系方式</label>
            <button class="button primary" data-action="unified-customer-search">${icon('search')}搜索</button>
          </div>
          <div class="source-legend"><span>${icon('database')}本地知识库</span><span>${icon('globe-2')}联网检索</span><small>联系方式包含邮箱、电话、网站及其它通讯方式</small></div>
          <div class="tool-results customer-result-list" id="unifiedCustomerResults">${renderUnifiedCustomerCards()}</div>
        </div>
      </article>

      <article class="panel tool-panel">
        <header class="tool-header"><div><h3>🌍 本地客户发现</h3><p>覆盖欧洲国家及中国，每次各选择一个国家、城市和客户类型，再交由 OpenClaw 分析。</p></div></header>
        <div class="panel-body">
          <div class="tool-form compact discovery-form">
            <select class="select" id="discoveryCountry">${Object.keys(discoveryCities).map(v=>`<option value="${v}" ${state.discoveryCountry===v?'selected':''}>${v}</option>`).join('')}</select>
            <select class="select" id="discoveryCity">${discoveryCities[state.discoveryCountry].map(v=>`<option value="${v}" ${state.discoveryCity===v?'selected':''}>${v}</option>`).join('')}</select>
            <select class="select" id="discoveryType">${[['Distributor','经销商'],['Importer','进口商'],['Dealer','车店'],['Brand','品牌'],['OEM','OEM']].map(([v,l])=>`<option value="${v}" ${state.discoveryType===v?'selected':''}>${l}</option>`).join('')}</select>
            <button class="button primary" data-action="local-discovery-search">${icon('radar')}开始发现</button>
          </div>
          <div class="agent-chain-note"><span class="agent-icon">${icon('bot')}</span><span><strong>OpenClaw · CustomerMeasurementAgent</strong><small>筛选条件 → 公开信息检索 → Agent 分析 → 统一客户字段展示</small></span></div>
          <div class="tool-results customer-result-list" id="localDiscoveryResults">${renderLocalDiscoveryCards()}</div>
        </div>
      </article>
    </section>
  </div>`;
}

function sortedOEMFactories() {
  const query = `${state.oemQuery} ${state.oemCategory === '全部骑行类目' ? '' : state.oemCategory}`.toLowerCase();
  const result = oemFactories.filter(f => {
    const matchesCategory = state.oemCategory === '全部骑行类目' || f.category.includes(state.oemCategory) || (state.oemCategory === '整车' && f.category.includes('整车'));
    const tokens = query.split(/\s+/).filter(token => token.length > 1 && !['oem','台','组','套','个'].includes(token));
    const matchesQuery = !tokens.length || tokens.some(token => `${f.name} ${f.category} ${f.city} ${f.reason}`.toLowerCase().includes(token)) || state.oemQuery.includes('OEM');
    return matchesCategory && matchesQuery;
  });
  if (state.oemSort === 'score') result.sort((a,b)=>b.score-a.score);
  if (state.oemSort === 'capacity') result.sort((a,b)=>parseInt(b.capacity.replace(/\D/g,''))-parseInt(a.capacity.replace(/\D/g,'')));
  if (state.oemSort === 'moq') result.sort((a,b)=>parseInt(a.moq)-parseInt(b.moq));
  if (state.oemSort === 'source') result.sort((a,b)=>a.source.localeCompare(b.source,'zh-CN'));
  return result.slice(0, state.oemTop);
}

function renderOEMCards() {
  return `<div class="match-grid">${sortedOEMFactories().map((f,i)=>`<article class="match-card"><div class="match-rank"><strong>TOP ${i+1}</strong><span>Tier ${f.tier}</span></div><div class="match-copy"><div class="spread"><h4>${f.name}</h4><strong class="match-score">${f.score}<small>分</small></strong></div><p>${f.category} · ${f.city}</p><div class="match-facts"><span>产能 <strong>${f.capacity}</strong></span><span>MOQ <strong>${f.moq}</strong></span></div><div class="match-reason">${f.reason}</div><div class="spread"><span class="badge ${f.source.includes('联网')?'blue':'green'}">${f.source}</span><button class="link-button" data-action="oem-detail" data-name="${f.name}">查看详情</button></div></div></article>`).join('')}</div><div class="tool-result-footer"><span>已统一为工厂、类目、城市、产能、MOQ、匹配分和来源字段</span><button class="button small" data-action="oem-export">${icon('file-down')}导出匹配报告</button></div>`;
}

function renderUnifiedCustomerCards() {
  let rows = unifiedSearchCustomers.filter(r => {
    const text = `${r.name} ${r.country} ${r.type} ${r.business} ${r.contact}`.toLowerCase();
    const tokens = state.customerSearchQuery.toLowerCase().split(/\s+/).filter(token => token.length > 1);
    const matchesQuery = !tokens.length || tokens.some(token => text.includes(token));
    const matchesContact = !state.customerHasContact || Boolean(r.contact);
    const matchesMode = state.customerSearchMode === 'hybrid' || (state.customerSearchMode === 'local' ? r.source === '本地知识库' : r.source.includes('联网检索'));
    return matchesQuery && matchesContact && matchesMode;
  });
  return rows.slice(0,3).map(r=>`<article class="customer-match-row"><span class="match-score">${r.score}<small>分</small></span><span class="customer-match-copy"><strong>${r.name}</strong><small>${r.country} · ${r.type} · ${r.business}</small><em>${icon('contact-round')} ${r.contact}</em></span><span class="badge ${r.source.includes('联网检索')?'blue':'green'}">${r.source}</span><button class="table-icon" data-action="unified-customer-detail" data-name="${r.name}" title="查看详情">${icon('arrow-up-right')}</button></article>`).join('');
}

function renderLocalDiscoveryCards() {
  return localDiscoveryLeads.filter(r => r.country === state.discoveryCountry && r.city === state.discoveryCity && r.type === state.discoveryType).slice(0,3).map(r=>`<article class="customer-match-row"><span class="match-score">${r.score}<small>分</small></span><span class="customer-match-copy"><strong>${r.name}</strong><small>${r.country} · ${r.city} · ${r.type}</small><em>${icon('contact-round')} ${r.contact}</em></span><span class="badge blue">OpenClaw 返回</span><button class="table-icon" data-action="local-lead-detail" data-name="${r.name}" title="查看详情">${icon('arrow-up-right')}</button></article>`).join('') || `<div class="tool-empty">当前筛选暂无候选客户，正式版将使用本地知识库检索后调用 OpenClaw 返回结果。</div>`;
}

function oemDetail(name) {
  const factory = oemFactories.find(f => f.name === name);
  if (!factory) return;
  openDrawer({ title: factory.name, eyebrow: `OEM 匹配 / ${factory.source}`, body: `<div class="spread"><span class="badge green">匹配 ${factory.score} 分</span><span class="secondary-text">原型数据</span></div><div class="detail-grid" style="margin-top:15px">${[['骑行类目',factory.category],['所在城市',factory.city],['供应等级',`Tier ${factory.tier}`],['月产能',factory.capacity],['起订量 MOQ',factory.moq],['数据来源',factory.source],['匹配说明',factory.reason]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${v}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">数据链路</div><div class="agent-chain-note"><span class="agent-icon">${icon('route')}</span><span><strong>本地知识库 + 联网检索（RAG/公开网络采集）→ 字段归一化 → 评分排序</strong><small>RAG 是联网资料的检索方式；页面按本地知识库与联网检索两类来源展示。</small></span></div>` });
}

function unifiedCustomerDetail(name) {
  const customer = unifiedSearchCustomers.find(r => r.name === name);
  if (!customer) return;
  openDrawer({ title: customer.name, eyebrow: `客户统一搜索 / ${customer.source}`, body: `<div class="spread"><span class="badge ${customer.source.includes('联网检索')?'blue':'green'}">${customer.source}</span><span class="secondary-text">匹配 ${customer.score} 分</span></div><div class="detail-grid" style="margin-top:15px">${[['国家',customer.country],['客户类型',customer.type],['业务方向',customer.business],['联系方式',customer.contact],['搜索模式',state.customerSearchMode === 'local'?'本地知识库':state.customerSearchMode === 'rag'?'联网检索':'本地知识库 + 联网检索'],['联系方式过滤',state.customerHasContact?'已启用':'未启用']].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${v}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">返回说明</div><p class="secondary-text" style="line-height:1.8">正式版将返回统一客户字段，并保留每个字段的来源、抓取时间和引用链接；联网检索结果不会覆盖本地客户主档。</p>` });
}

function localLeadDetail(name) {
  const lead = localDiscoveryLeads.find(r => r.name === name);
  if (!lead) return;
  openDrawer({ title: lead.name, eyebrow: '本地客户发现 / OpenClaw', body: `<div class="spread"><span class="badge blue">CustomerMeasurementAgent</span><span class="secondary-text">匹配 ${lead.score} 分</span></div><div class="detail-grid" style="margin-top:15px">${[['国家',lead.country],['城市',lead.city],['客户类型',lead.type],['联系方式',lead.contact],['OpenClaw 返回理由',lead.reason]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${v}</strong></div>`).join('')}</div><div class="agent-chain-note" style="margin-top:16px"><span class="agent-icon">${icon('bot')}</span><span><strong>筛选条件 → 本地/公开信息查询 → CustomerMeasurementAgent</strong><small>正式版返回前会记录查询条件、来源链接和结果时间。</small></span></div>` });
}

function oemExport() {
  toast('匹配报告已生成', `已按 ${state.oemSort === 'score' ? '匹配度' : state.oemSort} 和 Top ${state.oemTop} 生成原型报告；正式版导出格式待确认。`);
}

function offlineUpgradeModal() {
  openModal({ title: '导入离线升级包', eyebrow: '版本升级 / 管理员操作', body: `<input id="upgradeFileInput" type="file" accept=".zip" hidden><div class="upload-zone"><div><span class="upload-icon">${icon('package-open')}</span><h3>选择 STA-100 离线升级包</h3><p>仅支持 .zip；导入后将校验版本、ARM64 架构、签名、完整性、兼容性和可用空间。</p><button class="button primary" data-action="choose-upgrade-package">${icon('file-up')}选择 .zip 升级包</button></div></div><div class="model-warning" style="margin-top:14px"><span>${icon('shield-check')} 安装前自动创建 SQLite、配置、私有文件索引和 Agent 数据快照；失败时保留当前版本。</span></div>`, footer: `<button class="button" data-action="close-modal">取消</button>` });
  document.getElementById('upgradeFileInput').addEventListener('change', e => verifyUpgradePackage(e.target.files[0]));
}

function verifyUpgradePackage(file) {
  if (!file) return;
  const body = document.getElementById('modalBody');
  const footer = document.getElementById('modalFooter');
  const valid = file.name.toLowerCase().endsWith('.zip');
  body.innerHTML = `<div class="upgrade-file"><span class="upload-icon">${icon(valid ? 'file-check-2' : 'file-warning')}</span><div><strong>${file.name}</strong><small>${(file.size / 1024 / 1024).toFixed(2)} MB · 本机文件</small></div></div><div class="upgrade-checks"><div>${icon(valid?'check-circle-2':'circle-x')}<span>文件格式</span><strong>${valid?'通过':'失败'}</strong></div><div>${icon('cpu')}<span>目标架构</span><strong>${valid?'待解析 ARM64':'未检查'}</strong></div><div>${icon('shield-check')}<span>签名与完整性</span><strong>${valid?'待校验':'未检查'}</strong></div><div>${icon('hard-drive')}<span>磁盘空间</span><strong>${valid?'可用 78.4 GB':'未检查'}</strong></div></div><div class="model-warning" style="margin-top:14px"><span>${icon('info')} 原型仅模拟校验结果；正式版必须由 Go 后端解析升级包清单并验证签名。</span></div>`;
  footer.innerHTML = `<button class="button" data-action="close-modal">取消</button>${valid ? `<button class="button primary" data-action="offline-install">${icon('download')}校验并安装</button>` : ''}`;
  applyIcons();
}

function installOfflineUpgrade() {
  const body = document.getElementById('modalBody');
  const footer = document.getElementById('modalFooter');
  body.innerHTML = `<div class="upgrade-progress"><span class="upload-icon">${icon('loader-circle')}</span><h3>正在安装离线升级包</h3><p id="upgradeProgressText">正在创建数据快照...</p><div class="progress"><span id="upgradeProgressBar" style="width:12%"></span></div><small>安装完成后服务将自动重启，页面会短暂不可访问。</small></div>`;
  footer.innerHTML = '';
  applyIcons();
  const steps = ['正在创建数据快照...','正在校验签名和数据库迁移...','正在替换应用文件...','正在执行健康检查...','升级完成，设备将在 3 秒后自动重启...'];
  let index = 0;
  const timer = setInterval(() => {
    index += 1;
    const text = document.getElementById('upgradeProgressText');
    const bar = document.getElementById('upgradeProgressBar');
    if (!text || !bar) return clearInterval(timer);
    text.textContent = steps[Math.min(index, steps.length - 1)];
    bar.style.width = `${Math.min(100, 12 + index * 22)}%`;
    if (index >= steps.length - 1) {
      clearInterval(timer);
      setTimeout(() => { closeModal(); toast('升级流程完成', '原型已模拟自动重启；正式版将由 Go 服务执行进程重启。'); }, 1100);
    }
  }, 650);
}

function renderAgents() {
  const filtered = agents.map((agent, index) => ({ agent, index })).filter(({ agent }) => state.agentCategory === 'all' || agent[2] === state.agentCategory);
  const managedAgents = state.openClawAgents?.filter(agent => !agent.isDefault) || [];
  return `<div class="page-stack">
    ${!state.modelConfigured ? `<div class="model-warning"><span>${icon('triangle-alert')} 当前尚未完成模型配置，智能体不能发起真实调用。</span><button class="button small" data-page="settings">进入模型设置</button></div>` : ''}
    <section class="agents-summary panel">
      <div><h2>🤖 24 个专业智能体，对应 24 个 OpenClaw Agent</h2><p>每个智能体保留独立初始化配置、技能和会话记录，回答区分本地知识库与联网检索来源。</p></div>
      <div class="inline-actions"><button class="button" data-action="agent-manage">${icon('sliders-horizontal')}智能体管理</button><button class="button primary" data-action="weekly-report">${icon('file-clock')}生成本周周报</button></div>
    </section>
    <div class="toolbar">
      <div class="filter-row">
        ${[['all','全部 24'],['trade','贸易与出口 8'],['retail','门店与产品 6'],['market','市场分析 4'],['support','客户与支持 6']].map(([k,l]) => `<button class="filter-chip ${state.agentCategory===k?'active':''}" data-agent-category="${k}">${l}</button>`).join('')}
      </div>
      <span class="spacer"></span>
      <label class="field-search">${icon('search')}<input id="agentSearch" placeholder="搜索智能体名称或能力"></label>
    </div>
    <section class="agent-grid" id="agentGrid">
      ${filtered.map(({ agent: a, index }) => { const live = managedAgents.find(item => item.identityName === a[0]); const available = state.openClawAgents === null || Boolean(live); return `<article class="agent-card" data-agent-search="${a[0]} ${a[1]} ${a[4]}"><div class="agent-top"><span class="agent-icon agent-emoji" aria-hidden="true">${agentEmojis[index]}</span><div class="agent-title"><h3>${a[0]}</h3><span>${a[1]}${live?.model ? ` · ${escapeHTML(live.model)}` : ''}</span></div><span class="status-dot ${state.modelConfigured&&available?'online':'warning'}"></span></div><p>${a[4]}</p><div class="prompt-list">${a[5].map(p => `<button data-action="agent-chat" data-agent="${index}" data-prompt="${p}">${p}</button>`).join('')}</div><div class="agent-foot"><span class="secondary-text">OpenClaw · ${state.modelConfigured&&available?'已注册':'待核验'}</span><button class="link-button" data-action="agent-chat" data-agent="${index}">开始对话</button></div></article>`; }).join('')}
    </section>
  </div>`;
}

function customerRows() {
  const rows = customers.filter(c => {
    if (c.archived) return false;
    const query = state.customerSearch.toLowerCase();
    return (!query || Object.values(c).join(' ').toLowerCase().includes(query)) && (state.customerType === 'all' || c.type === state.customerType) && (state.customerCountry === 'all' || c.country === state.customerCountry);
  });
  return sortRows(rows, state.customerSort, {
    orders: customer => Number(customer.orders || 0),
    total: customer => moneyNumber(customer.total),
    updated: customer => customer.updated,
  });
}

function renderCustomers() {
  const rows = customerRows();
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>👥 客户档案</h2><p>客户是报价单、订单、单据和产品关系的业务中心。</p></div><div class="toolbar"><button class="button" data-action="export-customers">${icon('file-down')}导出 Excel</button><button class="button primary" data-action="new-customer">${icon('plus')}新建客户</button></div></div>
    <div class="toolbar">
      <label class="field-search">${icon('search')}<input id="customerSearch" value="${state.customerSearch}" placeholder="搜索客户、联系人、国家"></label>
      <select class="select" id="customerType"><option value="all">全部类型</option>${['Distributor','Importer','Customer','Reseller','Integrator'].map(v=>`<option ${state.customerType===v?'selected':''}>${v}</option>`).join('')}</select>
      <select class="select" id="customerCountry"><option value="all">全部国家</option>${[...new Set(customers.map(c=>c.country))].map(v=>`<option ${state.customerCountry===v?'selected':''}>${v}</option>`).join('')}</select>
      <span class="spacer"></span><span class="selection-count">已选 ${state.selectedRows.customers.size} 项</span><button class="button ghost" data-action="column-settings">${icon('columns-3')}列表字段</button>
    </div>
    <div class="data-wrap"><table class="data-table"><thead><tr><th>${selectAllCheckbox('customers','客户')}</th><th>客户</th><th>类型</th><th>国家</th><th>联系人</th><th>${sortHeader('订单数量','customer','orders',state.customerSort)}</th><th>${sortHeader('累计金额','customer','total',state.customerSort)}</th><th>评级</th><th>${sortHeader('最近更新','customer','updated',state.customerSort)}</th><th>操作</th></tr></thead><tbody>
      ${rows.map(c=>`<tr><td>${rowCheckbox('customers',c.id,c.name)}</td><td><button class="link-button primary-cell" data-action="customer-detail" data-id="${escapeAttr(c.id)}"><span class="avatar">${escapeHTML(c.name.slice(0,2).toUpperCase())}</span><span><strong>${escapeHTML(c.name)}</strong><small>${escapeHTML(c.id)}</small></span></button></td><td>${escapeHTML(c.type)}</td><td>${escapeHTML(c.country)}</td><td><span class="primary-cell"><span><strong>${escapeHTML(c.contact)}</strong><small>${escapeHTML(c.email)}</small></span></span></td><td>${c.orders}</td><td>${escapeHTML(c.total)}</td><td>${badge(c.rating)}</td><td>${escapeHTML(c.updated)}</td><td><span class="table-actions"><button class="table-icon" data-action="customer-detail" data-id="${escapeAttr(c.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-customer" data-id="${escapeAttr(c.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="customer-more" data-id="${escapeAttr(c.id)}" title="更多">${icon('ellipsis')}</button></span></td></tr>`).join('') || `<tr><td colspan="10"><div class="empty-state">${icon('search-x')}<div><h3>未找到客户</h3><p>请调整搜索词或筛选条件。</p></div></div></td></tr>`}
    </tbody></table></div>
    <div class="pagination"><span>共 ${rows.length} 条记录 · 每页 20 条</span><div><button class="button small ghost" disabled>${icon('chevron-left')}</button><button class="button small" data-action="pagination-current">1</button><button class="button small ghost" disabled>${icon('chevron-right')}</button></div></div>
  </div>`;
}

function renderQuotes() {
  const query = state.quoteSearch.trim().toLowerCase();
  const visible = sortRows(quotes.filter(q => (!query || [q.id, q.subject, q.customer, q.products].join(' ').toLowerCase().includes(query)) && (state.quoteStatus === 'all' || q.status === state.quoteStatus)), state.quoteSort, { value: quote => moneyNumber(quote.value) });
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>📄 报价单管理</h2><p>从客户和产品生成报价，接受后可完整转为订单。</p></div><div class="toolbar"><button class="button" data-action="template-center" data-kind="quote">${icon('layout-template')}模板管理</button><button class="button primary" data-action="new-quote">${icon('plus')}新建报价单</button></div></div>
    <section class="metric-grid" style="grid-template-columns:repeat(4,1fr)">${[['全部报价',42,'files','all'],['草稿',8,'file-pen-line','Draft'],['待客户确认',12,'send','Delivered'],['本月已接受',14,'badge-check','Accepted']].map(([l,v,i,status])=>`<button class="metric-button" data-action="quote-metric-filter" data-status="${status}"><span class="metric-icon">${icon(i)}</span><span><strong class="metric-number">${v}</strong><span class="metric-label">${l}</span></span></button>`).join('')}</section>
    <div class="toolbar"><label class="field-search">${icon('search')}<input id="quoteSearch" value="${escapeAttr(state.quoteSearch)}" placeholder="搜索报价编号、客户或产品"></label><select class="select" id="quoteStatus"><option value="all">全部状态</option>${[['Draft','草稿'],['Delivered','已发送'],['Accepted','已接受'],['Rejected','已拒绝']].map(([v,l])=>`<option value="${v}" ${state.quoteStatus===v?'selected':''}>${l}</option>`).join('')}</select><button class="button ghost" data-action="quote-date-filter">${icon('calendar-days')}有效期</button><span class="result-count">${visible.length} 条 · 已选 ${state.selectedRows.quotes.size} 项</span><span class="spacer"></span><div class="segmented"><button class="${state.quoteView==='table'?'active':''}" data-quote-view="table" title="表格视图">${icon('list')}</button><button class="${state.quoteView==='kanban'?'active':''}" data-quote-view="kanban" title="看板视图">${icon('columns-3')}</button></div></div>
    ${state.quoteView === 'table' ? renderQuoteTable(visible) : renderQuoteKanban(visible)}
  </div>`;
}

function renderQuoteTable(rows=quotes) {
  return `<div class="data-wrap"><table class="data-table"><thead><tr><th>${selectAllCheckbox('quotes','报价单')}</th><th>报价编号</th><th>主题 / 客户</th><th>产品</th><th>${sortHeader('金额','quote','value',state.quoteSort)}</th><th>有效期</th><th>负责人</th><th>状态</th><th>操作</th></tr></thead><tbody>${rows.map(q=>`<tr><td>${rowCheckbox('quotes',q.id,q.id)}</td><td><button class="link-button" data-action="quote-detail" data-id="${escapeAttr(q.id)}">${escapeHTML(q.id)}</button></td><td><span class="primary-cell"><span><strong>${escapeHTML(q.subject)}</strong><small>${escapeHTML(q.customer)}</small></span></span></td><td>${escapeHTML(q.products)}</td><td><strong>${escapeHTML(q.value)}</strong></td><td>${escapeHTML(q.valid)}</td><td>${escapeHTML(q.owner)}</td><td>${badge(q.status)}</td><td><span class="table-actions"><button class="table-icon" data-action="quote-detail" data-id="${escapeAttr(q.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="download-quote" data-id="${escapeAttr(q.id)}" title="下载">${icon('download')}</button><button class="table-icon" data-action="edit-quote" data-id="${escapeAttr(q.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="delete-quote" data-id="${escapeAttr(q.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="9"><div class="empty-state">${icon('search-x')}<div><h3>未找到报价单</h3><p>请调整搜索或状态筛选。</p></div></div></td></tr>`}</tbody></table></div>`;
}

function renderQuoteKanban(rows=quotes) {
  return `<section class="kanban">${[['Draft','草稿'],['Delivered','已发送'],['Accepted','已接受'],['Rejected','已拒绝']].map(([status,label])=>`<div class="kanban-column"><header class="kanban-head"><strong>${label}</strong><span>${rows.filter(q=>q.status===status).length}</span></header><div class="kanban-body">${rows.filter(q=>q.status===status).map(q=>`<button class="kanban-card" data-action="quote-detail" data-id="${escapeAttr(q.id)}"><h4>${escapeHTML(q.subject)}</h4><p>${escapeHTML(q.customer)} · ${escapeHTML(q.id)}</p><span class="amount">${escapeHTML(q.value)}</span></button>`).join('') || `<div class="empty-state" style="min-height:100px;padding:10px"><p>暂无记录</p></div>`}</div></div>`).join('')}</section>`;
}

function renderOrders() {
  const query = state.orderSearch.trim().toLowerCase();
  const visible = sortRows(orders.filter(o => (!query || [o.id, o.customer, o.quote, o.products].join(' ').toLowerCase().includes(query)) && (state.orderStatus === 'all' || o.status === state.orderStatus)), state.orderSort, { value: order => moneyNumber(order.value) });
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>📦 订单生命周期</h2><p>订单可由已接受报价转化，也可手动创建；不直接连接第三方平台下单。</p></div><div class="toolbar"><button class="button" data-action="template-center" data-kind="order">${icon('layout-template')}模板管理</button><button class="button primary" data-action="new-order">${icon('plus')}新建订单</button></div></div>
    <div class="toolbar"><label class="field-search">${icon('search')}<input id="orderSearch" value="${escapeAttr(state.orderSearch)}" placeholder="搜索订单、客户或产品"></label><select class="select" id="orderStatus"><option value="all">全部状态</option>${[['Confirmed','已确认'],['Production','生产中'],['Shipped','已发运'],['Completed','已完成']].map(([v,l])=>`<option value="${v}" ${state.orderStatus===v?'selected':''}>${l}</option>`).join('')}</select><button class="button ghost" data-action="order-date-filter">${icon('calendar-range')}交付日期</button><span class="result-count">${visible.length} 条 · 已选 ${state.selectedRows.orders.size} 项</span><span class="spacer"></span><span class="badge blue">${visible.filter(o=>o.status!=='Completed').length} 个进行中订单</span></div>
    <div class="data-wrap"><table class="data-table"><thead><tr><th>${selectAllCheckbox('orders','订单')}</th><th>订单编号</th><th>客户</th><th>来源报价</th><th>产品</th><th>${sortHeader('金额','order','value',state.orderSort)}</th><th>交付日期</th><th>进度</th><th>状态</th><th>操作</th></tr></thead><tbody>${visible.map(o=>`<tr><td>${rowCheckbox('orders',o.id,o.id)}</td><td><button class="link-button" data-action="order-detail" data-id="${escapeAttr(o.id)}">${escapeHTML(o.id)}</button></td><td>${escapeHTML(o.customer)}</td><td>${escapeHTML(o.quote)}</td><td>${escapeHTML(o.products)}</td><td><strong>${escapeHTML(o.value)}</strong></td><td>${escapeHTML(o.delivery)}</td><td><span style="display:grid;grid-template-columns:70px 30px;align-items:center;gap:5px"><span class="progress"><span style="width:${o.progress}%"></span></span><small>${o.progress}%</small></span></td><td>${badge(o.status)}</td><td><span class="table-actions"><button class="table-icon" data-action="order-detail" data-id="${escapeAttr(o.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-order" data-id="${escapeAttr(o.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="generate-docs" data-id="${escapeAttr(o.id)}" title="生成单据">${icon('files')}</button><button class="table-icon" data-action="delete-order" data-id="${escapeAttr(o.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="10"><div class="empty-state">${icon('search-x')}<div><h3>未找到订单</h3><p>请调整搜索或状态筛选。</p></div></div></td></tr>`}</tbody></table></div>
  </div>`;
}

function renderDocuments() {
  const query = state.documentSearch.trim().toLowerCase();
  const visible = documents.filter(document => {
    const matchesSearch = !query || [document.id, document.type, document.customer, document.order, document.template].join(' ').toLowerCase().includes(query);
    const matchesType = state.documentType === 'all' || document.type === state.documentType;
    const matchesStatus = state.documentStatus === 'all' || document.status === state.documentStatus;
    return matchesSearch && matchesType && matchesStatus;
  });
  const filtering = Boolean(query || state.documentType !== 'all' || state.documentStatus !== 'all');
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>🧾 外贸单据</h2><p>单据从客户、订单和产品数据生成，并保留模板版本和人工确认记录。</p></div><div class="toolbar"><button class="button" data-action="template-center" data-kind="document">${icon('layout-template')}模板管理</button><button class="button" data-action="generate-docs" data-id="all">${icon('copy-plus')}一键生成全套</button><button class="button primary" data-action="new-document">${icon('plus')}生成单据</button></div></div>
    <div class="toolbar document-toolbar">
      <label class="field-search">${icon('search')}<input id="documentSearch" type="search" value="${escapeAttr(state.documentSearch)}" placeholder="搜索编号、客户、订单或模板"></label>
      <select class="select" id="documentStatus"><option value="all">全部状态</option>${[['Draft','草稿'],['Review','待复核'],['Confirmed','已确认']].map(([value,label])=>`<option value="${value}" ${state.documentStatus===value?'selected':''}>${label}</option>`).join('')}</select>
      <div class="filter-row document-type-filter">${[['all','全部类型'],['PI','PI'],['CI','CI'],['PL','PL'],['报关单','报关单']].map(([value,label])=>`<button class="filter-chip ${state.documentType===value?'active':''}" data-document-type="${value}">${label}</button>`).join('')}</div>
      <span class="spacer"></span><span class="result-count">${filtering ? `筛选到 ${visible.length} / ${documents.length} 条` : `共 ${documents.length} 条`}</span>
      ${filtering ? `<button class="button ghost small" data-action="clear-document-filters">${icon('rotate-ccw')}重置</button>` : ''}
    </div>
    <div class="data-wrap"><table class="data-table"><thead><tr><th>单据编号</th><th>类型</th><th>客户</th><th>关联订单</th><th>使用模板</th><th>状态</th><th>最近更新</th><th>操作</th></tr></thead><tbody>${visible.map(d=>`<tr><td><button class="link-button" data-action="document-detail" data-id="${escapeAttr(d.id)}">${escapeHTML(d.id)}</button></td><td><span class="badge neutral">${escapeHTML(d.type)}</span></td><td>${escapeHTML(d.customer)}</td><td>${escapeHTML(d.order)}</td><td>${escapeHTML(d.template)}</td><td>${badge(d.status)}</td><td>${escapeHTML(d.updated)}</td><td><span class="table-actions"><button class="table-icon" data-action="document-detail" data-id="${escapeAttr(d.id)}" title="预览">${icon('eye')}</button><button class="table-icon" data-action="download-document" data-id="${escapeAttr(d.id)}" title="下载 PDF">${icon('file-down')}</button><button class="table-icon" data-action="edit-document" data-id="${escapeAttr(d.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="delete-document" data-id="${escapeAttr(d.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="8"><div class="empty-state">${icon('file-search')}<div><h3>没有符合条件的单据</h3><p>请调整搜索词、单据类型或状态。</p><button class="button small" data-action="clear-document-filters">重置筛选</button></div></div></td></tr>`}</tbody></table></div>
  </div>`;
}

function renderProducts() {
  const visible = productRows();
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>🚲 产品主数据</h2><p>产品被报价、订单和单据引用；停用不会影响已锁定的历史明细。</p></div><div class="toolbar"><button class="button" data-action="import-products">${icon('file-up')}批量导入</button><button class="button primary" data-action="new-product">${icon('plus')}新建产品</button></div></div>
    <div class="toolbar"><label class="field-search">${icon('search')}<input id="productSearch" value="${escapeAttr(state.productSearch)}" placeholder="搜索产品名、编码或 HS CODE"></label><select class="select" id="productCategory"><option value="all">全部类别</option>${['智能设备','智能骑行','整车方案','配件','服务'].map(v=>`<option value="${v}" ${state.productCategory===v?'selected':''}>${v}</option>`).join('')}</select><button class="button ghost" data-action="toggle-product-sort">${icon('arrow-up-down')}${state.productSort==='stockAsc'?'库存升序':'库存降序'}</button><span class="result-count">${visible.length} 条</span><span class="spacer"></span><div class="segmented"><button class="${state.productView==='grid'?'active':''}" data-product-view="grid" title="卡片视图">${icon('grid-2x2')}</button><button class="${state.productView==='table'?'active':''}" data-product-view="table" title="表格视图">${icon('list')}</button></div></div>
    ${state.productView==='grid' ? `<section class="product-grid" id="productGrid">${visible.map((p,i)=>`<article class="product-card" data-product-search="${escapeAttr(`${p.name} ${p.id} ${p.hs}`)}"><button class="product-visual" data-action="product-detail" data-id="${escapeAttr(p.id)}">${icon(i===0?'cpu':i===1?'gauge':'bike')}</button><div class="product-body"><span class="badge ${p.status==='Active'?'green':'amber'}">${p.status==='Active'?'已启用':'待审核'}</span><h3>${escapeHTML(p.name)}</h3><span class="secondary-text">${escapeHTML(p.id)} · HS ${escapeHTML(p.hs)}</span><p>${escapeHTML(p.desc)}</p><div class="product-price"><strong>${escapeHTML(p.price)}</strong><span class="stock">库存 ${p.stock}</span></div></div></article>`).join('') || `<div class="empty-state"><p>未找到产品</p></div>`}</section>` : renderProductTable(visible)}
  </div>`;
}

function productRows() {
  const query = state.productSearch.trim().toLowerCase();
  const rows = products.filter(p => {
    const text = Object.values(p).join(' ').toLowerCase();
    return (!query || text.includes(query)) && (state.productCategory === 'all' || p.category === state.productCategory);
  });
  if (state.productSort === 'stockAsc') rows.sort((a, b) => Number(a.stock) - Number(b.stock));
  if (state.productSort === 'stockDesc') rows.sort((a, b) => Number(b.stock) - Number(a.stock));
  return rows;
}

function renderProductTable(rows=productRows()) {
  return `<div class="data-wrap"><table class="data-table"><thead><tr><th>产品编码</th><th>产品名称</th><th>类别</th><th>HS CODE</th><th>销售价</th><th>库存</th><th>状态</th><th>操作</th></tr></thead><tbody>${rows.map(p=>`<tr><td>${escapeHTML(p.id)}</td><td><button class="link-button" data-action="product-detail" data-id="${escapeAttr(p.id)}">${escapeHTML(p.name)}</button></td><td>${escapeHTML(p.category)}</td><td>${escapeHTML(p.hs)}</td><td><strong>${escapeHTML(p.price)}</strong></td><td>${p.stock}</td><td>${badge(p.status)}</td><td><span class="table-actions"><button class="table-icon" data-action="product-detail" data-id="${escapeAttr(p.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-product" data-id="${escapeAttr(p.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="delete-product" data-id="${escapeAttr(p.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="8"><div class="empty-state"><p>未找到产品</p></div></td></tr>`}</tbody></table></div>`;
}

function supplierRows() {
  const query = state.supplierSearch.trim().toLowerCase();
  const rows = suppliers.filter(s => !query || Object.values(s).join(' ').toLowerCase().includes(query));
  if (state.supplierSort === 'company') rows.sort((a, b) => a.company.localeCompare(b.company, 'zh-CN'));
  if (state.supplierSort === 'updated') rows.sort((a, b) => b.updated.localeCompare(a.updated));
  return rows;
}

function renderSuppliers() {
  const rows = supplierRows();
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>🏭 供应商档案</h2><p>维护供应商联系方式、产品能力、规格、报价和来源信息，后续可与 OEM 匹配结果关联。</p></div><div class="toolbar"><button class="button" data-action="export-suppliers">${icon('file-down')}导出</button><button class="button primary" data-action="new-supplier">${icon('plus')}新建供应商</button></div></div>
    <div class="toolbar"><label class="field-search">${icon('search')}<input id="supplierSearch" value="${escapeAttr(state.supplierSearch)}" placeholder="搜索公司、联系人、产品、来源"></label><select class="select" id="supplierSort"><option value="updated" ${state.supplierSort==='updated'?'selected':''}>最近更新</option><option value="company" ${state.supplierSort==='company'?'selected':''}>公司名称</option></select><span class="result-count">${rows.length} 条</span></div>
    <div class="data-wrap"><table class="data-table"><thead><tr><th>公司</th><th>电话</th><th>联系人</th><th>邮件</th><th>产品</th><th>规格</th><th>报价</th><th>来源</th><th>备注</th><th>操作</th></tr></thead><tbody>${rows.map(s=>`<tr><td><button class="link-button primary-cell" data-action="supplier-detail" data-id="${escapeAttr(s.id)}"><span class="avatar">${escapeHTML(s.company.slice(0,2))}</span><span><strong>${escapeHTML(s.company)}</strong><small>${escapeHTML(s.id)}</small></span></button></td><td>${escapeHTML(s.phone)}</td><td>${escapeHTML(s.contact)}</td><td>${escapeHTML(s.email)}</td><td>${escapeHTML(s.product)}</td><td>${escapeHTML(s.specification)}</td><td><strong>${escapeHTML(s.quote)}</strong></td><td><span class="badge blue">${escapeHTML(s.source)}</span></td><td>${escapeHTML(s.notes)}</td><td><span class="table-actions"><button class="table-icon" data-action="supplier-detail" data-id="${escapeAttr(s.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-supplier" data-id="${escapeAttr(s.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="delete-supplier" data-id="${escapeAttr(s.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="10"><div class="empty-state">${icon('search-x')}<div><h3>未找到供应商</h3><p>请调整搜索条件。</p></div></div></td></tr>`}</tbody></table></div>
  </div>`;
}

function newSupplierForm(supplier) {
  const s=supplier||{}; state.formContext={type:'supplier',id:s.id||''};
  openModal({title:supplier?'编辑供应商':'新建供应商',eyebrow:'供应商档案',body:`<div class="form-grid">${inputField('公司',s.company||'',true,false,'text','supplierCompany')}${inputField('电话',s.phone||'',false,false,'tel','supplierPhone')}${inputField('联系人',s.contact||'',false,false,'text','supplierContact')}${inputField('邮件',s.email||'',false,false,'email','supplierEmail')}${inputField('产品',s.product||'',false,false,'text','supplierProduct')}${inputField('规格',s.specification||'',false,false,'text','supplierSpecification')}${inputField('报价',s.quote||'',false,false,'text','supplierQuote')}${selectField('来源',['展会','电话','朋友介绍','拜访','互联网线索','客户转介绍','其它'],false,'supplierSource',s.source||'其它')}${inputField('备注',s.notes||'',false,true,'text','supplierNotes')}</div>`,footer:formFooter(supplier?'保存修改':'创建供应商','save-supplier')});
}

function saveSupplier() {
  const company=formText('supplierCompany'); if(!company){toast('保存失败','公司名称为必填项。','warning');return;}
  const existing=suppliers.find(s=>s.id===state.formContext?.id); const record=existing||{id:nextRecordId('SUP',suppliers)};
  Object.assign(record,{company,phone:formText('supplierPhone'),contact:formText('supplierContact'),email:formText('supplierEmail'),product:formText('supplierProduct'),specification:formText('supplierSpecification'),quote:formText('supplierQuote'),source:formText('supplierSource'),notes:formText('supplierNotes'),updated:nowText()});
  if(!existing)suppliers.unshift(record); closeModal(); renderPage(); toast(existing?'供应商已更新':'供应商已创建',`${record.company} 已保存。`);
}

function supplierDetail(id) {
  const s=suppliers.find(item=>item.id===id); if(!s)return;
  openDrawer({title:s.company,eyebrow:`供应商 / ${s.id}`,body:`<div class="spread"><span class="badge blue">${escapeHTML(s.source)}</span><div class="inline-actions"><button class="button small" data-action="edit-supplier" data-id="${escapeAttr(s.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-supplier" data-id="${escapeAttr(s.id)}">${icon('trash-2')}删除</button></div></div><div class="detail-grid" style="margin-top:15px">${[['公司',s.company],['电话',s.phone],['联系人',s.contact],['邮件',s.email],['产品',s.product],['规格',s.specification],['报价',s.quote],['来源',s.source],['备注',s.notes],['更新时间',s.updated]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v||'未填写')}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">后续扩展字段</div><p class="secondary-text">可在正式业务库中继续扩展认证、产能、MOQ、合作状态、付款条件和历史报价。</p>`});
}

function deleteSupplier(id) { if(!window.confirm('确定删除该供应商吗？'))return; const index=suppliers.findIndex(s=>s.id===id); if(index>=0)suppliers.splice(index,1); closeDrawer(); renderPage(); toast('供应商已删除','前端示例数据已移除。'); }

function renderDatabase() {
  const categories = [['合同',23,'file-signature'],['报价单',46,'file-chart-column'],['产品手册',128,'book-open'],['法规',85,'scale'],['产品资料',216,'boxes'],['会议记录',32,'notebook-tabs'],['客户资料',174,'contact'],['图片',485,'images'],['其它',29,'folder']];
  const fileQuery = state.fileSearch.trim().toLowerCase();
  const visibleFiles = files.filter(f => !fileQuery || Object.values(f).join(' ').toLowerCase().includes(fileQuery));
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>📚 私有数据库</h2><p>出厂通用资料与客户私有数据分区存储，客户数据仅保存在本机。</p></div><div class="toolbar"><button class="button" data-action="agent-backup">${icon('archive')}备份智能体数据</button><button class="button primary" data-action="upload-file">${icon('upload')}上传文件</button></div></div>
    <section class="metric-grid" style="grid-template-columns:repeat(4,1fr)">${[['文件总数','1,218','files'],['已建立索引','1,164','database-zap'],['待人工确认','17','circle-help'],['本地占用','12.8 GB','hard-drive']].map(([l,v,i])=>`<div class="metric-button"><span class="metric-icon">${icon(i)}</span><span><strong class="metric-number">${v}</strong><span class="metric-label">${l}</span></span></div>`).join('')}</section>
    <div class="section-head"><div><h3>数据分类</h3><p>一个文件可以属于多个标签，分类用于主归档。</p></div><button class="button ghost small" data-action="tag-manage">${icon('tags')}标签管理</button></div>
    <section class="category-grid">${categories.map(([n,c,i])=>`<button class="category-card" data-action="open-category" data-category="${n}"><span class="category-icon">${icon(i)}</span><strong>${n}</strong><span>${c} 个文件</span></button>`).join('')}</section>
    <section class="panel"><header class="panel-head"><div><h3>最近文件</h3><p>展示解析、分类和索引结果</p></div><label class="field-search" style="height:32px">${icon('search')}<input id="fileSearch" value="${escapeAttr(state.fileSearch)}" placeholder="搜索文件或标签"></label></header><div class="data-wrap" style="border:0;border-radius:0"><table class="data-table"><thead><tr><th>文件名</th><th>主分类</th><th>标签</th><th>来源</th><th>大小</th><th>索引状态</th><th>更新时间</th><th>操作</th></tr></thead><tbody>${visibleFiles.map(f=>`<tr><td><button class="link-button" data-action="file-preview" data-name="${escapeAttr(f.name)}">${escapeHTML(f.name)}</button></td><td>${escapeHTML(f.category)}</td><td>${f.tags.map(t=>`<span class="badge neutral">${escapeHTML(t)}</span>`).join(' ')}</td><td>${escapeHTML(f.source)}</td><td>${escapeHTML(f.size)}</td><td>${badge(f.status)}</td><td>${escapeHTML(f.updated)}</td><td><span class="table-actions"><button class="table-icon" data-action="file-preview" data-name="${escapeAttr(f.name)}" title="预览">${icon('eye')}</button><button class="table-icon" data-action="file-edit" data-name="${escapeAttr(f.name)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="file-more" data-name="${escapeAttr(f.name)}" title="更多">${icon('ellipsis')}</button></span></td></tr>`).join('') || `<tr><td colspan="8"><div class="empty-state"><p>未找到匹配文件</p></div></td></tr>`}</tbody></table></div></section>
  </div>`;
}

function renderNews() {
  const filteredNews = state.newsCategory === '全部' ? news : news.filter(item => item.category === state.newsCategory);
  const configuredNews = filteredNews.slice(0, state.newsShowLimit);
  const list = state.newsExpanded ? configuredNews : configuredNews.slice(0,3);
  const feature = filteredNews[0] || news[0];
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>📰 行业新闻</h2><p>由独立推荐智能体按用户关注条件采集、去重和排序。</p></div><div class="toolbar"><button class="button" data-action="news-sources">${icon('rss')}来源管理</button><button class="button primary" data-action="refresh-news">${icon('refresh-cw')}更新新闻</button></div></div>
    <div class="filter-row">${['全部','欧洲市场','法规','智能骑行','渠道','产品'].map(v=>`<button class="filter-chip ${state.newsCategory===v?'active':''}" data-action="news-filter" data-category="${v}">${v}</button>`).join('')}</div>
    <section class="news-layout">
      <button class="news-feature panel" data-action="news-detail" data-title="${escapeAttr(feature.title)}"><div class="news-feature-content"><span class="badge green">${escapeHTML(feature.category)} · 相关度 ${escapeHTML(feature.relevance)}</span><h2>${escapeHTML(feature.title)}</h2><p>${escapeHTML(feature.summary)}</p><div class="source-line" style="margin-top:12px"><span class="mini-source">${escapeHTML(feature.source)}</span><span class="mini-source">${escapeHTML(feature.time)}</span></div></div></button>
      <aside class="panel"><header class="panel-head"><div><h3>关注主题</h3><p>每 ${state.newsFrequency} 获取，单次最多 ${state.newsShowLimit} 条</p></div></header><div class="panel-body"><div class="filter-row">${state.newsTopics.split(/[、,，]/).map(v=>v.trim()).filter(Boolean).slice(0,8).map(v=>`<span class="filter-chip active">${escapeHTML(v)}</span>`).join('')}</div><button class="button ghost small" style="margin-top:12px" data-action="news-sources">${icon('settings-2')}编辑关注条件</button></div></aside>
    </section>
    <section class="panel"><header class="panel-head"><div><h3>最新资讯</h3><p>本地已缓存，可离线查看已获取内容；超过 20 条时在列表内滚动</p></div><span class="meta">更新于 10:05 · 上限 ${state.newsShowLimit}</span></header><div class="news-list ${state.newsExpanded && state.newsShowLimit > 20 ? 'news-list-scroll' : ''}">${list.map(n=>`<button class="news-list-item" data-action="news-detail" data-title="${escapeAttr(n.title)}"><span class="badge neutral">${escapeHTML(n.category)}</span><h3>${escapeHTML(n.title)}</h3><p>${escapeHTML(n.summary)}</p><div class="source-line"><span class="mini-source">${escapeHTML(n.source)}</span><span class="mini-source">${escapeHTML(n.time)}</span><span class="mini-source">相关度 ${escapeHTML(n.relevance)}</span></div></button>`).join('') || `<div class="empty-state"><p>当前分类暂无资讯</p></div>`}</div><div class="panel-body" style="text-align:center"><button class="button ghost small" data-action="toggle-news">${icon(state.newsExpanded?'chevron-up':'chevron-down')}${state.newsExpanded?'收起':'MORE'}</button></div></section>
  </div>`;
}

function openNewsSettings() {
  const frequencies = ['1小时','2小时','3小时','6小时','8小时','12小时','24小时'];
  openModal({title:'新闻与推荐设置',eyebrow:'行业新闻',body:`<div class="form-grid"><div class="form-field full"><label>关注国家 <span class="required">*</span></label><input class="input" id="newsCountries" value="${escapeAttr(state.newsCountries)}"><small>多个国家使用顿号或逗号分隔，最多 200 个字符。</small></div><div class="form-field full"><label>关注主题 <span class="required">*</span></label><input class="input" id="newsTopics" value="${escapeAttr(state.newsTopics)}"><small>多个主题使用顿号或逗号分隔，最多 200 个字符。</small></div><div class="form-field"><label>每次展示数量 <span class="required">*</span></label><input class="input" id="newsShowLimit" type="number" min="1" max="100" step="1" inputmode="numeric" value="${state.newsShowLimit}"><small>必须是 1-100 的整数。</small></div>${selectField('获取频率',frequencies,false,'newsFrequency',state.newsFrequency)}<div class="form-field full"><label>指定来源 <span class="required">*</span></label><textarea class="textarea" id="newsSources" placeholder="每行一个来源">${escapeHTML(state.newsSources)}</textarea><small>最多 1,000 个字符，每行最多 200 个字符。</small></div><div class="form-field full"><div class="display-rule-note">${icon('layout-list')}<span><strong>展示规则</strong><small>概览固定显示相关度最高的 3 条；MORE 进入行业新闻页。行业新闻页纵向展示本次获取结果，超过 20 条时列表内部滚动，不会无限撑高页面；单次最多展示 100 条。</small></span></div></div></div>`,footer:formFooter('保存设置','save-news-settings')});
}

function invalidNewsField(input, message) {
  input?.setCustomValidity(message);
  input?.reportValidity();
  input?.focus();
  toast('无法保存新闻设置',message,'warning');
}

function saveNewsSettings() {
  const countries = document.getElementById('newsCountries');
  const topics = document.getElementById('newsTopics');
  const limitInput = document.getElementById('newsShowLimit');
  const frequency = document.getElementById('newsFrequency');
  const sourcesInput = document.getElementById('newsSources');
  [countries,topics,limitInput,frequency,sourcesInput].forEach(input=>input?.setCustomValidity(''));
  const countriesValue = countries?.value.trim() || '';
  const topicsValue = topics?.value.trim() || '';
  const sourcesValue = sourcesInput?.value.trim() || '';
  const limitText = limitInput?.value.trim() || '';
  const limit = Number(limitText);
  const allowedFrequencies = ['1小时','2小时','3小时','6小时','8小时','12小时','24小时'];
  if (!countriesValue) return invalidNewsField(countries,'请填写至少一个关注国家。');
  if (countriesValue.length > 200) return invalidNewsField(countries,'关注国家最多 200 个字符。');
  if (!topicsValue) return invalidNewsField(topics,'请填写至少一个关注主题。');
  if (topicsValue.length > 200) return invalidNewsField(topics,'关注主题最多 200 个字符。');
  if (!/^\d+$/.test(limitText) || !Number.isInteger(limit) || limit < 1 || limit > 100) return invalidNewsField(limitInput,'每次展示数量必须是 1-100 的整数。');
  if (!allowedFrequencies.includes(frequency?.value || '')) return invalidNewsField(frequency,'请选择允许的获取频率。');
  if (!sourcesValue) return invalidNewsField(sourcesInput,'请填写至少一个新闻来源。');
  if (sourcesValue.length > 1000) return invalidNewsField(sourcesInput,'指定来源最多 1,000 个字符。');
  if (sourcesValue.split(/\r?\n/).some(line=>line.length > 200)) return invalidNewsField(sourcesInput,'指定来源每行最多 200 个字符。');
  state.newsCountries = countriesValue;
  state.newsTopics = topicsValue;
  state.newsShowLimit = limit;
  state.newsFrequency = frequency.value;
  state.newsSources = sourcesValue;
  localStorage.setItem('sta100-news-settings',JSON.stringify({countries:countriesValue,topics:topicsValue,showLimit:limit,frequency:frequency.value,sources:sourcesValue}));
  closeModal();
  renderPage();
  toast('新闻设置已保存',`每 ${state.newsFrequency} 获取，单次最多展示 ${state.newsShowLimit} 条。`);
}

function renderSettings() {
  const tabs = [['model','模型设置','brain-circuit'],['plugins','插件绑定','plug-zap'],['scheduler','定时任务','clock-3'],['backup','智能体备份','archive'],['security','数据安全','shield-check'],['system','系统信息','monitor-cog'],['upgrade','版本升级','package-open']];
  return `<div class="settings-layout">
    <nav class="settings-nav panel">${tabs.map(([k,l,i])=>`<button class="${state.settingsTab===k?'active':''}" data-settings-tab="${k}">${icon(i)}${l}</button>`).join('')}</nav>
    <div class="settings-content">${renderSettingsContent()}</div>
  </div>`;
}

function renderSettingsContent() {
  const content = {
    model: renderModelSettings(),
    plugins: `<section class="panel"><header class="panel-head"><div><h3>插件绑定</h3><p>首期提供微信和飞书，能力与 OpenClaw 插件绑定方式保持一致</p></div></header>${[['微信','消息推送与交互入口','message-circle','未绑定'],['飞书','消息、文档和通知能力','send','已绑定']].map(([n,d,i,s])=>`<div class="setting-row"><span class="setting-icon">${icon(i)}</span><div class="setting-copy"><strong>${n}</strong><span>${d}</span></div>${s==='已绑定'?badge('Active'):'<span class="badge neutral">未绑定</span>'}<button class="button small" data-action="bind-plugin" data-plugin="${n}">${s==='已绑定'?'管理':'绑定'}</button></div>`).join('')}</section>`,
    scheduler: `<section class="panel"><header class="panel-head"><div><h3>定时任务</h3><p>内置任务不能删除，但可以编辑内容、频率和开关</p></div><button class="button small" data-action="new-schedule">${icon('plus')}新增任务</button></header>${[['每日推荐更新',`每 ${state.newsFrequency}`,'10:45','Active'],['智能体周报','每周五 18:00','2026-08-14 18:00','Active'],['行业新闻更新',`每 ${state.newsFrequency}`,'14:00','Active'],['数据索引维护','每天 02:30','2026-08-11 02:30','Draft']].map(([n,f,next,s])=>`<div class="setting-row"><span class="setting-icon">${icon('timer-reset')}</span><div class="setting-copy"><strong>${n}</strong><span>${f} · 下次运行 ${next}</span></div>${badge(s)}<button class="table-icon" data-action="edit-schedule" data-name="${n}" title="编辑">${icon('pencil')}</button></div>`).join('')}</section>`,
    backup: `<section class="panel"><header class="panel-head"><div><h3>智能体数据备份</h3><p>仅备份智能体初始化配置、技能、会话和用户操作记录</p></div></header><div class="setting-row"><span class="setting-icon">${icon('folder-cog')}</span><div class="setting-copy"><strong>备份目录</strong><span>/mnt/sta100-backup/agents · 可选择外置存储目录</span></div><button class="button small" data-action="choose-backup">选择目录</button></div><div class="setting-row"><span class="setting-icon">${icon('archive')}</span><div class="setting-copy"><strong>最近备份</strong><span>2026-08-09 23:30 · 1.42 GB · 校验成功</span></div><button class="button primary small" data-action="agent-backup">立即备份</button></div></section>`,
    security: `<section class="panel"><header class="panel-head"><div><h3>数据安全</h3><p>客户私有数据归客户所有，并保存在设备本地</p></div></header><div class="setting-row"><span class="setting-icon">${icon('user-round-cog')}</span><div class="setting-copy"><strong>本机登录账户</strong><span>当前用户名：${escapeHTML(authState.username)} · 密码以哈希形式保存在本机浏览器</span></div><button class="button small" data-action="account-settings">${icon('key-round')}修改用户名和密码</button></div>${[['本地存储区','客户上传文件、业务数据库和索引均存储在本机','hard-drive'],['敏感信息保护','密钥加密保存；列表和日志不显示完整密钥','lock-keyhole'],['联网调用边界','仅把完成当前任务所需的最小内容发送给已配置模型','network'],['操作记录','记录关键增删改、导出、模型和升级操作','scroll-text']].map(([n,d,i])=>`<div class="setting-row"><span class="setting-icon">${icon(i)}</span><div class="setting-copy"><strong>${n}</strong><span>${d}</span></div>${badge('Active')}</div>`).join('')}</section>`,
    system: renderSystemSettings(),
    upgrade: `<section class="panel"><header class="panel-head"><div><h3>版本升级</h3><p>仅支持管理员手动导入离线升级包，不启用在线热升级</p></div><span class="badge neutral">离线升级</span></header><div class="setting-row"><span class="setting-icon">${icon('package-check')}</span><div class="setting-copy"><strong>当前版本</strong><span>1.0.0 · ARM64 · 构建 20260810</span></div><span class="setting-value">运行正常</span></div><div class="setting-row"><span class="setting-icon">${icon('shield-check')}</span><div class="setting-copy"><strong>升级保护</strong><span>安装前校验签名、版本、架构和磁盘空间，并自动创建业务数据、配置及 Agent 数据快照</span></div>${badge('Active')}</div><div class="setting-row"><span class="setting-icon">${icon('upload')}</span><div class="setting-copy"><strong>导入离线升级包</strong><span>选择本机 .zip 包，校验通过并经管理员确认后安装；完成后设备自动重启</span></div><button class="button primary small" data-action="offline-upgrade">${icon('upload')}选择升级包</button></div><div class="setting-row"><span class="setting-icon">${icon('history')}</span><div class="setting-copy"><strong>最近升级记录</strong><span>暂无升级记录 · 日志和旧版本将保留用于失败回滚</span></div><button class="button small" data-action="upgrade-history">查看记录</button></div></section>`,
  };
  return content[state.settingsTab];
}

function renderModelSettings() {
  const data = state.openClawModels;
  const current = data?.resolvedDefault || data?.defaultModel || '正在读取 OpenClaw 配置';
  const providers = data?.providers || [];
  const configuredProviders = providers.filter(provider => provider.configured).map(provider => provider.provider);
  const providerText = configuredProviders.length ? configuredProviders.join('、') : '尚未读取到可用凭据';
  const modelCount = data?.models?.filter(model => model.available && !model.missing).length || 0;
  return `<section class="panel"><header class="panel-head"><div><h3>🤖 模型设置</h3><p>直接读取并编排本机 OpenClaw 的模型和凭据配置</p></div>${state.openClawModelsLoading ? '<span class="badge neutral">读取中</span>' : badge(data?.configured?'Active':'Review')}</header>
    ${data?.error ? `<div class="model-warning"><span>${icon('triangle-alert')} ${escapeHTML(data.error)}</span><button class="button small" data-action="refresh-openclaw-models">重试</button></div>` : ''}
    <div class="setting-row"><span class="setting-icon">${icon('cloud-cog')}</span><div class="setting-copy"><strong>当前默认模型</strong><span>由 OpenClaw 解析，未单独指定模型的 Agent 继承该值</span></div><span class="setting-value">${escapeHTML(current)}</span><button class="button small" data-action="configure-model" ${!data?'disabled':''}>配置</button></div>
    <div class="setting-row"><span class="setting-icon">${icon('key-round')}</span><div class="setting-copy"><strong>模型凭据</strong><span>通过标准输入写入 OpenClaw 凭据库；页面和接口均不回显密钥</span></div><span class="setting-value">${escapeHTML(providerText)}</span><button class="button small" data-action="configure-model" ${!data?'disabled':''}>更新</button></div>
    <div class="setting-row"><span class="setting-icon">${icon('list-checks')}</span><div class="setting-copy"><strong>可用模型</strong><span>来自 OpenClaw models list，不使用前端预置结果</span></div><span class="setting-value">${state.openClawModelsLoading?'读取中':`${modelCount} 个`}</span><button class="button small" data-action="refresh-openclaw-models">${icon('refresh-cw')}刷新</button></div>
  </section>`;
}

function renderSystemSettings() {
  const status = state.openClawStatus;
  const managedCount = state.openClawAgents?.filter(agent => !agent.isDefault).length;
  const openClawValue = state.openClawStatusLoading ? '正在读取服务状态' : status?.available ? `${status.version} · ${status.serviceStatus} · RPC ${status.rpcOK?'正常':'异常'}` : status?.error || '暂未读取';
  const agentValue = state.openClawAgentsLoading ? '正在读取 Agent' : managedCount === undefined ? '暂未读取' : `${managedCount} 个 STA-100 Agent 已注册`;
  const rows = [
    ['设备型号','STA-100 / ARM64','cpu',true],
    ['系统版本','STA-100 OS 1.0.0','monitor',true],
    ['OpenClaw',openClawValue,'bot',Boolean(status?.rpcOK)],
    ['Agent 编排',agentValue,'blocks',managedCount===24],
    ['本地数据库','SQLite · 正常 · 286 MB','database',true],
    ['存储空间','已用 12.8 GB / 可用 78.4 GB','hard-drive',true],
  ];
  return `<section class="panel"><header class="panel-head"><div><h3>⚙️ 系统信息</h3><p>硬件、服务和本地组件运行状态</p></div><button class="button small" data-action="refresh-openclaw-system">${icon('refresh-cw')}刷新</button></header>${rows.map(([name,value,iconName,ok])=>`<div class="setting-row"><span class="setting-icon">${icon(iconName)}</span><div class="setting-copy"><strong>${name}</strong><span>${escapeHTML(value)}</span></div>${badge(ok?'Active':'Review')}</div>`).join('')}</section>`;
}

function wirePageSpecific() {
  const sub = document.getElementById('subscriptionToggle');
  if (sub) sub.addEventListener('change', e => { state.subscription = e.target.checked; toast('订阅设置已更新', state.subscription ? `系统将每 ${state.newsFrequency} 更新一次推荐。` : '自动更新已暂停。'); renderPage(); });
  const agentSearch = document.getElementById('agentSearch');
  if (agentSearch) agentSearch.addEventListener('input', e => filterCards('#agentGrid .agent-card', 'agentSearch', e.target.value));
  const customerSearch = document.getElementById('customerSearch');
  if (customerSearch) customerSearch.addEventListener('input', e => { state.customerSearch = e.target.value; renderPage(); requestAnimationFrame(()=>{ const i=document.getElementById('customerSearch'); if(i){i.focus();i.setSelectionRange(i.value.length,i.value.length);} }); });
  const customerType = document.getElementById('customerType');
  if (customerType) customerType.addEventListener('change', e => { state.customerType=e.target.value; renderPage(); });
  const customerCountry = document.getElementById('customerCountry');
  if (customerCountry) customerCountry.addEventListener('change', e => { state.customerCountry=e.target.value; renderPage(); });
  const quoteSearch = document.getElementById('quoteSearch');
  if (quoteSearch) quoteSearch.addEventListener('input', e => { state.quoteSearch=e.target.value; renderPage(); requestAnimationFrame(()=>{ const input=document.getElementById('quoteSearch'); if(input){input.focus();input.setSelectionRange(input.value.length,input.value.length);} }); });
  const quoteStatus = document.getElementById('quoteStatus');
  if (quoteStatus) quoteStatus.addEventListener('change', e => { state.quoteStatus=e.target.value; renderPage(); });
  const orderSearch = document.getElementById('orderSearch');
  if (orderSearch) orderSearch.addEventListener('input', e => { state.orderSearch=e.target.value; renderPage(); requestAnimationFrame(()=>{ const input=document.getElementById('orderSearch'); if(input){input.focus();input.setSelectionRange(input.value.length,input.value.length);} }); });
  const orderStatus = document.getElementById('orderStatus');
  if (orderStatus) orderStatus.addEventListener('change', e => { state.orderStatus=e.target.value; renderPage(); });
  const documentSearch = document.getElementById('documentSearch');
  if (documentSearch) documentSearch.addEventListener('input', e => { state.documentSearch=e.target.value; renderPage(); requestAnimationFrame(()=>{ const input=document.getElementById('documentSearch'); if(input){input.focus();input.setSelectionRange(input.value.length,input.value.length);} }); });
  const documentStatus = document.getElementById('documentStatus');
  if (documentStatus) documentStatus.addEventListener('change', e => { state.documentStatus=e.target.value; renderPage(); });
  const productSearch = document.getElementById('productSearch');
  if (productSearch) productSearch.addEventListener('input', e => { state.productSearch=e.target.value; renderPage(); requestAnimationFrame(()=>{ const input=document.getElementById('productSearch'); if(input){input.focus();input.setSelectionRange(input.value.length,input.value.length);} }); });
  const productCategory = document.getElementById('productCategory');
  if (productCategory) productCategory.addEventListener('change', e => { state.productCategory=e.target.value; renderPage(); });
  const oemCategory = document.getElementById('oemCategory');
  if (oemCategory) oemCategory.addEventListener('change', e => { state.oemCategory = e.target.value; renderPage(); });
  const oemSort = document.getElementById('oemSort');
  if (oemSort) oemSort.addEventListener('change', e => { state.oemSort = e.target.value; renderPage(); });
  const oemTop = document.getElementById('oemTop');
  if (oemTop) oemTop.addEventListener('change', e => { state.oemTop = Number(e.target.value); renderPage(); });
  const oemQuery = document.getElementById('oemQuery');
  if (oemQuery) oemQuery.addEventListener('input', e => { state.oemQuery = e.target.value; });
  const unifiedMode = document.getElementById('unifiedSearchMode');
  if (unifiedMode) unifiedMode.addEventListener('change', e => { state.customerSearchMode = e.target.value; renderPage(); });
  const discoveryCountry = document.getElementById('discoveryCountry');
  if (discoveryCountry) discoveryCountry.addEventListener('change', e => { state.discoveryCountry = e.target.value; state.discoveryCity = discoveryCities[state.discoveryCountry][0]; renderPage(); });
  const discoveryCity = document.getElementById('discoveryCity');
  if (discoveryCity) discoveryCity.addEventListener('change', e => { state.discoveryCity = e.target.value; renderPage(); });
  const discoveryType = document.getElementById('discoveryType');
  if (discoveryType) discoveryType.addEventListener('change', e => { state.discoveryType = e.target.value; renderPage(); });
  const supplierSearch = document.getElementById('supplierSearch');
  if (supplierSearch) supplierSearch.addEventListener('input', e => { state.supplierSearch=e.target.value; renderPage(); requestAnimationFrame(()=>{ const input=document.getElementById('supplierSearch'); if(input){input.focus();input.setSelectionRange(input.value.length,input.value.length);} }); });
  const supplierSort = document.getElementById('supplierSort');
  if (supplierSort) supplierSort.addEventListener('change', e => { state.supplierSort=e.target.value; renderPage(); });
  const fileSearch = document.getElementById('fileSearch');
  if (fileSearch) fileSearch.addEventListener('input', e => { state.fileSearch=e.target.value; renderPage(); requestAnimationFrame(()=>{ const input=document.getElementById('fileSearch'); if(input){input.focus();input.setSelectionRange(input.value.length,input.value.length);} }); });
  document.querySelectorAll('[data-select-all]').forEach(header => {
    const module = header.dataset.module;
    const visible = [...document.querySelectorAll(`[data-row-select][data-module="${module}"]`)];
    const selectedCount = visible.filter(input => state.selectedRows[module].has(input.dataset.id)).length;
    header.checked = visible.length > 0 && selectedCount === visible.length;
    header.indeterminate = selectedCount > 0 && selectedCount < visible.length;
    header.addEventListener('change', () => {
      visible.forEach(input => header.checked ? state.selectedRows[module].add(input.dataset.id) : state.selectedRows[module].delete(input.dataset.id));
      renderPage();
    });
  });
  document.querySelectorAll('[data-row-select]').forEach(input => input.addEventListener('change', () => {
    const selected = state.selectedRows[input.dataset.module];
    input.checked ? selected.add(input.dataset.id) : selected.delete(input.dataset.id);
    renderPage();
  }));
}

function filterCards(selector, dataKey, query) {
  document.querySelectorAll(selector).forEach(card => card.hidden = !card.dataset[dataKey].toLowerCase().includes(query.toLowerCase()));
}

function openModal({ title, eyebrow='STA-100', body, footer='', wide=false }) {
  if (!document.getElementById('drawerBackdrop').hidden) closeDrawer();
  document.getElementById('modalTitle').textContent = title;
  document.getElementById('modalEyebrow').textContent = eyebrow;
  document.getElementById('modalBody').innerHTML = body;
  document.getElementById('modalFooter').innerHTML = footer;
  document.getElementById('modal').classList.toggle('wide', wide);
  document.getElementById('modalBackdrop').hidden = false;
  syncOverlayScroll();
  applyIcons();
}

function syncOverlayScroll() {
  const modalOpen = !document.getElementById('modalBackdrop').hidden;
  const drawerOpen = !document.getElementById('drawerBackdrop').hidden;
  const commandOpen = !document.getElementById('commandBackdrop').hidden;
  document.body.style.overflow = modalOpen || drawerOpen || commandOpen ? 'hidden' : '';
}
function closeModal() { document.getElementById('modalBackdrop').hidden = true; syncOverlayScroll(); }
function openDrawer({ title, eyebrow='记录详情', body }) {
  document.getElementById('drawerTitle').textContent = title;
  document.getElementById('drawerEyebrow').textContent = eyebrow;
  document.getElementById('drawerBody').innerHTML = body;
  document.getElementById('drawerBackdrop').hidden = false;
  syncOverlayScroll();
  applyIcons();
}
function closeDrawer() { document.getElementById('drawerBackdrop').hidden = true; syncOverlayScroll(); }

function formFooter(primary='保存', action='save-form') { return `<button class="button ghost" data-action="close-modal">取消</button><button class="button primary" data-action="${action}">${icon('check')}${primary}</button>`; }
function inputField(label, value='', required=false, full=false, type='text', id='') { return `<div class="form-field ${full?'full':''}"><label>${label}${required?' <span class="required">*</span>':''}</label><input class="input" ${id?`id="${id}"`:''} type="${type}" value="${escapeAttr(value)}"></div>`; }
function selectField(label, options, full=false, id='', selected='') { return `<div class="form-field ${full?'full':''}"><label>${label}</label><select class="select" ${id?`id="${id}"`:''}>${options.map(v=>`<option ${String(v)===String(selected)?'selected':''}>${escapeHTML(v)}</option>`).join('')}</select></div>`; }
function relationField(label, target, items, selected='') {
  const values = items.map(item => typeof item === 'string' ? item : item.label);
  return `<div class="form-field relation-field"><label>${label}</label><div class="relation-picker"><input class="input relation-input" id="${target}" data-relation-input="${target}" value="${escapeAttr(selected)}" autocomplete="off" placeholder="输入名称进行模糊匹配"><div class="relation-options" id="${target}Options">${relationOptions(target, values, selected)}</div></div></div>`;
}
function relationOptions(target, values, query='') {
  const needle = String(query || '').trim().toLowerCase();
  const matches = values.filter(value => !needle || value.toLowerCase().includes(needle)).slice(0, 8);
  return matches.map(value => `<button type="button" class="relation-option" data-action="relation-select" data-target="${target}" data-value="${escapeAttr(value)}">${escapeHTML(value)}</button>`).join('') || `<span class="relation-empty">未找到匹配项</span>`;
}
function formText(id) { return document.getElementById(id)?.value.trim() || ''; }
function formNumber(id) { return Number(document.getElementById(id)?.value || 0); }
function nowText() { return new Date().toISOString().slice(0, 16).replace('T', ' '); }
function nextRecordId(prefix, list) { return `${prefix}-${String(list.length + 1).padStart(4, '0')}`; }
function moneyNumber(value) { const number = Number(String(value || '').replace(/[^0-9.-]/g, '')); return Number.isFinite(number) ? number : 0; }
function formatMoney(value, currency='EUR') { return `${currency} ${Number(value || 0).toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`; }
function sortHeader(label, module, field, sort) {
  const active = sort.field === field;
  const direction = active ? sort.direction : '';
  const iconName = direction === 'asc' ? 'arrow-up' : direction === 'desc' ? 'arrow-down' : 'arrow-up-down';
  const directionText = direction === 'asc' ? '升序' : direction === 'desc' ? '降序' : '未排序';
  return `<button type="button" class="sort-button ${active?'active':''}" data-action="sort-table" data-module="${module}" data-field="${field}" title="${escapeAttr(label)}：${directionText}" aria-label="${escapeAttr(label)}，当前${directionText}"><span>${escapeHTML(label)}</span>${icon(iconName)}</button>`;
}
function sortRows(rows, sort, fields) {
  const getter = fields[sort.field];
  if (!getter) return rows;
  const direction = sort.direction === 'asc' ? 1 : -1;
  return rows.sort((a, b) => {
    const left = getter(a); const right = getter(b);
    if (typeof left === 'number' && typeof right === 'number') return (left - right) * direction;
    return String(left ?? '').localeCompare(String(right ?? ''), 'zh-CN', { numeric: true }) * direction;
  });
}
function rowCheckbox(module, id, label) {
  const checked = state.selectedRows[module]?.has(id) ? 'checked' : '';
  return `<input class="checkbox" type="checkbox" data-row-select data-module="${module}" data-id="${escapeAttr(id)}" aria-label="选择 ${escapeAttr(label)}" ${checked}>`;
}
function selectAllCheckbox(module, label) {
  return `<input class="checkbox" type="checkbox" data-select-all data-module="${module}" aria-label="全选当前${escapeAttr(label)}">`;
}
function toggleTableSort(module, field) {
  const key = `${module}Sort`;
  const current = state[key];
  state[key] = { field, direction: current.field === field && current.direction === 'desc' ? 'asc' : 'desc' };
  renderPage();
}
function productByID(id) { return products.find(product => product.id === id); }
function businessLines(record) {
  if (record?.lines?.length) return record.lines.map(line => ({ ...line }));
  const summary = String(record?.products || '');
  const product = products.find(item => summary.includes(item.id) || summary.includes(item.name) || (summary.includes('STA-100') && item.id === 'STA-100-EU')) || products.find(item => item.status === 'Active') || products[0];
  const quantity = Math.max(1, Number(summary.match(/[x×]\s*(\d+)/i)?.[1] || 1));
  const total = moneyNumber(record?.value);
  return product ? [{ productId: product.id, quantity, unitPrice: total > 0 ? total / quantity : moneyNumber(product.price), discount: 0 }] : [];
}
function lineProduct(line) { return productByID(line.productId) || products[0]; }
function lineSubtotal(line, withDiscount=false) {
  const discount = withDiscount ? Math.min(100, Math.max(0, Number(line.discount || 0))) : 0;
  return Math.max(0, Number(line.quantity || 0)) * Math.max(0, Number(line.unitPrice || 0)) * (1 - discount / 100);
}
function lineSummary(lines) {
  return lines.map(line => `${lineProduct(line)?.name || line.productId} x ${Number(line.quantity || 0)}`).join('；');
}
function detailLinesTable(lines, compact=false) {
  const rows = lines.map(line => `<tr><td>${escapeHTML(line.productName || lineProduct(line)?.name || line.productId)}</td><td>${Number(line.quantity)}</td><td>${formatMoney(line.unitPrice)}</td><td><strong>${formatMoney(line.amount ?? lineSubtotal(line))}</strong></td></tr>`).join('');
  return `<div class="data-wrap"><table class="data-table business-lines ${compact?'compact-lines':''}"><thead><tr><th>产品</th><th>数量</th><th>单价</th><th>小计</th></tr></thead><tbody>${rows}</tbody></table></div>`;
}
function productOptions(selected) {
  return products.filter(product => product.status === 'Active' || product.id === selected).map(product => `<option value="${escapeAttr(product.id)}" ${product.id===selected?'selected':''}>${escapeHTML(product.name)} (${escapeHTML(product.id)}) · 库存 ${product.stock}</option>`).join('');
}
function quoteDraftTotal() {
  return state.quoteDraftLines.reduce((total, line) => total + lineSubtotal(line, true), 0) + formNumber('quoteFreight') + formNumber('quoteTax');
}
function renderQuoteDraftLines() {
  const body = document.getElementById('quoteLinesBody');
  if (!body) return;
  body.innerHTML = state.quoteDraftLines.map((line, index) => {
    const product = lineProduct(line);
    return `<tr><td><select class="select line-product" data-quote-line-field="productId" data-index="${index}">${productOptions(line.productId)}</select><small class="line-stock">可用库存 ${product?.stock ?? 0}</small></td><td><input class="input line-number" type="number" min="1" max="${product?.stock ?? 1}" value="${line.quantity}" data-quote-line-field="quantity" data-index="${index}"></td><td><input class="input line-price" type="number" min="0" step="0.01" value="${Number(line.unitPrice).toFixed(2)}" data-quote-line-field="unitPrice" data-index="${index}"></td><td><input class="input line-number" type="number" min="0" max="100" step="0.1" value="${Number(line.discount||0)}" data-quote-line-field="discount" data-index="${index}"></td><td><strong data-quote-line-subtotal="${index}">${formatMoney(lineSubtotal(line,true))}</strong></td><td><button type="button" class="table-icon" data-action="remove-quote-line" data-index="${index}" title="删除明细">${icon('trash-2')}</button></td></tr>`;
  }).join('');
  const total = document.getElementById('quoteDraftTotal'); if (total) total.textContent = formatMoney(quoteDraftTotal());
  applyIcons();
}
function orderDraftTotal() { return state.orderDraftLines.reduce((total, line) => total + lineSubtotal(line), 0); }
function renderOrderDraftLines() {
  const body = document.getElementById('orderLinesBody');
  if (!body) return;
  body.innerHTML = state.orderDraftLines.map((line, index) => {
    const product = lineProduct(line);
    return `<tr><td><select class="select line-product" data-order-line-field="productId" data-index="${index}">${productOptions(line.productId)}</select><small class="line-stock">可用库存 ${product?.stock ?? 0}</small></td><td><input class="input line-number" type="number" min="1" max="${product?.stock ?? 1}" value="${line.quantity}" data-order-line-field="quantity" data-index="${index}"></td><td><input class="input line-price" type="number" min="0" step="0.01" value="${Number(line.unitPrice).toFixed(2)}" data-order-line-field="unitPrice" data-index="${index}"></td><td><strong data-order-line-subtotal="${index}">${formatMoney(lineSubtotal(line))}</strong></td><td><button type="button" class="table-icon" data-action="remove-order-line" data-index="${index}" title="删除明细">${icon('trash-2')}</button></td></tr>`;
  }).join('');
  const total = document.getElementById('orderDraftTotal'); if (total) total.textContent = formatMoney(orderDraftTotal());
  applyIcons();
}
function quoteLinesToOrder(lines) {
  return lines.map(line => ({ productId: line.productId, quantity: Number(line.quantity), unitPrice: Number(line.unitPrice) * (1 - Number(line.discount || 0) / 100) }));
}
function syncOrderFromQuote(value) {
  const quote = quotes.find(item => item.id === String(value || '').split(' · ')[0]);
  if (!quote) return;
  state.orderDraftLines = quoteLinesToOrder(businessLines(quote));
  const customerInput = document.getElementById('orderCustomer'); if (customerInput) customerInput.value = quote.customer;
  const source = document.getElementById('orderSource'); if (source) source.value = '从已接受报价创建';
  renderOrderDraftLines();
}
function downloadText(filename, content) {
  const blob = new Blob([content], { type: 'text/plain;charset=utf-8' });
  const link = document.createElement('a'); link.href = URL.createObjectURL(blob); link.download = filename; link.click();
  setTimeout(() => URL.revokeObjectURL(link.href), 1000);
}

function showMetric(key) {
  const m=metrics.find(x=>x.key===key);
  const items={tasks:['检查 VeloTrade 报价有效期','复核 2 份产品资料索引','确认 SIM 订单交付计划'],meetings:['09:30 内部产品会议 · 已完成','14:30 欧洲渠道会议','17:00 供应链交付沟通'],documents:['10 个文档已完成分类和索引','Shimano 兼容表待人工复核','1 份图片 OCR 字段待校正'],orders:['6 个待确认','9 个生产中','5 个已发运','4 个待收尾'],chats:['出口业务助手 12 次','客户洞察助手 9 次','知识检索助手 8 次','其他智能体 18 次'],news:['法规 4 条','渠道 5 条','智能骑行 3 条','产品更新 3 条']}[key];
  openModal({title:m.label,eyebrow:'今日业务摘要',body:`<div class="panel" style="margin-bottom:14px"><div class="panel-body"><span class="metric-label">当前数量</span><strong class="metric-number" style="display:block;margin-top:6px">${m.value}</strong><p class="secondary-text">${m.detail}</p></div></div><div class="timeline">${items.map((v,i)=>`<div class="timeline-item"><h4>${v}</h4><p>${i===0?'最近更新 10:05':'来自相关业务模块'}</p></div>`).join('')}</div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});
}

function getAgentAllowlist(agentID) {
  return state.agentInternetAllowlists[agentID] || defaultInternetAllowlist;
}

function renderAgentMessage(message, index) {
  const sources = (message.sources || []).map(escapeHTML).join('、');
  return `<div class="message-row ${message.role === 'user' ? 'user' : ''}"><div class="message ${message.role === 'user' ? 'user' : message.error ? 'error' : ''}">${escapeHTML(message.text).replace(/\n/g,'<br>')}${sources ? `<div class="sources">本次来源：${sources}</div>` : ''}${message.error ? `<button class="link-button message-retry" data-action="retry-chat" data-agent="${index}" data-message="${escapeAttr(message.retry || '')}">重新发送</button>` : ''}<time>${escapeHTML(message.time || '')}</time></div></div>`;
}

function renderAgentChatBody(index) {
  const a = agents[index];
  const agentID = agentIDs[index];
  const messages = state.agentChats[agentID] || [];
  const allowlist = getAgentAllowlist(agentID);
  const selected = state.agentSourceSelections[agentID] || ['本地业务数据库','客户私有知识库'];
  const checked = source => selected.includes(source) ? 'checked' : '';
  return `<div class="chat-layout"><aside class="chat-side"><div class="chat-side-head"><h3>本次可用来源</h3><span>按需勾选</span></div><label class="chat-source selectable"><input class="checkbox chat-source-option" type="checkbox" value="本地业务数据库" ${checked('本地业务数据库')}><span><strong>本地业务数据库</strong><small>业务数据接口待接入</small></span></label><label class="chat-source selectable"><input class="checkbox chat-source-option" type="checkbox" value="客户私有知识库" ${checked('客户私有知识库')}><span><strong>客户私有知识库</strong><small>本地向量索引</small></span></label><label class="chat-source selectable"><input class="checkbox chat-source-option" type="checkbox" value="联网检索" ${checked('联网检索')}><span><strong>联网检索</strong><small>${allowlist.length} 个允许域名</small></span></label><button class="button ghost small source-settings-button" data-action="agent-allowlist" data-agent="${index}">${icon('shield-check')}联网白名单</button><p class="chat-source-note">勾选结果和白名单会随本次消息提交给 OpenClaw Agent。工具层域名强制拦截将在联网工具接入时实现。</p></aside><section class="chat-main"><header class="chat-head"><div><strong>${agentEmojis[index]} ${escapeHTML(a[0])}</strong><span>${escapeHTML(agentID)}</span></div><span>${badge(state.modelConfigured?'Active':'Review')}</span></header><div class="chat-quick-prompts">${a[5].map(v=>`<button data-action="chat-quick-prompt" data-agent="${index}" data-prompt="${escapeAttr(v)}">${escapeHTML(v)}</button>`).join('')}</div><div class="chat-messages" id="chatMessages"><div class="message-row"><div class="message"><strong>${escapeHTML(a[0])}</strong><br>已连接 OpenClaw Agent。消息将进入该 Agent 的独立会话。<time>当前会话</time></div></div>${messages.map(m=>renderAgentMessage(m,index)).join('')}</div><div class="chat-status" id="chatStatus" aria-live="polite"></div><div class="chat-input-row"><textarea class="textarea" id="chatInput" maxlength="32768" rows="2" placeholder="向 ${escapeAttr(a[0])} 发送消息，Enter 发送，Shift+Enter 换行"></textarea><button class="icon-button chat-send" data-action="send-chat" data-agent="${index}" title="发送消息" aria-label="发送消息">${icon('send')}</button></div></section></div>`;
}

function wireChatInput(index) {
  const input = document.getElementById('chatInput');
  input?.addEventListener('keydown', event => {
    if (event.key === 'Enter' && !event.shiftKey && !event.isComposing) {
      event.preventDefault();
      void sendAgentMessage(index);
    }
  });
  input?.focus();
}

function showAgentChat(index,prompt='') {
  const a=agents[index];
  openModal({title:a[0],eyebrow:`OpenClaw Agent / ${agentIDs[index]}`,wide:true,body:renderAgentChatBody(index)});
  wireChatInput(index);
  if (prompt) setTimeout(()=>void sendAgentMessage(index,prompt),0);
}

function refreshAgentChat(index) {
  document.getElementById('modalBody').innerHTML = renderAgentChatBody(index);
  applyIcons();
  wireChatInput(index);
  const box = document.getElementById('chatMessages');
  if (box) box.scrollTop = box.scrollHeight;
}

async function sendAgentMessage(index, providedMessage='') {
  const input = document.getElementById('chatInput');
  const message = String(providedMessage || input?.value || '').trim();
  if (!message) {
    toast('请输入消息','消息不能为空。','warning');
    input?.focus();
    return;
  }
  if (message.length > 32768) {
    toast('消息过长','单次消息最多 32,768 个字符。','warning');
    return;
  }
  const agentID = agentIDs[index];
  const sources = [...document.querySelectorAll('.chat-source-option:checked')].map(item=>item.value);
  if (!sources.length) {
    toast('请选择来源','至少选择一个本次可用来源。','warning');
    return;
  }
  const allowlist = sources.includes('联网检索') ? getAgentAllowlist(agentID) : [];
  state.agentSourceSelections[agentID] = sources;
  const history = state.agentChats[agentID] || (state.agentChats[agentID] = []);
  history.push({role:'user',text:message,sources,time:new Date().toLocaleTimeString('zh-CN',{hour:'2-digit',minute:'2-digit'})});
  refreshAgentChat(index);
  const status = document.getElementById('chatStatus');
  const sendButton = document.querySelector('.chat-send');
  if (status) status.innerHTML = `${icon('loader-circle')} ${escapeHTML(agents[index][0])} 正在通过 OpenClaw 处理...`;
  if (sendButton) sendButton.disabled = true;
  applyIcons();
  try {
    const result = await apiFetch('/api/v1/agents/chat', {method:'POST', body:JSON.stringify({agentId:agentID,message,sessionKey:`sta100-${agentID}`,sources,allowlist})});
    history.push({role:'agent',text:result.text || 'OpenClaw 未返回文本内容。',sources,time:new Date().toLocaleTimeString('zh-CN',{hour:'2-digit',minute:'2-digit'})});
  } catch (error) {
    history.push({role:'agent',text:`调用失败：${error.message}`,error:true,retry:message,sources,time:new Date().toLocaleTimeString('zh-CN',{hour:'2-digit',minute:'2-digit'})});
  }
  refreshAgentChat(index);
}

function openAgentAllowlist(index) {
  const agentID = agentIDs[index];
  const domains = getAgentAllowlist(agentID).join('\n');
  openDrawer({title:'联网来源白名单',eyebrow:`${agents[index][0]} / ${agentID}`,body:`<div class="form-field"><label>允许访问的域名</label><textarea class="textarea allowlist-editor" id="agentAllowlist" maxlength="2000" placeholder="每行一个域名，不含协议和路径">${escapeHTML(domains)}</textarea><small>每行一个域名，例如 eur-lex.europa.eu。禁止填写 http://、路径、端口或通配符。</small></div><div class="source-policy-note">${icon('info')}当前白名单会作为来源约束提交给 Agent；联网工具接入后还需在工具调用层执行同一白名单，形成强制访问边界。</div><div class="inline-actions" style="margin-top:16px"><button class="button" data-action="close-drawer">取消</button><button class="button primary" data-action="save-agent-allowlist" data-agent="${index}">${icon('check')}保存白名单</button></div>`});
}

function saveAgentAllowlist(index) {
  const agentID = agentIDs[index];
  const domains = formText('agentAllowlist').split(/\r?\n/).map(v=>v.trim().toLowerCase()).filter(Boolean);
  const validDomain = /^(?=.{1,253}$)(?:[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?\.)+[a-z]{2,63}$/;
  if (!domains.length) {
    toast('白名单不能为空','如不需要联网，请在聊天窗口取消勾选“联网检索”。','warning');
    return;
  }
  const invalid = domains.find(domain=>!validDomain.test(domain));
  if (invalid) {
    toast('域名格式无效',`${invalid} 不是有效域名，请勿填写协议、端口或路径。`,'warning');
    return;
  }
  state.agentInternetAllowlists[agentID] = [...new Set(domains)];
  localStorage.setItem('sta100-agent-allowlists', JSON.stringify(state.agentInternetAllowlists));
  closeDrawer();
  toast('白名单已保存',`${agents[index][0]} 允许 ${state.agentInternetAllowlists[agentID].length} 个联网来源。`);
}

function newCustomerForm(customer) {
  const c=customer||{};
  state.formContext = { type: 'customer', id: c.id || '' };
  openModal({title:customer?'编辑客户':'新建客户',eyebrow:'客户档案',body:`<div class="form-grid"><div class="form-field full"><label>从名片或照片识别</label><button class="upload-zone" data-action="mock-ocr" style="min-height:92px">${icon('scan-line')}<span>选择图片后识别并填充字段</span></button></div><div class="form-section"><h3>基本信息</h3><p>列表字段可配置，详情页保留全部字段。</p></div>${inputField('客户名称',c.name||'',true,false,'text','customerName')}${selectField('客户类型',['Distributor','Importer','Customer','Reseller','Integrator','Supplier','Other'],false,'customerTypeForm',c.type||'Customer')}${inputField('主电话',c.phone||'',true,false,'tel','customerPhone')}${inputField('网站',c.website||'',false,false,'url','customerWebsite')}${inputField('账单国家',c.country||'',true,false,'text','customerCountryForm')}${inputField('账单城市',c.city||'',false,false,'text','customerCity')}${inputField('联系人',c.contact||'',false,false,'text','customerContact')}${inputField('联系邮箱',c.email||'',false,false,'email','customerEmail')}${selectField('评级',['Prospect','Active','Acquired','Market Failed'],false,'customerRating',c.rating||'Prospect')}${selectField('来源',['展会','电话','朋友介绍','拜访','互联网线索','客户转介绍','其它'],false,'customerSource',c.source||'其它')}${inputField('负责人',c.owner||'Donald',false,false,'text','customerOwner')}${inputField('描述',c.description||'',false,true,'text','customerDescription')}</div>`,footer:formFooter(customer?'保存修改':'创建客户','save-customer')});
}

function customerDetail(id, tab='overview') {
  const c=customers.find(x=>x.id===id);
  if (!c) return;
  const tabs = [['overview','概览'],['contacts','联系人'],['quotes','报价单'],['orders','订单'],['documents','单据'],['activity','沟通记录']];
  const tabBar = `<div class="tabs">${tabs.map(([key,label])=>`<button class="${tab===key?'active':''}" data-customer-tab="${key}" data-customer-id="${escapeAttr(c.id)}">${label}</button>`).join('')}</div>`;
  const relatedQuotes = quotes.filter(q=>q.customer===c.name);
  const relatedOrders = orders.filter(o=>o.customer===c.name);
  const relatedDocuments = documents.filter(d=>d.customer===c.name);
  let content = '';
  if (tab === 'quotes') content = `<div class="related-list">${relatedQuotes.map(q=>`<button class="related-record" data-action="quote-detail" data-id="${escapeAttr(q.id)}"><span>${icon('file-text')}<strong>${escapeHTML(q.id)}</strong></span><span>${escapeHTML(q.subject)}</span><span>${escapeHTML(q.value)}</span>${badge(q.status)}</button>`).join('') || `<div class="empty-state"><p>暂无关联报价单</p></div>`}</div>`;
  else if (tab === 'orders') content = `<div class="related-list">${relatedOrders.map(o=>`<button class="related-record" data-action="order-detail" data-id="${escapeAttr(o.id)}"><span>${icon('package')}<strong>${escapeHTML(o.id)}</strong></span><span>${escapeHTML(o.products)}</span><span>${escapeHTML(o.value)}</span>${badge(o.status)}</button>`).join('') || `<div class="empty-state"><p>暂无关联订单</p></div>`}</div>`;
  else if (tab === 'documents') content = `<div class="related-list">${relatedDocuments.map(d=>`<button class="related-record" data-action="document-detail" data-id="${escapeAttr(d.id)}"><span>${icon('file-check-2')}<strong>${escapeHTML(d.id)}</strong></span><span>${escapeHTML(d.order)}</span><span>${escapeHTML(d.template)}</span>${badge(d.status)}</button>`).join('') || `<div class="empty-state"><p>暂无关联单据</p></div>`}</div>`;
  else if (tab === 'contacts') content = `<div class="detail-grid">${[['联系人',c.contact],['电话',c.phone],['邮箱',c.email],['网站',c.website||'未填写'],['国家/城市',`${c.country}${c.city?' / '+c.city:''}`],['来源',c.source||'未填写']].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div>`;
  else if (tab === 'activity') content = `<div class="timeline"><div class="timeline-item"><h4>客户档案已更新</h4><p>${escapeHTML(c.updated)} · ${escapeHTML(c.owner)}</p></div><div class="timeline-item"><h4>关联业务记录 ${relatedQuotes.length + relatedOrders.length + relatedDocuments.length} 条</h4><p>报价单、订单和单据均可从对应标签进入详情。</p></div><div class="timeline-item"><h4>客户来源：${escapeHTML(c.source||'未填写')}</h4><p>来源字段用于后续渠道转化统计。</p></div></div>`;
  else content = `<div class="detail-grid">${[['客户编号',c.id],['客户类型',c.type],['国家',c.country],['城市',c.city||'未填写'],['负责人',c.owner],['主联系人',c.contact],['电话',c.phone],['邮箱',c.email],['来源',c.source||'未填写'],['订单数',String(c.orders)],['累计金额',c.total],['最近更新',c.updated]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 14px">最近业务记录</div><div class="filter-row"><button class="button small" data-customer-tab="quotes" data-customer-id="${escapeAttr(c.id)}">报价单 ${relatedQuotes.length}</button><button class="button small" data-customer-tab="orders" data-customer-id="${escapeAttr(c.id)}">订单 ${relatedOrders.length}</button><button class="button small" data-customer-tab="documents" data-customer-id="${escapeAttr(c.id)}">单据 ${relatedDocuments.length}</button></div>`;
  openDrawer({title:c.name,eyebrow:`客户 / ${c.id}`,body:`${tabBar}<div class="spread" style="margin-bottom:14px"><span>${badge(c.rating)}</span><div class="inline-actions"><button class="button small" data-action="edit-customer" data-id="${escapeAttr(c.id)}">${icon('pencil')}编辑</button><button class="button primary small" data-action="new-quote" data-customer="${escapeAttr(c.name)}">${icon('file-plus-2')}新建报价</button><button class="button danger small" data-action="delete-customer" data-id="${escapeAttr(c.id)}">${icon('trash-2')}删除</button></div></div>${content}`});
}

function newQuoteForm(quote, customerName='') {
  const q=quote||{}; state.formContext={type:'quote',id:q.id||''};
  state.quoteDraftLines = businessLines(q).map(line => ({ ...line, discount: Number(line.discount || 0) }));
  if (!state.quoteDraftLines.length) state.quoteDraftLines = [{ productId: products.find(product => product.status === 'Active')?.id, quantity: 1, unitPrice: moneyNumber(products[0]?.price), discount: 0 }];
  openModal({title:quote?'编辑报价单':'新建报价单',eyebrow:'报价单 / 草稿',wide:true,body:`<div class="form-grid">${inputField('报价主题',q.subject||'欧洲渠道设备报价',true,false,'text','quoteSubject')}${relationField('关联客户','quoteCustomer',customers.filter(c=>!c.archived).map(c=>c.name),customerName||q.customer||customers.find(c=>!c.archived)?.name||'')}${inputField('有效期',q.valid||'2026-09-10',true,false,'date','quoteValid')}${selectField('币种',['EUR','USD','CNY','GBP'],false,'quoteCurrency','EUR')}<div class="form-section"><h3>产品明细</h3><p>产品、库存和默认单价来自产品库；保存后写入报价快照，后续转订单时继续沿用。</p></div><div class="form-field full"><div class="data-wrap"><table class="data-table line-editor-table" style="min-width:820px"><thead><tr><th>产品</th><th>数量</th><th>单价</th><th>折扣 %</th><th>小计</th><th></th></tr></thead><tbody id="quoteLinesBody"></tbody><tfoot><tr><td colspan="4" style="text-align:right"><strong>报价合计</strong></td><td><strong id="quoteDraftTotal">${formatMoney(quoteDraftTotal())}</strong></td><td></td></tr></tfoot></table></div><button type="button" class="button ghost small" data-action="add-quote-line" style="margin-top:8px">${icon('plus')}添加产品</button></div>${inputField('运费','0.00',false,false,'number','quoteFreight')}${inputField('税费','0.00',false,false,'number','quoteTax')}${inputField('条款与条件','30% 预付款，70% 发货前支付。',false,true,'text','quoteTerms')}</div>`,footer:formFooter(quote?'保存修改':'保存草稿','save-quote')});
  renderQuoteDraftLines();
}

function quoteDetail(id) {
  const q=quotes.find(x=>x.id===id);
  if (!q) return;
  openDrawer({title:q.id,eyebrow:'报价单详情',body:`<div class="spread"><span>${badge(q.status)}</span><div class="inline-actions"><button class="button small" data-action="download-quote" data-id="${escapeAttr(q.id)}">${icon('download')}PDF</button><button class="button small" data-action="edit-quote" data-id="${escapeAttr(q.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-quote" data-id="${escapeAttr(q.id)}">${icon('trash-2')}删除</button>${q.status==='Accepted'?`<button class="button primary small" data-action="convert-order" data-id="${escapeAttr(q.id)}">${icon('arrow-right')}转为订单</button>`:''}</div></div><div class="detail-grid" style="margin-top:15px">${[['主题',q.subject],['客户',q.customer],['产品',q.products],['总金额',q.value],['有效期',q.valid],['负责人',q.owner]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 14px">状态记录</div><div class="timeline"><div class="timeline-item"><h4>${q.status==='Accepted'?'客户接受报价':'最近更新报价'}</h4><p>2026-08-10 09:40 · Donald</p></div><div class="timeline-item"><h4>创建报价草稿</h4><p>产品价格来自默认价格表</p></div></div>`});
}

function newOrderForm(order, quoteId='') {
  const o=order||{}; const quote=quotes.find(q=>q.id===(quoteId||o.quote)); state.formContext={type:'order',id:o.id||''};
  state.orderDraftLines = quote ? quoteLinesToOrder(businessLines(quote)) : businessLines(o);
  if (!state.orderDraftLines.length) { const product=products.find(item=>item.status==='Active')||products[0]; state.orderDraftLines=[{productId:product.id,quantity:1,unitPrice:moneyNumber(product.price)}]; }
  openModal({title:order?'编辑订单':'新建订单',eyebrow:'订单 / 草稿',wide:true,body:`<div class="form-grid">${selectField('来源方式',['从已接受报价创建','手动创建'],false,'orderSource',quote?'从已接受报价创建':(order?.quote?'从已接受报价创建':'手动创建'))}${relationField('关联报价','orderQuote',quotes.filter(q=>q.status==='Accepted').map(q=>`${q.id} · ${q.customer}`),quote?`${quote.id} · ${quote.customer}`:'')}${relationField('客户','orderCustomer',customers.filter(c=>!c.archived).map(c=>c.name),o.customer||quote?.customer||customers.find(c=>!c.archived)?.name||'')}${inputField('客户 PO 号',o.po||'',false,false,'text','orderPO')}${inputField('预计交付日期',o.delivery||'2026-09-15',true,false,'date','orderDelivery')}${selectField('贸易条款',['FOB Shenzhen','CIF Hamburg','DAP Customer'],false,'orderTerms',o.terms||'FOB Shenzhen')}<div class="form-section"><h3>订单明细</h3><p>产品选择、数量、库存和成交单价均来自产品库；保存为订单明细快照，生成单据时继续读取本快照。</p></div><div class="form-field full"><div class="data-wrap"><table class="data-table line-editor-table" style="min-width:820px"><thead><tr><th>产品</th><th>数量</th><th>成交单价</th><th>小计</th><th></th></tr></thead><tbody id="orderLinesBody"></tbody><tfoot><tr><td colspan="3" style="text-align:right"><strong>订单合计</strong></td><td><strong id="orderDraftTotal">${formatMoney(orderDraftTotal())}</strong></td><td></td></tr></tfoot></table></div><button type="button" class="button ghost small" data-action="add-order-line" style="margin-top:8px">${icon('plus')}添加产品</button></div></div>`,footer:formFooter(order?'保存修改':'创建订单','save-order')});
  renderOrderDraftLines();
}

function orderDetail(id) {
  const o=orders.find(x=>x.id===id);
  if (!o) return;
  openDrawer({title:o.id,eyebrow:'订单详情',body:`<div class="spread"><span>${badge(o.status)}</span><div class="inline-actions"><button class="button small" data-action="edit-order" data-id="${escapeAttr(o.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-order" data-id="${escapeAttr(o.id)}">${icon('trash-2')}删除</button><button class="button primary small" data-action="generate-docs" data-id="${escapeAttr(o.id)}">${icon('files')}生成单据</button></div></div><div class="detail-grid" style="margin-top:15px">${[['客户',o.customer],['来源报价',o.quote||'手动创建'],['产品明细数',String(businessLines(o).length)],['订单金额',o.value],['预计交付',o.delivery],['完成进度',`${o.progress}%`]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 14px">订单产品明细</div>${detailLinesTable(businessLines(o),true)}<div class="divider-title" style="margin:20px 0 14px">订单进度</div><div class="timeline"><div class="timeline-item"><h4>${o.status==='Shipped'?'已发运':'当前阶段：'+o.status}</h4><p>状态由本系统业务操作维护</p></div><div class="timeline-item"><h4>订单确认</h4><p>客户、产品与交付信息已保存为订单快照</p></div><div class="timeline-item"><h4>${o.quote?'报价转订单':'手动创建订单'}</h4><p>${o.quote?`${escapeHTML(o.quote)} · 保留完整来源`:'直接从产品库选择产品'}</p></div></div>`});
}

function generateDocs(orderId='', document) {
  const d=document||{}; state.formContext={type:'document',id:d.id||''};
  openModal({title:document?'编辑单据':'生成外贸单据',eyebrow:orderId==='all'?'一键生成全套':`订单 / ${orderId||d.order||''}`,body:`<div class="form-grid">${relationField('关联订单','documentOrder',orders.map(o=>`${o.id} · ${o.customer}`),d.order?`${d.order} · ${d.customer}`:orderId&&orderId!=='all'?`${orderId} · ${orders.find(o=>o.id===orderId)?.customer||''}`:'')}${selectField('单据类型',['PI','CI','PL','报关单'],false,'documentTypeForm',d.type||'PI')}${selectField('模板版本',['各类型默认模板','STRATRONIX 标准模板组 v3'],false,'documentTemplate',d.template||'各类型默认模板')}${selectField('输出语言',['英文','中文 / 英文双语'],false,'documentLanguage','英文')}<div class="form-section"><h3>生成规则</h3><p>系统读取客户、订单和产品快照，生成后先进入待复核状态；缺失字段会逐项提示。</p></div><div class="model-warning full"><span>${icon('info')} 报关单的申报要素、监管条件和最终格式需在正式数据提供后校准。</span></div></div>`,footer:formFooter(document?'保存修改':'生成并预览','save-document')});
}

function saveCustomer() {
  const name=formText('customerName');
  if (!name || !formText('customerCountryForm')) { toast('保存失败','客户名称和账单国家为必填项。','warning'); return; }
  const existing=customers.find(c=>c.id===state.formContext?.id);
  const record=existing||{id:nextRecordId('ACC',customers),orders:0,total:'EUR 0',updated:nowText()};
  Object.assign(record,{name,type:formText('customerTypeForm'),phone:formText('customerPhone'),website:formText('customerWebsite'),country:formText('customerCountryForm'),city:formText('customerCity'),contact:formText('customerContact'),email:formText('customerEmail'),rating:formText('customerRating'),source:formText('customerSource'),owner:formText('customerOwner'),description:formText('customerDescription'),updated:nowText(),archived:false});
  if (!existing) customers.unshift(record);
  closeModal(); renderPage(); toast(existing?'客户已更新':'客户已创建',`${record.name} 已保存到客户档案。`);
}

function saveQuote() {
  const subject=formText('quoteSubject'); const customer=formText('quoteCustomer');
  if (!subject || !customer || !customers.some(c=>c.name===customer&&!c.archived)) { toast('保存失败','请填写主题，并从客户搜索结果中选择有效客户。','warning'); return; }
  if (!state.quoteDraftLines.length || state.quoteDraftLines.some(line => Number(line.quantity) < 1 || Number(line.unitPrice) < 0)) { toast('保存失败','报价至少需要一条有效产品明细，数量和单价不能为负数。','warning'); return; }
  const existing=quotes.find(q=>q.id===state.formContext?.id);
  const currency=formText('quoteCurrency')||'EUR'; const total=quoteDraftTotal();
  const record=existing||{id:nextRecordId('QUO-2026',quotes),status:'Draft',owner:'Donald'};
  Object.assign(record,{subject,customer,valid:formText('quoteValid'),currency,lines:state.quoteDraftLines.map(line=>({...line,quantity:Number(line.quantity),unitPrice:Number(line.unitPrice),discount:Number(line.discount||0),amount:lineSubtotal(line,true)})),products:lineSummary(state.quoteDraftLines),value:formatMoney(total,currency),updated:nowText()});
  if (!existing) quotes.unshift(record);
  closeModal(); renderPage(); toast(existing?'报价单已更新':'报价草稿已创建',`${record.id} 已保存。`);
}

function saveOrder() {
  const customer=formText('orderCustomer'); const quoteValue=formText('orderQuote');
  if (!customer || !customers.some(c=>c.name===customer&&!c.archived)) { toast('保存失败','请从客户搜索结果中选择有效客户。','warning'); return; }
  if (!state.orderDraftLines.length || state.orderDraftLines.some(line => Number(line.quantity) < 1 || Number(line.unitPrice) < 0)) { toast('保存失败','订单至少需要一条有效产品明细，数量和成交单价不能为负数。','warning'); return; }
  const insufficient=state.orderDraftLines.find(line=>Number(line.quantity)>Number(lineProduct(line)?.stock||0));
  if (insufficient) { toast('保存失败',`${lineProduct(insufficient)?.name||'产品'} 数量超过当前库存 ${lineProduct(insufficient)?.stock||0}。`,'warning'); return; }
  const existing=orders.find(o=>o.id===state.formContext?.id);
  const record=existing||{id:nextRecordId('SO-2026',orders),status:'Confirmed',progress:0};
  Object.assign(record,{customer,quote:quoteValue.split(' · ')[0]||quoteValue,po:formText('orderPO'),delivery:formText('orderDelivery'),terms:formText('orderTerms'),lines:state.orderDraftLines.map(line=>({...line,quantity:Number(line.quantity),unitPrice:Number(line.unitPrice),amount:lineSubtotal(line)})),products:lineSummary(state.orderDraftLines),value:formatMoney(orderDraftTotal(),'EUR'),updated:nowText()});
  if (!existing) orders.unshift(record);
  refreshCustomerAggregates(); closeModal(); renderPage(); toast(existing?'订单已更新':'订单已创建',`${record.id} 已保存，已关联 ${record.lines.length} 个产品明细。`);
}

function refreshCustomerAggregates() {
  customers.forEach(customer => {
    const related = orders.filter(order => order.customer === customer.name && order.status !== 'Cancelled');
    if (!related.length) return;
    customer.orders = related.length;
    customer.total = formatMoney(related.reduce((sum, order) => sum + moneyNumber(order.value), 0), 'EUR');
  });
}

function saveDocument() {
  const orderValue=formText('documentOrder'); const orderId=orderValue.split(' · ')[0]; const order=orders.find(o=>o.id===orderId);
  if (!order) { toast('保存失败','请从订单搜索结果中选择有效订单。','warning'); return; }
  const existing=documents.find(d=>d.id===state.formContext?.id);
  const type=formText('documentTypeForm');
  const record=existing||{id:`${type}-${new Date().toISOString().slice(0,10).replaceAll('-','')}-${String(documents.length+1).padStart(3,'0')}`,status:'Review'};
  Object.assign(record,{type,customer:order.customer,order:order.id,template:formText('documentTemplate'),language:formText('documentLanguage'),lines:businessLines(order).map(line=>({...line,productName:lineProduct(line)?.name||line.productId,amount:line.amount??lineSubtotal(line)})),value:order.value,updated:nowText()});
  if (!existing) documents.unshift(record);
  closeModal(); renderPage(); toast(existing?'单据已更新':'单据已生成',`${record.id} 已保存为待复核状态。`);
}

function deleteCustomer(id) {
  const customer=customers.find(c=>c.id===id); if (!customer) return;
  if (!window.confirm(`确定归档客户 ${customer.name} 吗？关联报价、订单和单据会保留。`)) return;
  customer.archived=true; closeDrawer(); renderPage(); toast('客户已归档','关联业务记录仍然保留，可用于历史查询。');
}
function deleteQuote(id) { if (!window.confirm('确定删除这份报价单吗？')) return; const index=quotes.findIndex(q=>q.id===id); if(index>=0) quotes.splice(index,1); closeDrawer(); renderPage(); toast('报价单已删除','前端示例数据已移除。'); }
function deleteOrder(id) { if (!window.confirm('确定删除这份订单吗？')) return; const index=orders.findIndex(o=>o.id===id); if(index>=0) orders.splice(index,1); refreshCustomerAggregates(); closeDrawer(); renderPage(); toast('订单已删除','前端示例数据已移除。'); }
function deleteDocument(id) { if (!window.confirm('确定删除这份单据吗？')) return; const index=documents.findIndex(d=>d.id===id); if(index>=0) documents.splice(index,1); closeDrawer(); renderPage(); toast('单据已删除','前端示例数据已移除。'); }
function downloadQuote(id) { const q=quotes.find(item=>item.id===id); if(!q)return; downloadText(`${q.id}.txt`,`报价单 ${q.id}\n客户：${q.customer}\n主题：${q.subject}\n金额：${q.value}\n有效期：${q.valid}`); toast('报价单已下载',`${q.id} 已导出为文本预览。`); }
function downloadDocument(id) { const d=documents.find(item=>item.id===id); if(!d)return; downloadText(`${d.id}.txt`,`单据 ${d.id}\n类型：${d.type}\n客户：${d.customer}\n订单：${d.order}\n模板：${d.template}`); toast('单据已下载',`${d.id} 已导出为文本预览。`); }

function templateCenter(kind) {
  const name={quote:'报价单',order:'订单',document:'单据'}[kind]||'业务';
  state.templateKind=kind;
  const selected=state.templateUploads.filter(item=>item.kind===kind).slice(-1)[0];
  const selectedStatus=selected?`<section class="template-upload-status"><span class="setting-icon">${icon(selected.mode==='image'?'scan-line':'file-check-2')}</span><div><strong>${escapeHTML(selected.name)}</strong><span>${escapeHTML(selected.size)} · ${selected.mode==='image'?'图片识别模板':'文件模板'} · 本地校验通过</span><small>当前状态：等待后端上传接口，尚未上传、识别或发布。</small></div>${badge('Review')}</section>`:'';
  const interfaceRows=[
    ['POST /api/v1/templates/upload','接收模板文件，校验类型、大小、摘要和版本','未实现'],
    ['POST /api/v1/templates/image-recognition','OCR 与版式识别，输出字段坐标和置信度','未实现'],
    ['GET/PATCH /api/v1/templates/{id}','读取模板、维护字段映射和人工校正结果','未实现'],
    ['POST /api/v1/templates/{id}/publish','预览确认后发布版本并设置默认模板','未实现'],
  ];
  openModal({title:`${name}模板管理`,eyebrow:'模板中心',wide:true,body:`<input id="templateImageInput" type="file" accept="image/jpeg,image/png,image/webp" hidden><input id="templateFileInput" type="file" accept=".docx,.xlsx,.pdf,.html,.htm" hidden><div class="toolbar" style="margin-bottom:14px"><span class="badge green">2 个示例模板</span><span class="spacer"></span><button class="button" data-action="upload-template-image">${icon('image-up')}上传图片生成模板</button><button class="button primary" data-action="upload-template-file">${icon('file-up')}上传模板</button></div>${selectedStatus}<div class="data-wrap"><table class="data-table" style="min-width:700px"><thead><tr><th>模板名称</th><th>来源</th><th>版本</th><th>状态</th><th>更新时间</th><th>操作</th></tr></thead><tbody><tr><td>STRATRONIX 标准${name}</td><td>固定模板</td><td>v3</td><td>${badge('Active')}</td><td>2026-08-09</td><td><button class="button small" data-action="template-default">设为默认</button></td></tr><tr><td>欧洲渠道${name}</td><td>图片识别后校正</td><td>v2</td><td>${badge('Active')}</td><td>2026-08-08</td><td><button class="button small" data-action="template-edit">编辑</button></td></tr></tbody></table></div><div class="divider-title" style="margin:18px 0 10px">接口实现状态</div><div class="data-wrap"><table class="data-table template-api-table"><thead><tr><th>接口</th><th>功能</th><th>当前状态</th></tr></thead><tbody>${interfaceRows.map(([api,purpose,status])=>`<tr><td><code>${escapeHTML(api)}</code></td><td>${escapeHTML(purpose)}</td><td><span class="badge amber">${status}</span></td></tr>`).join('')}</tbody></table></div><div class="model-warning template-blocker" style="margin-top:14px"><span>${icon('triangle-alert')} 后端暂未实现：客户尚未提供正式模板样例、占位符和字段映射规范、OCR 引擎及精度要求、纸张/语言/签章规则。当前仅完成文件本地选择、格式校验和流程状态展示，不会伪装为已上传成功。</span></div><ol class="template-flow"><li>图片模板：选择图片 → 上传 → OCR/版式识别 → 字段映射 → 人工校正 → 预览发布</li><li>文件模板：选择 DOCX/XLSX/PDF/HTML → 上传校验 → 占位符识别 → 字段映射 → 预览发布</li></ol>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});
  document.getElementById('templateImageInput')?.addEventListener('change',event=>handleTemplateFile(event.target.files?.[0],'image',kind));
  document.getElementById('templateFileInput')?.addEventListener('change',event=>handleTemplateFile(event.target.files?.[0],'file',kind));
}

function handleTemplateFile(file,mode,kind) {
  if(!file)return;
  const allowed=mode==='image'?/\.(jpe?g|png|webp)$/i:/\.(docx|xlsx|pdf|html?)$/i;
  if(!allowed.test(file.name)){toast('文件格式不支持',mode==='image'?'请选择 JPG、PNG 或 WebP 图片。':'请选择 DOCX、XLSX、PDF 或 HTML 模板。','warning');return;}
  if(file.size>20*1024*1024){toast('文件过大','模板文件不能超过 20 MB。','warning');return;}
  const size=file.size>=1024*1024?`${(file.size/1024/1024).toFixed(1)} MB`:`${Math.max(1,Math.round(file.size/1024))} KB`;
  state.templateUploads.push({kind,mode,name:file.name,size,status:'waiting-api'});
  templateCenter(kind);
  toast('本地文件校验通过','文件尚未上传，正在等待模板后端接口接入。','warning');
}

function templateAction(action) {
  if (action === 'default') {
    toast('模板已设为默认', '后续生成单据将优先使用该模板。');
    return;
  }
  openModal({ title: '编辑模板', eyebrow: '模板中心 / 前端编辑', body: `<div class="form-grid">${inputField('模板名称', '欧洲渠道业务模板', true, true)}${selectField('适用单据', ['报价单','订单','PI','CI','PL','报关单'])}${inputField('版本', 'v2')}${selectField('状态', ['启用','待校正','停用'])}<div class="form-field full"><label>字段映射说明</label><textarea class="textarea">客户名称、地址、产品明细、币种、金额、交付和贸易条款</textarea></div></div>`, footer: formFooter('保存模板') });
}

function documentDetail(id) {
  const d=documents.find(x=>x.id===id);
  if(!d)return;
  const sourceOrder=orders.find(order=>order.id===d.order); const lines=d.lines?.length?d.lines:businessLines(sourceOrder||{});
  const previewRows=lines.map(line=>`<tr><td style="padding:9px 7px">${escapeHTML(line.productName||lineProduct(line)?.name||line.productId)}</td><td style="text-align:center">${Number(line.quantity)}</td><td style="text-align:center">${formatMoney(line.unitPrice)}</td><td style="text-align:right">${formatMoney(line.amount??lineSubtotal(line))}</td></tr>`).join('');
  openDrawer({title:d.id,eyebrow:`${d.type} / 单据预览`,body:`<div class="spread"><span>${badge(d.status)}</span><div class="inline-actions"><button class="button small" data-action="edit-document" data-id="${escapeAttr(d.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-document" data-id="${escapeAttr(d.id)}">${icon('trash-2')}删除</button><button class="button primary small" data-action="download-document" data-id="${escapeAttr(d.id)}">${icon('file-down')}下载 PDF</button></div></div><div class="panel document-preview-paper" style="margin-top:15px;background:#edf1f5;color:#1b2634;min-height:520px;padding:28px"><div style="border-bottom:2px solid #29394b;padding-bottom:14px"><strong style="font-size:18px">STA-100</strong><span style="float:right;font-size:18px">${escapeHTML(d.type)}</span></div><div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;margin-top:20px;font-size:10px"><div><strong>Bill To</strong><p>${escapeHTML(d.customer)}<br>Customer address snapshot</p></div><div><strong>Document</strong><p>${escapeHTML(d.id)}<br>Order ${escapeHTML(d.order)}</p></div></div><table style="width:100%;border-collapse:collapse;margin-top:22px;font-size:9px"><tr style="border-bottom:1px solid #8793a1"><th style="text-align:left;padding:7px">Item</th><th>Qty</th><th>Rate</th><th>Amount</th></tr>${previewRows}</table><div style="margin-top:18px;text-align:right"><strong>Total ${escapeHTML(d.value||sourceOrder?.value||formatMoney(lines.reduce((sum,line)=>sum+lineSubtotal(line),0)))}</strong></div></div>`});
}

function newProductForm(product) {
  const p=product||{};
  state.formContext={type:'product',id:p.id||''};
  openModal({title:product?'编辑产品':'新建产品',eyebrow:'产品主数据',body:`<div class="form-grid">${inputField('产品名称',p.name||'',true,false,'text','productName')}${inputField('产品编码',p.id||'',true,false,'text','productID')}${selectField('产品类别',['智能设备','智能骑行','整车方案','配件','服务'],false,'productCategory',p.category||'智能设备')}${inputField('制造商',p.manufacturer||'STRATRONIX',false,false,'text','productManufacturer')}${inputField('HS CODE',p.hs||'',true,false,'text','productHS')}${inputField('库存量',p.stock||'0',false,false,'number','productStock')}${inputField('默认单价',p.price||'',true,false,'text','productPrice')}${inputField('产品描述',p.desc||'',false,true,'text','productDescription')}${inputField('标签','欧洲 / 智能设备',false,true,'text','productTags')}</div>`,footer:formFooter(product?'保存修改':'创建产品','save-product')});
}

function saveProduct() {
  const name=formText('productName'); const id=formText('productID'); if(!name||!id){toast('保存失败','产品名称和编码为必填项。','warning');return;}
  const existing=products.find(p=>p.id===state.formContext?.id); const record=existing||{status:'Active'};
  Object.assign(record,{id,name,category:formText('productCategory'),manufacturer:formText('productManufacturer'),hs:formText('productHS'),stock:formNumber('productStock'),price:formText('productPrice'),desc:formText('productDescription'),tags:formText('productTags')});
  if(!existing)products.unshift(record); closeModal(); renderPage(); toast(existing?'产品已更新':'产品已创建',`${record.name} 已保存。`);
}
function deleteProduct(id) { if(!window.confirm('确定删除该产品吗？'))return; const index=products.findIndex(p=>p.id===id); if(index>=0)products.splice(index,1); closeDrawer(); renderPage(); toast('产品已删除','前端示例数据已移除。'); }

function productDetail(id) {
  const p=products.find(x=>x.id===id);
  if(!p)return;
  openDrawer({title:p.name,eyebrow:`产品 / ${p.id}`,body:`<div class="spread"><span>${badge(p.status)}</span><div class="inline-actions"><button class="button small" data-action="edit-product" data-id="${escapeAttr(p.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-product" data-id="${escapeAttr(p.id)}">${icon('trash-2')}删除</button></div></div><div class="product-visual panel" style="margin-top:14px;aspect-ratio:16/6">${icon('cpu')}</div><div class="detail-grid" style="margin-top:14px">${[['产品编码',p.id],['产品类别',p.category],['HS CODE',p.hs],['销售价',p.price],['当前库存',String(p.stock)],['状态',p.status]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">产品描述</div><p class="secondary-text" style="line-height:1.7">${escapeHTML(p.desc)}</p><div class="divider-title" style="margin:20px 0 12px">业务引用</div><div class="filter-row"><span class="badge blue">报价单 8</span><span class="badge amber">进行中订单 3</span><span class="badge green">历史单据 12</span></div>`});
}

function uploadFileModal() {
  openModal({title:'上传私有数据',eyebrow:'数据库 / 文件处理',body:`<div class="upload-zone" id="uploadZone"><div><span class="upload-icon">${icon('cloud-upload')}</span><h3>选择或拖入文件</h3><p>支持 PDF、DOCX、XLSX、CSV、TXT、MD、JPG、PNG；视频格式待调研确认。</p><button class="button primary" data-action="choose-file">选择文件</button></div></div><div class="form-grid" style="margin-top:14px">${selectField('数据区',['客户私有数据'])}${selectField('主分类',['自动识别','合同','报价单','产品手册','法规','产品资料','会议记录','客户资料','其它'])}${inputField('附加标签','',false,true)}<div class="form-field full"><small>上传后依次执行文件校验、文本/OCR 解析、AI 分类、标签建议、重复检查和索引；失败记录会保留并给出原因。</small></div></div>`,footer:formFooter('开始处理')});
  const zone=document.getElementById('uploadZone');
  ['dragenter','dragover'].forEach(ev=>zone.addEventListener(ev,e=>{e.preventDefault();zone.classList.add('dragover')}));
  ['dragleave','drop'].forEach(ev=>zone.addEventListener(ev,e=>{e.preventDefault();zone.classList.remove('dragover')}));
}

function filePreview(name) {
  const f=files.find(x=>x.name===name);
  if (!f) return;
  openDrawer({title:f.name,eyebrow:`${f.category} / 在线预览`,body:`<div class="spread"><span>${badge(f.status)}</span><div class="inline-actions"><button class="button small" data-action="file-download" data-name="${escapeAttr(f.name)}">${icon('download')}下载</button><button class="button small" data-action="file-edit" data-name="${escapeAttr(f.name)}">${icon('pencil')}编辑信息</button></div></div><div class="detail-grid" style="margin-top:15px">${[['主分类',f.category],['来源',f.source],['大小',f.size],['更新时间',f.updated],['标签',f.tags.join(' / ')],['索引状态',f.status]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="empty-state panel" style="margin-top:16px">${icon('file-search')}<div><h3>文件预览区</h3><p>正式接入后按文件格式调用本地预览器；不支持的格式提供下载和 AI 摘要。</p><button class="button small" data-action="file-summary" data-name="${escapeAttr(f.name)}">${icon('sparkles')}查看 AI 摘要</button></div></div>`});
}

function fileDownload(name) {
  const f = files.find(item => item.name === name);
  if (!f) return;
  downloadText(f.name + '.summary.txt', `文件：${f.name}\n分类：${f.category}\n来源：${f.source}\n索引状态：${f.status}`);
  toast('文件下载已开始', f.name);
}

function fileEdit(name) {
  const f = files.find(item => item.name === name);
  if (!f) return;
  openModal({ title: '编辑文件信息', eyebrow: '数据库 / 文件元数据', body: `<div class="form-grid">${inputField('文件名', f.name, true, true)}${selectField('主分类', ['产品手册','合同','法规','产品资料','会议记录','客户资料','其它'], false, '', f.category)}${inputField('标签', f.tags.join(' / '), false, true)}${inputField('来源', f.source, false, true)}</div>`, footer: formFooter('保存信息') });
}

function fileSummary(name) {
  const f = files.find(item => item.name === name);
  if (f) toast('AI 摘要已生成', `${f.name}：已提取文件分类、来源和索引状态；正式版将由知识检索 Agent 返回摘要。`);
}

function newsDetail(title) {
  const n=news.find(x=>x.title===title)||news[0];
  openDrawer({title:n.title,eyebrow:`${n.category} / ${n.source}`,body:`<div class="spread"><span class="badge green">相关度 ${n.relevance}</span><span class="secondary-text">${n.time}</span></div><p style="margin-top:20px;line-height:1.8;color:var(--text)">${n.summary}</p><div class="divider-title" style="margin:20px 0 12px">智能体摘要</div><p class="secondary-text" style="line-height:1.8">该信息与当前关注的欧洲渠道、智能骑行产品和合规主题相关。建议结合客户档案与产品库，检查受影响客户和产品后再形成行动项。</p><div class="divider-title" style="margin:20px 0 12px">来源信息</div><div class="detail-grid">${[['来源',n.source],['获取时间',n.time],['信息类别',n.category],['数据区域','互联网推荐数据']].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${v}</strong></div>`).join('')}</div><div class="inline-actions" style="margin-top:18px"><button class="button" data-action="news-source-link" data-title="${escapeAttr(n.title)}">${icon('external-link')}查看原文</button><button class="button primary" data-action="news-todo" data-title="${escapeAttr(n.title)}">${icon('list-plus')}生成待办</button></div>`});
}

function configureModel() {
  const data = state.openClawModels;
  if (!data) {
    toast('模型配置尚未就绪', '请等待 OpenClaw 模型信息读取完成。', 'warning');
    return;
  }
  const models = data.models || [];
  const current = data.resolvedDefault || data.defaultModel || '';
  const providers = [...new Set(['deepseek', ...(data.providers || []).map(item => item.provider)])];
  const modelOptions = models.map(model => `<option value="${escapeAttr(model.key)}" ${model.key===current?'selected':''} ${model.missing||!model.available?'disabled':''}>${escapeHTML(model.name)} · ${escapeHTML(model.key)}${model.missing||!model.available?'（不可用）':''}</option>`).join('');
  const configuredProviders=(data.providers||[]).filter(provider=>provider.configured).map(provider=>provider.provider);
  const storedKey=state.showApiKey ? '已配置（OpenClaw 不提供明文读取）' : '••••••••••••';
  openModal({title:'配置 OpenClaw 模型',eyebrow:'设置 / 真实配置',wide:true,body:`<div class="form-grid"><div class="form-field full"><label>默认模型 <span class="required">*</span></label><select class="select" id="modelID">${modelOptions}</select><small>下拉列表直接来自 OpenClaw <code>models list</code>；当前共有 ${models.length} 个模型，仅可选择 OpenClaw 标记为可用的模型。</small></div><div class="form-field"><label>凭据提供商</label><select class="select" id="modelProvider">${providers.map(provider=>`<option value="${escapeAttr(provider)}">${escapeHTML(provider)}${configuredProviders.includes(provider)?'（已配置）':''}</option>`).join('')}</select></div><div class="form-field"><label>已保存 API Key</label><div class="credential-field"><input class="input" readonly value="${escapeAttr(storedKey)}"><button class="icon-button" type="button" data-action="toggle-api-key" title="${state.showApiKey?'隐藏':'查看凭据状态'}">${icon(state.showApiKey?'eye-off':'eye')}</button></div></div><div class="form-field full"><label>更新 API Key</label><div class="credential-field"><input class="input" id="modelAPIKey" type="password" value="" autocomplete="new-password" placeholder="留空则保留现有凭据"><button class="icon-button" type="button" data-action="toggle-api-key-input" title="显示或隐藏本次输入">${icon('eye')}</button></div><small>可以查看本次输入的新 Key；已保存 Key 只显示配置状态。后端和 OpenClaw 不向页面返回凭据明文，避免本机页面、日志或接口泄露。</small></div><div class="form-field full"><label>OpenClaw 支持的模型</label><div class="data-wrap"><table class="data-table model-list-table"><thead><tr><th>模型</th><th>模型 ID</th><th>上下文</th><th>状态</th></tr></thead><tbody>${models.map(model=>`<tr><td>${escapeHTML(model.name)}</td><td>${escapeHTML(model.key)}</td><td>${Number(model.contextWindow||0).toLocaleString()}</td><td>${model.available&&!model.missing?badge('Active'):badge('Review')}</td></tr>`).join('')}</tbody></table></div></div></div>`,footer:`<button class="button ghost" data-action="close-modal">取消</button><button class="button primary" data-action="save-model-config">${icon('check')}保存到 OpenClaw</button>`});
}

async function saveModelConfiguration(button) {
  const model = document.getElementById('modelID')?.value || '';
  const provider = document.getElementById('modelProvider')?.value || '';
  const apiKey = document.getElementById('modelAPIKey')?.value.trim() || '';
  if (!model) {
    toast('无法保存', '请选择可用模型。', 'warning');
    return;
  }
  button.disabled = true;
  button.innerHTML = `${icon('loader-circle')}保存中`;
  applyIcons();
  try {
    if (apiKey) await apiFetch('/api/v1/openclaw/models/auth', { method: 'POST', body: JSON.stringify({ provider, apiKey }) });
    if (model !== state.openClawModels?.resolvedDefault) await apiFetch('/api/v1/openclaw/models/default', { method: 'PUT', body: JSON.stringify({ model }) });
    await loadOpenClawModels(true);
    closeModal();
    toast('OpenClaw 模型配置已更新', `${model} 已设为默认模型。`);
  } catch (error) {
    button.disabled = false;
    button.innerHTML = `${icon('check')}保存到 OpenClaw`;
    applyIcons();
    toast('模型配置更新失败', error.message, 'warning');
  }
}

function renderAgentManagerBody() {
  const rows = state.openClawAgents?.filter(agent => !agent.isDefault) || [];
  if (state.openClawAgentsLoading && !rows.length) return `<div class="empty-state compact-empty">${icon('loader-circle')}<div><h3>正在读取 OpenClaw Agent</h3><p>请稍候。</p></div></div>`;
  return `<div class="manager-summary"><span>${icon('check-circle-2')} 已注册 <strong>${rows.length}</strong> / 24 个 STA-100 Agent</span><span>数据来源：OpenClaw agents list</span></div><div class="data-wrap"><table class="data-table agent-manager-table"><thead><tr><th>业务智能体</th><th>Agent ID</th><th>模型</th><th>状态</th></tr></thead><tbody>${rows.map(agent=>`<tr><td><span class="agent-manager-name"><span class="agent-emoji" aria-hidden="true">${escapeHTML(agent.identityEmoji || '🤖')}</span><strong>${escapeHTML(agent.identityName || agent.name || agent.id)}</strong></span></td><td>${escapeHTML(agent.id)}</td><td>${escapeHTML(agent.model)}</td><td>${badge('Active')}</td></tr>`).join('') || `<tr><td colspan="4"><div class="empty-state"><p>未读取到 STA-100 Agent。</p></div></td></tr>`}</tbody></table></div>`;
}

async function openAgentManager() {
  openModal({title:'智能体管理',eyebrow:'OpenClaw / 真实 Agent',wide:true,body:renderAgentManagerBody(),footer:`<button class="button" data-action="close-modal">关闭</button><button class="button primary" data-action="sync-openclaw-agents">${icon('refresh-cw')}按 STA-100 清单同步</button>`});
  if (!state.openClawAgents) {
    await loadOpenClawAgents();
    if (!document.getElementById('modalBackdrop').hidden) {
      document.getElementById('modalBody').innerHTML = renderAgentManagerBody();
      applyIcons();
    }
  }
}

async function syncOpenClawAgents(button) {
  button.disabled = true;
  button.innerHTML = `${icon('loader-circle')}同步中`;
  applyIcons();
  try {
    const data = await apiFetch('/api/v1/openclaw/agents/sync', { method: 'POST' });
    state.openClawAgents = data.agents || [];
    document.getElementById('modalBody').innerHTML = renderAgentManagerBody();
    button.disabled = false;
    button.innerHTML = `${icon('refresh-cw')}按 STA-100 清单同步`;
    applyIcons();
    toast('Agent 同步完成', `OpenClaw 当前共 ${data.count} 个 Agent，其中 24 个由 STA-100 清单管理。`);
  } catch (error) {
    button.disabled = false;
    button.innerHTML = `${icon('refresh-cw')}重新同步`;
    applyIcons();
    toast('Agent 同步失败', error.message, 'warning');
  }
}

function schedulerForm(name='新增定时任务') {
  openModal({title:name,eyebrow:'设置 / 定时任务',body:`<div class="form-grid">${inputField('任务名称',name==='新增定时任务'?'':name,true,true)}${selectField('任务类型',['每日推荐更新','智能体周报','行业新闻更新','数据索引维护','自定义任务'],true)}${selectField('执行频率',['每 60 分钟','每天','每周','自定义 Cron'])}${inputField('执行时间','08:00')}<div class="form-field full"><label>任务内容</label><textarea class="textarea">根据用户关注条件更新推荐信息，去重后保存最近 20 条。</textarea></div><div class="form-field full"><label style="display:flex;align-items:center;gap:8px"><input type="checkbox" checked> 启用任务</label></div></div>`,footer:formFooter('保存任务')});
}

function toast(title, message, type='success') {
  const node=document.createElement('div');
  node.className=`toast ${type}`;
  node.innerHTML=`${icon(type==='warning'?'triangle-alert':'circle-check')}<div><strong>${title}</strong><span>${message}</span></div>`;
  document.getElementById('toastRegion').appendChild(node);
  applyIcons();
  setTimeout(()=>node.remove(),3600);
}

function openCommand(query='') {
  document.getElementById('commandBackdrop').hidden=false;
  const input=document.getElementById('commandInput');
  input.value=query;
  renderCommandResults(query);
  setTimeout(()=>input.focus(),0);
}
function closeCommand(){document.getElementById('commandBackdrop').hidden=true;}
function renderCommandResults(query='') {
  const entries=[
    ...customers.map(c=>({type:'客户',label:c.name,sub:`${c.id} · ${c.country}`,icon:'building-2',page:'customers'})),
    ...orders.map(o=>({type:'订单',label:o.id,sub:`${o.customer} · ${o.value}`,icon:'package-check',page:'orders'})),
    ...products.map(p=>({type:'产品',label:p.name,sub:`${p.id} · ${p.hs}`,icon:'boxes',page:'products'})),
    ...files.map(f=>({type:'文件',label:f.name,sub:`${f.category} · ${f.source}`,icon:'file',page:'database'})),
  ].filter(e=>!query||`${e.type} ${e.label} ${e.sub}`.toLowerCase().includes(query.toLowerCase())).slice(0,8);
  document.getElementById('commandResults').innerHTML=`<div class="command-results">${entries.length?entries.map(e=>`<button class="command-result" data-command-page="${e.page}"><span class="result-icon">${icon(e.icon)}</span><div><strong>${e.label}</strong><small>${e.type} · ${e.sub}</small></div>${icon('arrow-right')}</button>`).join(''):`<div class="command-empty">未找到匹配结果</div>`}</div>`;
  applyIcons();
}

document.addEventListener('input', e => {
  const quoteLine=e.target.closest('[data-quote-line-field]');
  if(quoteLine){
    const line=state.quoteDraftLines[Number(quoteLine.dataset.index)]; if(!line)return;
    line[quoteLine.dataset.quoteLineField]=quoteLine.dataset.quoteLineField==='productId'?quoteLine.value:Number(quoteLine.value);
    const subtotal=document.querySelector(`[data-quote-line-subtotal="${quoteLine.dataset.index}"]`); if(subtotal)subtotal.textContent=formatMoney(lineSubtotal(line,true));
    const total=document.getElementById('quoteDraftTotal'); if(total)total.textContent=formatMoney(quoteDraftTotal());
    return;
  }
  const orderLine=e.target.closest('[data-order-line-field]');
  if(orderLine){
    const line=state.orderDraftLines[Number(orderLine.dataset.index)]; if(!line)return;
    line[orderLine.dataset.orderLineField]=orderLine.dataset.orderLineField==='productId'?orderLine.value:Number(orderLine.value);
    const subtotal=document.querySelector(`[data-order-line-subtotal="${orderLine.dataset.index}"]`); if(subtotal)subtotal.textContent=formatMoney(lineSubtotal(line));
    const total=document.getElementById('orderDraftTotal'); if(total)total.textContent=formatMoney(orderDraftTotal());
    return;
  }
  if(e.target.id==='quoteFreight'||e.target.id==='quoteTax'){const total=document.getElementById('quoteDraftTotal');if(total)total.textContent=formatMoney(quoteDraftTotal());return;}
  const input=e.target.closest('[data-relation-input]');
  if(!input)return;
  const target=input.dataset.relationInput;
  const items=target==='quoteCustomer'||target==='orderCustomer' ? customers.filter(c=>!c.archived).map(c=>c.name) : target==='orderQuote' ? quotes.filter(q=>q.status==='Accepted').map(q=>`${q.id} · ${q.customer}`) : orders.map(o=>`${o.id} · ${o.customer}`);
  const options=document.getElementById(`${target}Options`);
  if(options)options.innerHTML=relationOptions(target,items,input.value);
});

document.addEventListener('change', e => {
  const quoteProduct=e.target.closest('[data-quote-line-field="productId"]');
  if(quoteProduct){const line=state.quoteDraftLines[Number(quoteProduct.dataset.index)];const product=productByID(quoteProduct.value);Object.assign(line,{productId:product.id,unitPrice:moneyNumber(product.price)});renderQuoteDraftLines();return;}
  const orderProduct=e.target.closest('[data-order-line-field="productId"]');
  if(orderProduct){const line=state.orderDraftLines[Number(orderProduct.dataset.index)];const product=productByID(orderProduct.value);Object.assign(line,{productId:product.id,unitPrice:moneyNumber(product.price)});renderOrderDraftLines();return;}
  if(e.target.id==='orderSource'&&e.target.value==='手动创建'){const quote=document.getElementById('orderQuote');if(quote){quote.value='';document.getElementById('orderQuoteOptions').innerHTML='';}}
});

document.addEventListener('click', e => {
  const pageTarget=e.target.closest('[data-page]');
  if(pageTarget){setPage(pageTarget.dataset.page);return;}
  const lang=e.target.closest('[data-lang]');
  if(lang){state.lang=lang.dataset.lang;document.querySelectorAll('[data-lang]').forEach(x=>x.classList.toggle('active',x===lang));applyTranslations();toast('语言已切换',state.lang==='zh'?'当前界面为中文。':'Interface language is now English.');return;}
  const cat=e.target.closest('[data-agent-category]'); if(cat){state.agentCategory=cat.dataset.agentCategory;renderPage();return;}
  const documentType=e.target.closest('[data-document-type]'); if(documentType){state.documentType=documentType.dataset.documentType;renderPage();return;}
  const qv=e.target.closest('[data-quote-view]'); if(qv){state.quoteView=qv.dataset.quoteView;renderPage();return;}
  const pv=e.target.closest('[data-product-view]'); if(pv){state.productView=pv.dataset.productView;renderPage();return;}
  const customerTab=e.target.closest('[data-customer-tab]'); if(customerTab){customerDetail(customerTab.dataset.customerId,customerTab.dataset.customerTab);return;}
  const st=e.target.closest('[data-settings-tab]'); if(st){state.settingsTab=st.dataset.settingsTab;renderPage();loadPageOpenClawData();return;}
  const cmd=e.target.closest('[data-command-page]'); if(cmd){closeCommand();setPage(cmd.dataset.commandPage);return;}
  const el=e.target.closest('[data-action]');
  if(!el){
    const passive=e.target.closest('button');
    if(passive && !passive.closest('.modal-header') && !passive.closest('.drawer-header')) toast('操作已触发','该按钮当前仅更新前端交互状态，正式接口接入后将执行完整业务操作。');
    return;
  }
  const action=el.dataset.action;
  const actions={
    'close-modal':closeModal,'close-drawer':closeDrawer,
    'open-sidebar':()=>document.getElementById('sidebar').classList.add('open'),
    'close-sidebar':()=>document.getElementById('sidebar').classList.remove('open'),
    'sort-table':()=>toggleTableSort(el.dataset.module,el.dataset.field),
    'refresh':()=>toast('刷新完成','页面数据已更新；当前原型未调用后端接口。'),'pagination-current':()=>toast('当前已是第 1 页','正式版接入分页接口后可切换更多页面。'),
    'notifications':()=>openDrawer({title:'通知',eyebrow:'最近 24 小时',body:`<div class="timeline"><div class="timeline-item"><h4>2 份文件需要人工复核</h4><p>数据库 · 10:02</p></div><div class="timeline-item"><h4>订单 SO-2026-0102 已发运</h4><p>订单 · 09:18</p></div><div class="timeline-item"><h4>模型连接测试成功</h4><p>设置 · 08:36</p></div></div>`}),
    'lock':logoutUser,
    'metric-detail':()=>showMetric(el.dataset.key),
    'toggle-recommendations':()=>{state.recExpanded=!state.recExpanded;renderPage();},
    'recommend-detail':()=>newsDetail(recommendations[Number(el.dataset.index)].title),
    'recommend-settings':()=>{setPage('news');setTimeout(()=>document.querySelector('[data-action="news-sources"]')?.click(),0);},
    'oem-preset':()=>{state.oemQuery=el.dataset.value;const input=document.getElementById('oemQuery');if(input)input.value=state.oemQuery;renderPage();},
    'oem-run':()=>{state.oemQuery=document.getElementById('oemQuery')?.value.trim()||'';state.oemCategory=document.getElementById('oemCategory')?.value||state.oemCategory;toast('OEM 匹配完成','已按当前骑行类目、排序和 Top 数量刷新候选工厂。');renderPage();},
    'oem-detail':()=>oemDetail(el.dataset.name),
    'oem-export':oemExport,
    'unified-customer-search':()=>{state.customerSearchQuery=document.getElementById('unifiedCustomerQuery')?.value.trim()||'';state.customerSearchMode=document.getElementById('unifiedSearchMode')?.value||state.customerSearchMode;state.customerHasContact=Boolean(document.getElementById('hasContactOnly')?.checked);toast('客户搜索完成',`${state.customerSearchMode==='local'?'本地知识库':state.customerSearchMode==='rag'?'联网检索':'本地知识库 + 联网检索'}结果已返回。`);renderPage();},
    'unified-customer-detail':()=>unifiedCustomerDetail(el.dataset.name),
    'local-discovery-search':()=>{state.discoveryCountry=document.getElementById('discoveryCountry')?.value||state.discoveryCountry;state.discoveryCity=document.getElementById('discoveryCity')?.value||state.discoveryCity;state.discoveryType=document.getElementById('discoveryType')?.value||state.discoveryType;toast('本地客户发现完成','已提交筛选条件，并展示 CustomerMeasurementAgent 返回的候选结果。');renderPage();},
    'local-lead-detail':()=>localLeadDetail(el.dataset.name),
    'agent-chat':()=>showAgentChat(Number(el.dataset.agent),el.dataset.prompt||''),
    'weekly-report':()=>openModal({title:'生成智能体周报',eyebrow:'全部智能体 / 本周',body:`<div class="form-grid">${inputField('统计周期','2026-08-03 至 2026-08-09',true,true)}${selectField('内容范围',['全部智能体','按分类选择'])}${selectField('输出语言',['中文','英文','中文 / 英文双语'])}${selectField('输出格式',['Markdown + PDF','Markdown'])}<div class="form-field full"><label>周报包括</label><div class="filter-row">${['使用概览','重要对话','完成事项','待跟进','引用来源'].map(v=>`<label class="filter-chip active"><input type="checkbox" checked>${v}</label>`).join('')}</div></div></div>`,footer:formFooter('生成周报')}),
    'agent-manage':openAgentManager,
    'sync-openclaw-agents':()=>syncOpenClawAgents(el),
    'new-customer':()=>newCustomerForm(), 'edit-customer':()=>newCustomerForm(customers.find(c=>c.id===el.dataset.id)), 'customer-detail':()=>customerDetail(el.dataset.id), 'delete-customer':()=>deleteCustomer(el.dataset.id), 'customer-more':()=>customerDetail(el.dataset.id,'activity'),
    'export-customers':()=>toast('导出任务已创建','客户列表将按当前筛选条件导出为 Excel。'),
    'column-settings':()=>openModal({title:'客户列表字段',eyebrow:'显示设置',body:`<div class="filter-row">${['客户编号','客户类型','国家','联系人','电话','订单数','累计金额','评级','更新时间'].map(v=>`<label class="filter-chip active"><input type="checkbox" checked>${v}</label>`).join('')}</div>`,footer:formFooter('应用')}),
    'new-quote':()=>newQuoteForm(null,el.dataset.customer||''),'quote-detail':()=>quoteDetail(el.dataset.id),'edit-quote':()=>newQuoteForm(quotes.find(q=>q.id===el.dataset.id)),'delete-quote':()=>deleteQuote(el.dataset.id),'download-quote':()=>downloadQuote(el.dataset.id),'convert-order':()=>{closeDrawer();newOrderForm(null,quotes.find(q=>q.id===el.dataset.id)?.id||'');},
    'quote-metric-filter':()=>{state.quoteStatus=el.dataset.status||'all';renderPage();},
    'new-order':()=>newOrderForm(),'edit-order':()=>newOrderForm(orders.find(o=>o.id===el.dataset.id)),'delete-order':()=>deleteOrder(el.dataset.id),'order-detail':()=>orderDetail(el.dataset.id),
    'generate-docs':()=>generateDocs(el.dataset.id),'new-document':()=>generateDocs(''),'edit-document':()=>generateDocs('',documents.find(d=>d.id===el.dataset.id)),'delete-document':()=>deleteDocument(el.dataset.id),'download-document':()=>downloadDocument(el.dataset.id),'template-center':()=>templateCenter(el.dataset.kind),'document-detail':()=>documentDetail(el.dataset.id),
    'clear-document-filters':()=>{state.documentSearch='';state.documentType='all';state.documentStatus='all';renderPage();},
    'upload-template-image':()=>document.getElementById('templateImageInput')?.click(),
    'upload-template-file':()=>document.getElementById('templateFileInput')?.click(),
    'new-product':()=>newProductForm(),'edit-product':()=>newProductForm(products.find(p=>p.id===el.dataset.id)),'product-detail':()=>productDetail(el.dataset.id),'delete-product':()=>deleteProduct(el.dataset.id),'save-product':saveProduct,'import-products':()=>uploadFileModal(),'toggle-product-sort':()=>{state.productSort=state.productSort==='stockAsc'?'stockDesc':'stockAsc';renderPage();},
    'new-supplier':()=>newSupplierForm(),'edit-supplier':()=>newSupplierForm(suppliers.find(s=>s.id===el.dataset.id)),'supplier-detail':()=>supplierDetail(el.dataset.id),'delete-supplier':()=>deleteSupplier(el.dataset.id),'export-suppliers':()=>downloadText('suppliers.csv',['公司,电话,联系人,邮件,产品,规格,报价,来源,备注',...suppliers.map(s=>[s.company,s.phone,s.contact,s.email,s.product,s.specification,s.quote,s.source,s.notes].map(v=>`"${String(v).replaceAll('"','""')}"`).join(','))].join('\n')),
    'upload-file':uploadFileModal,'choose-file':()=>toast('已选择示例文件','文件校验通过，等待开始处理。'),'file-preview':()=>filePreview(el.dataset.name),'file-download':()=>fileDownload(el.dataset.name),'file-edit':()=>fileEdit(el.dataset.name),'file-summary':()=>fileSummary(el.dataset.name),'file-more':()=>openModal({title:'文件更多操作',eyebrow:'数据库 / 文件操作',body:`<div class="filter-row"><button class="button" data-action="file-summary" data-name="${escapeAttr(el.dataset.name)}">${icon('sparkles')}生成摘要</button><button class="button" data-action="file-download" data-name="${escapeAttr(el.dataset.name)}">${icon('download')}下载文件</button><button class="button danger" data-action="file-archive" data-name="${escapeAttr(el.dataset.name)}">${icon('archive')}归档文件</button></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`}),'file-archive':()=>toast('归档操作已记录',`${el.dataset.name} 将进入归档队列，正式版由 Go 文件接口处理。`),
    'open-category':()=>openModal({title:el.dataset.category,eyebrow:'数据库 / 分类',body:`<div class="empty-state">${icon('folder-search')}<div><h3>${el.dataset.category}文件</h3><p>正式接口接入后显示该主分类及关联标签下的文件列表。</p><button class="button primary" data-action="upload-file">${icon('upload')}上传到此分类</button></div></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`}),
    'agent-backup':()=>toast('智能体数据备份已开始','将备份初始化配置、技能、会话和操作记录。'),
    'tag-manage':()=>openModal({title:'标签管理',eyebrow:'数据库',body:`<div class="filter-row">${['客户私有','通用资料','合同','报价单','产品','欧盟','E-bike','兼容'].map(v=>`<span class="filter-chip active">${v}</span>`).join('')}</div>`,footer:formFooter('保存标签')}),
    'news-detail':()=>newsDetail(el.dataset.title),'toggle-news':()=>{state.newsExpanded=!state.newsExpanded;renderPage();},'news-filter':()=>{state.newsCategory=el.dataset.category||'全部';state.newsExpanded=false;renderPage();},
    'news-source-link':()=>toast('原文链接待配置','正式版将打开已通过来源白名单校验的原文地址。'),'news-todo':()=>toast('待办已生成','已将该行业资讯加入今日日程，正式版将写入任务和提醒接口。'),
    'refresh-news':()=>toast('新闻更新任务已启动','推荐智能体将按来源白名单获取和去重。'),
    'news-sources':openNewsSettings,
    'save-news-settings':saveNewsSettings,
    'configure-model':configureModel,
    'account-settings':openAccountSettings,
    'save-model-config':()=>saveModelConfiguration(el),
    'save-account-settings':()=>void saveAccountSettings(),
    'toggle-api-key-input':()=>{const input=document.getElementById('modelAPIKey');if(input){input.type=input.type==='password'?'text':'password';el.innerHTML=icon(input.type==='password'?'eye':'eye-off');applyIcons();}},
    'test-model':async()=>{await loadOpenClawModels(true);toast(state.modelConfigured?'模型配置可用':'模型配置待处理',state.modelConfigured?'OpenClaw 已解析默认模型和对应凭据。':'请检查默认模型和提供商凭据。',state.modelConfigured?'success':'warning');},
    'refresh-openclaw-models':async()=>{await loadOpenClawModels(true);toast('模型信息已刷新',state.openClawModels?.error||`当前默认模型：${state.openClawModels?.resolvedDefault||'未配置'}`,state.openClawModels?.error?'warning':'success');},
    'refresh-openclaw-system':async()=>{await Promise.all([loadOpenClawStatus(true),loadOpenClawAgents(true)]);toast('系统状态已刷新',state.openClawStatus?.rpcOK?'OpenClaw 网关与 RPC 正常。':'OpenClaw 状态需要检查。',state.openClawStatus?.rpcOK?'success':'warning');},
    'bind-plugin':()=>toast(`${el.dataset.plugin}插件`,el.textContent.trim()==='绑定'?'已进入 OpenClaw 绑定流程。':'已打开插件管理。'),
    'new-schedule':()=>schedulerForm(),'edit-schedule':()=>schedulerForm(el.dataset.name),'choose-backup':()=>toast('备份目录已更新','已选择本机外置存储目录。'),
    'offline-upgrade':offlineUpgradeModal,'choose-upgrade-package':()=>document.getElementById('upgradeFileInput')?.click(),'offline-install':installOfflineUpgrade,
    'upgrade-history':()=>openModal({title:'升级记录',eyebrow:'版本升级 / 审计日志',body:`<div class="empty-state">${icon('history')}<div><h3>暂无升级记录</h3><p>正式版记录升级包版本、签名摘要、操作人、快照、结果和回滚信息。</p></div></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`}),
    'relation-select':()=>{const input=document.getElementById(el.dataset.target);if(input){input.value=el.dataset.value;const options=document.getElementById(`${el.dataset.target}Options`);if(options)options.innerHTML='';if(el.dataset.target==='orderQuote')syncOrderFromQuote(el.dataset.value);}},
    'save-customer':saveCustomer,'save-quote':saveQuote,'save-order':saveOrder,'save-document':saveDocument,'save-supplier':saveSupplier,
    'add-quote-line':()=>{const product=products.find(item=>item.status==='Active')||products[0];state.quoteDraftLines.push({productId:product.id,quantity:1,unitPrice:moneyNumber(product.price),discount:0});renderQuoteDraftLines();},
    'remove-quote-line':()=>{if(state.quoteDraftLines.length===1){toast('至少保留一条产品明细','正式报价单需要至少一个产品。','warning');return;}state.quoteDraftLines.splice(Number(el.dataset.index),1);renderQuoteDraftLines();},
    'add-order-line':()=>{const product=products.find(item=>item.status==='Active')||products[0];state.orderDraftLines.push({productId:product.id,quantity:1,unitPrice:moneyNumber(product.price)});renderOrderDraftLines();},
    'remove-order-line':()=>{if(state.orderDraftLines.length===1){toast('至少保留一条产品明细','订单需要至少一个产品。','warning');return;}state.orderDraftLines.splice(Number(el.dataset.index),1);renderOrderDraftLines();},
    'download-quote':()=>downloadQuote(el.dataset.id),'download-document':()=>downloadDocument(el.dataset.id),
    'toggle-api-key':()=>{state.showApiKey=!state.showApiKey;configureModel();},
    'template-default':()=>templateAction('default'),'template-edit':()=>templateAction('edit'),
    'quote-date-filter':()=>toast('有效期筛选','可在后端接入后按日期区间筛选当前报价单。'),'order-date-filter':()=>toast('交付日期筛选','可在后端接入后按日期区间筛选当前订单。'),
    'save-form':()=>{closeModal();toast('已保存','当前表单内容已记录；正式版将调用对应 Go REST API。');renderPage();},
    'send-chat':()=>void sendAgentMessage(Number(el.dataset.agent)),
    'retry-chat':()=>void sendAgentMessage(Number(el.dataset.agent),el.dataset.message||''),
    'chat-quick-prompt':()=>void sendAgentMessage(Number(el.dataset.agent),el.dataset.prompt||''),
    'agent-allowlist':()=>openAgentAllowlist(Number(el.dataset.agent)),
    'save-agent-allowlist':()=>saveAgentAllowlist(Number(el.dataset.agent)),
    'mock-ocr':()=>toast('图片识别完成','已识别公司名、联系人、电话和邮箱，请人工确认。'),
  };
  if(actions[action])actions[action]();
});

document.getElementById('globalSearch').addEventListener('focus',e=>{openCommand(e.target.value);e.target.blur();});
document.getElementById('globalSearch').addEventListener('keydown',e=>{if(e.key==='Enter')openCommand(e.target.value);});
document.getElementById('commandInput').addEventListener('input',e=>renderCommandResults(e.target.value));
document.getElementById('modalBackdrop').addEventListener('click',e=>{if(e.target===e.currentTarget)closeModal();});
document.getElementById('drawerBackdrop').addEventListener('click',e=>{if(e.target===e.currentTarget)closeDrawer();});
document.getElementById('commandBackdrop').addEventListener('click',e=>{if(e.target===e.currentTarget)closeCommand();});
window.addEventListener('hashchange',()=>{
  const page=location.hash.slice(1);
  if(pageMeta[page]&&page!==state.page)setPage(page);
});
document.addEventListener('keydown',e=>{
  if((e.ctrlKey||e.metaKey)&&e.key.toLowerCase()==='k'){e.preventDefault();openCommand();}
  if(e.key==='Escape'){closeModal();closeDrawer();closeCommand();document.getElementById('sidebar').classList.remove('open');}
});

const initialPage=location.hash.slice(1);
if(pageMeta[initialPage])state.page=initialPage;
try {
  const savedNews = JSON.parse(localStorage.getItem('sta100-news-settings') || 'null');
  if (savedNews) {
    state.newsCountries = savedNews.countries || state.newsCountries;
    state.newsTopics = savedNews.topics || state.newsTopics;
    state.newsShowLimit = Number.isInteger(savedNews.showLimit) && savedNews.showLimit >= 1 && savedNews.showLimit <= 100 ? savedNews.showLimit : state.newsShowLimit;
    state.newsFrequency = ['1小时','2小时','3小时','6小时','8小时','12小时','24小时'].includes(savedNews.frequency) ? savedNews.frequency : state.newsFrequency;
    state.newsSources = savedNews.sources || state.newsSources;
  }
  state.agentInternetAllowlists = JSON.parse(localStorage.getItem('sta100-agent-allowlists') || '{}') || {};
} catch {
  localStorage.removeItem('sta100-news-settings');
  localStorage.removeItem('sta100-agent-allowlists');
}
applyTranslations();
void initAuth();
