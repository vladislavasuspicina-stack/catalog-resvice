package main

const authHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Вход в магазин</title>
<style>
*{box-sizing:border-box}
body{margin:0;min-height:100vh;font-family:Segoe UI,Tahoma,Geneva,Verdana,sans-serif;background:linear-gradient(120deg,#1f2937,#111827);display:flex;align-items:center;justify-content:center;padding:20px;color:#111}
.card{width:420px;max-width:100%;background:#fff;border-radius:14px;box-shadow:0 18px 44px rgba(0,0,0,.35);padding:24px}
h1{margin:0 0 8px;font-size:28px}
.subtitle{margin:0 0 18px;color:#6b7280}
.tabs{display:flex;gap:8px;margin-bottom:16px}
.tab{flex:1;border:1px solid #d1d5db;background:#f3f4f6;border-radius:8px;padding:9px 10px;cursor:pointer;font-weight:600}
.tab.active{background:#111827;color:#fff;border-color:#111827}
form{display:none}
form.active{display:block}
label{display:block;font-size:13px;color:#374151;margin:10px 0 6px}
input{width:100%;padding:11px 12px;border:1px solid #d1d5db;border-radius:8px;font-size:14px}
input:focus{outline:none;border-color:#111827}
.submit{width:100%;margin-top:14px;padding:11px 12px;border:none;border-radius:8px;background:#111827;color:#fff;font-weight:600;cursor:pointer}
.hint{margin-top:10px;font-size:13px;color:#6b7280}
.msg{margin-top:10px;font-size:14px;min-height:18px}
.msg.error{color:#b91c1c}
.msg.ok{color:#166534}
</style>
</head>
<body>
<div class="card">
  <h1>Добро пожаловать</h1>
  <p class="subtitle">Войдите или зарегистрируйтесь, чтобы открыть магазин.</p>

  <div class="tabs">
    <button type="button" class="tab active" id="tabLogin">Вход</button>
    <button type="button" class="tab" id="tabRegister">Регистрация</button>
  </div>

  <form id="loginForm" class="active">
    <label for="loginUsername">Логин</label>
    <input id="loginUsername" type="text" required autocomplete="username">
    <label for="loginPassword">Пароль</label>
    <input id="loginPassword" type="password" required autocomplete="current-password">
    <button class="submit" type="submit">Войти в магазин</button>
  </form>

  <form id="registerForm">
    <label for="registerUsername">Логин</label>
    <input id="registerUsername" type="text" minlength="3" maxlength="32" required autocomplete="username">
    <label for="registerPassword">Пароль (минимум 6 символов)</label>
    <input id="registerPassword" type="password" minlength="6" required autocomplete="new-password">
    <button class="submit" type="submit">Зарегистрироваться</button>
  </form>

  <div id="msg" class="msg"></div>
  <div class="hint">После входа вы будете автоматически перенаправлены в магазин.</div>
</div>

<script>
const nextPath = {{printf "%q" .next}};
const targetPath = (nextPath === '/' || nextPath === '/shop' || nextPath === '/cart/view' || nextPath.startsWith('/product/')) ? nextPath : '/';
const tabLogin = document.getElementById('tabLogin');
const tabRegister = document.getElementById('tabRegister');
const loginForm = document.getElementById('loginForm');
const registerForm = document.getElementById('registerForm');
const msg = document.getElementById('msg');

function setMsg(text, ok){
  msg.textContent = text || '';
  msg.className = ok ? 'msg ok' : 'msg error';
}

function switchTab(mode){
  const login = mode === 'login';
  tabLogin.classList.toggle('active', login);
  tabRegister.classList.toggle('active', !login);
  loginForm.classList.toggle('active', login);
  registerForm.classList.toggle('active', !login);
  setMsg('', false);
}

async function postJSON(url, body){
  const res = await fetch(url, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    credentials: 'include',
    body: JSON.stringify(body)
  });
  const payload = await res.json().catch(() => ({}));
  if(!res.ok){
    throw new Error(payload.error || 'Ошибка запроса');
  }
  return payload;
}

tabLogin.addEventListener('click', ()=>switchTab('login'));
tabRegister.addEventListener('click', ()=>switchTab('register'));

loginForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const username = document.getElementById('loginUsername').value.trim();
  const password = document.getElementById('loginPassword').value;
  try{
    await postJSON('/auth/login', {username, password});
    setMsg('Успешный вход.', true);
    location.href = targetPath;
  }catch(err){
    setMsg(err.message || 'Ошибка входа', false);
  }
});

registerForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  const username = document.getElementById('registerUsername').value.trim();
  const password = document.getElementById('registerPassword').value;
  try{
    await postJSON('/auth/register', {username, password});
    alert('Вы успешно зарегистрировались. Теперь войдите под своими данными.');
    setMsg('Регистрация завершена. Выполните вход.', true);
    switchTab('login');
    document.getElementById('loginUsername').value = username;
    document.getElementById('loginPassword').value = '';
    document.getElementById('loginPassword').focus();
  }catch(err){
    setMsg(err.message || 'Ошибка регистрации', false);
  }
});
</script>
</body>
</html>`
