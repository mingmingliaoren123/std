import fs from 'node:fs/promises';

const endpoint = process.env.CDP_ENDPOINT || 'http://127.0.0.1:9224';
const appURL = process.env.STA100_URL || 'http://127.0.0.1:8080';
const outputDir = process.env.OUTPUT_DIR || '/tmp/sta100-token-visual';
const targets = await fetch(`${endpoint}/json/list`).then(response => response.json());
const target = targets.find(item => item.type === 'page');
if (!target) throw new Error('Chromium page target not found');

const socket = new WebSocket(target.webSocketDebuggerUrl);
let sequence = 0;
const pending = new Map();
const errors = [];
socket.addEventListener('message', event => {
  const message = JSON.parse(event.data);
  if (message.id && pending.has(message.id)) {
    const entry = pending.get(message.id);
    pending.delete(message.id);
    message.error ? entry.reject(new Error(message.error.message)) : entry.resolve(message.result);
  } else if (message.method === 'Runtime.exceptionThrown') errors.push(message.params.exceptionDetails.text);
});
await new Promise((resolve, reject) => {
  socket.addEventListener('open', resolve, { once: true });
  socket.addEventListener('error', reject, { once: true });
});
const command = (method, params = {}) => new Promise((resolve, reject) => {
  const id = ++sequence;
  pending.set(id, { resolve, reject });
  socket.send(JSON.stringify({ id, method, params }));
});
async function evaluate(expression) {
  const result = await command('Runtime.evaluate', { expression, awaitPromise: true, returnByValue: true });
  if (result.exceptionDetails) throw new Error(JSON.stringify(result.exceptionDetails));
  return result.result.value;
}
async function waitFor(expression, timeout = 20000) {
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
await waitFor("document.getElementById('loginPassword') || (document.getElementById('appShell') && !document.getElementById('appShell').hidden)");
if (await evaluate("Boolean(document.getElementById('loginPassword'))")) {
  await evaluate("document.getElementById('loginPassword').value='admin'; document.getElementById('loginForm').requestSubmit(); true");
}
await waitFor("!document.getElementById('appShell').hidden && document.querySelector('.topbar-brand-logo')");
await evaluate("document.querySelector('[data-action=token-usage]').click(); true");
await waitFor("!document.getElementById('modalBackdrop').hidden && document.querySelector('.token-metrics')");
await screenshot('desktop-token');
const desktop = await evaluate(`(() => {
  const logo = document.querySelector('.topbar-brand-logo').getBoundingClientRect();
  const token = document.querySelector('.token-usage-button').getBoundingClientRect();
  const user = document.querySelector('.topbar-user').getBoundingClientRect();
  return {logo:{x:logo.x,y:logo.y,width:logo.width,height:logo.height},token:{x:token.x,width:token.width},user:{x:user.x,width:user.width},overflow:document.body.scrollWidth>document.documentElement.clientWidth+1,modalText:document.getElementById('modalBody').innerText};
})()`);

await evaluate("document.querySelector('[data-action=close-modal]').click(); true");
await command('Emulation.setDeviceMetricsOverride', { width: 390, height: 844, deviceScaleFactor: 1, mobile: true });
await new Promise(resolve => setTimeout(resolve, 300));
await screenshot('mobile-header');
const mobile = await evaluate(`(() => {
  const logo = getComputedStyle(document.querySelector('.topbar-brand-logo'));
  const token = document.querySelector('.token-usage-button').getBoundingClientRect();
  return {logoDisplay:logo.display,token:{x:token.x,width:token.width},viewport:document.documentElement.clientWidth,bodyWidth:document.body.scrollWidth,overflow:document.body.scrollWidth>document.documentElement.clientWidth+1};
})()`);

socket.close();
console.log(JSON.stringify({ outputDir, desktop, mobile, errors }, null, 2));
if (desktop.overflow || mobile.overflow || errors.length) process.exitCode = 1;
