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
  animation:rise .32s ease both;
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
.category-native{display:none}
input:focus,select:focus,textarea:focus{
  border-color:#7bb8f4;
  box-shadow:0 0 0 4px rgba(97,174,255,0.2);
}
.category-picker{position:relative}
.category-picker-toggle{
  width:100%;
  min-height:42px;
  border:1px solid var(--line);
  border-radius:12px;
  background:#fff;
  color:#0f2132;
  padding:10px 40px 10px 12px;
  text-align:left;
  font:500 14px/1.35 'Manrope','Segoe UI',sans-serif;
  cursor:pointer;
  transition:border-color .18s ease,box-shadow .18s ease;
}
.category-picker-toggle:after{
  content:"";
  position:absolute;
  right:14px;
  top:17px;
  width:8px;
  height:8px;
  border-right:2px solid #486076;
  border-bottom:2px solid #486076;
  transform:rotate(45deg);
}
.category-picker.open .category-picker-toggle{border-color:#7bb8f4;box-shadow:0 0 0 4px rgba(97,174,255,0.2)}
.category-picker.open .category-picker-toggle:after{top:21px;transform:rotate(225deg)}
.category-picker-panel{
  position:absolute;
  left:0;
  right:0;
  top:calc(100% + 6px);
  z-index:50;
  display:none;
  max-height:230px;
  overflow:auto;
  border:1px solid var(--line);
  border-radius:14px;
  background:#fff;
  box-shadow:0 18px 36px rgba(16,52,82,.16);
  padding:8px;
}
.category-picker.open .category-picker-panel{display:grid;gap:4px}
.category-option{
  display:flex;
  align-items:center;
  gap:9px;
  min-height:36px;
  padding:8px 9px;
  border-radius:10px;
  color:#21384e;
  cursor:pointer;
  font-size:14px;
}
.category-option:hover{background:#f3f8fc}
.category-option input{width:16px;height:16px;margin:0;accent-color:var(--brand)}
.category-option span{min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
.category-option[data-level="1"]{padding-left:24px}
.category-option[data-level="2"]{padding-left:40px}
.category-empty{padding:10px;color:var(--muted);font-size:13px}
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
.category-list{
  display:grid;
  gap:10px;
}
.category-card{
  border:1px solid var(--line);
  border-radius:14px;
  background:#fff;
  padding:12px;
  display:grid;
  grid-template-columns:minmax(180px,1fr) auto;
  gap:10px;
  align-items:center;
  animation:rise .22s ease both;
  transition:transform .16s ease, box-shadow .16s ease, border-color .16s ease;
}
.category-card:hover{transform:translateY(-2px);border-color:#c6d9ea;box-shadow:0 10px 20px rgba(16,52,82,.08)}
.category-card.child{
  margin-left:24px;
  background:#f8fbff;
}
.category-name{
  font-weight:900;
}
.category-meta{
  color:var(--muted);
  font-size:13px;
  margin-top:3px;
}
.order-list{
  display:grid;
  gap:12px;
}
.user-list{
  display:grid;
  gap:12px;
}
.user-card{
  border:1px solid var(--line);
  border-radius:16px;
  background:#fff;
  padding:14px;
  display:grid;
  grid-template-columns:minmax(220px,1fr) 160px auto;
  gap:10px;
  align-items:end;
  animation:rise .22s ease both;
  transition:transform .16s ease, box-shadow .16s ease, border-color .16s ease;
}
.user-card:hover{transform:translateY(-2px);border-color:#c6d9ea;box-shadow:0 10px 22px rgba(16,52,82,.08)}
.user-title{
  font-weight:900;
  margin-bottom:6px;
}
.user-meta{
  color:var(--muted);
  font-size:13px;
}
.user-hash{
  margin-top:8px;
  max-width:100%;
  color:#5b6f84;
  font-family:Consolas,Menlo,monospace;
  font-size:12px;
  overflow:hidden;
  text-overflow:ellipsis;
  white-space:nowrap;
}
.secret-toggle{
  border:1px solid var(--line);
  border-radius:10px;
  background:#f4f8fc;
  color:#486076;
  min-height:40px;
  padding:0 10px;
  font-weight:800;
  transition:transform .16s ease, background .16s ease;
}
.secret-toggle:hover{transform:translateY(-1px);background:#eaf2fb}
.role-display{
  min-height:40px;
  display:flex;
  align-items:center;
  padding:0 12px;
  border:1px solid var(--line);
  border-radius:12px;
  background:#f8fbff;
  color:#486076;
  font-weight:800;
}
.user-actions{
  display:flex;
  gap:8px;
  flex-wrap:wrap;
  align-items:center;
}
.order-card{
  border:1px solid var(--line);
  border-radius:16px;
  background:#fff;
  padding:14px;
  animation:rise .22s ease both;
  transition:transform .16s ease, box-shadow .16s ease, border-color .16s ease;
}
.order-card:hover{transform:translateY(-2px);border-color:#c6d9ea;box-shadow:0 10px 22px rgba(16,52,82,.08)}
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
.modal-backdrop[style*="flex"]{animation:rise .16s ease both}
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
  animation:modalPop .2s ease both;
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
  transition:transform .16s ease, background .16s ease;
}
.close-btn:hover{transform:rotate(4deg);background:#e2edf8}
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
@keyframes modalPop{
  from{opacity:0;transform:translateY(10px) scale(.98)}
  to{opacity:1;transform:translateY(0) scale(1)}
}
@media(prefers-reduced-motion:reduce){
  *,*:before,*:after{animation:none!important;transition:none!important}
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
  .category-card{grid-template-columns:1fr}
  .category-card.child{margin-left:12px}
  .user-card{grid-template-columns:1fr}
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
      <div class="row">
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
          <label for="pCategory">Категории</label>
          <select id="pCategory" class="category-native" multiple><option value="">Выберите категории</option></select>
          <div class="category-picker" data-category-picker="pCategory"></div>
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
        <h2>Каталоги и подкаталоги</h2>
        <button id="btnReloadCategories" class="btn btn-light" type="button">Обновить</button>
      </div>
      <div id="catList" class="category-list"></div>
    </section>

    <section class="panel panel-list">
      <div class="panel-head">
        <h2>Заказы</h2>
        <button id="btnReloadOrders" class="btn btn-light" type="button">Обновить</button>
      </div>
      <div id="ordersList" class="order-list"></div>
    </section>

    <section class="panel panel-list">
      <div class="panel-head">
        <h2>Пользователи</h2>
        <button id="btnReloadUsers" class="btn btn-light" type="button">Обновить</button>
      </div>
      <div id="usersList" class="user-list"></div>
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
          <label for="eCategory">Категории</label>
          <select id="eCategory" class="category-native" multiple><option value="">Выберите категории</option></select>
          <div class="category-picker" data-category-picker="eCategory"></div>
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
  var rows = buildCategoryRows(cats);
  for(var i=0;i<rows.length;i++){
    var row = rows[i];
    html += '<option value="' + row.category.id + '">' + esc(Array(row.level + 1).join('- ') + row.category.name) + '</option>';
  }
  return html;
}
function selectedCategoryIDs(id){
  var select = document.getElementById(id);
  if(!select) return [];
  var ids = [];
  for(var i=0;i<select.options.length;i++){
    var option = select.options[i];
    if(option.selected && option.value){
      ids.push(parseInt(option.value, 10));
    }
  }
  return ids.filter(function(id){ return Number.isFinite(id) && id > 0; });
}
function setSelectedCategoryIDs(id, ids){
  var select = document.getElementById(id);
  if(!select) return;
  var set = {};
  (ids || []).forEach(function(value){ set[String(value)] = true; });
  for(var i=0;i<select.options.length;i++){
    select.options[i].selected = !!set[String(select.options[i].value)];
  }
  updateCategoryPicker(id);
}
function categoryLabelById(id){
  var select = document.getElementById(id);
  if(!select) return '';
  var names = [];
  for(var i=0;i<select.options.length;i++){
    var option = select.options[i];
    if(option.selected && option.value){
      names.push(option.textContent.replace(/^-+\s*/, ''));
    }
  }
  return names.length ? names.join(', ') : 'Выберите категории';
}
function updateCategoryPicker(id){
  var picker = document.querySelector('[data-category-picker="' + id + '"]');
  if(!picker) return;
  var text = picker.querySelector('[data-category-picker-text]');
  if(text) text.textContent = categoryLabelById(id);
  var ids = selectedCategoryIDs(id).map(String);
  picker.querySelectorAll('input[type="checkbox"]').forEach(function(input){
    input.checked = ids.indexOf(input.value) !== -1;
  });
}
function renderCategoryPicker(id, cats){
  var picker = document.querySelector('[data-category-picker="' + id + '"]');
  if(!picker) return;
  var rows = buildCategoryRows(cats);
  var body = rows.length ? rows.map(function(row){
    var name = esc(Array(row.level + 1).join('- ') + row.category.name);
    return '<label class="category-option" data-level="' + row.level + '">'
      + '<input type="checkbox" value="' + row.category.id + '">'
      + '<span>' + name + '</span>'
      + '</label>';
  }).join('') : '<div class="category-empty">Категорий пока нет</div>';
  picker.innerHTML = ''
    + '<button class="category-picker-toggle" type="button"><span data-category-picker-text>Выберите категории</span></button>'
    + '<div class="category-picker-panel">' + body + '</div>';
  updateCategoryPicker(id);
}
function buildCategoryRows(cats){
  var byParent = {};
  for(var i=0;i<cats.length;i++){
    var c = cats[i];
    var key = c.parent_id || 0;
    if(!byParent[key]) byParent[key] = [];
    byParent[key].push(c);
  }
  Object.keys(byParent).forEach(function(key){
    byParent[key].sort(function(a,b){ return String(a.name || '').localeCompare(String(b.name || ''), 'ru'); });
  });
  var rows = [];
  function walk(parentId, level){
    var list = byParent[parentId] || [];
    for(var j=0;j<list.length;j++){
      rows.push({category:list[j], level:level});
      walk(list[j].id, level + 1);
    }
  }
  walk(0, 0);
  return rows;
}
function renderCategoryList(cats){
  var list = document.getElementById('catList');
  if(!list) return;
  var rows = buildCategoryRows(cats);
  if(!rows.length){
    list.innerHTML = '<div class="empty-state">Каталогов пока нет. Создайте первый каталог выше.</div>';
    return;
  }
  list.innerHTML = rows.map(function(row){
    var c = row.category;
    var isChild = row.level > 0;
    var meta = isChild ? 'Подкаталог' : 'Основной каталог';
    return ''
      + '<article class="category-card' + (isChild ? ' child' : '') + '">'
      + '<div><div class="category-name">' + esc(c.name || '') + '</div>'
      + '<div class="category-meta">' + meta + ' | ID: ' + c.id + '</div></div>'
      + '<button class="btn btn-danger" type="button" data-delete-category="' + c.id + '">Удалить</button>'
      + '</article>';
  }).join('');
}
async function loadCategories(){
  var cats = await api('/categories', 'GET');
  if(!Array.isArray(cats)) cats = [];
  var html = categoryOptionsHtml(cats);
  var pSel = document.getElementById('pCategory');
  var eSel = document.getElementById('eCategory');
  pSel.innerHTML = html;
  eSel.innerHTML = html;
  renderCategoryPicker('pCategory', cats);
  renderCategoryPicker('eCategory', cats);
  renderCategoryList(cats);
  return cats;
}
function renderProductCard(p){
  var hasImage = !!p.image_url;
  var image = hasImage
    ? '<div class="card-image"><img src="' + esc(p.image_url) + '" alt=""></div>'
    : '<div class="card-image">Нет изображения</div>';
  var desc = p.description ? esc(p.description) : 'Нет описания';
  var cats = Array.isArray(p.categories) && p.categories.length ? p.categories.map(function(c){ return c.name; }).join(', ') : 'не указаны';
  var meta = 'Бренд: ' + (p.brand || 'не указан') + ' | Категории: ' + cats + ' | Цена: ' + (p.price || 0) + ' | Остаток: ' + (p.stock || 0);
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
    + '<button class="btn btn-danger" type="button" data-delete-order="' + o.id + '">Удалить</button>'
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
function renderUserCard(u){
  var role = u.role === 'admin' ? 'admin' : 'user';
  var isBuiltInAdmin = u.username === 'admin';
  var hash = u.password_hash || '';
  return ''
    + '<article class="user-card" data-user-id="' + u.id + '">'
    + '<div>'
    + '<div class="user-title">' + esc(u.username || '') + '</div>'
    + '<div class="user-meta">ID: ' + u.id + ' | роль: ' + role + '</div>'
    + '<div class="user-hash" title="' + esc(hash) + '" data-user-hash-text="' + esc(hash || 'нет') + '">Пароль: ' + maskedSecret(hash) + '</div>'
    + '<button class="secret-toggle" type="button" data-toggle-user-hash="' + u.id + '">Показать</button>'
    + '</div>'
    + '<div><label>Роль</label><div class="role-display">' + role + '</div></div>'
    + '<div class="user-actions">'
    + '<button class="btn btn-danger" type="button" data-delete-user="' + u.id + '" ' + (isBuiltInAdmin ? 'disabled' : '') + '>Удалить</button>'
    + '</div>'
    + '</article>';
}
function maskedSecret(value){
  if(!value) return 'нет';
  return '••••••••••••';
}
async function loadUsers(){
  var data = await api('/admin/api/users?page=1&per_page=100', 'GET');
  var items = data && Array.isArray(data.items) ? data.items : [];
  var list = document.getElementById('usersList');
  if(!items.length){
    list.innerHTML = '<div class="empty-state">Пользователей пока нет.</div>';
    return;
  }
  list.innerHTML = items.map(renderUserCard).join('');
}
async function saveOrderStatus(id){
  var select = document.querySelector('[data-order-status="' + id + '"]');
  if(!select) return;
  await api('/admin/api/orders/' + id + '/status', 'PUT', {status: select.value});
  await loadOrders();
}
async function deleteOrder(id){
  if(!confirm('Удалить заказ #' + id + '?')) return;
  await api('/admin/api/orders/' + id, 'DELETE');
  await loadOrders();
}
async function deleteUser(id){
  var card = document.querySelector('[data-user-id="' + id + '"]');
  var name = card ? card.querySelector('.user-title').textContent : id;
  if(!confirm('Удалить пользователя ' + name + '?')) return;
  await api('/admin/api/users/' + id, 'DELETE');
  await loadUsers();
}
async function deleteCategory(id){
  if(!confirm('Удалить каталог? Подкаталоги тоже будут удалены, товары останутся без категории.')) return;
  await api('/admin/api/categories/' + id, 'DELETE');
  await loadCategories();
  await loadProducts();
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
    var categoryIds = selectedCategoryIDs('pCategory');
    await api('/admin/api/products', 'POST', {
      name: document.getElementById('pName').value,
      price: toNumber(document.getElementById('pPrice').value),
      stock: toInt(document.getElementById('pStock').value),
      category_id: categoryIds.length ? categoryIds[0] : null,
      category_ids: categoryIds,
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
    setSelectedCategoryIDs('pCategory', []);
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
    setSelectedCategoryIDs('eCategory', prod.category_ids && prod.category_ids.length ? prod.category_ids : (prod.category_id ? [prod.category_id] : []));
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
    var categoryIds = selectedCategoryIDs('eCategory');
    var payload = {
      name: document.getElementById('eName').value,
      price: toNumber(document.getElementById('ePrice').value),
      stock: toInt(document.getElementById('eStock').value),
      category_id: categoryIds.length ? categoryIds[0] : null,
      category_ids: categoryIds,
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
document.getElementById('btnReloadCategories').addEventListener('click', async function(){
  try{
    await loadCategories();
  }catch(e){
    alert('Не удалось загрузить каталоги: ' + e.message);
  }
});
document.getElementById('btnReloadUsers').addEventListener('click', async function(){
  try{
    await loadUsers();
  }catch(e){
    alert('Не удалось загрузить пользователей: ' + e.message);
  }
});
document.getElementById('ordersList').addEventListener('click', async function(e){
  var btn = e.target.closest('[data-save-order-status], [data-delete-order]');
  if(!btn) return;
  try{
    btn.disabled = true;
    if(btn.dataset.saveOrderStatus){
      await saveOrderStatus(btn.dataset.saveOrderStatus);
    }else{
      await deleteOrder(btn.dataset.deleteOrder);
    }
  }catch(err){
    alert('Не удалось выполнить действие с заказом: ' + err.message);
  }finally{
    btn.disabled = false;
  }
});
document.getElementById('catList').addEventListener('click', async function(e){
  var btn = e.target.closest('[data-delete-category]');
  if(!btn) return;
  try{
    btn.disabled = true;
    await deleteCategory(btn.dataset.deleteCategory);
  }catch(err){
    alert('Не удалось удалить каталог: ' + err.message);
  }finally{
    btn.disabled = false;
  }
});
document.getElementById('usersList').addEventListener('click', async function(e){
  var hashToggle = e.target.closest('[data-toggle-user-hash]');
  if(hashToggle){
    var hashCard = hashToggle.closest('[data-user-id]');
    var hashEl = hashCard ? hashCard.querySelector('[data-user-hash-text]') : null;
    if(!hashEl) return;
    var showHash = hashToggle.dataset.visible !== '1';
    hashToggle.dataset.visible = showHash ? '1' : '0';
    hashEl.textContent = 'Пароль: ' + (showHash ? hashEl.dataset.userHashText : maskedSecret(hashEl.dataset.userHashText));
    hashToggle.textContent = showHash ? 'Скрыть' : 'Показать';
    return;
  }
  var btn = e.target.closest('[data-delete-user]');
  if(!btn) return;
  try{
    btn.disabled = true;
    await deleteUser(btn.dataset.deleteUser);
  }catch(err){
    alert('Не удалось выполнить действие с пользователем: ' + err.message);
  }finally{
    btn.disabled = false;
  }
});
document.addEventListener('click', function(e){
  var toggle = e.target.closest('.category-picker-toggle');
  var picker = e.target.closest('.category-picker');
  document.querySelectorAll('.category-picker.open').forEach(function(openPicker){
    if(openPicker !== picker) openPicker.classList.remove('open');
  });
  if(toggle && picker){
    picker.classList.toggle('open');
    return;
  }
  var option = e.target.closest('.category-option input');
  if(option && picker){
    var id = picker.dataset.categoryPicker;
    var select = document.getElementById(id);
    if(select){
      for(var i=0;i<select.options.length;i++){
        if(select.options[i].value === option.value){
          select.options[i].selected = option.checked;
          break;
        }
      }
      updateCategoryPicker(id);
    }
    return;
  }
  if(!picker){
    document.querySelectorAll('.category-picker.open').forEach(function(openPicker){
      openPicker.classList.remove('open');
    });
  }
});
async function initAdmin(){
  await loadCategories();
  await loadProducts();
  await loadOrders();
  await loadUsers();
}
initAdmin();
</script>
</body>
</html>`
