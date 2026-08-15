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
  quoteDateFrom: '',
  quoteDateTo: '',
  orderSearch: '',
  orderStatus: 'all',
  orderSort: { field: 'value', direction: 'desc' },
  orderDateFrom: '',
  orderDateTo: '',
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
  openClawChannels: null,
  openClawChannelsLoading: false,
  openClawChannelsError: '',
  channelQR: null,
  channelQRPollTimer: null,
  modelProviderSelection: '',
  modelSearch: '',
  modelDraftMode: 'create',
  modelDraftFamilyKey: '',
  modelDraftKey: '',
  modelDraftOriginalKey: '',
  channelSearch: '',
  modelTestResult: null,
  modelTestLoading: false,
  testingModelKey: '',
  openClawAgents: null,
  openClawAgentsLoading: false,
  agentChats: {},
  agentChatHistoryLoading: {},
  agentChatHistoryLoaded: {},
  agentChatAtBottom: {},
  agentInternetAllowlists: {},
  agentModelSelections: {},
  agentSourceSelections: {},
  agentChatProgress: {},
  chatAttachments: [],
  showApiKey: false,
  formContext: null,
  selectedRows: {
    customers: new Set(),
    quotes: new Set(),
    orders: new Set(),
  },
  customerVisibleColumns: new Set(['customer','type','country','contact','orders','total','rating','updated']),
  customerCommunications: {},
  quoteDraftLines: [],
  orderDraftLines: [],
  templateUploads: [],
  templateKind: 'document',
  selectedUploadFile: null,
  selectedUpgradeFile: null,
  systemHealth: null,
  systemHealthLoading: false,
  lastAgentBackup: null,
  overviewDataStatus: '',
  lastWeeklyReport: null,
  assistantResults: {},
  tokenUsage: null,
  commandSearchSeq: 0,
  businessDataLoaded: false,
};

const authState = {
  username: 'admin',
  mode: 'login',
  message: '',
  authenticated: false,
  masterPassword: '',
};

const chatProgressSteps = [
  ['received','接收消息'],
  ['local-retrieval','本地检索'],
  ['attachments','附件处理'],
  ['knowledge-agent','知识整理'],
  ['domain-agents','业务 Agent'],
  ['coordinator-agent','统一汇总'],
];
const chatProgressTimers = {};

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

const scheduledJobs = [];

const news = [
  { category: '欧洲市场', title: '欧洲自行车产业进入补库存周期，渠道更关注小批量和快速交付', summary: '多家欧洲经销商在 2026 年下半年调整采购节奏，订单结构从大批量预采转向小批量、多批次。', source: 'Bike Europe', time: '2026-08-10 09:10', relevance: '96%' },
  { category: '法规', title: '欧盟更新电池尽职调查实施指引，E-bike 供应链资料需同步准备', summary: '新指引强化材料来源、碳足迹和供应链证明要求。', source: 'EUR-Lex', time: '2026-08-10 08:35', relevance: '93%' },
  { category: '智能骑行', title: '无线电子变速与功率数据融合成为高端整车配置趋势', summary: '整车厂正在把兼容性和训练数据完整度作为高端产品卖点。', source: 'Cycling Industry News', time: '2026-08-09 17:40', relevance: '89%' },
  { category: '渠道', title: '北欧经销商加快建设线上线下一体的维修服务网络', summary: '服务能力和备件响应速度正在影响品牌进入门槛。', source: 'Nordic Cycling', time: '2026-08-09 15:20', relevance: '84%' },
  { category: '产品', title: '欧洲城市通勤市场对轻量化 E-bike 的关注持续上升', summary: '重量、可维护性和电池合规成为渠道选品主要指标。', source: 'E-bike News', time: '2026-08-09 11:05', relevance: '82%' },
];

// 业务内容只以登录后 bootstrap 返回的数据为准，接口失败时保持空状态。
[recommendations, customers, quotes, orders, documents, products, suppliers, files, news].forEach(records => records.splice(0));
metrics.forEach(metric => { metric.value = 0; });

const unifiedSearchCustomers = [];
const localDiscoveryLeads = [];

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
  screen.innerHTML = `<section class="auth-card"><header class="auth-head"><div class="auth-brand"><span class="auth-brand-icon"><img src="assets/cycling-agent-icon.jpg" alt="STA-100"></span><div><strong>STA-100</strong><span>骑行行业智能工作台</span></div></div><h1 id="authTitle">${loginView ? '登录工作台' : resetView ? '重置本机账户' : '账户恢复'}</h1><p>${loginView ? `请输入 ${escapeHTML(authState.username)} 的密码继续使用。` : resetView ? '请设置新的用户名和登录密码。' : '请输入设备维护万能密码进入重置流程。'}</p></header><div class="auth-body">${body}</div></section>`;
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
    if (status.authenticated) await showAuthenticatedApp();
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
async function showAuthenticatedApp() {
  authState.authenticated = true;
  document.getElementById('authScreen').hidden = true;
  document.getElementById('appShell').hidden = false;
  updateSidebarIdentity();
  setPage(state.page);
  await Promise.all([loadBusinessData(), loadTokenUsage()]);
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
    await showAuthenticatedApp();
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
    await showAuthenticatedApp();
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
  if (options.body && !(options.body instanceof FormData)) headers['Content-Type'] = 'application/json';
  if (options.method && options.method !== 'GET') headers['X-STA100-Request'] = '1';
  const response = await fetch(path, { ...options, headers });
  let data = null;
  try { data = await response.json(); } catch { data = null; }
  if (!response.ok) throw new Error(data?.error?.message || `请求失败 (${response.status})`);
  return data;
}

function replaceRecords(target, records) {
  target.splice(0, target.length, ...(Array.isArray(records) ? records : []));
}

function upsertRecord(records, record) {
  const index = records.findIndex(item => item.id === record.id);
  if (index >= 0) records.splice(index, 1, record);
  else records.unshift(record);
}

function removeRecord(records, id) {
  const index = records.findIndex(item => item.id === id);
  if (index >= 0) records.splice(index, 1);
}

async function loadBusinessData(force = false) {
  if (state.businessDataLoaded && !force) return;
  try {
    const data = await apiFetch('/api/v1/bootstrap');
    replaceRecords(customers, data.customers);
    replaceRecords(quotes, data.quotes);
    replaceRecords(orders, data.orders);
    replaceRecords(documents, data.documents);
    replaceRecords(products, data.products);
    replaceRecords(suppliers, data.suppliers);
    replaceRecords(files, data.files);
    replaceRecords(news, data.news);
    replaceRecords(recommendations, data.recommendations);
    replaceRecords(scheduledJobs, data.jobs);
    replaceRecords(unifiedSearchCustomers, []);
    replaceRecords(localDiscoveryLeads, []);
    const overview = data.overview || {};
    metrics.forEach(metric => {
      if (Number.isFinite(Number(overview[metric.key]))) metric.value = Number(overview[metric.key]);
    });
    state.overviewDataStatus = overview.dataStatus || '';
    const preferences = data.preferences || {};
    state.subscription = preferences.recommendationEnabled ?? state.subscription;
    state.newsShowLimit = preferences.newsShowLimit || state.newsShowLimit;
    state.newsFrequency = preferences.newsFrequency || state.newsFrequency;
    state.newsCountries = preferences.newsCountries || state.newsCountries;
    state.newsTopics = preferences.newsTopics || state.newsTopics;
    state.newsSources = preferences.newsSources || state.newsSources;
    state.agentInternetAllowlists = preferences.agentAllowlists || state.agentInternetAllowlists;
    state.agentModelSelections = preferences.agentModelOverrides || state.agentModelSelections;
    state.businessDataLoaded = true;
    renderPage();
  } catch (error) {
    toast('业务数据读取失败', error.message, 'warning');
  }
}

function currentPreferences() {
  return { recommendationEnabled:state.subscription, newsShowLimit:state.newsShowLimit, newsFrequency:state.newsFrequency, newsCountries:state.newsCountries, newsTopics:state.newsTopics, newsSources:state.newsSources, agentAllowlists:state.agentInternetAllowlists, agentModelOverrides:state.agentModelSelections };
}

async function savePreferences() {
  return apiFetch('/api/v1/settings/preferences',{method:'PATCH',body:JSON.stringify(currentPreferences())});
}

async function refreshBusinessData() {
  try { await Promise.all([loadBusinessData(true), loadTokenUsage()]); toast('刷新完成','页面已从本地业务数据库重新读取。'); }
  catch(error) { toast('刷新失败',error.message,'warning'); }
}

function formatTokenCount(value) {
  return Number(value || 0).toLocaleString('zh-CN');
}

function applyTokenUsage(usage) {
  if (usage) state.tokenUsage = usage;
  const target = document.getElementById('tokenUsageTotal');
  if (target) target.textContent = formatTokenCount(state.tokenUsage?.total || 0);
}

async function loadTokenUsage() {
  try {
    applyTokenUsage(await apiFetch('/api/v1/agent-token-usage'));
  } catch (error) {
    applyTokenUsage(null);
    toast('Token 统计读取失败', error.message, 'warning');
  }
}

function openTokenUsage() {
  const usage = state.tokenUsage || {input:0,output:0,cacheRead:0,cacheWrite:0,total:0,calls:0,measuredCalls:0,unavailableCalls:0,currentRequestTotal:0,byAgent:[]};
  const agents = usage.byAgent || [];
  openModal({title:'智能体 Token 测试统计',eyebrow:'测试工具 / 实际模型用量',wide:true,body:`
    <div class="token-metrics">
      <div class="token-metric"><span>累计 Token</span><strong>${formatTokenCount(usage.total)}</strong></div>
      <div class="token-metric"><span>输入 Token</span><strong>${formatTokenCount(usage.input)}</strong></div>
      <div class="token-metric"><span>输出 Token</span><strong>${formatTokenCount(usage.output)}</strong></div>
      <div class="token-metric"><span>最近一次请求</span><strong>${formatTokenCount(usage.currentRequestTotal)}</strong></div>
    </div>
    <div class="token-test-note">该面板仅用于前期测试。统计值来自 OpenClaw 返回的模型 usage，不根据文本长度估算；缓存读取 ${formatTokenCount(usage.cacheRead)}，缓存写入 ${formatTokenCount(usage.cacheWrite)}。共 ${formatTokenCount(usage.calls)} 次 Agent 调用，其中 ${formatTokenCount(usage.measuredCalls)} 次返回有效用量，${formatTokenCount(usage.unavailableCalls)} 次未返回用量。</div>
    <div class="token-agent-list">
      <div class="token-agent-row header"><span>Agent</span><span>调用</span><span>输入</span><span>输出</span><span>总计</span></div>
      ${agents.length ? agents.map(item=>`<div class="token-agent-row"><strong>${escapeHTML(item.agentId)}</strong><span>${formatTokenCount(item.calls)}${item.unavailableCalls?` <small class="unavailable">(${formatTokenCount(item.unavailableCalls)} 未返回)</small>`:''}</span><span>${formatTokenCount(item.input)}</span><span>${formatTokenCount(item.output)}</span><span>${formatTokenCount(item.total)}</span></div>`).join('') : '<div class="empty-state"><div><h3>暂无 Token 统计</h3><p>完成一次可用的智能体调用后显示提供商返回的实际用量。</p></div></div>'}
    </div>`,footer:`<button class="button danger" data-action="clear-token-usage">${icon('trash-2')}清空测试统计</button><button class="button" data-action="refresh-token-usage">${icon('refresh-cw')}刷新</button><button class="button primary" data-action="close-modal">关闭</button>`});
}

async function clearTokenUsage() {
  if (!window.confirm('确定清空当前 Token 测试统计吗？该操作不会删除 Agent 会话。')) return;
  try {
    await apiFetch('/api/v1/agent-token-usage',{method:'DELETE'});
    await loadTokenUsage();
    closeModal();
    toast('Token 测试统计已清空','Agent 会话和业务数据未受影响。');
  } catch(error) { toast('清空失败',error.message,'warning'); }
}

async function runOEMMatch() {
  state.oemQuery=document.getElementById('oemQuery')?.value.trim()||'';
  state.oemCategory=document.getElementById('oemCategory')?.value||state.oemCategory;
  if(!state.oemQuery){toast('请输入匹配需求','产品、数量、市场或规格至少需要一项。','warning');return;}
  const target=document.getElementById('oemResults');if(target)target.innerHTML=`<div class="tool-empty">${icon('loader-circle')} 正在整理本地证据并调用专业智能体...</div>`;applyIcons();
  try {
    const result=await apiFetch('/api/v1/assistant/query',{method:'POST',body:JSON.stringify({page:'overview',feature:'oem-match',message:state.oemQuery,sessionKey:'sta100-overview-oem',context:{category:state.oemCategory,sort:state.oemSort,top:state.oemTop}})});
    state.assistantResults.oem=result;applyTokenUsage(result.tokenUsage);renderPage();toast(result.partial?'OEM 匹配返回部分结果':'OEM 匹配完成',`${result.usedAgents.length} 个 Agent 参与处理。`,result.partial?'warning':'success');
  } catch(error) { toast('OEM 匹配失败',error.message,'warning'); }
}

async function runUnifiedCustomerSearch() {
  state.customerSearchQuery=document.getElementById('unifiedCustomerQuery')?.value.trim()||'';
  state.customerHasContact=Boolean(document.getElementById('hasContactOnly')?.checked);
  const target=document.getElementById('unifiedCustomerResults');if(target)target.innerHTML=`<div class="tool-empty">${icon('loader-circle')} 正在执行统一智能搜索...</div>`;applyIcons();
  try {
    const result=await apiFetch('/api/v1/assistant/query',{method:'POST',body:JSON.stringify({page:'overview',feature:'customer-search',message:state.customerSearchQuery||'查询客户',sessionKey:'sta100-overview-customer-search',context:{hasContact:state.customerHasContact}})});
    replaceRecords(unifiedSearchCustomers,result.items);state.assistantResults.customerSearch=result;applyTokenUsage(result.tokenUsage);renderPage();toast(result.partial?'客户搜索返回部分结果':'客户搜索完成',`返回 ${result.items.length} 条记录，${result.usedAgents.length} 个 Agent 参与。`,result.partial?'warning':'success');
  } catch(error) { toast('客户搜索失败',error.message,'warning'); }
}

async function runLocalDiscovery() {
  state.discoveryCountry=document.getElementById('discoveryCountry')?.value||state.discoveryCountry;
  state.discoveryCity=document.getElementById('discoveryCity')?.value||state.discoveryCity;
  state.discoveryType=document.getElementById('discoveryType')?.value||state.discoveryType;
  const target=document.getElementById('localDiscoveryResults');if(target)target.innerHTML=`<div class="tool-empty">${icon('loader-circle')} 正在整理本地证据并分析客户线索...</div>`;applyIcons();
  try {
    const message=`发现 ${state.discoveryCountry} ${state.discoveryCity} 的 ${state.discoveryType} 客户`;
    const result=await apiFetch('/api/v1/assistant/query',{method:'POST',body:JSON.stringify({page:'overview',feature:'customer-discovery',message,sessionKey:'sta100-overview-customer-discovery',context:{country:state.discoveryCountry,city:state.discoveryCity,type:state.discoveryType,hasContact:false}})});
    replaceRecords(localDiscoveryLeads,result.items);state.assistantResults.customerDiscovery=result;applyTokenUsage(result.tokenUsage);renderPage();toast(result.partial?'客户发现返回部分结果':'客户发现完成',`${result.usedAgents.length} 个 Agent 参与处理。`,result.partial?'warning':'success');
  } catch(error) { toast('客户发现失败',error.message,'warning'); }
}

async function backupAgents() {
  try { const result=await apiFetch('/api/v1/agent-backups',{method:'POST',body:'{}'}); state.lastAgentBackup=result; if(state.page==='settings')renderPage(); toast('智能体备份完成',`${result.path} · ${formatBytes(result.bytes)}`); }
  catch(error) { toast('智能体备份失败',error.message,'warning'); }
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
    state.openClawModels = await apiFetch('/api/v1/settings/model');
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

async function loadOpenClawChannels(force = false) {
  if (state.openClawChannelsLoading || (state.openClawChannels && !force)) return state.openClawChannels;
  state.openClawChannelsLoading = true;
  state.openClawChannelsError = '';
  if (state.page === 'settings' && state.settingsTab === 'channels') renderPage();
  try {
    const data = await apiFetch('/api/v1/openclaw/channels');
    state.openClawChannels = data.channels || [];
  } catch (error) {
    state.openClawChannels = [];
    state.openClawChannelsError = error.message || '未知错误';
    if (state.page === 'settings' && state.settingsTab === 'channels') toast('OpenClaw 通道读取失败', error.message, 'warning');
  } finally {
    state.openClawChannelsLoading = false;
    if (state.page === 'settings' && state.settingsTab === 'channels') renderPage();
  }
  return state.openClawChannels;
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

async function loadSystemHealth(force = false) {
  if (state.systemHealthLoading || (state.systemHealth && !force)) return state.systemHealth;
  state.systemHealthLoading = true;
  if (state.page === 'settings' && state.settingsTab === 'system') renderPage();
  try {
    state.systemHealth = await apiFetch('/api/v1/system/health');
  } catch (error) {
    state.systemHealth = { status: 'error', error: error.message };
  } finally {
    state.systemHealthLoading = false;
    if (state.page === 'settings' && state.settingsTab === 'system') renderPage();
  }
  return state.systemHealth;
}

function loadPageOpenClawData() {
  if (state.page === 'agents') void loadOpenClawAgents();
  if (state.page === 'agents') void loadOpenClawModels();
  if (state.page !== 'settings') return;
  if (state.settingsTab === 'model') void loadOpenClawModels();
  if (state.settingsTab === 'channels') void loadOpenClawChannels();
  if (state.settingsTab === 'scheduler') void loadOpenClawAgents();
  if (state.settingsTab === 'system') {
    void loadOpenClawStatus();
    void loadOpenClawAgents();
    void loadSystemHealth();
  }
}

function setPage(page) {
  if (!pageMeta[page]) return;
  state.page = page;
  const [title, eyebrow, emoji] = pageMeta[page];
  const pageTitle = document.getElementById('pageTitle');
  const pageEyebrow = document.getElementById('pageEyebrow');
  if (pageTitle) pageTitle.textContent = `${emoji} ${title}`;
  if (pageEyebrow) pageEyebrow.textContent = eyebrow;
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

    <div class="section-head tool-section-head"><div><h3>🛠️ 智能业务工具</h3><p>由本地证据整理器、领域智能体和任务协调器共同完成信息整合。</p></div><span class="meta">统一结果 · 冲突信息并列保留</span></div>
    <section class="tool-grid">
      <article class="panel tool-panel tool-panel-wide">
        <header class="tool-header"><div><h3>🏭 OEM 工厂智能匹配</h3><p>用户输入需求后，系统先整理本地证据，再由协调器分发领域智能体并汇总。</p></div><span class="badge amber">数据待补充</span></header>
        <div class="panel-body">
          <div class="filter-row tool-presets">${['公路整车 OEM 1000 台','E-bike 电池 OEM 100 组','中置电机 500 套','头盔 MIPS 500 个'].map(v=>`<button class="filter-chip" data-action="oem-preset" data-value="${v}">${v}</button>`).join('')}</div>
          <div class="tool-form">
            <label class="field-search tool-query">${icon('search')}<input id="oemQuery" value="${escapeAttr(state.oemQuery)}" placeholder="输入产品、数量、市场和要求"></label>
            <select class="select" id="oemCategory"><option ${state.oemCategory==='全部骑行类目'?'selected':''}>全部骑行类目</option><option ${state.oemCategory==='整车'?'selected':''}>整车</option><option ${state.oemCategory==='E-bike 电池'?'selected':''}>E-bike 电池</option><option ${state.oemCategory==='电机'?'selected':''}>电机</option><option ${state.oemCategory==='链条/传动'?'selected':''}>链条/传动</option><option ${state.oemCategory==='轮胎'?'selected':''}>轮胎</option><option ${state.oemCategory==='头盔'?'selected':''}>头盔</option><option ${state.oemCategory==='码表/智能设备'?'selected':''}>码表/智能设备</option><option ${state.oemCategory==='功率计'?'selected':''}>功率计</option></select>
            <select class="select" id="oemSort"><option value="score" ${state.oemSort==='score'?'selected':''}>按匹配度排序</option><option value="capacity" ${state.oemSort==='capacity'?'selected':''}>按产能排序</option><option value="moq" ${state.oemSort==='moq'?'selected':''}>按 MOQ 排序</option><option value="source" ${state.oemSort==='source'?'selected':''}>按数据来源排序</option></select>
            <select class="select" id="oemTop"><option value="3" ${state.oemTop===3?'selected':''}>Top 3</option><option value="5" ${state.oemTop===5?'selected':''}>Top 5</option><option value="10" ${state.oemTop===10?'selected':''}>Top 10</option></select>
            <button class="button primary" data-action="oem-run">${icon('scan-search')}开始匹配</button>
          </div>
          <div class="agent-chain-note"><span class="agent-icon">${icon('workflow')}</span><span><strong>本地证据 → 任务协调器 → OEM 领域智能体 → 统一汇总</strong><small>正式工厂数据、分类和评分规则待提供；冲突信息全部保留。</small></span></div>
          <div class="tool-results" id="oemResults">${renderOEMCards()}</div>
        </div>
      </article>

      <article class="panel tool-panel">
        <header class="tool-header"><div><h3>🔍 客户统一搜索</h3><p>统一检索本地业务记录和后续可接入的私有知识，结果由智能体整合展示。</p></div></header>
        <div class="panel-body">
          <div class="tool-form compact">
            <label class="field-search tool-query">${icon('search')}<input id="unifiedCustomerQuery" value="${escapeAttr(state.customerSearchQuery)}" placeholder="国家、公司、业务、邮箱或电话"></label>
            <label class="contact-check"><input class="checkbox" id="hasContactOnly" type="checkbox" ${state.customerHasContact?'checked':''}> 必有联系方式</label>
            <button class="button primary" data-action="unified-customer-search">${icon('search')}搜索</button>
          </div>
          <div class="agent-chain-note"><span class="agent-icon">${icon('workflow')}</span><span><strong>用户输入 → 本地证据整理 → 协调器分发 → 客户结果汇总</strong><small>联系方式包含邮箱、电话、网站及其它通讯方式。</small></span></div>
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

function renderOEMCards() {
  const result=state.assistantResults.oem;
  if(!result)return `<div class="tool-empty">正式工厂原始数据、骑行类目和评分规则尚未提供。可先提交需求，由系统基于现有证据和专业智能体返回部分分析。</div>`;
  return renderAssistantSummary(result,'OEM 匹配分析');
}

function renderAssistantSummary(result,title) {
  const agents=(result.usedAgents||[]).join('、')||'无成功调用';
  return `<article class="assistant-summary"><div class="spread"><strong>${escapeHTML(title)}</strong><span class="badge ${result.partial?'amber':'green'}">${result.partial?'部分结果':'已完成'}</span></div><p>${escapeHTML(result.text||'暂无汇总文本').replace(/\n/g,'<br>')}</p><div class="source-line"><span class="mini-source">参与 Agent <strong>${escapeHTML(agents)}</strong></span><span class="mini-source">本地证据 <strong>${result.evidence?.length||0} 条</strong></span><span class="mini-source">冲突 <strong>${result.conflicts?.length||0} 组</strong></span></div></article>`;
}

function renderUnifiedCustomerCards() {
  let rows = unifiedSearchCustomers.filter(r => {
    const text = `${r.name} ${r.country} ${r.type} ${r.business} ${r.contact}`.toLowerCase();
    const tokens = state.customerSearchQuery.toLowerCase().split(/\s+/).filter(token => token.length > 1);
    const matchesQuery = !tokens.length || tokens.some(token => text.includes(token));
    const matchesContact = !state.customerHasContact || Boolean(r.contact);
    return matchesQuery && matchesContact;
  });
  const summary=state.assistantResults.customerSearch?renderAssistantSummary(state.assistantResults.customerSearch,'客户搜索摘要'):'';
  return `${summary}${rows.slice(0,3).map(r=>`<article class="customer-match-row"><span class="match-score">${r.score}<small>分</small></span><span class="customer-match-copy"><strong>${r.name}</strong><small>${r.country} · ${r.type} · ${r.business||'暂无业务描述'}</small><em>${icon('contact-round')} ${r.contact||'未填写'}</em></span><span class="badge blue">已整合</span><button class="table-icon" data-action="unified-customer-detail" data-name="${r.name}" title="查看详情">${icon('arrow-up-right')}</button></article>`).join('')}`;
}

function renderLocalDiscoveryCards() {
  const rows=localDiscoveryLeads.slice(0,3);
  const summary=state.assistantResults.customerDiscovery?renderAssistantSummary(state.assistantResults.customerDiscovery,'客户发现摘要'):'';
  return `${summary}${rows.map(r=>`<article class="customer-match-row"><span class="match-score">${r.score}<small>分</small></span><span class="customer-match-copy"><strong>${r.name}</strong><small>${r.country} · ${r.city||'城市待核实'} · ${r.type}</small><em>${icon('contact-round')} ${r.contact||'未填写'}</em></span><span class="badge blue">已整合</span><button class="table-icon" data-action="local-lead-detail" data-name="${r.name}" title="查看详情">${icon('arrow-up-right')}</button></article>`).join('')||(!summary?`<div class="tool-empty">客户原始数据格式和允许的公开来源尚未提供，可先执行查询获取部分分析。</div>`:'')}`;
}

function unifiedCustomerDetail(name) {
  const customer = unifiedSearchCustomers.find(r => r.name === name);
  if (!customer) return;
  openDrawer({ title: customer.name, eyebrow: '客户统一搜索 / 已整合', body: `<div class="spread"><span class="badge blue">已整合</span><span class="secondary-text">匹配 ${customer.score} 分</span></div><div class="detail-grid" style="margin-top:15px">${[['客户编号',customer.id||'未提供'],['国家',customer.country],['客户类型',customer.type],['业务方向',customer.business||'未填写'],['联系方式',customer.contact||'未填写'],['联系方式过滤',state.customerHasContact?'已启用':'未启用']].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(String(v))}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">证据说明</div><p class="secondary-text" style="line-height:1.8">页面展示协调器整合后的结果。每条底层证据保留记录编号、更新时间和来源元数据；发生冲突时全部并列展示。</p>` });
}

function localLeadDetail(name) {
  const lead = localDiscoveryLeads.find(item => item.name === name);
  if (!lead) return;
  openDrawer({ title: lead.name, eyebrow: '本地客户发现 / OpenClaw', body: `<div class="spread"><span class="badge blue">CustomerMeasurementAgent</span><span class="secondary-text">匹配 ${escapeHTML(String(lead.score ?? '待核实'))} 分</span></div><div class="detail-grid" style="margin-top:15px">${[['国家',lead.country],['城市',lead.city],['客户类型',lead.type],['联系方式',lead.contact||'未提供'],['OpenClaw 返回理由',lead.reason||'未提供']].map(([label,value])=>`<div class="detail-field"><label>${label}</label><strong>${escapeHTML(String(value||'未提供'))}</strong></div>`).join('')}</div><div class="agent-chain-note" style="margin-top:16px"><span class="agent-icon">${icon('bot')}</span><span><strong>筛选条件 → 允许的数据源 → CustomerMeasurementAgent</strong><small>结果必须同时保留来源链接、抓取时间和查询条件。</small></span></div>` });
}

async function oemExport() {
  try {
    await apiFetch('/api/v1/overview/oem-matches/export',{method:'POST',body:JSON.stringify({query:state.oemQuery,category:state.oemCategory,sort:state.oemSort,top:state.oemTop})});
  } catch (error) {
    toast('OEM 报告暂不可生成',error.message,'warning');
  }
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
  state.selectedUpgradeFile=valid?file:null;
  body.innerHTML = `<div class="upgrade-file"><span class="upload-icon">${icon(valid ? 'file-check-2' : 'file-warning')}</span><div><strong>${escapeHTML(file.name)}</strong><small>${formatBytes(file.size)} · 本机文件</small></div></div><div class="upgrade-checks"><div>${icon(valid?'check-circle-2':'circle-x')}<span>文件扩展名</span><strong>${valid?'通过本地检查':'失败'}</strong></div><div>${icon('cpu')}<span>目标架构</span><strong>等待后端解析</strong></div><div>${icon('shield-check')}<span>签名与完整性</span><strong>等待后端校验</strong></div><div>${icon('hard-drive')}<span>磁盘空间</span><strong>等待后端校验</strong></div></div><div class="model-warning" style="margin-top:14px"><span>${icon('info')} 当前不会模拟安装。后端需要升级包 manifest、签名算法、兼容范围、迁移和回滚规范后才允许导入。</span></div>`;
  footer.innerHTML = `<button class="button" data-action="close-modal">取消</button>${valid ? `<button class="button primary" data-action="import-upgrade-package">${icon('shield-check')}提交后端校验</button>` : ''}`;
  applyIcons();
}

async function importOfflineUpgrade() {
  if(!state.selectedUpgradeFile)return;
  const form=new FormData();form.append('file',state.selectedUpgradeFile);
  try {await apiFetch('/api/v1/system/upgrade/import',{method:'POST',body:form});}
  catch(error){toast('升级包暂不能导入',error.message,'warning');}
}

async function showUpgradeHistory() {
  try {const result=await apiFetch('/api/v1/system/upgrade/history');openModal({title:'升级记录',eyebrow:'版本升级 / 审计日志',body:result.items?.length?`<pre>${escapeHTML(JSON.stringify(result.items,null,2))}</pre>`:`<div class="empty-state">${icon('history')}<div><h3>暂无升级记录</h3><p>尚未执行过后端离线升级。</p></div></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});}
  catch(error){toast('升级记录读取失败',error.message,'warning');}
}

function renderAgents() {
  const filtered = agents.map((agent, index) => ({ agent, index })).filter(({ agent }) => state.agentCategory === 'all' || agent[2] === state.agentCategory);
  const managedAgents = state.openClawAgents?.filter(agent => !agent.isDefault) || [];
  return `<div class="page-stack">
    ${!state.modelConfigured ? `<div class="model-warning"><span>${icon('triangle-alert')} 当前尚未完成模型配置，智能体不能发起真实调用。</span><button class="button small" data-page="settings">进入模型设置</button></div>` : ''}
    <section class="agents-summary panel">
      <div><h2>🤖 24 个专业智能体，对应 24 个 OpenClaw Agent</h2><p>每个智能体保留独立初始化配置、技能和会话记录，由系统协调器统一组织本地证据与专业回复。</p></div>
      <div class="inline-actions"><button class="button" data-action="agent-manage">${icon('sliders-horizontal')}智能体管理</button><button class="button primary" data-action="weekly-report">${icon('file-clock')}生成本周周报</button></div>
    </section>
    <div class="toolbar">
      <div class="filter-row agent-category-tabs">
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
    <div class="data-wrap"><table class="data-table" id="customerTable"><thead><tr><th>${selectAllCheckbox('customers','客户')}</th><th>客户</th><th>类型</th><th>国家</th><th>联系人</th><th>${sortHeader('订单数量','customer','orders',state.customerSort)}</th><th>${sortHeader('累计金额','customer','total',state.customerSort)}</th><th>评级</th><th>${sortHeader('最近更新','customer','updated',state.customerSort)}</th><th>操作</th></tr></thead><tbody>
      ${rows.map(c=>`<tr><td>${rowCheckbox('customers',c.id,c.name)}</td><td><button class="link-button primary-cell" data-action="customer-detail" data-id="${escapeAttr(c.id)}"><span class="avatar">${escapeHTML(c.name.slice(0,2).toUpperCase())}</span><span><strong>${escapeHTML(c.name)}</strong><small>${escapeHTML(c.id)}</small></span></button></td><td>${escapeHTML(c.type)}</td><td>${escapeHTML(c.country)}</td><td><span class="primary-cell"><span><strong>${escapeHTML(c.contact)}</strong><small>${escapeHTML(c.email)}</small></span></span></td><td>${c.orders}</td><td>${escapeHTML(c.total)}</td><td>${badge(c.rating)}</td><td>${escapeHTML(c.updated)}</td><td><span class="table-actions"><button class="table-icon" data-action="customer-detail" data-id="${escapeAttr(c.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-customer" data-id="${escapeAttr(c.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="customer-more" data-id="${escapeAttr(c.id)}" title="更多">${icon('ellipsis')}</button></span></td></tr>`).join('') || `<tr><td colspan="10"><div class="empty-state">${icon('search-x')}<div><h3>未找到客户</h3><p>请调整搜索词或筛选条件。</p></div></div></td></tr>`}
    </tbody></table></div>
    <div class="pagination"><span>共 ${rows.length} 条记录 · 每页 20 条</span><div><button class="button small ghost" disabled>${icon('chevron-left')}</button><button class="button small" data-action="pagination-current">1</button><button class="button small ghost" disabled>${icon('chevron-right')}</button></div></div>
  </div>`;
}

function renderQuotes() {
  const query = state.quoteSearch.trim().toLowerCase();
  const visible = sortRows(quotes.filter(q => (!query || [q.id, q.subject, q.customer, q.products].join(' ').toLowerCase().includes(query)) && (state.quoteStatus === 'all' || q.status === state.quoteStatus) && (!state.quoteDateFrom || q.valid >= state.quoteDateFrom) && (!state.quoteDateTo || q.valid <= state.quoteDateTo)), state.quoteSort, { value: quote => moneyNumber(quote.value) });
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
  const visible = sortRows(orders.filter(o => (!query || [o.id, o.customer, o.quote, o.products].join(' ').toLowerCase().includes(query)) && (state.orderStatus === 'all' || o.status === state.orderStatus) && (!state.orderDateFrom || o.delivery >= state.orderDateFrom) && (!state.orderDateTo || o.delivery <= state.orderDateTo)), state.orderSort, { value: order => moneyNumber(order.value) });
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

async function saveSupplier() {
  const company=formText('supplierCompany'); if(!company){toast('保存失败','公司名称为必填项。','warning');return;}
  const existing=suppliers.find(s=>s.id===state.formContext?.id);
  const payload={...(existing||{}),company,phone:formText('supplierPhone'),contact:formText('supplierContact'),email:formText('supplierEmail'),product:formText('supplierProduct'),specification:formText('supplierSpecification'),quote:formText('supplierQuote'),source:formText('supplierSource'),notes:formText('supplierNotes')};
  try {
    const record=await apiFetch(existing?`/api/v1/suppliers/${encodeURIComponent(existing.id)}`:'/api/v1/suppliers',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)});
    upsertRecord(suppliers,record); closeModal(); renderPage(); toast(existing?'供应商已更新':'供应商已创建',`${record.company} 已保存到本地数据库。`);
  } catch(error) { toast('保存失败',error.message,'warning'); }
}

function supplierDetail(id) {
  const s=suppliers.find(item=>item.id===id); if(!s)return;
  openDrawer({title:s.company,eyebrow:`供应商 / ${s.id}`,body:`<div class="spread"><span class="badge blue">${escapeHTML(s.source)}</span><div class="inline-actions"><button class="button small" data-action="edit-supplier" data-id="${escapeAttr(s.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-supplier" data-id="${escapeAttr(s.id)}">${icon('trash-2')}删除</button></div></div><div class="detail-grid" style="margin-top:15px">${[['公司',s.company],['电话',s.phone],['联系人',s.contact],['邮件',s.email],['产品',s.product],['规格',s.specification],['报价',s.quote],['来源',s.source],['备注',s.notes],['更新时间',s.updated]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v||'未填写')}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">后续扩展字段</div><p class="secondary-text">可在正式业务库中继续扩展认证、产能、MOQ、合作状态、付款条件和历史报价。</p>`});
}

async function deleteSupplier(id) { if(!window.confirm('确定归档该供应商吗？'))return; try { await apiFetch(`/api/v1/suppliers/${encodeURIComponent(id)}`,{method:'DELETE'}); removeRecord(suppliers,id); closeDrawer(); renderPage(); toast('供应商已归档','供应商记录已写入本地数据库。'); } catch(error) { toast('归档失败',error.message,'warning'); } }

function renderDatabase() {
  const categoryIcons = {合同:'file-signature',报价单:'file-chart-column',产品手册:'book-open',法规:'scale',产品资料:'boxes',会议记录:'notebook-tabs',客户资料:'contact',图片:'images',其它:'folder','待分类':'folder-search'};
  const categoryCounts = files.reduce((result,file)=>{result[file.category]=(result[file.category]||0)+1;return result;},{});
  const categories = Object.entries(categoryCounts).sort(([left],[right])=>left.localeCompare(right,'zh-CN')).map(([name,count])=>[name,count,categoryIcons[name]||'folder']);
  const fileQuery = state.fileSearch.trim().toLowerCase();
  const visibleFiles = files.filter(f => !fileQuery || Object.values(f).join(' ').toLowerCase().includes(fileQuery));
  const totalBytes = files.reduce((sum,file)=>sum+Number(file.bytes||0),0);
  const indexedCount = files.filter(file=>file.status==='Indexed').length;
  const reviewCount = files.filter(file=>!['Indexed'].includes(file.status)).length;
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>📚 私有数据库</h2><p>出厂通用资料与客户私有数据分区存储，客户数据仅保存在本机。</p></div><div class="toolbar"><button class="button" data-action="agent-backup">${icon('archive')}备份智能体数据</button><button class="button primary" data-action="upload-file">${icon('upload')}上传文件</button></div></div>
    <section class="metric-grid" style="grid-template-columns:repeat(4,1fr)">${[['文件总数',String(files.length),'files'],['已建立索引',String(indexedCount),'database-zap'],['待处理/确认',String(reviewCount),'circle-help'],['本地占用',formatBytes(totalBytes),'hard-drive']].map(([l,v,i])=>`<div class="metric-button"><span class="metric-icon">${icon(i)}</span><span><strong class="metric-number">${v}</strong><span class="metric-label">${l}</span></span></div>`).join('')}</section>
    <div class="section-head"><div><h3>数据分类</h3><p>一个文件可以属于多个标签，分类用于主归档。</p></div><button class="button ghost small" data-action="tag-manage">${icon('tags')}标签管理</button></div>
    <section class="category-grid">${categories.map(([n,c,i])=>`<button class="category-card" data-action="open-category" data-category="${n}"><span class="category-icon">${icon(i)}</span><strong>${n}</strong><span>${c} 个文件</span></button>`).join('')}</section>
    <section class="panel"><header class="panel-head"><div><h3>最近文件</h3><p>展示解析、分类和索引结果</p></div><label class="field-search" style="height:32px">${icon('search')}<input id="fileSearch" value="${escapeAttr(state.fileSearch)}" placeholder="搜索文件或标签"></label></header><div class="data-wrap" style="border:0;border-radius:0"><table class="data-table"><thead><tr><th>文件名</th><th>主分类</th><th>标签</th><th>来源</th><th>大小</th><th>索引状态</th><th>更新时间</th><th>操作</th></tr></thead><tbody>${visibleFiles.map(f=>`<tr><td><button class="link-button" data-action="file-preview" data-id="${escapeAttr(f.id)}">${escapeHTML(f.name)}</button></td><td>${escapeHTML(f.category)}</td><td>${(f.tags||[]).map(t=>`<span class="badge neutral">${escapeHTML(t)}</span>`).join(' ')}</td><td>${escapeHTML(f.source)}</td><td>${escapeHTML(f.size)}</td><td>${badge(f.status)}</td><td>${escapeHTML(f.updated)}</td><td><span class="table-actions"><button class="table-icon" data-action="file-preview" data-id="${escapeAttr(f.id)}" title="预览">${icon('eye')}</button><button class="table-icon" data-action="file-edit" data-id="${escapeAttr(f.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="file-more" data-id="${escapeAttr(f.id)}" title="更多">${icon('ellipsis')}</button></span></td></tr>`).join('') || `<tr><td colspan="8"><div class="empty-state"><p>未找到匹配文件</p></div></td></tr>`}</tbody></table></div></section>
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

async function saveNewsSettings() {
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
  try {
    await savePreferences();
    localStorage.removeItem('sta100-news-settings');
    closeModal(); renderPage(); toast('新闻设置已保存',`每 ${state.newsFrequency} 获取，单次最多展示 ${state.newsShowLimit} 条。`);
  } catch(error) { toast('新闻设置保存失败',error.message,'warning'); }
}

function renderSettings() {
  const tabs = [['model','模型设置','brain-circuit'],['channels','通道绑定','message-square'],['scheduler','定时任务','clock-3'],['backup','智能体备份','archive'],['security','数据安全','shield-check'],['system','系统信息','monitor-cog'],['upgrade','版本升级','package-open']];
  return `<div class="settings-layout">
    <nav class="settings-nav panel">${tabs.map(([k,l,i])=>`<button class="${state.settingsTab===k?'active':''}" data-settings-tab="${k}">${icon(i)}${l}</button>`).join('')}</nav>
    <div class="settings-content">${renderSettingsContent()}</div>
  </div>`;
}

function renderSettingsContent() {
  const content = {
    model: renderModelSettings(),
    channels: renderChannelSettings(),
    scheduler: renderSchedulerSettings(),
    backup: `<section class="panel"><header class="panel-head"><div><h3>智能体数据备份</h3><p>仅备份智能体初始化配置、技能、会话和用户操作记录</p></div></header><div class="setting-row"><span class="setting-icon">${icon('folder-cog')}</span><div class="setting-copy"><strong>备份目录</strong><span>当前由设备后端写入 STA-100 本地备份区；外置目录授权方式待部署确认</span></div><button class="button small" data-action="choose-backup">目录说明</button></div><div class="setting-row"><span class="setting-icon">${icon('archive')}</span><div class="setting-copy"><strong>本次会话最近备份</strong><span>${state.lastAgentBackup?`${escapeHTML(state.lastAgentBackup.path)} · ${formatBytes(state.lastAgentBackup.bytes)}`:'尚未执行备份'}</span></div><button class="button primary small" data-action="agent-backup">立即备份</button></div></section>`,
    security: `<section class="panel"><header class="panel-head"><div><h3>数据安全</h3><p>客户私有数据归客户所有，并保存在设备本地</p></div></header><div class="setting-row"><span class="setting-icon">${icon('user-round-cog')}</span><div class="setting-copy"><strong>本机登录账户</strong><span>当前用户名：${escapeHTML(authState.username)} · 密码以随机盐哈希保存在设备本机配置文件</span></div><button class="button small" data-action="account-settings">${icon('key-round')}修改用户名和密码</button></div>${[['本地存储区','客户上传文件、业务数据库和索引均存储在本机','hard-drive'],['敏感信息保护','模型密钥由 OpenClaw 凭据库保存；页面和日志不回显完整密钥','lock-keyhole'],['联网调用边界','仅把完成当前任务所需的最小内容发送给已配置模型','network'],['操作记录','记录关键增删改、导出、模型和升级操作','scroll-text']].map(([n,d,i])=>`<div class="setting-row"><span class="setting-icon">${icon(i)}</span><div class="setting-copy"><strong>${n}</strong><span>${d}</span></div>${badge('Active')}</div>`).join('')}</section>`,
    system: renderSystemSettings(),
    upgrade: `<section class="panel"><header class="panel-head"><div><h3>版本升级</h3><p>仅支持管理员手动导入离线升级包，不启用在线热升级</p></div><span class="badge neutral">离线升级</span></header><div class="setting-row"><span class="setting-icon">${icon('package-check')}</span><div class="setting-copy"><strong>当前版本</strong><span>1.0.0 · ARM64 · 构建 20260810</span></div><span class="setting-value">运行正常</span></div><div class="setting-row"><span class="setting-icon">${icon('shield-check')}</span><div class="setting-copy"><strong>升级保护</strong><span>安装前校验签名、版本、架构和磁盘空间，并自动创建业务数据、配置及 Agent 数据快照</span></div>${badge('Active')}</div><div class="setting-row"><span class="setting-icon">${icon('upload')}</span><div class="setting-copy"><strong>导入离线升级包</strong><span>选择本机 .zip 包，校验通过并经管理员确认后安装；完成后设备自动重启</span></div><button class="button primary small" data-action="offline-upgrade">${icon('upload')}选择升级包</button></div><div class="setting-row"><span class="setting-icon">${icon('history')}</span><div class="setting-copy"><strong>最近升级记录</strong><span>暂无升级记录 · 日志和旧版本将保留用于失败回滚</span></div><button class="button small" data-action="upgrade-history">查看记录</button></div></section>`,
  };
  return content[state.settingsTab] || content.model;
}

function renderChannelSettings() {
  const channels = state.openClawChannels || [];
  if (state.openClawChannelsLoading && !channels.length) return `<section class="panel"><div class="empty-state">${icon('loader-circle')}<div><h3>正在读取 OpenClaw 通道目录</h3><p>读取当前固定版本支持的聊天通道和账号状态。</p></div></div></section>`;
  if (state.openClawChannelsError && !channels.length) return `<section class="panel"><header class="panel-head"><div><h3>💬 通道绑定</h3><p>OpenClaw 通道目录未能加载。</p></div><button class="icon-button" data-action="refresh-openclaw-channels" title="重试读取通道" aria-label="重试读取通道">${icon('refresh-cw')}</button></header><div class="model-warning error"><span>${icon('triangle-alert')} 通道读取失败：${escapeHTML(state.openClawChannelsError)}</span></div><div class="empty-state"><p>请重试；如果仍失败，请检查 OpenClaw CLI 和服务状态。</p></div></section>`;
  const query = normalizeSearch(state.channelSearch);
  const visibleChannels = query ? channels.filter(channel=>channelSearchText(channel).includes(query)) : channels;
  const configured = channels.filter(channel=>Boolean(channel.configured) || Number(channel.accountCount||0)>0).length;
  return `<section class="panel"><header class="panel-head"><div><h3>💬 通道绑定</h3><p>通道用于让飞书、企业微信、微信、Telegram、Slack 等入口真实接入 OpenClaw。客户侧只需要配置模型和通道，系统内部能力由应用自动处理。</p></div><button class="icon-button" data-action="refresh-openclaw-channels" title="刷新通道状态" aria-label="刷新通道状态">${icon('refresh-cw')}</button></header><div class="model-warning"><span>${icon('info')} 页面状态区分通道包、插件启用、账号绑定和运行状态；扫码绑定会直接写入 OpenClaw，成功后自动重启网关并复核。</span></div><div class="settings-filter-row"><label class="field-search">${icon('search')}<input id="channelSearch" value="${escapeAttr(state.channelSearch)}" placeholder="搜索飞书、微信、企业微信、Telegram、Slack 等通道"></label><span class="result-count">已配置 ${configured} 个 · 当前显示 ${visibleChannels.length}/${channels.length}</span></div><div class="settings-subsection">${visibleChannels.map(renderChannelSettingRow).join('')||`<div class="empty-state"><p>未找到匹配通道</p></div>`}</div></section>`;
}

function channelSearchText(channel) {
  return normalizeSearch([channel.id,channel.name,channel.description,channel.origin,channel.status].join(' '));
}

function renderChannelSettingRow(channel) {
  const installed = Boolean(channel.installed);
  const accountCount = Number(channel.accountCount||0);
  const status = channelBusinessStatus(channel);
  const bindingButton = installed
    ? `<button class="button primary small" data-action="open-channel-binding" data-channel="${escapeAttr(channel.id)}">${icon('link')}绑定</button>`
    : `<button class="button small" disabled title="需要先安装 OpenClaw 通道包">${icon('package-x')}未安装</button>`;
  return `<div class="setting-row"><span class="setting-icon">${icon('message-circle')}</span><div class="setting-copy"><strong>${escapeHTML(channel.name)} <small class="secondary-text">${escapeHTML(channel.id)}</small></strong><span>${escapeHTML(channel.description||'OpenClaw 聊天通道')} · ${escapeHTML(channel.origin||'installable')} · 已绑定账号 ${accountCount} 个<br><small class="secondary-text">${escapeHTML(status.description)}</small></span></div><span class="badge ${status.className}">${status.label}</span><div class="inline-actions compact-actions"><button class="button small" data-action="channel-status" data-channel="${escapeAttr(channel.id)}">${icon('activity')}状态</button>${bindingButton}</div></div>`;
}

function channelBusinessStatus(channel) {
  if (!channel.installed) return { label: '未安装', className: 'neutral', description: '当前通道包未安装，不能直接绑定。' };
  if (!channel.enabled) return { label: '未启用', className: 'amber', description: '通道包已存在，但 OpenClaw 配置中尚未启用。' };
  if (channel.configured && channel.running) return { label: '已连接', className: 'green', description: '账号已配置，通道正在运行。' };
  if (channel.configured && channel.lastError && channel.lastError !== 'not configured') return { label: '绑定异常', className: 'red', description: `账号已配置，但通道运行异常：${channel.lastError}` };
  if (channel.configured) return { label: '已绑定', className: 'green', description: '账号凭据已写入 OpenClaw，当前未运行。' };
  return { label: '可绑定', className: 'blue', description: '通道已启用但尚未绑定账号。' };
}

function renderSchedulerSettings() {
  return `<section class="panel"><header class="panel-head"><div><h3>定时任务</h3><p>内置任务不能删除，但可以编辑说明、Prompt、Agent、频率和开关；尚未接入真实数据源的任务会保持待确认状态。</p></div><button class="button small" data-action="new-schedule">${icon('plus')}新增任务</button></header>${scheduledJobs.map(job=>`<div class="setting-row"><span class="setting-icon">${icon('timer-reset')}</span><div class="setting-copy"><strong>${escapeHTML(job.name)}</strong><span>${escapeHTML(job.description||'未填写任务说明')} · Agent：${escapeHTML(job.agentId||'未指定')} · ${escapeHTML(job.schedule||'未设置')}${job.error?` · ${escapeHTML(job.error)}`:''}</span></div>${badge(job.enabled?(job.status==='Ready'?'Active':'Review'):'Draft')}<button class="table-icon" data-action="edit-schedule" data-id="${escapeAttr(job.id)}" title="编辑">${icon('pencil')}</button></div>`).join('')||`<div class="empty-state"><p>未读取到定时任务。</p></div>`}</section>`;
}

function renderModelSettings() {
  const data = state.openClawModels;
  const configuredModelList = getConfiguredModelList();
  const catalogModels = getModelCatalog();
  const current = currentDefaultModelKey();
  const configuredKeys = configuredModelKeySet();
  const test = state.modelTestResult;
  const testText = state.modelTestLoading ? '正在写入 OpenClaw 并发起真实模型调用。' : test ? `${test.ok?'正常':'异常'}：${escapeHTML(test.message||'验证完成')} · ${Number(test.durationMs||0).toLocaleString()} ms` : '新建或编辑模型时，可先测试连通性再保存；已保存密钥不会回传明文。';
  const testedCount = configuredModelList.filter(model=>model.lastTestStatus==='passed').length;
  const unconfiguredCount = catalogModels.filter(model=>!configuredKeys.has(model.key)).length;
  return `<section class="panel model-config-panel"><header class="panel-head"><div><h3>🤖 模型设置</h3><p>用于配置多个可用模型。未手动设置默认模型时，系统会使用第一个配置成功的模型作为默认模型。</p></div><button class="icon-button" data-action="refresh-openclaw-models" title="刷新配置状态" aria-label="刷新配置状态">${icon('refresh-cw')}</button></header>
    ${data?.error ? `<div class="model-warning"><span>${icon('triangle-alert')} ${escapeHTML(data.error)}</span><button class="button small" data-action="refresh-openclaw-models">重试</button></div>` : ''}
    <div class="model-config-body">
      <div class="model-warning"><span>${icon('info')} API Key 会按所选模型自动写入 OpenClaw 对应凭据；模型列表只展示已配置记录，新增时只能选择一个未配置模型。</span></div>
      <div class="model-list-shell">
        <div class="model-list-head"><strong>模型列表</strong><span>${configuredModelList.length} 个已保存 · ${testedCount} 个正常 · ${unconfiguredCount} 个待配置</span><button class="button small" data-action="model-draft-new">${icon('plus')}新建模型</button></div>
        <div class="model-status-list model-status-list-wide">${configuredModelList.map(model=>renderConfiguredModelRow(model,current)).join('')||`<div class="empty-state compact-empty"><p>当前还没有已配置模型。</p></div>`}</div>
      </div>
    </div>
  </section>`;
}

function normalizeSearch(value) { return String(value||'').trim().toLowerCase(); }
function modelSearchText(model) { return normalizeSearch([model.key,model.name,model.input,model.source,...(model.tags||[])].join(' ')); }

function modelProvider(modelID) { return String(modelID||'').split('/')[0] || ''; }
function getModelCatalog() { return state.openClawModels?.catalogModels||state.openClawModels?.models||[]; }
function getConfiguredModelList() { return state.openClawModels?.configuredModels||[]; }
function configuredModelKeySet() { return new Set(getConfiguredModelList().map(model=>model.key)); }
function configuredModelKeys() { return getConfiguredModelList().map(model=>model.key).filter(Boolean); }
function currentDefaultModelKey() { return state.openClawModels?.resolvedDefault || state.openClawModels?.defaultModel || ''; }
function configuredModelEntry(modelKey) { return getConfiguredModelList().find(model=>model.key===modelKey) || null; }
function testedConfiguredModelKeys() { return getConfiguredModelList().filter(model=>model.lastTestStatus==='passed').map(model=>model.key); }
function normalizeAgentMessageModel(model='') {
  const selected = String(model || '').trim();
  const defaultModel = currentDefaultModelKey();
  return selected && selected !== defaultModel ? selected : '';
}
function agentSelectedModel(agentID) {
  return normalizeAgentMessageModel(state.agentModelSelections?.[agentID] || '');
}
function providerEndpointSummary(provider) {
  const providerKey = String(provider || '').toLowerCase();
  const defaults = {
    deepseek: '官方入口 api.deepseek.com',
    mistral: '官方入口 api.mistral.ai',
    moonshot: '官方入口 api.moonshot.cn',
    novita: '官方入口 api.novita.ai',
    nvidia: '官方入口 integrate.api.nvidia.com',
    together: '官方入口 api.together.xyz',
    anthropic: '官方入口 api.anthropic.com',
    cohere: '官方入口 api.cohere.com',
    volcengine: '火山方舟国内入口',
    'volcengine-plan': '火山方舟国内入口',
    byteplus: 'BytePlus 国际入口',
    'byteplus-plan': 'BytePlus 国际入口',
  };
  if (providerKey !== 'minimax') return defaults[providerKey] || 'OpenClaw 默认入口';
  const url = state.openClawModels?.providerBaseUrls?.[providerKey] || '';
  if (url.includes('minimaxi.com')) return '国内站 api.minimaxi.com';
  if (url.includes('minimax.io')) return '国际站 api.minimax.io';
  return '自动识别';
}
function modelEndpointModeFor(provider) {
  const providerKey = String(provider || '').toLowerCase();
  const url = state.openClawModels?.providerBaseUrls?.[providerKey] || '';
  if (url.includes('minimaxi.com')) return 'domestic';
  if (url.includes('minimax.io')) return 'international';
  return 'auto';
}
function findModelByKey(key) {
  const models = [...getModelCatalog(), ...getConfiguredModelList()];
  return models.find(model=>model.key===key) || null;
}
function draftModelKey() {
  const catalog = getModelCatalog();
  if (!catalog.length) return state.modelDraftKey || '';
  const configured = configuredModelKeySet();
  const key = state.modelDraftKey;
  if (state.modelDraftMode === 'edit' && key && configured.has(key)) return key;
  if (state.modelDraftMode !== 'edit' && key && catalog.some(model=>model.key===key) && !configured.has(key)) return key;
  const next = catalog.find(model=>!configured.has(model.key)) || catalog[0];
  state.modelDraftMode = 'create';
  state.modelDraftKey = next?.key || '';
  return state.modelDraftKey;
}
function defaultModelForKeys(keys, preferred='') {
  keys = [...new Set((keys||[]).filter(Boolean))];
  if (!keys.length) return '';
  const current = currentDefaultModelKey();
  if (current && keys.includes(current)) return current;
  if (preferred && keys.includes(preferred) && state.modelTestResult?.ok) return preferred;
  const passed = getConfiguredModelList().find(model=>keys.includes(model.key)&&model.lastTestStatus==='passed');
  return passed?.key || keys[0];
}
function modelFamilyName(model) {
  return String(model?.name||model?.key||'').replace(/\s+(M|V|R|K|Opus|Sonnet|Haiku|Flash|Pro|Mini|Nano|Large|Medium|Small|Highspeed|Coding|Thinking|Preview|latest|\\d|\\.)[\w\s.:-]*$/i,'').trim() || String(model?.name||model?.key||'模型');
}
function modelFamilyKey(model) {
  const name = modelFamilyName(model).toLowerCase();
  return name || String(model?.key||'').split('/').slice(-1)[0].toLowerCase();
}
function modelFamilies(models=[]) {
  const groups = new Map();
  models.forEach(model=>{
    const key = modelFamilyKey(model);
    if (!groups.has(key)) groups.set(key,{key,name:modelFamilyName(model),models:[]});
    groups.get(key).models.push(model);
  });
  return [...groups.values()].sort((a,b)=>a.name.localeCompare(b.name,'zh-Hans-CN')).map(group=>({...group,models:group.models.sort((a,b)=>a.name.localeCompare(b.name,'zh-Hans-CN'))}));
}
function modelSeriesLabel(model) {
  const name = String(model?.name||'').trim();
  const key = String(model?.key||'').trim();
  const provider = key.split('/')[0] || '';
  const text = `${name} ${key}`.toLowerCase();
  const known = [
    ['minimax','MiniMax'],['deepseek','DeepSeek'],['claude','Claude'],['kimi','Kimi'],
    ['moonshot','Kimi'],['glm','GLM'],['doubao','Doubao'],['gemini','Gemini'],
    ['mistral','Mistral'],['codestral','Mistral'],['devstral','Mistral'],['magistral','Mistral'],
    ['llama','Llama'],['gpt','GPT'],['codex','Codex'],['nemotron','Nemotron'],
    ['xiaomi','Xiaomi'],['mimo','Xiaomi'],['ark','Ark'],['command','Command'],['seed','Seed'],
  ];
  const match = known.find(([token])=>text.includes(token));
  if (match) return match[1];
  const firstWord = name.split(/\s+/)[0] || provider || '模型';
  return firstWord.replace(/[-_]?v?\d[\w.:/-]*$/i,'').replace(/[-_]+$/,'') || firstWord;
}
function modelSeriesKey(model) {
  return normalizeSearch(modelSeriesLabel(model));
}
function modelVersionLabel(model) {
  const name = String(model?.name||model?.key||'').trim();
  const series = modelSeriesLabel(model);
  const escaped = series.replace(/[.*+?^${}()|[\]\\]/g,'\\$&');
  const value = name.replace(new RegExp(`^${escaped}[\\s:_-]*`,'i'),'').trim();
  return value || name;
}
function modelSeriesGroups(models=[]) {
  const groups = new Map();
  models.forEach(model=>{
    const key = modelSeriesKey(model);
    if (!groups.has(key)) groups.set(key,{key,name:modelSeriesLabel(model),models:[]});
    groups.get(key).models.push(model);
  });
  return [...groups.values()]
    .sort((a,b)=>a.name.localeCompare(b.name,'zh-Hans-CN'))
    .map(group=>({...group,models:group.models.sort((a,b)=>{
      const preferred = {minimax:'minimax',deepseek:'deepseek',claude:'anthropic',mistral:'mistral',kimi:'moonshot',glm:'byteplus',doubao:'volcengine',xiaomi:'xiaomi'}[group.key] || group.key;
      const providerA = modelProvider(a.key);
      const providerB = modelProvider(b.key);
      const directA = providerA===preferred ? 0 : 1;
      const directB = providerB===preferred ? 0 : 1;
      if (directA !== directB) return directA-directB;
      return modelVersionLabel(a).localeCompare(modelVersionLabel(b),'zh-Hans-CN');
    })}));
}
function renderModelSelectOptions(selected='', options={}) {
  const configured = configuredModelKeySet();
  let models = getModelCatalog().filter(model=>!options.onlyUnconfigured || !configured.has(model.key) || model.key===options.includeKey);
  if (selected && !models.some(model=>model.key===selected)) {
    const selectedModel = findModelByKey(selected);
    if (selectedModel) models = [selectedModel, ...models];
  }
  return (modelFamilies(models).map(group=>`<optgroup label="${escapeAttr(group.name)}">${group.models.map(model=>{
    const testLabel = model.lastTestStatus==='passed'?'正常':model.lastTestStatus==='failed'?'异常':'未测试';
    const configuredLabel = configured.has(model.key)?'已配置':'待配置';
    return `<option value="${escapeAttr(model.key)}" ${model.key===selected?'selected':''}>${escapeHTML(model.name)} · ${escapeHTML(model.key)}（${configuredLabel} / ${testLabel}）</option>`;
  }).join('')}</optgroup>`).join('') || '<option value="">暂无模型</option>');
}
function renderModelTree(models, configuredKeys, current) {
  return modelFamilies(models).map(group=>{
    const allChecked = group.models.length > 0 && group.models.every(model=>configuredKeys.has(model.key));
    const someChecked = group.models.some(model=>configuredKeys.has(model.key));
    return `<div class="model-family">
      <label class="model-family-head"><input class="checkbox" type="checkbox" data-model-family="${escapeAttr(group.key)}" ${allChecked?'checked':''} data-indeterminate="${someChecked&&!allChecked?'true':'false'}"><span>${icon(someChecked?'folder-check':'folder')}</span><strong>${escapeHTML(group.name)}</strong><small>${group.models.length} 个版本</small></label>
      <div class="model-variant-list">${group.models.map(model=>renderModelVariant(model,configuredKeys.has(model.key),model.key===current)).join('')}</div>
    </div>`;
  }).join('') || `<div class="empty-state compact-empty"><p>未读取到模型目录。</p></div>`;
}
function renderModelVariant(model, checked, isDefault) {
  const testLabel = model.lastTestStatus==='passed'?'正常':model.lastTestStatus==='failed'?'异常':'未测试';
  return `<label class="model-variant ${isDefault?'is-default':''}">
    <input class="checkbox" type="checkbox" data-model-key="${escapeAttr(model.key)}" ${checked?'checked':''}>
    <span class="model-variant-copy"><strong>${escapeHTML(model.name)}</strong><small>${escapeHTML(model.key)}</small></span>
    <span class="model-variant-meta">${escapeHTML(model.input||'text')} · ${model.contextWindow?Number(model.contextWindow).toLocaleString()+' ctx · ':''}${escapeHTML(testLabel)}</span>
    ${isDefault?'<span class="badge green">默认</span>':''}
  </label>`;
}
function renderConfiguredModelRow(model,current) {
  const statusText = modelStatusText(model.lastTestStatus);
  const active = state.modelDraftMode === 'edit' && state.modelDraftKey === model.key;
  const testing = state.modelTestLoading && state.testingModelKey === model.key;
  const endpoint = providerEndpointSummary(modelProvider(model.key));
  const defaultBadge = model.key===current ? `<span class="badge blue">默认</span>` : '';
  const healthBadge = testing
    ? `<span class="badge amber">检测中</span>`
    : model.lastTestStatus==='passed'
      ? `<span class="badge green">正常</span>`
      : model.lastTestStatus==='failed'
        ? `<span class="badge red">异常</span>`
        : `<span class="badge neutral">未测试</span>`;
  const statusLine = testing
    ? '正在写入 OpenClaw 并验证连通性...'
    : `${escapeHTML(statusText)} · ${escapeHTML(endpoint)}${model.lastTestAt?` · ${escapeHTML(formatLocalizedDateTime(model.lastTestAt))}`:''}${model.lastTestMessage?` · ${escapeHTML(shortModelTestMessage(model.lastTestMessage))}`:''}`;
  return `<div class="setting-row configured-model-row ${active?'is-editing':''} ${testing?'is-testing':''}"><span class="setting-icon">${icon(testing?'loader-circle':'bot')}</span><div class="setting-copy"><strong>${escapeHTML(model.name)} <small class="secondary-text">${escapeHTML(model.key)}</small></strong><span>${statusLine}</span></div>${defaultBadge}${healthBadge}<div class="inline-actions compact-actions"><button class="button small" data-action="test-configured-model" data-model="${escapeAttr(model.key)}" title="使用已保存 OpenClaw 配置测试" ${state.modelTestLoading?'disabled':''}>${icon(testing?'loader-circle':'play')}${testing?'测试中':'测试'}</button><button class="button small" data-action="edit-model-config" data-model="${escapeAttr(model.key)}" ${state.modelTestLoading?'disabled':''}>${icon('pencil')}编辑</button><button class="button small" data-action="select-default-model" data-model="${escapeAttr(model.key)}" ${model.key===current||state.modelTestLoading?'disabled':''}>设为默认</button><button class="button danger small" data-action="delete-model-config" data-model="${escapeAttr(model.key)}" ${state.modelTestLoading?'disabled':''}>${icon('trash-2')}删除</button></div></div>`;
}

function modelStatusText(status) {
  if (status === 'passed') return '正常';
  if (status === 'failed') return '异常';
  return '未测试';
}

function shortModelTestMessage(message) {
  const text = String(message || '').trim();
  if (!text) return '';
  return text.length > 42 ? `${text.slice(0, 42)}...` : text;
}

function formatLocalizedDateTime(value) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  const locale = state.lang === 'en' ? 'en-US' : 'zh-CN';
  return date.toLocaleString(locale, {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  });
}

function formatLocalizedTime(value) {
  if (!value) return '';
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return String(value);
  const locale = state.lang === 'en' ? 'en-US' : 'zh-CN';
  return date.toLocaleTimeString(locale, {
    hour: '2-digit',
    minute: '2-digit',
  });
}

function applyModelTestResultToState(result) {
  const modelKey = result?.model;
  if (!modelKey || !state.openClawModels) return;
  const status = result.ok ? 'passed' : 'failed';
  const message = result.message || (result.ok ? 'API Key 已写入 OpenClaw，并完成真实模型调用验证' : '模型验证未通过');
  const testedAt = result.testedAt || new Date().toISOString();
  for (const listName of ['models','catalogModels','configuredModels']) {
    const list = state.openClawModels[listName];
    if (!Array.isArray(list)) continue;
    list.forEach(model=>{
      if (model.key === modelKey) {
        model.lastTestStatus = status;
        model.lastTestMessage = message;
        model.lastTestAt = testedAt;
      }
    });
  }
}

function draftModelCatalog() {
  const configured = configuredModelKeySet();
  const editing = state.modelDraftMode === 'edit';
  let models = getModelCatalog().filter(model=>!configured.has(model.key) || (editing && model.key===state.modelDraftOriginalKey));
  if (editing && state.modelDraftKey && !models.some(model=>model.key===state.modelDraftKey)) {
    const model = findModelByKey(state.modelDraftKey);
    if (model) models = [model, ...models];
  }
  return models;
}

function draftModelFamilyKey() {
  const groups = modelSeriesGroups(draftModelCatalog());
  if (state.modelDraftFamilyKey && groups.some(group=>group.key===state.modelDraftFamilyKey)) return state.modelDraftFamilyKey;
  const selected = findModelByKey(state.modelDraftKey) || draftModelCatalog()[0];
  return selected ? modelSeriesKey(selected) : groups[0]?.key || '';
}

function renderModelFamilySelectOptions(selectedFamily='') {
  return modelSeriesGroups(draftModelCatalog()).map(group=>`<option value="${escapeAttr(group.key)}" ${group.key===selectedFamily?'selected':''}>${escapeHTML(group.name)}</option>`).join('') || '<option value="">暂无可选模型</option>';
}

function renderModelVersionSelectOptions(familyKey='', selected='') {
  const groups = modelSeriesGroups(draftModelCatalog());
  const group = groups.find(item=>item.key===familyKey) || groups[0];
  if (!group) return '<option value="">暂无可选版本</option>';
  const modelKey = selected && group.models.some(model=>model.key===selected) ? selected : group.models[0]?.key || '';
  return group.models.map(model=>`<option value="${escapeAttr(model.key)}" ${model.key===modelKey?'selected':''}>${escapeHTML(modelVersionLabel(model))} · ${escapeHTML(model.key)}</option>`).join('');
}

function renderModelConfigurationForm() {
  const editing = state.modelDraftMode === 'edit';
  const selectedKey = draftModelKey();
  const familyKey = draftModelFamilyKey();
  const selectedProvider = modelProvider(selectedKey);
  const endpointMode = modelEndpointModeFor(selectedProvider);
  const test = state.modelTestResult;
  const testText = state.modelTestLoading ? '正在写入 OpenClaw 并发起真实模型调用。' : test ? `${test.ok?'正常':'异常'}：${escapeHTML(test.message||'验证完成')} · ${Number(test.durationMs||0).toLocaleString()} ms` : '建议先测试连通性，再保存配置。列表中的测试会使用已保存的 API Key 真实验证。';
  return `<div class="form-grid model-modal-form">
    <div class="form-field"><label for="modelFamilySelect">模型</label><select class="select" id="modelFamilySelect">${renderModelFamilySelectOptions(familyKey)}</select><small>先选模型系列。</small></div>
    <div class="form-field"><label for="modelVersionSelect">模型小版本</label><select class="select" id="modelVersionSelect">${renderModelVersionSelectOptions(familyKey, selectedKey)}</select><small>每次只选一个版本。</small></div>
    <div class="form-field full"><label for="modelEndpointMode">API Key 接入区域</label><select class="select" id="modelEndpointMode"><option value="auto" ${endpointMode==='auto'?'selected':''}>自动识别</option><option value="domestic" ${endpointMode==='domestic'?'selected':''}>国内站 · api.minimaxi.com</option><option value="international" ${endpointMode==='international'?'selected':''}>国际站 · api.minimax.io</option></select><small>MiniMax 的 sk-cp 开头 Key 建议使用国内站；非 MiniMax 模型保持自动即可。保存后会同步到 OpenClaw 当前配置。</small></div>
    <div class="form-field full"><label for="modelAPIKey">API Key</label><div class="credential-field"><input class="input" id="modelAPIKey" type="password" autocomplete="new-password" placeholder="${editing?'留空则使用已保存 API Key 测试或保存':'请输入该模型对应的 API Key'}"><button class="icon-button" type="button" data-action="toggle-api-key-input" title="显示 API Key" aria-label="显示或隐藏 API Key">${icon('eye')}</button></div><small>可临时显示本次输入；已保存密钥由 OpenClaw 存储且不会回传明文。</small></div>
    <div class="model-test-summary full ${test?.ok?'success':test?'warning':''}">${icon(test?.ok?'badge-check':test?'triangle-alert':'info')}${testText}<small>测试连通性会先把 API Key 和接入区域写入 OpenClaw，再向所选模型发起真实调用；只有真实调用成功才会标记为正常。</small></div>
  </div>`;
}

function openModelConfigurationForm(mode='create', model='') {
  state.modelDraftMode = mode === 'edit' ? 'edit' : 'create';
  state.modelDraftOriginalKey = state.modelDraftMode === 'edit' ? model : '';
  state.modelDraftKey = model || '';
  const selected = findModelByKey(state.modelDraftKey) || draftModelCatalog()[0];
  state.modelDraftFamilyKey = selected ? modelSeriesKey(selected) : '';
  state.modelTestResult = null;
  const title = state.modelDraftMode === 'edit' ? '编辑模型' : '新建模型';
  openModal({
    title,
    eyebrow:'设置 / 模型',
    wide:true,
    body:renderModelConfigurationForm(),
    footer:`<button class="button" data-action="close-modal">取消</button><button class="button primary" data-action="test-model-config" ${state.modelTestLoading?'disabled':''}>${icon(state.modelTestLoading?'loader-circle':'play')}${state.modelTestLoading?'测试中':'测试连通性'}</button><button class="button" data-action="save-model-config" ${state.modelTestLoading?'disabled':''}>${icon('save')}保存</button>`
  });
}

function refreshModelConfigurationModal() {
  if (document.getElementById('modalBackdrop')?.hidden) return;
  const body = document.getElementById('modalBody');
  const footer = document.getElementById('modalFooter');
  if (!body || !footer || !document.getElementById('modelVersionSelect')) return;
  body.innerHTML = renderModelConfigurationForm();
  footer.innerHTML = `<button class="button" data-action="close-modal">取消</button><button class="button primary" data-action="test-model-config" ${state.modelTestLoading?'disabled':''}>${icon(state.modelTestLoading?'loader-circle':'play')}${state.modelTestLoading?'测试中':'测试连通性'}</button><button class="button" data-action="save-model-config" ${state.modelTestLoading?'disabled':''}>${icon('save')}保存</button>`;
  applyIcons();
}

function syncModelVersionSelect() {
  const family = document.getElementById('modelFamilySelect')?.value || '';
  state.modelDraftFamilyKey = family;
  const version = document.getElementById('modelVersionSelect');
  if (!version) return;
  version.innerHTML = renderModelVersionSelectOptions(family, '');
  state.modelDraftKey = version.value || '';
  const endpoint = document.getElementById('modelEndpointMode');
  if (endpoint) endpoint.value = modelEndpointModeFor(modelProvider(state.modelDraftKey));
  state.modelTestResult = null;
  applyIcons();
}
function syncModelFamilyIndeterminate() {
  document.querySelectorAll('[data-model-family]').forEach(input=>{
    const familyBlock = input.closest('.model-family');
    const items = familyBlock ? [...familyBlock.querySelectorAll('[data-model-key]')] : [];
    const checked = items.filter(item=>item.checked).length;
    input.checked = items.length > 0 && checked === items.length;
    input.indeterminate = checked > 0 && checked < items.length;
  });
}
function selectedModelKeys() {
  return [...document.querySelectorAll('[data-model-key]:checked')].map(input=>input.dataset.modelKey).filter(Boolean);
}
function toggleModelFamily(input) {
  const familyBlock = input.closest('.model-family');
  if (!familyBlock) return;
  familyBlock.querySelectorAll('[data-model-key]').forEach(item=>{item.checked=input.checked;});
  input.indeterminate = false;
  state.modelTestResult=null;
  syncModelFamilyIndeterminate();
}
function updateModelSelection() {
  const modelSelect = document.getElementById('modelVersionSelect');
  if (modelSelect) state.modelDraftKey = modelSelect.value;
  const selected = findModelByKey(state.modelDraftKey);
  if (selected) state.modelDraftFamilyKey = modelSeriesKey(selected);
  const endpoint = document.getElementById('modelEndpointMode');
  if (endpoint) endpoint.value = modelEndpointModeFor(modelProvider(state.modelDraftKey));
  state.modelTestResult=null;
}

function updateModelFamilySelection() {
  syncModelVersionSelect();
}

async function selectDefaultModel(model) {
  try {
    const selectedModels = new Set(configuredModelKeys());
    if (!selectedModels.has(model)) selectedModels.add(model);
    await apiFetch('/api/v1/settings/model', { method: 'PATCH', body: JSON.stringify({ model: '', defaultModel: model, selectedModels:[...selectedModels], apiKey: '' }) });
    await loadOpenClawModels(true);
    toast('默认模型已更新', model, 'success');
  } catch (error) {
    toast('默认模型设置失败', error.message, 'warning');
  }
}

function openChannelBinding(channel) {
  const info = (state.openClawChannels||[]).find(item=>item.id===channel);
  const installed = info ? Boolean(info.installed) : true;
  stopChannelQRPoll();
  state.channelQR = null;
  openModal({
    title:`${info?.name||channel} 通道绑定`,
    eyebrow:'OpenClaw / 通道',
    body:`<div class="form-grid">${!installed?`<div class="model-warning full"><span>${icon('triangle-alert')} 当前通道包尚未安装，不能开始扫码绑定。</span></div>`:''}<div class="form-field"><label for="channelAccount">账号</label><input class="input" id="channelAccount" value="default" placeholder="扫码登录使用的账号标识"><small>默认账号会写入 channels.feishu；也支持填写 OpenClaw 账号 ID。</small></div><div class="form-field"><label for="channelDomain">接入区域</label><select class="select" id="channelDomain"><option value="feishu">飞书（中国大陆）</option><option value="lark">Lark（国际版）</option></select><small>扫码时会按所选区域请求 OpenClaw 官方授权服务。</small></div><div class="model-warning full"><span>${icon('info')} 点击“开始扫码”后生成二维码。扫码授权成功后，应用会自动写入 OpenClaw、重启网关并复核状态。</span></div><div class="channel-qr-panel full" id="channelQRPanel"><div class="channel-qr-empty">${icon('qrcode')}<strong>尚未生成二维码</strong><span>二维码只在本次绑定会话内有效。</span></div></div><pre class="channel-status-box full" id="channelStatusBox">尚未读取状态</pre></div>`,
    footer:`<button class="button" data-action="channel-status" data-channel="${escapeAttr(channel)}">${icon('activity')}查看状态</button><button class="button primary" ${installed?'':'disabled'} data-action="save-channel-binding" data-channel="${escapeAttr(channel)}">${icon('qrcode')}开始扫码</button>`
  });
  void refreshChannelStatus(channel);
}

async function refreshChannelStatus(channel) {
  const target=document.getElementById('channelStatusBox');
  if(target)target.textContent='正在读取 OpenClaw 通道状态...';
  try {
    const status=await apiFetch(`/api/v1/openclaw/channels/${encodeURIComponent(channel)}/status`);
    if(target)target.textContent=JSON.stringify(status,null,2);
  } catch(error) {
    if(target)target.textContent=`状态读取失败：${error.message}`;
    toast('通道状态读取失败',error.message,'warning');
  }
}

function stopChannelQRPoll() {
  if (state.channelQRPollTimer) {
    clearTimeout(state.channelQRPollTimer);
    state.channelQRPollTimer = null;
  }
}

function renderChannelQRPanel(data) {
  const target=document.getElementById('channelQRPanel');
  if (!target) return;
  if (!data) {
    target.innerHTML=`<div class="channel-qr-empty">${icon('qrcode')}<strong>尚未生成二维码</strong><span>二维码只在本次绑定会话内有效。</span></div>`;
    applyIcons();
    return;
  }
  const terminal = data.status === 'pending'
    ? `<div class="channel-qr-progress"><div class="progress"><span style="width:62%"></span></div><span>${escapeHTML(data.message||'等待扫码确认')} · 自动检查中</span></div>`
    : `<div class="channel-qr-result ${data.status==='success'?'success':'error'}">${icon(data.status==='success'?'check-circle-2':'triangle-alert')}<span>${escapeHTML(data.message||'绑定流程已结束')}</span></div>`;
  target.innerHTML=`<div class="channel-qr-content">${data.qrDataUrl?`<img class="channel-qr-image" src="${escapeAttr(data.qrDataUrl)}" alt="飞书扫码二维码">`:`<div class="channel-qr-fallback">${icon('link')}<span>当前环境未生成图片二维码，请打开授权地址扫码：</span><a href="${escapeAttr(data.qrUrl||'')}" target="_blank" rel="noopener">${escapeHTML(data.qrUrl||'')}</a></div>`}<div class="channel-qr-meta"><strong>请使用飞书移动端扫描</strong><span>账号：${escapeHTML(data.account||'default')} · 区域：${escapeHTML(data.domain||'feishu')}</span><small>有效期至 ${escapeHTML(data.expiresAt||'')}</small>${terminal}</div></div>`;
  applyIcons();
}

async function pollChannelQRCode(channel, sessionID) {
  try {
    const data=await apiFetch(`/api/v1/openclaw/channels/${encodeURIComponent(channel)}/qr/${encodeURIComponent(sessionID)}/status`);
    state.channelQR=data;
    renderChannelQRPanel(data);
    if (data.status==='pending') {
      state.channelQRPollTimer=setTimeout(()=>void pollChannelQRCode(channel,sessionID),Math.max(Number(data.interval||5),3)*1000);
      return;
    }
    if (data.status==='success') {
      await loadOpenClawChannels(true);
      await refreshChannelStatus(channel);
      toast('通道绑定成功',`${channel} 已写入 OpenClaw 并完成网关重启。`,'success');
    } else {
      toast('通道绑定未完成',data.message||'请重新生成二维码。','warning');
    }
  } catch(error) {
    const target=document.getElementById('channelQRPanel');
    if(target)target.innerHTML=`<div class="channel-qr-result error">${icon('triangle-alert')}<span>轮询绑定状态失败：${escapeHTML(error.message)}</span></div>`;
    applyIcons();
    toast('通道绑定状态读取失败',error.message,'warning');
  }
}

async function saveChannelBinding(channel) {
  stopChannelQRPoll();
  try {
    const data=await apiFetch(`/api/v1/openclaw/channels/${encodeURIComponent(channel)}/qr/start`,{method:'POST',body:JSON.stringify({
      account:formText('channelAccount')||'default',
      domain:document.getElementById('channelDomain')?.value||'feishu',
    })});
    state.channelQR=data;
    renderChannelQRPanel(data);
    toast('二维码已生成','请使用飞书移动端扫描并确认授权。');
    state.channelQRPollTimer=setTimeout(()=>void pollChannelQRCode(channel,data.sessionId),Math.max(Number(data.interval||5),3)*1000);
  } catch(error) {
    toast('通道绑定失败',error.message,'warning');
  }
}

function renderSystemSettings() {
  const status = state.openClawStatus;
  const health = state.systemHealth;
  const managedCount = state.openClawAgents?.filter(agent => !agent.isDefault).length;
  const openClawValue = state.openClawStatusLoading ? '正在读取服务状态' : status?.available ? `${status.version} · ${status.serviceStatus} · RPC ${status.rpcOK?'正常':'异常'}` : status?.error || '暂未读取';
  const agentValue = state.openClawAgentsLoading ? '正在读取 Agent' : managedCount === undefined ? '暂未读取' : `${managedCount} 个 STA-100 Agent 已注册`;
  const rows = [
    ['运行环境',health?.runtime?`${health.runtime.os} / ${health.runtime.arch} / ${health.runtime.go}`:'正在读取','cpu',Boolean(health?.runtime)],
    ['服务健康',health?.status||health?.error||'正在读取','monitor',health?.status==='ok'],
    ['OpenClaw',openClawValue,'bot',Boolean(status?.rpcOK)],
    ['Agent 编排',agentValue,'blocks',managedCount===24],
    ['本地数据库',health?.database?`SQLite · ${health.database.ok?'正常':'异常'} · ${formatBytes(health.database.bytes)}`:'正在读取','database',Boolean(health?.database?.ok)],
    ['存储空间',health?.storage?`私有文件 ${formatBytes(health.storage.privateFileBytes)} / 可用 ${formatBytes(health.storage.availableBytes)}`:'正在读取','hard-drive',Boolean(health?.storage?.ok)],
    ['知识索引',health?.index?`${health.index.indexed} / ${health.index.files} 个文件已索引 · ${health.index.status}`:'正在读取','database-zap',health?.index?.status==='ready'],
  ];
  return `<section class="panel"><header class="panel-head"><div><h3>⚙️ 系统信息</h3><p>数值来自 Go 健康检查、SQLite 和 OpenClaw，不使用预置容量</p></div><button class="button small" data-action="refresh-openclaw-system">${icon('refresh-cw')}刷新</button></header>${rows.map(([name,value,iconName,ok])=>`<div class="setting-row"><span class="setting-icon">${icon(iconName)}</span><div class="setting-copy"><strong>${name}</strong><span>${escapeHTML(value)}</span></div>${badge(ok?'Active':'Review')}</div>`).join('')}</section>`;
}

function wirePageSpecific() {
  const sub = document.getElementById('subscriptionToggle');
  if (sub) sub.addEventListener('change', async e => { const previous=state.subscription; state.subscription=e.target.checked; try { await apiFetch('/api/v1/overview/subscription',{method:'PATCH',body:JSON.stringify({enabled:state.subscription})}); toast('订阅设置已更新',state.subscription?`系统将每 ${state.newsFrequency} 更新一次推荐。`:'自动更新已暂停。'); } catch(error) { state.subscription=previous; toast('订阅设置失败',error.message,'warning'); } renderPage(); });
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
  const modelFamilySelect = document.getElementById('modelFamilySelect');
  if (modelFamilySelect) modelFamilySelect.addEventListener('change', updateModelFamilySelection);
  const modelVersionSelect = document.getElementById('modelVersionSelect');
  if (modelVersionSelect) modelVersionSelect.addEventListener('change', updateModelSelection);
  const channelSearch = document.getElementById('channelSearch');
  if (channelSearch) channelSearch.addEventListener('input', e => { state.channelSearch=e.target.value; renderPage(); requestAnimationFrame(()=>{ const input=document.getElementById('channelSearch'); if(input){input.focus();input.setSelectionRange(input.value.length,input.value.length);} }); });
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
  applyCustomerColumnVisibility();
}

const customerColumnDefinitions = [
  ['customer','客户'],['type','客户类型'],['country','国家'],['contact','联系人'],['orders','订单数'],['total','累计金额'],['rating','评级'],['updated','更新时间'],
];

function openCustomerColumnSettings() {
  openModal({title:'客户列表字段',eyebrow:'显示设置',body:`<div class="filter-row">${customerColumnDefinitions.map(([key,label])=>`<label class="filter-chip ${state.customerVisibleColumns.has(key)?'active':''}"><input type="checkbox" data-customer-column="${key}" ${state.customerVisibleColumns.has(key)?'checked':''}>${label}</label>`).join('')}</div>`,footer:formFooter('应用','save-customer-columns')});
}

function saveCustomerColumns() {
  const selected=[...document.querySelectorAll('[data-customer-column]:checked')].map(input=>input.dataset.customerColumn);
  if(!selected.includes('customer'))selected.unshift('customer');
  state.customerVisibleColumns=new Set(selected);closeModal();renderPage();
}

function applyCustomerColumnVisibility() {
  const table=document.getElementById('customerTable');if(!table)return;
  const indexes={customer:2,type:3,country:4,contact:5,orders:6,total:7,rating:8,updated:9};
  Object.entries(indexes).forEach(([key,index])=>table.querySelectorAll(`tr > *:nth-child(${index})`).forEach(cell=>{cell.hidden=!state.customerVisibleColumns.has(key);}));
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
function closeModal() { stopChannelQRPoll(); document.getElementById('modalBackdrop').hidden = true; syncOverlayScroll(); }
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
function formatBytes(value) {
  const bytes=Number(value||0);
  if(bytes<1024)return `${bytes} B`;
  if(bytes<1024*1024)return `${(bytes/1024).toFixed(1)} KB`;
  if(bytes<1024*1024*1024)return `${(bytes/1024/1024).toFixed(1)} MB`;
  return `${(bytes/1024/1024/1024).toFixed(1)} GB`;
}
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

function formatAttachmentSize(size=0) {
  const value = Number(size || 0);
  if (!value) return '';
  if (value < 1024) return `${value} B`;
  if (value < 1024 * 1024) return `${(value/1024).toFixed(1)} KB`;
  return `${(value/1024/1024).toFixed(1)} MB`;
}

function isImageAttachment(attachment={}) {
  return String(attachment.mime || '').startsWith('image/') || /\.(png|jpe?g|webp|gif|bmp)$/i.test(String(attachment.name || ''));
}

function renderMessageAttachments(message={}) {
  const attachments = Array.isArray(message.attachments) ? message.attachments : [];
  if (!attachments.length) return '';
  const status = message.attachmentStatus || '已随消息提交';
  return `<div class="message-attachments">${attachments.map(attachment=>{
    const meta = [attachment.mime, formatAttachmentSize(attachment.size)].filter(Boolean).join(' · ');
    if (isImageAttachment(attachment) && attachment.previewUrl) {
      return `<figure class="message-attachment image-attachment"><img src="${escapeAttr(attachment.previewUrl)}" alt="${escapeAttr(attachment.name || '图片附件')}"><figcaption><strong>${escapeHTML(attachment.name || '图片附件')}</strong>${meta?`<small>${escapeHTML(meta)}</small>`:''}</figcaption></figure>`;
    }
    return `<div class="message-attachment file-attachment">${icon('paperclip')}<span><strong>${escapeHTML(attachment.name || '附件')}</strong>${meta?`<small>${escapeHTML(meta)}</small>`:''}</span></div>`;
  }).join('')}<div class="attachment-status">${icon(message.attachmentStatus==='发送失败'?'triangle-alert':'check-circle-2')}${escapeHTML(status)}</div></div>`;
}

function renderComposerAttachment(attachment, index, attachmentIndex) {
  const thumb = isImageAttachment(attachment) && attachment.previewUrl ? `<img src="${escapeAttr(attachment.previewUrl)}" alt="${escapeAttr(attachment.name || '图片附件')}">` : icon('paperclip');
  return `<span class="composer-attachment-item">${thumb}<strong>${escapeHTML(attachment.name)}</strong><small>${escapeHTML(attachment.status || '待发送')}</small><button class="table-icon" data-action="remove-chat-attachment" data-agent="${index}" data-index="${attachmentIndex}" title="移除附件">${icon('x')}</button></span>`;
}

function normalizePipelineStage(stage={}) {
  const key = String(stage.stage || '').trim();
  const label = {
    received: '接收消息',
    'local-retrieval': '本地检索',
    attachments: '附件处理',
    'knowledge-agent': '知识整理',
    'domain-agents': '业务 Agent',
    'coordinator-agent': '统一汇总',
  }[key] || key || '处理阶段';
  const status = String(stage.status || '').toLowerCase();
  return { key, label, status, detail: stage.detail || '' };
}

function renderMessagePipeline(message={}) {
  const pipeline = Array.isArray(message.pipeline) ? message.pipeline.map(normalizePipelineStage) : [];
  if (!pipeline.length) return '';
  return `<div class="message-pipeline">${pipeline.map(stage=>`<span class="pipeline-chip ${escapeAttr(stage.status || 'pending')}" title="${escapeAttr(stage.detail || '')}">${icon(stage.status==='failed'?'circle-x':stage.status==='partial'?'triangle-alert':'check-circle-2')}${escapeHTML(stage.label)}</span>`).join('')}</div>`;
}

function renderAgentChatProgress(agentID) {
  const progress = state.agentChatProgress[agentID];
  if (!progress) return '';
  const activeIndex = Math.max(0, Number(progress.index || 0));
  const failed = progress.status === 'failed';
  const done = progress.status === 'done';
  const steps = progress.steps || chatProgressSteps.map(([key,label])=>({key,label}));
  const elapsedSeconds = Number.isFinite(Number(progress.elapsedSeconds))
    ? Number(progress.elapsedSeconds)
    : progress.startedAt
      ? Math.max(0, Math.floor((Date.now() - progress.startedAt) / 1000))
      : 0;
  const elapsedLabel = elapsedSeconds < 60
    ? `${elapsedSeconds} 秒`
    : `${Math.floor(elapsedSeconds / 60)} 分 ${elapsedSeconds % 60} 秒`;
  const timeState = failed ? 'failed' : done ? 'done' : 'running';
  return `<div class="chat-progress ${failed?'failed':done?'done':''}"><div class="progress-line">${steps.map((step,index)=>{
    const stateClass = failed && index===activeIndex ? 'failed' : index < activeIndex || done ? 'done' : index === activeIndex ? 'active' : 'pending';
    return `<span class="${stateClass}"><i></i>${escapeHTML(step.label)}</span>`;
  }).join('')}</div><div class="progress-time"><span>生成进度</span><div class="progress-time-track"><i class="${timeState}"></i></div><strong>${elapsedLabel}</strong></div><div class="progress-detail"><span class="progress-detail-copy">${icon(failed?'circle-x':done?'check-circle-2':'loader-circle')}${escapeHTML(progress.detail || (done ? 'OpenClaw 已返回结果' : 'OpenClaw 正在处理'))}</span><span class="progress-elapsed">${icon('clock-3')}已用时 ${elapsedLabel}</span></div></div>`;
}

function renderAgentMessage(message, index) {
  const modelMeta = message.model ? ` · ${message.provider ? `${escapeHTML(message.provider)}/` : ''}${escapeHTML(message.model)}` : '';
  return `<div class="message-row ${message.role === 'user' ? 'user' : ''}"><div class="message ${message.role === 'user' ? 'user' : message.error ? 'error' : ''}">${escapeHTML(message.text).replace(/\n/g,'<br>')}${renderMessageAttachments(message)}${message.error && message.retry ? `<button class="link-button message-retry" data-action="retry-chat" data-agent="${index}" data-message="${escapeAttr(message.retry || '')}">重新发送</button>` : ''}<time>${escapeHTML(message.time || '')}${modelMeta}</time></div></div>`;
}

function renderAgentChatBody(index) {
  const a = agents[index];
  const agentID = agentIDs[index];
  const messages = state.agentChats[agentID] || [];
  const historyLoading = Boolean(state.agentChatHistoryLoading[agentID]);
  const testedModels = (state.openClawModels?.configuredModels||[]).filter(model=>model.lastTestStatus==='passed');
  const selectedModel = agentSelectedModel(agentID);
  const defaultModel = currentDefaultModelKey();
  const modelOptions = testedModels.map(model=>`<option value="${escapeAttr(model.key)}" ${selectedModel===model.key?'selected':''}>${escapeHTML(model.name)} · ${escapeHTML(model.key)}</option>`).join('');
  const attachments = state.chatAttachments.map((attachment,attachmentIndex)=>renderComposerAttachment(attachment,index,attachmentIndex)).join('');
  const modelHint = state.openClawModelsLoading ? '<div class="chat-model-loading"><span class="secondary-text">正在读取 OpenClaw 已配置模型...</span></div>' : !testedModels.length ? '<div class="chat-model-loading"><span class="secondary-text">暂无可选模型，请先到设置中配置并测试模型。</span></div>' : '';
  return `<div class="chat-layout"><aside class="chat-side"><div class="chat-side-head"><h3>统一智能处理</h3><span>系统编排</span></div><div class="chat-source policy-fixed"><span class="agent-icon">${icon('workflow')}</span><span><strong>本地证据 → 协调器 → 领域 Agent</strong><small>页面不区分本地或联网结果，系统按规则自动整合。</small></span></div><button class="button ghost small source-settings-button" data-action="agent-allowlist" data-agent="${index}">${icon('shield-check')}联网白名单</button><p class="chat-source-note">白名单属于后台安全配置，不是本次聊天的来源选择。冲突数据会全部保留并标记。</p></aside><section class="chat-main"><header class="chat-head"><div><strong>${agentEmojis[index]} ${escapeHTML(a[0])}</strong><span>${escapeHTML(agentID)}</span></div><span>${badge(state.modelConfigured?'Active':'Review')}</span></header><div class="chat-quick-prompts">${a[5].map(v=>`<button data-action="chat-quick-prompt" data-agent="${index}" data-prompt="${escapeAttr(v)}">${escapeHTML(v)}</button>`).join('')}</div><div class="chat-messages" id="chatMessages"><div class="message-row"><div class="message"><strong>${escapeHTML(a[0])}</strong><br>已连接 STA-100 统一任务协调器，消息将按页面规则分发给专业智能体。<time>当前会话</time></div></div>${historyLoading ? `<div class="chat-history-loading">${icon('loader-circle')}正在从 OpenClaw 读取历史消息...</div>` : ''}${messages.map(m=>renderAgentMessage(m,index)).join('')}</div><div class="chat-status" id="chatStatus" aria-live="polite">${renderAgentChatProgress(agentID)}</div><div class="chat-composer"><div id="chatAttachmentList" class="chat-attachments composer-attachments">${attachments}</div><div class="chat-compose-tools"><div class="chat-attachment-tools"><button class="icon-button" data-action="choose-chat-image" title="上传图片" aria-label="上传图片">${icon('image-up')}</button><button class="icon-button" data-action="choose-chat-file" title="上传文件" aria-label="上传文件">${icon('file-up')}</button></div><span class="chat-compose-spacer"></span><label class="chat-model-picker"><span>模型</span><select class="select" id="chatModelSelect"><option value="">使用默认模型${defaultModel?`（${escapeHTML(defaultModel)}）`:''}</option>${modelOptions}</select></label><input id="chatImageInput" type="file" accept="image/*" hidden><input id="chatFileInput" type="file" hidden></div>${modelHint}<div class="chat-input-row"><textarea class="textarea" id="chatInput" maxlength="32768" rows="2" placeholder="向 ${escapeAttr(a[0])} 发送消息，Enter 发送，Shift+Enter 换行"></textarea><button class="icon-button chat-send" data-action="send-chat" data-agent="${index}" title="发送消息" aria-label="发送消息">${icon('send')}</button></div></div></section></div>`;
}

function wireChatInput(index) {
  const input = document.getElementById('chatInput');
  input?.addEventListener('keydown', event => {
    if (event.key === 'Enter' && !event.shiftKey && !event.isComposing) {
      event.preventDefault();
      void sendAgentMessage(index);
    }
  });
  document.getElementById('chatModelSelect')?.addEventListener('change', event=>{
    const agentID = agentIDs[index];
    const selected = normalizeAgentMessageModel(event.target.value);
    if (selected) state.agentModelSelections[agentID]=selected;
    else delete state.agentModelSelections[agentID];
    void savePreferences();
  });
  document.getElementById('chatImageInput')?.addEventListener('change', event=>void uploadChatAttachment(event.target.files?.[0], index));
  document.getElementById('chatFileInput')?.addEventListener('change', event=>void uploadChatAttachment(event.target.files?.[0], index));
  const box = document.getElementById('chatMessages');
  box?.addEventListener('scroll', () => {
    const distance = box.scrollHeight - box.scrollTop - box.clientHeight;
    state.agentChatAtBottom[agentIDs[index]] = distance < 40;
  }, {passive:true});
  input?.focus();
}

async function uploadChatAttachment(file,index) {
  if(!file)return;
  if(file.size>25*1024*1024){toast('附件过大','单个附件不能超过 25 MB。','warning');return;}
  const form=new FormData();
  form.append('file',file);
  try {
    const data=await apiFetch('/api/v1/assistant/attachments',{method:'POST',body:form});
    const attachment = data.attachment || {};
    attachment.localId = `att-${Date.now()}-${Math.random().toString(16).slice(2)}`;
    attachment.previewUrl = file.type.startsWith('image/') ? URL.createObjectURL(file) : '';
    attachment.status = '已上传到 STA-100，待随消息发送';
    state.chatAttachments.push(attachment);
    refreshAgentChat(index);
  } catch(error) { toast('附件上传失败',error.message,'warning'); }
}

async function showAgentChat(index,prompt='') {
  const a=agents[index];
  const agentID = agentIDs[index];
  state.agentChatAtBottom[agentID] = true;
  openModal({title:a[0],eyebrow:`OpenClaw Agent / ${agentIDs[index]}`,wide:true,body:renderAgentChatBody(index)});
  wireChatInput(index);
  requestAnimationFrame(()=>scrollAgentChatToBottom(agentID));
  void loadAgentChatHistory(index);
  if (!state.openClawModels && !state.openClawModelsLoading) {
    void loadOpenClawModels(true).then(() => {
      if (document.getElementById('modalBody')) refreshAgentChat(index);
    });
  } else if (!state.openClawModelsLoading && state.page === 'agents') {
    void loadOpenClawModels().then(() => {
      if (document.getElementById('modalBody')) refreshAgentChat(index);
    });
  }
  if (prompt) setTimeout(()=>void sendAgentMessage(index,prompt),0);
}

function refreshAgentChat(index, options={}) {
  const agentID = agentIDs[index];
  const previous = document.getElementById('chatMessages');
  const previousTop = previous?.scrollTop || 0;
  const wasAtBottom = previous ? (previous.scrollHeight - previous.scrollTop - previous.clientHeight < 40) : true;
  document.getElementById('modalBody').innerHTML = renderAgentChatBody(index);
  applyIcons();
  wireChatInput(index);
  const box = document.getElementById('chatMessages');
  if (!box) return;
  const shouldStick = Boolean(options.forceBottom || wasAtBottom || state.agentChatAtBottom[agentID] !== false);
  if (shouldStick) {
    box.scrollTop = box.scrollHeight;
    state.agentChatAtBottom[agentID] = true;
    requestAnimationFrame(()=>scrollAgentChatToBottom(agentID));
  } else {
    box.scrollTop = previousTop;
    state.agentChatAtBottom[agentID] = false;
  }
}

function scrollAgentChatToBottom(agentID) {
  const box = document.getElementById('chatMessages');
  if (!box) return;
  box.scrollTop = box.scrollHeight;
  state.agentChatAtBottom[agentID] = true;
}

function historyMessageTime(value) {
  return formatLocalizedDateTime(value);
}

function normalizeAgentHistoryMessage(message={}) {
  const role = message.role === 'user' ? 'user' : 'agent';
  const error = Boolean(message.error);
  return {
    role,
    text: String(message.text || '').trim() || (error ? `OpenClaw 调用失败：${message.error}` : 'OpenClaw 未返回文本内容。'),
    error,
    model: message.model || '',
    provider: message.provider || '',
    time: historyMessageTime(message.createdAt),
  };
}

async function loadAgentChatHistory(index) {
  const agentID = agentIDs[index];
  if (state.agentChatHistoryLoaded[agentID] || state.agentChatHistoryLoading[agentID]) return;
  state.agentChatHistoryLoading[agentID] = true;
  refreshAgentChat(index, {forceBottom:true});
  try {
    const data = await apiFetch(`/api/v1/assistant/history?agentId=${encodeURIComponent(agentID)}&sessionKey=${encodeURIComponent(`sta100-${agentID}`)}&limit=200`);
    const persisted = (data.messages || []).map(normalizeAgentHistoryMessage).filter(message=>message.text);
    const current = state.agentChats[agentID] || [];
    const existingKeys = new Set(persisted.map(message=>`${message.role}|${message.text}|${message.time}`));
    const pending = current.filter(message=>!existingKeys.has(`${message.role}|${message.text}|${message.time}`));
    state.agentChats[agentID] = persisted.concat(pending);
    state.agentChatHistoryLoaded[agentID] = true;
    state.agentChatHistoryLoading[agentID] = false;
    refreshAgentChat(index, {forceBottom:true});
  } catch (error) {
    state.agentChatHistoryLoading[agentID] = false;
    state.agentChatHistoryLoaded[agentID] = true;
    refreshAgentChat(index, {forceBottom:true});
    toast('历史消息读取失败', error.message, 'warning');
  }
}

function cloneChatAttachments(attachments=[]) {
  return attachments.map(attachment=>({
    name: attachment.name,
    path: attachment.path,
    mime: attachment.mime,
    size: attachment.size,
    previewUrl: attachment.previewUrl || '',
    localId: attachment.localId || '',
  }));
}

function payloadChatAttachments(attachments=[]) {
  return attachments.map(attachment=>({
    name: attachment.name,
    path: attachment.path,
    mime: attachment.mime,
    size: attachment.size,
  }));
}

function removeChatAttachment(index, attachmentIndex) {
  const [attachment] = state.chatAttachments.splice(Number(attachmentIndex), 1);
  if (attachment?.previewUrl) URL.revokeObjectURL(attachment.previewUrl);
  refreshAgentChat(Number(index || 0));
}

function startAgentProgress(agentID,index,hasAttachments=false) {
  stopAgentProgress(agentID);
  const steps = chatProgressSteps
    .filter(([key])=>hasAttachments || key !== 'attachments')
    .map(([key,label])=>({key,label}));
  state.agentChatProgress[agentID] = {index:0,status:'running',steps,startedAt:Date.now(),detail:'STA-100 已接收消息，准备提交 OpenClaw'};
  refreshAgentChat(index);
  chatProgressTimers[agentID] = window.setInterval(()=>{
    const progress = state.agentChatProgress[agentID];
    if (!progress || progress.status !== 'running') return;
    const maxIndex = Math.max(0, progress.steps.length - 2);
    progress.index = Math.min(maxIndex, Number(progress.index || 0) + 1);
    progress.detail = progress.steps[progress.index]?.label ? `${progress.steps[progress.index].label}处理中...` : 'OpenClaw 正在处理';
    const status = document.getElementById('chatStatus');
    if (status) {
      status.innerHTML = renderAgentChatProgress(agentID);
      applyIcons();
    }
  }, 1800);
}

function stopAgentProgress(agentID) {
  if (chatProgressTimers[agentID]) {
    window.clearInterval(chatProgressTimers[agentID]);
    delete chatProgressTimers[agentID];
  }
}

function finishAgentProgress(agentID,index,result,errorMessage='') {
  stopAgentProgress(agentID);
  const failed = Boolean(errorMessage);
  const previous = state.agentChatProgress[agentID] || {};
  const elapsedSeconds = previous.startedAt ? Math.max(0, Math.floor((Date.now() - previous.startedAt) / 1000)) : 0;
  const pipeline = Array.isArray(result?.pipeline) && result.pipeline.length ? result.pipeline.map(normalizePipelineStage) : null;
  const steps = pipeline ? pipeline.map(stage=>({key:stage.key,label:stage.label})) : (previous.steps || chatProgressSteps.map(([key,label])=>({key,label})));
  const failedIndex = pipeline ? Math.max(0,pipeline.findIndex(stage=>stage.status==='failed')) : -1;
  state.agentChatProgress[agentID] = {
    steps,
    index: failed ? (failedIndex >= 0 ? failedIndex : Math.max(0,steps.length-1)) : Math.max(0,steps.length-1),
    status: failed ? 'failed' : 'done',
    elapsedSeconds,
    detail: failed ? errorMessage : 'OpenClaw 已返回结果，消息已完成汇总',
  };
  refreshAgentChat(index);
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
  const model = normalizeAgentMessageModel(agentSelectedModel(agentID));
  const sources = ['本地业务数据库','客户私有知识库','联网检索'];
  const allowlist = getAgentAllowlist(agentID);
  const history = state.agentChats[agentID] || (state.agentChats[agentID] = []);
  const payloadAttachments = payloadChatAttachments(state.chatAttachments);
  const displayAttachments = cloneChatAttachments(state.chatAttachments);
  state.chatAttachments = [];
  const userMessage = {role:'user',text:message,sources,attachments:displayAttachments,attachmentStatus:displayAttachments.length?'正在提交给 OpenClaw':'',time:formatLocalizedTime(new Date().toISOString())};
  history.push(userMessage);
  state.agentChatAtBottom[agentID] = true;
  if (input && !providedMessage) input.value = '';
  startAgentProgress(agentID,index,payloadAttachments.length>0);
  const sendButton = document.querySelector('.chat-send');
  if (sendButton) sendButton.disabled = true;
  applyIcons();
  try {
    const requestBody = {page:'agents',feature:'agent-chat',message,attachments:payloadAttachments,sessionKey:`sta100-${agentID}`,context:{targetAgent:agentID,allowlist}};
    if (model) requestBody.model = model;
    const result = await apiFetch('/api/v1/assistant/query', {method:'POST', body:JSON.stringify(requestBody)});
    userMessage.attachmentStatus = displayAttachments.length ? `已随消息提交给 OpenClaw（${result.attachments?.length || displayAttachments.length} 个附件）` : '';
    applyTokenUsage(result.tokenUsage);
    history.push({role:'agent',text:result.text || 'OpenClaw 未返回文本内容。',sources:result.usedAgents||[],pipeline:result.pipeline||[],time:formatLocalizedTime(new Date().toISOString())});
    finishAgentProgress(agentID,index,result);
  } catch (error) {
    userMessage.attachmentStatus = displayAttachments.length ? '发送失败' : '';
    history.push({role:'agent',text:`调用失败：${error.message}`,error:true,retry:message,sources,pipeline:[{stage:'openclaw-agent',status:'failed',detail:error.message}],time:formatLocalizedTime(new Date().toISOString())});
    finishAgentProgress(agentID,index,null,error.message);
  }
  refreshAgentChat(index, {forceBottom:true});
}

function openAgentAllowlist(index) {
  const agentID = agentIDs[index];
  const domains = getAgentAllowlist(agentID).join('\n');
  openDrawer({title:'联网来源白名单',eyebrow:`${agents[index][0]} / ${agentID}`,body:`<div class="form-field"><label>允许访问的域名</label><textarea class="textarea allowlist-editor" id="agentAllowlist" maxlength="2000" placeholder="每行一个域名，不含协议和路径">${escapeHTML(domains)}</textarea><small>每行一个域名，例如 eur-lex.europa.eu。禁止填写 http://、路径、端口或通配符。</small></div><div class="source-policy-note">${icon('info')}当前白名单会作为来源约束提交给 Agent；联网工具接入后还需在工具调用层执行同一白名单，形成强制访问边界。</div><div class="inline-actions" style="margin-top:16px"><button class="button" data-action="close-drawer">取消</button><button class="button primary" data-action="save-agent-allowlist" data-agent="${index}">${icon('check')}保存白名单</button></div>`});
}

async function saveAgentAllowlist(index) {
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
  try { await savePreferences(); localStorage.removeItem('sta100-agent-allowlists'); closeDrawer(); toast('白名单已保存',`${agents[index][0]} 允许 ${state.agentInternetAllowlists[agentID].length} 个联网来源。`); }
  catch(error) { toast('白名单保存失败',error.message,'warning'); }
}

function newCustomerForm(customer) {
  const c=customer||{};
  state.formContext = { type: 'customer', id: c.id || '' };
  openModal({title:customer?'编辑客户':'新建客户',eyebrow:'客户档案',body:`<div class="form-grid"><div class="form-field full"><label>从名片或照片识别</label><button class="upload-zone" data-action="mock-ocr" style="min-height:92px">${icon('scan-line')}<span>选择图片后识别并填充字段</span></button></div><div class="form-section"><h3>基本信息</h3><p>列表字段可配置，详情页保留全部字段。</p></div>${inputField('客户名称',c.name||'',true,false,'text','customerName')}${selectField('客户类型',['Distributor','Importer','Customer','Reseller','Integrator','Supplier','Other'],false,'customerTypeForm',c.type||'Customer')}${inputField('主电话',c.phone||'',true,false,'tel','customerPhone')}${inputField('网站',c.website||'',false,false,'url','customerWebsite')}${inputField('账单国家',c.country||'',true,false,'text','customerCountryForm')}${inputField('账单城市',c.city||'',false,false,'text','customerCity')}${inputField('联系人',c.contact||'',false,false,'text','customerContact')}${inputField('联系邮箱',c.email||'',false,false,'email','customerEmail')}${selectField('评级',['Prospect','Active','Acquired','Market Failed'],false,'customerRating',c.rating||'Prospect')}${selectField('来源',['展会','电话','朋友介绍','拜访','互联网线索','客户转介绍','其它'],false,'customerSource',c.source||'其它')}${inputField('负责人',c.owner||'Donald',false,false,'text','customerOwner')}${inputField('描述',c.description||'',false,true,'text','customerDescription')}${customer?`<div class="form-section"><h3>历史沟通记录</h3><p>沟通记录独立保存且只能追加，编辑客户不会覆盖历史。</p></div><div class="form-field full"><button type="button" class="button" data-action="customer-communications" data-id="${escapeAttr(c.id)}">${icon('messages-square')}查看或新增沟通记录</button></div>`:''}</div>`,footer:formFooter(customer?'保存修改':'创建客户','save-customer')});
}

async function customerDetail(id, tab='overview') {
  const c=customers.find(x=>x.id===id);
  if (!c) return;
  const tabs = [['overview','概览'],['contacts','联系人'],['quotes','报价单'],['orders','订单'],['documents','单据'],['activity','沟通记录']];
  const tabBar = `<div class="tabs">${tabs.map(([key,label])=>`<button class="${tab===key?'active':''}" data-customer-tab="${key}" data-customer-id="${escapeAttr(c.id)}">${label}</button>`).join('')}</div>`;
  const relatedQuotes = quotes.filter(q=>q.customer===c.name);
  const relatedOrders = orders.filter(o=>o.customer===c.name);
  const relatedDocuments = documents.filter(d=>d.customer===c.name);
  if (tab === 'activity' && !Array.isArray(state.customerCommunications[id])) {
    try {
      const result = await apiFetch(`/api/v1/accounts/${encodeURIComponent(id)}/communications`);
      state.customerCommunications[id] = result.items || [];
    } catch (error) {
      toast('沟通记录加载失败',error.message,'warning');
      state.customerCommunications[id] = [];
    }
  }
  const communications = state.customerCommunications[id] || [];
  let content = '';
  if (tab === 'quotes') content = `<div class="related-list">${relatedQuotes.map(q=>`<button class="related-record" data-action="quote-detail" data-id="${escapeAttr(q.id)}"><span>${icon('file-text')}<strong>${escapeHTML(q.id)}</strong></span><span>${escapeHTML(q.subject)}</span><span>${escapeHTML(q.value)}</span>${badge(q.status)}</button>`).join('') || `<div class="empty-state"><p>暂无关联报价单</p></div>`}</div>`;
  else if (tab === 'orders') content = `<div class="related-list">${relatedOrders.map(o=>`<button class="related-record" data-action="order-detail" data-id="${escapeAttr(o.id)}"><span>${icon('package')}<strong>${escapeHTML(o.id)}</strong></span><span>${escapeHTML(o.products)}</span><span>${escapeHTML(o.value)}</span>${badge(o.status)}</button>`).join('') || `<div class="empty-state"><p>暂无关联订单</p></div>`}</div>`;
  else if (tab === 'documents') content = `<div class="related-list">${relatedDocuments.map(d=>`<button class="related-record" data-action="document-detail" data-id="${escapeAttr(d.id)}"><span>${icon('file-check-2')}<strong>${escapeHTML(d.id)}</strong></span><span>${escapeHTML(d.order)}</span><span>${escapeHTML(d.template)}</span>${badge(d.status)}</button>`).join('') || `<div class="empty-state"><p>暂无关联单据</p></div>`}</div>`;
  else if (tab === 'contacts') content = `<div class="detail-grid">${[['联系人',c.contact],['电话',c.phone],['邮箱',c.email],['网站',c.website||'未填写'],['国家/城市',`${c.country}${c.city?' / '+c.city:''}`],['来源',c.source||'未填写']].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div>`;
  else if (tab === 'activity') content = `<div class="spread communication-head"><div><strong>历史沟通记录</strong><p class="secondary-text">记录只能追加，不能修改或删除。</p></div><button class="button primary small" data-action="new-customer-communication" data-id="${escapeAttr(c.id)}">${icon('message-square-plus')}新增沟通</button></div>${communications.length?`<div class="timeline communication-timeline">${communications.map(item=>`<div class="timeline-item"><div class="spread"><h4>${escapeHTML(item.subject||item.type)}</h4><span class="badge blue">${escapeHTML(item.type)}</span></div><p class="communication-content">${escapeHTML(item.content)}</p><small>${escapeHTML(String(item.occurredAt||'').replace('T',' '))}${item.contact?` · ${escapeHTML(item.contact)}`:''} · 由 ${escapeHTML(item.createdBy)} 记录</small></div>`).join('')}</div>`:`<div class="empty-state panel">${icon('messages-square')}<div><h3>暂无沟通记录</h3><p>新增后将永久保留在本机业务数据库中。</p></div></div>`}`;
  else content = `<div class="detail-grid">${[['客户编号',c.id],['客户类型',c.type],['国家',c.country],['城市',c.city||'未填写'],['负责人',c.owner],['主联系人',c.contact],['电话',c.phone],['邮箱',c.email],['来源',c.source||'未填写'],['订单数',String(c.orders)],['累计金额',c.total],['最近更新',c.updated]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 14px">最近业务记录</div><div class="filter-row"><button class="button small" data-customer-tab="quotes" data-customer-id="${escapeAttr(c.id)}">报价单 ${relatedQuotes.length}</button><button class="button small" data-customer-tab="orders" data-customer-id="${escapeAttr(c.id)}">订单 ${relatedOrders.length}</button><button class="button small" data-customer-tab="documents" data-customer-id="${escapeAttr(c.id)}">单据 ${relatedDocuments.length}</button></div>`;
  openDrawer({title:c.name,eyebrow:`客户 / ${c.id}`,body:`${tabBar}<div class="spread" style="margin-bottom:14px"><span>${badge(c.rating)}</span><div class="inline-actions"><button class="button small" data-action="edit-customer" data-id="${escapeAttr(c.id)}">${icon('pencil')}编辑</button><button class="button primary small" data-action="new-quote" data-customer="${escapeAttr(c.name)}">${icon('file-plus-2')}新建报价</button><button class="button danger small" data-action="delete-customer" data-id="${escapeAttr(c.id)}">${icon('trash-2')}删除</button></div></div>${content}`});
}

function customerCommunicationForm(customerID) {
  const customer=customers.find(item=>item.id===customerID);
  if(!customer)return;
  closeDrawer();
  const now=new Date(Date.now()-new Date().getTimezoneOffset()*60000).toISOString().slice(0,16);
  openModal({title:'新增沟通记录',eyebrow:`${customer.name} / 只能追加`,body:`<div class="form-grid">${selectField('沟通方式',['电话','邮件','微信','会议','拜访','短信','其它'],false,'communicationType','电话')}${inputField('沟通时间',now,true,false,'datetime-local','communicationOccurredAt')}${inputField('主题','',false,true,'text','communicationSubject')}${inputField('联系人',customer.contact||'',false,true,'text','communicationContact')}<div class="form-field full"><label>沟通内容 <span class="required">*</span></label><textarea class="textarea" id="communicationContent" maxlength="10000" placeholder="记录客户反馈、结论和后续事项"></textarea><small>保存后不可修改或删除，请确认内容准确。</small></div></div>`,footer:`<button class="button" data-action="cancel-customer-communication" data-id="${escapeAttr(customerID)}">取消</button><button class="button primary" data-action="save-customer-communication" data-id="${escapeAttr(customerID)}">${icon('plus')}追加记录</button>`});
}

async function saveCustomerCommunication(customerID) {
  const content=formText('communicationContent');
  const occurredAt=formText('communicationOccurredAt');
  if(!content||!occurredAt){toast('保存失败','沟通时间和沟通内容为必填项。','warning');return;}
  const payload={type:formText('communicationType'),subject:formText('communicationSubject'),content,contact:formText('communicationContact'),occurredAt};
  try {
    const record=await apiFetch(`/api/v1/accounts/${encodeURIComponent(customerID)}/communications`,{method:'POST',body:JSON.stringify(payload)});
    state.customerCommunications[customerID]=[record,...(state.customerCommunications[customerID]||[])];
    closeModal();
    await customerDetail(customerID,'activity');
    toast('沟通记录已追加','该记录已永久保留，不能修改或删除。');
  } catch(error) { toast('沟通记录保存失败',error.message,'warning'); }
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
  openModal({title:document?'编辑单据':'生成外贸单据',eyebrow:orderId==='all'?'一键生成全套':`订单 / ${orderId||d.order||''}`,body:`<div class="form-grid">${relationField('关联订单','documentOrder',orders.map(o=>`${o.id} · ${o.customer}`),d.order?`${d.order} · ${d.customer}`:orderId&&orderId!=='all'?`${orderId} · ${orders.find(o=>o.id===orderId)?.customer||''}`:'')}${selectField('单据类型',document?['PI','CI','PL','报关单']:['PI','CI','PL','报关单','全套（PI + CI + PL + 报关单）'],false,'documentTypeForm',d.type||(orderId==='all'?'全套（PI + CI + PL + 报关单）':'PI'))}${selectField('模板版本',['各类型默认模板','STRATRONIX 标准模板组 v3'],false,'documentTemplate',d.template||'各类型默认模板')}${selectField('输出语言',['英文','中文 / 英文双语'],false,'documentLanguage','英文')}<div class="form-section"><h3>生成规则</h3><p>系统读取客户、订单和产品快照，生成后先进入待复核状态；缺失字段会逐项提示。</p></div><div class="model-warning full"><span>${icon('info')} 报关单的申报要素、监管条件和最终格式需在正式数据提供后校准。</span></div></div>`,footer:formFooter(document?'保存修改':'生成并预览','save-document')});
}

async function saveCustomer() {
  const name=formText('customerName');
  if (!name || !formText('customerCountryForm')) { toast('保存失败','客户名称和账单国家为必填项。','warning'); return; }
  const existing=customers.find(c=>c.id===state.formContext?.id);
  const payload={...(existing||{}),name,type:formText('customerTypeForm'),phone:formText('customerPhone'),website:formText('customerWebsite'),country:formText('customerCountryForm'),city:formText('customerCity'),contact:formText('customerContact'),email:formText('customerEmail'),rating:formText('customerRating'),source:formText('customerSource'),owner:formText('customerOwner'),description:formText('customerDescription'),archived:false};
  try {
    const record=await apiFetch(existing?`/api/v1/accounts/${encodeURIComponent(existing.id)}`:'/api/v1/accounts',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)});
    upsertRecord(customers,record); closeModal(); renderPage(); toast(existing?'客户已更新':'客户已创建',`${record.name} 已保存到本地数据库。`);
  } catch(error) { toast('保存失败',error.message,'warning'); }
}

async function saveQuote() {
  const subject=formText('quoteSubject'); const customer=formText('quoteCustomer');
  if (!subject || !customer || !customers.some(c=>c.name===customer&&!c.archived)) { toast('保存失败','请填写主题，并从客户搜索结果中选择有效客户。','warning'); return; }
  if (!state.quoteDraftLines.length || state.quoteDraftLines.some(line => Number(line.quantity) < 1 || Number(line.unitPrice) < 0)) { toast('保存失败','报价至少需要一条有效产品明细，数量和单价不能为负数。','warning'); return; }
  const existing=quotes.find(q=>q.id===state.formContext?.id);
  const currency=formText('quoteCurrency')||'EUR';
  const payload={...(existing||{}),subject,customer,valid:formText('quoteValid'),currency,freight:formNumber('quoteFreight'),tax:formNumber('quoteTax'),terms:formText('quoteTerms'),lines:state.quoteDraftLines.map(line=>({...line,quantity:Number(line.quantity),unitPrice:Number(line.unitPrice),discount:Number(line.discount||0) }))};
  try {
    const record=await apiFetch(existing?`/api/v1/quotes/${encodeURIComponent(existing.id)}`:'/api/v1/quotes',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)});
    upsertRecord(quotes,record); closeModal(); renderPage(); toast(existing?'报价单已更新':'报价草稿已创建',`${record.id} 已保存到本地数据库。`);
  } catch(error) { toast('保存失败',error.message,'warning'); }
}

async function saveOrder() {
  const customer=formText('orderCustomer'); const quoteValue=formText('orderQuote');
  if (!customer || !customers.some(c=>c.name===customer&&!c.archived)) { toast('保存失败','请从客户搜索结果中选择有效客户。','warning'); return; }
  if (!state.orderDraftLines.length || state.orderDraftLines.some(line => Number(line.quantity) < 1 || Number(line.unitPrice) < 0)) { toast('保存失败','订单至少需要一条有效产品明细，数量和成交单价不能为负数。','warning'); return; }
  const insufficient=state.orderDraftLines.find(line=>Number(line.quantity)>Number(lineProduct(line)?.stock||0));
  if (insufficient) { toast('保存失败',`${lineProduct(insufficient)?.name||'产品'} 数量超过当前库存 ${lineProduct(insufficient)?.stock||0}。`,'warning'); return; }
  const existing=orders.find(o=>o.id===state.formContext?.id);
  const payload={...(existing||{}),customer,quote:quoteValue.split(' · ')[0]||quoteValue,po:formText('orderPO'),delivery:formText('orderDelivery'),terms:formText('orderTerms'),currency:existing?.currency||'EUR',status:existing?.status||'Confirmed',progress:Number(existing?.progress||0),lines:state.orderDraftLines.map(line=>({...line,quantity:Number(line.quantity),unitPrice:Number(line.unitPrice)}))};
  try {
    const record=await apiFetch(existing?`/api/v1/orders/${encodeURIComponent(existing.id)}`:'/api/v1/orders',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)});
    upsertRecord(orders,record); await loadBusinessData(true); closeModal(); renderPage(); toast(existing?'订单已更新':'订单已创建',`${record.id} 已保存，已关联 ${record.lines.length} 个产品明细。`);
  } catch(error) { toast('保存失败',error.message,'warning'); }
}

function refreshCustomerAggregates() {
  customers.forEach(customer => {
    const related = orders.filter(order => order.customer === customer.name && order.status !== 'Cancelled');
    if (!related.length) return;
    customer.orders = related.length;
    customer.total = formatMoney(related.reduce((sum, order) => sum + moneyNumber(order.value), 0), 'EUR');
  });
}

async function saveDocument() {
  const orderValue=formText('documentOrder'); const orderId=orderValue.split(' · ')[0]; const order=orders.find(o=>o.id===orderId);
  if (!order) { toast('保存失败','请从订单搜索结果中选择有效订单。','warning'); return; }
  const existing=documents.find(d=>d.id===state.formContext?.id);
  const type=formText('documentTypeForm');
  const payload={...(existing||{}),type,customer:order.customer,order:order.id,template:formText('documentTemplate'),language:formText('documentLanguage'),status:existing?.status||'Review',lines:businessLines(order),value:order.value};
  try {
    if(!existing&&type.startsWith('全套')){
      const result=await apiFetch(`/api/v1/orders/${encodeURIComponent(order.id)}/documents`,{method:'POST',body:JSON.stringify({types:['PI','CI','PL','报关单'],template:payload.template,language:payload.language})});
      result.items.forEach(record=>upsertRecord(documents,record));closeModal();renderPage();toast('全套单据已生成',`${result.total} 份单据已进入待复核状态。`);return;
    }
    const record=await apiFetch(existing?`/api/v1/documents/${encodeURIComponent(existing.id)}`:'/api/v1/documents',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)});
    upsertRecord(documents,record); closeModal(); renderPage(); toast(existing?'单据已更新':'单据已生成',`${record.id} 已保存为待复核状态。`);
  } catch(error) { toast('保存失败',error.message,'warning'); }
}

async function deleteCustomer(id) {
  const customer=customers.find(c=>c.id===id); if (!customer) return;
  if (!window.confirm(`确定归档客户 ${customer.name} 吗？关联报价、订单和单据会保留。`)) return;
  try { await apiFetch(`/api/v1/accounts/${encodeURIComponent(id)}`,{method:'DELETE'}); customer.archived=true; closeDrawer(); renderPage(); toast('客户已归档','关联业务记录仍然保留，可用于历史查询。'); } catch(error) { toast('归档失败',error.message,'warning'); }
}
async function deleteQuote(id) { if (!window.confirm('确定归档这份报价单吗？')) return; try { await apiFetch(`/api/v1/quotes/${encodeURIComponent(id)}`,{method:'DELETE'}); removeRecord(quotes,id); closeDrawer(); renderPage(); toast('报价单已归档','变更已写入本地数据库。'); } catch(error) { toast('归档失败',error.message,'warning'); } }
async function deleteOrder(id) { if (!window.confirm('确定归档这份订单吗？')) return; try { await apiFetch(`/api/v1/orders/${encodeURIComponent(id)}`,{method:'DELETE'}); removeRecord(orders,id); await loadBusinessData(true); closeDrawer(); renderPage(); toast('订单已归档','变更已写入本地数据库。'); } catch(error) { toast('归档失败',error.message,'warning'); } }
async function deleteDocument(id) { if (!window.confirm('确定归档这份单据吗？')) return; try { await apiFetch(`/api/v1/documents/${encodeURIComponent(id)}`,{method:'DELETE'}); removeRecord(documents,id); closeDrawer(); renderPage(); toast('单据已归档','变更已写入本地数据库。'); } catch(error) { toast('归档失败',error.message,'warning'); } }
async function downloadQuote(id) {
  try { await apiFetch(`/api/v1/quotes/${encodeURIComponent(id)}/download`); }
  catch(error) { toast('报价 PDF 暂不可生成',error.message,'warning'); }
}
async function downloadDocument(id) {
  try { await apiFetch(`/api/v1/documents/${encodeURIComponent(id)}/download`); }
  catch(error) { toast('单据文件暂不可生成',error.message,'warning'); }
}

function openNotifications() {
  const pendingFiles=files.filter(file=>file.status!=='Indexed');
  const shippedOrders=orders.filter(order=>order.status==='Shipped');
  const openClawOK=Boolean(state.systemHealth?.openclaw?.ok);
  const entries=[];
  if(pendingFiles.length) entries.push([`${pendingFiles.length} 份文件等待解析或复核`,'数据库 · 当前状态']);
  if(shippedOrders.length) entries.push([`${shippedOrders.length} 个订单处于已发运状态`,'订单 · 当前状态']);
  entries.push([openClawOK?'OpenClaw Gateway 当前可用':'OpenClaw Gateway 状态待检查','设置 · 实时健康状态']);
  openDrawer({title:'通知',eyebrow:'当前系统状态',body:`<div class="timeline">${entries.map(([title,detail])=>`<div class="timeline-item"><h4>${escapeHTML(title)}</h4><p>${escapeHTML(detail)}</p></div>`).join('')}</div>`});
}

async function convertQuoteToOrder(id) {
  try { const record=await apiFetch(`/api/v1/quotes/${encodeURIComponent(id)}/convert-order`,{method:'POST',body:'{}'}); upsertRecord(orders,record); await loadBusinessData(true); closeDrawer(); setPage('orders'); toast('订单已创建',`${record.id} 已保留报价来源和产品快照。`); }
  catch(error) { toast('转换失败',error.message,'warning'); }
}

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

async function handleTemplateFile(file,mode,kind) {
  if(!file)return;
  const allowed=mode==='image'?/\.(jpe?g|png|webp)$/i:/\.(docx|xlsx|pdf|html?)$/i;
  if(!allowed.test(file.name)){toast('文件格式不支持',mode==='image'?'请选择 JPG、PNG 或 WebP 图片。':'请选择 DOCX、XLSX、PDF 或 HTML 模板。','warning');return;}
  if(file.size>20*1024*1024){toast('文件过大','模板文件不能超过 20 MB。','warning');return;}
  const size=file.size>=1024*1024?`${(file.size/1024/1024).toFixed(1)} MB`:`${Math.max(1,Math.round(file.size/1024))} KB`;
  state.templateUploads.push({kind,mode,name:file.name,size,status:'waiting-api'});
  templateCenter(kind);
  const form=new FormData();form.append('file',file);form.append('kind',kind);
  try {await apiFetch(mode==='image'?'/api/v1/templates/image-recognition':'/api/v1/templates/upload',{method:'POST',body:form});}
  catch(error){toast('模板暂不能处理',error.message,'warning');}
}

async function templateAction(action) {
  try {await apiFetch(`/api/v1/templates/${encodeURIComponent(`${state.templateKind}-${action}`)}`,{method:'PATCH',body:JSON.stringify({action,kind:state.templateKind})});}
  catch(error){toast('模板操作暂不可用',error.message,'warning');}
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

async function saveProduct() {
  const name=formText('productName'); const id=formText('productID'); if(!name||!id){toast('保存失败','产品名称和编码为必填项。','warning');return;}
  const existing=products.find(p=>p.id===state.formContext?.id);
  const payload={...(existing||{}),id,name,category:formText('productCategory'),manufacturer:formText('productManufacturer'),hs:formText('productHS'),stock:formNumber('productStock'),price:formText('productPrice'),desc:formText('productDescription'),tags:formText('productTags'),status:existing?.status||'Active'};
  try { const record=await apiFetch(existing?`/api/v1/products/${encodeURIComponent(existing.id)}`:'/api/v1/products',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)}); upsertRecord(products,record); closeModal(); renderPage(); toast(existing?'产品已更新':'产品已创建',`${record.name} 已保存到本地数据库。`); }
  catch(error) { toast('保存失败',error.message,'warning'); }
}
async function deleteProduct(id) { if(!window.confirm('确定停用该产品吗？'))return; try { await apiFetch(`/api/v1/products/${encodeURIComponent(id)}`,{method:'DELETE'}); const product=products.find(item=>item.id===id); if(product)product.status='Inactive'; closeDrawer(); renderPage(); toast('产品已停用','产品主数据已更新。'); } catch(error) { toast('停用失败',error.message,'warning'); } }

function productImportModal() {
  openModal({title:'批量导入产品',eyebrow:'产品库 / 文件导入',body:`<input id="productImportInput" type="file" accept=".xlsx,.csv" hidden><div class="upload-zone"><div><span class="upload-icon">${icon('file-up')}</span><h3>选择产品导入文件</h3><p>支持 XLSX 或 CSV。导入模板、必填字段、重复编码和错误回执规则尚待确认。</p><button type="button" class="button primary" data-action="choose-product-import">${icon('folder-open')}选择文件</button></div></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});
  document.getElementById('productImportInput')?.addEventListener('change',async event=>{
    const file=event.target.files?.[0];if(!file)return;
    const form=new FormData();form.append('file',file);
    try {await apiFetch('/api/v1/products/import',{method:'POST',body:form});}
    catch(error){toast('产品批量导入暂不可用',error.message,'warning');}
  });
}

function customerOCRModal() {
  openModal({title:'识别客户名片或照片',eyebrow:'客户档案 / OCR',body:`<input id="customerOCRInput" type="file" accept="image/jpeg,image/png,image/webp" hidden><div class="upload-zone"><div><span class="upload-icon">${icon('scan-line')}</span><h3>选择客户名片或照片</h3><p>图片只会提交到本机 Go 接口；OCR 引擎、字段映射和置信度规则尚待确认。</p><button type="button" class="button primary" data-action="choose-customer-ocr">${icon('image-up')}选择图片</button></div></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});
  document.getElementById('customerOCRInput')?.addEventListener('change',async event=>{
    const file=event.target.files?.[0];if(!file)return;
    const form=new FormData();form.append('file',file);
    try {await apiFetch('/api/v1/accounts/ocr',{method:'POST',body:form});}
    catch(error){toast('图片识别暂不可用',error.message,'warning');}
  });
}

function dateFilterForm(module) {
  const isQuote=module==='quote';
  const from=isQuote?state.quoteDateFrom:state.orderDateFrom;
  const to=isQuote?state.quoteDateTo:state.orderDateTo;
  openModal({title:isQuote?'报价有效期筛选':'订单交付日期筛选',eyebrow:'日期范围',body:`<div class="form-grid">${inputField('开始日期',from,false,false,'date','dateFilterFrom')}${inputField('结束日期',to,false,false,'date','dateFilterTo')}</div>`,footer:`<button class="button" data-action="close-modal">取消</button><button class="button ghost" data-action="apply-date-filter" data-module="${module}" data-clear="true">清除</button><button class="button primary" data-action="apply-date-filter" data-module="${module}">${icon('check')}应用</button>`});
}

function applyDateFilter(module, clear=false) {
  const from=clear?'':formText('dateFilterFrom');
  const to=clear?'':formText('dateFilterTo');
  if(from&&to&&from>to){toast('日期范围无效','开始日期不能晚于结束日期。','warning');return;}
  if(module==='quote'){state.quoteDateFrom=from;state.quoteDateTo=to;}else{state.orderDateFrom=from;state.orderDateTo=to;}
  closeModal();renderPage();
}

function productDetail(id) {
  const p=products.find(x=>x.id===id);
  if(!p)return;
  openDrawer({title:p.name,eyebrow:`产品 / ${p.id}`,body:`<div class="spread"><span>${badge(p.status)}</span><div class="inline-actions"><button class="button small" data-action="edit-product" data-id="${escapeAttr(p.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-product" data-id="${escapeAttr(p.id)}">${icon('trash-2')}删除</button></div></div><div class="product-visual panel" style="margin-top:14px;aspect-ratio:16/6">${icon('cpu')}</div><div class="detail-grid" style="margin-top:14px">${[['产品编码',p.id],['产品类别',p.category],['HS CODE',p.hs],['销售价',p.price],['当前库存',String(p.stock)],['状态',p.status]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">产品描述</div><p class="secondary-text" style="line-height:1.7">${escapeHTML(p.desc)}</p><div class="divider-title" style="margin:20px 0 12px">业务引用</div><div class="filter-row"><span class="badge blue">报价单 8</span><span class="badge amber">进行中订单 3</span><span class="badge green">历史单据 12</span></div>`});
}

function uploadFileModal() {
  state.selectedUploadFile=null;
  openModal({title:'上传私有数据',eyebrow:'数据库 / 文件处理',body:`<input id="privateFileInput" type="file" accept=".pdf,.docx,.xlsx,.csv,.txt,.md,.jpg,.jpeg,.png" hidden><div class="upload-zone" id="uploadZone"><div><span class="upload-icon">${icon('cloud-upload')}</span><h3>选择或拖入文件</h3><p>支持 PDF、DOCX、XLSX、CSV、TXT、MD、JPG、PNG；单文件最大 50 MB。</p><button type="button" class="button primary" data-action="choose-file">选择文件</button><div id="selectedPrivateFile"></div></div></div><div class="form-grid" style="margin-top:14px"><div class="form-field"><label>数据区</label><input class="input" value="客户私有数据（本机）" readonly></div>${selectField('主分类',['自动识别','合同','报价单','产品手册','法规','产品资料','会议记录','客户资料','其它'],false,'privateFileCategory','自动识别')}${inputField('附加标签','',false,true,'text','privateFileTags')}<div class="form-field full"><small>当前后端完成格式、大小、SHA-256 去重和本机安全保存。文本/OCR 解析、AI 分类与索引等待客户原始数据格式后启用。</small></div></div>`,footer:formFooter('上传到本机','upload-private-file')});
  const zone=document.getElementById('uploadZone');
  ['dragenter','dragover'].forEach(ev=>zone.addEventListener(ev,e=>{e.preventDefault();zone.classList.add('dragover')}));
  zone.addEventListener('dragleave',e=>{e.preventDefault();zone.classList.remove('dragover')});
  zone.addEventListener('drop',e=>{e.preventDefault();zone.classList.remove('dragover');selectPrivateFile(e.dataTransfer?.files?.[0])});
  document.getElementById('privateFileInput').addEventListener('change',e=>selectPrivateFile(e.target.files?.[0]));
}

function selectPrivateFile(file) {
  if(!file)return;
  const allowed=/\.(pdf|docx|xlsx|csv|txt|md|jpe?g|png)$/i;
  if(!allowed.test(file.name)){toast('文件格式不支持','请选择页面列出的文件格式。','warning');return;}
  if(file.size>50*1024*1024){toast('文件过大','单文件不能超过 50 MB。','warning');return;}
  state.selectedUploadFile=file;
  const target=document.getElementById('selectedPrivateFile');
  if(target)target.innerHTML=`<div class="upgrade-file" style="margin-top:12px"><span class="upload-icon">${icon('file-check-2')}</span><div><strong>${escapeHTML(file.name)}</strong><small>${formatBytes(file.size)} · 等待上传</small></div></div>`;
  applyIcons();
}

async function uploadPrivateFile() {
  const file=state.selectedUploadFile;
  if(!file){toast('请选择文件','需要先选择或拖入一个文件。','warning');return;}
  const form=new FormData();
  form.append('file',file);
  form.append('category',formText('privateFileCategory'));
  form.append('tags',formText('privateFileTags'));
  try {
    const result=await apiFetch('/api/v1/private-files/upload',{method:'POST',body:form});
    upsertRecord(files,result.item);
    state.selectedUploadFile=null;
    closeModal();renderPage();toast('文件已保存到本机',`${result.item.name} 当前状态：${result.item.status}`);
  } catch(error) { toast('文件上传失败',error.message,'warning'); }
}

function filePreview(id) {
  const f=files.find(x=>x.id===id);
  if (!f) return;
  const previewable=f.mime==='application/pdf'||String(f.mime||'').startsWith('image/');
  const preview=previewable?`<iframe title="${escapeAttr(f.name)}" src="/api/v1/private-files/${encodeURIComponent(f.id)}/content" style="width:100%;height:520px;border:1px solid var(--line);margin-top:16px;background:white"></iframe>`:`<div class="empty-state panel" style="margin-top:16px">${icon('file-search')}<div><h3>当前格式不支持内嵌预览</h3><p>可以下载原文件；摘要和索引需等待客户数据格式确认。</p><button class="button small" data-action="file-summary" data-id="${escapeAttr(f.id)}">${icon('sparkles')}查看摘要状态</button></div></div>`;
  openDrawer({title:f.name,eyebrow:`${f.category} / 在线预览`,body:`<div class="spread"><span>${badge(f.status)}</span><div class="inline-actions"><button class="button small" data-action="file-download" data-id="${escapeAttr(f.id)}">${icon('download')}下载</button><button class="button small" data-action="file-edit" data-id="${escapeAttr(f.id)}">${icon('pencil')}编辑信息</button></div></div><div class="detail-grid" style="margin-top:15px">${[['主分类',f.category],['来源',f.source],['大小',f.size],['更新时间',f.updated],['标签',(f.tags||[]).join(' / ')||'无'],['索引状态',f.status]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div>${preview}`});
}

function fileDownload(id) {
  const f = files.find(item => item.id === id);
  if (!f) return;
  window.location.assign(`/api/v1/private-files/${encodeURIComponent(f.id)}/download`);
}

function fileEdit(id) {
  const f = files.find(item => item.id === id);
  if (!f) return;
  state.formContext={type:'private-file',id:f.id};
  openModal({ title: '编辑文件信息', eyebrow: '数据库 / 文件元数据', body: `<div class="form-grid">${inputField('文件名',f.name,true,true,'text','privateFileName')}${selectField('主分类',['待分类','产品手册','合同','报价单','法规','产品资料','会议记录','客户资料','其它'],false,'privateFileEditCategory',f.category)}${inputField('标签',(f.tags||[]).join(' / '),false,true,'text','privateFileEditTags')}${inputField('来源',f.source,false,true,'text','privateFileSource')}</div>`, footer: formFooter('保存信息','save-file-metadata') });
}

async function saveFileMetadata() {
  const f=files.find(item=>item.id===state.formContext?.id);if(!f)return;
  const payload={...f,name:formText('privateFileName'),category:formText('privateFileEditCategory'),tags:formText('privateFileEditTags').split(/[、,，/]/).map(value=>value.trim()).filter(Boolean),source:formText('privateFileSource')};
  if(!payload.name||!payload.category){toast('保存失败','文件名和主分类不能为空。','warning');return;}
  try {const record=await apiFetch(`/api/v1/private-files/${encodeURIComponent(f.id)}`,{method:'PATCH',body:JSON.stringify(payload)});upsertRecord(files,record);closeModal();closeDrawer();renderPage();toast('文件信息已更新',record.name);}
  catch(error){toast('保存失败',error.message,'warning');}
}

async function archivePrivateFile(id) {
  const f=files.find(item=>item.id===id);if(!f||!window.confirm(`确定删除本机文件“${f.name}”吗？此操作会同时删除保存的原始文件。`))return;
  try {await apiFetch(`/api/v1/private-files/${encodeURIComponent(id)}`,{method:'DELETE'});removeRecord(files,id);closeModal();closeDrawer();renderPage();toast('文件已删除',f.name);}
  catch(error){toast('删除失败',error.message,'warning');}
}

async function reindexPrivateFile(id) {
  try {const result=await apiFetch(`/api/v1/private-files/${encodeURIComponent(id)}/reindex`,{method:'POST',body:'{}'});upsertRecord(files,result.item);closeModal();renderPage();toast('重新索引请求已记录',result.todo,'warning');}
  catch(error){toast('重新索引失败',error.message,'warning');}
}

async function fileSummary(id) {
  const f=files.find(item=>item.id===id);if(!f)return;
  try {await apiFetch(`/api/v1/private-files/${encodeURIComponent(id)}/summary`);}
  catch(error){toast('摘要暂不可用',error.message,'warning');}
}

function newsDetail(title) {
  const n=news.find(x=>x.title===title)||news[0];
  openDrawer({title:n.title,eyebrow:`${n.category} / ${n.source}`,body:`<div class="spread"><span class="badge green">相关度 ${n.relevance}</span><span class="secondary-text">${n.time}</span></div><p style="margin-top:20px;line-height:1.8;color:var(--text)">${n.summary}</p><div class="divider-title" style="margin:20px 0 12px">智能体摘要</div><p class="secondary-text" style="line-height:1.8">该信息与当前关注的欧洲渠道、智能骑行产品和合规主题相关。建议结合客户档案与产品库，检查受影响客户和产品后再形成行动项。</p><div class="divider-title" style="margin:20px 0 12px">来源信息</div><div class="detail-grid">${[['来源',n.source],['获取时间',n.time],['信息类别',n.category],['数据区域','互联网推荐数据']].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${v}</strong></div>`).join('')}</div><div class="inline-actions" style="margin-top:18px"><button class="button" data-action="news-source-link" data-title="${escapeAttr(n.title)}">${icon('external-link')}查看原文</button><button class="button primary" data-action="news-todo" data-title="${escapeAttr(n.title)}">${icon('list-plus')}生成待办</button></div>`});
}

async function refreshNews() {
  try {await apiFetch('/api/v1/news/refresh',{method:'POST',body:'{}'});}
  catch(error){toast('新闻暂不可更新',error.message,'warning');}
}

function openNewsSource(title) {
  const item=news.find(entry=>entry.title===title);
  if(item?.sourceUrl){window.open(item.sourceUrl,'_blank','noopener,noreferrer');return;}
  toast('原文链接未配置','当前缓存记录没有来源 URL，待来源白名单和抓取规则确认后补充。','warning');
}

async function createNewsTodo(title) {
  try {await apiFetch('/api/v1/tasks',{method:'POST',body:JSON.stringify({title:`跟进资讯：${title}`,source:'industry_news'})});}
  catch(error){toast('待办暂不可生成',error.message,'warning');}
}

function currentModelConfiguration() {
  const model = document.getElementById('modelVersionSelect')?.value || state.modelDraftKey || '';
  const apiKey = document.getElementById('modelAPIKey')?.value.trim() || '';
  const endpointMode = document.getElementById('modelEndpointMode')?.value || modelEndpointModeFor(modelProvider(model));
  const selectedModels = configuredModelKeys().filter(key=>key!==state.modelDraftOriginalKey);
  if (model && !selectedModels.includes(model)) selectedModels.push(model);
  const preferredDefault = state.modelDraftMode === 'edit' && state.modelDraftOriginalKey && currentDefaultModelKey() === state.modelDraftOriginalKey ? model : '';
  const defaultModel = defaultModelForKeys(selectedModels, preferredDefault || model, preferredDefault === model);
  return { model, defaultModel, apiKey, endpointMode, selectedModels };
}

function validateModelConfiguration({model}) {
  if (!model) {
    toast('无法保存', '请选择模型。', 'warning');
    return false;
  }
  return true;
}

async function saveModelConfiguration(button) {
  const configuration = currentModelConfiguration();
  if (!validateModelConfiguration(configuration)) return;
  button.disabled = true;
  button.innerHTML = `${icon('loader-circle')}保存中`;
  applyIcons();
  try {
    const saved = await apiFetch('/api/v1/settings/model', { method: 'PATCH', body: JSON.stringify(configuration) });
    state.modelTestResult=null;
    state.modelDraftMode='create';
    state.modelDraftFamilyKey='';
    state.modelDraftKey='';
    state.modelDraftOriginalKey='';
    closeModal();
    await loadOpenClawModels(true);
    toast('OpenClaw 模型配置已保存', saved.message || `${configuration.model} 已保存。`, 'success');
  } catch (error) {
    button.disabled = false;
    button.innerHTML = `${icon('save')}保存`;
    applyIcons();
    toast(state.modelDraftMode === 'edit' ? '模型配置更新失败' : '模型配置保存失败', error.message, 'warning');
  }
}

async function testCurrentModelConfiguration() {
  const configuration = currentModelConfiguration();
  if (!validateModelConfiguration(configuration)) return;
  await testModelConnection(configuration);
}

function startNewModelConfiguration() {
  openModelConfigurationForm('create', '');
}

function editModelConfiguration(model) {
  openModelConfigurationForm('edit', model);
}

async function deleteModelConfiguration(model) {
  const modelInfo = findModelByKey(model);
  if (!window.confirm(`确认删除模型配置：${modelInfo?.name||model}？\n该操作只会从当前应用/OpenClaw 默认模型列表移除，不会回传或展示已保存 API Key。`)) return;
  const selectedModels = configuredModelKeys().filter(key=>key!==model);
  const defaultModel = defaultModelForKeys(selectedModels);
  try {
    await apiFetch('/api/v1/settings/model', { method:'PATCH', body: JSON.stringify({ model:'', defaultModel, apiKey:'', selectedModels }) });
    if (state.modelDraftKey===model) {
      state.modelDraftMode='create';
      state.modelDraftFamilyKey='';
      state.modelDraftKey='';
      state.modelDraftOriginalKey='';
    }
    state.modelTestResult=null;
    await loadOpenClawModels(true);
    toast('模型配置已删除', modelInfo?.name||model, 'success');
  } catch(error) {
    toast('模型配置删除失败', error.message, 'warning');
  }
}

async function testConfiguredModel(model) {
  const selectedModels = configuredModelKeys();
  const defaultModel = defaultModelForKeys(selectedModels, model);
  await testModelConnection({ model, defaultModel, apiKey:'', endpointMode:modelEndpointModeFor(modelProvider(model)), selectedModels });
}

async function testModelConnection(configuration={}) {
  if (state.modelTestLoading) return;
  state.modelTestLoading = true;
  state.testingModelKey = configuration.model || '';
  state.modelTestResult = null;
  renderPage();
  try {
    const result=await apiFetch('/api/v1/settings/model/test',{method:'POST',body:JSON.stringify(configuration)});
    state.modelTestResult=result;
    applyModelTestResultToState(result);
    toast(result.ok?'模型真实调用通过':'模型真实调用未通过',result.message||'模型验证完成',result.ok?'success':'warning');
  } catch(error) {
    state.modelTestResult={ok:false,message:error.message,testedAt:new Date().toISOString(),durationMs:0,model:configuration.model||state.openClawModels?.resolvedDefault||state.openClawModels?.defaultModel||''};
    applyModelTestResultToState(state.modelTestResult);
    toast('模型验证请求失败',error.message,'warning');
  } finally {
    state.modelTestLoading=false;
    state.testingModelKey='';
    await loadOpenClawModels(true);
    applyModelTestResultToState(state.modelTestResult);
    if(state.page==='settings'&&state.settingsTab==='model'){
      if (document.getElementById('modalBackdrop') && !document.getElementById('modalBackdrop').hidden) refreshModelConfigurationModal();
      else renderPage();
    }
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

function schedulerForm(id='') {
  const job=scheduledJobs.find(item=>item.id===id)||{};
  state.formContext={type:'job',id:job.id||''};
  const agentOptions = (state.openClawAgents||[]).filter(agent=>agent.visibility!=='system').map(agent=>agent.id).concat(['sta100-coordinator','sta100-knowledge','market-analyzer']).filter((value,index,array)=>array.indexOf(value)===index);
  openModal({title:job.id?'编辑定时任务':'新增定时任务',eyebrow:'设置 / 定时任务',body:`<div class="form-grid">${inputField('任务名称',job.name||'',true,true,'text','jobName')}${selectField('任务类型',['recommendations','weekly_report','news','index','custom'],true,'jobKind',job.kind||'custom')}${inputField('执行频率 / Cron',job.schedule||'每天 08:00',true,false,'text','jobSchedule')}${inputField('任务说明',job.description||'',true,true,'text','jobDescription')}${selectField('执行 Agent',agentOptions,true,'jobAgent',job.agentId||'sta100-coordinator')}<div class="form-field full"><label for="jobPrompt">执行 Prompt</label><textarea class="textarea" id="jobPrompt" maxlength="8000" rows="5" placeholder="描述任务每次实际需要执行的内容">${escapeHTML(job.prompt||'')}</textarea><small>Prompt 会作为任务定义保存；当前定时调度执行器仍在接入，未执行的任务不会伪装成已完成。</small></div><div class="form-field full"><label style="display:flex;align-items:center;gap:8px"><input id="jobEnabled" type="checkbox" ${job.id?!job.enabled?'':'checked':'checked'}> 启用任务</label></div>${job.id&&!job.builtIn?`<div class="form-field full"><button type="button" class="button danger" data-action="delete-schedule" data-id="${escapeAttr(job.id)}">${icon('trash-2')}删除自定义任务</button></div>`:''}</div>`,footer:formFooter('保存任务','save-schedule')});
}

async function saveSchedule() {
  const existing=scheduledJobs.find(item=>item.id===state.formContext?.id);
  const payload={...(existing||{}),name:formText('jobName'),kind:formText('jobKind'),schedule:formText('jobSchedule'),description:formText('jobDescription'),agentId:formText('jobAgent'),prompt:formText('jobPrompt'),enabled:Boolean(document.getElementById('jobEnabled')?.checked),status:existing?.status||'Ready'};
  if(!payload.name||!payload.kind||!payload.schedule){toast('保存失败','任务名称、类型和执行频率不能为空。','warning');return;}
  try {const record=await apiFetch(existing?'/api/v1/jobs':'/api/v1/jobs',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)});upsertRecord(scheduledJobs,record);closeModal();renderPage();toast(existing?'任务已更新':'任务已创建',record.name);}
  catch(error){toast('任务保存失败',error.message,'warning');}
}

async function deleteSchedule(id) {
  const job=scheduledJobs.find(item=>item.id===id);if(!job||!window.confirm(`确定删除任务“${job.name}”吗？`))return;
  try {await apiFetch(`/api/v1/jobs/${encodeURIComponent(id)}`,{method:'DELETE'});removeRecord(scheduledJobs,id);closeModal();renderPage();toast('任务已删除',job.name);}
  catch(error){toast('任务删除失败',error.message,'warning');}
}

async function generateWeeklyReport() {
  openModal({title:'正在生成智能体周报',eyebrow:'全部智能体 / 最近 7 天',body:`<div class="empty-state">${icon('loader-circle')}<div><h3>正在汇总本机记录</h3><p>读取 Agent 会话和业务审计日志。</p></div></div>`});
  try {const result=await apiFetch('/api/v1/agents/weekly-report',{method:'POST',body:'{}'});document.getElementById('modalBody').innerHTML=`<div class="model-warning"><span>${icon('info')} ${escapeHTML(result.todo||'')}</span></div><pre style="white-space:pre-wrap;line-height:1.7;margin-top:14px">${escapeHTML(result.markdown)}</pre>`;document.getElementById('modalFooter').innerHTML=`<button class="button" data-action="close-modal">关闭</button><button class="button primary" data-action="download-weekly-report">${icon('download')}下载 Markdown</button>`;state.lastWeeklyReport=result;applyIcons();}
  catch(error){closeModal();toast('周报生成失败',error.message,'warning');}
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
async function renderCommandResults(query='') {
  const target=document.getElementById('commandResults');if(!target)return;
  const trimmed=String(query||'').trim();
  if(!trimmed){target.innerHTML=`<div class="command-empty">输入客户、订单、产品或文件关键字</div>`;return;}
  const sequence=++state.commandSearchSeq;
  target.innerHTML=`<div class="command-empty">正在搜索本地业务数据库...</div>`;
  try {
    const data=await apiFetch(`/api/v1/search?q=${encodeURIComponent(trimmed)}`);
    if(sequence!==state.commandSearchSeq)return;
    const iconFor={客户:'building-2',订单:'package-check',产品:'boxes',文件:'file'};
    target.innerHTML=`<div class="command-results">${data.items.length?data.items.map(entry=>`<button class="command-result" data-command-page="${escapeAttr(entry.page)}"><span class="result-icon">${icon(iconFor[entry.type]||'search')}</span><div><strong>${escapeHTML(entry.label)}</strong><small>${escapeHTML(entry.type)} · ${escapeHTML(entry.sub)}</small></div>${icon('arrow-right')}</button>`).join(''):`<div class="command-empty">未找到匹配结果</div>`}</div>`;
    applyIcons();
  } catch(error) {
    if(sequence===state.commandSearchSeq)target.innerHTML=`<div class="command-empty">${escapeHTML(error.message)}</div>`;
  }
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
  if(e.target.id==='modelFamilySelect'){updateModelFamilySelection();return;}
  if(e.target.id==='modelVersionSelect'){updateModelSelection();return;}
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
  if(!el)return;
  const action=el.dataset.action;
  const actions={
    'close-modal':closeModal,'close-drawer':closeDrawer,
    'open-sidebar':()=>document.getElementById('sidebar').classList.add('open'),
    'close-sidebar':()=>document.getElementById('sidebar').classList.remove('open'),
    'sort-table':()=>toggleTableSort(el.dataset.module,el.dataset.field),
    'refresh':()=>void refreshBusinessData(),'pagination-current':()=>toast('当前已是第 1 页','当前数据量只有一页。'),
    'notifications':openNotifications,
    'token-usage':openTokenUsage,
    'refresh-token-usage':async()=>{await loadTokenUsage();openTokenUsage();},
    'clear-token-usage':()=>void clearTokenUsage(),
    'lock':logoutUser,
    'metric-detail':()=>showMetric(el.dataset.key),
    'toggle-recommendations':()=>{state.recExpanded=!state.recExpanded;renderPage();},
    'recommend-detail':()=>newsDetail(recommendations[Number(el.dataset.index)].title),
    'recommend-settings':()=>{setPage('news');setTimeout(()=>document.querySelector('[data-action="news-sources"]')?.click(),0);},
    'oem-preset':()=>{state.oemQuery=el.dataset.value;const input=document.getElementById('oemQuery');if(input)input.value=state.oemQuery;renderPage();},
    'oem-run':()=>void runOEMMatch(),
    'oem-export':oemExport,
    'unified-customer-search':()=>void runUnifiedCustomerSearch(),
    'unified-customer-detail':()=>unifiedCustomerDetail(el.dataset.name),
    'local-discovery-search':()=>void runLocalDiscovery(),
    'local-lead-detail':()=>localLeadDetail(el.dataset.name),
    'agent-chat':()=>showAgentChat(Number(el.dataset.agent),el.dataset.prompt||''),
    'weekly-report':()=>void generateWeeklyReport(),
    'download-weekly-report':()=>state.lastWeeklyReport&&downloadText(`STA100-Agent-Weekly-${new Date().toISOString().slice(0,10)}.md`,state.lastWeeklyReport.markdown),
    'agent-manage':openAgentManager,
    'sync-openclaw-agents':()=>syncOpenClawAgents(el),
    'new-customer':()=>newCustomerForm(), 'edit-customer':()=>newCustomerForm(customers.find(c=>c.id===el.dataset.id)), 'customer-detail':()=>void customerDetail(el.dataset.id), 'delete-customer':()=>deleteCustomer(el.dataset.id), 'customer-more':()=>void customerDetail(el.dataset.id,'activity'),
    'customer-communications':()=>{closeModal();void customerDetail(el.dataset.id,'activity');},'new-customer-communication':()=>customerCommunicationForm(el.dataset.id),'save-customer-communication':()=>void saveCustomerCommunication(el.dataset.id),'cancel-customer-communication':()=>{closeModal();void customerDetail(el.dataset.id,'activity');},
    'export-customers':()=>window.location.assign('/api/v1/accounts/export'),
    'column-settings':openCustomerColumnSettings,'save-customer-columns':saveCustomerColumns,
    'new-quote':()=>newQuoteForm(null,el.dataset.customer||''),'quote-detail':()=>quoteDetail(el.dataset.id),'edit-quote':()=>newQuoteForm(quotes.find(q=>q.id===el.dataset.id)),'delete-quote':()=>void deleteQuote(el.dataset.id),'download-quote':()=>void downloadQuote(el.dataset.id),'convert-order':()=>void convertQuoteToOrder(el.dataset.id),
    'quote-metric-filter':()=>{state.quoteStatus=el.dataset.status||'all';renderPage();},
    'new-order':()=>newOrderForm(),'edit-order':()=>newOrderForm(orders.find(o=>o.id===el.dataset.id)),'delete-order':()=>deleteOrder(el.dataset.id),'order-detail':()=>orderDetail(el.dataset.id),
    'generate-docs':()=>generateDocs(el.dataset.id),'new-document':()=>generateDocs(''),'edit-document':()=>generateDocs('',documents.find(d=>d.id===el.dataset.id)),'delete-document':()=>deleteDocument(el.dataset.id),'download-document':()=>void downloadDocument(el.dataset.id),'template-center':()=>templateCenter(el.dataset.kind),'document-detail':()=>documentDetail(el.dataset.id),
    'clear-document-filters':()=>{state.documentSearch='';state.documentType='all';state.documentStatus='all';renderPage();},
    'upload-template-image':()=>document.getElementById('templateImageInput')?.click(),
    'upload-template-file':()=>document.getElementById('templateFileInput')?.click(),
    'new-product':()=>newProductForm(),'edit-product':()=>newProductForm(products.find(p=>p.id===el.dataset.id)),'product-detail':()=>productDetail(el.dataset.id),'delete-product':()=>deleteProduct(el.dataset.id),'save-product':saveProduct,'import-products':productImportModal,'choose-product-import':()=>document.getElementById('productImportInput')?.click(),'toggle-product-sort':()=>{state.productSort=state.productSort==='stockAsc'?'stockDesc':'stockAsc';renderPage();},
    'new-supplier':()=>newSupplierForm(),'edit-supplier':()=>newSupplierForm(suppliers.find(s=>s.id===el.dataset.id)),'supplier-detail':()=>supplierDetail(el.dataset.id),'delete-supplier':()=>deleteSupplier(el.dataset.id),'export-suppliers':()=>window.location.assign('/api/v1/suppliers/export'),
    'upload-file':uploadFileModal,'choose-file':()=>document.getElementById('privateFileInput')?.click(),'upload-private-file':()=>void uploadPrivateFile(),'save-file-metadata':()=>void saveFileMetadata(),'file-preview':()=>filePreview(el.dataset.id),'file-download':()=>fileDownload(el.dataset.id),'file-edit':()=>fileEdit(el.dataset.id),'file-summary':()=>void fileSummary(el.dataset.id),'file-more':()=>openModal({title:'文件更多操作',eyebrow:'数据库 / 文件操作',body:`<div class="filter-row"><button class="button" data-action="file-summary" data-id="${escapeAttr(el.dataset.id)}">${icon('sparkles')}摘要状态</button><button class="button" data-action="file-download" data-id="${escapeAttr(el.dataset.id)}">${icon('download')}下载文件</button><button class="button" data-action="file-reindex" data-id="${escapeAttr(el.dataset.id)}">${icon('refresh-cw')}重新索引</button><button class="button danger" data-action="file-archive" data-id="${escapeAttr(el.dataset.id)}">${icon('trash-2')}删除文件</button></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`}),'file-archive':()=>void archivePrivateFile(el.dataset.id),'file-reindex':()=>void reindexPrivateFile(el.dataset.id),
    'open-category':()=>{state.fileSearch=el.dataset.category;renderPage();},
    'agent-backup':()=>void backupAgents(),
    'tag-manage':()=>openModal({title:'当前文件标签',eyebrow:'数据库',body:`<div class="filter-row">${[...new Set(files.flatMap(file=>file.tags||[]))].map(v=>`<span class="filter-chip active">${escapeHTML(v)}</span>`).join('')||'<span class="secondary-text">暂无标签</span>'}</div><p class="secondary-text" style="margin-top:14px">标签通过每个文件的“编辑信息”维护，修改后立即保存到本机数据库。</p>`,footer:`<button class="button" data-action="close-modal">关闭</button>`}),
    'news-detail':()=>newsDetail(el.dataset.title),'toggle-news':()=>{state.newsExpanded=!state.newsExpanded;renderPage();},'news-filter':()=>{state.newsCategory=el.dataset.category||'全部';state.newsExpanded=false;renderPage();},
    'news-source-link':()=>openNewsSource(el.dataset.title),'news-todo':()=>void createNewsTodo(el.dataset.title),
    'refresh-news':()=>void refreshNews(),
    'news-sources':openNewsSettings,
    'save-news-settings':()=>void saveNewsSettings(),
    'account-settings':openAccountSettings,
    'model-draft-new':startNewModelConfiguration,
    'edit-model-config':()=>editModelConfiguration(el.dataset.model),
    'delete-model-config':()=>void deleteModelConfiguration(el.dataset.model),
    'test-configured-model':()=>void testConfiguredModel(el.dataset.model),
    'cancel-model-edit':closeModal,
    'save-model-config':()=>saveModelConfiguration(el),
    'test-model-config':()=>void testCurrentModelConfiguration(),
    'select-default-model':()=>void selectDefaultModel(el.dataset.model),
    'save-account-settings':()=>void saveAccountSettings(),
    'toggle-api-key-input':()=>{const input=document.getElementById('modelAPIKey');if(input){input.type=input.type==='password'?'text':'password';el.innerHTML=icon(input.type==='password'?'eye':'eye-off');applyIcons();}},
    'refresh-openclaw-models':async()=>{await loadOpenClawModels(true);toast('模型信息已刷新',state.openClawModels?.error||`当前默认模型：${state.openClawModels?.resolvedDefault||'未配置'}`,state.openClawModels?.error?'warning':'success');},
    'refresh-openclaw-channels':async()=>{await loadOpenClawChannels(true);toast('通道清单已刷新',`当前支持 ${state.openClawChannels?.length||0} 个 OpenClaw 通道。`);},
    'refresh-openclaw-system':async()=>{await Promise.all([loadOpenClawStatus(true),loadOpenClawAgents(true),loadSystemHealth(true)]);toast('系统状态已刷新',state.systemHealth?.status==='ok'?'Go 服务、SQLite 与 OpenClaw 状态正常。':'部分组件需要检查。',state.systemHealth?.status==='ok'?'success':'warning');},
	    'open-channel-binding':()=>openChannelBinding(el.dataset.channel),
	    'channel-status':()=>document.getElementById('channelStatusBox')?void refreshChannelStatus(el.dataset.channel):openChannelBinding(el.dataset.channel),
	    'save-channel-binding':()=>void saveChannelBinding(el.dataset.channel),
	    'new-schedule':()=>schedulerForm(),'edit-schedule':()=>schedulerForm(el.dataset.id),'save-schedule':()=>void saveSchedule(),'delete-schedule':()=>void deleteSchedule(el.dataset.id),'choose-backup':()=>toast('备份目录待部署确认','浏览器不能直接授权后端写入任意外置路径；需确定盒子挂载点和目录白名单。','warning'),
    'offline-upgrade':offlineUpgradeModal,'choose-upgrade-package':()=>document.getElementById('upgradeFileInput')?.click(),'import-upgrade-package':()=>void importOfflineUpgrade(),
    'upgrade-history':()=>void showUpgradeHistory(),
    'relation-select':()=>{const input=document.getElementById(el.dataset.target);if(input){input.value=el.dataset.value;const options=document.getElementById(`${el.dataset.target}Options`);if(options)options.innerHTML='';if(el.dataset.target==='orderQuote')syncOrderFromQuote(el.dataset.value);}},
    'save-customer':()=>void saveCustomer(),'save-quote':()=>void saveQuote(),'save-order':()=>void saveOrder(),'save-document':()=>void saveDocument(),'save-supplier':()=>void saveSupplier(),
    'add-quote-line':()=>{const product=products.find(item=>item.status==='Active')||products[0];state.quoteDraftLines.push({productId:product.id,quantity:1,unitPrice:moneyNumber(product.price),discount:0});renderQuoteDraftLines();},
    'remove-quote-line':()=>{if(state.quoteDraftLines.length===1){toast('至少保留一条产品明细','正式报价单需要至少一个产品。','warning');return;}state.quoteDraftLines.splice(Number(el.dataset.index),1);renderQuoteDraftLines();},
    'add-order-line':()=>{const product=products.find(item=>item.status==='Active')||products[0];state.orderDraftLines.push({productId:product.id,quantity:1,unitPrice:moneyNumber(product.price)});renderOrderDraftLines();},
    'remove-order-line':()=>{if(state.orderDraftLines.length===1){toast('至少保留一条产品明细','订单需要至少一个产品。','warning');return;}state.orderDraftLines.splice(Number(el.dataset.index),1);renderOrderDraftLines();},
    'template-default':()=>void templateAction('default'),'template-edit':()=>void templateAction('edit'),
    'quote-date-filter':()=>dateFilterForm('quote'),'order-date-filter':()=>dateFilterForm('order'),'apply-date-filter':()=>applyDateFilter(el.dataset.module,el.dataset.clear==='true'),
    'send-chat':()=>void sendAgentMessage(Number(el.dataset.agent)),
    'choose-chat-image':()=>document.getElementById('chatImageInput')?.click(),
    'choose-chat-file':()=>document.getElementById('chatFileInput')?.click(),
    'remove-chat-attachment':()=>removeChatAttachment(el.dataset.agent,el.dataset.index),
    'retry-chat':()=>void sendAgentMessage(Number(el.dataset.agent),el.dataset.message||''),
    'chat-quick-prompt':()=>void sendAgentMessage(Number(el.dataset.agent),el.dataset.prompt||''),
    'agent-allowlist':()=>openAgentAllowlist(Number(el.dataset.agent)),
    'save-agent-allowlist':()=>void saveAgentAllowlist(Number(el.dataset.agent)),
    'mock-ocr':customerOCRModal,'choose-customer-ocr':()=>document.getElementById('customerOCRInput')?.click(),
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
