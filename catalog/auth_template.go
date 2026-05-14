package main

const authHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Вход в Shop Go</title>
<style>
:root{--blue:#005bff;--blue-dark:#003fb8;--teal:#00a38f;--pink:#f91155;--ink:#07192f;--muted:#667085;--line:#dfe7f1;--soft:#f4f7fc}
*{box-sizing:border-box}
@keyframes authRise{from{opacity:0;transform:translateY(18px) scale(.985)}to{opacity:1;transform:translateY(0) scale(1)}}
@keyframes sweep{from{transform:translateX(-120%)}to{transform:translateX(120%)}}
@keyframes floatCard{0%,100%{transform:translateY(0)}50%{transform:translateY(-6px)}}
body{margin:0;min-height:100vh;font-family:Inter,Segoe UI,Arial,sans-serif;background:linear-gradient(145deg,#eef5ff 0%,#f8fbff 44%,#eefaf7 100%);display:grid;place-items:center;padding:24px;color:var(--ink);overflow-x:hidden}
body:before{content:"";position:fixed;inset:0;background:linear-gradient(115deg,rgba(0,91,255,.12),transparent 32%,rgba(0,163,143,.1) 68%,transparent);pointer-events:none}
.wrap{position:relative;width:min(1030px,100%);display:grid;grid-template-columns:minmax(0,1.1fr) 390px;gap:20px;align-items:stretch}
.promo{position:relative;border-radius:28px;padding:34px;min-height:540px;overflow:hidden;background:linear-gradient(135deg,#091b34 0%,#0d4ed9 58%,#00a38f 130%);color:#fff;display:flex;flex-direction:column;justify-content:space-between;box-shadow:0 28px 70px rgba(13,55,130,.24);animation:authRise .36s ease both}
.promo:before{content:"";position:absolute;inset:0;background:linear-gradient(90deg,transparent,rgba(255,255,255,.18),transparent);animation:sweep 5.5s ease-in-out infinite}
.promo>*{position:relative}
.promo h1{margin:0;font-size:46px;line-height:.98;max-width:540px}
.promo p{font-size:17px;line-height:1.55;opacity:.9;max-width:560px;margin:18px 0 0}
.promo-showcase{display:grid;grid-template-columns:1fr;gap:14px;margin-top:34px;margin-bottom:22px;max-width:340px}
.showcase-card{min-height:118px;border:1px solid rgba(255,255,255,.2);border-radius:18px;background:rgba(255,255,255,.12);padding:18px;backdrop-filter:blur(8px);transition:transform .18s ease,background .18s ease}
.showcase-card:hover{transform:translateY(-3px);background:rgba(255,255,255,.17)}
.showcase-card strong{display:block;font-size:18px;margin-bottom:8px}
.showcase-card span{display:block;font-size:13px;line-height:1.45;opacity:.86}
.delivery-card{border:1px solid rgba(255,255,255,.2);border-radius:22px;padding:20px;background:rgba(255,255,255,.12);backdrop-filter:blur(10px);animation:floatCard 5s ease-in-out infinite}
.delivery-title{font-weight:900;font-size:20px;margin-bottom:16px}
.delivery-options{display:grid;grid-template-columns:1fr 1fr;gap:14px}
.delivery-option{background:rgba(255,255,255,.12);border-radius:14px;padding:16px;transition:transform .16s ease,background .16s ease}
.delivery-option:hover{transform:translateY(-2px);background:rgba(255,255,255,.2)}
.delivery-option b{display:block;margin-bottom:5px}
.delivery-option span{font-size:13px;opacity:.88;line-height:1.35}
.card{position:relative;background:rgba(255,255,255,.94);border:1px solid rgba(255,255,255,.82);border-radius:28px;box-shadow:0 26px 70px rgba(16,24,40,.13);padding:26px;animation:authRise .42s ease both;backdrop-filter:blur(10px)}
.brand{display:flex;align-items:center;gap:11px;font-size:25px;font-weight:900;margin-bottom:22px}
.mark{width:40px;height:40px;border-radius:13px;background:linear-gradient(135deg,var(--blue),#11b8a4);color:#fff;display:grid;place-items:center;box-shadow:0 10px 22px rgba(0,91,255,.28)}
.tabs{display:grid;grid-template-columns:1fr 1fr;gap:6px;background:#edf2f8;border-radius:16px;padding:5px;margin-bottom:20px}
.tab{border:0;background:transparent;border-radius:12px;padding:11px;font-weight:900;color:var(--muted);transition:background .16s ease,color .16s ease,transform .16s ease,box-shadow .16s ease}
.tab:hover{transform:translateY(-1px)}
.tab.active{background:#fff;color:var(--ink);box-shadow:0 6px 18px rgba(16,24,40,.1)}
form{display:none}
form.active{display:block;animation:authRise .2s ease both}
.field{margin-bottom:13px}
label{display:block;font-size:13px;color:var(--muted);font-weight:800;margin-bottom:7px}
input{width:100%;border:1px solid var(--line);border-radius:14px;padding:13px 14px;outline:none;font-size:15px;background:#f8fbff;transition:border-color .16s ease,box-shadow .16s ease,background .16s ease}
input:focus{border-color:var(--blue);box-shadow:0 0 0 4px rgba(0,91,255,.11);background:#fff}
.password-wrap{position:relative}
.password-wrap input{padding-right:94px}
.toggle-password{position:absolute;right:7px;top:7px;height:34px;border:0;border-radius:10px;background:#eaf1fa;color:#53677e;padding:0 10px;font-size:13px;font-weight:900;cursor:pointer;transition:background .16s ease,color .16s ease}
.toggle-password:hover{background:#dde9f6;color:var(--ink)}
.submit{width:100%;border:0;border-radius:14px;background:linear-gradient(135deg,var(--blue),var(--blue-dark));color:#fff;padding:14px 16px;font-weight:900;margin-top:8px;cursor:pointer;transition:transform .16s ease,box-shadow .16s ease,filter .16s ease}
.submit:hover{transform:translateY(-2px);box-shadow:0 14px 28px rgba(0,91,255,.24);filter:saturate(1.08)}
.msg{min-height:22px;margin-top:12px;font-weight:700}
.msg.error{color:#b42318}.msg.ok{color:#067647}
.hint{margin-top:34px;color:var(--muted);font-size:13px;line-height:1.5}
@media(max-width:900px){.wrap{grid-template-columns:1fr}.promo{min-height:auto}.promo h1{font-size:36px}.promo-showcase{grid-template-columns:1fr;max-width:none}}
@media(max-width:560px){body{padding:14px}.promo,.card{border-radius:22px;padding:22px}.promo h1{font-size:32px}.promo-showcase{grid-template-columns:1fr;margin-bottom:18px}.delivery-options{grid-template-columns:1fr}}
@media(prefers-reduced-motion:reduce){*,*:before,*:after{animation:none!important;transition:none!important}}
</style>
</head>
<body>
<main class="wrap">
  <section class="promo">
    <div>
      <h1>Добро пожаловать в Shop Go</h1>
      <p>Войдите в аккаунт, чтобы выбрать товары, оформить заказ и получить его удобным способом.</p>
      <div class="promo-showcase">
        <div class="showcase-card"><strong>Быстрые заказы</strong><span>Корзина, избранное и история покупок остаются под рукой.</span></div>
      </div>
    </div>
    <div class="delivery-card">
      <div class="delivery-title">Доставка и получение</div>
      <div class="delivery-options">
        <div class="delivery-option"><b>Самовывоз</b><span>Выберите удобный пункт выдачи рядом с вами.</span></div>
        <div class="delivery-option"><b>Курьер</b><span>Оформите доставку до нужного адреса.</span></div>
      </div>
    </div>
  </section>
  <section class="card">
    <div class="brand"><span class="mark">S</span><span>Shop Go</span></div>
    <div class="tabs">
      <button type="button" class="tab active" id="tabLogin">Вход</button>
      <button type="button" class="tab" id="tabRegister">Регистрация</button>
    </div>
    <form id="loginForm" class="active">
      <div class="field"><label for="loginUsername">Логин</label><input id="loginUsername" required autocomplete="username"></div>
      <div class="field"><label for="loginPassword">Пароль</label><div class="password-wrap"><input id="loginPassword" type="password" required autocomplete="current-password"><button class="toggle-password" type="button" data-toggle-password="loginPassword">Показать</button></div></div>
      <button class="submit" type="submit">Войти</button>
    </form>
    <form id="registerForm">
      <div class="field"><label for="registerUsername">Логин</label><input id="registerUsername" minlength="3" maxlength="32" required autocomplete="username"></div>
      <div class="field"><label for="registerPassword">Пароль</label><div class="password-wrap"><input id="registerPassword" type="password" minlength="6" required autocomplete="new-password"><button class="toggle-password" type="button" data-toggle-password="registerPassword">Показать</button></div></div>
      <button class="submit" type="submit">Создать аккаунт</button>
    </form>
    <div id="msg" class="msg"></div>
    <div class="hint">После входа вы автоматически попадёте в магазин, сможете оформить заказ и отслеживать покупки.</div>
  </section>
</main>
<script>
const nextPath = {{.next}};
const targetPath = (nextPath === '/' || nextPath === '/shop' || nextPath === '/cart/view' || nextPath === '/favorites' || nextPath === '/orders' || nextPath.startsWith('/product/') || nextPath === '/admin') ? nextPath : '/';
const tabLogin = document.getElementById('tabLogin');
const tabRegister = document.getElementById('tabRegister');
const loginForm = document.getElementById('loginForm');
const registerForm = document.getElementById('registerForm');
const msg = document.getElementById('msg');
function setMsg(text, ok){msg.textContent=text||'';msg.className=ok?'msg ok':'msg error'}
function switchTab(mode){const login=mode==='login';tabLogin.classList.toggle('active',login);tabRegister.classList.toggle('active',!login);loginForm.classList.toggle('active',login);registerForm.classList.toggle('active',!login);setMsg('',false)}
async function postJSON(url, body){
  const res = await fetch(url,{method:'POST',headers:{'Content-Type':'application/json'},credentials:'include',body:JSON.stringify(body)});
  const payload = await res.json().catch(()=>({}));
  if(!res.ok) throw new Error(payload.error || 'Ошибка запроса');
  return payload;
}
tabLogin.addEventListener('click',()=>switchTab('login'));
tabRegister.addEventListener('click',()=>switchTab('register'));
document.querySelectorAll('[data-toggle-password]').forEach(btn=>{
  btn.addEventListener('click',()=>{
    const input=document.getElementById(btn.dataset.togglePassword);
    if(!input) return;
    const show=input.type==='password';
    input.type=show?'text':'password';
    btn.textContent=show?'Скрыть':'Показать';
  });
});
loginForm.addEventListener('submit',async e=>{
  e.preventDefault();
  try{await postJSON('/auth/login',{username:document.getElementById('loginUsername').value.trim(),password:document.getElementById('loginPassword').value});setMsg('Вход выполнен',true);location.href=targetPath}catch(err){setMsg(err.message,false)}
});
registerForm.addEventListener('submit',async e=>{
  e.preventDefault();
  const username=document.getElementById('registerUsername').value.trim();
  try{await postJSON('/auth/register',{username,password:document.getElementById('registerPassword').value});setMsg('Аккаунт создан.',true);location.href=targetPath}catch(err){setMsg(err.message,false)}
});
</script>
</body>
</html>`
