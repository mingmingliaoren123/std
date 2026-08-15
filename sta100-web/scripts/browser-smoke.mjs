import process from 'node:process';

const endpoint = process.env.CDP_ENDPOINT || 'http://127.0.0.1:9223';
const appURL = process.env.STA100_URL || 'http://127.0.0.1:18080';
const password = process.env.STA100_PASSWORD || 'admin';
const viewportWidth = Number(process.env.VIEWPORT_WIDTH || 1440);
const viewportHeight = Number(process.env.VIEWPORT_HEIGHT || 1000);
const pages = ['overview', 'agents', 'customers', 'quotes', 'orders', 'documents', 'products', 'suppliers', 'database', 'news', 'settings'];
const settingsTabs = ['model', 'channels', 'scheduler', 'backup', 'security', 'system', 'upgrade'];

const targets = await fetch(`${endpoint}/json/list`).then(response => response.json());
const target = targets.find(item => item.type === 'page');
if (!target) throw new Error('Chromium page target not found');

const socket = new WebSocket(target.webSocketDebuggerUrl);
let sequence = 0;
const pending = new Map();
const consoleErrors = [];

socket.addEventListener('message', event => {
  const message = JSON.parse(event.data);
  if (message.id && pending.has(message.id)) {
    const { resolve, reject } = pending.get(message.id);
    pending.delete(message.id);
    if (message.error) reject(new Error(message.error.message));
    else resolve(message.result);
    return;
  }
  if (message.method === 'Runtime.exceptionThrown') consoleErrors.push(message.params.exceptionDetails.text);
  if (message.method === 'Runtime.consoleAPICalled' && message.params.type === 'error') {
    consoleErrors.push(message.params.args.map(item => item.value || item.description || '').join(' '));
  }
});

await new Promise((resolve, reject) => {
  socket.addEventListener('open', resolve, { once: true });
  socket.addEventListener('error', reject, { once: true });
});

function command(method, params = {}) {
  const id = ++sequence;
  socket.send(JSON.stringify({ id, method, params }));
  return new Promise((resolve, reject) => pending.set(id, { resolve, reject }));
}

async function evaluate(expression) {
  const result = await command('Runtime.evaluate', { expression, awaitPromise: true, returnByValue: true });
  if (result.exceptionDetails) throw new Error(JSON.stringify(result.exceptionDetails));
  return result.result.value;
}

async function waitFor(expression, timeout = 15000) {
  const started = Date.now();
  while (Date.now() - started < timeout) {
    if (await evaluate(`Boolean(${expression})`)) return;
    await new Promise(resolve => setTimeout(resolve, 100));
  }
  throw new Error(`Timed out waiting for ${expression}`);
}

await command('Page.enable');
await command('Runtime.enable');
await command('Emulation.setDeviceMetricsOverride', { width: viewportWidth, height: viewportHeight, deviceScaleFactor: 1, mobile: viewportWidth < 600 });
await command('Page.navigate', { url: appURL });
await waitFor("document.getElementById('loginPassword') || (document.getElementById('appShell') && !document.getElementById('appShell').hidden)");
if (await evaluate("Boolean(document.getElementById('loginPassword'))")) {
  await evaluate(`(() => { const input=document.getElementById('loginPassword'); input.value=${JSON.stringify(password)}; document.getElementById('loginForm').requestSubmit(); return true; })()`);
}
await waitFor("!document.getElementById('appShell').hidden && document.querySelector('.metric-grid')", 20000);

const results = [];
for (const page of pages) {
  await evaluate(`location.hash=${JSON.stringify(page)}`);
  await waitFor(`location.hash === ${JSON.stringify(`#${page}`)} && document.getElementById('pageTitle').textContent.trim().length > 0`);
  await new Promise(resolve => setTimeout(resolve, 150));
  results.push(await evaluate(`(() => ({
    page: ${JSON.stringify(page)},
    title: document.getElementById('pageTitle').textContent.trim(),
    textLength: document.getElementById('pageRoot').innerText.trim().length,
    viewport: document.documentElement.clientWidth,
    bodyWidth: document.body.scrollWidth,
    overflow: document.body.scrollWidth > document.documentElement.clientWidth + 1
  }))()`));
}

const settingsResults = [];
for (const tab of settingsTabs) {
  await evaluate(`document.querySelector('[data-settings-tab=${JSON.stringify(tab)}]').click()`);
  await waitFor(`document.querySelector('[data-settings-tab=${JSON.stringify(tab)}]').classList.contains('active')`);
  await new Promise(resolve => setTimeout(resolve, 100));
  settingsResults.push(await evaluate(`(() => ({
    tab: ${JSON.stringify(tab)},
    textLength: document.querySelector('.settings-content').innerText.trim().length,
    overflow: document.body.scrollWidth > document.documentElement.clientWidth + 1
  }))()`));
}

socket.close();
const failed = [...results.filter(item => item.textLength < 20 || item.overflow), ...settingsResults.filter(item => item.textLength < 20 || item.overflow)];
console.log(JSON.stringify({ viewport: { width: viewportWidth, height: viewportHeight }, pages: results, settings: settingsResults, consoleErrors, failed }, null, 2));
if (failed.length || consoleErrors.length) process.exitCode = 1;
