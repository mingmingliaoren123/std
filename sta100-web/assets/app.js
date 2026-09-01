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
  leadSearch: '',
  leadType: 'all',
  leadCountry: 'all',
  leadStatus: 'all',
  leadSort: { field: 'updated', direction: 'desc' },
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
  newsRefreshLoading: false,
  recommendRefreshLoading: false,
  recommendationShowLimit: 5,
  newsFetchLimit: 5,
  newsShowLimit: 20,
  newsFrequency: '1小时',
  newsCountries: '德国、法国、波兰、瑞典',
  newsTopics: 'E-bike、智能骑行、经销商、欧盟法规',
  newsSources: 'EUR-Lex\nBike Europe\nCycling Industry News\nEurobike',
  newsMediaCatalog: [],
  newsMediaCategories: [],
  newsMediaIDs: [],
  newsCustomSources: '',
  newsMediaFilterCategories: [],
  settingsTab: 'model',
  modelConfigured: true,
  oemCategories: [],
  oemCategory: '全部骑行类目',
  oemQuery: 'E-bike 电池 OEM 100 组',
  customerSearchMode: 'local',
  customerSearchQuery: '德国 E-bike 经销商',
  customerSearchInternet: false,
  customerHasContact: true,
  discoveryCountry: '德国',
  discoveryCities: ['柏林'],
  discoveryTypes: ['Distributor'],
  discoveryShowLimit: 10,
  openClawStatus: null,
  openClawStatusLoading: false,
  openClawModels: null,
  openClawModelsLoading: false,
  openClawChannels: null,
  openClawChannelsLoading: false,
  openClawChannelsError: '',
  openClawJobs: null,
  openClawJobsLoading: false,
  openClawJobsError: '',
  openClawCronStatus: null,
  overviewAutomation: null,
  channelActionLoading: {},
  scheduleActionLoading: {},
  channelQR: null,
  channelQRPollTimer: null,
  modelProviderSelection: '',
  modelSearch: '',
  modelDraftMode: 'create',
  modelDraftFamilyKey: '',
  modelDraftKey: '',
  modelDraftOriginalKey: '',
  channelSearch: '',
  channelSkillRoutes: [],
  channelSkillDefinitions: [],
  channelSkillSessions: [],
  channelSkillLoading: false,
  channelSkillError: '',
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
  agentChatPending: {},
  agentChatControllers: {},
  activeAgentChatID: '',
  chatAttachmentsByAgent: {},
  chatAttachments: [],
  tasksLoading: false,
  tasksLoaded: false,
  showApiKey: false,
  formContext: null,
  selectedRows: {
    customers: new Set(),
    quotes: new Set(),
    orders: new Set(),
  },
  customerVisibleColumns: new Set(['customer','type','country','city','contact','orders','total','rating','updated']),
  customerCommunications: {},
  leadCommunications: {},
  supplierCommunications: {},
  quoteDraftLines: [],
  orderDraftLines: [],
  templateUploads: [],
  templates: [],
  templatesLoading: false,
  templateKind: 'document',
  templateSort: { field: 'updated', direction: 'desc' },
  templatePage: 1,
  templatePageSize: 10,
  selectedUploadFile: null,
  selectedUpgradeFile: null,
  systemHealth: null,
  systemHealthLoading: false,
  lastAgentBackup: null,
  overviewDataStatus: '',
  lastWeeklyReport: null,
  assistantResults: {},
  overviewPending: {},
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
  ['agent-knowledge-context','Agent 知识库'],
  ['attachments','附件处理'],
  ['target-agent','当前智能体'],
];
const chatProgressTimers = {};

const pageMeta = {
  overview: ['概览', 'STA-100 / 工作台', '📊'],
  agents: ['智能体', 'STA-100 / 智能协作', '🤖'],
  customers: ['客户', 'STA-100 / 业务管理', '👥'],
  leads: ['线索', 'STA-100 / 业务管理', '🎯'],
  quotes: ['报价单', 'STA-100 / 业务管理', '📄'],
  orders: ['订单', 'STA-100 / 业务管理', '📦'],
  documents: ['单据', 'STA-100 / 外贸单据', '🧾'],
  products: ['产品库', 'STA-100 / 业务管理', '🚲'],
  suppliers: ['供应商', 'STA-100 / 业务管理', '🏭'],
  database: ['数据库', 'STA-100 / 私有知识', '📚'],
  news: ['新闻设置', 'STA-100 / 行业情报', '📰'],
  settings: ['设置', 'STA-100 / 系统管理', '⚙️'],
};

const translations = {
  zh: {
    productSubtitle: '骑行行业智能工作台', navOverview: '概览', navAgents: '智能体', navCustomers: '客户', navLeads: '线索',
    navQuotes: '报价单', navOrders: '订单', navDocuments: '单据', navProducts: '产品库', navSuppliers: '供应商', navDatabase: '数据库',
    navNews: '行业新闻', navSettings: '设置', serverOnline: '服务在线', serverRegion: '欧洲节点', localAccount: '本机账户',
  },
  en: {
    productSubtitle: 'Cycling Industry Workspace', navOverview: 'Overview', navAgents: 'Agents', navCustomers: 'Customers', navLeads: 'Leads',
    navQuotes: 'Quotes', navOrders: 'Orders', navDocuments: 'Documents', navProducts: 'Products', navSuppliers: 'Suppliers', navDatabase: 'Database',
    navNews: 'Industry News', navSettings: 'Settings', serverOnline: 'Service online', serverRegion: 'Europe node', localAccount: 'Local account',
  },
};

const metrics = [
  { key: 'tasks', label: '今日待办', value: 0, icon: 'list-checks', detail: '等待业务数据接入' },
  { key: 'meetings', label: '今日会议', value: 0, icon: 'calendar-clock', detail: '等待日程数据接入' },
  { key: 'documents', label: '今日处理文档', value: 0, icon: 'file-check-2', detail: '等待文档数据接入' },
  { key: 'orders', label: '进行中订单', value: 0, icon: 'package-open', detail: '订单模块实时聚合' },
  { key: 'chats', label: '今日对话', value: 0, icon: 'messages-square', detail: '等待智能体会话数据' },
  { key: 'news', label: '行业资讯', value: 0, icon: 'radio-tower', detail: '等待行业新闻数据' },
];

const recommendations = [];
const tasks = [];

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
const agentFallbackIcons = ['briefcase-business','landmark','ship','file-spreadsheet','shield-check','scale','mail-plus','chart-no-axes-combined','warehouse','package-x','store','refresh-cw','wrench','combine','chart-spline','map','calendar-search','trophy','factory','badge-check','scan-search','book-open-check','pen-tool','route'];
const defaultInternetAllowlist = ['eur-lex.europa.eu','bike-eu.com','cyclingindustry.news','eurobike.com'];

const customers = [];

const leads = [];

const quotes = [];

const orders = [];

const documents = [];

const products = [];

const suppliers = [];

const files = [];

const scheduledJobs = [];

const news = [];

// 业务内容只以登录后 bootstrap 返回的数据为准，接口失败时保持空状态。
[recommendations, customers, quotes, orders, documents, products, suppliers, files, news].forEach(records => records.splice(0));
metrics.forEach(metric => { metric.value = 0; });

const unifiedSearchCustomers = [];
const localDiscoveryLeads = [];

function stableKeyPart(value) {
  const text = String(value || '').trim().toLowerCase();
  if (!text) return 'all';
  let hash = 2166136261;
  for (let i = 0; i < text.length; i++) {
    hash ^= text.charCodeAt(i);
    hash = Math.imul(hash, 16777619);
  }
  return `k${(hash >>> 0).toString(16)}`;
}

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

// 客户发现使用静态地域目录，避免每次打开页面都依赖网络查询国家和城市。
Object.assign(discoveryCities, {
  全球: ['全球'],
  阿富汗: ['喀布尔','赫拉特','坎大哈'], 阿尔及利亚: ['阿尔及尔','奥兰','君士坦丁'],
  安哥拉: ['罗安达','本格拉','洛比托'], 安提瓜和巴布达: ['圣约翰'], 阿根廷: ['布宜诺斯艾利斯','科尔多瓦','罗萨里奥','门多萨'],
  澳大利亚: ['堪培拉','悉尼','墨尔本','布里斯班','珀斯','阿德莱德'], 巴哈马: ['拿骚','弗里波特'],
  巴林: ['麦纳麦','里法'], 孟加拉国: ['达卡','吉大港','库尔纳'], 巴巴多斯: ['布里奇敦'],
  伯利兹: ['贝尔莫潘','伯利兹城'], 贝宁: ['波多诺伏','科托努'], 不丹: ['廷布','帕罗'],
  玻利维亚: ['苏克雷','拉巴斯','圣克鲁斯'], 博茨瓦纳: ['哈博罗内','弗朗西斯敦'],
  巴西: ['巴西利亚','圣保罗','里约热内卢','贝洛奥里藏特','库里蒂巴','阿雷格里港'],
  文莱: ['斯里巴加湾市'], 布基纳法索: ['瓦加杜古','博博迪乌拉索'], 布隆迪: ['基特加','布琼布拉'],
  佛得角: ['普拉亚','明德卢'], 柬埔寨: ['金边','暹粒','西哈努克市'], 喀麦隆: ['雅温得','杜阿拉'],
  加拿大: ['渥太华','多伦多','蒙特利尔','温哥华','卡尔加里','埃德蒙顿'], 中非共和国: ['班吉'],
  乍得: ['恩贾梅纳','蒙杜'], 智利: ['圣地亚哥','瓦尔帕莱索','康塞普西翁'], 哥伦比亚: ['波哥大','麦德林','卡利','巴兰基亚'],
  科摩罗: ['莫罗尼'], 刚果共和国: ['布拉柴维尔','黑角'], 刚果民主共和国: ['金沙萨','卢本巴希','戈马'],
  哥斯达黎加: ['圣何塞','利蒙'], 科特迪瓦: ['亚穆苏克罗','阿比让'], 古巴: ['哈瓦那','圣地亚哥'],
  吉布提: ['吉布提市'], 多米尼克: ['罗索'], 多米尼加: ['圣多明各','圣地亚哥'],
  厄瓜多尔: ['基多','瓜亚基尔','昆卡'], 埃及: ['开罗','亚历山大','吉萨'], 萨尔瓦多: ['圣萨尔瓦多','圣安娜'],
  赤道几内亚: ['马拉博','巴塔'], 厄立特里亚: ['阿斯马拉','马萨瓦'], 埃斯瓦蒂尼: ['姆巴巴内','曼齐尼'],
  埃塞俄比亚: ['亚的斯亚贝巴','德雷达瓦'], 斐济: ['苏瓦','楠迪'], 加蓬: ['利伯维尔','让蒂尔港'],
  冈比亚: ['班珠尔','塞雷昆达'], 加纳: ['阿克拉','库马西','塔马利'], 格林纳达: ['圣乔治'],
  危地马拉: ['危地马拉城','克萨尔特南戈'], 几内亚: ['科纳克里','坎康'], 几内亚比绍: ['比绍'],
  圭亚那: ['乔治敦'], 海地: ['太子港','海地角'], 洪都拉斯: ['特古西加尔巴','圣佩德罗苏拉'],
  印度: ['新德里','孟买','班加罗尔','海得拉巴','金奈','加尔各答'], 印度尼西亚: ['雅加达','泗水','万隆','棉兰','巴厘巴板'],
  伊朗: ['德黑兰','马什哈德','伊斯法罕','设拉子'], 伊拉克: ['巴格达','巴士拉','摩苏尔'],
  以色列: ['耶路撒冷','特拉维夫','海法'], 牙买加: ['金斯敦','蒙特哥贝'], 日本: ['东京','大阪','名古屋','横滨','京都','福冈'],
  约旦: ['安曼','亚喀巴','伊尔比德'], 肯尼亚: ['内罗毕','蒙巴萨','基苏木'], 基里巴斯: ['南塔拉瓦'],
  朝鲜: ['平壤','清津'], 韩国: ['首尔','釜山','仁川','大邱','大田'], 科威特: ['科威特城','艾哈迈迪'],
  吉尔吉斯斯坦: ['比什凯克','奥什'], 老挝: ['万象','琅勃拉邦','巴色'], 黎巴嫩: ['贝鲁特','的黎波里'],
  莱索托: ['马塞卢','马费滕'], 利比里亚: ['蒙罗维亚','布坎南'], 利比亚: ['的黎波里','班加西'],
  马达加斯加: ['塔那那利佛','图阿马西纳'], 马拉维: ['利隆圭','布兰太尔'], 马来西亚: ['吉隆坡','乔治市','新山','怡保','马六甲'],
  马尔代夫: ['马累'], 马里: ['巴马科','廷巴克图'], 毛里塔尼亚: ['努瓦克肖特','努瓦迪布'],
  毛里求斯: ['路易港','居尔皮普'], 墨西哥: ['墨西哥城','瓜达拉哈拉','蒙特雷','普埃布拉','蒂华纳'],
  密克罗尼西亚: ['帕利基尔','科洛尼亚'], 蒙古: ['乌兰巴托','额尔登特'], 摩洛哥: ['拉巴特','卡萨布兰卡','马拉喀什','丹吉尔'],
  莫桑比克: ['马普托','贝拉','楠普拉'], 缅甸: ['内比都','仰光','曼德勒'], 纳米比亚: ['温得和克','鲸湾港'],
  瑙鲁: ['亚伦'], 尼泊尔: ['加德满都','博卡拉'], 新西兰: ['惠灵顿','奥克兰','基督城','汉密尔顿'],
  尼加拉瓜: ['马那瓜','莱昂'], 尼日尔: ['尼亚美','津德尔'], 尼日利亚: ['阿布贾','拉各斯','伊巴丹','卡诺','哈科特港'],
  阿曼: ['马斯喀特','塞拉莱'], 巴基斯坦: ['伊斯兰堡','卡拉奇','拉合尔','费萨拉巴德','白沙瓦'],
  帕劳: ['梅莱凯奥克','科罗尔'], 巴拿马: ['巴拿马城','科隆'], 巴布亚新几内亚: ['莫尔兹比港','莱城'],
  巴拉圭: ['亚松森','东方市'], 秘鲁: ['利马','阿雷基帕','库斯科'], 菲律宾: ['马尼拉','宿务','达沃','奎松市'],
  卡塔尔: ['多哈','赖扬'], 萨摩亚: ['阿皮亚'], 沙特阿拉伯: ['利雅得','吉达','麦加','达曼'],
  塞内加尔: ['达喀尔','图巴'], 塞舌尔: ['维多利亚'], 塞拉利昂: ['弗里敦','博城'],
  新加坡: ['新加坡'], 所罗门群岛: ['霍尼亚拉'], 索马里: ['摩加迪沙','哈尔格萨'],
  南非: ['比勒陀利亚','开普敦','约翰内斯堡','德班','伊丽莎白港'], 南苏丹: ['朱巴'],
  斯里兰卡: ['斯里贾亚瓦德纳普拉科特','科伦坡','康提'], 苏丹: ['喀土穆','恩图曼','苏丹港'],
  苏里南: ['帕拉马里博'], 叙利亚: ['大马士革','阿勒颇','霍姆斯'], 台湾: ['台北','高雄','台中','台南'],
  塔吉克斯坦: ['杜尚别','苦盏'], 坦桑尼亚: ['多多马','达累斯萨拉姆','阿鲁沙'], 泰国: ['曼谷','清迈','芭堤雅','普吉'],
  东帝汶: ['帝力'], 多哥: ['洛美','索科德'], 汤加: ['努库阿洛法'], 特立尼达和多巴哥: ['西班牙港','查瓜拉马斯'],
  突尼斯: ['突尼斯市','苏塞','斯法克斯'], 土库曼斯坦: ['阿什哈巴德','土库曼巴希'], 乌干达: ['坎帕拉','恩德培'],
  阿联酋: ['阿布扎比','迪拜','沙迦','阿治曼'], 美国: ['华盛顿','纽约','洛杉矶','旧金山','西雅图','芝加哥','波士顿','迈阿密'],
  乌拉圭: ['蒙得维的亚','萨尔托'], 乌兹别克斯坦: ['塔什干','撒马尔罕','布哈拉'], 瓦努阿图: ['维拉港','卢甘维尔'],
  委内瑞拉: ['加拉加斯','马拉开波','巴伦西亚'], 越南: ['河内','胡志明市','海防','岘港','芽庄'], 也门: ['萨那','亚丁'],
  赞比亚: ['卢萨卡','恩多拉','基特韦'], 津巴布韦: ['哈拉雷','布拉瓦约'],
  巴勒斯坦: ['拉姆安拉','加沙','伯利恒'], 马绍尔群岛: ['马朱罗','夸贾林'], 圣基茨和尼维斯: ['巴斯特尔'],
  圣卢西亚: ['卡斯特里'], 圣文森特和格林纳丁斯: ['金斯敦'], 圣多美和普林西比: ['圣多美'], 图瓦卢: ['富纳富提'],
});

const customerTypeLabels = {
  Distributor: ['经销商', 'Distributor'], Importer: ['进口商', 'Importer'], Customer: ['客户', 'Customer'],
  Reseller: ['转售商', 'Reseller'], Integrator: ['系统集成商', 'Integrator'], Supplier: ['供应商', 'Supplier'], Other: ['其它', 'Other'],
};
const customerRatingLabels = {
  Prospect: ['潜在客户', 'Prospect'], Active: ['活跃客户', 'Active'], Acquired: ['已成交', 'Acquired'], 'Market Failed': ['市场淘汰', 'Market Failed'],
};
const countryEnglishLabels = {
  中国: 'China', 德国: 'Germany', 法国: 'France', 意大利: 'Italy', 西班牙: 'Spain', 荷兰: 'Netherlands', 瑞典: 'Sweden',
  瑞士: 'Switzerland', 奥地利: 'Austria', 比利时: 'Belgium', 波兰: 'Poland', 丹麦: 'Denmark', 芬兰: 'Finland', 挪威: 'Norway',
  英国: 'United Kingdom', 爱尔兰: 'Ireland', 葡萄牙: 'Portugal', 希腊: 'Greece', 捷克: 'Czechia', 匈牙利: 'Hungary', 罗马尼亚: 'Romania',
  保加利亚: 'Bulgaria', 克罗地亚: 'Croatia', 斯洛伐克: 'Slovakia', 斯洛文尼亚: 'Slovenia', 爱沙尼亚: 'Estonia', 拉脱维亚: 'Latvia',
  立陶宛: 'Lithuania', 卢森堡: 'Luxembourg', 马耳他: 'Malta', 塞浦路斯: 'Cyprus', 冰岛: 'Iceland', 塞尔维亚: 'Serbia', 乌克兰: 'Ukraine',
};
Object.assign(countryEnglishLabels, {
  全球: 'Global', 阿富汗: 'Afghanistan', 阿尔及利亚: 'Algeria', 安哥拉: 'Angola', 安提瓜和巴布达: 'Antigua and Barbuda',
  阿根廷: 'Argentina', 澳大利亚: 'Australia', 巴哈马: 'Bahamas', 巴林: 'Bahrain', 孟加拉国: 'Bangladesh', 巴巴多斯: 'Barbados',
  伯利兹: 'Belize', 贝宁: 'Benin', 不丹: 'Bhutan', 玻利维亚: 'Bolivia', 博茨瓦纳: 'Botswana', 巴西: 'Brazil',
  文莱: 'Brunei', 布基纳法索: 'Burkina Faso', 布隆迪: 'Burundi', 佛得角: 'Cabo Verde', 柬埔寨: 'Cambodia', 喀麦隆: 'Cameroon',
  加拿大: 'Canada', 中非共和国: 'Central African Republic', 乍得: 'Chad', 智利: 'Chile', 哥伦比亚: 'Colombia', 科摩罗: 'Comoros',
  刚果共和国: 'Republic of the Congo', 刚果民主共和国: 'Democratic Republic of the Congo', 哥斯达黎加: 'Costa Rica', 科特迪瓦: 'Côte d’Ivoire',
  古巴: 'Cuba', 吉布提: 'Djibouti', 多米尼克: 'Dominica', 多米尼加: 'Dominican Republic', 厄瓜多尔: 'Ecuador', 埃及: 'Egypt',
  萨尔瓦多: 'El Salvador', 赤道几内亚: 'Equatorial Guinea', 厄立特里亚: 'Eritrea', 埃斯瓦蒂尼: 'Eswatini', 埃塞俄比亚: 'Ethiopia',
  斐济: 'Fiji', 加蓬: 'Gabon', 冈比亚: 'Gambia', 加纳: 'Ghana', 格林纳达: 'Grenada', 危地马拉: 'Guatemala',
  几内亚: 'Guinea', 几内亚比绍: 'Guinea-Bissau', 圭亚那: 'Guyana', 海地: 'Haiti', 洪都拉斯: 'Honduras', 印度: 'India',
  印度尼西亚: 'Indonesia', 伊朗: 'Iran', 伊拉克: 'Iraq', 以色列: 'Israel', 牙买加: 'Jamaica', 日本: 'Japan', 约旦: 'Jordan',
  肯尼亚: 'Kenya', 基里巴斯: 'Kiribati', 朝鲜: 'North Korea', 韩国: 'South Korea', 科威特: 'Kuwait', 吉尔吉斯斯坦: 'Kyrgyzstan',
  老挝: 'Laos', 黎巴嫩: 'Lebanon', 莱索托: 'Lesotho', 利比里亚: 'Liberia', 利比亚: 'Libya', 马达加斯加: 'Madagascar',
  马拉维: 'Malawi', 马来西亚: 'Malaysia', 马尔代夫: 'Maldives', 马里: 'Mali', 毛里塔尼亚: 'Mauritania', 毛里求斯: 'Mauritius',
  墨西哥: 'Mexico', 密克罗尼西亚: 'Micronesia', 蒙古: 'Mongolia', 摩洛哥: 'Morocco', 莫桑比克: 'Mozambique', 缅甸: 'Myanmar',
  纳米比亚: 'Namibia', 瑙鲁: 'Nauru', 尼泊尔: 'Nepal', 新西兰: 'New Zealand', 尼加拉瓜: 'Nicaragua', 尼日尔: 'Niger',
  尼日利亚: 'Nigeria', 阿曼: 'Oman', 巴基斯坦: 'Pakistan', 帕劳: 'Palau', 巴拿马: 'Panama', 巴布亚新几内亚: 'Papua New Guinea',
  巴拉圭: 'Paraguay', 秘鲁: 'Peru', 菲律宾: 'Philippines', 卡塔尔: 'Qatar', 萨摩亚: 'Samoa', 沙特阿拉伯: 'Saudi Arabia',
  塞内加尔: 'Senegal', 塞舌尔: 'Seychelles', 塞拉利昂: 'Sierra Leone', 新加坡: 'Singapore', 所罗门群岛: 'Solomon Islands',
  索马里: 'Somalia', 南非: 'South Africa', 南苏丹: 'South Sudan', 斯里兰卡: 'Sri Lanka', 苏丹: 'Sudan', 苏里南: 'Suriname',
  叙利亚: 'Syria', 台湾: 'Taiwan', 塔吉克斯坦: 'Tajikistan', 坦桑尼亚: 'Tanzania', 泰国: 'Thailand', 东帝汶: 'Timor-Leste',
  多哥: 'Togo', 汤加: 'Tonga', 特立尼达和多巴哥: 'Trinidad and Tobago', 突尼斯: 'Tunisia', 土库曼斯坦: 'Turkmenistan',
  乌干达: 'Uganda', 阿联酋: 'United Arab Emirates', 美国: 'United States', 乌拉圭: 'Uruguay', 乌兹别克斯坦: 'Uzbekistan',
  瓦努阿图: 'Vanuatu', 委内瑞拉: 'Venezuela', 越南: 'Vietnam', 也门: 'Yemen', 赞比亚: 'Zambia', 津巴布韦: 'Zimbabwe',
});

function bilingualLabel(pair) { return state.lang === 'en' ? `${pair[1]} / ${pair[0]}` : `${pair[0]} / ${pair[1]}`; }
function localizedCustomerType(value) { return customerTypeLabels[value] ? bilingualLabel(customerTypeLabels[value]) : value || '未设置'; }
function localizedCustomerRating(value) { return customerRatingLabels[value] ? bilingualLabel(customerRatingLabels[value]) : value || '未设置'; }
function localizedCountry(value) { return value ? `${state.lang === 'en' ? (countryEnglishLabels[value] || value) : value}${state.lang === 'en' && countryEnglishLabels[value] ? ` / ${value}` : state.lang === 'zh' && countryEnglishLabels[value] ? ` / ${countryEnglishLabels[value]}` : ''}` : '未设置'; }
function customerTypeOptions(selected='Customer') { return Object.keys(customerTypeLabels).map(value => ({ value, label: bilingualLabel(customerTypeLabels[value]) })); }
function customerRatingOptions() { return Object.keys(customerRatingLabels).map(value => ({ value, label: bilingualLabel(customerRatingLabels[value]) })); }
function countryOptions(selected='') { return Object.keys(discoveryCities).sort((a,b)=>a.localeCompare(b,'zh-CN')).map(value => ({ value, label: localizedCountry(value) })); }
function preferredAssistantCountries() {
  const out = [];
  if (state.discoveryCountry && typeof state.discoveryCountry === 'string') out.push(state.discoveryCountry);
  if (Array.isArray(state.discoveryCities)) out.push(...state.discoveryCities.filter(Boolean));
  if (Array.isArray(state.discoveryTypes) && state.discoveryTypes.length) out.push(state.discoveryTypes.join(' / '));
  return Array.from(new Set(out.map(v => String(v).trim()).filter(Boolean)));
}
function cityOptions(country, selected='') { return (discoveryCities[country] || []).map(value => ({ value, label: value })); }
function productStatusLabel(status) { return state.lang === 'en' ? ({ Active: 'Active', Review: 'Review', Inactive: 'Inactive' }[status] || status || 'Unset') : ({ Active: '已启用', Review: '待审核', Inactive: '已停用' }[status] || status || '未设置'); }
function refreshCustomerCityOptions(country, selected='') {
  const city = document.getElementById('customerCity');
  if (!city) return;
  const entries = cityOptions(country);
  const nextSelected = selected && entries.some(entry => entry.value === selected) ? selected : (entries[0]?.value || '');
  city.innerHTML = entries.map(entry => `<option value="${escapeAttr(entry.value)}" ${entry.value === nextSelected ? 'selected' : ''}>${escapeHTML(entry.label)}</option>`).join('');
}

function icon(name) { return `<i data-lucide="${name}"></i>`; }
function emojiIcon(emoji, fallbackIcon = 'circle') {
  return `<span class="emoji-stack" aria-hidden="true"><span class="emoji-char">${escapeHTML(emoji)}</span><i class="emoji-fallback-icon" data-lucide="${escapeAttr(fallbackIcon)}"></i></span>`;
}
function escapeAttr(value) { return String(value).replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;'); }
function escapeHTML(value) { return String(value ?? '').replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;').replace(/'/g, '&#39;'); }
function safeExternalURL(value) { try { const url = new URL(String(value || '')); return ['http:','https:'].includes(url.protocol) ? url.href : ''; } catch { return ''; } }
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
  await loadNewsMediaCatalog();
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
  const map = {
    Active: ['green', '活跃', 'Active'], Prospect: ['blue', '潜在', 'Prospect'], Customer: ['green', '客户', 'Customer'],
    Draft: ['neutral', '草稿', 'Draft'], Delivered: ['blue', '已发送', 'Delivered'], Accepted: ['green', '已接受', 'Accepted'],
    Rejected: ['red', '已拒绝', 'Rejected'], Confirmed: ['green', '已确认', 'Confirmed'], Production: ['amber', '生产中', 'Production'],
    Shipped: ['blue', '已发运', 'Shipped'], Completed: ['green', '已完成', 'Completed'], Review: ['amber', '待复核', 'Review'],
    Indexed: ['green', '已索引', 'Indexed'], '已启用': ['green', '已启用', 'Active'], '已停用': ['neutral', '已停用', 'Inactive'],
    '执行中': ['blue', '执行中', 'Running'], '执行成功': ['green', '执行成功', 'Succeeded'], '执行失败': ['red', '执行失败', 'Failed'],
    '等待首次运行': ['amber', '等待首次运行', 'Waiting'], '未同步': ['red', '未同步', 'Unsynced'],
  };
  const [cls, zh, en] = map[status] || ['neutral', status, status];
  const label = state.lang === 'en' ? en : zh;
  return `<span class="badge ${cls}">${label}</span>`;
}

function applyIcons() {
  if (window.lucide) window.lucide.createIcons({ attrs: { 'aria-hidden': 'true' } });
}

function detectEmojiRendering() {
  const canvas = document.createElement('canvas');
  canvas.width = 36;
  canvas.height = 36;
  const context = canvas.getContext('2d', { willReadFrequently: true });
  if (!context) return false;
  context.clearRect(0, 0, canvas.width, canvas.height);
  context.font = '28px "STA100 Emoji", "Noto Color Emoji", "Apple Color Emoji", "Segoe UI Emoji", sans-serif';
  context.textBaseline = 'top';
  context.fillText('📊', 2, 2);
  const data = context.getImageData(0, 0, canvas.width, canvas.height).data;
  let painted = 0;
  let colored = 0;
  for (let index = 0; index < data.length; index += 4) {
    const alpha = data[index + 3];
    if (alpha < 16) continue;
    painted++;
    const r = data[index];
    const g = data[index + 1];
    const b = data[index + 2];
    if (Math.max(r, g, b) - Math.min(r, g, b) > 20) colored++;
  }
  return painted > 20 && colored > 8;
}

function applyEmojiRenderingMode() {
  const supportsEmoji = detectEmojiRendering();
  document.body.classList.toggle('emoji-fallback', !supportsEmoji);
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

async function loadNewsMediaCatalog() {
  if (state.newsMediaCatalog.length) return;
  try {
    const data = await fetch('assets/bicycle-media-sources.json', { headers: { Accept: 'application/json' } }).then(response => {
      if (!response.ok) throw new Error(`HTTP ${response.status}`);
      return response.json();
    });
    state.newsMediaCatalog = Array.isArray(data.items) ? data.items : [];
  } catch (error) {
    state.newsMediaCatalog = [];
    console.warn('news media catalog load failed', error);
  }
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
    replaceRecords(leads, data.leads || []);
    await loadTasks(true);
    replaceRecords(unifiedSearchCustomers, []);
    replaceRecords(localDiscoveryLeads, []);
    const overview = data.overview || {};
    metrics.forEach(metric => {
      if (Number.isFinite(Number(overview[metric.key]))) metric.value = Number(overview[metric.key]);
    });
    state.overviewDataStatus = overview.dataStatus || '';
    state.overviewAutomation = data.automation || null;
    state.oemCategories = Array.isArray(data.oemCategories) ? data.oemCategories : [];
    applyPreferencesPayload(data.preferences || {});
    state.businessDataLoaded = true;
    renderPage();
  } catch (error) {
    toast('业务数据读取失败', error.message, 'warning');
  }
}

async function loadTasks(force = false) {
  if (state.tasksLoading || (!force && state.tasksLoaded)) return tasks;
  state.tasksLoading = true;
  try {
    const data = await apiFetch('/api/v1/tasks');
    replaceRecords(tasks, data.items || []);
    const metric = metrics.find(item => item.key === 'tasks');
    if (metric) {
      metric.value = tasks.filter(item => !['done', 'closed', 'completed'].includes(String(item.status || '').toLowerCase())).length;
    }
  } catch (error) {
    if (force) console.warn('task list load failed', error);
  } finally {
    state.tasksLoading = false;
    state.tasksLoaded = true;
  }
  return tasks;
}

function applyPreferencesPayload(preferences = {}) {
  state.subscription = preferences.recommendationEnabled ?? state.subscription;
  state.recommendationShowLimit = Number.isInteger(preferences.recommendationShowLimit) && preferences.recommendationShowLimit >= 1 && preferences.recommendationShowLimit <= 20 ? preferences.recommendationShowLimit : state.recommendationShowLimit;
  state.discoveryShowLimit = Number.isInteger(preferences.discoveryShowLimit) && preferences.discoveryShowLimit >= 1 && preferences.discoveryShowLimit <= 100 ? preferences.discoveryShowLimit : state.discoveryShowLimit;
  state.discoveryCountry = preferences.discoveryCountry || state.discoveryCountry;
  state.discoveryCities = Array.isArray(preferences.discoveryCities) && preferences.discoveryCities.length ? preferences.discoveryCities : [(discoveryCities[state.discoveryCountry] || [])[0] || ''].filter(Boolean);
  state.discoveryTypes = Array.isArray(preferences.discoveryTypes) && preferences.discoveryTypes.length ? preferences.discoveryTypes : ['Distributor'];
  state.newsFetchLimit = Number.isInteger(preferences.newsFetchLimit) && preferences.newsFetchLimit >= 1 && preferences.newsFetchLimit <= 100 ? preferences.newsFetchLimit : state.newsFetchLimit;
  state.newsShowLimit = preferences.newsShowLimit || state.newsShowLimit;
  state.newsFrequency = preferences.newsFrequency || state.newsFrequency;
  state.newsCountries = preferences.newsCountries || state.newsCountries;
  state.newsTopics = preferences.newsTopics || state.newsTopics;
  state.newsSources = preferences.newsSources || state.newsSources;
  state.newsMediaCategories = Array.isArray(preferences.newsMediaCategories) ? preferences.newsMediaCategories : state.newsMediaCategories;
  state.newsMediaIDs = Array.isArray(preferences.newsMediaIds) ? preferences.newsMediaIds : state.newsMediaIDs;
  state.newsCustomSources = preferences.newsCustomSources ?? state.newsCustomSources;
  if (!state.newsMediaIDs.length && !state.newsCustomSources && state.newsSources) state.newsCustomSources = state.newsSources;
  state.agentInternetAllowlists = preferences.agentAllowlists || state.agentInternetAllowlists;
  state.agentModelSelections = preferences.agentModelOverrides || state.agentModelSelections;
}

function currentPreferences() {
  return { recommendationEnabled:state.subscription, recommendationShowLimit:state.recommendationShowLimit, discoveryShowLimit:state.discoveryShowLimit, discoveryCountry:state.discoveryCountry, discoveryCities:state.discoveryCities, discoveryTypes:state.discoveryTypes, newsFetchLimit:state.newsFetchLimit, newsShowLimit:state.newsShowLimit, newsFrequency:state.newsFrequency, newsCountries:state.newsCountries, newsTopics:state.newsTopics, newsSources:state.newsSources, newsMediaCategories:state.newsMediaCategories, newsMediaIds:state.newsMediaIDs, newsCustomSources:state.newsCustomSources, agentAllowlists:state.agentInternetAllowlists, agentModelOverrides:state.agentModelSelections };
}

async function savePreferences() {
  const response = await apiFetch('/api/v1/settings/preferences',{method:'PATCH',body:JSON.stringify(currentPreferences())});
  applyPreferencesPayload(response.preferences || response || {});
  if (response.automation) state.overviewAutomation = response.automation;
  return response;
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
  if (state.overviewPending.oem) {
    toast('上一次匹配仍在处理', '请等待上个消息完成再进行操作。', 'warning');
    return;
  }
  state.oemQuery=document.getElementById('oemQuery')?.value.trim()||'';
  state.oemCategory=document.getElementById('oemCategory')?.value||state.oemCategory;
  if(!state.oemQuery){toast('请输入匹配需求','产品、数量、市场或规格至少需要一项。','warning');return;}
  state.overviewPending.oem = true;
  renderPage();
  const target=document.getElementById('oemResults');if(target)target.innerHTML=`<div class="tool-empty">${icon('loader-circle')} 正在整理本地证据并调用专业智能体...</div>`;applyIcons();
  try {
    const result=await apiFetch('/api/v1/assistant/query',{method:'POST',body:JSON.stringify({page:'overview',feature:'oem-match',message:state.oemQuery,sessionKey:'sta100-overview-oem',context:{category:state.oemCategory==='全部'?'':state.oemCategory}})});
    state.assistantResults.oem=result;applyTokenUsage(result.tokenUsage);renderPage();toast(result.partial?'OEM 匹配返回部分结果':'OEM 匹配完成',`${result.usedAgents.length} 个 Agent 参与处理。`,result.partial?'warning':'success');
  } catch(error) { toast('OEM 匹配失败',error.message,'warning'); }
  finally {
    delete state.overviewPending.oem;
    if (state.page === 'overview') renderPage();
  }
}

async function runUnifiedCustomerSearch() {
  if (state.overviewPending.customerSearch) {
    toast('上一次搜索仍在处理', '请等待上个消息完成再进行操作。', 'warning');
    return;
  }
  state.customerSearchQuery=document.getElementById('unifiedCustomerQuery')?.value.trim()||'';
  state.customerHasContact=Boolean(document.getElementById('hasContactOnly')?.checked);
  state.overviewPending.customerSearch = true;
  renderPage();
  const target=document.getElementById('unifiedCustomerResults');if(target)target.innerHTML=`<div class="tool-empty">${icon('loader-circle')} 正在执行统一智能搜索...</div>`;applyIcons();
  try {
    state.customerSearchInternet=Boolean(document.getElementById('customerSearchInternet')?.checked);
    const result=await apiFetch('/api/v1/assistant/query',{method:'POST',body:JSON.stringify({page:'overview',feature:'customer-search',message:state.customerSearchQuery||'查询客户',sessionKey:'sta100-overview-customer-search',context:{hasContact:state.customerHasContact,includeInternet:state.customerSearchInternet}})});
    replaceRecords(unifiedSearchCustomers, result.items || []);state.assistantResults.customerSearch=result;applyTokenUsage(result.tokenUsage);renderPage();toast(result.partial?'客户搜索返回部分结果':'客户搜索完成',`返回 ${(result.items || []).length} 条记录，${(result.usedAgents || []).length} 个 Agent 参与。`,result.partial?'warning':'success');
  } catch(error) { toast('客户搜索失败',error.message,'warning'); }
  finally {
    delete state.overviewPending.customerSearch;
    if (state.page === 'overview') renderPage();
  }
}

async function runLocalDiscovery() {
  if (state.overviewPending.customerDiscovery) {
    toast('上一次客户发现仍在处理', '请等待上个消息完成再进行操作。', 'warning');
    return;
  }
  state.discoveryCountry=document.getElementById('discoveryCountry')?.value||state.discoveryCountry;
  const selectedCities=(state.discoveryCities||[]).join('、');
  const selectedTypes=(state.discoveryTypes||[]).join('、');
  const limitInput = document.getElementById('discoveryShowLimit');
  const limitText = limitInput?.value.trim() || String(state.discoveryShowLimit);
  const limit = Number(limitText);
  if (!/^\d+$/.test(limitText) || !Number.isInteger(limit) || limit < 1 || limit > 100) {
    limitInput?.setCustomValidity('返回数量必须是 1-100 的整数。');
    limitInput?.reportValidity();
    limitInput?.focus();
    toast('客户发现数量无效','请输入 1-100 的整数。','warning');
    return;
  }
  if (!(state.discoveryCities||[]).length) { toast('客户发现城市无效','请至少选择一个城市。','warning'); return; }
  if (!(state.discoveryTypes||[]).length) { toast('客户发现类型无效','请至少选择一个客户类型。','warning'); return; }
  state.discoveryShowLimit = limit;
  state.overviewPending.customerDiscovery = true;
  renderPage();
  const target=document.getElementById('localDiscoveryResults');if(target)target.innerHTML=`<div class="tool-empty">${icon('loader-circle')} 正在提交筛选条件并调用客户发现 Agent...</div>`;applyIcons();
  try {
    const message=`请根据以下硬性筛选条件发现客户线索：国家=${state.discoveryCountry}；城市=${selectedCities}；客户类型=${selectedTypes}。最多返回 ${state.discoveryShowLimit} 条。只返回同时满足筛选条件的相关客户或线索；如无法核验，请说明来源和时间，不要编造。`;
    const sessionKey = ['sta100-overview-customer-discovery', stableKeyPart(state.discoveryCountry), stableKeyPart(selectedCities), stableKeyPart(selectedTypes), state.discoveryShowLimit].join('-');
    const result=await apiFetch('/api/v1/overview/customer-discovery',{method:'POST',body:JSON.stringify({page:'overview',feature:'customer-discovery',message,sessionKey,context:{country:state.discoveryCountry,cities:state.discoveryCities,types:state.discoveryTypes,limit:state.discoveryShowLimit,hasContact:false,targetAgent:'customer-measurement-agent',strictDiscovery:true,lang:state.lang}})});
    replaceRecords(localDiscoveryLeads,result.items);state.assistantResults.customerDiscovery=result;applyTokenUsage(result.tokenUsage);renderPage();
    const accessIssue = hasDiscoveryAccessIssue(result);
    toast(accessIssue?'客户发现未完成':(result.partial?'客户发现返回部分结果':'客户发现完成'),accessIssue?'OpenClaw 客户发现 Agent 尚未获取到可核验公开来源，当前不能生成可靠客户线索。':`${result.usedAgents.length} 个 Agent 参与处理。`,result.partial?'warning':'success');
  } catch(error) { toast('客户发现失败',error.message,'warning'); }
  finally {
    delete state.overviewPending.customerDiscovery;
    if (state.page === 'overview') renderPage();
  }
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
    state.modelConfigured = isCurrentDefaultModelUsable();
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

async function loadChannelSkillData(force = false) {
  if (state.channelSkillLoading || (!force && state.channelSkillDefinitions.length && state.channelSkillRoutes.length !== undefined)) return;
  state.channelSkillLoading = true;
  state.channelSkillError = '';
  if (state.page === 'settings' && state.settingsTab === 'channels') renderPage();
  try {
    const [skills, routes, sessions] = await Promise.all([
      apiFetch('/api/v1/channel-skill/skills'),
      apiFetch('/api/v1/channel-skill/routes'),
      apiFetch('/api/v1/channel-skill/sessions'),
    ]);
    state.channelSkillDefinitions = Array.isArray(skills?.skills) ? skills.skills : [];
    state.channelSkillRoutes = Array.isArray(routes?.items) ? routes.items : [];
    state.channelSkillSessions = Array.isArray(sessions?.items) ? sessions.items : [];
  } catch (error) {
    state.channelSkillError = error.message || '通道 Skill 数据读取失败';
  } finally {
    state.channelSkillLoading = false;
    if (state.page === 'settings' && state.settingsTab === 'channels') renderPage();
  }
}

async function loadOpenClawJobs(force = false) {
  if (state.openClawJobsLoading || (state.openClawJobs && !force)) return state.openClawJobs;
  state.openClawJobsLoading = true;
  state.openClawJobsError = '';
  if (state.page === 'settings' && state.settingsTab === 'scheduler') renderPage();
  try {
    const data = await apiFetch('/api/v1/jobs/runtime');
    state.openClawJobs = data.jobs || [];
    state.openClawCronStatus = data.openclawStatus || null;
    if (data.automation) state.overviewAutomation = data.automation;
    state.openClawJobsError = data.synced === false ? (data.syncError || 'OpenClaw Cron 未完成同步，任务列表为本地状态。') : '';
    replaceRecords(scheduledJobs, state.openClawJobs);
  } catch (error) {
    state.openClawJobs = [];
    state.openClawJobsError = error.message || '未知错误';
    if (state.page === 'settings' && state.settingsTab === 'scheduler') toast('定时任务读取失败', error.message, 'warning');
  } finally {
    state.openClawJobsLoading = false;
    if (state.page === 'overview' || (state.page === 'settings' && state.settingsTab === 'scheduler')) renderPage();
  }
  return state.openClawJobs;
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
  if (state.page === 'overview') void loadOpenClawJobs();
  if (state.page === 'overview') void loadTasks();
  if (state.page !== 'settings') return;
  if (state.settingsTab === 'model') void loadOpenClawModels();
  if (state.settingsTab === 'channels') {
    void loadOpenClawChannels();
    void loadChannelSkillData();
  }
  if (state.settingsTab === 'scheduler') {
    void loadOpenClawAgents();
    void loadOpenClawJobs();
  }
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
  const renderers = { overview: renderOverview, agents: renderAgents, customers: renderCustomers, leads: renderLeads, quotes: renderQuotes, orders: renderOrders, documents: renderDocuments, products: renderProducts, suppliers: renderSuppliers, database: renderDatabase, news: renderNews, settings: renderSettings };
  document.getElementById('pageRoot').innerHTML = renderers[state.page]();
  applyIcons();
  wirePageSpecific();
}

function automationByKey(key) {
  return (state.overviewAutomation?.items || []).find(item => item.key === key) || {};
}

function channelActionKey(channel, action) {
  return `${String(channel || '').toLowerCase()}::${String(action || '').toLowerCase()}`;
}

function isChannelActionLoading(channel, action) {
  return Boolean(state.channelActionLoading[channelActionKey(channel, action)]);
}

function channelActionLoadingDetail(channel, action) {
  return state.channelActionLoading[channelActionKey(channel, action)]?.detail || '';
}

function setChannelActionLoading(channel, action, detail = '') {
  const key = channelActionKey(channel, action);
  state.channelActionLoading = { ...state.channelActionLoading, [key]: { channel, action, detail, startedAt: Date.now() } };
  renderPage();
}

function clearChannelActionLoading(channel, action) {
  const key = channelActionKey(channel, action);
  if (!state.channelActionLoading[key]) return;
  const next = { ...state.channelActionLoading };
  delete next[key];
  state.channelActionLoading = next;
  renderPage();
}

function scheduleActionKey(id, action) {
  return `${String(id || '')}::${String(action || '')}`;
}

function isScheduleActionLoading(id, action = '') {
  const prefix = `${String(id || '')}::`;
  if (!action) return Object.keys(state.scheduleActionLoading).some(key => key.startsWith(prefix));
  return Boolean(state.scheduleActionLoading[scheduleActionKey(id, action)]);
}

function setScheduleActionLoading(id, action) {
  state.scheduleActionLoading = { ...state.scheduleActionLoading, [scheduleActionKey(id, action)]: true };
  renderPage();
}

function clearScheduleActionLoading(id, action) {
  const key = scheduleActionKey(id, action);
  if (!state.scheduleActionLoading[key]) return;
  const next = { ...state.scheduleActionLoading };
  delete next[key];
  state.scheduleActionLoading = next;
  renderPage();
}

function businessStatusMeta(status, enabled = true) {
  if (enabled === false) return { label: '已暂停', className: 'neutral', detail: '自动更新已停用，页面继续展示已缓存数据。' };
  switch (String(status || '').toLowerCase()) {
    case 'updated':
      return { label: '数据已更新', className: 'green', detail: '任务结果已写入本机业务库。' };
    case 'syncing':
      return { label: '正在整理', className: 'blue', detail: 'OpenClaw 已执行，正在整理可展示数据。' };
    case 'needs_review':
      return { label: '待人工复核', className: 'amber', detail: '任务已返回，但结果未自动入库。' };
    case 'failed':
      return { label: '更新失败', className: 'red', detail: '任务执行或业务写入失败。' };
    default:
      return { label: '等待首次执行', className: 'amber', detail: '任务已配置，等待 OpenClaw 首次触发。' };
  }
}

function overviewContentStatusLabel(item={}) {
  if (item.businessStatus === 'needs_review' && String(item.businessMessage || '').includes('数据块')) {
    return { label: '格式待复核', className: 'amber', detail: 'OpenClaw 已返回内容，但缺少自动入库格式，页面继续展示缓存数据。' };
  }
  return businessStatusMeta(item.businessStatus, item.enabled);
}

function businessOpenClawAgents() {
  return (state.openClawAgents || []).filter(agent => !agent.isDefault && agent.visibility !== 'system');
}

function humanDateTime(value, fallback = '暂无记录') {
  if (!value) return fallback;
  return formatLocalizedDateTime(String(value).replace(' ', 'T'));
}

function overviewTime(value) {
  if (!value) return '暂无时间';
  const text = String(value);
  const relative = text.match(/^(\d+)\s*(分钟|小时|天|周|月)前$/);
  if (relative) {
    if (state.lang === 'en') {
      const unit = { '分钟': 'minute', '小时': 'hour', '天': 'day', '周': 'week', '月': 'month' }[relative[2]];
      return `${relative[1]} ${unit}${Number(relative[1]) > 1 ? 's' : ''} ago`;
    }
    return text;
  }
  if (text === '刚刚') return state.lang === 'en' ? 'just now' : text;
  const date = new Date(text.replace(' ', 'T'));
  if (!Number.isNaN(date.getTime())) return formatLocalizedDateTime(text);
  if (state.lang === 'en') {
    const englishRelative = text.match(/^(\d+)\s+(minute|minutes|hour|hours|day|days|week|weeks|month|months)\s+ago$/i);
    if (englishRelative) return text;
  }
  return text;
}

function renderOverviewAutomationStrip() {
  const automation = state.overviewAutomation || {};
  const overall = businessStatusMeta(automation.status, state.subscription);
  const items = ['recommendations','news'].map(automationByKey).filter(item => item.key);
  return `<section class="automation-strip panel">
    <div class="automation-main">
      <span class="badge ${overall.className}">${overall.label}</span>
      <div><strong>自动更新状态</strong><p>${escapeHTML(automation.message || overall.detail)}</p></div>
    </div>
    <div class="automation-grid">
      ${items.map(item => {
        const meta = businessStatusMeta(item.businessStatus, item.enabled);
        const moreAction = item.key === 'recommendations' ? 'open-overview-recommendations' : 'open-overview-news';
        const dataCount = item.key === 'recommendations' ? visibleRecommendations().length : visibleNewsItems().length;
        return `<div class="automation-item">
          <div class="automation-card-head">
            <span class="badge ${meta.className}">${meta.label}</span>
            <button class="link-button" data-action="${moreAction}">更多</button>
          </div>
          <strong>${escapeHTML(item.label)}</strong>
          <small>${escapeHTML(item.businessMessage || meta.detail)}</small>
          <em>更新：${escapeHTML(humanDateTime(item.businessUpdatedAt, '等待首次更新'))} · ${escapeHTML(dataCount)} 条</em>
        </div>`;
      }).join('')}
    </div>
  </section>`;
}

function renderOverviewContentStatus(item, countLabel) {
  const meta = overviewContentStatusLabel(item);
  return `<div class="overview-status-line">
    <span class="badge ${meta.className}">${meta.label}</span>
    <span>${escapeHTML(countLabel)}</span>
    <span>最近业务更新：${escapeHTML(humanDateTime(item.businessUpdatedAt, '等待首次更新'))}</span>
    <span>下次调度：${escapeHTML(humanDateTime(item.nextRun, item.enabled === false ? '已暂停' : '等待 OpenClaw 计算'))}</span>
  </div>`;
}

function renderOverview() {
  const displayRecommendations = visibleRecommendations();
  const displayNews = visibleNewsItems();
  const recommendationOverviewLimit = 3;
  const recs = state.recExpanded ? displayRecommendations : displayRecommendations.slice(0, recommendationOverviewLimit);
  const recommendAuto = automationByKey('recommendations');
  const newsAuto = automationByKey('news');
  const nextRunText = humanDateTime(state.overviewAutomation?.nextRun || recommendAuto.nextRun || newsAuto.nextRun, state.subscription ? '等待 OpenClaw 计算' : '不会自动更新');
  return `<div class="page-stack">
    <section class="hero-strip panel">
      <div>
        <span class="badge green">📊 日报订阅</span>
        <h2>行业信息和业务进展，按你的关注条件持续更新</h2>
        <p>系统每 ${escapeHTML(state.newsFrequency)} 由 OpenClaw Agent 执行推荐和新闻更新；结果通过结构化校验后才写入本机业务库。</p>
      </div>
      <div class="subscribe-area">
        <div class="subscribe-meta"><strong>${state.subscription ? '订阅已开启' : '订阅已暂停'}</strong><span>${escapeHTML(state.subscription ? `下次调度 ${nextRunText}` : '不会自动更新')}</span></div>
        <label class="toggle"><input id="subscriptionToggle" type="checkbox" ${state.subscription ? 'checked' : ''}><span></span></label>
      </div>
    </section>

    <section class="metric-grid" aria-label="今日业务摘要">
      ${metrics.map(m => `<button class="metric-button" data-action="metric-detail" data-key="${m.key}"><span class="metric-icon">${icon(m.icon)}</span><span><strong class="metric-number">${m.value}</strong><span class="metric-label">${m.label}</span></span></button>`).join('')}
    </section>

    <section class="content-grid">
      <div class="panel">
        <header class="panel-head"><div><h3>为你推荐</h3><p>根据关注条件、用户操作和智能体记录生成</p></div><div class="inline-actions"><button class="link-button" data-action="recommend-settings">推荐设置</button></div></header>
        ${renderOverviewContentStatus(recommendAuto, `当前可查看 ${displayRecommendations.length} 条推荐`)}
        <div class="recommend-sync-row"><button class="button small" data-action="refresh-recommendations" ${state.recommendRefreshLoading?'disabled':''}>${icon(state.recommendRefreshLoading?'loader-circle':'refresh-cw')}${state.recommendRefreshLoading?'同步中':'同步刷新'}</button><span>${escapeHTML(recommendAuto.businessMessage || '按推荐设置调用 OpenClaw 推荐 Agent 获取真实内容。')}</span></div>
        <div class="recommendation-list ${state.recExpanded ? 'recommendation-list-expanded' : ''}">
          ${recs.map((r, i) => `<article class="recommendation-item"><span class="recommendation-rank">${String(i + 1).padStart(2, '0')}</span><div class="recommendation-copy"><button class="recommendation-title" data-action="recommend-detail" data-id="${escapeAttr(r.id)}">${escapeHTML(recommendationDisplayTitle(r))}</button><p class="recommendation-preview">${escapeHTML(recommendationPreview(r))}</p><div class="source-line"><span class="mini-source">来源 <strong>${escapeHTML(r.source || '未返回')}</strong></span><span class="mini-source">类型 <strong>${escapeHTML(r.type || '推荐')}</strong></span><span class="mini-source">${escapeHTML(overviewTime(r.updatedAt || r.time))}</span></div></div><button class="table-icon" data-action="recommend-detail" data-id="${escapeAttr(r.id)}" aria-label="查看详情" title="查看详情">${icon('arrow-up-right')}</button></article>`).join('') || `<div class="empty-state compact-empty"><p>暂无有效推荐。当前任务只返回了说明性内容或没有可靠结果，保留的历史缓存仍可在任务记录中复核。</p></div>`}
        </div>
        ${displayRecommendations.length > recommendationOverviewLimit ? `<div class="panel-body recommendation-more-row"><button class="button ghost small" data-action="toggle-recommendations">${icon(state.recExpanded ? 'chevron-up' : 'chevron-down')}${state.recExpanded ? '收起' : '更多'}</button>${state.recExpanded ? `<span class="secondary-text">可在列表内滑动查看 ${displayRecommendations.length} 条推荐</span>` : ''}</div>` : ''}
      </div>
      <aside class="panel">
        <header class="panel-head"><div><h3>行业新闻</h3><p>默认展示相关度最高的 3 条</p></div><div class="inline-actions"><button class="link-button" data-action="news-sources">新闻设置</button></div></header>
        ${renderOverviewContentStatus(newsAuto, `当前可查看 ${displayNews.length} 条资讯`)}
        <div class="news-mini">${displayNews.slice(0,3).map(n => `<button class="news-mini-item" data-action="news-detail" data-title="${escapeAttr(n.title)}"><span class="news-mini-meta"><span class="badge neutral">${escapeHTML(n.category || '行业资讯')}</span><span class="news-time">${escapeHTML(overviewTime(n.updatedAt || n.time))}</span></span><span><h4>${escapeHTML(cleanVisibleText(n.title))}</h4><p>${escapeHTML(newsPreview(n, 120))}</p></span></button>`).join('') || `<div class="empty-state compact-empty"><p>暂无有效行业新闻。请先执行新闻设置中的同步任务。</p></div>`}</div>
        <div class="news-more-row"><button class="button ghost small" data-page="news">${icon('chevron-down')}更多</button></div>
      </aside>
    </section>

    <div class="section-head tool-section-head"><div><h3>🛠️ 智能业务工具</h3><p>不同工具按自身流程调用 OpenClaw Agent，输出统一展示。</p></div><span class="meta">统一结果 · 冲突信息并列保留</span></div>
    <section class="tool-grid">
      <article class="panel tool-panel tool-panel-wide">
        <header class="tool-header"><div><h3>🏭 OEM 工厂智能匹配</h3><p>用户输入需求后，系统整理已索引证据并调用 OEM 工厂匹配助手返回结果。</p></div></header>
        <div class="panel-body">
          <div class="filter-row tool-presets">${['公路整车 OEM 1000 台','E-bike 电池 OEM 100 组','中置电机 500 套','头盔 MIPS 500 个'].map(v=>`<button class="filter-chip" data-action="oem-preset" data-value="${v}">${v}</button>`).join('')}</div>
          <div class="filter-row oem-category-tags" aria-label="知识库分类">${[{value:'全部',label:'全部已索引分类'},...state.oemCategories].map(category=>`<button class="filter-chip ${state.oemCategory===category.value||(category.value==='全部'&&state.oemCategory==='全部骑行类目')?'active':''}" data-action="oem-category" data-value="${escapeAttr(category.value)}">${escapeHTML(category.label)}${category.documents?` (${category.documents})`:''}</button>`).join('') || '<span class="secondary-text">知识库同步完成后将展示可选分类</span>'}</div>
          <div class="tool-form">
            <label class="field-search tool-query">${icon('search')}<input id="oemQuery" value="${escapeAttr(state.oemQuery)}" placeholder="输入产品、数量、市场和要求"></label>
            <button class="button primary" data-action="oem-run" aria-busy="${state.overviewPending.oem?'true':'false'}">${icon(state.overviewPending.oem?'loader-circle':'scan-search')}${state.overviewPending.oem?'匹配中':'开始匹配'}</button>
          </div>
          <div class="agent-chain-note"><span class="agent-icon">${icon('workflow')}</span><span><strong>本地证据 → OEM 工厂匹配助手 → 统一汇总</strong><small>分类仅展示已索引的 OEM 产品品类；可按标签筛选，输入需求支持模糊匹配；冲突信息全部保留。</small></span></div>
          <div class="tool-results" id="oemResults">${renderOEMCards()}</div>
        </div>
      </article>

      <article class="panel tool-panel">
        <header class="tool-header"><div><h3>🔍 客户统一搜索</h3><p>优先检索本地客户、线索和已索引知识库；可按需补充互联网公开来源。</p></div></header>
        <div class="panel-body">
          <div class="tool-form compact">
            <label class="field-search tool-query">${icon('search')}<input id="unifiedCustomerQuery" value="${escapeAttr(state.customerSearchQuery)}" placeholder="国家、公司、业务、邮箱或电话"></label>
            <label class="contact-check"><input class="checkbox" id="hasContactOnly" type="checkbox" ${state.customerHasContact?'checked':''}> 必有联系方式</label>
            <label class="contact-check"><input class="checkbox" id="customerSearchInternet" type="checkbox" ${state.customerSearchInternet?'checked':''}> 补充互联网</label>
            <button class="button primary" data-action="unified-customer-search" aria-busy="${state.overviewPending.customerSearch?'true':'false'}">${icon(state.overviewPending.customerSearch?'loader-circle':'search')}${state.overviewPending.customerSearch?'搜索中':'搜索'}</button>
          </div>
          <div class="agent-chain-note"><span class="agent-icon">${icon('workflow')}</span><span><strong>用户输入 → 本地客户 / 线索 / 知识库 → 按需补充互联网</strong><small>“代理商、经销商、分销商”会匹配 Dealer、Distributor 等渠道类型；所有来源分别保留。</small></span></div>
          <div class="tool-results customer-result-list unified-customer-result-list" id="unifiedCustomerResults">${renderUnifiedCustomerCards()}</div>
        </div>
      </article>

      <article class="panel tool-panel">
          <header class="tool-header"><div><h3>🌍 本地客户发现</h3><p>选择国家（单选）、城市和客户类型（可多选），OpenClaw 客户发现 Agent 将搜索符合条件的公开来源客户线索。</p></div></header>
          <div class="panel-body">
          <div class="tool-form compact discovery-form">
            <label class="discovery-field"><span>国家</span><select class="select" id="discoveryCountry">${Object.keys(discoveryCities).sort((a,b)=>a.localeCompare(b,'zh-CN')).map(v=>`<option value="${escapeAttr(v)}" ${state.discoveryCountry===v?'selected':''}>${escapeHTML(localizedCountry(v))}</option>`).join('')}</select></label>
            <label class="field-number discovery-limit"><span>返回条数</span><input class="input" id="discoveryShowLimit" type="number" min="1" max="100" step="1" value="${state.discoveryShowLimit}"></label>
            <div class="discovery-picker discovery-picker-cities"><div class="discovery-picker-head"><strong>城市</strong><span>${(state.discoveryCities||[]).length} 个已选</span></div><div class="discovery-chip-grid" id="discoveryCitiesRow">${(discoveryCities[state.discoveryCountry]||[]).map(v=>`<button type="button" class="filter-chip ${(state.discoveryCities||[]).includes(v)?'active':''}" data-discovery-city="${escapeAttr(v)}">${escapeHTML(v)}</button>`).join('')}</div></div>
            <div class="discovery-picker"><div class="discovery-picker-head"><strong>客户类型</strong><span>${(state.discoveryTypes||[]).length} 个已选</span></div><div class="discovery-chip-grid discovery-type-grid" id="discoveryTypesRow">${[['Distributor','经销商'],['Importer','进口商'],['Dealer','车店'],['Brand','品牌'],['OEM','OEM']].map(([v,l])=>`<button type="button" class="filter-chip ${(state.discoveryTypes||[]).includes(v)?'active':''}" data-discovery-type="${v}">${l}</button>`).join('')}</div></div>
            <div class="discovery-actions"><div class="agent-chain-note discovery-agent-note"><span class="agent-icon">${icon('bot')}</span><span><strong>OpenClaw · CustomerMeasurementAgent</strong><small>国家单选，城市和类型多选；无可靠公开来源时不编造结果。</small></span></div><div class="discovery-action-buttons"><button class="button" data-action="save-discovery-settings">${icon('save')}保存默认</button><button class="button primary" data-action="local-discovery-search" aria-busy="${state.overviewPending.customerDiscovery?'true':'false'}">${icon(state.overviewPending.customerDiscovery?'loader-circle':'radar')}${state.overviewPending.customerDiscovery?'发现中':'开始发现'}</button></div></div>
          </div>
          <div class="tool-results customer-result-list" id="localDiscoveryResults">${renderLocalDiscoveryCards()}</div>
        </div>
      </article>
    </section>
  </div>`;
}

async function saveDiscoverySettings() {
  const input = document.getElementById('discoveryShowLimit');
  const text = input?.value.trim() || '';
  const value = Number(text);
  if (!/^\d+$/.test(text) || !Number.isInteger(value) || value < 1 || value > 100) {
    input?.setCustomValidity('返回数量必须是 1-100 的整数。');
    input?.reportValidity();
    input?.focus();
    toast('默认数量无效', '请输入 1-100 的整数。', 'warning');
    return;
  }
  input?.setCustomValidity('');
  state.discoveryShowLimit = value;
  const button = document.querySelector('[data-action="save-discovery-settings"]');
  if (button) { button.disabled = true; button.innerHTML = `${icon('loader-circle')}保存中`; applyIcons(); }
  try {
    await savePreferences();
    renderPage();
    toast('客户发现默认设置已保存', `以后每次默认最多返回 ${value} 条客户线索。`, 'success');
  } catch (error) {
    if (button) { button.disabled = false; button.innerHTML = `${icon('save')}保存默认`; applyIcons(); }
    toast('默认设置保存失败', error.message, 'warning');
  }
}

function renderOEMCards() {
  const result=state.assistantResults.oem;
  if(!result)return `<div class="tool-empty">输入需求后，系统将从已索引的 Agent 知识库整理证据，再调用 OEM 专业 Agent 返回结果。</div>`;
  const rows = Array.isArray(result.items) ? result.items : [];
  const summary = renderAssistantSummary(result,'OEM 匹配分析');
  if (!rows.length) return summary;
  return `${summary}<div class="oem-match-list">${rows.map((item,index)=>`<article class="customer-match-row oem-match-row"><span class="match-score">${escapeHTML(String(item.score || index + 1))}<small>${item.score?'分':'位'}</small></span><span class="customer-match-copy"><strong>${escapeHTML(item.title || item.name || '候选工厂')}</strong><small>${escapeHTML([item.category,item.capacity,item.moq].filter(Boolean).join(' · ') || '能力信息待核实')}</small><em>${icon('factory')} ${escapeHTML(item.reason || item.detail || 'OpenClaw 已返回结构化候选结果')}</em></span><span class="badge blue">${escapeHTML(item.source || '来源待核实')}</span></article>`).join('')}</div>`;
}

function renderAssistantSummary(result,title) {
  const agents=(result.usedAgents||[]).join('、')||'无成功调用';
  const text = cleanVisibleText(result.text || '暂无汇总文本', '暂无汇总文本');
  return `<article class="assistant-summary"><div class="spread"><strong>${escapeHTML(title)}</strong><span class="badge ${result.partial?'amber':'green'}">${result.partial?'部分结果':'已完成'}</span></div><p>${escapeHTML(text).replace(/\n/g,'<br>')}</p><div class="source-line"><span class="mini-source">参与 Agent <strong>${escapeHTML(agents)}</strong></span><span class="mini-source">本地证据 <strong>${result.evidence?.length||0} 条</strong></span><span class="mini-source">结构化结果 <strong>${Array.isArray(result.items)?result.items.length:0} 条</strong></span><span class="mini-source">冲突 <strong>${result.conflicts?.length||0} 组</strong></span></div></article>`;
}

function renderCompactAssistantSummary(result,title,emptyText='暂无明确结果') {
  if (!result) return '';
  const text = compactAssistantText(result.text || emptyText);
  const issue = title.includes('客户发现') && hasDiscoveryAccessIssue(result);
  return `<article class="assistant-summary compact-summary ${issue?'assistant-summary-warning':''}"><div class="spread"><strong>${escapeHTML(title)}</strong><span class="badge ${issue||result.partial?'amber':'green'}">${issue?'未完成':(result.partial?'部分':'完成')}</span></div><p>${escapeHTML(text).replace(/\n/g,'<br>')}</p></article>`;
}

function compactAssistantText(value) {
  const text = cleanVisibleText(value, '暂无明确结果')
    .replace(/请结合证据记录编号和更新时间人工复核。?/g, '')
    .replace(/具体原因可在消息下方阶段结果中查看。?/g, '')
    .replace(/已从本机业务数据库检索到\s*\d+\s*条证据和\s*\d+\s*条可展示记录，?/g, '')
    .replace(/最终协调 Agent 未完成，当前为部分结果。?/g, '')
    .trim();
  if (text.length <= 360) return text;
  return `${text.slice(0, 360)}…`;
}

function renderUnifiedCustomerCards() {
  const rows = unifiedSearchCustomers.filter(r => !state.customerHasContact || Boolean([r.contact,r.phone,r.email,r.website].find(Boolean)));
  const summary=state.assistantResults.customerSearch?renderCompactAssistantSummary(state.assistantResults.customerSearch,'搜索结果'):'';
  return `${summary}${rows.map(r=>{const sourceType=r.sourceType||'internet';const sourceLabel=r.sourceLabel||r.source||'来源待核实';const sourceClass=['local_business','local_lead','local_knowledge'].includes(sourceType)?'green':'blue';const sourceName=sourceType==='local_business'?'本地业务数据':sourceType==='local_lead'?'本地线索库':sourceType==='local_knowledge'?'本地知识库':'互联网公开来源';const sourceURL=safeExternalURL(r.sourceUrl);return `<article class="customer-match-row"><span class="match-score">${escapeHTML(String(r.score||0))}<small>分</small></span><span class="customer-match-copy"><strong>${escapeHTML(r.name||'未命名客户')}</strong><small>${escapeHTML([r.country,r.city,r.type].filter(Boolean).join(' · ')||'区域和类型待核实')} · ${escapeHTML(r.business||'暂无业务描述')}</small><em>${icon('contact-round')} ${escapeHTML([r.contact,r.phone,r.email].filter(Boolean).join(' · ')||'未填写')}</em><span class="source-line"><span class="badge ${sourceClass}">${sourceName}</span><span class="mini-source">${escapeHTML(sourceLabel)}</span>${r.sourceUpdatedAt||r.updatedAt?`<span class="mini-source">更新 ${escapeHTML(overviewTime(r.sourceUpdatedAt||r.updatedAt))}</span>`:''}${sourceURL?`<a class="mini-source" href="${escapeAttr(sourceURL)}" target="_blank" rel="noreferrer">查看原文</a>`:''}</span></span><span class="badge ${sourceClass}" title="${escapeAttr(sourceType==='internet'?'互联网公开来源':'本地来源')}">${sourceType==='internet'?'互联网':'本地'}</span><button class="table-icon" data-action="unified-customer-detail" data-name="${escapeAttr(r.name||'')}" title="查看详情">${icon('arrow-up-right')}</button></article>`;}).join('')}`;
}

function renderLocalDiscoveryCards() {
  const rows=localDiscoveryLeads.slice(0, Math.max(1, Number(state.discoveryShowLimit) || 10));
  const result = state.assistantResults.customerDiscovery;
  const accessIssue = hasDiscoveryAccessIssue(result);
  const failed = hasAssistantFailure(result);
  // 简洁展示：只显示筛选条件和状态，带一键转化按钮
  const addAllBtn = rows.length ? `<button class="button small" data-action="add-all-discovery-leads">${icon('plus')}一键转化全部</button>` : '';
  const filterInfo = `<div class="discovery-filter-info"><div class="discovery-filter-copy"><span class="badge ${accessIssue||failed?'amber':(result?'green':'neutral')}">${accessIssue||failed?'未完成':(result?'已完成':'待查询')}</span><span class="secondary-text">${state.discoveryCountry} · ${(state.discoveryCities||[]).join('、')} · ${(state.discoveryTypes||[]).join('、')}</span><span class="secondary-text">${rows.length} 条结果</span></div>${addAllBtn}</div>`;
  if (accessIssue || failed) return `${filterInfo}<div class="tool-empty tool-empty-warning">${icon('wifi-off')} ${escapeHTML(discoveryFailureText(result, accessIssue))}</div>`;
  if (!result) return `${filterInfo}<div class="tool-empty">选择国家、城市和客户类型后，系统会直接调用 OpenClaw 客户发现 Agent 搜索相关线索。</div>`;
  if (!rows.length) return `${filterInfo}<div class="tool-empty">OpenClaw 客户发现 Agent 已完成检索，但没有返回同时满足当前国家、城市和客户类型条件的可核验客户线索。</div>`;
  const rowsHTML = rows.map(r=>`<article class="customer-match-row discovery-match-row"><span class="match-score">${r.score}<small>分</small></span><span class="customer-match-copy"><strong>${r.name}</strong><small>${r.country} · ${Array.isArray(r.city)?r.city.join('、'):r.city} · ${Array.isArray(r.type)?r.type.join('、'):r.type}</small><em>${icon('contact-round')} ${r.contact||'未填写'}</em></span><span class="match-status-actions"><span class="badge blue">已整合</span><button class="table-icon" data-action="add-to-leads" data-name="${r.name}" title="加入线索库">${icon('plus')}</button></span><button class="table-icon" data-action="local-lead-detail" data-name="${r.name}" title="查看详情">${icon('arrow-up-right')}</button></article>`).join('');
  return `${filterInfo}${rowsHTML}`;
}

function hasAssistantFailure(result) {
  if (!result) return false;
  if (Array.isArray(result.items) && result.items.length) return false;
  return Boolean(result.partial) || (result.pipeline||[]).some(stage => stage.status === 'failed') || (result.agentOutputs||[]).some(output => output.error);
}

function discoveryFailureText(result, accessIssue=false) {
  if (!result) return 'OpenClaw 客户发现 Agent 暂不可用。';
  if (accessIssue) return 'OpenClaw 客户发现 Agent 尚未获取到可核验公开来源。请确认该 Agent 的公开来源能力后，再按国家、城市和客户类型重新发现。';
  const failedStage = (result.pipeline||[]).find(stage => stage.status === 'failed');
  const outputError = (result.agentOutputs||[]).find(output => output.error)?.error;
  return cleanVisibleText(failedStage?.detail || outputError || result.text || 'OpenClaw 客户发现 Agent 调用失败，请先检查模型配置和 OpenClaw 网关状态。', 'OpenClaw 客户发现 Agent 调用失败，请先检查模型配置和 OpenClaw 网关状态。');
}

function hasDiscoveryAccessIssue(result) {
  if (!result) return false;
  const text = [result.text, ...(result.pipeline||[]).map(stage=>`${stage.detail||''} ${stage.reason||''} ${stage.data||''}`), ...(result.agentOutputs||[]).map(output=>`${output.text||''} ${output.error||''}`)].join(' ').toLowerCase();
  return text.includes('公开来源能力不可用') ||
    text.includes('尚未获取到可核验公开来源') ||
    text.includes('联网检索或网页读取工具不可用') ||
    (text.includes('web_search') && (text.includes('disabled') || text.includes('no provider') || text.includes('禁用'))) ||
    (text.includes('web_fetch') && (text.includes('private/internal/special-use') || text.includes('timeout') || text.includes('超时') || text.includes('被拦截'))) ||
    text.includes('未取到任何页面正文');
}

function unifiedCustomerDetail(name) {
  const customer = unifiedSearchCustomers.find(r => r.name === name);
  if (!customer) return;
  const sourceType = String(customer.sourceType || 'internet');
  const sourceURL = safeExternalURL(customer.sourceUrl);
  const sourceName = sourceType === 'local_business' ? '本地业务数据库' : sourceType === 'local_lead' ? '本地线索库' : sourceType === 'local_knowledge' ? '本地知识库' : '互联网公开来源';
  const value = (item, fallback = '未提供') => String(item ?? '').trim() || fallback;
  const fields = [['客户编号', customer.id], ['国家', customer.country], ['城市', customer.city], ['客户类型', customer.type], ['业务方向', customer.business], ['联系人', customer.contact], ['电话', customer.phone], ['邮箱', customer.email], ['官网', customer.website], ['地址', customer.address], ['来源', value(customer.source, sourceName)], ['来源文件', customer.sourcePath], ['更新时间', customer.sourceUpdatedAt || customer.updatedAt], ['数据新鲜度', customer.freshness], ['匹配理由', customer.reason], ['联系方式过滤', state.customerHasContact ? '已启用' : '未启用']];
  const content = value(customer.content || customer.business, '暂无可展示的知识库正文片段');
  openDrawer({ title: value(customer.name, '客户详情'), eyebrow: `客户统一搜索 / ${sourceName}`, body: `<div class="spread"><span class="badge ${sourceType === 'internet' ? 'blue' : 'green'}">${sourceType === 'internet' ? '互联网' : '本地来源'}</span><span class="secondary-text">匹配 ${escapeHTML(String(customer.score ?? '待核实'))} 分</span></div><div class="detail-grid" style="margin-top:15px">${fields.map(([label,item])=>`<div class="detail-field"><label>${label}</label><strong>${escapeHTML(value(item))}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">详情内容</div><div class="recommendation-detail-summary" style="white-space:pre-wrap;line-height:1.75">${escapeHTML(content)}</div><div class="divider-title" style="margin:20px 0 12px">来源信息</div><p class="secondary-text" style="line-height:1.8">${escapeHTML(value(customer.sourceLabel, sourceName))}；更新时间：${escapeHTML(overviewTime(customer.sourceUpdatedAt || customer.updatedAt, '未提供'))}。</p>${sourceURL ? `<div class="inline-actions" style="margin-top:16px"><a class="button" href="${escapeAttr(sourceURL)}" target="_blank" rel="noreferrer">${icon('external-link')}查看原文</a></div>` : ''}` });
}

function localLeadDetail(name) {
  const lead = localDiscoveryLeads.find(item => item.name === name);
  if (!lead) return;
  const sourceButton = lead.sourceUrl ? `<div class="inline-actions" style="margin-top:16px"><button class="button" onclick="window.open('${escapeAttr(lead.sourceUrl)}','_blank','noopener,noreferrer')">${icon('external-link')}查看原文</button></div>` : '';
  openDrawer({ title: lead.name, eyebrow: '本地客户发现 / OpenClaw', body: `<div class="spread"><span class="badge blue">CustomerMeasurementAgent</span><span class="secondary-text">匹配 ${escapeHTML(String(lead.score ?? '待核实'))} 分</span></div><div class="detail-grid" style="margin-top:15px">${[['国家',lead.country],['城市',Array.isArray(lead.city)?lead.city.join('、'):lead.city],['客户类型',Array.isArray(lead.type)?lead.type.join('、'):lead.type],['业务方向',lead.business||'未提供'],['联系人',lead.contact||'未提供'],['电话',lead.phone||'未提供'],['邮箱',lead.email||'未提供'],['官网',lead.website||'未提供'],['地址',lead.address||'未提供'],['来源',lead.source||'未提供'],['更新时间',lead.updatedAt||'未提供'],['匹配理由',lead.reason||'未提供']].map(([label,value])=>`<div class="detail-field"><label>${label}</label><strong>${escapeHTML(String(value||'未提供'))}</strong></div>`).join('')}</div>${sourceButton}<div class="agent-chain-note" style="margin-top:16px"><span class="agent-icon">${icon('bot')}</span><span><strong>筛选条件 → CustomerMeasurementAgent → 结构化客户线索</strong><small>只展示满足国家、城市和客户类型条件的候选客户；无可靠结果时不编造。</small></span></div><div style="margin-top:16px"><button class="button primary" data-action="add-to-leads" data-name="${lead.name}">${icon('plus')}加入线索库</button></div>` });
}

async function addDiscoveryLeadToLeads(name) {
  const lead = localDiscoveryLeads.find(item => item.name === name);
  if (!lead) return;
  const payload = {
    name: lead.name,
    type: Array.isArray(lead.type) ? lead.type[0] : (lead.type || 'Distributor'),
    country: lead.country || '',
    city: Array.isArray(lead.city) ? lead.city.join('、') : (lead.city || ''),
    contact: lead.contact || '',
    phone: lead.phone || '',
    email: lead.email || '',
    website: lead.website || '',
    address: lead.address || '',
    business: lead.business || '',
    source: lead.source || 'OpenClaw 客户发现',
    sourceUrl: lead.sourceUrl || '',
    score: lead.score || 0,
    reason: lead.reason || '',
  };
  try {
    const response = await apiFetch('/api/v1/leads', { method: 'POST', body: JSON.stringify(payload) });
    leads.push(response);
    toast('已加入线索库', `线索 "${lead.name}" 已保存到线索库。`);
    renderPage();
  } catch (error) {
    toast('加入线索库失败', error.message, 'warning');
  }
}

async function addAllDiscoveryLeadsToLeads() {
  if (!localDiscoveryLeads.length) { toast('无可加入线索', '当前没有客户发现结果。', 'warning'); return; }
  const existingNames = new Set(leads.map(l => l.name));
  const newLeads = localDiscoveryLeads.filter(lead => !existingNames.has(lead.name));
  if (!newLeads.length) { toast('已是线索', '所有发现的客户都已在线索库中。', 'warning'); return; }
  let added = 0;
  for (const lead of newLeads) {
    const payload = {
      name: lead.name,
      type: Array.isArray(lead.type) ? lead.type[0] : (lead.type || 'Distributor'),
      country: lead.country || '',
      city: Array.isArray(lead.city) ? lead.city.join('、') : (lead.city || ''),
      contact: lead.contact || '',
      phone: lead.phone || '',
      email: lead.email || '',
      website: lead.website || '',
      address: lead.address || '',
      business: lead.business || '',
      source: lead.source || 'OpenClaw 客户发现',
      sourceUrl: lead.sourceUrl || '',
      score: lead.score || 0,
      reason: lead.reason || '',
    };
    try {
      const response = await apiFetch('/api/v1/leads', { method: 'POST', body: JSON.stringify(payload) });
      leads.push(response);
      added++;
    } catch (error) { /* skip failed ones */ }
  }
  toast('已加入线索库', `已将 ${added} 条客户发现结果保存到线索库。`);
  renderPage();
}

async function oemExport() {
  try {
    await apiFetch('/api/v1/overview/oem-matches/export',{method:'POST',body:JSON.stringify({query:state.oemQuery,category:state.oemCategory})});
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
  const managedAgents = businessOpenClawAgents();
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
      ${filtered.map(({ agent: a, index }) => { const live = managedAgents.find(item => item.identityName === a[0]); const available = state.openClawAgents === null || Boolean(live); return `<article class="agent-card" data-agent-search="${a[0]} ${a[1]} ${a[4]}"><div class="agent-top"><span class="agent-icon agent-emoji" aria-hidden="true">${emojiIcon(agentEmojis[index], agentFallbackIcons[index] || a[3] || 'bot')}</span><div class="agent-title"><h3>${a[0]}</h3><span>${a[1]}${live?.model ? ` · ${escapeHTML(live.model)}` : ''}</span></div><span class="status-dot ${state.modelConfigured&&available?'online':'warning'}"></span></div><p>${a[4]}</p><div class="prompt-list">${a[5].map(p => `<button data-action="agent-chat" data-agent="${index}" data-prompt="${p}">${p}</button>`).join('')}</div><div class="agent-foot"><span class="secondary-text">OpenClaw · ${state.modelConfigured&&available?'已注册':'待核验'}</span><button class="link-button" data-action="agent-chat" data-agent="${index}">开始对话</button></div></article>`; }).join('')}
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
      <select class="select" id="customerType"><option value="all">全部类型</option>${customerTypeOptions().map(({value,label})=>`<option value="${escapeAttr(value)}" ${state.customerType===value?'selected':''}>${escapeHTML(label)}</option>`).join('')}</select>
      <select class="select" id="customerCountry"><option value="all">全部国家</option>${countryOptions().map(({value,label})=>`<option value="${escapeAttr(value)}" ${state.customerCountry===value?'selected':''}>${escapeHTML(label)}</option>`).join('')}</select>
      <span class="spacer"></span><span class="selection-count">已选 ${state.selectedRows.customers.size} 项</span><button class="button ghost" data-action="column-settings">${icon('columns-3')}列表字段</button>
    </div>
    <div class="data-wrap"><table class="data-table" id="customerTable"><thead><tr><th>${selectAllCheckbox('customers','客户')}</th><th>客户</th><th>类型</th><th>国家</th><th>城市</th><th>联系人</th><th>电话</th><th>邮箱</th><th>网站</th><th>负责人</th><th>来源</th><th>${sortHeader('订单数量','customer','orders',state.customerSort)}</th><th>${sortHeader('累计金额','customer','total',state.customerSort)}</th><th>评级</th><th>${sortHeader('最近更新','customer','updated',state.customerSort)}</th><th>操作</th></tr></thead><tbody>
      ${rows.map(c=>`<tr><td>${rowCheckbox('customers',c.id,c.name)}</td><td><button class="link-button primary-cell" data-action="customer-detail" data-id="${escapeAttr(c.id)}"><span class="avatar">${escapeHTML(c.name.slice(0,2).toUpperCase())}</span><span><strong>${escapeHTML(c.name)}</strong><small>${escapeHTML(c.id)}</small></span></button></td><td>${escapeHTML(localizedCustomerType(c.type))}</td><td>${escapeHTML(localizedCountry(c.country))}</td><td>${escapeHTML(c.city || '未填写')}</td><td><span class="primary-cell"><span><strong>${escapeHTML(c.contact)}</strong><small>${escapeHTML(c.email)}</small></span></span></td><td>${escapeHTML(c.phone || '未填写')}</td><td>${escapeHTML(c.email || '未填写')}</td><td>${escapeHTML(c.website || '未填写')}</td><td>${escapeHTML(c.owner || '未填写')}</td><td>${escapeHTML(c.source || '未填写')}</td><td>${c.orders}</td><td>${escapeHTML(c.total)}</td><td>${badge(c.rating)}</td><td>${escapeHTML(formatLocalizedDateTime(c.updated) || '未记录')}</td><td><span class="table-actions"><button class="table-icon" data-action="customer-detail" data-id="${escapeAttr(c.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-customer" data-id="${escapeAttr(c.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="customer-more" data-id="${escapeAttr(c.id)}" title="更多">${icon('ellipsis')}</button></span></td></tr>`).join('') || `<tr><td colspan="16"><div class="empty-state">${icon('search-x')}<div><h3>未找到客户</h3><p>请调整搜索词或筛选条件。</p></div></div></td></tr>`}
    </tbody></table></div>
    <div class="pagination"><span>共 ${rows.length} 条记录 · 每页 20 条</span><div><button class="button small ghost" disabled>${icon('chevron-left')}</button><button class="button small" data-action="pagination-current">1</button><button class="button small ghost" disabled>${icon('chevron-right')}</button></div></div>
  </div>`;
}

function renderLeads() {
  const query = state.leadSearch.trim().toLowerCase();
  const rows = sortRows(leads.filter(l => (!l.archived) && (!query || [l.id, l.name, l.contact, l.email, l.phone, l.country].join(' ').toLowerCase().includes(query)) && (state.leadType === 'all' || l.type === state.leadType) && (state.leadCountry === 'all' || l.country === state.leadCountry) && (state.leadStatus === 'all' || (state.leadStatus === 'converted' ? l.converted : !l.converted))), state.leadSort, { createdAt: item => timestampValue(item.createdAt), updated: item => timestampValue(item.updated), score: item => Number(item.score || 0) });
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>🎯 线索库</h2><p>从本地客户发现或公开来源整合的客户线索，支持转化为正式客户。</p></div><div class="toolbar"><button class="button primary" data-action="new-lead">${icon('plus')}新建线索</button></div></div>
    <div class="toolbar">
      <label class="field-search">${icon('search')}<input id="leadSearch" value="${escapeAttr(state.leadSearch)}" placeholder="搜索线索、国家、联系人"></label>
      <select class="select" id="leadType"><option value="all">全部类型</option>${customerTypeOptions().map(({value,label})=>`<option value="${escapeAttr(value)}" ${state.leadType===value?'selected':''}>${escapeHTML(label)}</option>`).join('')}</select>
      <select class="select" id="leadCountry"><option value="all">全部国家</option>${countryOptions().map(({value,label})=>`<option value="${escapeAttr(value)}" ${state.leadCountry===value?'selected':''}>${escapeHTML(label)}</option>`).join('')}</select>
      <select class="select" id="leadStatus"><option value="all">全部状态</option><option value="pending" ${state.leadStatus==='pending'?'selected':''}>待跟进</option><option value="converted" ${state.leadStatus==='converted'?'selected':''}>已转换</option></select>
    </div>
    <div class="data-wrap"><table class="data-table"><thead><tr><th>线索</th><th>类型</th><th>国家</th><th>城市</th><th>联系人</th><th>电话</th><th>邮箱</th><th>网站</th><th>来源</th><th>${sortHeader('评分','lead','score',state.leadSort)}</th><th>${sortHeader('创建时间','lead','createdAt',state.leadSort)}</th><th>${sortHeader('更新时间','lead','updated',state.leadSort)}</th><th>状态</th><th>操作</th></tr></thead><tbody>
      ${rows.map(l=>`<tr><td><button class="link-button primary-cell" data-action="lead-detail" data-id="${escapeAttr(l.id)}"><span class="avatar">${escapeHTML(l.name.slice(0,2).toUpperCase())}</span><span><strong>${escapeHTML(l.name)}</strong><small>${escapeHTML(l.id)}</small></span></button></td><td>${escapeHTML(localizedCustomerType(l.type) || '未填写')}</td><td>${escapeHTML(localizedCountry(l.country) || '未填写')}</td><td>${escapeHTML(l.city || '未填写')}</td><td>${escapeHTML(l.contact || '未填写')}</td><td>${escapeHTML(l.phone || '未填写')}</td><td>${escapeHTML(l.email || '未填写')}</td><td>${escapeHTML(l.website || '未填写')}</td><td>${escapeHTML(l.source || '未填写')}</td><td>${l.score ? l.score + ' 分' : '待核实'}</td><td>${escapeHTML(formatLocalizedDateTime(l.createdAt) || '未记录')}</td><td>${escapeHTML(formatLocalizedDateTime(l.updated) || '未记录')}</td><td>${l.converted ? '<span class="badge green">已转换</span>' : '<span class="badge neutral">待跟进</span>'}</td><td><span class="table-actions"><button class="table-icon" data-action="lead-detail" data-id="${escapeAttr(l.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-lead" data-id="${escapeAttr(l.id)}" title="编辑">${icon('pencil')}</button>${!l.converted ? `<button class="table-icon" data-action="convert-lead" data-id="${escapeAttr(l.id)}" title="转为客户">${icon('arrow-right')}</button>` : ''}<button class="table-icon" data-action="delete-lead" data-id="${escapeAttr(l.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="14"><div class="empty-state">${icon('search-x')}<div><h3>未找到线索</h3><p>请调整搜索词或筛选条件。</p></div></div></td></tr>`}
    </tbody></table></div>
  </div>`;
}

function renderQuotes() {
  const query = state.quoteSearch.trim().toLowerCase();
  const visible = sortRows(quotes.filter(q => (!query || [q.id, q.subject, q.customer, q.products].join(' ').toLowerCase().includes(query)) && (state.quoteStatus === 'all' || q.status === state.quoteStatus) && (!state.quoteDateFrom || q.valid >= state.quoteDateFrom) && (!state.quoteDateTo || q.valid <= state.quoteDateTo)), state.quoteSort, { value: quote => moneyNumber(quote.value), updated: quote => quote.updated });
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>📄 报价单管理</h2><p>从客户和产品生成报价，接受后可完整转为订单。</p></div><div class="toolbar"><button class="button" data-action="template-center" data-kind="quote">${icon('layout-template')}模板管理</button><button class="button primary" data-action="new-quote">${icon('plus')}新建报价单</button></div></div>
    <section class="metric-grid" style="grid-template-columns:repeat(4,1fr)">${[['全部报价',42,'files','all'],['草稿',8,'file-pen-line','Draft'],['待客户确认',12,'send','Delivered'],['本月已接受',14,'badge-check','Accepted']].map(([l,v,i,status])=>`<button class="metric-button" data-action="quote-metric-filter" data-status="${status}"><span class="metric-icon">${icon(i)}</span><span><strong class="metric-number">${v}</strong><span class="metric-label">${l}</span></span></button>`).join('')}</section>
    <div class="toolbar"><label class="field-search">${icon('search')}<input id="quoteSearch" value="${escapeAttr(state.quoteSearch)}" placeholder="搜索报价编号、客户或产品"></label><select class="select" id="quoteStatus"><option value="all">全部状态</option>${[['Draft','草稿'],['Delivered','已发送'],['Accepted','已接受'],['Rejected','已拒绝']].map(([v,l])=>`<option value="${v}" ${state.quoteStatus===v?'selected':''}>${l}</option>`).join('')}</select><button class="button ghost" data-action="quote-date-filter">${icon('calendar-days')}有效期</button><span class="result-count">${visible.length} 条 · 已选 ${state.selectedRows.quotes.size} 项</span><span class="spacer"></span><div class="segmented"><button class="${state.quoteView==='table'?'active':''}" data-quote-view="table" title="表格视图">${icon('list')}</button><button class="${state.quoteView==='kanban'?'active':''}" data-quote-view="kanban" title="看板视图">${icon('columns-3')}</button></div></div>
    ${state.quoteView === 'table' ? renderQuoteTable(visible) : renderQuoteKanban(visible)}
  </div>`;
}

function renderQuoteTable(rows=quotes) {
  return `<div class="data-wrap"><table class="data-table"><thead><tr><th>${selectAllCheckbox('quotes','报价单')}</th><th>报价编号</th><th>主题 / 客户</th><th>产品</th><th>${sortHeader('金额','quote','value',state.quoteSort)}</th><th>有效期</th><th>更新时间</th><th>负责人</th><th>状态</th><th>操作</th></tr></thead><tbody>${rows.map(q=>`<tr><td>${rowCheckbox('quotes',q.id,q.id)}</td><td><button class="link-button" data-action="quote-detail" data-id="${escapeAttr(q.id)}">${escapeHTML(q.id)}</button></td><td><span class="primary-cell"><span><strong>${escapeHTML(q.subject)}</strong><small>${escapeHTML(q.customer)}</small></span></span></td><td>${escapeHTML(q.products)}</td><td><strong>${escapeHTML(q.value)}</strong></td><td>${escapeHTML(q.valid)}</td><td>${escapeHTML(q.updated)}</td><td>${escapeHTML(q.owner)}</td><td>${badge(q.status)}</td><td><span class="table-actions"><button class="table-icon" data-action="quote-detail" data-id="${escapeAttr(q.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="download-quote" data-id="${escapeAttr(q.id)}" title="下载">${icon('download')}</button><button class="table-icon" data-action="edit-quote" data-id="${escapeAttr(q.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="delete-quote" data-id="${escapeAttr(q.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="10"><div class="empty-state">${icon('search-x')}<div><h3>未找到报价单</h3><p>请调整搜索或状态筛选。</p></div></div></td></tr>`}</tbody></table></div>`;
}

function renderQuoteKanban(rows=quotes) {
  return `<section class="kanban">${[['Draft','草稿'],['Delivered','已发送'],['Accepted','已接受'],['Rejected','已拒绝']].map(([status,label])=>`<div class="kanban-column"><header class="kanban-head"><strong>${label}</strong><span>${rows.filter(q=>q.status===status).length}</span></header><div class="kanban-body">${rows.filter(q=>q.status===status).map(q=>`<button class="kanban-card" data-action="quote-detail" data-id="${escapeAttr(q.id)}"><h4>${escapeHTML(q.subject)}</h4><p>${escapeHTML(q.customer)} · ${escapeHTML(q.id)}</p><span class="amount">${escapeHTML(q.value)}</span></button>`).join('') || `<div class="empty-state" style="min-height:100px;padding:10px"><p>暂无记录</p></div>`}</div></div>`).join('')}</section>`;
}

function renderOrders() {
  const query = state.orderSearch.trim().toLowerCase();
  const visible = sortRows(orders.filter(o => (!query || [o.id, o.customer, o.quote, o.products].join(' ').toLowerCase().includes(query)) && (state.orderStatus === 'all' || o.status === state.orderStatus) && (!state.orderDateFrom || o.delivery >= state.orderDateFrom) && (!state.orderDateTo || o.delivery <= state.orderDateTo)), state.orderSort, { value: order => moneyNumber(order.value), updated: order => order.updated });
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>📦 订单生命周期</h2><p>订单只承载业务流程与数据快照，不单独维护模板；生成单据请到单据模块关联订单后操作。</p></div><div class="toolbar"><button class="button" data-action="template-center" data-kind="order">${icon('info')}订单说明</button><button class="button primary" data-action="new-order">${icon('plus')}新建订单</button></div></div>
    <div class="toolbar"><label class="field-search">${icon('search')}<input id="orderSearch" value="${escapeAttr(state.orderSearch)}" placeholder="搜索订单、客户或产品"></label><select class="select" id="orderStatus"><option value="all">全部状态</option>${[['Confirmed','已确认'],['Paid','已付款'],['Production','生产中'],['Shipped','已发运'],['Completed','已完成']].map(([v,l])=>`<option value="${v}" ${state.orderStatus===v?'selected':''}>${l}</option>`).join('')}</select><button class="button ghost" data-action="order-date-filter">${icon('calendar-range')}交付日期</button><span class="result-count">${visible.length} 条 · 已选 ${state.selectedRows.orders.size} 项</span><span class="spacer"></span><span class="badge blue">${visible.filter(o=>o.status!=='Completed').length} 个进行中订单</span></div>
    <div class="data-wrap"><table class="data-table"><thead><tr><th>${selectAllCheckbox('orders','订单')}</th><th>订单编号</th><th>客户</th><th>来源报价</th><th>产品</th><th>${sortHeader('金额','order','value',state.orderSort)}</th><th>交付日期</th><th>更新时间</th><th>状态</th><th>操作</th></tr></thead><tbody>${visible.map(o=>`<tr><td>${rowCheckbox('orders',o.id,o.id)}</td><td><button class="link-button" data-action="order-detail" data-id="${escapeAttr(o.id)}">${escapeHTML(o.id)}</button></td><td>${escapeHTML(o.customer)}</td><td>${escapeHTML(o.quote||'手动创建')}</td><td>${escapeHTML(o.products)}</td><td><strong>${escapeHTML(o.value)}</strong></td><td>${escapeHTML(o.delivery)}</td><td>${escapeHTML(o.updated)}</td><td>${badge(o.status)}</td><td><span class="table-actions"><button class="table-icon" data-action="order-detail" data-id="${escapeAttr(o.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-order" data-id="${escapeAttr(o.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="generate-docs" data-id="${escapeAttr(o.id)}" title="生成单据">${icon('files')}</button><button class="table-icon" data-action="delete-order" data-id="${escapeAttr(o.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="10"><div class="empty-state">${icon('search-x')}<div><h3>未找到订单</h3><p>请调整搜索或状态筛选。</p></div></div></td></tr>`}</tbody></table></div>
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
      <div class="filter-row document-type-filter">${[['all','全部类型'],['PI','PI'],['CI','CI'],['PL','PL'],['报关单','报关单'],['合同','合同']].map(([value,label])=>`<button class="filter-chip ${state.documentType===value?'active':''}" data-document-type="${value}">${label}</button>`).join('')}</div>
      <span class="spacer"></span><span class="result-count">${filtering ? `筛选到 ${visible.length} / ${documents.length} 条` : `共 ${documents.length} 条`}</span>
      ${filtering ? `<button class="button ghost small" data-action="clear-document-filters">${icon('rotate-ccw')}重置</button>` : ''}
    </div>
    <div class="data-wrap"><table class="data-table"><thead><tr><th>单据编号</th><th>类型</th><th>客户</th><th>关联订单</th><th>使用模板</th><th>状态</th><th>最近更新</th><th>操作</th></tr></thead><tbody>${visible.map(d=>`<tr><td><button class="link-button" data-action="document-detail" data-id="${escapeAttr(d.id)}">${escapeHTML(d.id)}</button></td><td><span class="badge neutral">${escapeHTML(d.type)}</span></td><td>${escapeHTML(d.customer)}</td><td>${escapeHTML(d.order)}</td><td>${escapeHTML(d.template)}</td><td>${badge(d.status)}</td><td>${escapeHTML(d.updated)}</td><td><span class="table-actions"><button class="table-icon" data-action="document-detail" data-id="${escapeAttr(d.id)}" title="预览">${icon('eye')}</button><button class="table-icon" data-action="download-document" data-id="${escapeAttr(d.id)}" title="下载文件">${icon('file-down')}</button><button class="table-icon" data-action="edit-document" data-id="${escapeAttr(d.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="delete-document" data-id="${escapeAttr(d.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="8"><div class="empty-state">${icon('file-search')}<div><h3>没有符合条件的单据</h3><p>请调整搜索词、单据类型或状态。</p><button class="button small" data-action="clear-document-filters">重置筛选</button></div></div></td></tr>`}</tbody></table></div>
  </div>`;
}

function renderProducts() {
  const visible = productRows();
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>🚲 产品主数据</h2><p>产品被报价、订单和单据引用；停用不会影响已锁定的历史明细。</p></div><div class="toolbar"><button class="button" data-action="import-products">${icon('file-up')}批量导入</button><button class="button primary" data-action="new-product">${icon('plus')}新建产品</button></div></div>
    <div class="toolbar"><label class="field-search">${icon('search')}<input id="productSearch" value="${escapeAttr(state.productSearch)}" placeholder="搜索产品名、编码或 HS CODE"></label><select class="select" id="productCategory"><option value="all">全部类别</option>${['智能设备','智能骑行','整车方案','配件','服务'].map(v=>`<option value="${v}" ${state.productCategory===v?'selected':''}>${v}</option>`).join('')}</select><button class="button ghost" data-action="toggle-product-sort">${icon('arrow-up-down')}${state.productSort==='stockAsc'?'库存升序':'库存降序'}</button><span class="result-count">${visible.length} 条</span><span class="spacer"></span><div class="segmented"><button class="${state.productView==='grid'?'active':''}" data-product-view="grid" title="卡片视图">${icon('grid-2x2')}</button><button class="${state.productView==='table'?'active':''}" data-product-view="table" title="表格视图">${icon('list')}</button></div></div>
    ${state.productView==='grid' ? `<section class="product-grid" id="productGrid">${visible.map((p,i)=>`<article class="product-card" data-product-search="${escapeAttr(`${p.name} ${p.id} ${p.hs}`)}"><button class="product-visual" data-action="product-detail" data-id="${escapeAttr(p.id)}">${icon(i===0?'cpu':i===1?'gauge':'bike')}</button><div class="product-body"><span class="badge ${p.status==='Active'?'green':p.status==='Inactive'?'neutral':'amber'}">${productStatusLabel(p.status)}</span><h3>${escapeHTML(p.name)}</h3><span class="secondary-text">${escapeHTML(p.id)} · HS ${escapeHTML(p.hs)}</span><p>${escapeHTML(p.desc)}</p><div class="product-price"><strong>${escapeHTML(p.price)}</strong><span class="stock">库存 ${p.stock}</span></div></div></article>`).join('') || `<div class="empty-state"><p>未找到产品</p></div>`}</section>` : renderProductTable(visible)}
  </div>`;
}

function productRows() {
  const query = state.productSearch.trim().toLowerCase();
  const rows = products.filter(p => {
    const text = Object.values(p).join(' ').toLowerCase();
    return (!query || text.includes(query)) && (state.productCategory === 'all' || p.category === state.productCategory);
  });
  if (state.productSort === 'stockAsc') rows.sort((a, b) => Number(a.stock) - Number(b.stock));
  else if (state.productSort === 'stockDesc') rows.sort((a, b) => Number(b.stock) - Number(a.stock));
  else rows.sort((a, b) => String(b.updated || '').localeCompare(String(a.updated || '')));
  return rows;
}

function renderProductTable(rows=productRows()) {
  return `<div class="data-wrap"><table class="data-table"><thead><tr><th>产品编码</th><th>产品名称</th><th>类别</th><th>HS CODE</th><th>销售价</th><th>库存</th><th>更新时间</th><th>状态</th><th>操作</th></tr></thead><tbody>${rows.map(p=>`<tr><td>${escapeHTML(p.id)}</td><td><button class="link-button" data-action="product-detail" data-id="${escapeAttr(p.id)}">${escapeHTML(p.name)}</button></td><td>${escapeHTML(p.category)}</td><td>${escapeHTML(p.hs)}</td><td><strong>${escapeHTML(p.price)}</strong></td><td>${p.stock}</td><td>${escapeHTML(p.updated || '未填写')}</td><td>${badge(p.status)}</td><td><span class="table-actions"><button class="table-icon" data-action="product-detail" data-id="${escapeAttr(p.id)}" title="查看">${icon('eye')}</button><button class="table-icon" data-action="edit-product" data-id="${escapeAttr(p.id)}" title="编辑">${icon('pencil')}</button><button class="table-icon" data-action="delete-product" data-id="${escapeAttr(p.id)}" title="删除">${icon('trash-2')}</button></span></td></tr>`).join('') || `<tr><td colspan="9"><div class="empty-state"><p>未找到产品</p></div></td></tr>`}</tbody></table></div>`;
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

async function supplierDetail(id, tab='overview') {
  const s=suppliers.find(item=>item.id===id); if(!s)return;
  const tabs = [['overview','概览'],['activity','沟通记录']];
  const tabBar = `<div class="tabs">${tabs.map(([key,label])=>`<button class="${tab===key?'active':''}" data-action="supplier-tab" data-supplier-tab="${key}" data-supplier-id="${escapeAttr(s.id)}">${label}</button>`).join('')}</div>`;
  if (tab === 'activity' && !Array.isArray(state.supplierCommunications[id])) {
    try {
      const result = await apiFetch(`/api/v1/suppliers/${encodeURIComponent(id)}/communications`);
      state.supplierCommunications[id] = result.items || [];
    } catch (error) {
      toast('沟通记录加载失败',error.message,'warning');
      state.supplierCommunications[id] = [];
    }
  }
  const communications = state.supplierCommunications[id] || [];
  let content = '';
  if (tab === 'overview') {
    content = `<div class="detail-grid" style="margin-top:15px">${[['公司',s.company],['电话',s.phone],['联系人',s.contact],['邮件',s.email],['产品',s.product],['规格',s.specification],['报价',s.quote],['来源',s.source],['备注',s.notes],['更新时间',s.updated]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v||'未填写')}</strong></div>`).join('')}</div>`;
  } else if (tab === 'activity') {
    content = `<div class="spread communication-head"><div><strong>历史沟通记录</strong><p class="secondary-text">记录只能追加，不能修改或删除。</p></div><button class="button primary small" data-action="new-supplier-communication" data-id="${escapeAttr(s.id)}">${icon('message-square-plus')}新增沟通</button></div>${communications.length?`<div class="timeline communication-timeline">${communications.map(item=>`<div class="timeline-item"><div class="spread"><h4>${escapeHTML(item.subject||item.type)}</h4><span class="badge blue">${escapeHTML(item.type)}</span></div><p class="communication-content">${escapeHTML(item.content)}</p><small>${escapeHTML(String(item.occurredAt||'').replace('T',' '))}${item.contact?` · ${escapeHTML(item.contact)}`:''} · 由 ${escapeHTML(item.createdBy)} 记录</small></div>`).join('')}</div>`:`<div class="empty-state panel">${icon('messages-square')}<div><h3>暂无沟通记录</h3><p>新增后将永久保留在本机业务数据库中。</p></div></div>`}`;
  }
  openDrawer({title:s.company,eyebrow:`供应商 / ${s.id}`,body:`${tabBar}<div class="spread" style="margin-bottom:14px"><span class="badge blue">${escapeHTML(s.source)}</span><div class="inline-actions"><button class="button small" data-action="edit-supplier" data-id="${escapeAttr(s.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-supplier" data-id="${escapeAttr(s.id)}">${icon('trash-2')}删除</button></div></div>${content}`});
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
  const availableNews = visibleNewsItems();
  const filteredNews = state.newsCategory === '全部' ? availableNews : availableNews.filter(item => newsItemMatchesFilter(item, state.newsCategory));
  const configuredNews = filteredNews.slice(0, state.newsShowLimit);
  const list = state.newsExpanded ? configuredNews : configuredNews.slice(0,3);
  const latest = filteredNews[0] || availableNews[0] || {};
  return `<div class="page-stack">
    <div class="page-intro"><div><h2>📰 新闻设置</h2><p>由独立推荐智能体按新闻设置采集、去重和排序。</p></div><div class="toolbar"><button class="button" data-action="news-sources">${icon('settings-2')}新闻设置</button><button class="button primary" data-action="refresh-news" ${state.newsRefreshLoading?'disabled':''}>${icon(state.newsRefreshLoading?'loader-circle':'refresh-cw')}${state.newsRefreshLoading?'更新中':'更新新闻'}</button></div></div>
    <div class="filter-row news-category-tabs">${['全部','欧洲市场','法规','智能骑行','渠道','产品'].map(v=>`<button class="filter-chip ${state.newsCategory===v?'active':''}" data-action="news-filter" data-category="${v}">${v}</button>`).join('')}</div>
    <section class="panel"><header class="panel-head"><div><h3>最新资讯</h3><p>本地已缓存，可离线查看已获取内容；单次获取 ${state.newsFetchLimit} 条，展开后按页面上限展示</p></div><span class="meta">更新于 ${escapeHTML(overviewTime(latest.updatedAt || latest.time))} · 获取 ${state.newsFetchLimit} 条 · 展示上限 ${state.newsShowLimit}</span></header><div class="news-list ${state.newsExpanded && state.newsShowLimit > 20 ? 'news-list-scroll' : ''}">${list.map(n=>`<button class="news-list-item" data-action="news-detail" data-title="${escapeAttr(n.title)}"><span class="badge neutral">${escapeHTML(n.category || '行业资讯')}</span><h3>${escapeHTML(cleanVisibleText(n.title))}</h3><p>${escapeHTML(newsPreview(n, 180))}</p><div class="source-line"><span class="mini-source">${escapeHTML(n.source)}</span><span class="mini-source">${escapeHTML(overviewTime(n.updatedAt || n.time))}</span><span class="mini-source">相关度 ${escapeHTML(n.relevance)}</span></div></button>`).join('') || `<div class="empty-state"><p>当前分类暂无资讯</p></div>`}</div><div class="panel-body" style="text-align:center"><button class="button ghost small" data-action="toggle-news">${icon(state.newsExpanded?'chevron-up':'chevron-down')}${state.newsExpanded?'收起':'更多'}</button></div></section>
  </div>`;
}

function openNewsSettings() {
  const frequencies = ['1小时','2小时','3小时','6小时','8小时','12小时','24小时'];
  const categories = newsMediaCategories();
  const filters = state.newsMediaFilterCategories.length ? state.newsMediaFilterCategories : categories.map(item=>item.key);
  const visibleMedia = state.newsMediaCatalog.filter(item=>!filters.length || filters.includes(item.categoryKey || item.category));
  const selectedCount = state.newsMediaIDs.length;
  const catalogBlock = state.newsMediaCatalog.length
    ? `<div class="form-section full"><h3>媒体来源</h3><p>从客户提供的自行车媒体联系表生成，可多选；类型类别会作为媒体描述展示。</p></div><div class="form-field full"><label>类别筛选</label><div class="filter-row media-filter-row">${categories.map(item=>`<label class="filter-chip ${filters.includes(item.key)?'active':''}"><input type="checkbox" data-news-media-category-filter="${escapeAttr(item.key)}" ${filters.includes(item.key)?'checked':''}>${escapeHTML(newsMediaCategoryLabel(item))}</label>`).join('')}</div><small>用于筛选下面的媒体列表，不等同于最终选中。</small></div><div class="form-field full"><div class="spread"><label>媒体名称多选 <span class="required">*</span></label><small>已选 ${selectedCount} 个 · 当前显示 ${visibleMedia.length} 个</small></div><div class="media-source-list">${visibleMedia.map(item=>`<label class="media-source-option ${state.newsMediaIDs.includes(item.id)?'active':''}"><input type="checkbox" data-news-media-id="${escapeAttr(item.id)}" ${state.newsMediaIDs.includes(item.id)?'checked':''}><span><strong>${escapeHTML(newsMediaName(item))}</strong><small>${escapeHTML(newsMediaDescription(item))}</small></span></label>`).join('')}</div></div>`
    : `<div class="form-field full"><div class="display-rule-note">${icon('triangle-alert')}<span><strong>媒体来源表暂未加载</strong><small>仍可通过自定义来源手动填写来源名称或 URL。</small></span></div></div>`;
  openModal({title:'新闻设置',eyebrow:'行业新闻 / 新闻设置',body:`<div class="form-grid"><div class="form-field full"><label>关注国家 <span class="required">*</span></label><input class="input" id="newsCountries" value="${escapeAttr(state.newsCountries)}"><small>多个国家使用顿号或逗号分隔，最多 200 个字符。</small></div><div class="form-field full"><label>关注主题 <span class="required">*</span></label><input class="input" id="newsTopics" value="${escapeAttr(state.newsTopics)}"><small>多个主题使用顿号或逗号分隔，最多 200 个字符。</small></div><div class="form-field"><label>每次获取数量 <span class="required">*</span></label><input class="input" id="newsFetchLimit" type="number" min="1" max="100" step="1" inputmode="numeric" value="${state.newsFetchLimit}"><small>必须是 1-100 的整数。这里控制新闻 Agent 每次拉取多少条。</small></div><div class="form-field"><label>每次展示数量 <span class="required">*</span></label><input class="input" id="newsShowLimit" type="number" min="1" max="100" step="1" inputmode="numeric" value="${state.newsShowLimit}"><small>必须是 1-100 的整数。这里控制页面最多展示多少条。</small></div>${selectField('获取频率',frequencies,false,'newsFrequency',state.newsFrequency)}${catalogBlock}<div class="form-field full"><label>自定义来源</label><textarea class="textarea" id="newsCustomSources" placeholder="每行一个补充来源，可填写媒体名或 URL">${escapeHTML(state.newsCustomSources)}</textarea><small>用于补充媒体表之外的来源；每行最多 200 个字符，总长度最多 1,000 个字符。</small></div><div class="form-field full"><div class="display-rule-note">${icon('layout-list')}<span><strong>展示规则</strong><small>获取数量决定新闻 Agent 每次拉取多少新闻，展示数量决定页面最多显示多少条；超过 20 条时列表内部滚动。</small></span></div></div></div>`,footer:formFooter('保存设置','save-news-settings')});
}

function newsMediaCategories() {
  const map = new Map();
  state.newsMediaCatalog.forEach(item => {
    const key = item.categoryKey || item.category || '其它';
    if (!map.has(key)) map.set(key, { key, category: item.category || key, categoryEn: item.categoryEn || key });
  });
  return [...map.values()].sort((a,b)=>newsMediaCategoryLabel(a).localeCompare(newsMediaCategoryLabel(b),'zh-CN'));
}

function newsMediaCategoryLabel(item) {
  return state.lang === 'en' ? (item.categoryEn || item.category || item.key) : (item.category || item.key);
}

function newsMediaName(item) {
  return state.lang === 'en' ? (item.nameEn || item.nameZh || item.id) : (item.nameZh || item.nameEn || item.id);
}

function newsMediaDescription(item) {
  const category = state.lang === 'en' ? (item.categoryEn || item.category) : (item.category || item.categoryEn);
  const type = state.lang === 'en' ? (item.typeEn || item.type) : (item.type || item.typeEn);
  return [category, type, item.countryRegion, item.traffic].filter(Boolean).join(' · ');
}

function selectedNewsMediaSources() {
  return state.newsMediaCatalog.filter(item=>state.newsMediaIDs.includes(item.id));
}

function buildNewsSourcesText() {
  const mediaLines = selectedNewsMediaSources().map(item => {
    const name = `${newsMediaName(item)}${item.nameEn && item.nameZh && state.lang !== 'en' ? ` / ${item.nameEn}` : ''}`;
    const url = item.newsUrl || item.homeUrl || item.reviewUrl || item.youtube || '';
    const desc = newsMediaDescription(item);
    return [name, desc, url].filter(Boolean).join(' | ');
  });
  const customLines = state.newsCustomSources.split(/\r?\n/).map(line=>line.trim()).filter(Boolean);
  return [...mediaLines, ...customLines].join('\n');
}

function openRecommendationSettings() {
  const frequencies = ['1小时','2小时','3小时','6小时','8小时','12小时','24小时'];
  openModal({
    title: '推荐设置',
    eyebrow: '概览 / 为你推荐',
    body: `<div class="form-grid">
      <div class="form-field full"><label class="toggle-field"><input id="recommendationEnabled" type="checkbox" ${state.subscription ? 'checked' : ''}><span>启用为你推荐自动更新</span></label><small>关闭后不再按频率生成推荐，已缓存的有效推荐仍可查看。</small></div>
      <div class="form-field full"><label>关注国家 <span class="required">*</span></label><input class="input" id="recommendCountries" value="${escapeAttr(state.newsCountries)}"><small>推荐 Agent 会基于这些国家筛选市场、渠道和客户相关信息。</small></div>
      <div class="form-field full"><label>关注主题 <span class="required">*</span></label><input class="input" id="recommendTopics" value="${escapeAttr(state.newsTopics)}"><small>多个主题使用顿号或逗号分隔，最多 200 个字符。</small></div>
      <div class="form-field"><label>每次查询数量 <span class="required">*</span></label><input class="input" id="recommendShowLimit" type="number" min="1" max="20" step="1" inputmode="numeric" value="${state.recommendationShowLimit}"><small>默认 5 条，最多 20 条。</small></div>
      ${selectField('更新频率', frequencies, false, 'recommendFrequency', state.newsFrequency)}
      <div class="form-field full"><div class="display-rule-note">${icon('sparkles')}<span><strong>推荐内容范围</strong><small>推荐由 OpenClaw 推荐 Agent 根据关注国家、关注主题和本机业务数据生成；可展示内容包括市场机会、法规提醒、展会情报、客户线索和产品趋势。没有可靠结果时不会把提示语当成推荐。</small></span></div></div>
    </div>`,
    footer: formFooter('保存推荐设置', 'save-recommendation-settings'),
  });
}

async function saveRecommendationSettings() {
  const saveButton = document.querySelector('[data-action="save-recommendation-settings"]');
  const countries = document.getElementById('recommendCountries');
  const topics = document.getElementById('recommendTopics');
  const limitInput = document.getElementById('recommendShowLimit');
  const frequency = document.getElementById('recommendFrequency');
  [countries, topics, limitInput, frequency].forEach(input => input?.setCustomValidity(''));
  const countriesValue = countries?.value.trim() || '';
  const topicsValue = topics?.value.trim() || '';
  const limitText = limitInput?.value.trim() || '';
  const limit = Number(limitText);
  const allowedFrequencies = ['1小时','2小时','3小时','6小时','8小时','12小时','24小时'];
  if (!countriesValue) return invalidNewsField(countries, '请填写至少一个关注国家。');
  if (countriesValue.length > 200) return invalidNewsField(countries, '关注国家最多 200 个字符。');
  if (!topicsValue) return invalidNewsField(topics, '请填写至少一个关注主题。');
  if (topicsValue.length > 200) return invalidNewsField(topics, '关注主题最多 200 个字符。');
  if (!/^\d+$/.test(limitText) || !Number.isInteger(limit) || limit < 1 || limit > 20) return invalidNewsField(limitInput, '为你推荐每次查询数量必须是 1-20 的整数。');
  if (!allowedFrequencies.includes(frequency?.value || '')) return invalidNewsField(frequency, '请选择允许的更新频率。');
  state.subscription = Boolean(document.getElementById('recommendationEnabled')?.checked);
  state.recommendationShowLimit = limit;
  state.newsCountries = countriesValue;
  state.newsTopics = topicsValue;
  state.newsFrequency = frequency.value;
  if (saveButton) {
    saveButton.disabled = true;
    saveButton.innerHTML = `${icon('loader-circle')}保存中`;
    applyIcons();
  }
  try {
    const response = await savePreferences();
    closeModal();
    renderPage();
    toast('推荐设置已保存', `为你推荐已${state.subscription ? '启用' : '暂停'}，每 ${state.newsFrequency} 更新一次。`);
    if (response?.automationMessage) toast('自动任务同步', response.automationMessage, response.automationSynced === false ? 'warning' : 'success');
  } catch (error) {
    if (saveButton) {
      saveButton.disabled = false;
      saveButton.innerHTML = `${icon('save')}保存推荐设置`;
      applyIcons();
    }
    toast('推荐设置保存失败', error.message, 'warning');
  }
}

function invalidNewsField(input, message) {
  input?.setCustomValidity(message);
  input?.reportValidity();
  input?.focus();
  toast('无法保存新闻设置',message,'warning');
}

async function saveNewsSettings() {
  const saveButton = document.querySelector('[data-action="save-news-settings"]');
  const countries = document.getElementById('newsCountries');
  const topics = document.getElementById('newsTopics');
  const fetchLimitInput = document.getElementById('newsFetchLimit');
  const limitInput = document.getElementById('newsShowLimit');
  const frequency = document.getElementById('newsFrequency');
  const customSourcesInput = document.getElementById('newsCustomSources');
  [countries,topics,fetchLimitInput,limitInput,frequency,customSourcesInput].forEach(input=>input?.setCustomValidity(''));
  const countriesValue = countries?.value.trim() || '';
  const topicsValue = topics?.value.trim() || '';
  const fetchLimitText = fetchLimitInput?.value.trim() || '';
  const fetchLimit = Number(fetchLimitText);
  const customSourcesValue = customSourcesInput?.value.trim() || '';
  const limitText = limitInput?.value.trim() || '';
  const limit = Number(limitText);
  const allowedFrequencies = ['1小时','2小时','3小时','6小时','8小时','12小时','24小时'];
  if (!countriesValue) return invalidNewsField(countries,'请填写至少一个关注国家。');
  if (countriesValue.length > 200) return invalidNewsField(countries,'关注国家最多 200 个字符。');
  if (!topicsValue) return invalidNewsField(topics,'请填写至少一个关注主题。');
  if (topicsValue.length > 200) return invalidNewsField(topics,'关注主题最多 200 个字符。');
  if (!/^\d+$/.test(fetchLimitText) || !Number.isInteger(fetchLimit) || fetchLimit < 1 || fetchLimit > 100) return invalidNewsField(fetchLimitInput,'每次获取数量必须是 1-100 的整数。');
  if (!/^\d+$/.test(limitText) || !Number.isInteger(limit) || limit < 1 || limit > 100) return invalidNewsField(limitInput,'每次展示数量必须是 1-100 的整数。');
  if (!allowedFrequencies.includes(frequency?.value || '')) return invalidNewsField(frequency,'请选择允许的获取频率。');
  if (!state.newsMediaIDs.length && !customSourcesValue) return invalidNewsField(customSourcesInput,'请至少选择一个媒体来源，或填写一个自定义来源。');
  if (customSourcesValue.length > 1000) return invalidNewsField(customSourcesInput,'自定义来源最多 1,000 个字符。');
  if (customSourcesValue.split(/\r?\n/).some(line=>line.length > 200)) return invalidNewsField(customSourcesInput,'自定义来源每行最多 200 个字符。');
  state.newsCountries = countriesValue;
  state.newsTopics = topicsValue;
  state.newsFetchLimit = fetchLimit;
  state.newsShowLimit = limit;
  state.newsFrequency = frequency.value;
  state.newsCustomSources = customSourcesValue;
  state.newsMediaCategories = [...new Set(selectedNewsMediaSources().map(item=>item.categoryKey || item.category).filter(Boolean))];
  state.newsSources = buildNewsSourcesText();
  if (saveButton) {
    saveButton.disabled = true;
    saveButton.innerHTML = `${icon('loader-circle')}保存中`;
    applyIcons();
  }
  try {
    const response = await savePreferences();
    localStorage.setItem('sta100-news-settings', JSON.stringify({countries:state.newsCountries, topics:state.newsTopics, fetchLimit:state.newsFetchLimit, showLimit:state.newsShowLimit, frequency:state.newsFrequency, sources:state.newsSources, mediaIds:state.newsMediaIDs, mediaCategories:state.newsMediaCategories, customSources:state.newsCustomSources}));
    closeModal(); renderPage(); toast('新闻设置已保存',`每 ${state.newsFrequency} 获取，单次最多展示 ${state.newsShowLimit} 条。`);
    if (response?.automationMessage) toast('自动任务同步', response.automationMessage, response.automationSynced === false ? 'warning' : 'success');
  } catch(error) {
    if (saveButton) {
      saveButton.disabled = false;
      saveButton.innerHTML = `${icon('save')}保存设置`;
      applyIcons();
    }
    toast('新闻设置保存失败',error.message,'warning');
  }
}

function renderSettings() {
  const tabs = [['model','模型设置','brain-circuit'],['channels','通道绑定','message-square'],['email','邮件发送','mail'],['scheduler','定时任务','clock-3'],['backup','智能体备份','archive'],['security','数据安全','shield-check'],['system','系统信息','monitor-cog'],['upgrade','版本升级','package-open']];
  return `<div class="settings-layout">
    <nav class="settings-nav panel">${tabs.map(([k,l,i])=>`<button class="${state.settingsTab===k?'active':''}" data-settings-tab="${k}">${icon(i)}${l}</button>`).join('')}</nav>
    <div class="settings-content">${renderSettingsContent()}</div>
  </div>`;
}

function renderSettingsContent() {
  const content = {
    model: renderModelSettings(),
    channels: renderChannelSettings(),
    email: `<section class="panel"><header class="panel-head"><div><h3>邮件发送</h3><p>使用本机配置的 SMTP 发件人发送报价单和订单邮件，支持国内外标准邮箱。</p></div><button class="button primary small" data-action="open-email-settings">${icon('settings')}设置发件人</button></header><div class="setting-row"><span class="setting-icon">${icon('mail')}</span><div class="setting-copy"><strong>发件人和 SMTP</strong><span>配置 SMTP 主机、端口、加密方式和应用专用密码后，报价单和订单详情可直接发送邮件并附加 PDF 单据。</span></div></div><div class="setting-row"><span class="setting-icon">${icon('languages')}</span><div class="setting-copy"><strong>业务邮件模板</strong><span>报价单和订单入口分别提供中文、英文标准模板；主题、正文、收件人、抄送和附件均可在发送前编辑。</span></div></div></section>`,
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
  return `<section class="panel"><header class="panel-head"><div><h3>💬 通道绑定</h3><p>通道用于让飞书、企业微信、微信、Telegram、Slack 等入口真实接入 OpenClaw。客户侧只需要配置模型和通道，系统内部能力由应用自动处理。</p></div><button class="icon-button" data-action="refresh-openclaw-channels" title="刷新通道状态" aria-label="刷新通道状态">${icon('refresh-cw')}</button></header><div class="model-warning"><span>${icon('info')} 页面状态区分通道包、插件启用、账号绑定和运行状态；扫码绑定会直接写入 OpenClaw，成功后自动重启网关并复核。</span></div><div class="settings-filter-row"><label class="field-search">${icon('search')}<input id="channelSearch" value="${escapeAttr(state.channelSearch)}" placeholder="搜索飞书、微信、企业微信、Telegram、Slack 等通道"></label><span class="result-count">已配置 ${configured} 个 · 当前显示 ${visibleChannels.length}/${channels.length}</span></div><div class="settings-subsection">${visibleChannels.map(renderChannelSettingRow).join('')||`<div class="empty-state"><p>未找到匹配通道</p></div>`}</div></section>${renderChannelSkillManagement()}`;
}

function channelLabel(channel) {
  return (state.openClawChannels || []).find(item => item.id === channel)?.name || channel;
}

function channelSkillAgentNames() {
  const seen = new Set();
  return (state.channelSkillDefinitions || []).filter(skill => {
    if (!skill.agentId || seen.has(skill.agentId)) return false;
    seen.add(skill.agentId);
    return true;
  });
}

function routeListText(values) {
  return Array.isArray(values) && values.length ? values.join('、') : '未配置';
}

function renderChannelSkillManagement() {
  const routes = state.channelSkillRoutes || [];
  const sessions = (state.channelSkillSessions || []).slice().sort((a,b) => String(b.updatedAt || '').localeCompare(String(a.updatedAt || ''))).slice(0, 8);
  const routeRows = routes.map(route => `<div class="setting-row channel-skill-route-row"><span class="setting-icon">${icon('workflow')}</span><div class="setting-copy"><strong>${escapeHTML(channelLabel(route.channel))} · ${escapeHTML(route.account)} <small class="secondary-text">${escapeHTML(route.agentId)}</small></strong><span>发送人：${escapeHTML(routeListText(route.allowedSenders))} · 会话：${escapeHTML(routeListText(route.allowedChats))}</span></div><span class="badge ${route.enabled ? 'green' : 'neutral'}">${route.enabled ? '启用' : '停用'}</span><div class="inline-actions compact-actions"><button class="button small" data-action="edit-channel-skill-route" data-id="${escapeAttr(route.id)}">${icon('pencil')}编辑</button><button class="button small" data-action="toggle-channel-skill-route" data-id="${escapeAttr(route.id)}" data-enabled="${route.enabled ? 'true' : 'false'}">${icon(route.enabled ? 'power-off' : 'power')}${route.enabled ? '停用' : '启用'}</button><button class="table-icon" data-action="delete-channel-skill-route" data-id="${escapeAttr(route.id)}" title="删除" aria-label="删除">${icon('trash-2')}</button></div></div>`).join('');
  const sessionRows = sessions.map(session => `<div class="channel-skill-session"><div><strong>${escapeHTML(channelLabel(session.channel))} · ${escapeHTML(session.sender || session.conversation)}</strong><span>${escapeHTML(session.skillId || '尚未选择功能')} · ${escapeHTML(session.state || '未知')}</span></div><small>${escapeHTML(formatLocalizedDateTime(session.updatedAt))}</small></div>`).join('');
  if (state.channelSkillLoading && !state.channelSkillDefinitions.length && !routes.length) return `<section class="panel channel-skill-panel"><div class="empty-state">${icon('loader-circle')}<div><h3>正在读取通道 Skill</h3><p>正在读取功能目录、入站授权和会话状态。</p></div></div></section>`;
  return `<section class="panel channel-skill-panel"><header class="panel-head"><div><h3>通道定制化 Skill</h3><p>配置“通道收到指令 → 选择功能 → 补充信息 → 确认执行”的入站授权规则。</p></div><div class="inline-actions"><button class="button small" data-action="refresh-channel-skill">${icon('refresh-cw')}刷新</button><button class="button primary small" data-action="new-channel-skill-route">${icon('plus')}新增路由</button></div></header>${state.channelSkillError ? `<div class="model-warning error"><span>${icon('triangle-alert')} ${escapeHTML(state.channelSkillError)}</span></div>` : ''}<div class="model-warning"><span>${icon('shield-check')} 必须明确允许的发送人或会话；使用 <code>*</code> 表示允许全部。执行结果只回传原消息会话。</span></div><div class="channel-skill-section"><div class="section-head"><div><h4>入站路由</h4><p>每个通道账号可指定一个默认 Agent，并设置授权范围。</p></div></div>${routeRows || `<div class="empty-state compact-empty"><p>尚未配置入站路由。绑定通道后新增一条路由即可启用定制化 Skill。</p></div>`}</div><div class="channel-skill-section"><div class="section-head"><div><h4>最近 Skill 会话</h4><p>只展示最近 8 条，完整状态保存在本机数据库。</p></div></div>${sessionRows || `<div class="empty-state compact-empty"><p>暂无通道 Skill 会话。</p></div>`}</div></section>`;
}

function channelSkillRouteForm(route = {}) {
  const agents = channelSkillAgentNames();
  const channelsOptions = (state.openClawChannels || []).map(channel => ({value: channel.id, label: `${channel.name}（${channel.id}）`}));
  const agentOptions = agents.map(skill => ({value: skill.agentId, label: `${skill.agentId}（${skill.name}）`}));
  openModal({title: route.id ? '编辑通道 Skill 路由' : '新增通道 Skill 路由', eyebrow: '通道绑定 / 入站授权', body: `<div class="form-grid">${selectField('通道', channelsOptions.length ? channelsOptions : [{value:'feishu',label:'飞书（feishu）'}], false, 'channelSkillChannel', route.channel || 'feishu')}${inputField('账号', route.account || 'default', true, false, 'text', 'channelSkillAccount')}${selectField('默认 Agent', agentOptions, false, 'channelSkillAgent', route.agentId || agents[0]?.agentId || '')}<div class="form-field full"><label for="channelSkillSenders">允许的发送人</label><textarea class="textarea" id="channelSkillSenders" rows="3" placeholder="每行一个 sender ID；全部允许填写 *">${escapeHTML((route.allowedSenders || []).join('\n'))}</textarea></div><div class="form-field full"><label for="channelSkillChats">允许的会话</label><textarea class="textarea" id="channelSkillChats" rows="3" placeholder="每行一个 conversation/chat ID；全部允许填写 *">${escapeHTML((route.allowedChats || []).join('\n'))}</textarea></div><div class="form-field full"><label class="toggle-field"><input id="channelSkillEnabled" type="checkbox" ${route.id ? (route.enabled ? 'checked' : '') : 'checked'}><span>启用该路由</span></label><small>通道账号、发送人或会话必须与入站消息一致，才会进入 Skill 流程。</small></div></div>`, footer: `<button class="button ghost" data-action="close-modal">取消</button><button class="button primary" data-action="save-channel-skill-route" data-id="${escapeAttr(route.id || '')}">${icon('save')}保存路由</button>`});
}

function channelSkillLines(id) {
  return formText(id).split(/\r?\n|[,，]/).map(value => value.trim()).filter(Boolean);
}

async function saveChannelSkillRoute(id = '') {
  const route = {id, channel: formText('channelSkillChannel'), account: formText('channelSkillAccount'), agentId: formText('channelSkillAgent'), allowedSenders: channelSkillLines('channelSkillSenders'), allowedChats: channelSkillLines('channelSkillChats'), enabled: Boolean(document.getElementById('channelSkillEnabled')?.checked)};
  if (!route.account || (!route.allowedSenders.length && !route.allowedChats.length)) { toast('路由信息不完整', '请填写账号，并至少配置一个允许的发送人或会话。', 'warning'); return; }
  try { const data = await apiFetch('/api/v1/channel-skill/routes', {method:'PUT', body:JSON.stringify(route)}); upsertRecord(state.channelSkillRoutes, data.item); closeModal(); renderPage(); toast('Skill 路由已保存', '入站授权规则已写入本机配置。', 'success'); }
  catch (error) { toast('Skill 路由保存失败', error.message, 'warning'); }
}

async function toggleChannelSkillRoute(id, enabled) {
  const route = state.channelSkillRoutes.find(item => item.id === id); if (!route) return;
  try { const data = await apiFetch('/api/v1/channel-skill/routes', {method:'PUT', body:JSON.stringify({...route, enabled: !enabled})}); upsertRecord(state.channelSkillRoutes, data.item); renderPage(); toast(data.item.enabled ? '路由已启用' : '路由已停用', '后续入站消息会按新授权规则处理。', 'success'); }
  catch (error) { toast('路由状态更新失败', error.message, 'warning'); }
}

async function deleteChannelSkillRoute(id) {
  if (!window.confirm('确定删除这条通道 Skill 路由吗？删除后该账号将不再接收定制化 Skill 指令。')) return;
  try { await apiFetch(`/api/v1/channel-skill/routes/${encodeURIComponent(id)}`, {method:'DELETE'}); removeRecord(state.channelSkillRoutes, id); renderPage(); toast('路由已删除', '通道入站授权已撤销。', 'success'); }
  catch (error) { toast('路由删除失败', error.message, 'warning'); }
}

function channelSearchText(channel) {
  return normalizeSearch([channel.id,channel.name,channel.description,channel.origin,channel.status].join(' '));
}

function isChannelNoiseError(value) {
  const text = String(value || '').trim().toLowerCase();
  if (!text) return true;
  return text === 'not configured' || text === 'not configure' || text === 'not connected' || text.includes('not configured') || text.includes('not configure') || text.includes('not connected');
}

function renderChannelSettingRow(channel) {
  const installed = Boolean(channel.installed);
  const accountCount = Number(channel.accountCount||0);
  const installing = isChannelActionLoading(channel.id, 'install');
  const uninstalling = isChannelActionLoading(channel.id, 'uninstall');
  const binding = isChannelActionLoading(channel.id, 'bind');
  const busy = installing || uninstalling || binding;
  const installState = installed ? (channel.enabled ? '已安装 · 已启用' : '已安装 · 未启用') : '未安装';
  const lifecycleButton = installed
    ? `<button class="button danger small" data-action="uninstall-channel" data-channel="${escapeAttr(channel.id)}" ${busy?'disabled':''}>${uninstalling?`${icon('loader-circle')}卸载中`:`${icon('package-minus')}卸载`}</button>`
    : `<button class="button primary small" data-action="install-channel" data-channel="${escapeAttr(channel.id)}" ${busy?'disabled':''}>${installing?`${icon('loader-circle')}安装中`:`${icon('package-plus')}安装`}</button>`;
  const bindingButton = installed
    ? `<button class="button primary small" data-action="open-channel-binding" data-channel="${escapeAttr(channel.id)}" ${busy?'disabled':''}>${icon(binding?'loader-circle':channel.bindingMode==='qr'?'qrcode':channel.bindingMode==='login'?'log-in':'link')}${binding?'绑定中':channel.bindingMode==='qr'?'扫码绑定':channel.bindingMode==='login'?'登录绑定':'配置'}</button>`
    : '';
  const busyAction = installing ? 'install' : uninstalling ? 'uninstall' : 'bind';
  return `<div class="setting-row channel-setting-row ${busy ? 'is-loading' : ''}"><span class="setting-icon">${icon('message-circle')}</span><div class="setting-copy"><strong>${escapeHTML(channel.name)} <small class="secondary-text">${escapeHTML(channel.id)}</small></strong><span>${escapeHTML(channel.description||'OpenClaw 聊天通道')} · ${escapeHTML(installState)} · 已绑定账号 ${accountCount} 个</span>${busy?`<div class="channel-install-progress"><div class="progress"><span></span></div><small>${escapeHTML(channelActionLoadingDetail(channel.id, busyAction) || '正在处理，请稍候。')}</small></div>`:''}</div><div class="inline-actions compact-actions"><button class="button small" data-action="channel-status" data-channel="${escapeAttr(channel.id)}" ${busy?'disabled':''}>${icon('activity')}查看连接状态</button>${lifecycleButton}${bindingButton}</div></div>`;
}

function channelBusinessStatus(channel) {
  if (!channel.installed) return { label: '未安装', className: 'neutral', description: '当前通道包未安装，不能直接绑定。' };
  if (!channel.enabled) return { label: '未启用', className: 'amber', description: '通道包已存在，但 OpenClaw 配置中尚未启用。' };
  if (channel.configured && channel.running) return { label: '已连接', className: 'green', description: '账号已配置，通道正在运行。' };
  if (channel.configured && channel.lastError && !isChannelNoiseError(channel.lastError)) return { label: '绑定异常', className: 'red', description: `账号已配置，但通道运行异常：${channel.lastError}` };
  if (channel.configured) return { label: '已绑定', className: 'green', description: '账号凭据已写入 OpenClaw，当前未运行。' };
  return { label: '可绑定', className: 'blue', description: '通道已启用但尚未绑定账号。' };
}

function humanScheduleValue(value) {
  const text = String(value || '').trim();
  const match = text.match(/^(\d+)([smhd])$/i);
  if (!match) return text;
  const unit = { s: '秒', m: '分钟', h: '小时', d: '天' }[match[2].toLowerCase()];
  return `${match[1]} ${unit}`;
}

function schedulerScheduleText(job) {
  if (job.scheduleKind && job.scheduleValue) {
    if (job.scheduleKind === 'every') return `每 ${humanScheduleValue(job.scheduleValue)}`;
    if (job.scheduleKind === 'at') return `一次性：${job.scheduleValue}`;
    return `Cron：${job.scheduleValue}${job.timezone ? ` · ${job.timezone}` : ''}`;
  }
  return job.schedule || '未设置';
}

function schedulerVisibleMeta(job) {
  return [
    ['任务调度', schedulerScheduleText(job)],
    ['最近执行', humanDateTime(job.lastRun, '尚未执行')],
    ['下次调度', humanDateTime(job.nextRun, job.enabled ? '等待计算' : '已关闭')],
  ];
}

function renderSchedulerSettings() {
  if (state.openClawJobsLoading && !scheduledJobs.length) {
    return `<section class="panel"><div class="empty-state">${icon('loader-circle')}<div><h3>正在读取 OpenClaw 定时任务</h3><p>正在复核调度器、内置任务和运行状态。</p></div></div></section>`;
  }
  const cronEnabled = state.openClawCronStatus?.enabled;
  const headerStatus = cronEnabled === true ? '调度器已启用' : cronEnabled === false ? '调度器已停用' : '调度器状态未知';
  const rows = scheduledJobs.map(job => {
    const finalStatus = schedulerFinalStatus(job);
    const detail = schedulerFinalDetail(job, finalStatus);
    const visibleMeta = schedulerVisibleMeta(job).map(([label, value]) => `<span><strong>${escapeHTML(label)}</strong>${escapeHTML(value)}</span>`).join('');
    const busy = isScheduleActionLoading(job.id);
    const toggling = isScheduleActionLoading(job.id, 'toggle');
    const running = isScheduleActionLoading(job.id, 'run');
    const nextEnabled = !job.enabled;
    return `<div class="setting-row scheduler-row"><span class="setting-icon">${icon(job.builtIn?'timer-reset':'calendar-clock')}</span><div class="setting-copy"><strong class="scheduler-job-name">${escapeHTML(job.name)}${job.builtIn?'<small class="scheduler-job-badge">内置</small>':''}</strong><div class="scheduler-visible-meta">${visibleMeta}</div></div><div class="inline-actions compact-actions scheduler-actions"><span class="button small scheduler-state-button ${finalStatus.className}" title="${escapeAttr(detail)}" aria-label="${escapeAttr(`${finalStatus.label}：${detail}`)}">${escapeHTML(finalStatus.label)}</span><button class="button small" data-action="run-schedule" data-id="${escapeAttr(job.id)}" ${job.enabled&&!busy?'':'disabled'}>${icon(running?'loader-circle':'play')}${running?'执行中':'立即执行'}</button><button class="button small" data-action="schedule-runs" data-id="${escapeAttr(job.id)}" ${busy?'disabled':''}>${icon('history')}运行记录</button><button class="button small" data-action="toggle-schedule" data-id="${escapeAttr(job.id)}" data-enabled="${nextEnabled?'true':'false'}" ${busy?'disabled':''}>${icon(toggling?'loader-circle':nextEnabled?'power':'power-off')}${toggling?'处理中':nextEnabled?'开启':'关闭'}</button><button class="table-icon" data-action="edit-schedule" data-id="${escapeAttr(job.id)}" title="编辑" ${busy?'disabled':''}>${icon('pencil')}</button>${job.builtIn?'':`<button class="table-icon" data-action="delete-schedule" data-id="${escapeAttr(job.id)}" title="删除" ${busy?'disabled':''}>${icon('trash-2')}</button>`}</div></div>`;
  }).join('') || `<div class="empty-state"><p>未读取到定时任务。</p></div>`;
  return `<section class="panel"><header class="panel-head"><div><h3>定时任务</h3><p>任务定义、启停、立即执行和运行记录均由 OpenClaw Cron 实际执行。</p></div><div class="inline-actions"><span class="badge ${cronEnabled === true ? 'green' : cronEnabled === false ? 'red' : 'amber'}">${headerStatus}</span><button class="icon-button" data-action="refresh-openclaw-jobs" title="刷新并复核任务" aria-label="刷新并复核任务">${icon('refresh-cw')}</button><button class="button small" data-action="new-schedule">${icon('plus')}新增任务</button></div></header>${state.openClawJobsError ? `<div class="model-warning error"><span>${icon('triangle-alert')} ${escapeHTML(state.openClawJobsError)}</span><button class="button small" data-action="refresh-openclaw-jobs">重试</button></div>` : ''}<div class="scheduler-summary"><span>${icon('list-checks')} 共 ${scheduledJobs.length} 个任务</span><span>${scheduledJobs.filter(job=>job.builtIn).length} 个内置任务</span><span>${scheduledJobs.filter(job=>job.enabled).length} 个已开启</span></div>${rows}</section>`;
}

function schedulerFinalStatus(job) {
  const status = String(job.status || '').toLowerCase();
  const syncStatus = String(job.syncStatus || '').toLowerCase();
  const businessStatus = String(job.businessStatus || '').toLowerCase();
  if (syncStatus === 'unavailable' || syncStatus === 'error' || syncStatus === 'missing' || status === 'unsynced' || businessStatus === 'failed' || status === 'failed' || status === 'error') {
    return { label: '异常', className: 'red', summary: '悬浮查看异常详情。' };
  }
  if (!job.enabled || status === 'disabled') return { label: '关闭', className: 'neutral', summary: '任务已关闭，不会按计划触发。' };
  if (status === 'running' || businessStatus === 'syncing') return { label: '开启', className: 'blue', summary: '任务已开启，正在执行或同步。' };
  if (businessStatus === 'updated' || status === 'success' || status === 'ok') return { label: '正常', className: 'green', summary: '任务已开启，最近运行正常。' };
  if (businessStatus === 'needs_review') return { label: '正常', className: 'green', summary: '任务已开启，结果可在运行记录复核。' };
  return { label: '开启', className: 'blue', summary: '任务已开启，等待首次执行。' };
}

function schedulerFinalDetail(job, finalStatus) {
  const businessMessage = String(job.businessMessage || '').trim();
  const syncMessage = String(job.syncMessage || '').trim();
  const lastResult = String(job.lastResult || '').trim();
  const errorReason = [job.error, syncMessage, businessMessage, lastResult].map(value => String(value || '').trim()).find(Boolean);
  const statusText = finalStatus.label === '异常'
    ? `异常详情：${errorReason || '任务执行异常，请查看运行记录。'}`
    : finalStatus.summary;
  return [
    statusText,
    `任务调度：${schedulerScheduleText(job)}`,
    `最近执行：${humanDateTime(job.lastRun, '尚未运行')}`,
    `下次调度：${humanDateTime(job.nextRun, '暂未计算')}`,
  ].filter(Boolean).join('\n');
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
  return `<section class="panel model-config-panel"><header class="panel-head"><div><h3>🤖 模型设置</h3><p>用于配置多个可用模型。默认模型必须先真实测试通过，否则部分智能体、概览和定时任务无法正常使用。</p></div><button class="icon-button" data-action="refresh-openclaw-models" title="刷新配置状态" aria-label="刷新配置状态">${icon('refresh-cw')}</button></header>
    ${data?.error ? `<div class="model-warning"><span>${icon('triangle-alert')} ${escapeHTML(data.error)}</span><button class="button small" data-action="refresh-openclaw-models">重试</button></div>` : ''}
    <div class="model-config-body">
      <div class="model-warning model-warning-inline"><span>${icon('info')}</span><span>API Key 写入所选模型凭据；新增只选未配置模型，默认模型请先测试通过。</span></div>
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
function isCurrentDefaultModelUsable() {
  const defaultModel = currentDefaultModelKey();
  if (!defaultModel) return false;
  const entry = configuredModelEntry(defaultModel);
  return Boolean(entry && entry.lastTestStatus === 'passed');
}
function normalizeAgentMessageModel(model='') {
  const selected = String(model || '').trim();
  const defaultModel = currentDefaultModelKey();
  return selected && selected !== defaultModel ? selected : '';
}
function agentSelectedModel(agentID) {
  return normalizeAgentMessageModel(state.agentModelSelections?.[agentID] || '');
}

function chatModelOptions(defaultModel='') {
  const seen = new Set();
  return (state.openClawModels?.configuredModels||[])
    .filter(model=>model.lastTestStatus==='passed' && model.key && model.key !== defaultModel)
    .filter(model=>{
      if (seen.has(model.key)) return false;
      seen.add(model.key);
      return true;
    });
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
  if (state.modelDraftMode !== 'edit') {
    state.modelDraftKey = '';
    return '';
  }
  state.modelDraftKey = '';
  return '';
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
  let date = new Date(value);
  if (Number.isNaN(date.getTime())) date = new Date(String(value).replace(' ', 'T'));
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

function formatDurationMs(value) {
  const ms = Number(value || 0);
  if (!Number.isFinite(ms) || ms <= 0) return '';
  if (ms < 1000) return `${Math.round(ms)} ms`;
  const seconds = ms / 1000;
  if (seconds < 60) return `${seconds.toFixed(seconds >= 10 ? 0 : 1)} 秒`;
  const minutes = Math.floor(seconds / 60);
  const rest = Math.round(seconds % 60);
  return `${minutes} 分 ${rest} 秒`;
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
  const selected = findModelByKey(state.modelDraftKey);
  return selected ? modelSeriesKey(selected) : '';
}

function renderModelFamilySelectOptions(selectedFamily='') {
  const options = modelSeriesGroups(draftModelCatalog()).map(group=>`<option value="${escapeAttr(group.key)}" ${group.key===selectedFamily?'selected':''}>${escapeHTML(group.name)}</option>`).join('');
  return `<option value="">请选择模型系列</option>${options || '<option value="" disabled>暂无可选模型</option>'}`;
}

function renderModelVersionSelectOptions(familyKey='', selected='') {
  if (!familyKey) return '<option value="">请先选择模型系列</option>';
  const groups = modelSeriesGroups(draftModelCatalog());
  const group = groups.find(item=>item.key===familyKey);
  if (!group) return '<option value="">请先选择模型系列</option>';
  const modelKey = selected && group.models.some(model=>model.key===selected) ? selected : '';
  return `<option value="">请选择模型小版本</option>${group.models.map(model=>`<option value="${escapeAttr(model.key)}" ${model.key===modelKey?'selected':''}>${escapeHTML(modelVersionLabel(model))} · ${escapeHTML(model.key)}</option>`).join('')}`;
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
    <div class="form-field"><label for="modelVersionSelect">模型小版本</label><select class="select" id="modelVersionSelect" ${familyKey?'':'disabled'}>${renderModelVersionSelectOptions(familyKey, selectedKey)}</select><small>每次只选一个版本。</small></div>
    <div class="form-field full"><label for="modelEndpointMode">API Key 接入区域</label><select class="select" id="modelEndpointMode"><option value="auto" ${endpointMode==='auto'?'selected':''}>自动识别</option><option value="domestic" ${endpointMode==='domestic'?'selected':''}>国内站 · api.minimaxi.com</option><option value="international" ${endpointMode==='international'?'selected':''}>国际站 · api.minimax.io</option></select><small>MiniMax 的 sk-cp 开头 Key 建议使用国内站；非 MiniMax 模型保持自动即可。保存后会同步到 OpenClaw 当前配置。</small></div>
    <div class="form-field full"><label for="modelAPIKey">API Key</label><div class="credential-field"><input class="input" id="modelAPIKey" type="password" autocomplete="new-password" placeholder="${editing?'留空则使用已保存 API Key 测试或保存':'请输入该模型对应的 API Key'}"><button class="icon-button" type="button" data-action="toggle-api-key-input" title="显示 API Key" aria-label="显示或隐藏 API Key">${icon('eye')}</button></div><small>可临时显示本次输入；已保存密钥由 OpenClaw 存储且不会回传明文。</small></div>
    <div class="model-test-summary full ${test?.ok?'success':test?'warning':''}">${icon(test?.ok?'badge-check':test?'triangle-alert':'info')}${testText}<small>测试连通性会先把 API Key 和接入区域写入 OpenClaw，再向所选模型发起真实调用；只有真实调用成功才会标记为正常。</small></div>
  </div>`;
}

function openModelConfigurationForm(mode='create', model='') {
  state.modelDraftMode = mode === 'edit' ? 'edit' : 'create';
  state.modelDraftOriginalKey = state.modelDraftMode === 'edit' ? model : '';
  state.modelDraftKey = state.modelDraftMode === 'edit' ? model : '';
  const selected = state.modelDraftMode === 'edit' ? (findModelByKey(state.modelDraftKey) || null) : null;
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
  version.disabled = !family;
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
    const entry = configuredModelEntry(model);
    if (!entry || entry.lastTestStatus !== 'passed') {
      toast('默认模型未更新', '默认模型必须先通过真实连通性测试，否则部分智能体和概览功能无法正常使用。', 'warning');
      return;
    }
    if (!window.confirm(`确认将 ${entry.name || model} 设为默认模型？\n默认模型会被未单独选择模型的智能体、概览功能和定时任务使用，请确认它已经真实测试通过。`)) return;
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
  const mode = info?.bindingMode || (channel === 'feishu' ? 'qr' : 'config');
  const isQR = mode === 'qr' && channel === 'feishu';
  const isLogin = mode === 'login';
  const alreadyBound = Number(info?.accountCount || 0) > 0 || Boolean(info?.configured);
  stopChannelQRPoll();
  state.channelQR = null;
  openModal({
    title:`${info?.name||channel} 通道绑定`,
    eyebrow:'OpenClaw / 通道',
    body:`<div class="channel-binding-layout">
      ${!installed?`<div class="model-warning"><span>${icon('triangle-alert')} 当前通道包尚未安装，请先返回列表安装。</span></div>`:''}
      <section class="channel-binding-instructions"><div class="channel-section-title"><strong>绑定方式</strong><span>${isQR?'扫码授权':isLogin?'OpenClaw 登录':'账号配置'}</span></div><p>${isQR?'请使用飞书移动端扫描下方二维码并确认授权。授权完成后，系统会自动写入 OpenClaw、重启网关并复核连接状态。':isLogin?'点击登录后，系统会调用 OpenClaw 的通道登录流程；如果 OpenClaw 返回二维码、网页登录或终端提示，请按返回提示完成。':'填写该通道所需的账号凭据，系统会调用 OpenClaw channels add 写入配置。'}</p>${alreadyBound?`<div class="model-warning"><span>${icon('info')} 当前通道已有绑定账号；重新绑定会按 OpenClaw 通道规则覆盖或新增账号，建议先确认是否需要更换账号。</span></div>`:''}</section>
      <div class="form-grid channel-binding-fields">
        <div class="form-field"><label for="channelAccount">账号</label><input class="input" id="channelAccount" value="default" placeholder="OpenClaw 账号标识"><small>默认账号为 default。</small></div>
        ${isQR?`<div class="form-field"><label for="channelDomain">接入区域</label><select class="select" id="channelDomain"><option value="feishu">飞书（中国大陆）</option><option value="lark">Lark（国际版）</option></select><small>按账号所在区域选择。</small></div>`:''}
        ${!isQR&&!isLogin?`<div class="form-field"><label for="channelToken">Token</label><input class="input" id="channelToken" type="password" autocomplete="off" placeholder="按通道要求填写"><small>仅在当前请求中提交，不在页面回显。</small></div><div class="form-field"><label for="channelSecret">Secret</label><input class="input" id="channelSecret" type="password" autocomplete="off" placeholder="可选"><small>通道没有 Secret 时留空。</small></div>`:''}
      </div>
      ${isQR?`<section class="channel-qr-section"><div class="channel-section-title"><strong>二维码</strong><span>自动生成</span></div><div class="channel-qr-panel" id="channelQRPanel"><div class="channel-qr-empty">${icon('qrcode')}<strong>正在生成二维码...</strong><span>二维码只在本次绑定会话内有效。</span></div></div></section>`:''}
      ${isLogin?`<section class="channel-qr-section"><div class="channel-section-title"><strong>登录绑定</strong><span>${alreadyBound?'可重新绑定':'等待操作'}</span></div><div class="channel-qr-panel" id="channelLoginPanel"><div class="channel-qr-empty">${icon('log-in')}<strong>${alreadyBound?'当前已有绑定账号':'点击下方登录绑定'}</strong><span>${alreadyBound?'如需更换账号，点击重新绑定后按 OpenClaw 返回提示操作。':'未绑定账号时不会先报错，点击登录绑定后再读取 OpenClaw 返回的扫码或网页登录提示。'}</span></div></div></section>`:''}
      <section class="channel-status-section"><div class="channel-section-title"><strong>连接状态</strong><button class="button small" data-action="channel-status" data-channel="${escapeAttr(channel)}">${icon('refresh-cw')}刷新状态</button></div><div class="channel-status-box" id="channelStatusBox"><div class="channel-status-loading">正在读取 OpenClaw 通道状态...</div></div></section>
    </div>`,
    footer:`<button class="button" data-action="close-modal">关闭</button>${isLogin?`<button class="button primary" ${installed?'':'disabled'} data-action="login-channel" data-channel="${escapeAttr(channel)}">${icon(alreadyBound?'refresh-cw':'log-in')}${alreadyBound?'重新绑定':'登录绑定'}</button>`:!isQR?`<button class="button primary" ${installed?'':'disabled'} data-action="save-channel-account" data-channel="${escapeAttr(channel)}">${icon('save')}保存配置</button>`:''}`
  });
  void refreshChannelStatus(channel);
  if (isQR && installed) setTimeout(()=>void saveChannelBinding(channel), 0);
}

async function refreshChannelStatus(channel) {
  const target=document.getElementById('channelStatusBox');
  if(target)target.innerHTML='<div class="channel-status-loading">正在读取 OpenClaw 通道状态...</div>';
  try {
    const status=await apiFetch(`/api/v1/openclaw/channels/${encodeURIComponent(channel)}/status`);
    if(target)target.innerHTML=renderChannelStatusSummary(status);
  } catch(error) {
    if(target)target.innerHTML=`<div class="channel-status-error">${icon('triangle-alert')}<span>状态读取失败：${escapeHTML(error.message)}</span></div>`;
    toast('通道状态读取失败',error.message,'warning');
  }
}

function renderChannelStatusSummary(status={}) {
  const label = String(status.label || '未知');
  const statusClass = label === '已连接' ? 'green'
    : label === '通道异常' || label === '状态查询失败' ? 'red'
    : label === '已配置' || label === '运行中' ? 'blue'
    : label === '未配置' ? 'neutral'
    : 'amber';
  const accounts = Array.isArray(status.accounts) ? status.accounts : [];
  const queryState = String(status.queryState || 'ok');
  const queryHint = queryState === 'timeout'
    ? '状态探测耗时较长，未将查询超时判定为通道异常。'
    : queryState !== 'ok' ? '这是状态查询问题，不代表通道凭据一定失效。' : '';
  const lastError = isChannelNoiseError(status.lastError) ? '' : status.lastError;
  return `<div class="channel-status-summary"><div class="channel-status-main"><span class="badge ${statusClass}">${escapeHTML(label)}</span><strong>${escapeHTML(status.message||'暂无状态说明')}</strong></div><div class="channel-status-facts"><span>账号数：${accounts.length}</span><span>最近检查：${escapeHTML(formatLocalizedDateTime(status.checkedAt))}</span>${queryState!=='ok'?`<span>查询结果：${escapeHTML(queryState==='timeout'?'超时':'失败')}</span>`:''}</div>${queryHint?`<div class="channel-status-note">${icon('info')}<span>${escapeHTML(queryHint)}</span></div>`:''}${lastError?`<div class="channel-status-error">${icon('triangle-alert')}<span>${escapeHTML(lastError)}</span></div>`:''}${accounts.length?`<div class="channel-account-list">${accounts.map(account=>`<div><strong>${escapeHTML(String(account.accountId||'default'))}</strong><span>${escapeHTML(account.label||'未配置')}</span>${account.lastError&&!isChannelNoiseError(account.lastError)?`<small>${escapeHTML(account.lastError)}</small>`:''}</div>`).join('')}</div>`:'<small class="secondary-text">当前没有可展示的账号明细。</small>'}</div>`;
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
  target.innerHTML=`<div class="channel-qr-content">${data.qrDataUrl?`<img class="channel-qr-image" src="${escapeAttr(data.qrDataUrl)}" alt="飞书扫码二维码">`:`<div class="channel-qr-fallback">${icon('link')}<span>当前环境未生成图片二维码，请打开授权地址扫码：</span><a href="${escapeAttr(data.qrUrl||'')}" target="_blank" rel="noopener">${escapeHTML(data.qrUrl||'')}</a></div>`}<div class="channel-qr-meta"><strong>请使用飞书移动端扫描</strong><span>账号：${escapeHTML(data.account||'default')} · 区域：${escapeHTML(data.domain||'feishu')}</span><small>有效期至 ${escapeHTML(formatLocalizedDateTime(data.expiresAt))}</small>${terminal}</div></div>`;
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
    const target=document.getElementById('channelQRPanel');
    if(target)target.innerHTML=`<div class="channel-qr-result error">${icon('triangle-alert')}<span>二维码生成失败：${escapeHTML(error.message)}</span></div>`;
    applyIcons();
  }
}

async function saveChannelAccount(channel) {
  try {
    const data = await apiFetch(`/api/v1/openclaw/channels/${encodeURIComponent(channel)}`, {method:'POST', body:JSON.stringify({
      account:formText('channelAccount')||'default',
      token:formText('channelToken'),
      secret:formText('channelSecret'),
    })});
    toast('通道配置已保存', data.message || 'OpenClaw 通道配置已更新。', 'success');
    await loadOpenClawChannels(true);
    await refreshChannelStatus(channel);
  } catch(error) {
    toast('通道配置失败', error.message, 'warning');
  }
}

async function loginChannel(channel) {
  if (isChannelActionLoading(channel, 'bind')) return;
  setChannelActionLoading(channel, 'bind', '正在调用 OpenClaw 通道登录流程');
  const panel=document.getElementById('channelLoginPanel');
  if(panel){panel.innerHTML=`<div class="channel-qr-empty">${icon('loader-circle')}<strong>正在请求 OpenClaw 登录绑定</strong><span>如果该通道需要扫码或网页确认，成功返回后会在这里展示提示。</span></div>`;applyIcons();}
  try {
    const data = await apiFetch(`/api/v1/openclaw/channels/${encodeURIComponent(channel)}/login`, {method:'POST', body:JSON.stringify({account:formText('channelAccount')||'default'})});
    if(panel)panel.innerHTML=renderChannelLoginResult(data);
    toast('通道登录请求已提交', data.message || '请按 OpenClaw 返回的提示完成登录。', 'success');
    await loadOpenClawChannels(true);
    await refreshChannelStatus(channel);
  } catch(error) {
    if(panel){panel.innerHTML=`<div class="channel-qr-result error">${icon('triangle-alert')}<span>${escapeHTML(error.message)}</span></div>`;applyIcons();}
    toast('通道登录失败', error.message, 'warning');
  } finally {
    clearChannelActionLoading(channel, 'bind');
  }
}

function renderChannelLoginResult(data={}) {
  const output = String(data.output || data.message || 'OpenClaw 已接收登录请求，请按通道要求完成授权。').trim();
  const url = output.match(/https?:\/\/[^\s)）]+/)?.[0] || '';
  return `<div class="channel-qr-content"><div class="channel-qr-meta"><strong>OpenClaw 登录提示</strong><span>账号：${escapeHTML(data.account||'default')}</span>${url?`<a class="button small" href="${escapeAttr(url)}" target="_blank" rel="noopener">${icon('external-link')}打开授权地址</a>`:''}</div><pre class="channel-login-output">${escapeHTML(output)}</pre></div>`;
}

function installChannel(channel) {
  if (isChannelActionLoading(channel, 'install')) return;
  const info = (state.openClawChannels||[]).find(item=>item.id===channel);
  const packageSpec = info?.installSpec || '';
  if (packageSpec) {
    openModal({
      title: `安装 ${info.name||channel} 通道`,
      eyebrow: `OpenClaw / ${channel}`,
      body: `<div class="form-field"><label for="channelPackageSpec">OpenClaw 默认安装包</label><input class="input" id="channelPackageSpec" value="${escapeAttr(packageSpec)}" placeholder="OpenClaw npm spec"><small>该安装包来自当前固定版本 OpenClaw 官方通道目录。安装完成后会自动复核插件目录和启用状态。</small></div>`,
      footer: `<button class="button" data-action="close-modal">取消</button><button class="button primary" data-action="confirm-channel-install" data-channel="${escapeAttr(channel)}">${icon('package-plus')}安装并复核</button>`,
    });
    return;
  }
  openModal({
    title: '安装通道插件',
    eyebrow: `OpenClaw / ${channel}`,
    body: `<div class="form-field"><label for="channelPackageSpec">插件包名或 npm spec</label><input class="input" id="channelPackageSpec" placeholder="例如 @openclaw/example-channel"><small>当前固定版本官方目录没有为该通道声明默认包名，需要填写 OpenClaw 支持的插件包名。</small></div>`,
    footer: `<button class="button" data-action="close-modal">取消</button><button class="button primary" data-action="confirm-channel-install" data-channel="${escapeAttr(channel)}">${icon('package-plus')}安装并复核</button>`,
  });
}

async function executeChannelInstall(channel, packageSpec, triggerButton) {
  if (isChannelActionLoading(channel, 'install')) return;
  setChannelActionLoading(channel, 'install', '正在安装并复核 OpenClaw 通道插件');
  if (triggerButton) {
    triggerButton.disabled = true;
    triggerButton.innerHTML = `${icon('loader-circle')}安装中`;
    applyIcons();
  }
  try {
    const data = await apiFetch(`/api/v1/openclaw/channels/${encodeURIComponent(channel)}/install`, {method:'POST', body:JSON.stringify({packageSpec})});
    toast('通道安装完成', data.message || 'OpenClaw 通道插件已更新。', 'success');
    clearChannelActionLoading(channel, 'install');
    await loadOpenClawChannels(true);
    closeModal();
  } catch(error) {
    toast('通道安装失败', error.message, 'warning');
    if (triggerButton) {
      triggerButton.disabled = false;
      triggerButton.innerHTML = `${icon('package-plus')}安装并复核`;
      applyIcons();
    }
  } finally {
    clearChannelActionLoading(channel, 'install');
  }
}

async function uninstallChannel(channel) {
  const info = (state.openClawChannels||[]).find(item=>item.id===channel);
  if (!window.confirm(`确定卸载 ${info?.name||channel} 通道吗？内置通道会改为停用，外部通道会从 OpenClaw 插件目录卸载。`)) return;
  if (isChannelActionLoading(channel, 'uninstall')) return;
  setChannelActionLoading(channel, 'uninstall', '正在卸载并复核 OpenClaw 通道插件');
  try {
    const data = await apiFetch(`/api/v1/openclaw/channels/${encodeURIComponent(channel)}/uninstall`, {method:'POST', body:'{}'});
    toast('通道处理完成', data.message || 'OpenClaw 通道状态已更新。', 'success');
    clearChannelActionLoading(channel, 'uninstall');
    await loadOpenClawChannels(true);
  } catch(error) {
    toast('通道卸载失败', error.message, 'warning');
  } finally {
    clearChannelActionLoading(channel, 'uninstall');
  }
}

function renderSystemSettings() {
  const status = state.openClawStatus;
  const health = state.systemHealth;
  const managedCount = state.openClawAgents ? businessOpenClawAgents().length : undefined;
  const openClawValue = state.openClawStatusLoading ? '正在读取服务状态' : status?.available ? `${status.version} · ${status.serviceStatus} · RPC ${status.rpcOK?'正常':'异常'}` : status?.error || '暂未读取';
  const agentValue = state.openClawAgentsLoading ? '正在读取 Agent' : managedCount === undefined ? '暂未读取' : `${managedCount} 个 STA-100 业务 Agent 已注册`;
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
  if (sub) sub.addEventListener('change', async e => { const previous=state.subscription; state.subscription=e.target.checked; try { const response = await apiFetch('/api/v1/overview/subscription',{method:'PATCH',body:JSON.stringify({enabled:state.subscription})}); if (response?.automation) state.overviewAutomation = response.automation; toast('订阅设置已更新',response?.automationMessage || (state.subscription?`系统将每 ${state.newsFrequency} 更新一次推荐。`:'自动更新已暂停。'), response?.automationSynced === false ? 'warning' : 'success'); } catch(error) { state.subscription=previous; toast('订阅设置失败',error.message,'warning'); } renderPage(); });
  const agentSearch = document.getElementById('agentSearch');
  if (agentSearch) agentSearch.addEventListener('input', e => filterCards('#agentGrid .agent-card', 'agentSearch', e.target.value));
  const customerSearch = document.getElementById('customerSearch');
  if (customerSearch) customerSearch.addEventListener('input', e => { state.customerSearch = e.target.value; renderPage(); requestAnimationFrame(()=>{ const i=document.getElementById('customerSearch'); if(i){i.focus();i.setSelectionRange(i.value.length,i.value.length);} }); });
  const customerType = document.getElementById('customerType');
  if (customerType) customerType.addEventListener('change', e => { state.customerType=e.target.value; renderPage(); });
  const customerCountry = document.getElementById('customerCountry');
  if (customerCountry) customerCountry.addEventListener('change', e => { state.customerCountry=e.target.value; renderPage(); });
  const leadSearch = document.getElementById('leadSearch');
  if (leadSearch) leadSearch.addEventListener('input', e => { state.leadSearch = e.target.value; renderPage(); });
  const leadType = document.getElementById('leadType');
  if (leadType) leadType.addEventListener('change', e => { state.leadType=e.target.value; renderPage(); });
  const leadCountry = document.getElementById('leadCountry');
  if (leadCountry) leadCountry.addEventListener('change', e => { state.leadCountry=e.target.value; renderPage(); });
  const leadStatus = document.getElementById('leadStatus');
  if (leadStatus) leadStatus.addEventListener('change', e => { state.leadStatus=e.target.value; renderPage(); });
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
  const oemQuery = document.getElementById('oemQuery');
  if (oemQuery) oemQuery.addEventListener('input', e => { state.oemQuery = e.target.value; });
  const discoveryCountry = document.getElementById('discoveryCountry');
  if (discoveryCountry) discoveryCountry.addEventListener('change', e => {
    state.discoveryCountry = e.target.value;
    state.discoveryCities = [(discoveryCities[state.discoveryCountry] || [])[0] || ''].filter(Boolean);
    renderPage();
  });
  document.getElementById('discoveryCitiesRow')?.addEventListener('click', e => {
    const btn = e.target.closest('[data-discovery-city]');
    if (!btn) return;
    const city = btn.getAttribute('data-discovery-city');
    state.discoveryCities = state.discoveryCities || [];
    if (state.discoveryCities.includes(city)) {
      state.discoveryCities = state.discoveryCities.filter(c => c !== city);
    } else {
      state.discoveryCities = [...state.discoveryCities, city];
    }
    renderPage();
  });
  document.getElementById('discoveryTypesRow')?.addEventListener('click', e => {
    const btn = e.target.closest('[data-discovery-type]');
    if (!btn) return;
    const type = btn.getAttribute('data-discovery-type');
    state.discoveryTypes = state.discoveryTypes || [];
    if (state.discoveryTypes.includes(type)) {
      state.discoveryTypes = state.discoveryTypes.filter(t => t !== type);
    } else {
      state.discoveryTypes = [...state.discoveryTypes, type];
    }
    renderPage();
  });
  const discoveryShowLimit = document.getElementById('discoveryShowLimit');
  if (discoveryShowLimit) discoveryShowLimit.addEventListener('change', e => {
    const value = Number(e.target.value);
    if (Number.isInteger(value) && value >= 1 && value <= 100) {
      e.target.setCustomValidity('');
      state.discoveryShowLimit = value;
    }
  });
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
  ['customer','客户'],['type','客户类型'],['country','国家'],['city','城市'],['contact','联系人'],['phone','电话'],['email','邮箱'],['website','网站'],['owner','负责人'],['source','来源'],['orders','订单数'],['total','累计金额'],['rating','评级'],['updated','更新时间'],
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
  const indexes={customer:2,type:3,country:4,city:5,contact:6,phone:7,email:8,website:9,owner:10,source:11,orders:12,total:13,rating:14,updated:15};
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

// confirmChoice 弹出"是/否"选择对话框，返回 Promise<boolean>。
// 是 -> resolve(true)（保留原记录）；否 -> resolve(false)（删除原记录）。
// 取消（关闭弹窗）-> resolve(null)。
function confirmChoice({ title, question, yesLabel='是', noLabel='否', description='' }) {
  return new Promise(resolve => {
    let resolved = false;
    let cleanup = () => {};
    const finish = (value) => {
      if (resolved) return;
      resolved = true;
      cleanup();
      closeModal();
      resolve(value);
    };
    const body = `<div class="confirm-dialog"><p class="confirm-question">${escapeHTML(question)}</p>${description ? `<p class="secondary-text" style="margin-top:8px">${escapeHTML(description)}</p>` : ''}</div>`;
    const footer = `<button class="button" data-action="confirm-choice-cancel">取消</button><button class="button danger" data-action="confirm-choice-no">${escapeHTML(noLabel)}</button><button class="button primary" data-action="confirm-choice-yes">${escapeHTML(yesLabel)}</button>`;
    if (!document.getElementById('drawerBackdrop').hidden) closeDrawer();
    document.getElementById('modalTitle').textContent = title;
    document.getElementById('modalEyebrow').textContent = '确认操作';
    document.getElementById('modalBody').innerHTML = body;
    document.getElementById('modalFooter').innerHTML = footer;
    document.getElementById('modal').classList.remove('wide');
    document.getElementById('modalBackdrop').hidden = false;
    syncOverlayScroll();
    applyIcons();
    const onChoice = (event) => {
      const target = event.target.closest('[data-action]');
      if (!target) return;
      const action = target.dataset.action;
      if (action === 'confirm-choice-yes') finish(true);
      else if (action === 'confirm-choice-no') finish(false);
      else if (action === 'confirm-choice-cancel' || action === 'close-modal') finish(null);
    };
    const footerEl = document.getElementById('modalFooter');
    footerEl.addEventListener('click', onChoice);
    // 同步拦截 backdrop 上的取消
    const backdrop = document.getElementById('modalBackdrop');
    const backdropHandler = (event) => {
      if (event.target === backdrop) finish(null);
    };
    backdrop.addEventListener('click', backdropHandler);
    cleanup = () => {
      footerEl.removeEventListener('click', onChoice);
      backdrop.removeEventListener('click', backdropHandler);
    };
  });
}
function closeModal() { stopChannelQRPoll(); state.activeAgentChatID=''; document.getElementById('modalBackdrop').hidden = true; syncOverlayScroll(); }
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
function selectField(label, options, full=false, id='', selected='') {
  return `<div class="form-field ${full?'full':''}"><label>${label}</label><select class="select" ${id?`id="${id}"`:''}>${options.map(option => {
    const value = typeof option === 'object' ? option.value : option;
    const text = typeof option === 'object' ? option.label : option;
    return `<option value="${escapeAttr(value)}" ${String(value)===String(selected)?'selected':''}>${escapeHTML(text)}</option>`;
  }).join('')}</select></div>`;
}
function relationField(label, target, items, selected='') {
  const values = items.map(item => typeof item === 'string' ? item : item.label);
  return `<div class="form-field relation-field"><label>${label}</label><div class="relation-picker"><input class="input relation-input" id="${target}" data-relation-input="${target}" value="${escapeAttr(selected)}" autocomplete="off" placeholder="输入名称进行模糊匹配"><div class="relation-options" id="${target}Options" hidden></div></div></div>`;
}
function relationOptions(target, values, query='') {
  const needle = String(query || '').trim().toLowerCase();
  if (!needle || needle === '__all__') return values.slice(0, 20).map(value => `<button type="button" class="relation-option" data-action="relation-select" data-target="${target}" data-value="${escapeAttr(value)}">${escapeHTML(value)}</button>`).join('') || `<span class="relation-empty">未找到匹配项</span>`;
  const matches = values.filter(value => value.toLowerCase().includes(needle)).slice(0, 8);
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
function timestampValue(value) {
  if (!value) return 0;
  const text = String(value).trim();
  const parsed = new Date(text).getTime();
  if (Number.isFinite(parsed)) return parsed;
  const normalized = new Date(text.replace(' ', 'T')).getTime();
  return Number.isFinite(normalized) ? normalized : 0;
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
function quoteDraftSubtotal() { return state.quoteDraftLines.reduce((total, line) => total + lineSubtotal(line, true), 0); }
function quoteDraftTotal() { return quoteDraftSubtotal() + formNumber('quoteFreight'); }
function refreshQuoteDraftTotals(currency=formText('quoteCurrency')||'EUR') {
  const subtotal = document.getElementById('quoteDraftSubtotal'); if (subtotal) subtotal.textContent = formatMoney(quoteDraftSubtotal(), currency);
  const freight = document.getElementById('quoteDraftFreight'); if (freight) freight.textContent = formatMoney(formNumber('quoteFreight'), currency);
  const tax = document.getElementById('quoteDraftTax'); if (tax) tax.textContent = formatMoney(formNumber('quoteTax'), currency);
  const total = document.getElementById('quoteDraftTotal'); if (total) total.textContent = formatMoney(quoteDraftTotal(), currency);
}
function renderQuoteDraftLines() {
  const body = document.getElementById('quoteLinesBody');
  if (!body) return;
  body.innerHTML = state.quoteDraftLines.map((line, index) => {
    const product = lineProduct(line);
    return `<tr><td><select class="select line-product" data-quote-line-field="productId" data-index="${index}">${productOptions(line.productId)}</select><small class="line-stock">可用库存 ${product?.stock ?? 0}</small></td><td><input class="input line-number" type="number" min="1" max="${product?.stock ?? 1}" value="${line.quantity}" data-quote-line-field="quantity" data-index="${index}"></td><td><input class="input line-price" type="number" min="0" step="0.01" value="${Number(line.unitPrice).toFixed(2)}" data-quote-line-field="unitPrice" data-index="${index}"></td><td><input class="input line-number" type="number" min="0" max="100" step="0.1" value="${Number(line.discount||0)}" data-quote-line-field="discount" data-index="${index}"></td><td><strong data-quote-line-subtotal="${index}">${formatMoney(lineSubtotal(line,true))}</strong></td><td><button type="button" class="table-icon" data-action="remove-quote-line" data-index="${index}" title="删除明细">${icon('trash-2')}</button></td></tr>`;
  }).join('');
  refreshQuoteDraftTotals();
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
  const items={tasks:tasks.slice(0,8).map(item=>`${item.title} · ${item.status || 'open'} · ${formatLocalizedDateTime(item.updatedAt || item.createdAt) || '未记录'}`),meetings:['暂无会议明细'],documents:['暂无文档处理明细'],orders:['暂无订单明细'],chats:['暂无智能体会话明细'],news:['暂无行业资讯明细']}[key] || ['暂无明细'];
  openModal({title:m.label,eyebrow:'今日业务摘要',body:`<div class="panel" style="margin-bottom:14px"><div class="panel-body"><span class="metric-label">当前数量</span><strong class="metric-number" style="display:block;margin-top:6px">${m.value}</strong><p class="secondary-text">${m.detail}</p></div></div><div class="timeline">${items.map(v=>`<div class="timeline-item"><h4>${escapeHTML(v)}</h4><p>来自 /api/v1/tasks 的实时待办列表</p></div>`).join('')}</div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});
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
    'local-retrieval': 'Agent 知识库',
    'agent-knowledge-context': 'Agent 知识库',
    'query-fields': '查询条件',
    attachments: '附件处理',
    queue: '排队等待',
    'target-agent': '当前智能体',
    'openclaw-agent': 'OpenClaw 调用',
    'knowledge-agent': 'Agent 知识库',
    'domain-agents': '业务 Agent',
    'coordinator-agent': '统一汇总',
  }[key] || key || '处理阶段';
  const status = String(stage.status || '').toLowerCase();
  return { key, label, status, detail: stage.detail || '', durationMs: Number(stage.durationMs || 0), reason: stage.reason || '', data: stage.data || '' };
}

function renderMessagePipeline(message={}) {
  const pipeline = Array.isArray(message.pipeline) ? message.pipeline.map(normalizePipelineStage) : [];
  if (!pipeline.length) return '';
  return `<div class="message-pipeline">${pipeline.map(stage=>{
    const tip = pipelineStageTooltip(stage);
    const stageIcon = stage.status === 'failed' ? 'circle-x' : stage.status === 'cancelled' ? 'square' : stage.status === 'partial' ? 'triangle-alert' : 'check-circle-2';
    return `<span class="pipeline-chip ${escapeAttr(stage.status || 'pending')}" title="${escapeAttr(tip)}">${icon(stageIcon)}${escapeHTML(stage.label)}${stage.durationMs?`<small>${escapeHTML(formatDurationMs(stage.durationMs))}</small>`:''}</span>`;
  }).join('')}</div>`;
}

function renderMessageTiming(message={}) {
  const total = Number(message.totalDurationMs || 0);
  if (!total) return '';
  return `<div class="message-timing">${icon('clock-3')}<span>总耗时 ${escapeHTML(formatDurationMs(total))}</span></div>`;
}

function pipelineStageTooltip(stage={}) {
  return stage.data || stage.detail || '该阶段未返回可展示的数据。';
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
  const activeStep = steps[Math.min(activeIndex, Math.max(0, steps.length - 1))] || {};
  const progressDetail = progress.detail || (done ? 'OpenClaw 已返回结果' : `${activeStep.label || 'OpenClaw'}处理中...`);
  return `<div class="chat-progress ${failed?'failed':done?'done':''}" title="${escapeAttr(progressDetail)}"><div class="progress-line">${steps.map((step,index)=>{
    const stepStatus = String(step.status || '').toLowerCase();
    const stateClass = stepStatus === 'failed' || stepStatus === 'error'
      ? 'failed'
      : stepStatus === 'partial'
        ? 'partial'
        : failed && index===activeIndex ? 'failed' : index < activeIndex || done ? 'done' : index === activeIndex ? 'active' : 'pending';
    const tip = progressStepTooltip(step, stateClass, progress);
    return `<span class="${stateClass}" title="${escapeAttr(tip)}"><i></i>${escapeHTML(step.label)}${step.durationMs?`<small>${escapeHTML(formatDurationMs(step.durationMs))}</small>`:''}</span>`;
  }).join('')}</div><div class="progress-time"><span>${escapeHTML(done?'已完成':failed?'生成失败':'生成中')}</span><div class="progress-time-track"><i class="${timeState}"></i></div><strong>${elapsedLabel}</strong></div></div>`;
}

function progressStepTooltip(step={}, stateClass='', progress={}) {
  return step.data || step.detail || (stateClass === 'active' ? progress.detail : '') || '该阶段暂未返回可展示的数据。';
}

function renderAgentMessage(message, index, messageIndex) {
  const modelMeta = message.model ? ` · ${message.provider ? `${escapeHTML(message.provider)}/` : ''}${escapeHTML(message.model)}` : '';
  const caseAction = message.role !== 'user' && !message.error && !message.caseConfirmed ? `<button class="link-button" data-action="confirm-agent-answer" data-agent="${index}" data-message-index="${messageIndex}">确认并存入案例</button>` : '';
  return `<div class="message-row ${message.role === 'user' ? 'user' : ''}"><div class="message ${message.role === 'user' ? 'user' : message.error ? 'error' : ''}">${escapeHTML(message.text).replace(/\n/g,'<br>')}${renderMessageAttachments(message)}${renderMessagePipeline(message)}${renderMessageTiming(message)}${message.error && message.retry ? `<button class="link-button message-retry" data-action="retry-chat" data-agent="${index}" data-message="${escapeAttr(message.retry || '')}">重新发送</button>` : ''}${caseAction}<time>${escapeHTML(message.time || '')}${modelMeta}</time></div></div>`;
}

function renderAgentChatBody(index) {
  const a = agents[index];
  const agentID = agentIDs[index];
  const messages = state.agentChats[agentID] || [];
  const historyLoading = Boolean(state.agentChatHistoryLoading[agentID]);
  const selectedModel = agentSelectedModel(agentID);
  const defaultModel = currentDefaultModelKey();
  const defaultUsable = isCurrentDefaultModelUsable();
  const testedModels = chatModelOptions(defaultModel);
  const modelOptions = testedModels.map(model=>`<option value="${escapeAttr(model.key)}" ${selectedModel===model.key?'selected':''}>${escapeHTML(model.name)} · ${escapeHTML(model.key)}</option>`).join('');
  const agentAttachments = state.chatAttachmentsByAgent?.[agentID] || [];
  const attachments = agentAttachments.map((attachment,attachmentIndex)=>renderComposerAttachment(attachment,index,attachmentIndex)).join('');
  const modelHint = state.openClawModelsLoading
    ? '<div class="chat-model-loading"><span class="secondary-text">正在读取 OpenClaw 已配置模型...</span></div>'
    : !defaultUsable && !testedModels.length
      ? '<div class="chat-model-loading"><span class="secondary-text">暂无可用默认模型，请先到设置中配置并测试模型。</span></div>'
      : defaultModel && !defaultUsable
        ? '<div class="chat-model-loading"><span class="secondary-text">当前默认模型尚未测试通过，请选择已测试模型或到设置中重新测试。</span></div>'
        : '';
  const busy = Boolean(state.agentChatPending[agentID]);
  const quickPrompts = a[5].slice(0,3).map(v=>`<button ${busy?'disabled':''} data-action="chat-quick-prompt" data-agent="${index}" data-prompt="${escapeAttr(v)}">${escapeHTML(v)}</button>`).join('');
  return `<div class="chat-layout"><aside class="chat-side"><div class="chat-side-head"><h3>统一智能处理</h3><span>系统编排</span></div><div class="chat-source policy-fixed"><span class="agent-icon">${icon('workflow')}</span><span><strong>本地证据 → 协调器 → 领域 Agent</strong><small>页面不区分本地或联网结果，系统按规则自动整合。</small></span></div><button class="button ghost small source-settings-button" data-action="agent-allowlist" data-agent="${index}">${icon('shield-check')}联网白名单</button><p class="chat-source-note">白名单属于后台安全配置，不是本次聊天的来源选择。冲突数据会全部保留并标记。</p></aside><section class="chat-main"><header class="chat-head"><div><strong>${agentEmojis[index]} ${escapeHTML(a[0])}</strong><span>${escapeHTML(agentID)}</span></div><span>${badge(state.modelConfigured?'Active':'Review')}</span></header><div class="chat-messages" id="chatMessages"><div class="message-row"><div class="message"><strong>${escapeHTML(a[0])}</strong><br>已连接 STA-100 统一任务协调器，消息将按页面规则分发给专业智能体。<time>当前会话</time></div></div>${historyLoading ? `<div class="chat-history-loading">${icon('loader-circle')}正在从 OpenClaw 读取历史消息...</div>` : ''}${messages.map((m,messageIndex)=>renderAgentMessage(m,index,messageIndex)).join('')}</div><div class="chat-status" id="chatStatus" data-agent="${escapeAttr(agentID)}" aria-live="polite">${renderAgentChatProgress(agentID)}</div><div class="chat-composer"><div id="chatAttachmentList" class="chat-attachments composer-attachments">${attachments}</div><div class="chat-compose-tools"><div class="chat-attachment-tools"><button class="icon-button" data-action="choose-chat-image" title="上传图片" aria-label="上传图片">${icon('image-up')}</button><button class="icon-button" data-action="choose-chat-file" title="上传文件" aria-label="上传文件">${icon('file-up')}</button></div><span class="chat-compose-spacer"></span><label class="chat-model-picker"><span>模型</span><select class="select" id="chatModelSelect"><option value="">使用默认模型${defaultModel?`（${escapeHTML(defaultModel)}）`:''}</option>${modelOptions}</select></label><input id="chatImageInput" type="file" accept="image/*" hidden><input id="chatFileInput" type="file" hidden></div>${modelHint}<div class="chat-input-row"><textarea class="textarea" id="chatInput" maxlength="32768" rows="2" placeholder="${busy?'当前消息处理中，可点击右侧停止':'向 '+escapeAttr(a[0])+' 发送消息，Enter 发送，Shift+Enter 换行'}"></textarea><button class="icon-button chat-send ${busy?'chat-stop':''}" data-action="${busy?'stop-chat':'send-chat'}" data-agent="${index}" title="${busy?'停止当前 Agent 调用':'发送消息'}" aria-label="${busy?'停止当前 Agent 调用':'发送消息'}">${icon(busy?'square':'send')}</button></div><div class="chat-quick-prompts inline">${quickPrompts}</div></div></section></div>`;
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

async function assistantQueueNotice(index) {
  try {
    const status = await apiFetch('/api/v1/assistant/queue');
    const running = Array.isArray(status.running) ? status.running : [];
    const limit = Number(status.limit || 2);
    if (limit > 0 && running.length >= limit) {
      const names = running.map(item=>item.name||item.agentId||'Agent 任务').filter(Boolean).slice(0, limit);
      const agentName = agents[index]?.[0] || agentIDs[index] || '当前 Agent';
      return {
        names,
        message: `前面已有 ${names.join('、')} 正在运行，${agentName} 会进入排队队列，可能稍慢。`,
      };
    }
  } catch {
    return null;
  }
  return null;
}

function markAgentProgressQueued(agentID, index, notice) {
  if (!notice?.message) return;
  const progress = state.agentChatProgress[agentID];
  if (!progress) return;
  progress.detail = notice.message;
  if (!progress.steps?.some(step=>step.key === 'queue')) {
    progress.steps = [{key:'queue',label:'排队等待',status:'active',detail:notice.message,data:notice.names?.join('\n')||''}, ...(progress.steps||[])];
  }
  progress.index = 0;
  refreshAgentChat(index);
}

async function uploadChatAttachment(file,index) {
  if(!file)return;
  const agentID = agentIDs[index];
  if(file.size>25*1024*1024){toast('附件过大','单个附件不能超过 25 MB。','warning');return;}
  const form=new FormData();
  form.append('file',file);
  try {
    const data=await apiFetch('/api/v1/assistant/attachments',{method:'POST',body:form});
    const attachment = data.attachment || {};
    attachment.localId = `att-${Date.now()}-${Math.random().toString(16).slice(2)}`;
    attachment.previewUrl = file.type.startsWith('image/') ? URL.createObjectURL(file) : '';
    attachment.status = '已上传到 STA-100，待随消息发送';
    if (!state.chatAttachmentsByAgent[agentID]) state.chatAttachmentsByAgent[agentID] = [];
    state.chatAttachmentsByAgent[agentID].push(attachment);
    refreshAgentChat(index);
  } catch(error) { toast('附件上传失败',error.message,'warning'); }
}

async function showAgentChat(index,prompt='') {
  const a=agents[index];
  const agentID = agentIDs[index];
  state.activeAgentChatID = agentID;
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
  if (state.activeAgentChatID && state.activeAgentChatID !== agentID) return;
  if (!document.getElementById('modalBody')) return;
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

function isAbortError(error) {
  return error?.name === 'AbortError' || /abort|aborted|cancel/i.test(String(error?.message || ''));
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
    retry: message.retry || '',
    pipeline: Array.isArray(message.pipeline) ? message.pipeline : [],
    totalDurationMs: Number(message.totalDurationMs || 0),
    sources: Array.isArray(message.sources) ? message.sources : [],
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
  const agentID = agentIDs[index];
  const attachments = state.chatAttachmentsByAgent[agentID] || [];
  const [attachment] = attachments.splice(Number(attachmentIndex), 1);
  if (attachment?.previewUrl) URL.revokeObjectURL(attachment.previewUrl);
  refreshAgentChat(Number(index || 0));
}

function startAgentProgress(agentID,index,hasAttachments=false,startedAt=Date.now()) {
  stopAgentProgress(agentID);
  const steps = chatProgressSteps
    .filter(([key])=>hasAttachments || key !== 'attachments')
    .map(([key,label])=>({key,label}));
  state.agentChatProgress[agentID] = {index:0,status:'running',steps,startedAt,detail:'正在读取当前 Agent 已同步知识库上下文'};
  refreshAgentChat(index);
  chatProgressTimers[agentID] = window.setInterval(()=>{
    const progress = state.agentChatProgress[agentID];
    if (!progress || progress.status !== 'running') return;
    const maxIndex = Math.max(0, progress.steps.length - 2);
    progress.index = Math.min(maxIndex, Number(progress.index || 0) + 1);
    progress.detail = progress.steps[progress.index]?.label ? `${progress.steps[progress.index].label}处理中...` : 'OpenClaw 正在处理';
    const status = document.getElementById('chatStatus');
    if (status && state.activeAgentChatID === agentID && status.dataset.agent === agentID) {
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

function finishAgentProgress(agentID,index,result,errorMessage='',startedAt=0) {
  stopAgentProgress(agentID);
  const failed = Boolean(errorMessage);
  const previous = state.agentChatProgress[agentID] || {};
  const startMs = Number(startedAt || previous.startedAt || 0);
  const elapsedSeconds = startMs ? Math.max(1, Math.ceil((Date.now() - startMs) / 1000)) : 0;
  const pipeline = Array.isArray(result?.pipeline) && result.pipeline.length ? result.pipeline.map(normalizePipelineStage) : null;
  const steps = pipeline ? pipeline.map(stage=>({key:stage.key,label:stage.label,detail:stage.detail,reason:stage.reason,durationMs:stage.durationMs,status:stage.status,data:stage.data})) : (previous.steps || chatProgressSteps.map(([key,label])=>({key,label})));
  const failedIndex = pipeline ? Math.max(0,pipeline.findIndex(stage=>stage.status==='failed')) : -1;
  state.agentChatProgress[agentID] = {
    steps,
    index: failed ? (failedIndex >= 0 ? failedIndex : Math.max(0,steps.length-1)) : Math.max(0,steps.length-1),
    status: failed ? 'failed' : 'done',
    elapsedSeconds,
    detail: failed ? errorMessage : 'OpenClaw 已返回结果，消息已完成汇总',
  };
  refreshAgentChat(index);
  window.setTimeout(() => {
    const current = state.agentChatProgress[agentID];
    if (current && (current.status === 'failed' || current.status === 'done')) {
      delete state.agentChatProgress[agentID];
      if (state.activeAgentChatID === agentID) refreshAgentChat(index);
    }
  }, 5000);
}

async function sendAgentMessage(index, providedMessage='') {
  const input = document.getElementById('chatInput');
  const agentID = agentIDs[index];
  if (state.agentChatPending[agentID]) {
    toast('上一条消息仍在处理','请等待当前回复完成后再发送下一条消息，或点击停止取消当前调用。','warning');
    input?.focus();
    return;
  }
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
  const queueNotice = await assistantQueueNotice(index);
  if (queueNotice) toast('Agent 任务已排队', queueNotice.message, 'warning');
  const model = normalizeAgentMessageModel(agentSelectedModel(agentID));
  const sources = ['Agent 知识库'];
  const allowlist = getAgentAllowlist(agentID);
  const history = state.agentChats[agentID] || (state.agentChats[agentID] = []);
  const pendingAttachments = state.chatAttachmentsByAgent[agentID] || [];
  const payloadAttachments = payloadChatAttachments(pendingAttachments);
  const displayAttachments = cloneChatAttachments(pendingAttachments);
  state.chatAttachmentsByAgent[agentID] = [];
  const userMessage = {role:'user',text:message,sources,attachments:displayAttachments,attachmentStatus:displayAttachments.length?'正在提交给 OpenClaw':'',time:formatLocalizedTime(new Date().toISOString())};
  history.push(userMessage);
  state.agentChatAtBottom[agentID] = true;
  if (input && !providedMessage) input.value = '';
  const requestStartedAt = Date.now();
  state.agentChatPending[agentID] = true;
  const controller = new AbortController();
  state.agentChatControllers[agentID] = controller;
  startAgentProgress(agentID,index,payloadAttachments.length>0,requestStartedAt);
  markAgentProgressQueued(agentID,index,queueNotice);
  refreshAgentChat(index, {forceBottom:true});
  applyIcons();
  try {
    const requestBody = {page:'agents',feature:'agent-chat',message,attachments:payloadAttachments,sessionKey:`sta100-${agentID}`,context:{targetAgent:agentID,allowlist,lang:state.lang}};
    if (model) requestBody.model = model;
    const result = await apiFetch('/api/v1/assistant/query', {method:'POST', signal:controller.signal, body:JSON.stringify(requestBody)});
    if (result.queue?.queued && result.queue.message) toast('Agent 排队后已开始执行', result.queue.message, 'warning');
    userMessage.attachmentStatus = displayAttachments.length ? `已随消息提交给 OpenClaw（${result.attachments?.length || displayAttachments.length} 个附件）` : '';
    applyTokenUsage(result.tokenUsage);
    history.push({role:'agent',text:result.text || 'OpenClaw 未返回文本内容。',sources:result.usedAgents||[],pipeline:result.pipeline||[],totalDurationMs:result.totalDurationMs||0,time:formatLocalizedTime(new Date().toISOString())});
    finishAgentProgress(agentID,index,result,'',requestStartedAt);
  } catch (error) {
    userMessage.attachmentStatus = displayAttachments.length ? '发送失败' : '';
    const aborted = isAbortError(error) || controller.signal.aborted;
    const text = aborted ? '当前 Agent 调用已停止。' : `调用失败：${error.message}`;
    history.push({role:'agent',text,error:!aborted,retry:aborted?'':message,sources,pipeline:[{stage:'openclaw-agent',status:aborted?'cancelled':'failed',detail:aborted?'用户主动停止当前 Agent 调用':error.message}],time:formatLocalizedTime(new Date().toISOString())});
    finishAgentProgress(agentID,index,null,aborted?'已停止当前 Agent 调用':error.message,requestStartedAt);
    if (aborted) toast('已停止当前 Agent 调用','只取消当前智能体本次消息，不影响其它智能体。','success');
  } finally {
    state.agentChatPending[agentID] = false;
    if (state.agentChatControllers[agentID] === controller) delete state.agentChatControllers[agentID];
  }
  refreshAgentChat(index, {forceBottom:true});
}

async function confirmAgentAnswer(index, messageIndex) {
  const agentID = agentIDs[index];
  const history = state.agentChats[agentID] || [];
  const answer = history[Number(messageIndex)];
  if (!answer || answer.role === 'user' || !String(answer.text || '').trim()) return;
  let question = '';
  for (let cursor = Number(messageIndex) - 1; cursor >= 0; cursor--) {
    if (history[cursor]?.role === 'user') { question = history[cursor].text; break; }
  }
  if (!question) { toast('无法入库案例','未找到对应的用户问题。','warning'); return; }
  try {
    await apiFetch('/api/v1/assistant/cases',{method:'POST',body:JSON.stringify({agentId:agentID,title:question,question,answer:answer.text})});
    answer.caseConfirmed = true;
    refreshAgentChat(index);
    toast('案例已入库',`已保存到 ${agentID} 的本机知识库。`,'success');
  } catch(error) { toast('案例入库失败',error.message,'warning'); }
}

function stopAgentMessage(index) {
  const agentID = agentIDs[index];
  const controller = state.agentChatControllers[agentID];
  if (!state.agentChatPending[agentID] || !controller) {
    toast('当前没有可停止的消息','该智能体没有正在调用的消息。','warning');
    return;
  }
  controller.abort();
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
  const selectedCountry = c.country || '德国';
  const selectedCity = c.city || (cityOptions(selectedCountry)[0]?.value || '');
  state.formContext = { type: 'customer', id: c.id || '' };
  openModal({title:customer?'编辑客户':'新建客户',eyebrow:'客户档案',body:`<div class="form-grid"><div class="form-field full"><label>从名片或照片识别</label><button class="upload-zone" data-action="mock-ocr" style="min-height:92px">${icon('scan-line')}<span>选择图片后识别并填充字段</span></button></div><div class="form-section"><h3>基本信息</h3><p>列表字段可配置，详情页保留全部字段。</p></div>${inputField('客户名称',c.name||'',true,false,'text','customerName')}${selectField('客户类型',customerTypeOptions(),false,'customerTypeForm',c.type||'Customer')}${inputField('主电话',c.phone||'',true,false,'tel','customerPhone')}${inputField('网站',c.website||'',false,false,'url','customerWebsite')}${selectField('国家',countryOptions(),false,'customerCountryForm',selectedCountry)}${selectField('城市',cityOptions(selectedCountry),false,'customerCity',selectedCity)}${inputField('联系人',c.contact||'',false,false,'text','customerContact')}${inputField('联系邮箱',c.email||'',false,false,'email','customerEmail')}${selectField('评级',customerRatingOptions(),false,'customerRating',c.rating||'Prospect')}${selectField('来源',['展会','电话','朋友介绍','拜访','互联网线索','客户转介绍','线索','其它'],false,'customerSource',c.source||'其它')}${inputField('负责人',c.owner||'Donald',false,false,'text','customerOwner')}${inputField('描述',c.description||'',false,true,'text','customerDescription')}${customer?`<div class="form-section"><h3>历史沟通记录</h3><p>沟通记录独立保存且只能追加，编辑客户不会覆盖历史。</p></div><div class="form-field full"><button type="button" class="button" data-action="customer-communications" data-id="${escapeAttr(c.id)}">${icon('messages-square')}查看或新增沟通记录</button></div>`:''}</div>`,footer:formFooter(customer?'保存修改':'创建客户','save-customer')});
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
  else if (tab === 'contacts') content = `<div class="detail-grid">${[['联系人',c.contact],['电话',c.phone],['邮箱',c.email],['网站',c.website||'未填写'],['国家/城市',`${localizedCountry(c.country)}${c.city?' / '+c.city:''}`],['来源',c.source||'未填写']].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div>`;
  else if (tab === 'activity') content = `<div class="spread communication-head"><div><strong>历史沟通记录</strong><p class="secondary-text">记录只能追加，不能修改或删除。</p></div><button class="button primary small" data-action="new-customer-communication" data-id="${escapeAttr(c.id)}">${icon('message-square-plus')}新增沟通</button></div>${communications.length?`<div class="timeline communication-timeline">${communications.map(item=>`<div class="timeline-item"><div class="spread"><h4>${escapeHTML(item.subject||item.type)}</h4><span class="badge blue">${escapeHTML(item.type)}</span></div><p class="communication-content">${escapeHTML(item.content)}</p><small>${escapeHTML(String(item.occurredAt||'').replace('T',' '))}${item.contact?` · ${escapeHTML(item.contact)}`:''} · 由 ${escapeHTML(item.createdBy)} 记录</small></div>`).join('')}</div>`:`<div class="empty-state panel">${icon('messages-square')}<div><h3>暂无沟通记录</h3><p>新增后将永久保留在本机业务数据库中。</p></div></div>`}`;
  else content = `<div class="detail-grid">${[['客户编号',c.id],['客户类型',localizedCustomerType(c.type)],['国家',localizedCountry(c.country)],['城市',c.city||'未填写'],['评级',localizedCustomerRating(c.rating)],['负责人',c.owner],['主联系人',c.contact],['电话',c.phone],['邮箱',c.email],['来源',c.source||'未填写'],['订单数',String(c.orders)],['累计金额',c.total],['最近更新',c.updated]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 14px">最近业务记录</div><div class="filter-row"><button class="button small" data-customer-tab="quotes" data-customer-id="${escapeAttr(c.id)}">报价单 ${relatedQuotes.length}</button><button class="button small" data-customer-tab="orders" data-customer-id="${escapeAttr(c.id)}">订单 ${relatedOrders.length}</button><button class="button small" data-customer-tab="documents" data-customer-id="${escapeAttr(c.id)}">单据 ${relatedDocuments.length}</button></div>`;
  openDrawer({title:c.name,eyebrow:`客户 / ${c.id}`,body:`${tabBar}<div class="spread" style="margin-bottom:14px"><span>${badge(c.rating)}</span><div class="inline-actions"><button class="button small" data-action="edit-customer" data-id="${escapeAttr(c.id)}">${icon('pencil')}编辑</button><button class="button primary small" data-action="new-quote" data-customer="${escapeAttr(c.name)}">${icon('file-plus-2')}新建报价</button><button class="button small" data-action="customer-to-lead" data-id="${escapeAttr(c.id)}">${icon('arrow-left')}转为线索</button><button class="button danger small" data-action="delete-customer" data-id="${escapeAttr(c.id)}">${icon('trash-2')}删除</button></div></div>${content}`});
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

function supplierCommunicationForm(supplierID) {
  const supplier=suppliers.find(item=>item.id===supplierID);
  if(!supplier)return;
  closeDrawer();
  const now=new Date(Date.now()-new Date().getTimezoneOffset()*60000).toISOString().slice(0,16);
  openModal({title:'新增沟通记录',eyebrow:`${supplier.company} / 只能追加`,body:`<div class="form-grid">${selectField('沟通方式',['电话','邮件','微信','会议','拜访','短信','其它'],false,'supplierCommType','电话')}${inputField('沟通时间',now,true,false,'datetime-local','supplierCommOccurredAt')}${inputField('主题','',false,true,'text','supplierCommSubject')}${inputField('联系人',supplier.contact||'',false,true,'text','supplierCommContact')}<div class="form-field full"><label>沟通内容 <span class="required">*</span></label><textarea class="textarea" id="supplierCommContent" maxlength="10000" placeholder="记录供应商反馈、结论和后续事项"></textarea><small>保存后不可修改或删除，请确认内容准确。</small></div></div>`,footer:`<button class="button" data-action="cancel-supplier-communication" data-id="${escapeAttr(supplierID)}">取消</button><button class="button primary" data-action="save-supplier-communication" data-id="${escapeAttr(supplierID)}">${icon('plus')}追加记录</button>`});
}

async function saveSupplierCommunication(supplierID) {
  const content=formText('supplierCommContent');
  const occurredAt=formText('supplierCommOccurredAt');
  if(!content||!occurredAt){toast('保存失败','沟通时间和沟通内容为必填项。','warning');return;}
  const payload={type:formText('supplierCommType'),subject:formText('supplierCommSubject'),content,contact:formText('supplierCommContact'),occurredAt};
  try {
    const record=await apiFetch(`/api/v1/suppliers/${encodeURIComponent(supplierID)}/communications`,{method:'POST',body:JSON.stringify(payload)});
    state.supplierCommunications[supplierID]=[record,...(state.supplierCommunications[supplierID]||[])];
    closeModal();
    await supplierDetail(supplierID,'activity');
    toast('沟通记录已追加','该记录已永久保留，不能修改或删除。');
  } catch(error) { toast('沟通记录保存失败',error.message,'warning'); }
}

async function newQuoteForm(quote, customerName='') {
  await loadTemplates('quote', true);
  const q=quote||{}; state.formContext={type:'quote',id:q.id||''};
  state.quoteDraftLines = businessLines(q).map(line => ({ ...line, discount: Number(line.discount || 0) }));
  if (!state.quoteDraftLines.length) state.quoteDraftLines = [{ productId: products.find(product => product.status === 'Active')?.id, quantity: 1, unitPrice: moneyNumber(products[0]?.price), discount: 0 }];
  const selectedCustomer = customers.find(c=>!c.archived && c.name === (customerName||q.customer)) || customers.find(c=>!c.archived);
  const templateFields = quoteTemplateFieldDefaults(q, selectedCustomer);
  openModal({title:quote?'编辑报价单':'新建报价单',eyebrow:'报价单 / 标准 PDF 模板',wide:true,body:`<div class="form-grid">${inputField('报价主题',q.subject||'欧洲渠道设备报价',true,false,'text','quoteSubject')}${relationField('关联客户','quoteCustomer',customers.filter(c=>!c.archived).map(c=>c.name),customerName||q.customer||selectedCustomer?.name||'')}${selectField('状态',[{value:'Draft',label:'草稿 Draft'},{value:'Delivered',label:'已发送 Delivered'},{value:'Accepted',label:'已接受 Accepted'},{value:'Rejected',label:'已拒绝 Rejected'}],false,'quoteStatusForm',q.status||'Draft')}${inputField('有效期',q.valid||'2026-09-10',true,false,'date','quoteValid')}${selectField('币种',[{value:'EUR',label:'EUR / 欧元'},{value:'USD',label:'USD / 美元'},{value:'CNY',label:'CNY / 人民币'},{value:'GBP',label:'GBP / 英镑'}],false,'quoteCurrency',q.currency||'EUR')}<div class="form-section"><h3>产品明细</h3><p>产品、库存和默认单价来自产品库；保存后写入报价快照，后续转订单时继续沿用。</p></div><div class="form-field full"><div class="data-wrap"><table class="data-table line-editor-table" style="min-width:820px"><thead><tr><th>产品</th><th>数量</th><th>单价</th><th>折扣 %</th><th>小计</th><th></th></tr></thead><tbody id="quoteLinesBody"></tbody><tfoot><tr><td colspan="4" style="text-align:right">产品小计</td><td><strong id="quoteDraftSubtotal">${formatMoney(quoteDraftSubtotal(),q.currency||'EUR')}</strong></td><td></td></tr><tr><td colspan="4" style="text-align:right">运费</td><td><strong id="quoteDraftFreight">${formatMoney(q.freight||0,q.currency||'EUR')}</strong></td><td></td></tr><tr><td colspan="4" style="text-align:right">税费</td><td><strong id="quoteDraftTax">${formatMoney(q.tax||0,q.currency||'EUR')}</strong></td><td></td></tr><tr><td colspan="4" style="text-align:right"><strong>报价合计（不含税）</strong></td><td><strong id="quoteDraftTotal">${formatMoney(quoteDraftTotal(),q.currency||'EUR')}</strong></td><td></td></tr></tfoot></table></div><button type="button" class="button ghost small" data-action="add-quote-line" style="margin-top:8px">${icon('plus')}添加产品</button></div>${inputField('运费（本单）',q.freight===undefined?'':q.freight,false,false,'number','quoteFreight')}${inputField('税费',q.tax===undefined?'':q.tax,false,false,'number','quoteTax')}${inputField('条款与条件',q.terms||templateFields['quote.paymentTerms'],false,true,'text','quoteTerms')}<div class="form-section full"><h3>标准模板填充字段</h3><p>客户和产品从业务数据实时读取；下面字段来自标准模板默认值，也可针对本张报价单单独调整。</p></div>${inputField('客户国家',templateFields['customer.country'],false,false,'text','quoteCustomerCountry')}${inputField('客户联系人',templateFields['customer.contact'],false,false,'text','quoteCustomerContact')}${inputField('客户邮箱',templateFields['customer.email'],false,false,'email','quoteCustomerEmail')}${inputField('客户地址',templateFields['customer.address'],false,true,'text','quoteCustomerAddress')}${inputField('价格条件',templateFields['quote.priceTerms'],false,false,'text','quotePriceTerms')}${inputField('付款条件',templateFields['quote.paymentTerms'],false,false,'text','quotePaymentTerms')}${inputField('运输说明',templateFields['quote.shipping'],false,false,'text','quoteShipping')}${inputField('目的港/目的地',templateFields['quote.destination'],false,false,'text','quoteDestination')}${inputField('交期',templateFields['quote.leadTime'],false,false,'text','quoteLeadTime')}${inputField('质保',templateFields['quote.warranty'],false,false,'text','quoteWarranty')}${inputField('认证',templateFields['quote.certification'],false,false,'text','quoteCertification')}${inputField('包装',templateFields['quote.packaging'],false,false,'text','quotePackaging')}${inputField('备注',templateFields['quote.note'],false,true,'text','quoteNote')}</div>`,footer:formFooter(quote?'保存修改':'保存草稿','save-quote')});
  renderQuoteDraftLines();
}

function quoteTemplateFieldDefaults(q={}, customer={}) {
  const fields = q.templateFields || {};
  const tpl=state.templates.find(t=>t.kind==='quote'&&t.default) || state.templates.find(t=>t.kind==='quote') || {};
  const defaults=quoteTemplateDefaultValuesJS(tpl.defaultValues);
  const pick = (key, fallback='') => fields[key] || fallback || '';
  return {
    'customer.country': pick('customer.country', customer?.country || ''),
    'customer.contact': pick('customer.contact', customer?.contact || ''),
    'customer.email': pick('customer.email', customer?.email || ''),
    'customer.address': pick('customer.address', customer?.description || ''),
    'quote.priceTerms': pick('quote.priceTerms', defaults['quote.priceTerms']),
    'quote.paymentTerms': pick('quote.paymentTerms', q.terms || defaults['quote.paymentTerms']),
    'quote.shipping': pick('quote.shipping', defaults['quote.shipping']),
    'quote.destination': pick('quote.destination', customer?.city || customer?.country || defaults['quote.destination']),
    'quote.leadTime': pick('quote.leadTime', defaults['quote.leadTime']),
    'quote.warranty': pick('quote.warranty', defaults['quote.warranty']),
    'quote.certification': pick('quote.certification', defaults['quote.certification']),
    'quote.packaging': pick('quote.packaging', defaults['quote.packaging']),
    'quote.note': pick('quote.note', defaults['quote.note']),
  };
}

function syncQuoteTemplateCustomerFields(customerName='') {
  const customer = customers.find(c=>!c.archived && c.name === customerName);
  if (!customer) return;
  const mapping = {
    quoteCustomerCountry: customer.country || '',
    quoteCustomerContact: customer.contact || '',
    quoteCustomerEmail: customer.email || '',
    quoteCustomerAddress: customer.description || '',
    quoteDestination: customer.city || customer.country || '',
  };
  Object.entries(mapping).forEach(([id,value])=>{
    const input = document.getElementById(id);
    if (input && !input.value.trim()) input.value = value;
  });
}

function quoteDetail(id) {
  const q=quotes.find(x=>x.id===id);
  if (!q) return;
  openDrawer({title:q.id,eyebrow:'报价单详情',body:`<div class="spread"><span>${badge(q.status)}</span><div class="inline-actions"><button class="button small" data-action="download-quote" data-id="${escapeAttr(q.id)}">${icon('download')}PDF</button><button class="button primary small" data-action="send-quote-email" data-id="${escapeAttr(q.id)}">${icon('mail')}发送邮件</button><button class="button small" data-action="edit-quote" data-id="${escapeAttr(q.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-quote" data-id="${escapeAttr(q.id)}">${icon('trash-2')}删除</button>${q.status==='Accepted'?`<button class="button primary small" data-action="convert-order" data-id="${escapeAttr(q.id)}">${icon('arrow-right')}转为订单</button>`:''}</div></div><div class="detail-grid" style="margin-top:15px">${[['主题',q.subject],['客户',q.customer],['产品',q.products],['总金额',q.value],['有效期',q.valid],['更新时间',q.updated],['负责人',q.owner],['状态',q.status]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 14px">状态记录</div><div class="timeline"><div class="timeline-item"><h4>${q.status==='Accepted'?'客户接受报价':'最近更新报价'}</h4><p>${escapeHTML(q.updated || '暂无更新时间')} · ${escapeHTML(q.owner || '系统')}</p></div><div class="timeline-item"><h4>创建报价草稿</h4><p>产品价格来自默认价格表</p></div></div>`});
}

function newOrderForm(order, quoteId='') {
  const o=order||{}; const quote=quotes.find(q=>q.id===(quoteId||o.quote)); state.formContext={type:'order',id:o.id||''};
  state.orderDraftLines = quote ? quoteLinesToOrder(businessLines(quote)) : businessLines(o);
  if (!state.orderDraftLines.length) { const product=products.find(item=>item.status==='Active')||products[0]; state.orderDraftLines=[{productId:product.id,quantity:1,unitPrice:moneyNumber(product.price)}]; }
  openModal({title:order?'编辑订单':'新建订单',eyebrow:'订单 / 草稿',wide:true,body:`<div class="form-grid">${selectField('来源方式',[{value:'从已接受报价创建',label:'从已接受报价创建'},{value:'手动创建',label:'手动创建'}],false,'orderSource',quote?'从已接受报价创建':(order?.quote?'从已接受报价创建':'手动创建'))}${relationField('关联报价','orderQuote',quotes.map(q=>`${q.id} · ${q.customer}`),quote?`${quote.id} · ${quote.customer}`:(o.quote?`${o.quote} · ${o.customer||''}`:''))}${relationField('客户','orderCustomer',customers.filter(c=>!c.archived).map(c=>c.name),o.customer||quote?.customer||customers.find(c=>!c.archived)?.name||'')}${selectField('状态',[{value:'Draft',label:'草稿 Draft'},{value:'Confirmed',label:'已确认 Confirmed'},{value:'Paid',label:'已付款 Paid'},{value:'Production',label:'生产中 Production'},{value:'Shipped',label:'已发运 Shipped'},{value:'Completed',label:'已完成 Completed'},{value:'Cancelled',label:'已取消 Cancelled'}],false,'orderStatusForm',o.status||'Confirmed')}${inputField('客户 PO 号',o.po||'',false,false,'text','orderPO')}${inputField('预计交付日期',o.delivery||'2026-09-15',true,false,'date','orderDelivery')}${selectField('贸易条款',[{value:'FOB Shenzhen',label:'FOB Shenzhen'},{value:'CIF Hamburg',label:'CIF Hamburg'},{value:'DAP Customer',label:'DAP Customer'}],false,'orderTerms',o.terms||'FOB Shenzhen')}<div class="form-section"><h3>订单明细</h3><p>产品选择、数量、库存和成交单价均来自产品库；保存为订单明细快照，生成单据时继续读取本快照。</p></div><div class="form-field full"><div class="data-wrap"><table class="data-table line-editor-table" style="min-width:820px"><thead><tr><th>产品</th><th>数量</th><th>成交单价</th><th>小计</th><th></th></tr></thead><tbody id="orderLinesBody"></tbody><tfoot><tr><td colspan="3" style="text-align:right"><strong>订单合计</strong></td><td><strong id="orderDraftTotal">${formatMoney(orderDraftTotal())}</strong></td><td></td></tr></tfoot></table></div><button type="button" class="button ghost small" data-action="add-order-line" style="margin-top:8px">${icon('plus')}添加产品</button></div></div>`,footer:formFooter(order?'保存修改':'创建订单','save-order')});
  renderOrderDraftLines();
}

function orderDetail(id) {
  const o=orders.find(x=>x.id===id);
  if (!o) return;
  openDrawer({title:o.id,eyebrow:'订单详情',body:`<div class="spread"><span>${badge(o.status)}</span><div class="inline-actions"><button class="button primary small" data-action="send-order-email" data-id="${escapeAttr(o.id)}">${icon('mail')}发送邮件</button><button class="button small" data-action="edit-order" data-id="${escapeAttr(o.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-order" data-id="${escapeAttr(o.id)}">${icon('trash-2')}删除</button><button class="button primary small" data-action="generate-docs" data-id="${escapeAttr(o.id)}">${icon('files')}生成单据</button></div></div><div class="detail-grid" style="margin-top:15px">${[['客户',o.customer],['来源报价',o.quote||'手动创建'],['产品明细数',String(businessLines(o).length)],['订单金额',o.value],['预计交付',o.delivery],['更新时间',o.updated],['状态',o.status]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 14px">订单产品明细</div>${detailLinesTable(businessLines(o),true)}<div class="divider-title" style="margin:20px 0 14px">订单流转</div><div class="timeline"><div class="timeline-item"><h4>${o.status==='Shipped'?'已发运':'当前阶段：'+o.status}</h4><p>${escapeHTML(o.updated || '暂无更新时间')} · 状态由本系统业务操作维护</p></div><div class="timeline-item"><h4>订单确认</h4><p>客户、产品与交付信息已保存为订单快照</p></div><div class="timeline-item"><h4>${o.quote?'报价转订单':'手动创建订单'}</h4><p>${o.quote?`${escapeHTML(o.quote)} · 保留完整来源`:'直接从产品库选择产品'}</p></div></div>`});
}

function emailDefaultBody(kind, record, language='en') {
  const customer=customers.find(c=>!c.archived&&c.name===record.customer)||{};
  const number=record.id||''; const amount=record.value||''; const date=record.updated||new Date().toISOString().slice(0,10);
  if (language==='zh') return `尊敬的 ${customer.name||record.customer}，\n\n感谢您的支持。附件为 ${kind==='quote'?'报价单':'订单及相关单据'}，请查收。\n\n单据号：${number}\n金额：${amount}\n日期：${date}\n\n如有任何问题，请与我们联系。\n\n此致\nStratronix\nStratronix Technology (Shenzhen) Company, Limited`;
  return `Dear ${customer.name||record.customer},\n\nThank you for your business. Please find attached your ${kind==='quote'?'quotation':'order documents'} for review, printing and download.\n\n${kind==='quote'?'Quotation':'Order'} No. ${number}\nAmount ${amount}\nDate ${date}\n\nPlease contact us if you have any questions.\n\nRegards,\nStratronix\nStratronix Technology (Shenzhen) Company, Limited`;
}

function openEmailSettings(afterOpen) {
  apiFetch('/api/v1/email').then(settings=>{
    openModal({title:'邮件发送设置',eyebrow:'SMTP / 发件人',wide:true,body:`<div class="form-grid">${inputField('SMTP 主机',settings.host||'',true,false,'text','emailHost')}${inputField('端口',String(settings.port||587),true,false,'number','emailPort')}${selectField('安全方式',[{value:'starttls',label:'STARTTLS（推荐）'},{value:'tls',label:'TLS / SSL'},{value:'none',label:'无加密'}],false,'emailSecurity',settings.security||'starttls')}${inputField('发件人名称',settings.fromName||'Stratronix',true,false,'text','emailFromName')}${inputField('发件人邮箱',settings.fromEmail||'',true,false,'email','emailFrom')}${inputField('SMTP 用户名',settings.username||'',false,false,'text','emailUsername')}${inputField('应用专用密码', '',false,false,'password','emailPassword')}<div class="form-section full"><p>支持国内外标准 SMTP 邮箱。多数邮箱需要应用专用密码；密码仅保存于本机。</p></div></div>`,footer:`<button class="button" data-action="cancel-email-settings">取消</button><button class="button primary" data-action="save-email-settings" data-after="${afterOpen||''}">${icon('save')}保存</button>`});
  }).catch(error=>toast('读取邮件设置失败',error.message,'warning'));
}

async function saveEmailSettings(button) {
  const payload={host:formText('emailHost'),port:Number(formText('emailPort')||587),security:formText('emailSecurity'),fromName:formText('emailFromName'),fromEmail:formText('emailFrom'),username:formText('emailUsername'),password:formText('emailPassword')};
  button.disabled=true;button.innerHTML=`${icon('loader-circle')}保存中`;applyIcons();
  try { await apiFetch('/api/v1/email',{method:'PATCH',body:JSON.stringify(payload)}); closeModal(); toast('邮件设置已保存','SMTP 发件人已更新。','success'); }
  catch(error){button.disabled=false;button.innerHTML=`${icon('save')}保存`;applyIcons();toast('邮件设置保存失败',error.message,'warning');}
}

function openBusinessEmail(kind,id) {
  const record=(kind==='quote'?quotes:orders).find(item=>item.id===id); if(!record)return;
  const customer=customers.find(c=>!c.archived&&c.name===record.customer)||{};
  const relatedDocs=kind==='order'?documents.filter(doc=>doc.order===record.id):[];
  const body=emailDefaultBody(kind,record,'en');
  openModal({title:kind==='quote'?'发送报价单邮件':'发送订单邮件',eyebrow:`${kind==='quote'?'报价单':'订单'} / ${record.id}`,wide:true,body:`<div class="form-grid">${inputField('收件人',customer.email||'',true,false,'email','businessEmailTo')}${inputField('抄送', '',false,false,'text','businessEmailCc')}${selectField('邮件语言',[{value:'en',label:'English'},{value:'zh',label:'中文'}],false,'businessEmailLanguage','en')}${inputField('主题',`${kind==='quote'?'Quotation':'Order'} #${record.id}`,true,false,'text','businessEmailSubject')}<div class="form-field full"><label>邮件正文 <span class="required">*</span></label><textarea class="textarea" id="businessEmailBody" rows="12">${escapeHTML(body)}</textarea></div>${kind==='quote'?`<label class="check-field full"><input type="checkbox" id="businessAttachRecord" checked> 附加当前报价单 PDF</label>`:`<div class="form-field full"><label>附加已生成单据</label><div class="check-list">${relatedDocs.map(doc=>`<label class="check-field"><input type="checkbox" name="businessDocumentId" value="${escapeAttr(doc.id)}" checked> ${escapeHTML(doc.type)} · ${escapeHTML(doc.id)}（PDF）</label>`).join('')||'<small>当前订单尚未生成单据，可先生成 PI、CI、PL、合同或报关单。</small>'}</div></div>`}</div>`,footer:`<button class="button" data-action="open-email-settings">${icon('settings')}发件人设置</button><button class="button" data-action="close-modal">取消</button><button class="button primary" data-action="confirm-business-email" data-kind="${kind}" data-id="${escapeAttr(id)}">${icon('send')}发送</button>`});
}

async function submitBusinessEmail(kind,id,button) {
  const to=formText('businessEmailTo').split(/[;,，\s]+/).filter(Boolean); const cc=formText('businessEmailCc').split(/[;,，\s]+/).filter(Boolean);
  const payload={to,cc,subject:formText('businessEmailSubject'),body:formText('businessEmailBody'),language:formText('businessEmailLanguage'),attachRecord:Boolean(document.getElementById('businessAttachRecord')?.checked),documentIds:[...document.querySelectorAll('input[name="businessDocumentId"]:checked')].map(input=>input.value)};
  if(!to.length||!payload.subject||!payload.body){toast('邮件内容不完整','请填写收件人、主题和正文。','warning');return;}
  button.disabled=true; button.innerHTML=`${icon('loader-circle')}发送中`;applyIcons();
  try { const result=await apiFetch(`/api/v1/${kind==='quote'?'quotes':'orders'}/${encodeURIComponent(id)}/deliver`,{method:'POST',body:JSON.stringify(payload)}); if(kind==='quote'&&result.quote)upsertRecord(quotes,result.quote);if(kind==='order'&&result.order)upsertRecord(orders,result.order);closeModal();renderPage();toast('邮件已发送',result.message||'已通过 SMTP 投递。','success'); } catch(error){button.disabled=false;button.innerHTML=`${icon('send')}发送`;applyIcons();toast('邮件发送失败',error.message,'warning');}
}

function generateDocs(orderId='', document) {
  const d=document||{}; state.formContext={type:'document',id:d.id||''};
  openModal({title:document?'编辑单据':'生成外贸单据',eyebrow:orderId==='all'?'一键生成全套':`订单 / ${orderId||d.order||''}`,body:`<div class="form-grid">${relationField('关联订单','documentOrder',orders.map(o=>`${o.id} · ${o.customer}`),d.order?`${d.order} · ${d.customer}`:orderId&&orderId!=='all'?`${orderId} · ${orders.find(o=>o.id===orderId)?.customer||''}`:'')}${selectField('单据类型',document?['PI','CI','PL','报关单','合同']:['PI','CI','PL','报关单','合同','全套（PI + CI + PL + 报关单 + 合同）'],false,'documentTypeForm',d.type||(orderId==='all'?'全套（PI + CI + PL + 报关单 + 合同）':'PI'))}${selectField('模板版本',['系统内置标准模板'],false,'documentTemplate',d.template||'系统内置标准模板')}${selectField('输出语言',['英文','中文 / 英文双语'],false,'documentLanguage','英文')}<div class="form-section full"><h3>生成方式</h3><p>选择订单和单据类型，系统自动填充字段；下载时选择所需格式。</p></div></div>`,footer:formFooter(document?'保存修改':'生成并预览','save-document')});
}

async function saveCustomer() {
  const name=formText('customerName');
  if (!name || !formText('customerCountryForm')) { toast('保存失败','客户名称和国家为必填项。','warning'); return; }
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
  const templateFields={
    'customer.country': formText('quoteCustomerCountry'),
    'customer.contact': formText('quoteCustomerContact'),
    'customer.email': formText('quoteCustomerEmail'),
    'customer.address': formText('quoteCustomerAddress'),
    'quote.priceTerms': formText('quotePriceTerms'),
    'quote.paymentTerms': formText('quotePaymentTerms') || formText('quoteTerms'),
    'quote.shipping': formText('quoteShipping'),
    'quote.destination': formText('quoteDestination'),
    'quote.leadTime': formText('quoteLeadTime'),
    'quote.warranty': formText('quoteWarranty'),
    'quote.certification': formText('quoteCertification'),
    'quote.packaging': formText('quotePackaging'),
    'quote.note': formText('quoteNote'),
  };
  const payload={...(existing||{}),subject,customer,valid:formText('quoteValid'),currency,status:formText('quoteStatusForm')||existing?.status||'Draft',freight:formNumber('quoteFreight'),tax:formNumber('quoteTax'),terms:formText('quoteTerms'),templateId:'',templateFields,lines:state.quoteDraftLines.map(line=>({...line,quantity:Number(line.quantity),unitPrice:Number(line.unitPrice),discount:Number(line.discount||0) }))};
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
  const quoteId=quoteValue.split(' · ')[0]||quoteValue;
  const matchedQuote=quotes.find(q=>q.id===quoteId);
  if (quoteId && !matchedQuote) { toast('保存失败','请从报价搜索结果中选择有效报价。','warning'); return; }
  const linkedCustomer = matchedQuote?.customer || customer;
  const payload={...(existing||{}),customer:linkedCustomer,quote:matchedQuote?.id||quoteId,po:formText('orderPO'),delivery:formText('orderDelivery'),terms:formText('orderTerms'),currency:existing?.currency||'EUR',status:formText('orderStatusForm')||existing?.status||'Confirmed',progress:Number(existing?.progress||0),lines:state.orderDraftLines.map(line=>({...line,quantity:Number(line.quantity),unitPrice:Number(line.unitPrice)}))};
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
      const result=await apiFetch(`/api/v1/orders/${encodeURIComponent(order.id)}/documents`,{method:'POST',body:JSON.stringify({types:['PI','CI','PL','报关单','合同'],template:payload.template,language:payload.language})});
      result.items.forEach(record=>upsertRecord(documents,record));closeModal();renderPage();toast('全套单据已生成',`${result.total} 份单据已进入待复核状态。`);return;
    }
    const record=await apiFetch(existing?`/api/v1/documents/${encodeURIComponent(existing.id)}`:'/api/v1/documents',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)});
    upsertRecord(documents,record); closeModal(); renderPage(); toast(existing?'单据已更新':'单据已生成',`${record.id} 已保存为待复核状态。`);
  } catch(error) { toast('保存失败',error.message,'warning'); }
}

async function deleteCustomer(id) {
  const customer=customers.find(c=>c.id===id); if (!customer) return;
  const choice = await confirmChoice({
    title: `删除客户 ${customer.name}`,
    question: `确定要删除客户 "${customer.name}" 吗？`,
    description: '选择"是"：软删除（归档），关联业务记录仍可用于历史查询；选择"否"：硬删除（物理删除），关联业务记录也会被一并清理。',
    yesLabel: '是，软删除',
    noLabel: '否，硬删除',
  });
  if (choice === null) return;
  try {
    await apiFetch(`/api/v1/accounts/${encodeURIComponent(id)}`,{method:'DELETE',body:JSON.stringify({ hardDelete: choice === false })});
    if (choice === false) {
      removeRecord(customers, id);
    } else {
      customer.archived = true;
    }
    closeDrawer();
    renderPage();
    toast(choice === true ? '客户已归档' : '客户已硬删除', choice === true ? '关联业务记录仍然保留，可用于历史查询。' : '客户及其关联业务记录已清理。');
  } catch(error) { toast('删除失败',error.message,'warning'); }
}

function newLeadForm(lead) {
  const l = lead || {};
  state.formContext = { type: 'lead', id: l.id || '' };
  openModal({title: lead ? '编辑线索' : '新建线索', eyebrow: '线索库', body: `<div class="form-grid">${inputField('线索客户名称', l.name || '', true, false, 'text', 'leadName')}${selectField('客户类型', customerTypeOptions(), false, 'leadType', l.type || 'Distributor')}${inputField('国家', l.country || '', true, false, 'text', 'leadCountry')}${inputField('城市', l.city || '', false, false, 'text', 'leadCity')}${inputField('联系人', l.contact || '', false, false, 'text', 'leadContact')}${inputField('电话', l.phone || '', false, false, 'tel', 'leadPhone')}${inputField('邮箱', l.email || '', false, false, 'email', 'leadEmail')}${inputField('官网', l.website || '', false, false, 'url', 'leadWebsite')}${inputField('地址', l.address || '', false, false, 'text', 'leadAddress')}${inputField('业务方向', l.business || '', false, true, 'text', 'leadBusiness')}${inputField('来源', l.source || '', false, false, 'text', 'leadSource')}${inputField('来源链接', l.sourceUrl || '', false, false, 'url', 'leadSourceUrl')}</div>`, footer: formFooter(lead ? '保存修改' : '创建线索', 'save-lead')});
}

async function saveLead() {
  const payload = {
    name: document.getElementById('leadName')?.value?.trim(),
    type: document.getElementById('leadType')?.value,
    country: document.getElementById('leadCountry')?.value?.trim(),
    city: document.getElementById('leadCity')?.value?.trim(),
    contact: document.getElementById('leadContact')?.value?.trim(),
    phone: document.getElementById('leadPhone')?.value?.trim(),
    email: document.getElementById('leadEmail')?.value?.trim(),
    website: document.getElementById('leadWebsite')?.value?.trim(),
    address: document.getElementById('leadAddress')?.value?.trim(),
    business: document.getElementById('leadBusiness')?.value?.trim(),
    source: document.getElementById('leadSource')?.value?.trim(),
    sourceUrl: document.getElementById('leadSourceUrl')?.value?.trim(),
  };
  if (!payload.name) { toast('保存失败', '线索客户名称不能为空', 'warning'); return; }
  try {
    const id = state.formContext.id;
    if (id) {
      const response = await apiFetch(`/api/v1/leads/${encodeURIComponent(id)}`, { method: 'PATCH', body: JSON.stringify(payload) });
      const index = leads.findIndex(l => l.id === id);
      if (index !== -1) leads[index] = { ...leads[index], ...response };
      toast('线索已保存', '变更已写入本地数据库。');
    } else {
      const response = await apiFetch('/api/v1/leads', { method: 'POST', body: JSON.stringify(payload) });
      leads.push(response);
      toast('线索已创建', '新线索已添加到线索库。');
    }
    closeModal();
    renderPage();
  } catch (error) {
    toast('保存失败', error.message, 'warning');
  }
}

async function leadDetail(id, tab='overview') {
  const l = leads.find(x => x.id === id);
  if (!l) return;
  const sourceButton = l.sourceUrl ? `<div class="inline-actions"><button class="button" onclick="window.open('${escapeAttr(l.sourceUrl)}','_blank','noopener,noreferrer')">${icon('external-link')}查看来源</button></div>` : '';
  if (tab === 'activity' && !Array.isArray(state.leadCommunications[id])) {
    try {
      const result = await apiFetch(`/api/v1/leads/${encodeURIComponent(id)}/communications`);
      state.leadCommunications[id] = result.items || [];
    } catch (error) {
      toast('跟进记录加载失败', error.message, 'warning');
      state.leadCommunications[id] = [];
    }
  }
  const tabs = [['overview','概览'],['activity','跟进记录']];
  const tabBar = `<div class="tabs">${tabs.map(([key,label])=>`<button class="${tab===key?'active':''}" data-lead-tab="${key}" data-lead-id="${escapeAttr(l.id)}">${label}</button>`).join('')}</div>`;
  const communications = state.leadCommunications[id] || [];
  const content = tab === 'activity'
    ? `<div class="spread communication-head"><div><strong>历史跟进记录</strong><p class="secondary-text">规则与客户跟进记录一致，只能追加，不能修改或删除。</p></div><button class="button primary small" data-action="new-lead-communication" data-id="${escapeAttr(l.id)}">${icon('message-square-plus')}新增跟进</button></div>${communications.length?`<div class="timeline communication-timeline">${communications.map(item=>`<div class="timeline-item"><div class="spread"><h4>${escapeHTML(item.subject||item.type)}</h4><span class="badge blue">${escapeHTML(item.type)}</span></div><p class="communication-content">${escapeHTML(item.content)}</p><small>${escapeHTML(formatLocalizedDateTime(item.occurredAt)||String(item.occurredAt||''))}${item.contact?` · ${escapeHTML(item.contact)}`:''} · 由 ${escapeHTML(item.createdBy)} 记录</small></div>`).join('')}</div>`:`<div class="empty-state panel">${icon('messages-square')}<div><h3>暂无跟进记录</h3><p>新增后将永久保留，线索与客户互转时会一并复制。</p></div></div>`}`
    : `<div class="detail-grid" style="margin-top:15px">${[['类型', localizedCustomerType(l.type)], ['国家', localizedCountry(l.country)], ['城市', l.city], ['联系人', l.contact], ['电话', l.phone], ['邮箱', l.email], ['官网', l.website], ['地址', l.address], ['业务方向', l.business], ['来源', l.source], ['评分', l.score ? l.score + ' 分' : '待核实'], ['发现理由', l.reason], ['创建时间', formatLocalizedDateTime(l.createdAt)], ['更新时间', formatLocalizedDateTime(l.updated)]].map(([label, value]) => `<div class="detail-field"><label>${label}</label><strong>${escapeHTML(String(value || '未填写'))}</strong></div>`).join('')}</div>${sourceButton}`;
  openDrawer({ title: l.name, eyebrow: `线索 / ${l.id}`, body: `${tabBar}<div class="spread" style="margin-bottom:14px"><span class="badge ${l.converted ? 'green' : 'neutral'}">${l.converted ? '已转换' : '待跟进'}</span><div class="inline-actions"><button class="button" data-action="edit-lead" data-id="${escapeAttr(l.id)}">${icon('pencil')}编辑</button>${!l.converted ? `<button class="button primary" data-action="convert-lead" data-id="${escapeAttr(l.id)}">${icon('arrow-right')}转为客户</button>` : ''}<button class="button danger" data-action="delete-lead" data-id="${escapeAttr(l.id)}">${icon('trash-2')}删除</button></div></div>${content}` });
}

function leadCommunicationForm(leadID) {
  const lead=leads.find(item=>item.id===leadID);
  if(!lead)return;
  closeDrawer();
  const now=new Date(Date.now()-new Date().getTimezoneOffset()*60000).toISOString().slice(0,16);
  openModal({title:'新增跟进记录',eyebrow:`${lead.name} / 只能追加`,body:`<div class="form-grid">${selectField('跟进方式',['电话','邮件','微信','会议','拜访','短信','其它'],false,'leadCommType','电话')}${inputField('跟进时间',now,true,false,'datetime-local','leadCommOccurredAt')}${inputField('主题','',false,true,'text','leadCommSubject')}${inputField('联系人',lead.contact||'',false,true,'text','leadCommContact')}<div class="form-field full"><label>跟进内容 <span class="required">*</span></label><textarea class="textarea" id="leadCommContent" maxlength="10000" placeholder="记录线索反馈、结论和后续事项"></textarea><small>保存后不可修改或删除；转为客户时会一并复制。</small></div></div>`,footer:`<button class="button" data-action="cancel-lead-communication" data-id="${escapeAttr(leadID)}">取消</button><button class="button primary" data-action="save-lead-communication" data-id="${escapeAttr(leadID)}">${icon('plus')}追加记录</button>`});
}

async function saveLeadCommunication(leadID) {
  const content=formText('leadCommContent');
  const occurredAt=formText('leadCommOccurredAt');
  if(!content||!occurredAt){toast('保存失败','跟进时间和跟进内容为必填项。','warning');return;}
  const payload={type:formText('leadCommType'),subject:formText('leadCommSubject'),content,contact:formText('leadCommContact'),occurredAt};
  try {
    const record=await apiFetch(`/api/v1/leads/${encodeURIComponent(leadID)}/communications`,{method:'POST',body:JSON.stringify(payload)});
    state.leadCommunications[leadID]=[record,...(state.leadCommunications[leadID]||[])];
    closeModal();
    await leadDetail(leadID,'activity');
    toast('跟进记录已追加','该记录已永久保留，不能修改或删除。');
  } catch(error) { toast('跟进记录保存失败',error.message,'warning'); }
}

async function deleteLead(id) { if(!window.confirm('确定删除该线索吗？删除后将不再出现在线索列表中。'))return; try { await apiFetch(`/api/v1/leads/${encodeURIComponent(id)}`,{method:'DELETE'}); removeRecord(leads,id); closeDrawer(); renderPage(); toast('线索已删除','变更已写入本地数据库。'); } catch(error) { toast('删除失败',error.message,'warning'); } }

async function customerToLead(id) {
  const c = customers.find(x => x.id === id);
  if (!c) return;
  const choice = await confirmChoice({
    title: `将客户转化为线索`,
    question: `确定将客户 "${c.name}" 转化为线索吗？转化后将创建一条新线索。`,
    description: '选择"是"：保留当前客户并新增线索；选择"否"：转化后直接删除当前客户（关联业务记录会被删除）。',
    yesLabel: '是，保留客户',
    noLabel: '否，删除客户',
  });
  if (choice === null) return;
  const payload = {
    name: c.name,
    type: c.type,
    country: c.country,
    city: c.city,
    contact: c.contact,
    phone: c.phone,
    email: c.email,
    website: c.website,
    address: '',
    business: c.description,
    source: '客户转化线索',
    sourceUrl: '',
    score: 0,
    reason: `由客户 ${c.id} 反向转化`,
  };
  try {
    let communications = state.customerCommunications[id];
    if (!Array.isArray(communications)) {
      try {
        const result = await apiFetch(`/api/v1/accounts/${encodeURIComponent(id)}/communications`);
        communications = result.items || [];
        state.customerCommunications[id] = communications;
      } catch {
        communications = [];
      }
    }
    const lead = await apiFetch('/api/v1/leads', { method: 'POST', body: JSON.stringify(payload) });
    const copied = [];
    for (const item of communications) {
      try {
        const record = await apiFetch(`/api/v1/leads/${encodeURIComponent(lead.id)}/communications`, { method: 'POST', body: JSON.stringify({ type:item.type, subject:item.subject, content:item.content, contact:item.contact, occurredAt:item.occurredAt }) });
        copied.push(record);
      } catch (error) {
        toast('跟进记录复制失败', error.message, 'warning');
      }
    }
    if (copied.length) state.leadCommunications[lead.id] = copied;
    leads.push(lead);
    if (choice === false) {
      try {
        await apiFetch(`/api/v1/accounts/${encodeURIComponent(id)}`, { method: 'DELETE', body: JSON.stringify({ hardDelete: true }) });
      } catch (e) {
        toast('线索已创建，但客户删除失败', e.message, 'warning');
      }
      removeRecord(customers, id);
    }
    toast(choice === true ? '客户已转为线索并保留客户' : '客户已转为线索，客户已删除', `线索编号: ${lead.id}`);
    closeDrawer();
    renderPage();
  } catch (error) {
    toast('转化失败', error.message, 'warning');
  }
}

async function convertLead(id) {
  const lead = leads.find(l => l.id === id);
  if (!lead) return;
  if (lead.converted) { toast('无法转化', '该线索已转化为客户', 'warning'); return; }
  const choice = await confirmChoice({
    title: `将线索转化为客户`,
    question: `确定将线索 "${lead.name}" 转化为客户吗？转化后将创建一个新客户。`,
    description: '选择"是"：保留当前线索并标记为已转化；选择"否"：转化后直接删除当前线索。',
    yesLabel: '是，保留线索',
    noLabel: '否，删除线索',
  });
  if (choice === null) return;
  try {
    const response = await apiFetch(`/api/v1/leads/${encodeURIComponent(id)}/convert`, { method: 'POST', body: JSON.stringify({ keepSource: choice === true }) });
    const customer = response.customer;
    if (choice === true) {
      lead.converted = true;
      lead.convertedAt = customer.createdAt;
      lead.convertedId = customer.id;
    } else {
      removeRecord(leads, id);
    }
    customers.push(customer);
    closeDrawer();
    renderPage();
    toast(choice === true ? '线索已转化为客户并保留线索' : '线索已转化为客户，线索已删除', `客户编号: ${customer.id}`);
  } catch (error) {
    toast('转化失败', error.message, 'warning');
  }
}
async function deleteQuote(id) { if (!window.confirm('确定归档这份报价单吗？')) return; try { await apiFetch(`/api/v1/quotes/${encodeURIComponent(id)}`,{method:'DELETE'}); removeRecord(quotes,id); closeDrawer(); renderPage(); toast('报价单已归档','变更已写入本地数据库。'); } catch(error) { toast('归档失败',error.message,'warning'); } }
async function deleteOrder(id) { if (!window.confirm('确定归档这份订单吗？')) return; try { await apiFetch(`/api/v1/orders/${encodeURIComponent(id)}`,{method:'DELETE'}); removeRecord(orders,id); await loadBusinessData(true); closeDrawer(); renderPage(); toast('订单已归档','变更已写入本地数据库。'); } catch(error) { toast('归档失败',error.message,'warning'); } }
async function deleteDocument(id) { if (!window.confirm('确定归档这份单据吗？')) return; try { await apiFetch(`/api/v1/documents/${encodeURIComponent(id)}`,{method:'DELETE'}); removeRecord(documents,id); closeDrawer(); renderPage(); toast('单据已归档','变更已写入本地数据库。'); } catch(error) { toast('归档失败',error.message,'warning'); } }
async function downloadQuote(id) {
  window.location.assign(`/api/v1/quotes/${encodeURIComponent(id)}/download`);
}
async function downloadOrder(id) {
  window.location.assign(`/api/v1/orders/${encodeURIComponent(id)}/download`);
}
function documentDownloadFormats(documentType) {
  if (documentType === 'PI') return [{ value: 'pdf', label: 'PDF' }];
  if (documentType === '报关单') return [{ value: 'pdf', label: 'PDF' }, { value: 'xml', label: 'XML' }];
  return [{ value: 'pdf', label: 'PDF' }, { value: 'word', label: 'Word' }, { value: 'excel', label: 'Excel（保留标准模板）' }];
}
function downloadDocument(id, format='') {
  const document = documents.find(item => item.id === id);
  if (!document) return;
  if (!format) {
    const options = documentDownloadFormats(document.type);
    if (options.length === 1) {
      window.location.assign(`/api/v1/documents/${encodeURIComponent(id)}/download?format=${encodeURIComponent(options[0].value)}`);
      return;
    }
    openModal({
      title: `下载 ${document.type}`,
      eyebrow: `${document.id} / 选择输出格式`,
      body: `<div class="form-section"><h3>选择文件格式</h3><p>系统会先按当前订单、客户和产品快照填充字段，再按所选格式导出。PI 只支持 PDF；CI、PL、合同支持 PDF、Word、Excel；报关单支持 PDF 和 XML。</p></div><div class="filter-row document-download-formats">${options.map(option => `<button class="button ${option.value === 'pdf' ? 'primary' : ''}" data-action="download-document-format" data-id="${escapeAttr(id)}" data-format="${escapeAttr(option.value)}">${icon(option.value === 'pdf' ? 'file-text' : option.value === 'xml' ? 'file-code-2' : option.value === 'excel' ? 'file-spreadsheet' : 'file-type-2')}${escapeHTML(option.label)}</button>`).join('')}</div>`,
      footer: `<button class="button" data-action="close-modal">取消</button>`
    });
    return;
  }
  window.location.assign(`/api/v1/documents/${encodeURIComponent(id)}/download?format=${encodeURIComponent(format)}`);
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

async function loadTemplates(kind=state.templateKind, force=false) {
  if (state.templatesLoading && !force) return;
  state.templatesLoading = true;
  try {
    const data = await apiFetch(`/api/v1/templates?kind=${encodeURIComponent(kind)}`);
    state.templates = Array.isArray(data.items) ? data.items : [];
  } catch (error) {
    toast('模板读取失败', error.message, 'warning');
    state.templates = [];
  } finally {
    state.templatesLoading = false;
  }
}

async function templateCenter(kind) {
  const name={quote:'报价单',order:'订单',document:'单据'}[kind]||'业务';
  state.templateKind=kind;
  state.templatePage = 1;
  if (kind === 'order') {
    orderTemplateExplanationCenter();
    return;
  }
  await loadTemplates(kind,true);
  if (kind === 'quote') {
    quoteTemplateConfigCenter();
    return;
  }
  if (kind === 'document') {
    documentTemplateConfigCenter();
    return;
  }
  const rows=sortRows(state.templates.filter(item=>item.kind===kind), state.templateSort, { updated: item => timestampValue(item.updated), name: item => String(item.name || '') });
  const pageSizeOptions = [10,20,50,100];
  const pageSize = pageSizeOptions.includes(Number(state.templatePageSize)) ? Number(state.templatePageSize) : 10;
  const totalPages = Math.max(1, Math.ceil(rows.length / pageSize));
  if (state.templatePage > totalPages) state.templatePage = totalPages;
  const start = (state.templatePage - 1) * pageSize;
  const pageRows = rows.slice(start, start + pageSize);
  const selected=state.templateUploads.filter(item=>item.kind===kind).slice(-1)[0];
  const uploadControls = kind === 'document' ? '' : `<input id="templateImageInput" type="file" accept="image/jpeg,image/png,image/webp" hidden><input id="templateFileInput" type="file" accept=".docx,.xlsx,.pdf,.html,.htm,.txt,.md" hidden><button class="button" data-action="upload-template-image">${icon('image-up')}上传图片生成模板</button><button class="button primary" data-action="upload-template-file">${icon('file-up')}上传模板</button>`;
  const uploadStatus = kind === 'document' ? '' : selected ? `<section class="template-upload-status"><span class="setting-icon">${icon(selected.mode==='image'?'scan-line':'file-check-2')}</span><div><strong>${escapeHTML(selected.name)}</strong><span>${escapeHTML(selected.size)} · ${selected.mode==='image'?'图片模板':'文件模板'} · ${escapeHTML(selected.status||'处理中')}</span><small>${escapeHTML(selected.message||'正在同步到后端模板库。')}</small></div>${badge(selected.status==='已上传'?'Active':'Review')}</section>` : '';
  openModal({title:`${name}模板管理`,eyebrow:'模板中心 / 实时接口',wide:true,body:`${kind === 'document' ? '' : '<input id="templateImageInput" type="file" accept="image/jpeg,image/png,image/webp" hidden><input id="templateFileInput" type="file" accept=".docx,.xlsx,.pdf,.html,.htm,.txt,.md" hidden>'}<div class="toolbar" style="margin-bottom:14px"><span class="badge ${rows.length?'green':'neutral'}">${rows.length ? `${rows.length} 个模板` : '暂无模板'}</span><span class="spacer"></span><select class="select" id="templatePageSize">${pageSizeOptions.map(size=>`<option value="${size}" ${size===pageSize?'selected':''}>每页 ${size} 条</option>`).join('')}</select><button class="button" data-action="refresh-templates">${icon('refresh-cw')}刷新</button>${uploadControls}</div>${uploadStatus}<div class="data-wrap"><table class="data-table" style="min-width:820px"><thead><tr><th>模板名称</th><th>来源</th><th>版本</th><th>文件</th><th>状态</th><th>${sortHeader('更新时间','template','updated',state.templateSort)}</th><th>操作</th></tr></thead><tbody>${pageRows.map(t=>templateRow(t)).join('') || `<tr><td colspan="7"><div class="empty-state">${icon('layout-template')}<div><h3>暂无${escapeHTML(name)}模板</h3><p>只保留默认值和字段映射，生成时按单据类型导出。</p></div></div></td></tr>`}</tbody></table></div><div class="pagination"><span>共 ${rows.length} 条记录 · 第 ${state.templatePage} / ${totalPages} 页</span><div class="inline-actions"><button class="button small ghost" data-action="template-page-prev" ${state.templatePage<=1?'disabled':''}>${icon('chevron-left')}</button><button class="button small" data-action="template-page-current">${state.templatePage}</button><button class="button small ghost" data-action="template-page-next" ${state.templatePage>=totalPages?'disabled':''}>${icon('chevron-right')}</button></div></div><div class="model-warning" style="margin-top:14px"><span>${icon('info')} 只维护默认值和字段映射，出单时按订单生成。</span></div><ol class="template-flow"><li>识别模板 → 读取默认值</li><li>按订单、客户、产品实时填充</li><li>按所选格式下载</li></ol>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});
  if (kind !== 'document') {
    document.getElementById('templateImageInput')?.addEventListener('change',event=>handleTemplateFile(event.target.files?.[0],'image',kind));
    document.getElementById('templateFileInput')?.addEventListener('change',event=>handleTemplateFile(event.target.files?.[0],'file',kind));
  }
  document.getElementById('templatePageSize')?.addEventListener('change',e=>{state.templatePageSize=Number(e.target.value)||10;state.templatePage=1;void templateCenter(kind);});
}

function orderTemplateExplanationCenter() {
  openModal({
    title:'订单说明',
    eyebrow:'订单 / 单据关系',
    wide:true,
    body:`<div class="standard-template-layout">
      <section class="standard-template-summary">
        <span class="setting-icon">${icon('package-check')}</span>
        <div>
          <h3>订单不单独维护模板</h3>
          <p>订单只保存业务数据；需要出具 PI、CI、PL、合同、报关单时，到单据模块按订单生成。</p>
          <div class="mini-meta"><span>订单：业务数据源</span><span>单据：按订单生成文件</span><span>PI：属于单据类型</span></div>
        </div>
        ${badge('已确认')}
      </section>
      <div class="data-wrap"><table class="data-table" style="min-width:760px"><thead><tr><th>业务对象</th><th>职责</th><th>模板状态</th></tr></thead><tbody>
        <tr><td><strong>报价单</strong></td><td>维护报价明细，下载标准报价 PDF</td><td>${badge('已内置')}</td></tr>
        <tr><td><strong>订单</strong></td><td>维护订单数据和流程，不直接出模板文件</td><td>${badge('不需要')}</td></tr>
        <tr><td><strong>单据</strong></td><td>关联订单生成 PI、CI、PL、合同、报关单</td><td>${badge('已内置')}</td></tr>
      </tbody></table></div>
      <div class="model-warning"><span>${icon('info')} 订单只保留业务数据，单据从订单生成文件。</span></div>
    </div>`,
    footer:`<button class="button" data-action="close-modal">关闭</button><button class="button primary" data-page="documents">${icon('file-check-2')}去生成单据</button>`
  });
}

function documentTemplateConfigCenter() {
  const rows=state.templates.filter(item=>item.kind==='document'&&!String(item.status||'').toLowerCase().includes('archived'));
  const pi=rows.find(item=>item.id==='DOCUMENT-TPL-PI-20260625') || rows.find(item=>String(item.fileName||'').toLowerCase().endsWith('.pdf')) || {};
  const customs=rows.find(item=>item.id==='DOCUMENT-TPL-CUSTOMS-20260826') || {};
  const contract=rows.find(item=>item.id==='DOCUMENT-TPL-CONTRACT-20260826') || {};
  const invoice=rows.find(item=>item.id==='DOCUMENT-TPL-INVOICE-20260826') || {};
  const packing=rows.find(item=>item.id==='DOCUMENT-TPL-PACKING-20260826') || {};
  const tpl=pi.id ? pi : (invoice.id || packing.id || customs.id || contract.id ? (invoice.id ? invoice : packing.id ? packing : customs.id ? customs : contract) : {});
  const defaults=documentTemplateDefaultValuesJS(tpl.defaultValues);
  state.formContext={type:'template',id:tpl.id||'',documentTemplateIds:{pi:pi.id||'',customs:customs.id||'',contract:contract.id||'',invoice:invoice.id||'',packing:packing.id||''}};
  const body=`<div class="standard-template-layout">
    <section class="standard-template-summary">
      <span class="setting-icon">${icon('file-check-2')}</span>
      <div>
        <h3>单据标准模板</h3>
        <p>维护默认值；生成单据时按订单自动填充并下载。</p>
        <div class="mini-meta"><span>PI：${escapeHTML(pi.fileName||'未识别')}</span><span>CI：${escapeHTML(invoice.fileName||'未识别')}</span><span>PL：${escapeHTML(packing.fileName||'未识别')}</span><span>合同：${escapeHTML(contract.fileName||'未识别')}</span><span>报关单：${escapeHTML(customs.fileName||'未识别')}</span></div>
      </div>
      ${badge('默认')}
    </section>
    <div class="data-wrap"><table class="data-table" style="min-width:980px"><thead><tr><th>模板</th><th>用途</th><th>支持格式</th><th>状态</th><th>操作</th></tr></thead><tbody>
      ${pi.id?`<tr><td><strong>${escapeHTML(pi.name)}</strong></td><td>生成 PI 单</td><td>PDF</td><td>${badge(pi.status||'Active')}</td><td><button class="button small" data-action="template-preview" data-id="${escapeAttr(pi.id)}">${icon('eye')}预览</button></td></tr>`:''}
      ${invoice.id?`<tr><td><strong>${escapeHTML(invoice.name)}</strong></td><td>生成 CI / 发票</td><td>PDF / Word / Excel</td><td>${badge(invoice.status||'Active')}</td><td><button class="button small" data-action="template-preview" data-id="${escapeAttr(invoice.id)}">${icon('eye')}预览</button></td></tr>`:''}
      ${packing.id?`<tr><td><strong>${escapeHTML(packing.name)}</strong></td><td>生成 PL / 装箱单</td><td>PDF / Word / Excel</td><td>${badge(packing.status||'Active')}</td><td><button class="button small" data-action="template-preview" data-id="${escapeAttr(packing.id)}">${icon('eye')}预览</button></td></tr>`:''}
      ${contract.id?`<tr><td><strong>${escapeHTML(contract.name)}</strong></td><td>生成合同</td><td>PDF / Word / Excel</td><td>${badge(contract.status||'Active')}</td><td><button class="button small" data-action="template-preview" data-id="${escapeAttr(contract.id)}">${icon('eye')}预览</button></td></tr>`:''}
      ${customs.id?`<tr><td><strong>${escapeHTML(customs.name)}</strong></td><td>生成报关单</td><td>PDF / XML</td><td>${badge(customs.status||'Active')}</td><td><button class="button small" data-action="template-preview" data-id="${escapeAttr(customs.id)}">${icon('eye')}预览</button></td></tr>`:''}
    </tbody></table></div>
    <div class="form-grid">
      <div class="form-section full"><h3>供方默认信息</h3><p>客户、订单、产品明细由业务数据实时填充；供方和贸易默认值在这里统一维护。</p></div>
      ${inputField('供方公司',defaults['supplier.company'],false,false,'text','documentDefaultSupplierCompany')}
      ${inputField('供方地址',defaults['supplier.address'],false,false,'text','documentDefaultSupplierAddress')}
      ${inputField('供方国家',defaults['supplier.country'],false,false,'text','documentDefaultSupplierCountry')}
      ${inputField('供方电话',defaults['supplier.phone'],false,false,'text','documentDefaultSupplierPhone')}
      ${inputField('供方邮箱',defaults['supplier.email'],false,false,'email','documentDefaultSupplierEmail')}
      ${inputField('供方网址',defaults['supplier.website'],false,false,'text','documentDefaultSupplierWebsite')}
      <div class="form-section full"><h3>报关与贸易默认值</h3><p>用于合同、发票、装箱单和报关单的字段填充。</p></div>
      ${inputField('境内发货人',defaults['domestic.shipper'],false,false,'text','documentDefaultDomesticShipper')}
      ${inputField('境内发货人地址',defaults['domestic.shipper.address'],false,false,'text','documentDefaultDomesticShipperAddress')}
      ${inputField('海关编码',defaults['customs.code'],false,false,'text','documentDefaultCustomsCode')}
      ${inputField('出口海关',defaults['exit.customs'],false,false,'text','documentDefaultExitCustoms')}
      ${inputField('贸易方式',defaults['trade.mode'],false,false,'text','documentDefaultTradeMode')}
      ${inputField('贸易条款',defaults['trade.term'],false,false,'text','documentDefaultTradeTerm')}
      ${inputField('运输方式',defaults['transport.mode'],false,false,'text','documentDefaultTransportMode')}
      ${inputField('结算方式',defaults['settlement.method'],false,false,'text','documentDefaultSettlementMethod')}
      ${inputField('起运港',defaults['port.departure'],false,false,'text','documentDefaultPortDeparture')}
      ${inputField('成交地',defaults['contract.sign.place'],false,false,'text','documentDefaultSignPlace')}
      ${inputField('原产国',defaults['origin.country'],false,false,'text','documentDefaultOriginCountry')}
      ${inputField('包装类型',defaults['packaging.type'],false,false,'text','documentDefaultPackagingType')}
      ${inputField('默认箱数',defaults['package.count'],false,false,'number','documentDefaultPackageCount')}
      ${inputField('默认单位',defaults['quantity.unit'],false,false,'text','documentDefaultQuantityUnit')}
      ${inputField('毛重',defaults['gross.weight'],false,false,'text','documentDefaultGrossWeight')}
      ${inputField('净重',defaults['net.weight'],false,false,'text','documentDefaultNetWeight')}
      <div class="form-section full"><h3>PI 与备注默认值</h3></div>
      ${inputField('付款条款',defaults['payment.terms'],false,false,'text','documentDefaultPaymentTerms')}
      <div class="form-field full"><label>单据备注</label><textarea class="textarea" id="documentDefaultNotes" rows="3">${escapeHTML(defaults['document.notes'])}</textarea></div>
      <div class="form-field full"><label>银行账户</label><textarea class="textarea" id="documentDefaultBankAccount" rows="3">${escapeHTML(defaults['bank.account'])}</textarea></div>
      <div class="form-field full"><label>申报要素默认提示</label><textarea class="textarea" id="documentDefaultDeclarationElement" rows="3">${escapeHTML(defaults['declaration.element'])}</textarea></div>
    </div>
    <div class="model-warning" style="margin-top:14px"><span>${icon('info')} 生成时自动填充订单、客户和产品信息。</span></div>
  </div>`;
  openModal({title:'标准单据模板配置',eyebrow:'系统内置模板 / 默认值',wide:true,body,footer:`<button class="button" data-action="template-preview" data-id="${escapeAttr(pi.id||invoice.id||packing.id||contract.id||customs.id||'')}">${icon('eye')}预览模板</button><button class="button primary" data-action="save-template">${icon('save')}保存默认值</button>`});
}

function quoteTemplateConfigCenter() {
  const rows=state.templates.filter(item=>item.kind==='quote'&&!String(item.status||'').toLowerCase().includes('archived'));
  const tpl=rows.find(item=>item.default) || rows[0] || {};
  const defaults=quoteTemplateDefaultValuesJS(tpl.defaultValues);
  const body=`<div class="standard-template-layout">
    <section class="standard-template-summary">
      <span class="setting-icon">${icon('file-check-2')}</span>
      <div>
        <h3>${escapeHTML(tpl.name||'STRATRONIX 标准报价单 PDF 模板')}</h3>
        <p>报价单固定使用系统内置 PDF 标准模板，不再支持上传图片、上传模板或任意 PDF 自动识别。这里只维护生成时会复用的默认填充值。</p>
        <div class="mini-meta"><span>版本：${escapeHTML(tpl.version||'2026-08-19')}</span><span>状态：${escapeHTML(tpl.status||'Active')}</span><span>更新时间：${escapeHTML(tpl.updated||'')}</span></div>
      </div>
      ${badge('默认')}
    </section>
    <div class="form-grid">
      <div class="form-section full"><h3>供方默认信息</h3><p>用于填充 PDF 右侧供方区域；客户方和报价明细仍从客户、报价单、产品实时读取。</p></div>
      ${inputField('供方公司',defaults['supplier.company'],false,false,'text','templateDefaultSupplierCompany')}
      ${inputField('供方联系人',defaults['supplier.contact'],false,false,'text','templateDefaultSupplierContact')}
      ${inputField('供方电话',defaults['supplier.phone'],false,false,'text','templateDefaultSupplierPhone')}
      ${inputField('供方邮箱',defaults['supplier.email'],false,false,'email','templateDefaultSupplierEmail')}
      ${inputField('供方网址',defaults['supplier.website'],false,false,'text','templateDefaultSupplierWebsite')}
      <div></div>
      <div class="form-section full"><h3>报价条款默认值</h3><p>新建报价单会自动带入；生成 PDF 时报价单字段优先，没填写才使用这里的默认值。</p></div>
      ${inputField('默认价格条件',defaults['quote.priceTerms'],false,false,'text','templateDefaultPriceTerms')}
      ${inputField('默认付款条件',defaults['quote.paymentTerms'],false,false,'text','templateDefaultPaymentTerms')}
      ${inputField('默认目的港/目的地',defaults['quote.destination'],false,false,'text','templateDefaultDestination')}
      ${inputField('默认交期',defaults['quote.leadTime'],false,false,'text','templateDefaultLeadTime')}
      ${inputField('默认质保',defaults['quote.warranty'],false,false,'text','templateDefaultWarranty')}
      ${inputField('默认认证',defaults['quote.certification'],false,false,'text','templateDefaultCertification')}
      ${inputField('默认包装',defaults['quote.packaging'],false,false,'text','templateDefaultPackaging')}
      ${inputField('默认有效期天数',defaults['quote.validityDays'],false,false,'number','templateDefaultValidityDays')}
      <div class="form-field full"><label>默认备注</label><textarea class="textarea" id="templateDefaultNote" rows="3">${escapeHTML(defaults['quote.note'])}</textarea></div>
    </div>
    <div class="model-warning" style="margin-top:14px"><span>${icon('info')} 版式、字段坐标和模板文件由系统维护；用户只需要维护这些通用默认值。实时字段包括报价单号、日期、客户、产品、数量、单价、运费和总计。</span></div>
  </div>`;
  state.formContext={type:'template',id:tpl.id||''};
  openModal({title:'标准报价单配置',eyebrow:'系统内置模板 / 默认值',wide:true,body,footer:`<button class="button" data-action="template-preview" data-id="${escapeAttr(tpl.id||'')}">${icon('eye')}预览模板</button><button class="button primary" data-action="save-template">${icon('save')}保存默认值</button>`});
}

function quoteTemplateDefaultValuesJS(values={}) {
  const defaults={
    'supplier.company':'STRATRONIX / 鼎图',
    'supplier.contact':'Donald',
    'supplier.phone':'86-755-23086689',
    'supplier.email':'info@stratronix.ai',
    'supplier.website':'www.stratronix.ai',
    'quote.priceTerms':'FOB Shenzhen',
    'quote.paymentTerms':'T/T 30% deposit, 70% before shipment',
    'quote.shipping':'Shipping to be confirmed',
    'quote.destination':'Rotterdam',
    'quote.leadTime':'15-20 working days after deposit received',
    'quote.warranty':'12 months under STRATRONIX standard terms',
    'quote.certification':'CE / FCC / RoHS if applicable',
    'quote.packaging':'Standard export carton',
    'quote.note':'Quotation is subject to final quantity and specification confirmation.',
    'quote.validityDays':'30',
  };
  return {...defaults,...(values||{})};
}

function documentTemplateDefaultValuesJS(values={}) {
  const defaults={
    'supplier.company':'Stratronix Technology (Shenzhen) Company, Limited',
    'supplier.address':'Shenzhen China 528201',
    'supplier.country':'China',
    'supplier.phone':'86-755-23086689',
    'supplier.email':'info@stratronix.ai',
    'supplier.website':'www.stratronix.ai',
    'domestic.shipper':'Stratronix Technology (Shenzhen) Company, Limited',
    'domestic.shipper.address':'Shenzhen China 528201',
    'customs.code':'',
    'exit.customs':'皇岗海关',
    'trade.mode':'一般贸易',
    'trade.term':'FOB',
    'transport.mode':'航空',
    'settlement.method':'T/T',
    'port.departure':'Shenzhen',
    'contract.sign.place':'Shenzhen',
    'origin.country':'中国',
    'packaging.type':'纸箱',
    'package.count':'1',
    'quantity.unit':'台',
    'gross.weight':'',
    'net.weight':'',
    'payment.terms':'Due on Receipt',
    'document.notes':'1. STA-100, HDMI CABLE, MANUAL, GIFTBOX in total one box;\n2. Product details and customs declaration elements require final manual review.',
    'bank.account':'BANK ACCOUNT:',
    'declaration.element':'品牌类型、出口享惠情况、用途、材质、规格型号等申报要素待人工复核。',
  };
  return {...defaults,...(values||{})};
}

function templateRow(t={}) {
  const actions=[
    `<button class="button small" data-action="template-preview" data-id="${escapeAttr(t.id)}">${icon('eye')}预览</button>`,
    `<button class="button small" data-action="template-edit" data-id="${escapeAttr(t.id)}">${icon('pencil')}编辑</button>`,
    t.default ? `<span class="badge green">默认</span>` : `<button class="button small" data-action="template-default" data-id="${escapeAttr(t.id)}">${icon('badge-check')}设为默认</button>`,
    t.default ? '' : `<button class="button danger small" data-action="template-delete" data-id="${escapeAttr(t.id)}">${icon('trash-2')}删除</button>`,
  ].filter(Boolean).join('');
  return `<tr><td><strong>${escapeHTML(t.name)}</strong>${t.default?`<small class="secondary-text"> 当前默认模板</small>`:''}</td><td>${escapeHTML(t.source||'上传')}</td><td>${escapeHTML(t.version||'v1')}</td><td>${escapeHTML(t.fileName||'系统模板')}<br><small class="secondary-text">${escapeHTML(t.size||'')}</small></td><td>${badge(t.status||'Draft')}</td><td>${escapeHTML(t.updated||'')}</td><td><span class="table-actions">${actions}</span></td></tr>`;
}

async function handleTemplateFile(file,mode,kind) {
  if(!file)return;
  const allowed=mode==='image'?/\.(jpe?g|png|webp)$/i:/\.(docx|xlsx|pdf|html?|txt|md)$/i;
  if(!allowed.test(file.name)){toast('文件格式不支持',mode==='image'?'请选择 JPG、PNG 或 WebP 图片。':'请选择 DOCX、XLSX、PDF 或 HTML 模板。','warning');return;}
  if(file.size>20*1024*1024){toast('文件过大','模板文件不能超过 20 MB。','warning');return;}
  const size=file.size>=1024*1024?`${(file.size/1024/1024).toFixed(1)} MB`:`${Math.max(1,Math.round(file.size/1024))} KB`;
  state.templateUploads.push({kind,mode,name:file.name,size,status:'上传中',message:'正在上传并写入模板库。'});
  void templateCenter(kind);
  const form=new FormData();form.append('file',file);form.append('kind',kind);
  try {
    const record=await apiFetch(mode==='image'?'/api/v1/templates/image-recognition':'/api/v1/templates/upload',{method:'POST',body:form});
    state.templateUploads.push({kind,mode,name:file.name,size,status:'已上传',message:`${record.name} 已保存，可编辑字段映射或发布为默认。`});
    toast('模板已上传',`${record.name} 已写入模板库。`);
  }
  catch(error){state.templateUploads.push({kind,mode,name:file.name,size,status:'上传失败',message:error.message});toast('模板上传失败',error.message,'warning');}
  await templateCenter(kind);
}

function editTemplateForm(id) {
  const t=state.templates.find(item=>item.id===id);
  if(!t){toast('模板不存在','请刷新模板列表后重试。','warning');return;}
  if (t.kind === 'quote') {
    state.formContext={type:'template',id:t.id};
    quoteTemplateConfigCenter();
    return;
  }
  state.formContext={type:'template',id:t.id};
  const pdfFields = JSON.stringify(t.pdfFields||defaultPdfFieldsForTemplate(t.kind),null,2);
  const pdfFormFields = JSON.stringify(t.pdfFormFields||[],null,2);
  openModal({title:'编辑模板',eyebrow:`${templateKindLabelJS(t.kind)} / ${t.id}`,body:`<div class="form-grid">${inputField('模板名称',t.name||'',true,false,'text','templateName')}${inputField('版本号',t.version||'v1',true,false,'text','templateVersion')}${selectField('状态',[{value:'Draft',label:'草稿 Draft'},{value:'Review',label:'待复核 Review'},{value:'Active',label:'启用 Active'}],false,'templateStatus',t.status||'Draft')}${selectField('输出格式',[{value:'html',label:'HTML 文件'},{value:'pdf',label:'PDF 文件'}],false,'templateOutputFormat',t.outputFormat||'html')}<div class="form-section full"><h3>已识别字段</h3><p>${(t.placeholders||[]).map(v=>`<code>${escapeHTML(v)}</code>`).join(' ') || '未识别到占位符，将使用系统默认字段。'}</p></div><div class="form-section full"><h3>字段映射</h3><textarea id="templateFieldMapping" rows="8">${escapeHTML(JSON.stringify(t.fieldMapping||{},null,2))}</textarea><small class="secondary-text">JSON 左侧写模板字段名，右侧写业务字段路径，例如 <code>"customer.name"</code>。</small></div><div class="form-section full"><h3>PDF 字段坐标</h3><textarea id="templatePdfFields" rows="10">${escapeHTML(pdfFields)}</textarea><small class="secondary-text">x/y 以 PDF 左下角为原点，单位 points。Page 默认 1。若模板本身带表单字段，也可在下方填写 PDF 表单字段名列表。</small></div><div class="form-section full"><h3>PDF 表单字段名</h3><textarea id="templatePdfFormFields" rows="4">${escapeHTML(pdfFormFields)}</textarea><small class="secondary-text">仅当模板是 AcroForm PDF 时生效。留空则使用坐标叠加。</small></div><div class="model-warning full"><span>${icon('info')} 字段映射按 JSON 保存。业务生成时会把对应字段写入模板占位符，PDF 模板会优先按表单字段填充，再回退到坐标叠加。</span></div></div>`,footer:formFooter('保存模板','save-template')});
}

async function saveTemplate() {
  const id=state.formContext?.id;
  const current=state.templates.find(item=>item.id===id);
  if(!current)return;
  if (current.kind === 'quote' || current.kind === 'document') {
    const defaultValues=current.kind==='quote'?{
        'supplier.company':formText('templateDefaultSupplierCompany'),
        'supplier.contact':formText('templateDefaultSupplierContact'),
        'supplier.phone':formText('templateDefaultSupplierPhone'),
        'supplier.email':formText('templateDefaultSupplierEmail'),
        'supplier.website':formText('templateDefaultSupplierWebsite'),
        'quote.priceTerms':formText('templateDefaultPriceTerms'),
        'quote.paymentTerms':formText('templateDefaultPaymentTerms'),
        'quote.destination':formText('templateDefaultDestination'),
        'quote.leadTime':formText('templateDefaultLeadTime'),
        'quote.warranty':formText('templateDefaultWarranty'),
        'quote.certification':formText('templateDefaultCertification'),
        'quote.packaging':formText('templateDefaultPackaging'),
        'quote.note':formText('templateDefaultNote'),
        'quote.validityDays':formText('templateDefaultValidityDays') || '30',
      }:{
        'supplier.company':formText('documentDefaultSupplierCompany'),
        'supplier.address':formText('documentDefaultSupplierAddress'),
        'supplier.country':formText('documentDefaultSupplierCountry'),
        'supplier.phone':formText('documentDefaultSupplierPhone'),
        'supplier.email':formText('documentDefaultSupplierEmail'),
        'supplier.website':formText('documentDefaultSupplierWebsite'),
        'domestic.shipper':formText('documentDefaultDomesticShipper'),
        'domestic.shipper.address':formText('documentDefaultDomesticShipperAddress'),
        'customs.code':formText('documentDefaultCustomsCode'),
        'exit.customs':formText('documentDefaultExitCustoms'),
        'trade.mode':formText('documentDefaultTradeMode'),
        'trade.term':formText('documentDefaultTradeTerm'),
        'transport.mode':formText('documentDefaultTransportMode'),
        'settlement.method':formText('documentDefaultSettlementMethod'),
        'port.departure':formText('documentDefaultPortDeparture'),
        'contract.sign.place':formText('documentDefaultSignPlace'),
        'origin.country':formText('documentDefaultOriginCountry'),
        'packaging.type':formText('documentDefaultPackagingType'),
        'package.count':formText('documentDefaultPackageCount') || '1',
        'quantity.unit':formText('documentDefaultQuantityUnit') || '台',
        'gross.weight':formText('documentDefaultGrossWeight'),
        'net.weight':formText('documentDefaultNetWeight'),
        'payment.terms':formText('documentDefaultPaymentTerms'),
        'document.notes':formText('documentDefaultNotes'),
        'bank.account':formText('documentDefaultBankAccount'),
        'declaration.element':formText('documentDefaultDeclarationElement'),
      };
    const targets=current.kind==='document'
      ? state.templates.filter(item=>item.kind==='document'&&!String(item.status||'').toLowerCase().includes('archived'))
      : [current];
    try {
      for (const target of targets) {
        await apiFetch(`/api/v1/templates/${encodeURIComponent(target.id)}`,{method:'PATCH',body:JSON.stringify({...target,defaultValues})});
      }
      await loadTemplates(current.kind,true);
      closeModal();
      toast('默认值已保存',current.kind==='quote'?'后续生成报价单会使用新的标准模板默认值。':'后续生成单据会使用新的标准模板默认值。');
    } catch(error) {toast('默认值保存失败',error.message,'warning');}
    return;
  }
  let mapping=current.fieldMapping||{};
  let pdfFields=current.pdfFields||[];
  let pdfFormFields=current.pdfFormFields||[];
  try {mapping=JSON.parse(document.getElementById('templateFieldMapping')?.value||'{}');}
  catch {toast('保存失败','字段映射必须是合法 JSON。','warning');return;}
  try {pdfFields=JSON.parse(document.getElementById('templatePdfFields')?.value||'[]');}
  catch {toast('保存失败','PDF 字段坐标必须是合法 JSON。','warning');return;}
  try {pdfFormFields=JSON.parse(document.getElementById('templatePdfFormFields')?.value||'[]');}
  catch {toast('保存失败','PDF 表单字段名必须是合法 JSON。','warning');return;}
  const payload={...current,name:formText('templateName'),version:formText('templateVersion'),status:formText('templateStatus'),outputFormat:formText('templateOutputFormat'),fieldMapping:mapping,pdfFields:pdfFields,pdfFormFields:pdfFormFields};
  try {
    await apiFetch(`/api/v1/templates/${encodeURIComponent(id)}`,{method:'PATCH',body:JSON.stringify(payload)});
    closeModal();
    await templateCenter(current.kind);
    toast('模板已保存',`${payload.name} 已更新。`);
  } catch(error) {toast('模板保存失败',error.message,'warning');}
}

async function templateAction(action,id='') {
  const templateID=id||state.templates[0]?.id;
  if(!templateID){toast('请选择模板','当前没有可操作模板。','warning');return;}
  try {
    if(action==='default'){
      const record=await apiFetch(`/api/v1/templates/${encodeURIComponent(templateID)}/publish`,{method:'POST',body:'{}'});
      toast('默认模板已更新',`${record.name} 已设为默认。`);
      await templateCenter(record.kind);
    } else if(action==='preview') {
      window.location.assign(`/api/v1/templates/${encodeURIComponent(templateID)}/render`);
    } else if(action==='edit') {
      editTemplateForm(templateID);
    } else if(action==='delete') {
      await apiFetch(`/api/v1/templates/${encodeURIComponent(templateID)}`,{method:'DELETE'});
      toast('模板已删除','模板已从列表归档。');
      await templateCenter(state.templateKind);
    }
  }
  catch(error){toast('模板操作失败',error.message,'warning');}
}

function templateKindLabelJS(kind) {
  return {quote:'报价单',order:'订单',document:'单据'}[kind]||'业务';
}

function defaultPdfFieldsForTemplate(kind) {
  const common = [
    { key: 'company.name', label: '公司名称', page: 1, x: 42, y: 800, fontSize: 12, align: 'l' },
    { key: 'customer.name', label: '客户名称', page: 1, x: 42, y: 772, fontSize: 11, align: 'l' },
    { key: 'record.id', label: '单号', page: 1, x: 405, y: 800, fontSize: 11, align: 'l' },
    { key: 'record.date', label: '日期', page: 1, x: 405, y: 778, fontSize: 11, align: 'l' },
    { key: 'record.total', label: '总金额', page: 1, x: 405, y: 754, fontSize: 12, align: 'l' },
  ];
  if (kind === 'quote') {
    return [
      { key: 'pdf.record.id', label: '报价单号', page: 1, x: 171, y: 607, fontSize: 8, maxWidth: 100, align: 'l' },
      { key: 'record.date', label: '报价日期', page: 1, x: 388, y: 610, fontSize: 9, maxWidth: 80, align: 'l' },
      { key: 'quote.validity', label: '有效期', page: 1, x: 171, y: 574, fontSize: 8, maxWidth: 100, align: 'l' },
      { key: 'pdf.quote.paymentTerms', label: '付款条件', page: 1, x: 388, y: 574, fontSize: 8, maxWidth: 150, align: 'l' },
      { key: 'customer.name', label: '客户公司', page: 1, x: 128, y: 510, fontSize: 8, maxWidth: 140, align: 'l' },
      { key: 'customer.country', label: '客户国家', page: 1, x: 128, y: 495, fontSize: 8, maxWidth: 140, align: 'l' },
      { key: 'pdf.customer.contact', label: '客户联系人', page: 1, x: 128, y: 480, fontSize: 8, maxWidth: 140, align: 'l' },
      { key: 'customer.email', label: '客户邮箱', page: 1, x: 128, y: 465, fontSize: 8, maxWidth: 140, align: 'l' },
      { key: 'pdf.customer.address', label: '客户地址', page: 1, x: 128, y: 450, fontSize: 8, maxWidth: 145, align: 'l' },
      { key: 'line.1.no', label: '序号', page: 1, x: 67, y: 368, fontSize: 7, maxWidth: 20, align: 'c' },
      { key: 'line.1.model', label: '型号', page: 1, x: 91, y: 368, fontSize: 7, maxWidth: 52, align: 'l' },
      { key: 'line.1.description', label: '描述', page: 1, x: 159, y: 368, fontSize: 7, maxWidth: 170, align: 'l' },
      { key: 'line.1.quantity', label: '数量', page: 1, x: 349, y: 368, fontSize: 7, maxWidth: 36, align: 'l' },
      { key: 'pdf.line.1.unitPrice', label: '单价', page: 1, x: 430, y: 368, fontSize: 7, maxWidth: 52, align: 'r' },
      { key: 'pdf.line.1.amount', label: '小计', page: 1, x: 514, y: 368, fontSize: 7, maxWidth: 52, align: 'r' },
      { key: 'line.2.no', label: '序号2', page: 1, x: 67, y: 356, fontSize: 7, maxWidth: 20, align: 'c' },
      { key: 'line.2.model', label: '型号2', page: 1, x: 91, y: 356, fontSize: 7, maxWidth: 52, align: 'l' },
      { key: 'line.2.description', label: '描述2', page: 1, x: 159, y: 356, fontSize: 7, maxWidth: 170, align: 'l' },
      { key: 'line.2.quantity', label: '数量2', page: 1, x: 349, y: 356, fontSize: 7, maxWidth: 36, align: 'l' },
      { key: 'pdf.line.2.unitPrice', label: '单价2', page: 1, x: 430, y: 356, fontSize: 7, maxWidth: 52, align: 'r' },
      { key: 'pdf.line.2.amount', label: '小计2', page: 1, x: 514, y: 356, fontSize: 7, maxWidth: 52, align: 'r' },
      { key: 'line.3.no', label: '序号3', page: 1, x: 67, y: 344, fontSize: 7, maxWidth: 20, align: 'c' },
      { key: 'line.3.model', label: '型号3', page: 1, x: 91, y: 344, fontSize: 7, maxWidth: 52, align: 'l' },
      { key: 'line.3.description', label: '描述3', page: 1, x: 159, y: 344, fontSize: 7, maxWidth: 170, align: 'l' },
      { key: 'line.3.quantity', label: '数量3', page: 1, x: 349, y: 344, fontSize: 7, maxWidth: 36, align: 'l' },
      { key: 'pdf.line.3.unitPrice', label: '单价3', page: 1, x: 430, y: 344, fontSize: 7, maxWidth: 52, align: 'r' },
      { key: 'pdf.line.3.amount', label: '小计3', page: 1, x: 514, y: 344, fontSize: 7, maxWidth: 52, align: 'r' },
      { key: 'pdf.record.subtotal', label: '产品小计', page: 1, x: 526, y: 322, fontSize: 8, maxWidth: 82, align: 'r' },
      { key: 'pdf.record.freight', label: '运费', page: 1, x: 526, y: 294, fontSize: 8, maxWidth: 82, align: 'r' },
      { key: 'pdf.record.total', label: '总计', page: 1, x: 526, y: 242, fontSize: 10, maxWidth: 82, align: 'r' },
      { key: 'quote.priceTerms', label: '价格条件', page: 1, x: 71, y: 200, fontSize: 8, maxWidth: 430, align: 'l' },
      { key: 'quote.paymentTerms', label: '付款条件条款', page: 1, x: 71, y: 187, fontSize: 8, maxWidth: 430, align: 'l' },
      { key: 'quote.leadTime', label: '交期', page: 1, x: 71, y: 173, fontSize: 8, maxWidth: 430, align: 'l' },
      { key: 'quote.warranty', label: '质保', page: 1, x: 71, y: 160, fontSize: 8, maxWidth: 430, align: 'l' },
      { key: 'quote.certification', label: '认证', page: 1, x: 71, y: 147, fontSize: 8, maxWidth: 430, align: 'l' },
      { key: 'quote.packaging', label: '包装', page: 1, x: 71, y: 133, fontSize: 8, maxWidth: 430, align: 'l' },
      { key: 'quote.validity', label: '报价有效期条款', page: 1, x: 71, y: 120, fontSize: 8, maxWidth: 430, align: 'l' },
      { key: 'quote.note', label: '备注', page: 1, x: 71, y: 107, fontSize: 8, maxWidth: 430, align: 'l' },
    ];
  }
  if (kind === 'order') {
    return common.concat([
      { key: 'order.po', label: '客户 PO 号', page: 1, x: 42, y: 744, fontSize: 11, align: 'l' },
      { key: 'order.delivery', label: '交付日期', page: 1, x: 42, y: 720, fontSize: 11, align: 'l' },
      { key: 'order.status', label: '订单状态', page: 1, x: 42, y: 696, fontSize: 11, align: 'l' },
      { key: 'record.lines', label: '订单明细', page: 1, x: 40, y: 620, fontSize: 10, align: 'l', maxWidth: 500, multi: true },
    ]);
  }
  return common.concat([
    { key: 'document.type', label: '单据类型', page: 1, x: 42, y: 744, fontSize: 11, align: 'l' },
    { key: 'document.language', label: '输出语言', page: 1, x: 42, y: 720, fontSize: 11, align: 'l' },
    { key: 'order.id', label: '关联订单', page: 1, x: 42, y: 696, fontSize: 11, align: 'l' },
    { key: 'record.lines', label: '单据明细', page: 1, x: 40, y: 620, fontSize: 10, align: 'l', maxWidth: 500, multi: true },
  ]);
}

function documentDetail(id) {
  const d=documents.find(x=>x.id===id);
  if(!d)return;
  const sourceOrder=orders.find(order=>order.id===d.order); const lines=d.lines?.length?d.lines:businessLines(sourceOrder||{});
  const previewRows=lines.map(line=>`<tr><td style="padding:9px 7px">${escapeHTML(line.productName||lineProduct(line)?.name||line.productId)}</td><td style="text-align:center">${Number(line.quantity)}</td><td style="text-align:center">${formatMoney(line.unitPrice)}</td><td style="text-align:right">${formatMoney(line.amount??lineSubtotal(line))}</td></tr>`).join('');
  openDrawer({title:d.id,eyebrow:`${d.type} / 单据预览`,body:`<div class="spread"><span>${badge(d.status)}</span><div class="inline-actions"><button class="button small" data-action="edit-document" data-id="${escapeAttr(d.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-document" data-id="${escapeAttr(d.id)}">${icon('trash-2')}删除</button><button class="button primary small" data-action="download-document" data-id="${escapeAttr(d.id)}">${icon('file-down')}下载文件</button></div></div><div class="panel document-preview-paper" style="margin-top:15px;background:#edf1f5;color:#1b2634;min-height:520px;padding:28px"><div style="border-bottom:2px solid #29394b;padding-bottom:14px"><strong style="font-size:18px">STA-100</strong><span style="float:right;font-size:18px">${escapeHTML(d.type)}</span></div><div style="display:grid;grid-template-columns:1fr 1fr;gap:20px;margin-top:20px;font-size:10px"><div><strong>Bill To</strong><p>${escapeHTML(d.customer)}<br>Customer address snapshot</p></div><div><strong>Document</strong><p>${escapeHTML(d.id)}<br>Order ${escapeHTML(d.order)}</p></div></div><table style="width:100%;border-collapse:collapse;margin-top:22px;font-size:9px"><tr style="border-bottom:1px solid #8793a1"><th style="text-align:left;padding:7px">Item</th><th>Qty</th><th>Rate</th><th>Amount</th></tr>${previewRows}</table><div style="margin-top:18px;text-align:right"><strong>Total ${escapeHTML(d.value||sourceOrder?.value||formatMoney(lines.reduce((sum,line)=>sum+lineSubtotal(line),0)))}</strong></div></div>`});
}

function newProductForm(product) {
  const p=product||{};
  state.formContext={type:'product',id:p.id||''};
  openModal({title:product?'编辑产品':'新建产品',eyebrow:'产品主数据',body:`<div class="form-grid">${inputField('产品名称',p.name||'',true,false,'text','productName')}${inputField('产品编码',p.id||'',true,false,'text','productID')}${selectField('产品类别',['智能设备','智能骑行','整车方案','配件','服务'].map(value=>({value,label:value})),false,'productCategory',p.category||'智能设备')}${selectField('状态',[{value:'Active',label:'已启用 Active'},{value:'Review',label:'待审核 Review'},{value:'Inactive',label:'已停用 Inactive'}],false,'productStatus',p.status||'Active')}${inputField('制造商',p.manufacturer||'STRATRONIX',false,false,'text','productManufacturer')}${inputField('HS CODE',p.hs||'',true,false,'text','productHS')}${inputField('库存量',p.stock||'0',false,false,'number','productStock')}${inputField('默认单价',p.price||'',true,false,'text','productPrice')}${inputField('产品描述',p.desc||'',false,true,'text','productDescription')}${inputField('标签','欧洲 / 智能设备',false,true,'text','productTags')}</div>`,footer:formFooter(product?'保存修改':'创建产品','save-product')});
}

async function saveProduct() {
  const name=formText('productName'); const id=formText('productID'); if(!name||!id){toast('保存失败','产品名称和编码为必填项。','warning');return;}
  const existing=products.find(p=>p.id===state.formContext?.id);
  const payload={...(existing||{}),id,name,category:formText('productCategory'),manufacturer:formText('productManufacturer'),hs:formText('productHS'),stock:formNumber('productStock'),price:formText('productPrice'),desc:formText('productDescription'),tags:formText('productTags'),status:formText('productStatus')||existing?.status||'Active'};
  try { const record=await apiFetch(existing?`/api/v1/products/${encodeURIComponent(existing.id)}`:'/api/v1/products',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)}); upsertRecord(products,record); closeModal(); renderPage(); toast(existing?'产品已更新':'产品已创建',`${record.name} 已保存到本地数据库。`); }
  catch(error) { toast('保存失败',error.message,'warning'); }
}
async function deleteProduct(id) { if(!window.confirm('确定停用该产品吗？'))return; try { await apiFetch(`/api/v1/products/${encodeURIComponent(id)}`,{method:'DELETE'}); const product=products.find(item=>item.id===id); if(product)product.status='Inactive'; closeDrawer(); renderPage(); toast('产品已停用','产品主数据已更新。'); } catch(error) { toast('停用失败',error.message,'warning'); } }

function productImportModal() {
  openModal({title:'批量导入产品',eyebrow:'产品库 / 文件导入',body:`<input id="productImportInput" type="file" accept=".xlsx,.csv" hidden><div class="upload-zone"><div><span class="upload-icon">${icon('file-up')}</span><h3>选择产品导入文件</h3><p>支持 XLSX 或 CSV。建议字段：产品编码、产品名称、类别、HS CODE、销售价、库存、状态、制造商、描述、标签。</p><button type="button" class="button primary" data-action="choose-product-import">${icon('folder-open')}选择文件</button></div></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});
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
  const quoteRefs = quotes.filter(q => businessLines(q).some(line => line.productId === p.id || line.productName === p.name));
  const orderRefs = orders.filter(o => businessLines(o).some(line => line.productId === p.id || line.productName === p.name));
  const docRefs = documents.filter(d => businessLines(d).some(line => line.productId === p.id || line.productName === p.name));
  const renderRefs = (title, empty, rows) => `<div class="divider-title" style="margin:18px 0 10px">${title}</div><div class="related-list">${rows.length ? rows.slice(0,5).map(row => row).join('') : `<div class="empty-state"><p>${empty}</p></div>`}</div>`;
  openDrawer({title:p.name,eyebrow:`产品 / ${p.id}`,body:`<div class="spread"><span class="badge ${p.status==='Active'?'green':p.status==='Inactive'?'neutral':'amber'}">${escapeHTML(productStatusLabel(p.status))}</span><div class="inline-actions"><button class="button small" data-action="edit-product" data-id="${escapeAttr(p.id)}">${icon('pencil')}编辑</button><button class="button danger small" data-action="delete-product" data-id="${escapeAttr(p.id)}">${icon('trash-2')}删除</button></div></div><div class="product-visual panel" style="margin-top:14px;aspect-ratio:16/6">${icon('cpu')}</div><div class="detail-grid" style="margin-top:14px">${[['产品编码',p.id],['产品类别',p.category],['HS CODE',p.hs],['销售价',p.price],['当前库存',String(p.stock)],['更新时间',p.updated||'未填写'],['状态',productStatusLabel(p.status)]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">产品描述</div><p class="secondary-text" style="line-height:1.7">${escapeHTML(p.desc)}</p>${renderRefs('业务引用 - 报价单','暂无关联报价单',quoteRefs.map(q=>`<button class="related-record" data-action="quote-detail" data-id="${escapeAttr(q.id)}"><span>${icon('file-text')}<strong>${escapeHTML(q.id)}</strong></span><span>${escapeHTML(q.subject)}</span><span>${escapeHTML(q.value)}</span>${badge(q.status)}</button>`))}${renderRefs('业务引用 - 订单','暂无关联订单',orderRefs.map(o=>`<button class="related-record" data-action="order-detail" data-id="${escapeAttr(o.id)}"><span>${icon('package')}<strong>${escapeHTML(o.id)}</strong></span><span>${escapeHTML(o.customer)}</span><span>${escapeHTML(o.value)}</span>${badge(o.status)}</button>`))}${renderRefs('业务引用 - 单据','暂无关联单据',docRefs.map(d=>`<button class="related-record" data-action="document-detail" data-id="${escapeAttr(d.id)}"><span>${icon('file-check-2')}<strong>${escapeHTML(d.id)}</strong></span><span>${escapeHTML(d.type)}</span><span>${escapeHTML(d.status)}</span>${badge(d.status)}</button>`))}`});
}

function uploadFileModal() {
  state.selectedUploadFile=null;
  openModal({title:'上传私有数据',eyebrow:'数据库 / 文件处理',body:`<input id="privateFileInput" type="file" accept=".pdf,.docx,.xlsx,.csv,.txt,.md,.rtf,.xml,.html,.htm" hidden><div class="upload-zone" id="uploadZone"><div><span class="upload-icon">${icon('cloud-upload')}</span><h3>选择或拖入文件</h3><p>支持 PDF、DOCX、XLSX、CSV、TXT、MD、RTF、XML、HTML；单文件最大 25 MB。</p><button type="button" class="button primary" data-action="choose-file">选择文件</button><div id="selectedPrivateFile"></div></div></div><div class="form-grid" style="margin-top:14px"><div class="form-field"><label>数据区</label><input class="input" value="客户私有数据（本机）" readonly></div>${selectField('主分类',['自动识别','合同','报价单','产品手册','法规','产品资料','会议记录','客户资料','其它'],false,'privateFileCategory','自动识别')}${inputField('附加标签','',false,true,'text','privateFileTags')}<div class="form-field full"><small>文件会在本机解析、分块并建立向量索引；PDF 暂保存为待处理，需部署 OCR 或先转为可检索文本。</small></div></div>`,footer:formFooter('上传到本机','upload-private-file')});
  const zone=document.getElementById('uploadZone');
  ['dragenter','dragover'].forEach(ev=>zone.addEventListener(ev,e=>{e.preventDefault();zone.classList.add('dragover')}));
  zone.addEventListener('dragleave',e=>{e.preventDefault();zone.classList.remove('dragover')});
  zone.addEventListener('drop',e=>{e.preventDefault();zone.classList.remove('dragover');selectPrivateFile(e.dataTransfer?.files?.[0])});
  document.getElementById('privateFileInput').addEventListener('change',e=>selectPrivateFile(e.target.files?.[0]));
}

function selectPrivateFile(file) {
  if(!file)return;
  const allowed=/\.(pdf|docx|xlsx|csv|txt|md|rtf|xml|html?|htm)$/i;
  if(!allowed.test(file.name)){toast('文件格式不支持','请选择页面列出的文件格式。','warning');return;}
  if(file.size>25*1024*1024){toast('文件过大','单文件不能超过 25 MB。','warning');return;}
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
  try {
    const result=await apiFetch(`/api/v1/private-files/${encodeURIComponent(id)}/summary`);
    const document=result.document||{};
    const detail=document.chunks?`${document.chunks} 个索引片段，来源更新时间 ${document.sourceUpdatedAt||'待核实'}`:(result.message||'暂无索引详情');
    toast(result.status==='indexed'?'知识库索引正常':'知识库索引待处理',detail,result.status==='indexed'?'success':'warning');
  }
  catch(error){toast('摘要暂不可用',error.message,'warning');}
}

function stripSTA100Result(value) {
  return String(value || '')
    .replace(/\[STA100_RESULT\][\s\S]*?(?:\[\/STA100_RESULT\]?|$)/g, '')
    .replace(/\{\s*"type"\s*:\s*"(?:news|recommendations)"\s*,\s*"items"\s*:\s*\[[\s\S]*?\]\s*\}/g, '')
    .trim();
}

function cleanVisibleText(value, fallback='暂无内容') {
  const text = stripSTA100Result(value).replace(/\s{3,}/g, ' ').trim();
  return text || fallback;
}

function isRecommendationPlaceholder(item={}) {
  const title = String(item.title || '').trim();
  const type = String(item.type || '').trim();
  const desc = String(item.desc || '').trim();
  const content = cleanVisibleText(item.content || '', '');
  const dataStatus = String(item.dataStatus || '').toLowerCase();
  // A real narrative may mention local business data as its evidence source.
  // Do not discard it as a prompt/status placeholder when it has a long body.
  if ((dataStatus === 'generated_narrative' || dataStatus === 'generated_section') && content.length >= 80) return false;
  if (isMetadataTitle(title)) return true;
  if (dataStatus !== 'generated_narrative' && dataStatus !== 'generated_section' && content.length < 120) return true;
  if (!title || title.includes('[STA100_RESULT]') || title.includes('**')) return true;
  if (type.includes('证据缺失') || type.includes('格式说明')) return true;
  return [
    '检索源摘要',
    '运营信号 / 用户偏好',
    '本地业务数据',
    '关注国家 / 主题配置',
    '未随请求附带',
    '没有业务数据文件',
  ].some(marker => title.includes(marker) || desc.includes(marker)) || desc.includes('[STA100_RESULT]');
}

function isNewsPlaceholder(item={}) {
  const title = String(item.title || '').trim();
  const category = String(item.category || '').trim();
  const source = String(item.source || '').trim();
  const summary = String(item.summary || '').trim();
  const dataStatus = String(item.dataStatus || '').toLowerCase();
  if (isMetadataTitle(title)) return true;
  if (!String(item.content || '').trim() && summary.includes('详情请结合来源复核')) return true;
  if (dataStatus !== 'generated_narrative' && dataStatus !== 'generated_section') {
    const contentLength = cleanVisibleText(item.content || '', '').length;
    const summaryLength = cleanVisibleText(item.summary || '', '').length;
    if (contentLength < 80 && summaryLength < 30) return true;
  }
  if (!title || title.includes('[STA100_RESULT]') || summary.includes('[STA100_RESULT]')) return true;
  if (title.includes('以上内容仅为') || title.includes('不构成可独立发布')) return true;
  return category === '备注' && source === 'OpenClaw Agent' && !String(item.sourceUrl || '').trim();
}

function isMetadataTitle(value='') {
  const text = String(value || '').trim().replace(/^[-*•]\s*/, '').replace(/^[*_`]+|[*_`]+$/g, '').trim();
  if (!text) return true;
  return /^(链接|来源|时间|发布时间|相关度|类型|分类)\s*[:：]/i.test(text) || /^(url|link|source|time|relevance|type|category)\s*:/i.test(text) || /^https?:\/\//i.test(text);
}

function recommendationTimestamp(item = {}) {
  const value = item.updatedAt || item.time || '';
  const date = new Date(String(value).replace(' ', 'T'));
  return Number.isNaN(date.getTime()) ? 0 : date.getTime();
}

function recommendationContentPriority(item = {}) {
  const content = recommendationContent(item);
  const dataStatus = String(item.dataStatus || '').toLowerCase();
  if (dataStatus === 'generated_narrative') return 3;
  if (dataStatus === 'generated_section') return 2;
  if (content.length >= 120) return 2;
  if (content.length > 0) return 1;
  return 0;
}

function visibleRecommendations() {
  return recommendations
    .filter(item => !isRecommendationPlaceholder(item))
    .sort((a, b) => {
      const contentPriority = recommendationContentPriority(b) - recommendationContentPriority(a);
      if (contentPriority !== 0) return contentPriority;
      return recommendationTimestamp(b) - recommendationTimestamp(a);
    });
}

function visibleNewsItems() {
  return news.filter(item => !isNewsPlaceholder(item));
}

function newsItemMatchesFilter(item = {}, category = '全部') {
  const label = String(category || '全部').trim();
  if (label === '全部') return true;
  // 优先精确匹配 category 字段
  const itemCategory = String(item.category || '').trim().toLowerCase();
  if (itemCategory === label.toLowerCase()) return true;
  // 再匹配标题和摘要中的关键字
  const text = [item.title, item.summary, item.content, item.source, item.relevance]
    .map(value => cleanVisibleText(value, ''))
    .join(' ')
    .toLowerCase();
  const terms = {
    '欧洲市场': ['europe', '欧洲', '德国', '法国', '意大利', '西班牙', '荷兰', '波兰', '比利时', '瑞典', '奥地利', '挪威', '丹麦', '芬兰', '葡萄牙', '捷克', '匈牙利', '北欧', '欧盟'],
    '法规': ['法规', '规则', '监管', '合规', '政策', 'directive', 'regulation', 'eur-lex', '电池法', '反倾销', '认证', '尽职调查'],
    '智能骑行': ['智能骑行', 'e-bike', 'ebike', '电助力', '电动自行车', '智能设备', '码表', '功率计', 'gps', '传感器', '电子变速'],
    '渠道': ['渠道', '经销商', '分销', 'dealer', 'distribution', 'retail', '零售', '代理', '门店', '售后'],
    '产品': ['产品', '新品', 'launch', '发布', '组件', '电池', '电机', '整车', '配件', '车架', '轮胎', '传动', '头盔'],
  }[label] || [label];
  return terms.some(term => text.includes(String(term).toLowerCase()));
}

function newsContentBlocks(item = {}) {
  const summary = normalizeNewsBody(item.summary || '');
  let content = normalizeNewsBody(item.content || item.body || item.detail || '');
  if (!content) return { summary, content: '' };
  if (!summary) return { summary: '', content };

  const normalizedSummary = summary.replace(/\s+/g, ' ').trim();
  const normalizedContent = content.replace(/\s+/g, ' ').trim();
  if (normalizedContent === normalizedSummary || normalizedSummary.includes(normalizedContent)) {
    return { summary, content: '' };
  }
  // When the full body starts by repeating the summary, remove only that
  // repeated prefix and keep the remaining article detail.
  if (normalizedContent.startsWith(normalizedSummary)) {
    const compactSummary = normalizedSummary.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
    content = content.replace(new RegExp(`^\\s*${compactSummary}\\s*`, 'i'), '').trim();
  }
  if (!content || content.replace(/\s+/g, ' ').trim() === normalizedSummary) {
    return { summary, content: '' };
  }
  return { summary, content };
}

function normalizeNewsBody(value = '') {
  const text = cleanVisibleText(value, '')
    .replace(/(?:摘要|完整内容|详细内容|正文|内容|详情)\s*[:：]\s*/g, '\n')
    .trim();
  if (!text) return '';
  const lines = text.split(/\r?\n/);
  const seen = new Set();
  return lines.filter(line => {
    const normalized = line.replace(/\s+/g, ' ').trim();
    if (!normalized) return true;
    if (/^(?:摘要|完整内容|详细内容|正文|内容|详情)\s*[:：]?$/i.test(normalized)) return false;
    if (seen.has(normalized)) return false;
    seen.add(normalized);
    return true;
  }).join('\n').trim();
}

function detailTextBlock(value) {
  const text = cleanVisibleText(value);
  const blocks = [];
  let paragraph = [];
  let listType = '';
  let listItems = [];
  const flushParagraph = () => {
    if (paragraph.length) {
      blocks.push(`<p>${escapeHTML(paragraph.join('\n')).replace(/\n/g, '<br>')}</p>`);
      paragraph = [];
    }
  };
  const flushList = () => {
    if (!listItems.length) return;
    const tag = listType === 'ordered' ? 'ol' : 'ul';
    blocks.push(`<${tag} class="detail-content-list">${listItems.map(line => `<li>${escapeHTML(line)}</li>`).join('')}</${tag}>`);
    listType = '';
    listItems = [];
  };
  text.split('\n').forEach(rawLine => {
    const line = rawLine.trim();
    if (!line) {
      flushParagraph();
      flushList();
      return;
    }
    const heading = line.match(/^#{2,4}\s+(.+)$/);
    if (heading) {
      flushParagraph();
      flushList();
      blocks.push(`<h3 class="detail-content-heading">${escapeHTML(heading[1])}</h3>`);
      return;
    }
    if (/^-{3,}$/.test(line)) {
      flushParagraph();
      flushList();
      blocks.push('<hr class="detail-content-divider">');
      return;
    }
    const bullet = line.match(/^[-*•]\s+(.+)$/);
    const ordered = line.match(/^\d+[.)]\s+(.+)$/);
    if (bullet || ordered) {
      flushParagraph();
      const nextType = ordered ? 'ordered' : 'unordered';
      if (listType && listType !== nextType) flushList();
      listType = nextType;
      listItems.push((bullet || ordered)[1]);
      return;
    }
    flushList();
    paragraph.push(line);
  });
  flushParagraph();
  flushList();
  return blocks.join('');
}

function recommendationContent(item = {}) {
  const candidates = [item.detail, item.content, item.summary, item.body, item.why, item.reason, item.desc]
    .map(value => cleanVisibleText(value, ''))
    .filter(Boolean)
    .sort((a, b) => b.length - a.length);
  const merged = [];
  for (const candidate of candidates) {
    const normalized = candidate.replace(/[.…]+$/g, '').trim();
    if (merged.some(existing => existing.includes(candidate) || (normalized.length >= 24 && existing.includes(normalized)))) continue;
    merged.push(candidate);
  }
  return merged.join('\n\n') || '当前推荐未返回详细内容，请结合来源和关注条件人工复核。';
}

function recommendationPreview(item = {}) {
  const text = cleanVisibleText(recommendationContent(item), '当前推荐未返回详细内容。')
    .replace(/\s+/g, ' ')
    .trim();
  const limit = 220;
  return text.length > limit ? `${text.slice(0, limit)}…` : text;
}

function recommendationDisplayTitle(item = {}) {
  const title = cleanVisibleText(item.title, '推荐详情');
  if (title.length <= 48) return title;
  const split = title.match(/^([^:：]{2,24})[:：]\s*(.+)$/);
  if (split) return split[1];
  return `${title.slice(0, 46)}…`;
}

function newsFullContent(item = {}) {
  return cleanVisibleText(item.content || item.summary || item.body || item.detail, '当前新闻没有返回可展示的详细内容。');
}

function newsPreview(item = {}, limit = 180) {
  const text = cleanVisibleText(item.summary || item.content || item.body || item.detail, '暂无摘要。')
    .replace(/\s+/g, ' ')
    .trim();
  return text.length > limit ? `${text.slice(0, limit)}…` : text;
}

function humanDataStatusLabel(value) {
  switch (String(value || '').trim().toLowerCase()) {
    case 'generated_narrative':
      return '智能体长文摘要';
    case 'generated_section':
      return '智能体摘要拆分条目';
    case 'generated':
      return '智能体结构化结果';
    case 'system_defaults_only':
      return '系统默认配置';
    case 'local_business_data':
      return '本地业务数据';
    case 'cached':
      return '本机缓存';
    default:
      return cleanVisibleText(value, '未标记');
  }
}

function recommendationDetail(id) {
  const displayItems = visibleRecommendations();
  const r = displayItems.find(item => item.id === id) || displayItems[0];
  if (!r) return;
  const title = cleanVisibleText(r.title, '推荐详情');
  let why = cleanVisibleText(r.why || r.reason || '', '');
  let detail = cleanVisibleText(r.detail || r.content || '', '');
  if (why && detail && why.replace(/\s+/g, ' ').trim() === detail.replace(/\s+/g, ' ').trim()) {
    why = '';
  }
  const content = recommendationContent(r);
  const basis = [
    ['关注国家', state.newsCountries || '未配置'],
    ['关注主题', state.newsTopics || '未配置'],
    ['指定来源', state.newsSources || '未配置'],
    ['生成方式', r.generatedBy || '本地业务数据整理后调用 OpenClaw 推荐 Agent'],
  ];
  const sourceAction = r.sourceUrl ? `<div class="inline-actions" style="margin-top:16px"><button class="button" data-action="recommend-source-link" data-id="${escapeAttr(r.id)}">${icon('external-link')}查看原文</button></div>` : '';
  const contentBody = why || detail ? `<div class="recommendation-detail-stack">${why ? `<section class="recommendation-detail-block"><div class="detail-content-label">为什么推荐</div><div class="recommendation-detail-summary">${detailTextBlock(why)}</div></section>` : ''}${detail ? `<section class="recommendation-detail-block"><div class="detail-content-label">详情</div><div class="recommendation-detail-summary">${detailTextBlock(detail)}</div></section>` : ''}</div>` : `<div class="recommendation-detail-summary">${detailTextBlock(content)}</div>`;
  openDrawer({title,eyebrow:`${r.type || '为你推荐'} / ${r.source || 'OpenClaw Agent'}`,body:`<div class="recommendation-detail-topline"><span class="badge blue">${escapeHTML(r.type||'推荐')}</span><span class="secondary-text">${escapeHTML(overviewTime(r.updatedAt || r.time || ''))}</span></div><article class="recommendation-detail-content"><h2>${escapeHTML(title)}</h2><div class="detail-content-label">推荐内容</div>${contentBody}</article><div class="divider-title" style="margin:20px 0 12px">推荐依据</div><div class="detail-grid">${basis.map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div><div class="divider-title" style="margin:20px 0 12px">来源信息</div><div class="detail-grid">${[['来源',r.source||'未返回'],['更新时间',overviewTime(r.updatedAt || r.time || '')],['类型',r.type||'未返回'],['原文链接',r.sourceUrl||'未返回'],['生成状态',humanDataStatusLabel(r.dataStatus)]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v)}</strong></div>`).join('')}</div>${sourceAction}`});
}

function openRecommendationSource(id) {
  const item = visibleRecommendations().find(entry => entry.id === id);
  if (item?.sourceUrl) {
    window.open(item.sourceUrl, '_blank', 'noopener,noreferrer');
    return;
  }
  toast('原文链接未配置', '当前推荐没有返回来源 URL，已展示本机缓存的完整内容。', 'warning');
}

function newsDetail(title) {
  const displayItems = visibleNewsItems();
  const n=displayItems.find(x=>x.title===title)||displayItems[0];
  if (!n) return;
  const sourceAction = n.sourceUrl ? `<button class="button" data-action="news-source-link" data-title="${escapeAttr(n.title)}">${icon('external-link')}查看原文</button>` : `<span class="secondary-text">未返回原文链接，已展示本机缓存完整内容</span>`;
  const blocks = newsContentBlocks(n);
  const bodyParts = [];
  if (blocks.summary) {
    bodyParts.push(`<div class="recommendation-detail-content news-detail-content"><h2>${escapeHTML(cleanVisibleText(n.title,'新闻详情'))}</h2><div class="detail-content-label">摘要</div><div class="recommendation-detail-summary">${detailTextBlock(blocks.summary)}</div></div>`);
  }
  if (blocks.content) {
    bodyParts.push(`<div class="recommendation-detail-content news-detail-content"><div class="detail-content-label">完整内容</div><div class="recommendation-detail-summary">${detailTextBlock(blocks.content)}</div></div>`);
  }
  if (!bodyParts.length) {
    bodyParts.push(`<div class="recommendation-detail-content news-detail-content"><h2>${escapeHTML(cleanVisibleText(n.title,'新闻详情'))}</h2><div class="detail-content-label">内容</div><div class="recommendation-detail-summary">${detailTextBlock(cleanVisibleText(n.content || n.summary || '', '暂无可展示内容。'))}</div></div>`);
  }
  openDrawer({title:cleanVisibleText(n.title,'新闻详情'),eyebrow:`${n.category || '行业资讯'} / ${n.source || 'OpenClaw Agent'}`,body:`<div class="recommendation-detail-topline"><span class="badge green">相关度 ${escapeHTML(n.relevance||'待复核')}</span><span class="secondary-text">${escapeHTML(overviewTime(n.updatedAt || n.time || ''))}</span></div>${bodyParts.join('')}<div class="divider-title" style="margin:20px 0 12px">来源信息</div><div class="detail-grid">${[['来源',n.source||'未返回'],['获取时间',overviewTime(n.updatedAt || n.time || '')],['信息类别',n.category||'行业资讯'],['原文链接',n.sourceUrl||'未返回'],['生成方式',n.generatedBy||'OpenClaw 行业新闻 Agent'],['内容状态',humanDataStatusLabel(n.dataStatus)]].map(([l,v])=>`<div class="detail-field"><label>${l}</label><strong>${escapeHTML(v||'未返回')}</strong></div>`).join('')}</div><div class="inline-actions" style="margin-top:18px">${sourceAction}<button class="button primary" data-action="news-todo" data-title="${escapeAttr(n.title)}">${icon('list-plus')}生成待办</button></div>`});
}

async function refreshNews() {
  if (state.newsRefreshLoading) return;
  state.newsRefreshLoading = true;
  renderPage();
  try {
    const result = await apiFetch('/api/v1/news/refresh',{method:'POST',body:'{}'});
    if (Array.isArray(result.news)) replaceRecords(news, result.news);
    if (result.automation) state.overviewAutomation = result.automation;
    if (result.job) upsertRecord(scheduledJobs, result.job);
    renderPage();
    toast(result.completed === false ? '新闻刷新已提交' : '新闻刷新完成', result.message || '已同步 OpenClaw 行业新闻 Agent 返回内容。', result.completed === false ? 'warning' : 'success');
  } catch(error) {
    toast('新闻刷新失败', error.message, 'warning');
  } finally {
    state.newsRefreshLoading = false;
    renderPage();
  }
}

async function refreshRecommendations() {
  if (state.recommendRefreshLoading) return;
  state.recommendRefreshLoading = true;
  renderPage();
  try {
    const result = await apiFetch('/api/v1/overview/recommendations/refresh', { method:'POST', body:'{}' });
    if (Array.isArray(result.recommendations)) replaceRecords(recommendations, result.recommendations);
    if (result.automation) state.overviewAutomation = result.automation;
    if (result.job) upsertRecord(scheduledJobs, result.job);
    renderPage();
    toast(result.completed === false ? '推荐刷新已提交' : '推荐刷新完成', result.message || '已同步 OpenClaw 推荐 Agent 返回内容。', result.completed === false ? 'warning' : 'success');
  } catch(error) {
    toast('推荐刷新失败', error.message, 'warning');
  } finally {
    state.recommendRefreshLoading = false;
    renderPage();
  }
}

function openNewsSource(title) {
  const item=news.find(entry=>entry.title===title);
  if(item?.sourceUrl){window.open(item.sourceUrl,'_blank','noopener,noreferrer');return;}
  toast('原文链接未配置','当前缓存记录没有来源 URL，待来源白名单和抓取规则确认后补充。','warning');
}

async function createNewsTodo(title) {
  try {
    const result=await apiFetch('/api/v1/tasks',{method:'POST',body:JSON.stringify({title:`跟进资讯：${title}`,source:'industry_news'})});
    await loadTasks(true);
    if (state.page === 'overview') renderPage();
    toast('待办已生成',result.message||`已创建：${title}`,'success');
  }
  catch(error){toast('待办暂不可生成',error.message,'warning');}
}

function currentModelConfiguration() {
  const model = document.getElementById('modelVersionSelect')?.value || state.modelDraftKey || '';
  const apiKey = document.getElementById('modelAPIKey')?.value.trim() || '';
  const endpointMode = document.getElementById('modelEndpointMode')?.value || modelEndpointModeFor(modelProvider(model));
  const baseModels = state.modelDraftMode === 'edit'
    ? configuredModelKeys().filter(key => key && key !== state.modelDraftOriginalKey)
    : configuredModelKeys();
  const selectedModels = model ? [...new Set([...baseModels, model])] : baseModels;
  // Saving a model must not implicitly replace an already working default.
  // The default-model action in the list remains an explicit operation.
  const defaultModel = '';
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
  const rows = businessOpenClawAgents();
  if (state.openClawAgentsLoading && !rows.length) return `<div class="empty-state compact-empty">${icon('loader-circle')}<div><h3>正在读取 OpenClaw Agent</h3><p>请稍候。</p></div></div>`;
  return `<div class="manager-summary"><span>${icon('check-circle-2')} 已注册 <strong>${rows.length}</strong> / 24 个 STA-100 业务 Agent</span><span>数据来源：OpenClaw agents list</span></div><div class="data-wrap"><table class="data-table agent-manager-table"><thead><tr><th>业务智能体</th><th>Agent ID</th><th>模型</th><th>状态</th></tr></thead><tbody>${rows.map(agent=>`<tr><td><span class="agent-manager-name"><span class="agent-emoji" aria-hidden="true">${escapeHTML(agent.identityEmoji || '🤖')}</span><strong>${escapeHTML(agent.identityName || agent.name || agent.id)}</strong></span></td><td>${escapeHTML(agent.id)}</td><td>${escapeHTML(agent.model)}</td><td>${badge('Active')}</td></tr>`).join('') || `<tr><td colspan="4"><div class="empty-state"><p>未读取到 STA-100 业务 Agent。</p></div></td></tr>`}</tbody></table></div>`;
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
    toast('Agent 同步完成', `OpenClaw 当前共 ${data.count} 个业务 Agent，系统 Agent 已隐藏。`);
  } catch (error) {
    button.disabled = false;
    button.innerHTML = `${icon('refresh-cw')}重新同步`;
    applyIcons();
    toast('Agent 同步失败', error.message, 'warning');
  }
}

function schedulerAgentCatalog() {
  const knownNames = {
    'sta100-coordinator': 'STA-100 统一协调 Agent',
    'sta100-knowledge': 'STA-100 本地知识库 Agent',
    'sta100-recommend-curator': 'STA-100 推荐整理 Agent',
    'sta100-news-curator': 'STA-100 行业新闻 Agent',
    'market-analyzer': '市场分析 Agent',
  };
  const agents = businessOpenClawAgents().map(agent => ({
    id: agent.id,
    name: agent.identityName || agent.name || knownNames[agent.id] || agent.id,
  }));
  Object.entries(knownNames).forEach(([id, name]) => {
    if (!agents.some(agent => agent.id === id)) agents.push({ id, name });
  });
  return agents;
}

function schedulerAgentLabel(agentID) {
  const id = String(agentID || '').trim();
  if (!id) return '未指定 Agent';
  const agent = schedulerAgentCatalog().find(item => item.id === id);
  return agent ? `${agent.name}（${agent.id}）` : id;
}

function schedulerAgentSelect(selected) {
  const options = schedulerAgentCatalog().map(agent => `<option value="${escapeAttr(agent.id)}" ${agent.id === selected ? 'selected' : ''}>${escapeHTML(agent.name)}（${escapeHTML(agent.id)}）</option>`).join('');
  return `<div class="form-field"><label>执行 Agent <span class="required">*</span></label><select class="select" id="jobAgent">${options}</select><small>显示名称和 OpenClaw Agent ID，保存时写入 ID。</small></div>`;
}

function schedulerReadonlyField(label, value, id, full = false, note = '') {
  return `<div class="form-field ${full ? 'full' : ''}"><label for="${escapeAttr(id)}">${label}</label><input class="input" id="${escapeAttr(id)}" value="${escapeAttr(value)}" readonly>${note ? `<small>${escapeHTML(note)}</small>` : ''}</div>`;
}

function schedulerReadonlyPrompt(value) {
  return `<div class="form-field full"><label for="jobPrompt">执行 Prompt <span class="required">*</span></label><textarea class="textarea scheduler-readonly-prompt" id="jobPrompt" maxlength="8000" rows="5" readonly>${escapeHTML(value || '')}</textarea><small>内置任务的 Prompt 由系统固定；每日推荐和行业新闻请在概览页/行业新闻页的设置中修改业务条件。</small></div>`;
}

function schedulerScheduleValueOptions(kind, selected) {
  const every = [
    ['1h','每 1 小时'],['2h','每 2 小时'],['3h','每 3 小时'],['6h','每 6 小时'],['8h','每 8 小时'],['12h','每 12 小时'],['24h','每 24 小时'],
  ];
  const cron = [
    ['0 8 * * *','每天 08:00'],['0 9 * * *','每天 09:00'],['0 9 * * 1','每周一 09:00'],['0 8 1 * *','每月 1 日 08:00'],
  ];
  const at = [['+10m','10 分钟后执行一次'],['+1h','1 小时后执行一次'],['+24h','24 小时后执行一次']];
  const options = kind === 'every' ? every : kind === 'at' ? at : cron;
  const hasSelected = selected && options.some(([value])=>value===selected);
  return `${hasSelected?'':selected?`<option value="${escapeAttr(selected)}" selected>${escapeHTML(selected)}</option>`:''}${options.map(([value,label])=>`<option value="${escapeAttr(value)}" ${value===selected?'selected':''}>${escapeHTML(label)} · ${escapeHTML(value)}</option>`).join('')}`;
}

function schedulerScheduleValueField(kind, selected) {
  return `<div class="form-field"><label for="jobScheduleValue">调度值 <span class="required">*</span></label><select class="select" id="jobScheduleValue">${schedulerScheduleValueOptions(kind, selected)}</select><small>请选择固定周期，避免 Cron 格式手输错误。</small></div>`;
}

function schedulerForm(id='') {
  const job=scheduledJobs.find(item=>item.id===id)||{};
  state.formContext={type:'job',id:job.id||''};
  const builtIn = Boolean(job.id && job.builtIn);
  const scheduleKind = job.scheduleKind || (job.scheduleValue ? 'cron' : 'cron');
  const scheduleValue = job.scheduleValue || job.schedule || '0 8 * * *';
  const identityFields = builtIn
    ? `${schedulerReadonlyField('任务名称', job.name || '', 'jobName', true, '内置任务名称由系统固定。')}${schedulerReadonlyField('执行 Agent', schedulerAgentLabel(job.agentId), 'jobAgent', false, '内置任务 Agent 由系统固定。')}${schedulerReadonlyPrompt(job.prompt)}`
    : `${inputField('任务名称',job.name||'',true,true,'text','jobName')}${schedulerAgentSelect(job.agentId||'sta100-coordinator')}<div class="form-field full"><label for="jobPrompt">执行 Prompt <span class="required">*</span></label><textarea class="textarea" id="jobPrompt" maxlength="8000" rows="5" placeholder="描述任务每次实际需要执行的内容">${escapeHTML(job.prompt||'')}</textarea><small>该 Prompt 会原样作为 OpenClaw Agent Cron 的执行消息。</small></div>`;
  const builtInNote = builtIn ? `<div class="model-warning full"><span>${icon('lock-keyhole')} 这是内置定时任务。名称、执行 Agent 和 Prompt 仅供查看；每日推荐和行业新闻的关注条件、来源、数量和频率请通过对应页面设置修改。</span></div>` : '';
  openModal({title:job.id?'编辑定时任务':'新增定时任务',eyebrow:'设置 / 定时任务',body:`<div class="form-grid">${identityFields}${selectField('任务类型',['recommendations','weekly_report','news','index','custom'],true,'jobKind',job.kind||'custom')}${selectField('调度类型',['every','cron','at'],false,'jobScheduleKind',scheduleKind)}${schedulerScheduleValueField(scheduleKind,scheduleValue)}${inputField('时区（Cron 可选）',job.timezone||'',false,false,'text','jobTimezone')}${inputField('任务说明',job.description||'',true,true,'text','jobDescription')}${builtInNote}<div class="form-field full"><label style="display:flex;align-items:center;gap:8px"><input id="jobEnabled" type="checkbox" ${job.id?!job.enabled?'':'checked':'checked'}> 启用任务</label></div>${job.id&&!job.builtIn?`<div class="form-field full"><button type="button" class="button danger" data-action="delete-schedule" data-id="${escapeAttr(job.id)}">${icon('trash-2')}删除自定义任务</button></div>`:''}</div>`,footer:formFooter('保存任务','save-schedule')});
  document.getElementById('jobScheduleKind')?.addEventListener('change', e => {
    const wrapper = document.getElementById('jobScheduleValue')?.closest('.form-field');
    if (wrapper) wrapper.outerHTML = schedulerScheduleValueField(e.target.value, '');
  });
}

async function saveSchedule() {
  const existing=scheduledJobs.find(item=>item.id===state.formContext?.id);
  const scheduleKind=formText('jobScheduleKind');
  const scheduleValue=formText('jobScheduleValue');
  const payload={...(existing||{}),name:formText('jobName'),kind:formText('jobKind'),schedule:scheduleValue,scheduleKind,scheduleValue,timezone:formText('jobTimezone'),description:formText('jobDescription'),agentId:formText('jobAgent'),prompt:formText('jobPrompt'),enabled:Boolean(document.getElementById('jobEnabled')?.checked),status:existing?.status||'unsynced'};
  if(!payload.name||!payload.kind||!payload.scheduleKind||!payload.scheduleValue||!payload.prompt){toast('保存失败','任务名称、类型、调度值和执行 Prompt 不能为空。','warning');return;}
  try {const record=await apiFetch(existing?'/api/v1/jobs':'/api/v1/jobs',{method:existing?'PATCH':'POST',body:JSON.stringify(payload)});upsertRecord(scheduledJobs,record);closeModal();renderPage();toast(existing?'任务已更新':'任务已创建',record.syncMessage||record.name,record.syncStatus==='synced'?'success':'warning');}
  catch(error){toast('任务保存失败',error.message,'warning');}
}

async function deleteSchedule(id) {
  const job=scheduledJobs.find(item=>item.id===id);if(!job||!window.confirm(`确定删除任务“${job.name}”吗？`))return;
  try {await apiFetch(`/api/v1/jobs/${encodeURIComponent(id)}`,{method:'DELETE'});removeRecord(scheduledJobs,id);closeModal();renderPage();toast('任务已删除',job.name);}
  catch(error){toast('任务删除失败',error.message,'warning');}
}

async function toggleSchedule(id, enabled) {
  if (isScheduleActionLoading(id)) return;
  setScheduleActionLoading(id, 'toggle');
  try {
    const record = await apiFetch(`/api/v1/jobs/${encodeURIComponent(id)}/${enabled?'enable':'disable'}`, {method:'POST', body:'{}'});
    upsertRecord(scheduledJobs, record);
    renderPage();
    toast(enabled?'任务已开启':'任务已关闭', record.syncMessage || record.name, 'success');
  } catch (error) {
    toast(enabled?'任务开启失败':'任务关闭失败', error.message, 'warning');
  } finally {
    clearScheduleActionLoading(id, 'toggle');
  }
}

async function runSchedule(id) {
  const job = scheduledJobs.find(item=>item.id===id);
  if (!job) return;
  if (isScheduleActionLoading(id)) return;
  setScheduleActionLoading(id, 'run');
  try {
    const result = await apiFetch(`/api/v1/jobs/${encodeURIComponent(id)}/run`, {method:'POST', body:'{}'});
    if (result.job) upsertRecord(scheduledJobs, result.job);
    renderPage();
    toast('任务已提交执行', result.run?.message || '请稍后查看运行记录。', 'success');
  } catch (error) {
    toast('任务执行失败', error.message, 'warning');
  } finally {
    clearScheduleActionLoading(id, 'run');
  }
}

async function showScheduleRuns(id) {
  const job = scheduledJobs.find(item=>item.id===id);
  if (!job) return;
  openModal({title:`${job.name} · 运行记录`,eyebrow:'设置 / OpenClaw Cron',body:`<div class="empty-state">${icon('loader-circle')}<div><h3>正在读取运行记录</h3><p>数据来自 OpenClaw Cron。</p></div></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`});
  try {
    const result = await apiFetch(`/api/v1/jobs/${encodeURIComponent(id)}/runs`);
    const entries = result.entries || [];
    document.getElementById('modalBody').innerHTML = entries.length ? `<div class="scheduler-run-list">${entries.map(run=>`<div class="scheduler-run-item"><div><strong>${escapeHTML(run.statusLabel||run.status||'状态待确认')}</strong><span>${escapeHTML(formatLocalizedDateTime(run.runAtMs)||'未记录时间')}</span></div><p>${escapeHTML(run.message||'无摘要')}</p><small>${run.durationMs?`耗时 ${escapeHTML(formatDurationMs(run.durationMs))} · `:''}${escapeHTML(run.deliveryStatus||'未请求推送')}</small></div>`).join('')}</div>` : `<div class="empty-state"><p>该任务还没有运行记录。</p></div>`;
    applyIcons();
  } catch (error) {
    document.getElementById('modalBody').innerHTML = `<div class="model-warning error">${icon('triangle-alert')} ${escapeHTML(error.message)}</div>`;
    applyIcons();
  }
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
  setTimeout(()=>node.remove(),7000);
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
    refreshQuoteDraftTotals();
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
  if(e.target.id==='quoteFreight'||e.target.id==='quoteTax'||e.target.id==='quoteCurrency'){refreshQuoteDraftTotals();return;}
  const input=e.target.closest('[data-relation-input]');
  if(!input)return;
  const target=input.dataset.relationInput;
  const items=target==='quoteCustomer'||target==='orderCustomer' ? customers.filter(c=>!c.archived).map(c=>c.name) : target==='orderQuote' ? quotes.map(q=>`${q.id} · ${q.customer}`) : orders.map(o=>`${o.id} · ${o.customer}`);
  const options=document.getElementById(`${target}Options`);
  if(options){
    options.hidden = false;
    options.innerHTML=relationOptions(target,items,input.value);
  }
});

document.addEventListener('focusin', e => {
  const input = e.target.closest('[data-relation-input]');
  if (!input) return;
  const target = input.dataset.relationInput;
  const items = target==='quoteCustomer'||target==='orderCustomer'
    ? customers.filter(c=>!c.archived).map(c=>c.name)
    : target==='orderQuote'
      ? quotes.map(q=>`${q.id} · ${q.customer}`)
      : orders.map(o=>`${o.id} · ${o.customer}`);
  const options = document.getElementById(`${target}Options`);
  if (options) {
    options.hidden = false;
    options.innerHTML = relationOptions(target, items, '__all__');
  }
});

document.addEventListener('change', e => {
  if(e.target.id==='modelFamilySelect'){updateModelFamilySelection();return;}
  if(e.target.id==='modelVersionSelect'){updateModelSelection();return;}
  if(e.target.id==='customerCountryForm'){refreshCustomerCityOptions(e.target.value);return;}
  const mediaCategoryFilter=e.target.closest('[data-news-media-category-filter]');
  if(mediaCategoryFilter){
    const selected=[...document.querySelectorAll('[data-news-media-category-filter]:checked')].map(item=>item.dataset.newsMediaCategoryFilter);
    state.newsMediaFilterCategories=selected.length?selected:newsMediaCategories().map(item=>item.key);
    openNewsSettings();
    return;
  }
  const mediaSource=e.target.closest('[data-news-media-id]');
  if(mediaSource){
    const id=mediaSource.dataset.newsMediaId;
    state.newsMediaIDs=state.newsMediaIDs||[];
    if(mediaSource.checked&&!state.newsMediaIDs.includes(id))state.newsMediaIDs=[...state.newsMediaIDs,id];
    if(!mediaSource.checked)state.newsMediaIDs=state.newsMediaIDs.filter(item=>item!==id);
    openNewsSettings();
    return;
  }
  const quoteProduct=e.target.closest('[data-quote-line-field="productId"]');
  if(quoteProduct){const line=state.quoteDraftLines[Number(quoteProduct.dataset.index)];const product=productByID(quoteProduct.value);Object.assign(line,{productId:product.id,unitPrice:moneyNumber(product.price)});renderQuoteDraftLines();return;}
  const orderProduct=e.target.closest('[data-order-line-field="productId"]');
  if(orderProduct){const line=state.orderDraftLines[Number(orderProduct.dataset.index)];const product=productByID(orderProduct.value);Object.assign(line,{productId:product.id,unitPrice:moneyNumber(product.price)});renderOrderDraftLines();return;}
  if(e.target.id==='orderSource'&&e.target.value==='手动创建'){const quote=document.getElementById('orderQuote');if(quote){quote.value='';document.getElementById('orderQuoteOptions').innerHTML='';}}
});

document.addEventListener('click', e => {
  const pageTarget=e.target.closest('[data-page]');
  if(pageTarget){setPage(pageTarget.dataset.page);return;}
  if(!e.target.closest('.relation-picker')) document.querySelectorAll('.relation-options').forEach(options => { options.hidden = true; options.innerHTML = ''; });
  const lang=e.target.closest('[data-lang]');
  if(lang){state.lang=lang.dataset.lang;document.documentElement.lang=state.lang==='en'?'en':'zh-CN';document.querySelectorAll('[data-lang]').forEach(x=>x.classList.toggle('active',x===lang));applyTranslations();renderPage();toast('语言已切换',state.lang==='zh'?'当前界面为中文。':'Interface is now English.');return;}
  const cat=e.target.closest('[data-agent-category]'); if(cat){state.agentCategory=cat.dataset.agentCategory;renderPage();return;}
  const documentType=e.target.closest('[data-document-type]'); if(documentType){state.documentType=documentType.dataset.documentType;renderPage();return;}
  const qv=e.target.closest('[data-quote-view]'); if(qv){state.quoteView=qv.dataset.quoteView;renderPage();return;}
  const pv=e.target.closest('[data-product-view]'); if(pv){state.productView=pv.dataset.productView;renderPage();return;}
  const customerTab=e.target.closest('[data-customer-tab]'); if(customerTab){customerDetail(customerTab.dataset.customerId,customerTab.dataset.customerTab);return;}
  const leadTab=e.target.closest('[data-lead-tab]'); if(leadTab){leadDetail(leadTab.dataset.leadId,leadTab.dataset.leadTab);return;}
  const supplierTab=e.target.closest('[data-supplier-tab]'); if(supplierTab){supplierDetail(supplierTab.dataset.supplierId,supplierTab.dataset.supplierTab);return;}
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
    'open-overview-recommendations':()=>{state.recExpanded=true;renderPage();setTimeout(()=>document.querySelector('.recommendation-list')?.scrollIntoView({behavior:'smooth', block:'start'}),0);},
    'open-overview-news':()=>setPage('news'),
    'recommend-detail':()=>recommendationDetail(el.dataset.id),'recommend-source-link':()=>openRecommendationSource(el.dataset.id),
    'recommend-settings':openRecommendationSettings,
    'refresh-recommendations':()=>void refreshRecommendations(),
    'oem-preset':()=>{state.oemQuery=el.dataset.value;const input=document.getElementById('oemQuery');if(input)input.value=state.oemQuery;renderPage();},
    'oem-category':()=>{state.oemCategory=el.dataset.value||'全部';renderPage();},
    'oem-run':()=>void runOEMMatch(),
    'oem-export':oemExport,
    'unified-customer-search':()=>void runUnifiedCustomerSearch(),
    'unified-customer-detail':()=>unifiedCustomerDetail(el.dataset.name),
    'local-discovery-search':()=>void runLocalDiscovery(),
    'save-discovery-settings':()=>void saveDiscoverySettings(),
    'local-lead-detail':()=>localLeadDetail(el.dataset.name),'add-to-leads':()=>addDiscoveryLeadToLeads(el.dataset.name),'add-all-discovery-leads':()=>void addAllDiscoveryLeadsToLeads(),
    'agent-chat':()=>showAgentChat(Number(el.dataset.agent),el.dataset.prompt||''),
    'confirm-agent-answer':()=>void confirmAgentAnswer(Number(el.dataset.agent),Number(el.dataset.messageIndex)),
    'weekly-report':()=>void generateWeeklyReport(),
    'download-weekly-report':()=>state.lastWeeklyReport&&downloadText(`STA100-Agent-Weekly-${new Date().toISOString().slice(0,10)}.md`,state.lastWeeklyReport.markdown),
    'agent-manage':openAgentManager,
    'sync-openclaw-agents':()=>syncOpenClawAgents(el),
    'new-customer':()=>newCustomerForm(), 'edit-customer':()=>newCustomerForm(customers.find(c=>c.id===el.dataset.id)), 'customer-detail':()=>void customerDetail(el.dataset.id), 'delete-customer':()=>deleteCustomer(el.dataset.id), 'customer-more':()=>void customerDetail(el.dataset.id,'activity'), 'customer-to-lead':()=>customerToLead(el.dataset.id),
    'customer-communications':()=>{closeModal();void customerDetail(el.dataset.id,'activity');},'new-customer-communication':()=>customerCommunicationForm(el.dataset.id),'save-customer-communication':()=>void saveCustomerCommunication(el.dataset.id),'cancel-customer-communication':()=>{closeModal();void customerDetail(el.dataset.id,'activity');},
    'export-customers':()=>window.location.assign('/api/v1/accounts/export'),
    'column-settings':openCustomerColumnSettings,'save-customer-columns':saveCustomerColumns,
    'new-lead':()=>newLeadForm(),'edit-lead':()=>newLeadForm(leads.find(l=>l.id===el.dataset.id)),'lead-detail':()=>leadDetail(el.dataset.id),'convert-lead':()=>convertLead(el.dataset.id),'delete-lead':()=>deleteLead(el.dataset.id),'new-lead-communication':()=>leadCommunicationForm(el.dataset.id),'save-lead-communication':()=>void saveLeadCommunication(el.dataset.id),'cancel-lead-communication':()=>{closeModal();void leadDetail(el.dataset.id,'activity');},
    'new-quote':()=>void newQuoteForm(null,el.dataset.customer||''),'quote-detail':()=>quoteDetail(el.dataset.id),'edit-quote':()=>void newQuoteForm(quotes.find(q=>q.id===el.dataset.id)),'delete-quote':()=>void deleteQuote(el.dataset.id),'download-quote':()=>void downloadQuote(el.dataset.id),'convert-order':()=>void convertQuoteToOrder(el.dataset.id),'send-quote-email':()=>openBusinessEmail('quote',el.dataset.id),
    'quote-metric-filter':()=>{state.quoteStatus=el.dataset.status||'all';renderPage();},
    'new-order':()=>newOrderForm(),'edit-order':()=>newOrderForm(orders.find(o=>o.id===el.dataset.id)),'delete-order':()=>deleteOrder(el.dataset.id),'order-detail':()=>orderDetail(el.dataset.id),'download-order':()=>void downloadOrder(el.dataset.id),'send-order-email':()=>openBusinessEmail('order',el.dataset.id),
    'open-email-settings':()=>openEmailSettings(),'cancel-email-settings':closeModal,'save-email-settings':()=>void saveEmailSettings(el),'confirm-business-email':()=>void submitBusinessEmail(el.dataset.kind,el.dataset.id,el),
    'generate-docs':()=>generateDocs(el.dataset.id),'new-document':()=>generateDocs(''),'edit-document':()=>generateDocs('',documents.find(d=>d.id===el.dataset.id)),'delete-document':()=>deleteDocument(el.dataset.id),'download-document':()=>void downloadDocument(el.dataset.id),'download-document-format':()=>void downloadDocument(el.dataset.id,el.dataset.format),'template-center':()=>void templateCenter(el.dataset.kind),'document-detail':()=>documentDetail(el.dataset.id),
    'clear-document-filters':()=>{state.documentSearch='';state.documentType='all';state.documentStatus='all';renderPage();},
    'upload-template-image':()=>document.getElementById('templateImageInput')?.click(),
    'upload-template-file':()=>document.getElementById('templateFileInput')?.click(),
    'refresh-templates':()=>void templateCenter(state.templateKind),
    'new-product':()=>newProductForm(),'edit-product':()=>newProductForm(products.find(p=>p.id===el.dataset.id)),'product-detail':()=>productDetail(el.dataset.id),'delete-product':()=>deleteProduct(el.dataset.id),'save-product':saveProduct,'import-products':productImportModal,'choose-product-import':()=>document.getElementById('productImportInput')?.click(),'toggle-product-sort':()=>{state.productSort=state.productSort==='stockAsc'?'stockDesc':'stockAsc';renderPage();},
    'new-supplier':()=>newSupplierForm(),'edit-supplier':()=>newSupplierForm(suppliers.find(s=>s.id===el.dataset.id)),'supplier-detail':()=>supplierDetail(el.dataset.id),'delete-supplier':()=>deleteSupplier(el.dataset.id),'export-suppliers':()=>window.location.assign('/api/v1/suppliers/export'),'supplier-tab':()=>{supplierDetail(el.dataset.supplierId,el.dataset.supplierTab);},'new-supplier-communication':()=>supplierCommunicationForm(el.dataset.id),'save-supplier-communication':()=>void saveSupplierCommunication(el.dataset.id),'cancel-supplier-communication':()=>{closeModal();void supplierDetail(el.dataset.id,'activity');},
    'upload-file':uploadFileModal,'choose-file':()=>document.getElementById('privateFileInput')?.click(),'upload-private-file':()=>void uploadPrivateFile(),'save-file-metadata':()=>void saveFileMetadata(),'file-preview':()=>filePreview(el.dataset.id),'file-download':()=>fileDownload(el.dataset.id),'file-edit':()=>fileEdit(el.dataset.id),'file-summary':()=>void fileSummary(el.dataset.id),'file-more':()=>openModal({title:'文件更多操作',eyebrow:'数据库 / 文件操作',body:`<div class="filter-row"><button class="button" data-action="file-summary" data-id="${escapeAttr(el.dataset.id)}">${icon('sparkles')}摘要状态</button><button class="button" data-action="file-download" data-id="${escapeAttr(el.dataset.id)}">${icon('download')}下载文件</button><button class="button" data-action="file-reindex" data-id="${escapeAttr(el.dataset.id)}">${icon('refresh-cw')}重新索引</button><button class="button danger" data-action="file-archive" data-id="${escapeAttr(el.dataset.id)}">${icon('trash-2')}删除文件</button></div>`,footer:`<button class="button" data-action="close-modal">关闭</button>`}),'file-archive':()=>void archivePrivateFile(el.dataset.id),'file-reindex':()=>void reindexPrivateFile(el.dataset.id),
    'open-category':()=>{state.fileSearch=el.dataset.category;renderPage();},
    'agent-backup':()=>void backupAgents(),
    'tag-manage':()=>openModal({title:'当前文件标签',eyebrow:'数据库',body:`<div class="filter-row">${[...new Set(files.flatMap(file=>file.tags||[]))].map(v=>`<span class="filter-chip active">${escapeHTML(v)}</span>`).join('')||'<span class="secondary-text">暂无标签</span>'}</div><p class="secondary-text" style="margin-top:14px">标签通过每个文件的“编辑信息”维护，修改后立即保存到本机数据库。</p>`,footer:`<button class="button" data-action="close-modal">关闭</button>`}),
    'news-detail':()=>newsDetail(el.dataset.title),'toggle-news':()=>{state.newsExpanded=!state.newsExpanded;renderPage();},'news-filter':()=>{state.newsCategory=el.dataset.category||'全部';state.newsExpanded=false;renderPage();},
    'news-source-link':()=>openNewsSource(el.dataset.title),'news-todo':()=>void createNewsTodo(el.dataset.title),
    'refresh-news':()=>void refreshNews(),
    'news-sources':openNewsSettings,'save-recommendation-settings':()=>void saveRecommendationSettings(),
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
    'refresh-channel-skill':async()=>{await loadChannelSkillData(true);toast('通道 Skill 已刷新','已读取最新路由和会话状态。','success');},
    'new-channel-skill-route':()=>channelSkillRouteForm(),
    'edit-channel-skill-route':()=>channelSkillRouteForm(state.channelSkillRoutes.find(item=>item.id===el.dataset.id)),
    'save-channel-skill-route':()=>void saveChannelSkillRoute(el.dataset.id),
    'toggle-channel-skill-route':()=>void toggleChannelSkillRoute(el.dataset.id,el.dataset.enabled==='true'),
    'delete-channel-skill-route':()=>void deleteChannelSkillRoute(el.dataset.id),
    'refresh-openclaw-system':async()=>{await Promise.all([loadOpenClawStatus(true),loadOpenClawAgents(true),loadSystemHealth(true)]);toast('系统状态已刷新',state.systemHealth?.status==='ok'?'Go 服务、SQLite 与 OpenClaw 状态正常。':'部分组件需要检查。',state.systemHealth?.status==='ok'?'success':'warning');},
	    'open-channel-binding':()=>openChannelBinding(el.dataset.channel),
	    'channel-status':()=>document.getElementById('channelStatusBox')?void refreshChannelStatus(el.dataset.channel):openChannelBinding(el.dataset.channel),
	    'save-channel-binding':()=>void saveChannelBinding(el.dataset.channel),
	    'save-channel-account':()=>void saveChannelAccount(el.dataset.channel),
	    'login-channel':()=>void loginChannel(el.dataset.channel),
	    'install-channel':()=>installChannel(el.dataset.channel),
    'confirm-channel-install':()=>{const spec=formText('channelPackageSpec'); const info=(state.openClawChannels||[]).find(item=>item.id===el.dataset.channel); if(!spec && !info?.installSpec){toast('缺少插件包名','请填写 OpenClaw 插件包名或 npm spec。','warning');return;} void executeChannelInstall(el.dataset.channel,spec,el);},
	    'uninstall-channel':()=>void uninstallChannel(el.dataset.channel),
    'new-schedule':()=>schedulerForm(),'edit-schedule':()=>schedulerForm(el.dataset.id),'save-schedule':()=>void saveSchedule(),'delete-schedule':()=>void deleteSchedule(el.dataset.id),'toggle-schedule':()=>void toggleSchedule(el.dataset.id,el.dataset.enabled==='true'),'run-schedule':()=>void runSchedule(el.dataset.id),'schedule-runs':()=>void showScheduleRuns(el.dataset.id),'refresh-openclaw-jobs':()=>void loadOpenClawJobs(true),'choose-backup':()=>toast('备份目录待部署确认','浏览器不能直接授权后端写入任意外置路径；需确定盒子挂载点和目录白名单。','warning'),
    'offline-upgrade':offlineUpgradeModal,'choose-upgrade-package':()=>document.getElementById('upgradeFileInput')?.click(),'import-upgrade-package':()=>void importOfflineUpgrade(),
    'upgrade-history':()=>void showUpgradeHistory(),
    'relation-select':()=>{const input=document.getElementById(el.dataset.target);if(input){input.value=el.dataset.value;const options=document.getElementById(`${el.dataset.target}Options`);if(options){options.innerHTML='';options.hidden=true;}if(el.dataset.target==='orderQuote')syncOrderFromQuote(el.dataset.value);if(el.dataset.target==='quoteCustomer')syncQuoteTemplateCustomerFields(el.dataset.value);}},
    'save-customer':()=>void saveCustomer(),'save-lead':()=>void saveLead(),'save-quote':()=>void saveQuote(),'save-order':()=>void saveOrder(),'save-document':()=>void saveDocument(),'save-supplier':()=>void saveSupplier(),
    'add-quote-line':()=>{const product=products.find(item=>item.status==='Active')||products[0];state.quoteDraftLines.push({productId:product.id,quantity:1,unitPrice:moneyNumber(product.price),discount:0});renderQuoteDraftLines();},
    'remove-quote-line':()=>{if(state.quoteDraftLines.length===1){toast('至少保留一条产品明细','正式报价单需要至少一个产品。','warning');return;}state.quoteDraftLines.splice(Number(el.dataset.index),1);renderQuoteDraftLines();},
    'add-order-line':()=>{const product=products.find(item=>item.status==='Active')||products[0];state.orderDraftLines.push({productId:product.id,quantity:1,unitPrice:moneyNumber(product.price)});renderOrderDraftLines();},
    'remove-order-line':()=>{if(state.orderDraftLines.length===1){toast('至少保留一条产品明细','订单需要至少一个产品。','warning');return;}state.orderDraftLines.splice(Number(el.dataset.index),1);renderOrderDraftLines();},
    'template-default':()=>void templateAction('default',el.dataset.id),'template-edit':()=>void templateAction('edit',el.dataset.id),'template-preview':()=>void templateAction('preview',el.dataset.id),'template-delete':()=>void templateAction('delete',el.dataset.id),'save-template':()=>void saveTemplate(),
    'template-page-prev':()=>{state.templatePage=Math.max(1,state.templatePage-1);void templateCenter(state.templateKind);},
    'template-page-next':()=>{state.templatePage+=1;void templateCenter(state.templateKind);},
    'template-page-current':()=>toast('当前页',`第 ${state.templatePage} 页。`),
    'quote-date-filter':()=>dateFilterForm('quote'),'order-date-filter':()=>dateFilterForm('order'),'apply-date-filter':()=>applyDateFilter(el.dataset.module,el.dataset.clear==='true'),
    'send-chat':()=>void sendAgentMessage(Number(el.dataset.agent)),
    'stop-chat':()=>stopAgentMessage(Number(el.dataset.agent)),
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
    state.newsFetchLimit = Number.isInteger(savedNews.fetchLimit) && savedNews.fetchLimit >= 1 && savedNews.fetchLimit <= 100 ? savedNews.fetchLimit : state.newsFetchLimit;
    state.newsShowLimit = Number.isInteger(savedNews.showLimit) && savedNews.showLimit >= 1 && savedNews.showLimit <= 100 ? savedNews.showLimit : state.newsShowLimit;
    state.newsFrequency = ['1小时','2小时','3小时','6小时','8小时','12小时','24小时'].includes(savedNews.frequency) ? savedNews.frequency : state.newsFrequency;
    state.newsSources = savedNews.sources || state.newsSources;
    state.newsMediaIDs = Array.isArray(savedNews.mediaIds) ? savedNews.mediaIds : state.newsMediaIDs;
    state.newsMediaCategories = Array.isArray(savedNews.mediaCategories) ? savedNews.mediaCategories : state.newsMediaCategories;
    state.newsCustomSources = savedNews.customSources ?? state.newsCustomSources;
  }
  state.agentInternetAllowlists = JSON.parse(localStorage.getItem('sta100-agent-allowlists') || '{}') || {};
} catch {
  localStorage.removeItem('sta100-news-settings');
  localStorage.removeItem('sta100-agent-allowlists');
}
applyTranslations();
applyEmojiRenderingMode();
void initAuth();
