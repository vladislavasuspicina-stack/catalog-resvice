package main

const adminHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Админ-панель каталога</title>
<style>
@import url('https://fonts.googleapis.com/css2?family=Manrope:wght@400;500;700;800&family=Space+Grotesk:wght@500;700&display=swap');
:root{
  --bg-a:#eef4ff;
  --bg-b:#ecfff5;
  --ink:#152432;
  --muted:#5d7288;
  --card:#ffffff;
  --line:#d7e3ee;
  --brand:#0f7a70;
  --brand-2:#0aa792;
  --danger:#cc3e4c;
  --shadow:0 20px 45px rgba(17,38,58,0.12);
}
*{box-sizing:border-box}
body{
  margin:0;
  min-height:100vh;
  padding:24px;
  color:var(--ink);
  font-family:'Manrope','Segoe UI',sans-serif;
  background:
    radial-gradient(900px 500px at -10% -10%, rgba(14,122,112,0.12), transparent 65%),
    radial-gradient(700px 420px at 105% 0%, rgba(15,71,173,0.1), transparent 60%),
    linear-gradient(145deg, var(--bg-a) 0%, var(--bg-b) 100%);
}
.admin-shell{max-width:1240px;margin:0 auto}
.topbar{
  display:flex;
  justify-content:space-between;
  align-items:center;
  gap:14px;
  padding:18px 22px;
  border-radius:22px;
  border:1px solid rgba(255,255,255,0.6);
  background:linear-gradient(120deg, rgba(255,255,255,0.92), rgba(255,255,255,0.72));
  box-shadow:var(--shadow);
  backdrop-filter:blur(6px);
}
.title-wrap h1{
  margin:0;
  font-family:'Space Grotesk','Manrope',sans-serif;
  font-size:clamp(24px,3.2vw,34px);
  line-height:1.05;
  letter-spacing:0.02em;
}
.title-wrap p{margin:6px 0 0;color:var(--muted);font-size:14px}
.layout{
  margin-top:18px;
  display:grid;
  grid-template-columns:minmax(320px,1fr) minmax(320px,1fr);
  gap:16px;
}
.panel{
  background:var(--card);
  border:1px solid var(--line);
  border-radius:20px;
  padding:18px;
  box-shadow:0 10px 24px rgba(16,52,82,0.07);
  animation:rise .35s ease both;
}
.panel h2{margin:0 0 12px;font-size:19px}
.panel-list{grid-column:1 / -1}
.row{display:grid;gap:10px}
.row.two{grid-template-columns:1fr 1fr}
label{display:block;font-size:13px;font-weight:700;color:#395167;margin-bottom:6px}
input,select,textarea{
  width:100%;
  padding:10px 11px;
  border:1px solid var(--line);
  border-radius:12px;
  outline:none;
  background:#fff;
  color:#0f2132;
  font:500 14px/1.35 'Manrope','Segoe UI',sans-serif;
  transition:border-color .18s ease, box-shadow .18s ease;
}
textarea{min-height:96px;resize:vertical}
input:focus,select:focus,textarea:focus{
  border-color:#7bb8f4;
  box-shadow:0 0 0 4px rgba(97,174,255,0.2);
}
.btn{
  border:none;
  border-radius:12px;
  font:700 14px/1 'Space Grotesk','Manrope',sans-serif;
  letter-spacing:0.01em;
  padding:11px 14px;
  cursor:pointer;
  transition:transform .16s ease, box-shadow .16s ease, background .16s ease;
}
.btn:hover{transform:translateY(-1px)}
.btn:active{transform:translateY(0)}
.btn-primary{
  color:#fff;
  background:linear-gradient(130deg, var(--brand) 0%, var(--brand-2) 100%);
  box-shadow:0 9px 18px rgba(15,122,112,0.24);
}
.btn-light{background:#eaf2fb;color:#325676}
.btn-danger{background:#ffe7e7;color:var(--danger)}
.btn-ghost{background:#f1f6fb;color:#39546d}
.btn-wide{width:100%;margin-top:12px}
.inline{
  display:flex;
  gap:10px;
  align-items:center;
}
.inline input{flex:1}
.upload-wrap{
  border:1px dashed #9ec6e8;
  border-radius:14px;
  background:#f8fcff;
  padding:10px;
}
.file-input{
  display:inline-flex;
  align-items:center;
  justify-content:center;
  gap:8px;
  min-height:38px;
  padding:0 14px;
  border-radius:10px;
  background:#e8f2ff;
  color:#234f78;
  font:700 13px/1 'Space Grotesk','Manrope',sans-serif;
  cursor:pointer;
}
.file-input input{display:none}
.preview{
  margin-top:10px;
  min-height:132px;
  border:1px dashed #bfd6ec;
  border-radius:14px;
  background:#f8fbff;
  display:flex;
  align-items:center;
  justify-content:center;
  overflow:hidden;
  color:#5f768c;
  font-size:13px;
}
.preview img{width:100%;height:132px;object-fit:contain;background:#fff}
.helper{font-size:12px;color:#647a91}
.panel-head{
  display:flex;
  justify-content:space-between;
  align-items:flex-end;
  gap:12px;
  margin-bottom:12px;
}
.stats{font-size:13px;color:#5c7288}
.product-grid{
  display:grid;
  gap:12px;
  grid-template-columns:repeat(auto-fill,minmax(240px,1fr));
}
.order-list{
  display:grid;
  gap:12px;
}
.order-card{
  border:1px solid var(--line);
  border-radius:16px;
  background:#fff;
  padding:14px;
}
.order-head{
  display:flex;
  justify-content:space-between;
  gap:12px;
  align-items:flex-start;
  margin-bottom:10px;
}
.order-head h3{margin:0;font-size:18px}
.order-meta{font-size:13px;color:var(--muted);line-height:1.45}
.status-pill{
  display:inline-flex;
  align-items:center;
  min-height:28px;
  padding:0 10px;
  border-radius:999px;
  background:#eaf7f4;
  color:var(--brand);
  font-weight:800;
  font-size:12px;
  white-space:nowrap;
}
.order-items{margin:10px 0 0;padding-left:18px;color:#395167;font-size:14px}
.order-items li{margin:4px 0}
.order-total{font-weight:800;margin-top:10px}
.order-status-row{
  display:flex;
  align-items:center;
  gap:8px;
  flex-wrap:wrap;
  margin-top:12px;
}
.order-status-row select{max-width:220px}
.product-card{
  border:1px solid var(--line);
  border-radius:16px;
  background:#fff;
  overflow:hidden;
  cursor:pointer;
  transition:transform .16s ease, box-shadow .16s ease;
}
.product-card:hover{
  transform:translateY(-2px);
  box-shadow:0 12px 24px rgba(12,38,66,0.12);
}
.card-image{
  height:130px;
  background:linear-gradient(150deg,#eaf1fb,#d8f0f1);
  display:flex;
  align-items:center;
  justify-content:center;
  color:#60788f;
  font-size:13px;
}
.card-image img{width:100%;height:100%;object-fit:contain;padding:8px;background:#fff}
.card-body{padding:12px}
.card-title{
  margin:0 0 6px;
  font-size:17px;
  line-height:1.2;
  white-space:nowrap;
  overflow:hidden;
  text-overflow:ellipsis;
}
.card-meta{font-size:13px;color:#4f667f;margin-bottom:7px}
.card-desc{
  margin:0;
  font-size:13px;
  color:#60758b;
  min-height:34px;
  max-height:34px;
  overflow:hidden;
}
.card-actions{
  margin-top:10px;
  display:flex;
  gap:8px;
  justify-content:flex-end;
}
.card-actions button{padding:8px 10px;border-radius:10px}
.empty-state{
  grid-column:1/-1;
  border:1px dashed #a8c6df;
  border-radius:14px;
  padding:20px;
  text-align:center;
  color:#5f7a92;
  background:#f8fbff;
}
.modal-backdrop{
  position:fixed;
  inset:0;
  background:rgba(11,31,48,0.46);
  align-items:center;
  justify-content:center;
  padding:14px;
  z-index:2000;
}
.modal-card{
  width:860px;
  max-width:100%;
  max-height:92vh;
  overflow:auto;
  background:#fff;
  border-radius:18px;
  border:1px solid var(--line);
  box-shadow:0 25px 55px rgba(7,19,35,0.28);
  padding:16px;
}
.modal-head{
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap:10px;
  margin-bottom:12px;
}
.close-btn{
  border:none;
  width:34px;
  height:34px;
  border-radius:10px;
  background:#edf3fb;
  color:#2e4a63;
  font-size:18px;
  cursor:pointer;
}
.modal-grid{
  display:grid;
  gap:12px;
  grid-template-columns:1fr 1fr;
}
.modal-actions{
  margin-top:12px;
  display:flex;
  justify-content:flex-end;
  gap:8px;
}
@keyframes rise{
  from{opacity:0;transform:translateY(8px)}
  to{opacity:1;transform:translateY(0)}
}
@media (max-width:980px){
  .layout{grid-template-columns:1fr}
  .panel-list{grid-column:1}
}
@media (max-width:760px){
  body{padding:14px}
  .row.two{grid-template-columns:1fr}
  .modal-grid{grid-template-columns:1fr}
  .topbar{padding:14px}
}
</style>
</head>
<body>
<div class="admin-shell">
  <header class="topbar">
    <div class="title-wrap">
      <h1>Админ-панель каталога</h1>
      <p>Управление категориями и товарами в одном месте.</p>
    </div>
    <button id="btnLogout" class="btn btn-ghost" type="button">Выйти</button>
  </header>

  <main class="layout">
    <section class="panel">
      <h2>Создать категорию</h2>
      <div class="inline">
        <input id="catName" placeholder="Название категории" autocomplete="off">
        <button id="btnCreateCat" class="btn btn-primary" type="button">Добавить</button>
      </div>
    </section>

    <section class="panel">
      <h2>Создать товар</h2>
      <div class="row two">
        <div>
          <label for="pName">Название</label>
          <input id="pName" placeholder="Название товара" autocomplete="off">
        </div>
        <div>
          <label for="pCategory">Категория</label>
          <select id="pCategory"><option value="">Выберите категорию</option></select>
        </div>
        <div>
          <label for="pBrand">Бренд</label>
          <input id="pBrand" placeholder="Samsung / Xiaomi / Apple" autocomplete="off">
        </div>
        <div>
          <label for="pPrice">Цена</label>
          <input id="pPrice" placeholder="0" inputmode="decimal">
        </div>
        <div>
          <label for="pStock">Остаток</label>
          <input id="pStock" placeholder="0" inputmode="numeric">
        </div>
        <div>
          <label for="pColor">Цвет</label>
          <input id="pColor" placeholder="Цвет">
        </div>
        <div>
          <label for="pCondition">Состояние</label>
          <input id="pCondition" placeholder="новый / б/у">
        </div>
        <div>
          <label for="pCountry">Страна</label>
          <input id="pCountry" placeholder="Страна">
        </div>
        <div>
          <label for="pMaterial">Материал</label>
          <input id="pMaterial" placeholder="Материал">
        </div>
      </div>

      <div style="margin-top:10px">
        <label for="pDesc">Описание</label>
        <textarea id="pDesc" placeholder="Краткое описание"></textarea>
      </div>

      <div class="upload-wrap" style="margin-top:10px">
        <div class="inline" style="justify-content:space-between">
          <label class="file-input" for="pImage">
            <input id="pImage" type="file" accept="image/*">
            <span>Выбрать изображение</span>
          </label>
          <button id="pImageRemove" class="btn btn-danger" type="button" style="display:none">Удалить изображение</button>
        </div>
        <div id="pImagePreview" class="preview">Изображение не выбрано</div>
        <div class="helper">Формат карточки: 900×1100. Изображение будет вписано целиком.</div>
      </div>

      <button id="btnCreateProd" class="btn btn-primary btn-wide" type="button">Создать товар</button>
    </section>

    <section class="panel panel-list">
      <div class="panel-head">
        <h2>Товары</h2>
        <div id="statsLine" class="stats">0 товаров</div>
      </div>
      <div id="prodList" class="product-grid"></div>
    </section>

    <section class="panel panel-list">
      <div class="panel-head">
        <h2>Заказы</h2>
        <button id="btnReloadOrders" class="btn btn-light" type="button">Обновить</button>
      </div>
      <div id="ordersList" class="order-list"></div>
    </section>
  </main>
</div>

<div id="editModal" class="modal-backdrop" style="display:none">
  <div class="modal-card">
    <div class="modal-head">
      <h2 style="margin:0">Редактировать товар</h2>
      <button id="editClose" class="close-btn" type="button">×</button>
    </div>

    <div class="modal-grid">
      <div class="row">
        <div>
          <label for="eName">Название</label>
          <input id="eName" placeholder="Название товара">
        </div>
        <div>
          <label for="eCategory">Категория</label>
          <select id="eCategory"><option value="">Выберите категорию</option></select>
        </div>
        <div>
          <label for="eBrand">Бренд</label>
          <input id="eBrand" placeholder="Samsung / Xiaomi / Apple">
        </div>
        <div>
          <label for="ePrice">Цена</label>
          <input id="ePrice" placeholder="0" inputmode="decimal">
        </div>
        <div>
          <label for="eStock">Остаток</label>
          <input id="eStock" placeholder="0" inputmode="numeric">
        </div>
        <div>
          <label for="eColor">Цвет</label>
          <input id="eColor" placeholder="Цвет">
        </div>
        <div>
          <label for="eCondition">Состояние</label>
          <input id="eCondition" placeholder="новый / б/у">
        </div>
      </div>

      <div class="row">
        <div>
          <label for="eCountry">Страна</label>
          <input id="eCountry" placeholder="Страна">
        </div>
        <div>
          <label for="eMaterial">Материал</label>
          <input id="eMaterial" placeholder="Материал">
        </div>
        <div>
          <label for="eDesc">Описание</label>
          <textarea id="eDesc" placeholder="Краткое описание"></textarea>
        </div>
        <div class="upload-wrap">
          <div class="inline" style="justify-content:space-between">
            <button id="eImageReplaceBtn" class="btn btn-light" type="button">Загрузить новое изображение</button>
            <button id="eImageRemoveBtn" class="btn btn-danger" type="button" style="display:none">Удалить</button>
          </div>
          <input id="eImage" type="file" accept="image/*" style="display:none">
          <div id="eImagePreview" class="preview">Изображение не выбрано</div>
          <div class="helper">Изменения применяются после нажатия «Сохранить».</div>
        </div>
      </div>
    </div>

    <div class="modal-actions">
      <button id="cancelEdit" class="btn btn-light" type="button">Отмена</button>
      <button id="saveEdit" class="btn btn-primary" type="button">Сохранить</button>
    </div>
  </div>
</div>

<script>
async function api(path, method, body, isJson){
  if(!method) method = 'GET';
  if(isJson === undefined) isJson = true;
  var opts = {method: method, headers: {}, credentials: 'include'};
  if(isJson && body){ opts.headers['Content-Type'] = 'application/json'; opts.body = JSON.stringify(body); }
  if(!isJson && body){ opts.body = body; }
  var r = await fetch(path, opts);
  if(!r.ok){ throw new Error(await r.text()); }
  if(r.status === 204){ return {}; }
  return r.json();
}
function esc(v){
  return String(v || '')
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}
function toNumber(v){
  var n = Number(v);
  return Number.isFinite(n) ? n : 0;
}
function toInt(v){
  var n = parseInt(v, 10);
  return Number.isFinite(n) ? n : 0;
}
function money(v){
  return new Intl.NumberFormat('ru-RU', {style:'currency', currency:'RUB', maximumFractionDigits:0}).format(Number(v || 0));
}
function formatDate(v){
  if(!v) return '';
  var d = new Date(v);
  return Number.isNaN(d.getTime()) ? String(v) : d.toLocaleString('ru-RU');
}
function categoryOptionsHtml(cats){
  var html = '<option value="">Выберите категорию</option>';
  for(var i=0;i<cats.length;i++){
    var c = cats[i];
    html += '<option value="' + c.id + '">' + esc(c.name) + '</option>';
  }
  return html;
}
async function loadCategories(){
  var cats = await api('/categories', 'GET');
  if(!Array.isArray(cats)) cats = [];
  var html = categoryOptionsHtml(cats);
  var pSel = document.getElementById('pCategory');
  var eSel = document.getElementById('eCategory');
  pSel.innerHTML = html;
  eSel.innerHTML = html;
  return cats;
}
function renderProductCard(p){
  var hasImage = !!p.image_url;
  var image = hasImage
    ? '<div class="card-image"><img src="' + esc(p.image_url) + '" alt=""></div>'
    : '<div class="card-image">Нет изображения</div>';
  var desc = p.description ? esc(p.description) : 'Нет описания';
  var meta = 'Бренд: ' + (p.brand || 'не указан') + ' | Цена: ' + (p.price || 0) + ' | Остаток: ' + (p.stock || 0);
  return ''
    + '<article class="product-card" onclick="editProduct(' + p.id + ')">'
    + image
    + '<div class="card-body">'
    + '<h3 class="card-title">' + esc(p.name || '') + '</h3>'
    + '<div class="card-meta">' + meta + '</div>'
    + '<p class="card-desc">' + desc + '</p>'
    + '<div class="card-actions">'
    + '<button class="btn btn-light" type="button" onclick="event.stopPropagation(); editProduct(' + p.id + ')">Изменить</button>'
    + '<button class="btn btn-danger" type="button" onclick="event.stopPropagation(); delProduct(' + p.id + ')">Удалить</button>'
    + '</div>'
    + '</div>'
    + '</article>';
}
async function loadProducts(){
  var data = await api('/products?page=1&per_page=100', 'GET');
  var items = data && Array.isArray(data.items) ? data.items : [];
  var list = document.getElementById('prodList');
  list.innerHTML = '';
  if(!items.length){
    list.innerHTML = '<div class="empty-state">Товаров пока нет. Создайте первый товар выше.</div>';
  }else{
    for(var i=0;i<items.length;i++){
      list.insertAdjacentHTML('beforeend', renderProductCard(items[i]));
    }
  }
  document.getElementById('statsLine').textContent = items.length + ' товаров';
}
function loadCanvasImage(file){
  return new Promise(function(resolve, reject){
    var url = URL.createObjectURL(file);
    var img = new Image();
    img.onload = function(){
      URL.revokeObjectURL(url);
      resolve(img);
    };
    img.onerror = function(){
      URL.revokeObjectURL(url);
      reject(new Error('Не удалось прочитать изображение'));
    };
    img.src = url;
  });
}
async function normalizeProductImage(file){
  if(!file) return null;
  if(!file.type || file.type.indexOf('image/') !== 0){
    throw new Error('Загрузите изображение товара в формате JPG, PNG или WEBP.');
  }
  var img = await loadCanvasImage(file);
  var width = 900;
  var height = 1100;
  var padding = 44;
  var canvas = document.createElement('canvas');
  canvas.width = width;
  canvas.height = height;
  var ctx = canvas.getContext('2d');
  ctx.fillStyle = '#ffffff';
  ctx.fillRect(0, 0, width, height);
  var scale = Math.min((width - padding * 2) / img.naturalWidth, (height - padding * 2) / img.naturalHeight);
  var drawWidth = Math.round(img.naturalWidth * scale);
  var drawHeight = Math.round(img.naturalHeight * scale);
  var x = Math.round((width - drawWidth) / 2);
  var y = Math.round((height - drawHeight) / 2);
  ctx.imageSmoothingEnabled = true;
  ctx.imageSmoothingQuality = 'high';
  ctx.drawImage(img, x, y, drawWidth, drawHeight);
  var blob = await new Promise(function(resolve){
    canvas.toBlob(resolve, 'image/webp', 0.92);
  });
  if(!blob) throw new Error('Не удалось подготовить изображение товара');
  var base = String(file.name || 'product').replace(/\.[^.]+$/, '').replace(/[^a-zA-Z0-9_-]+/g, '-').replace(/^-+|-+$/g, '') || 'product';
  return new File([blob], base + '-900x1100.webp', {type: 'image/webp'});
}
function orderStatusLabel(status){
  if(status === 'pending') return 'Заказ принят';
  if(status === 'processing') return 'В обработке';
  if(status === 'shipped') return 'В пути';
  if(status === 'delivered') return 'Заказ доставлен';
  if(status === 'cancelled') return 'Отменен';
  return 'Заказ принят';
}
function orderStatusOptions(selected){
  var statuses = [
    ['pending', 'Заказ принят'],
    ['processing', 'В обработке'],
    ['shipped', 'В пути'],
    ['delivered', 'Заказ доставлен'],
    ['cancelled', 'Отменен']
  ];
  var html = '';
  for(var i=0;i<statuses.length;i++){
    var st = statuses[i];
    html += '<option value="' + st[0] + '"' + (st[0] === selected ? ' selected' : '') + '>' + st[1] + '</option>';
  }
  return html;
}
function renderOrderCard(o){
  var items = Array.isArray(o.items) ? o.items : [];
  var itemHtml = items.length
    ? items.map(function(it){
        return '<li>' + esc(it.product_name || 'Товар') + ' × ' + (it.quantity || 1) + ' · ' + money(it.price) + '</li>';
      }).join('')
    : '<li>Состав заказа не загружен</li>';
  var place = o.delivery_type === 'pickup' ? o.pickup_point : o.address;
  return ''
    + '<article class="order-card">'
    + '<div class="order-head">'
    + '<div><h3>Заказ №' + o.id + '</h3><div class="order-meta">' + formatDate(o.created_at) + '</div></div>'
    + '<span class="status-pill">' + orderStatusLabel(o.status) + '</span>'
    + '</div>'
    + '<div class="order-meta"><b>' + esc(o.customer_name || 'Покупатель') + '</b><br>'
    + esc(o.email || '') + (o.phone ? ' · ' + esc(o.phone) : '') + '<br>'
    + esc(o.delivery_type === 'pickup' ? 'Самовывоз' : 'Курьер') + (place ? ': ' + esc(place) : '') + '</div>'
    + '<ul class="order-items">' + itemHtml + '</ul>'
    + '<div class="order-total">Итого: ' + money(o.total) + '</div>'
    + '<div class="order-status-row">'
    + '<select data-order-status="' + o.id + '">' + orderStatusOptions(o.status) + '</select>'
    + '<button class="btn btn-primary" type="button" data-save-order-status="' + o.id + '">Сохранить статус</button>'
    + '</div>'
    + '</article>';
}
async function loadOrders(){
  var data = await api('/admin/api/orders?page=1&per_page=50', 'GET');
  var items = data && Array.isArray(data.items) ? data.items : [];
  var list = document.getElementById('ordersList');
  if(!items.length){
    list.innerHTML = '<div class="empty-state">Заказов пока нет. Когда пользователь оформит покупку, она появится здесь.</div>';
    return;
  }
  list.innerHTML = items.map(renderOrderCard).join('');
}
async function saveOrderStatus(id){
  var select = document.querySelector('[data-order-status="' + id + '"]');
  if(!select) return;
  await api('/admin/api/orders/' + id + '/status', 'PUT', {status: select.value});
  await loadOrders();
}
async function uploadFile(file){
  file = await normalizeProductImage(file);
  var fd = new FormData();
  fd.append('file', file);
  var r = await fetch('/admin/api/upload', {method: 'POST', body: fd, credentials: 'include'});
  if(!r.ok) throw new Error('Ошибка загрузки файла');
  return (await r.json()).url;
}

var pImageInput = document.getElementById('pImage');
var pImagePreview = document.getElementById('pImagePreview');
var pImageRemoveBtn = document.getElementById('pImageRemove');
function resetCreateImagePreview(){
  if(pImagePreview) pImagePreview.innerHTML = 'Изображение не выбрано';
  if(pImageRemoveBtn) pImageRemoveBtn.style.display = 'none';
}
if(pImageInput){
  pImageInput.addEventListener('change', function(){
    var f = this.files && this.files[0] ? this.files[0] : null;
    if(!f){ resetCreateImagePreview(); return; }
    var reader = new FileReader();
    reader.onload = function(e){
      pImagePreview.innerHTML = '<img src="' + e.target.result + '" alt="">';
      pImageRemoveBtn.style.display = 'inline-flex';
    };
    reader.readAsDataURL(f);
  });
}
if(pImageRemoveBtn){
  pImageRemoveBtn.addEventListener('click', function(){
    if(pImageInput) pImageInput.value = '';
    resetCreateImagePreview();
  });
}

document.getElementById('btnCreateCat').addEventListener('click', async function(){
  try{
    var name = document.getElementById('catName').value.trim();
    if(!name) return alert('Введите название категории.');
    await api('/admin/api/categories', 'POST', {name: name});
    document.getElementById('catName').value = '';
    await loadCategories();
    await loadProducts();
  }catch(e){ alert(e.message); }
});

document.getElementById('btnCreateProd').addEventListener('click', async function(){
  try{
    var file = pImageInput && pImageInput.files ? pImageInput.files[0] : null;
    var image_url = null;
    if(file) image_url = await uploadFile(file);
    var cat = document.getElementById('pCategory').value;
    await api('/admin/api/products', 'POST', {
      name: document.getElementById('pName').value,
      price: toNumber(document.getElementById('pPrice').value),
      stock: toInt(document.getElementById('pStock').value),
      category_id: cat ? parseInt(cat, 10) : null,
      brand: document.getElementById('pBrand').value || '',
      description: document.getElementById('pDesc').value,
      image_url: image_url,
      color: document.getElementById('pColor').value || '',
      condition: document.getElementById('pCondition').value || '',
      country: document.getElementById('pCountry').value || '',
      material: document.getElementById('pMaterial').value || ''
    });
    document.getElementById('pName').value = '';
    document.getElementById('pPrice').value = '';
    document.getElementById('pStock').value = '';
    document.getElementById('pBrand').value = '';
    document.getElementById('pDesc').value = '';
    document.getElementById('pColor').value = '';
    document.getElementById('pCondition').value = '';
    document.getElementById('pCountry').value = '';
    document.getElementById('pMaterial').value = '';
    document.getElementById('pCategory').value = '';
    if(pImageInput) pImageInput.value = '';
    resetCreateImagePreview();
    await loadProducts();
  }catch(e){ alert(e.message); }
});

var eRemoveImage = false;
var eSelectedFile = null;
var eImageInput = document.getElementById('eImage');
var eImagePreview = document.getElementById('eImagePreview');
var eImageRemoveBtn = document.getElementById('eImageRemoveBtn');
var eImageReplaceBtn = document.getElementById('eImageReplaceBtn');
function setEditImagePreview(url, removedText){
  if(removedText){
    eImagePreview.innerHTML = removedText;
    eImageRemoveBtn.style.display = 'none';
    return;
  }
  if(url){
    eImagePreview.innerHTML = '<img src="' + esc(url) + '" alt="">';
    eImageRemoveBtn.style.display = 'inline-flex';
  }else{
    eImagePreview.innerHTML = 'Изображение не выбрано';
    eImageRemoveBtn.style.display = 'none';
  }
}
if(eImageReplaceBtn){
  eImageReplaceBtn.addEventListener('click', function(){
    if(eImageInput) eImageInput.click();
  });
}
if(eImageInput){
  eImageInput.addEventListener('change', function(){
    var f = this.files && this.files[0] ? this.files[0] : null;
    eSelectedFile = f;
    eRemoveImage = false;
    if(!f){ setEditImagePreview('', ''); return; }
    var reader = new FileReader();
    reader.onload = function(e){
      setEditImagePreview(e.target.result, '');
    };
    reader.readAsDataURL(f);
  });
}
if(eImageRemoveBtn){
  eImageRemoveBtn.addEventListener('click', function(){
    eRemoveImage = true;
    eSelectedFile = null;
    if(eImageInput) eImageInput.value = '';
    setEditImagePreview('', 'Изображение будет удалено после сохранения.');
  });
}

async function editProduct(id){
  try{
    var prod = await api('/products/' + id, 'GET');
    await loadCategories();
    document.getElementById('eName').value = prod.name || '';
    document.getElementById('ePrice').value = prod.price || '';
    document.getElementById('eStock').value = prod.stock || '';
    document.getElementById('eCategory').value = prod.category_id || '';
    document.getElementById('eBrand').value = prod.brand || '';
    document.getElementById('eColor').value = prod.color || '';
    document.getElementById('eCondition').value = prod.condition || '';
    document.getElementById('eCountry').value = prod.country || '';
    document.getElementById('eMaterial').value = prod.material || '';
    document.getElementById('eDesc').value = prod.description || '';
    if(eImageInput) eImageInput.value = '';
    eSelectedFile = null;
    eRemoveImage = false;
    setEditImagePreview(prod.image_url || '', '');
    var modal = document.getElementById('editModal');
    modal.dataset.editing = id;
    modal.style.display = 'flex';
  }catch(e){
    alert('Не удалось загрузить товар: ' + e.message);
  }
}
function closeEditModal(){
  var modal = document.getElementById('editModal');
  modal.style.display = 'none';
  delete modal.dataset.editing;
  eSelectedFile = null;
  eRemoveImage = false;
}
document.getElementById('editClose').addEventListener('click', closeEditModal);
document.getElementById('cancelEdit').addEventListener('click', closeEditModal);
document.getElementById('editModal').addEventListener('click', function(e){
  if(e.target && e.target.id === 'editModal') closeEditModal();
});
document.getElementById('saveEdit').addEventListener('click', async function(){
  try{
    var id = document.getElementById('editModal').dataset.editing;
    if(!id) return alert('Товар не выбран.');
    var cat = document.getElementById('eCategory').value;
    var payload = {
      name: document.getElementById('eName').value,
      price: toNumber(document.getElementById('ePrice').value),
      stock: toInt(document.getElementById('eStock').value),
      category_id: cat ? parseInt(cat, 10) : null,
      brand: document.getElementById('eBrand').value || '',
      color: document.getElementById('eColor').value || '',
      condition: document.getElementById('eCondition').value || '',
      country: document.getElementById('eCountry').value || '',
      material: document.getElementById('eMaterial').value || '',
      description: document.getElementById('eDesc').value || ''
    };
    if(eSelectedFile){
      payload.image_url = await uploadFile(eSelectedFile);
    }else if(eRemoveImage){
      payload.image_url = null;
    }
    await api('/admin/api/products/' + id, 'PUT', payload);
    closeEditModal();
    await loadProducts();
  }catch(e){
    alert('Не удалось сохранить товар: ' + e.message);
  }
});
async function delProduct(id){
  try{
    if(!confirm('Удалить этот товар?')) return;
    await api('/admin/api/products/' + id, 'DELETE');
    await loadProducts();
  }catch(e){ alert(e.message); }
}
document.getElementById('btnLogout').addEventListener('click', async function(){
  await fetch('/auth/logout', {method: 'POST', credentials: 'include'});
  location.href = '/auth';
});
document.getElementById('btnReloadOrders').addEventListener('click', async function(){
  try{
    await loadOrders();
  }catch(e){
    alert('Не удалось загрузить заказы: ' + e.message);
  }
});
document.getElementById('ordersList').addEventListener('click', async function(e){
  var btn = e.target.closest('[data-save-order-status]');
  if(!btn) return;
  try{
    btn.disabled = true;
    await saveOrderStatus(btn.dataset.saveOrderStatus);
  }catch(err){
    alert('Не удалось изменить статус: ' + err.message);
  }finally{
    btn.disabled = false;
  }
});
async function initAdmin(){
  await loadCategories();
  await loadProducts();
  await loadOrders();
}
initAdmin();
</script>
</body>
</html>`
