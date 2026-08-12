const endpoint = process.env.CDP_ENDPOINT || 'http://127.0.0.1:9225';
const appURL = process.env.STA100_URL || 'http://127.0.0.1:8080';
const width = Number(process.env.VIEWPORT_WIDTH || 1440);
const height = Number(process.env.VIEWPORT_HEIGHT || 1000);
const pages = ['overview','agents','customers','quotes','orders','documents','products','suppliers','database','news','settings'];
const settingsTabs = ['model','plugins','scheduler','backup','security','system','upgrade'];
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
    const item = pending.get(message.id); pending.delete(message.id);
    message.error ? item.reject(new Error(message.error.message)) : item.resolve(message.result);
  } else if (message.method === 'Runtime.exceptionThrown') errors.push(message.params.exceptionDetails.text);
});
await new Promise((resolve,reject) => { socket.addEventListener('open',resolve,{once:true}); socket.addEventListener('error',reject,{once:true}); });
const command = (method, params={}) => new Promise((resolve,reject) => { const id=++sequence; pending.set(id,{resolve,reject}); socket.send(JSON.stringify({id,method,params})); });
async function evaluate(expression) {
  const result=await command('Runtime.evaluate',{expression,awaitPromise:true,returnByValue:true});
  if(result.exceptionDetails) throw new Error(JSON.stringify(result.exceptionDetails));
  return result.result.value;
}
async function waitFor(expression, timeout=20000) {
  const started=Date.now();
  while(Date.now()-started<timeout){if(await evaluate(`Boolean(${expression})`))return;await new Promise(r=>setTimeout(r,100));}
  throw new Error(`Timed out waiting for ${expression}`);
}
async function audit(label, rootExpression="document.getElementById('pageRoot')") {
  return evaluate(`(() => {
    const root=${rootExpression};
    const elements=[...root.querySelectorAll('*')].filter(el=>{
      const style=getComputedStyle(el), rect=el.getBoundingClientRect();
      return rect.width>0&&rect.height>0&&style.visibility!=='hidden'&&style.display!=='none'&&[...el.childNodes].some(n=>n.nodeType===Node.TEXT_NODE&&n.textContent.trim());
    });
    const rows=elements.map(el=>{const s=getComputedStyle(el),r=el.getBoundingClientRect();return {tag:el.tagName.toLowerCase(),cls:String(el.className||'').slice(0,80),text:el.innerText.trim().replace(/\\s+/g,' ').slice(0,90),font:Number.parseFloat(s.fontSize),line:Number.parseFloat(s.lineHeight)||0,color:s.color,width:Math.round(r.width),height:Math.round(r.height),horizontal:el.scrollWidth>el.clientWidth+1,vertical:el.scrollHeight>el.clientHeight+1,ellipsis:s.textOverflow==='ellipsis',clamped:Number.parseInt(s.webkitLineClamp)>0,whiteSpace:s.whiteSpace};});
    const counts={};rows.forEach(row=>counts[row.font]=(counts[row.font]||0)+1);
    return {label:${JSON.stringify(label)},count:rows.length,min:rows.length?Math.min(...rows.map(r=>r.font)):0,fontCounts:counts,under11:rows.filter(r=>r.font<11).slice(0,30),under12:rows.filter(r=>r.font<12).slice(0,30),overflow:rows.filter(r=>(r.horizontal||r.vertical)&&!r.ellipsis&&!r.clamped).slice(0,30),ellipsis:rows.filter(r=>(r.ellipsis&&r.horizontal)||r.clamped).slice(0,20),pageOverflow:document.body.scrollWidth>document.documentElement.clientWidth+1};
  })()`);
}

await command('Page.enable'); await command('Runtime.enable'); await command('Network.enable'); await command('Network.setCacheDisabled',{cacheDisabled:true});
await command('Emulation.setDeviceMetricsOverride',{width,height,deviceScaleFactor:1,mobile:width<600});
await command('Page.navigate',{url:`${appURL}/#overview`});
await waitFor("document.getElementById('loginPassword') || (document.getElementById('appShell') && !document.getElementById('appShell').hidden)");
if(await evaluate("Boolean(document.getElementById('loginPassword'))"))await evaluate("document.getElementById('loginPassword').value='admin';document.getElementById('loginForm').requestSubmit();true");
await waitFor("!document.getElementById('appShell').hidden && document.querySelector('.metric-grid')");
const results=[];
results.push(await audit('chrome:topbar',"document.querySelector('.topbar')"));
results.push(await audit('chrome:sidebar',"document.querySelector('.sidebar')"));
for(const page of pages){await evaluate(`location.hash=${JSON.stringify(page)}`);await waitFor(`location.hash===${JSON.stringify('#'+page)}`);await new Promise(r=>setTimeout(r,120));results.push(await audit(page));}
for(const tab of settingsTabs){await evaluate(`document.querySelector('[data-settings-tab=${JSON.stringify(tab)}]').click()`);await waitFor(`document.querySelector('[data-settings-tab=${JSON.stringify(tab)}]').classList.contains('active')`);results.push(await audit(`settings:${tab}`,"document.querySelector('.settings-content')"));}

await evaluate("document.querySelector('[data-action=token-usage]').click();true");await waitFor("!document.getElementById('modalBackdrop').hidden");results.push(await audit('modal:token',"document.getElementById('modal')"));await evaluate("document.querySelector('[data-action=close-modal]').click();true");
await evaluate("location.hash='agents'");await waitFor("document.querySelector('[data-action=agent-chat]')");await evaluate("document.querySelector('[data-action=agent-chat]').click();true");await waitFor("!document.getElementById('modalBackdrop').hidden");results.push(await audit('modal:agent-chat',"document.getElementById('modal')"));await evaluate("document.querySelector('[data-action=close-modal]').click();true");
await evaluate("location.hash='customers'");await waitFor("document.querySelector('[data-action=customer-detail]')");await evaluate("document.querySelector('[data-action=customer-detail]').click();true");await waitFor("!document.getElementById('drawerBackdrop').hidden");results.push(await audit('drawer:customer',"document.getElementById('drawer')"));

socket.close();
const failures=results.filter(r=>r.pageOverflow||r.under11.length||r.overflow.length);
console.log(JSON.stringify({viewport:{width,height},results,errors,failures:failures.map(r=>({label:r.label,under11:r.under11,overflow:r.overflow,pageOverflow:r.pageOverflow}))},null,2));
if(errors.length||results.some(r=>r.pageOverflow))process.exitCode=1;
