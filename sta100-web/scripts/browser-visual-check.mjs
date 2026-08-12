import fs from 'node:fs/promises';

const endpoint = process.env.CDP_ENDPOINT || 'http://127.0.0.1:9223';
const appURL = process.env.STA100_URL || 'http://127.0.0.1:8080';
const outputDir = process.env.OUTPUT_DIR || '/tmp/sta100-visual';
const targets = await fetch(`${endpoint}/json/list`).then(response => response.json());
const target = targets.find(item => item.type === 'page');
if (!target) throw new Error('Chromium page target not found');

const socket = new WebSocket(target.webSocketDebuggerUrl);
let sequence = 0;
const pending = new Map();
socket.addEventListener('message', event => {
  const message = JSON.parse(event.data);
  if (!message.id || !pending.has(message.id)) return;
  const { resolve, reject } = pending.get(message.id);
  pending.delete(message.id);
  if (message.error) reject(new Error(message.error.message));
  else resolve(message.result);
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
async function screenshot(name) {
  const result = await command('Page.captureScreenshot', { format: 'png', captureBeyondViewport: false });
  await fs.writeFile(`${outputDir}/${name}.png`, Buffer.from(result.data, 'base64'));
}

await fs.mkdir(outputDir, { recursive: true });
await command('Page.enable');
await command('Runtime.enable');
await command('Emulation.setDeviceMetricsOverride', { width: 1440, height: 1000, deviceScaleFactor: 1, mobile: false });
await command('Page.navigate', { url: `${appURL}/#overview` });
await waitFor("document.getElementById('loginPassword') || !document.getElementById('appShell').hidden");
if (await evaluate("Boolean(document.getElementById('loginPassword'))")) {
  await evaluate("document.getElementById('loginPassword').value='admin'; document.getElementById('loginForm').requestSubmit(); true");
}
await waitFor("!document.getElementById('appShell').hidden && document.querySelector('.brand-icon-frame')");
await screenshot('overview');

await evaluate("location.hash='customers'");
await waitFor("document.querySelector('[data-action=\"customer-detail\"]')");
await evaluate("document.querySelector('[data-action=\"customer-detail\"]').click(); true");
await waitFor("!document.getElementById('drawerBackdrop').hidden");
await screenshot('customer-detail');
const detail = await evaluate(`(() => {
  const frame = document.querySelector('.brand-icon-frame').getBoundingClientRect();
  const image = document.querySelector('.brand-app-icon').getBoundingClientRect();
  const user = document.querySelector('.topbar-user').getBoundingClientRect();
  const disclaimer = document.querySelector('.ai-disclaimer').getBoundingClientRect();
  const backdrop = getComputedStyle(document.getElementById('drawerBackdrop'));
  return {
    brandFrame: { x: frame.x, y: frame.y, width: frame.width, height: frame.height },
    brandImage: { x: image.x, y: image.y, width: image.width, height: image.height },
    user: { x: user.x, y: user.y, width: user.width, height: user.height },
    disclaimer: { x: disclaimer.x, y: disclaimer.y, width: disclaimer.width, height: disclaimer.height },
    drawerBackdropFilter: backdrop.backdropFilter,
    drawerBackground: backdrop.backgroundColor,
    bodyOverflow: document.body.scrollWidth > document.documentElement.clientWidth + 1
  };
})()`);

await evaluate("document.querySelector('[data-customer-tab=\"activity\"]').click(); true");
await waitFor("document.querySelector('[data-action=\"new-customer-communication\"]')");
await evaluate("document.querySelector('[data-action=\"new-customer-communication\"]').click(); true");
await waitFor("!document.getElementById('modalBackdrop').hidden && document.getElementById('communicationContent')");
await screenshot('communication-modal');
const modal = await evaluate(`(() => ({
  drawerHidden: document.getElementById('drawerBackdrop').hidden,
  modalVisible: !document.getElementById('modalBackdrop').hidden,
  modalBackdropFilter: getComputedStyle(document.getElementById('modalBackdrop')).backdropFilter,
  contentVisible: document.getElementById('communicationContent').getBoundingClientRect().height > 0
}))()`);

socket.close();
console.log(JSON.stringify({ outputDir, detail, modal }, null, 2));
