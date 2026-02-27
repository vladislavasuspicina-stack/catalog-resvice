package main

// ==================== TEMPLATES (inline) ====================
const indexHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Премиум Каталог</title>
<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" crossorigin="">
<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js" crossorigin=""></script>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'Segoe UI',Tahoma,Geneva,Verdana,sans-serif;background:linear-gradient(135deg,#f5f7fa 0%,#c3cfe2 100%);min-height:100vh;color:#2c3e50;padding:20px 0}
.navbar{background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;padding:0 20px;position:sticky;top:0;z-index:100;box-shadow:0 4px 15px rgba(0,0,0,0.15)}
.navbar-content{max-width:1100px;margin:0 auto;display:flex;justify-content:space-between;align-items:center;padding:16px 0}
.navbar-title{font-size:28px;font-weight:700;letter-spacing:1.2px}
.navbar-actions{display:flex;gap:16px;align-items:center}
.navbar-actions a{color:#fff;text-decoration:none;font-weight:500;transition:opacity .3s;padding:8px 12px;border-radius:6px}
.navbar-actions a:hover{opacity:0.85;background:rgba(255,255,255,0.1)}
.container{max-width:1100px;margin:0 auto;padding:32px 20px}
.header{margin-bottom:32px}
.header h1{font-size:32px;font-weight:700;margin-bottom:8px;color:#2c3e50}
.header p{color:#7f8c8d;font-size:16px}
.controls{background:#fff;padding:20px;border-radius:12px;box-shadow:0 4px 15px rgba(0,0,0,0.08);display:flex;gap:12px;align-items:center;flex-wrap:wrap;margin-bottom:32px}
.controls label{font-weight:500;color:#555;font-size:14px}
.controls input,.controls select{padding:10px 12px;border-radius:8px;border:2px solid #e8eef7;background:#fff;font-size:14px;transition:all .3s}
.controls input:focus,.controls select:focus{outline:none;border-color:#667eea;box-shadow:0 0 0 3px rgba(102,126,234,0.1)}
.controls button{background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:10px 20px;border-radius:8px;cursor:pointer;font-weight:600;transition:transform .2s,box-shadow .2s}
.controls button:hover{transform:translateY(-2px);box-shadow:0 4px 15px rgba(102,126,234,0.4)}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(260px,1fr));gap:24px;margin-bottom:32px}
.card{background:#fff;border-radius:12px;overflow:hidden;box-shadow:0 4px 15px rgba(0,0,0,0.08);transition:transform .3s,box-shadow .3s;cursor:pointer;display:flex;flex-direction:column}
.card:hover{transform:translateY(-8px);box-shadow:0 12px 35px rgba(0,0,0,0.15)}
.card-img-wrap{position:relative;width:100%;height:200px;background:#f8f9fa;overflow:hidden}
.card-img{width:100%;height:100%;object-fit:cover;transition:transform .4s}
.card:hover .card-img{transform:scale(1.08)}
.card-img-placeholder{width:100%;height:100%;display:flex;align-items:center;justify-content:center;color:#bdc3c7;font-size:48px}
.card-content{padding:16px;flex:1;display:flex;flex-direction:column}
.card-title{font-size:18px;font-weight:700;margin-bottom:8px;color:#2c3e50;line-height:1.3}
.card-desc{font-size:13px;color:#7f8c8d;margin-bottom:10px;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}
.card-specs{display:grid;grid-template-columns:1fr 1fr;gap:8px;margin-bottom:12px;font-size:12px;color:#555}
.card-specs span{background:#f8f9fa;padding:6px 8px;border-radius:6px}
.card-price{font-size:22px;font-weight:700;color:#667eea;margin-bottom:12px}
.card-footer{margin-top:auto}
.card-btn{width:100%;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:10px;border-radius:8px;cursor:pointer;font-weight:600;transition:opacity .2s}
.card-btn:hover{opacity:0.9}
.card-btn:disabled{background:#94a3b8;cursor:not-allowed;opacity:.85}
.card-btn:disabled:hover{opacity:.85}
.pager{display:flex;gap:8px;align-items:center;justify-content:center;margin-top:32px;flex-wrap:wrap}
/* flying cart and small animation */
.fly-cart{position:fixed;width:50px;height:50px;border-radius:50%;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,#f093fb 0%,#f5576c 100%);color:#fff;font-size:22px;z-index:2000;transition:left .6s cubic-bezier(.2,.9,.2,1), top .6s cubic-bezier(.2,.9,.2,1), transform .6s cubic-bezier(.2,.9,.2,1), opacity .6s ease;box-shadow:0 4px 15px rgba(245,87,108,0.4)}
.bump{animation:bump .38s ease}
@keyframes bump{0%{transform:scale(1)}50%{transform:scale(1.15)}100%{transform:scale(1)}}

/* cart button one-time shake */
.cart-btn{background:linear-gradient(135deg,#f093fb 0%,#f5576c 100%);color:#fff;border:none;padding:10px 16px;border-radius:8px;cursor:pointer;font-weight:600;transition:transform .2s,box-shadow .2s;box-shadow:0 4px 15px rgba(245,87,108,0.2)}
.cart-btn:hover{transform:translateY(-2px);box-shadow:0 6px 20px rgba(245,87,108,0.35)}
.cart-btn.shake{animation:shake .6s ease}
@keyframes shake{0%{transform:translateX(0)}20%{transform:translateX(-6px)}40%{transform:translateX(6px)}60%{transform:translateX(-4px)}80%{transform:translateX(4px)}100%{transform:translateX(0)}}

.pager button{background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:10px 16px;border-radius:8px;cursor:pointer;font-weight:600;transition:transform .2s,box-shadow .2s}
.pager button:hover:not(:disabled){transform:translateY(-2px);box-shadow:0 4px 15px rgba(102,126,234,0.4)}
.pager button:disabled{opacity:0.5;cursor:not-allowed}
.pager #pageInfo{font-weight:500;color:#2c3e50;min-width:150px;text-align:center}
.admin-secret{position:fixed;left:16px;bottom:16px;z-index:1200;display:none;align-items:center;justify-content:center;padding:10px 14px;border:none;border-radius:10px;background:linear-gradient(135deg,#0f766e 0%,#0ea5a1 100%);color:#fff;font-size:13px;font-weight:700;cursor:pointer;box-shadow:0 8px 22px rgba(15,118,110,.28);transition:transform .15s ease,box-shadow .15s ease}
.admin-secret:hover{transform:translateY(-1px);box-shadow:0 12px 26px rgba(15,118,110,.35)}

@media(max-width:800px){
  .grid{grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:16px}
  .controls{flex-direction:column;align-items:stretch}
  .controls input,.controls select,.controls button{width:100%}
}
</style>
</head>
<body>

<nav class="navbar">
  <div class="navbar-content">
    <div class="navbar-title">🛍️ Каталог</div>
    <div class="navbar-actions">
      <button id="btnLogoutNav" class="cart-btn" style="background:linear-gradient(135deg,#64748b 0%,#334155 100%);padding:8px 12px">Выйти</button>
      <button id="btnCart" class="cart-btn">🛒 Корзина</button>
    </div>
  </div>
</nav>
<button id="btnAdminHidden" class="admin-secret" type="button">Админ-панель</button>

<div class="container">
  <div class="header">
    <h1>Популярные товары</h1>
    <p>Выбирайте из широкого ассортимента качественных товаров</p>
  </div>

  <div class="controls">
    <input type="text" id="search" placeholder="🔍 Поиск по названию...">
    <select id="categorySelect">
      <option value="">Все категории</option>
      {{range .categories}}
        <option value="{{.ID}}">{{.Name}}</option>
      {{end}}
    </select>
    <input type="text" id="minPrice" placeholder="Мин. цена">
    <input type="text" id="maxPrice" placeholder="Макс. цена">
    <select id="sortSelect">
      <option value="id_asc">По умолчанию</option>
      <option value="price_asc">Цена ↑</option>
      <option value="price_desc">Цена ↓</option>
      <option value="name_asc">Название ↑</option>
      <option value="name_desc">Название ↓</option>
      <option value="newest">Сначала новые</option>
    </select>
    <button id="btnFilter">🔎 Применить фильтры</button>
  </div>

  <div id="productsArea" class="grid" aria-live="polite">
    {{range .products}}
    <div class="card">
      {{if .ImageURL}}
        <div class="card-img-wrap"><img class="card-img" src="{{.ImageURL}}" alt="{{.Name}}"></div>
      {{else}}
        <div class="card-img-wrap"><div class="card-img-placeholder">📦</div></div>
      {{end}}
      <div class="card-content">
        <a href="/product/{{.ID}}" style="text-decoration:none;color:inherit">
        <h3 class="card-title">{{.Name}}</h3>
        <p class="card-desc">{{.Description}}</p>
        <div class="card-specs">
          <span>💰 {{printf "%.0f" .Price}} ₽</span>
          <span>📦 {{.Stock}} шт</span>
          <span>🎨 {{.Color}}</span>
          <span>✓ {{.Condition}}</span>
        </div>
        </a>
        <div class="card-price">{{printf "%.2f" .Price}} ₽</div>
        <div class="card-footer">
          {{if le .Stock 0}}
            <button class="card-btn" type="button" disabled title="Товара нет в наличии">Товара нет в наличии</button>
          {{else}}
            <button class="card-btn" type="button" onclick="addToCart(event, {{.ID}})">Добавить в корзину</button>
          {{end}}
        </div>
      </div>
    </div>
    {{else}}
    <div style="grid-column:1/-1;text-align:center;padding:40px;color:#7f8c8d">
      <div style="font-size:48px;margin-bottom:16px">📭</div>
      <p>Товары не найдены</p>
    </div>
    {{end}}
  </div>

  <div class="pager">
    <button id="prev">← Назад</button>
    <div id="pageInfo"></div>
    <button id="next">Вперёд →</button>
  </div>
<div id="cartModal" style="display:none;position:fixed;right:20px;bottom:20px;width:340px;background:#fff;padding:20px;border-radius:12px;box-shadow:0 12px 35px rgba(0,0,0,0.15);max-height:80vh;overflow:hidden;z-index:1100">
  <button id="closeCart" style="float:right;background:none;border:none;font-size:20px;cursor:pointer;color:#999">✕</button>
  <h3 style="margin-bottom:16px;color:#2c3e50">🛒 Корзина</h3>
  <div id="cartContent" style="overflow-y:auto;max-height:calc(80vh - 210px);padding-right:4px"></div>
  <div id="cartTotalDiv" style="margin-top:10px;font-weight:700;color:#667eea;font-size:18px"></div>
  <div style="margin-top:12px;display:flex;gap:8px;flex-direction:column;border-top:1px solid #e8eef7;padding-top:12px">
    <button id="checkoutBtn" style="background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:12px;border-radius:8px;cursor:pointer;font-weight:600">Оформить заказ</button>
  </div>
</div>

<!-- Checkout Modal -->
<div id="checkoutModal" style="display:none;position:fixed;left:0;top:0;right:0;bottom:0;background:rgba(0,0,0,0.5);z-index:2000;align-items:center;justify-content:center;flex-direction:column">
  <div style="background:#fff;border-radius:14px;box-shadow:0 20px 60px rgba(0,0,0,0.3);width:420px;max-width:90%;max-height:90vh;overflow-y:auto;padding:24px;">
    <h2 style="margin-bottom:24px;color:#2c3e50;font-size:24px;text-align:center">📋 Оформление заказа</h2>
    <form id="checkoutForm" style="display:flex;flex-direction:column;gap:16px">
      <div>
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">👤 Ваше имя</label>
        <input id="coName" type="text" placeholder="Иван Петров" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box">
        <div id="coNameErr" style="color:#e74c3c;font-size:13px;margin-top:6px"></div>
      </div>
      <div>
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">✉️ Email</label>
        <input id="coEmail" type="email" placeholder="ivan@example.com" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box">
        <div id="coEmailErr" style="color:#e74c3c;font-size:13px;margin-top:6px"></div>
      </div>
      <div id="coAddressWrap">
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">🏠 Адрес доставки</label>
        <textarea id="coAddress" placeholder="Город, улица, дом, квартира" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;min-height:80px;resize:vertical;transition:all .3s;box-sizing:border-box;font-family:inherit"></textarea>
        <div id="coAddressErr" style="color:#e74c3c;font-size:13px;margin-top:6px"></div>
      </div>
      <div>
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">📞 Номер телефона</label>
        <input id="coPhoneFull" type="tel" placeholder="+7 (999) 123-45-67" style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box">
        <div id="coPhoneErr" style="color:#e74c3c;font-size:13px;margin-top:6px"></div>
      </div>
      <div style="border-top:2px solid #e8eef7;padding-top:16px">
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:12px">🚚 Способ доставки</label>
        <div style="display:flex;gap:16px;margin-bottom:12px">
          <label style="display:flex;align-items:center;gap:8px;cursor:pointer">
            <input type="radio" id="coDeliveryPickup" name="coDeliveryType" value="pickup" style="cursor:pointer">
            <span>📍 Самовывоз (бесплатно)</span>
          </label>
          <label style="display:flex;align-items:center;gap:8px;cursor:pointer">
            <input type="radio" id="coDeliveryCourier" name="coDeliveryType" value="courier" style="cursor:pointer">
            <span>🚗 Курьер (200₽ + 50₽/км)</span>
          </label>
        </div>
      </div>
      <div id="pickupSection" style="display:none">
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">📦 Пункт выдачи</label>
        <div id="pickupMap" style="width:100%;height:260px;border:2px solid #e8eef7;border-radius:10px;overflow:hidden;margin-bottom:10px"></div>
        <div id="pickupInfo" style="background:#f7f9ff;border:1px solid #dce6ff;border-radius:10px;padding:10px 12px;color:#334155;font-size:13px;line-height:1.4;margin-bottom:10px">
          Выберите пункт на карте, чтобы увидеть адрес, время работы и контакты.
        </div>
        <select id="coPickupPoint" style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;box-sizing:border-box">
          <option value="">Выберите пункт выдачи</option>
        </select>
      </div>
      <input id="coPickupPointId" type="hidden" value="0">
      <input id="pickupLat" type="hidden" value="0">
      <input id="pickupLng" type="hidden" value="0">
      </div>
      <div style="display:flex;gap:12px;margin-top:20px">
        <button type="button" id="coCancel" style="flex:1;background:#e8eef7;color:#667eea;border:none;padding:12px;border-radius:8px;font-weight:600;cursor:pointer;transition:all .2s">Отмена</button>
        <button type="submit" id="coSubmit" style="flex:1;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:12px;border-radius:8px;font-weight:600;cursor:pointer;transition:all .2s">✓ Подтвердить</button>
      </div>
    </form>
  </div>
</div>

<script>
let page = 1, per_page = 12, totalPages = 1;
async function fetchProducts() {
  const search = document.getElementById('search').value;
  const category = document.getElementById('categorySelect').value;
  const min = document.getElementById('minPrice').value;
  const max = document.getElementById('maxPrice').value;
  const sort = document.getElementById('sortSelect').value;
  var url = '/products?page=' + page + '&per_page=' + per_page + '&sort=' + encodeURIComponent(sort);
  if (search) url = url + '&search=' + encodeURIComponent(search);
  if (category) url = url + '&category=' + encodeURIComponent(category);
  if (min) url = url + '&min_price=' + encodeURIComponent(min);
  if (max) url = url + '&max_price=' + encodeURIComponent(max);
  const r = await fetch(url);
  const data = await r.json();
  const area = document.getElementById('productsArea');
  area.innerHTML = '';
  if (!data.items || data.items.length === 0) {
    area.innerHTML = '<div style="grid-column:1/-1;text-align:center;padding:40px;color:#7f8c8d"><div style="font-size:48px;margin-bottom:16px">📭</div><p>Товары не найдены</p></div>';
  } else {
		for (var i = 0; i < data.items.length; i++) {
      var p = data.items[i];
      var div = document.createElement('div');
      div.className = 'card';
      var imghtml = '';
      if (p.image_url) imghtml = '<div class="card-img-wrap"><img class="card-img" src="' + p.image_url + '" alt=""></div>';
      else imghtml = '<div class="card-img-wrap"><div class="card-img-placeholder">📦</div></div>';
			var specs = '<div class="card-specs"><span>💰 ' + Number(p.price).toFixed(0) + ' ₽</span><span>📦 ' + (p.stock||0) + ' шт</span><span>🎨 ' + escapeHtml(p.color || '') + '</span><span>✓ ' + escapeHtml(p.condition || '') + '</span></div>';
      var outOfStock = Number(p.stock || 0) <= 0;
      var buyBtn = outOfStock
        ? '<button class="card-btn" type="button" disabled title="Товара нет в наличии">Товара нет в наличии</button>'
        : '<button class="card-btn" type="button" onclick="addToCart(event, ' + p.id + ')">Добавить в корзину</button>';
			div.innerHTML = imghtml + '<div class="card-content"><a href="/product/' + p.id + '" style="text-decoration:none;color:inherit"><h3 class="card-title">' + escapeHtml(p.name) + '</h3>' +
										'<p class="card-desc">' + escapeHtml(p.description || '') + '</p>' + specs + '</a>' +
										'<div class="card-price">' + Number(p.price).toFixed(2) + ' ₽</div>' +
										'<div class="card-footer">' + buyBtn + '</div></div>'; 
      area.appendChild(div);
    }
  }
  document.getElementById('pageInfo').innerText = 'Страница ' + data.page + ' / ' + data.total_page + ' (всего: ' + data.total + ' товаров)';
  totalPages = data.total_page;
  document.getElementById('prev').disabled = page <= 1;
  document.getElementById('next').disabled = page >= totalPages;
}
function escapeHtml(s){ return String(s).replace(/[&<>"'\/]/g, function(c){ return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;','/':'&#x2F;'}[c]; });}

let pickupMap = null;
let pickupPointsCache = [];
let pickupMarkerById = {};
let activePickupMarker = null;
let selectedPickupPoint = null;

function normalizePickupPoints(payload){
  if(Array.isArray(payload)) return payload;
  if(payload && Array.isArray(payload.points)) return payload.points;
  return [];
}

function pickupInfoHtml(point){
  if(!point){
    return 'Выберите пункт на карте, чтобы увидеть адрес, время работы и контакты.';
  }
  const city = escapeHtml(point.city || 'Россия');
  const name = escapeHtml(point.name || '');
  const address = escapeHtml(point.address || 'Адрес не указан');
  const hours = escapeHtml(point.working_hours || 'не указано');
  const phone = escapeHtml(point.phone || 'не указан');
  const details = escapeHtml(point.details || '');
  return '<strong>' + name + '</strong><br>' +
    'Город: ' + city + '<br>' +
    'Адрес: ' + address + '<br>' +
    'Время работы: ' + hours + '<br>' +
    'Контакт: ' + phone +
    (details ? '<br>Дополнительно: ' + details : '');
}

const RUSSIA_BOUNDS = {
  minLat: 41.0,
  maxLat: 82.0,
  minLng: 19.0,
  maxLng: 180.0
};

function isPointInRussia(pt){
  const lat = Number(pt && pt.latitude);
  const lng = Number(pt && pt.longitude);
  if(!Number.isFinite(lat) || !Number.isFinite(lng)) return false;
  return lat >= RUSSIA_BOUNDS.minLat && lat <= RUSSIA_BOUNDS.maxLat &&
    lng >= RUSSIA_BOUNDS.minLng && lng <= RUSSIA_BOUNDS.maxLng;
}

function addRussianTileLayer(map){
  const yandexLayer = L.tileLayer(
    'https://core-renderer-tiles.maps.yandex.net/tiles?l=map&x={x}&y={y}&z={z}&scale=1&lang=ru_RU',
    {
      maxZoom: 18,
      attribution: '&copy; Яндекс Карты'
    }
  );
  const osmFallback = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 19,
    attribution: '&copy; OpenStreetMap contributors'
  });

  let fallbackActivated = false;
  yandexLayer.on('tileerror', function(){
    if(fallbackActivated) return;
    fallbackActivated = true;
    if(map.hasLayer(yandexLayer)) map.removeLayer(yandexLayer);
    osmFallback.addTo(map);
  });

  yandexLayer.addTo(map);
}

function ensurePickupMap(){
  if(pickupMap || typeof L === 'undefined') return;
  const bounds = L.latLngBounds(
    [RUSSIA_BOUNDS.minLat, RUSSIA_BOUNDS.minLng],
    [RUSSIA_BOUNDS.maxLat, RUSSIA_BOUNDS.maxLng]
  );
  pickupMap = L.map('pickupMap', {
    worldCopyJump: false,
    minZoom: 3,
    maxBounds: bounds,
    maxBoundsViscosity: 0.85
  }).setView([61.0, 100.0], 3);
  addRussianTileLayer(pickupMap);
}

function setPickupInfo(point){
  const info = document.getElementById('pickupInfo');
  if(info) info.innerHTML = pickupInfoHtml(point);
}

function selectPickupPoint(pointId, fromMap){
  const id = parseInt(pointId || 0, 10);
  const selected = pickupPointsCache.find(function(pt){ return Number(pt.id) === id; }) || null;
  selectedPickupPoint = selected;

  const select = document.getElementById('coPickupPoint');
  const hiddenId = document.getElementById('coPickupPointId');
  const latEl = document.getElementById('pickupLat');
  const lngEl = document.getElementById('pickupLng');

  if(select && String(select.value) !== String(id) && id > 0){
    select.value = String(id);
  }
  if(hiddenId) hiddenId.value = selected ? String(selected.id) : '0';
  if(latEl) latEl.value = selected ? String(selected.latitude || 0) : '0';
  if(lngEl) lngEl.value = selected ? String(selected.longitude || 0) : '0';
  setPickupInfo(selected);

  if(!pickupMap || !selected) return;
  if(activePickupMarker) activePickupMarker.setZIndexOffset(0);
  const marker = pickupMarkerById[selected.id];
  if(marker){
    marker.setZIndexOffset(1000);
    marker.openPopup();
    activePickupMarker = marker;
    if(!fromMap) pickupMap.panTo([selected.latitude, selected.longitude]);
  }
}

function renderPickupPoints(points){
  pickupPointsCache = points || [];
  selectedPickupPoint = null;
  const hiddenId = document.getElementById('coPickupPointId');
  const latEl = document.getElementById('pickupLat');
  const lngEl = document.getElementById('pickupLng');
  if(hiddenId) hiddenId.value = '0';
  if(latEl) latEl.value = '0';
  if(lngEl) lngEl.value = '0';
  const select = document.getElementById('coPickupPoint');
  if(select){
    select.innerHTML = '<option value="">Выберите пункт выдачи</option>';
    pickupPointsCache.forEach(function(pt){
      const opt = document.createElement('option');
      opt.value = String(pt.id);
      opt.textContent = (pt.name || 'Пункт выдачи') + ' - ' + (pt.address || '');
      select.appendChild(opt);
    });
  }
  setPickupInfo(null);

  if(typeof L === 'undefined'){
    const mapBox = document.getElementById('pickupMap');
    if(mapBox){
      mapBox.innerHTML = '<div style="padding:10px;color:#dc2626;font-size:13px">Карта недоступна. Выберите пункт из списка ниже.</div>';
    }
    return;
  }

  ensurePickupMap();
  if(!pickupMap) return;

  Object.keys(pickupMarkerById).forEach(function(key){
    pickupMap.removeLayer(pickupMarkerById[key]);
  });
  pickupMarkerById = {};
  activePickupMarker = null;

  const mapPoints = pickupPointsCache.filter(isPointInRussia);
  const pointsForMap = mapPoints.length > 0 ? mapPoints : pickupPointsCache;

  pointsForMap.forEach(function(pt){
    const lat = Number(pt.latitude);
    const lng = Number(pt.longitude);
    if(!lat || !lng) return;
    const popup = '<strong>' + escapeHtml(pt.name || 'Пункт выдачи') + '</strong><br>' +
      escapeHtml(pt.address || '') + '<br>' +
      'Время работы: ' + escapeHtml(pt.working_hours || 'не указано') + '<br>' +
      'Телефон: ' + escapeHtml(pt.phone || 'не указан');
    const marker = L.marker([lat, lng]).addTo(pickupMap).bindPopup(popup);
    marker.on('click', function(){ selectPickupPoint(pt.id, true); });
    pickupMarkerById[pt.id] = marker;
  });

  if(pointsForMap.length > 0){
    const bounds = L.latLngBounds(
      pointsForMap
        .filter(function(pt){ return Number(pt.latitude) && Number(pt.longitude); })
        .map(function(pt){ return [pt.latitude, pt.longitude]; })
    );
    if(bounds.isValid()){
      pickupMap.fitBounds(bounds.pad(0.12), {padding:[24,24], maxZoom: 10});
    }
  }
  setTimeout(function(){ if(pickupMap) pickupMap.invalidateSize(); }, 120);
}

async function loadPickupPoints(){
  try {
    const r = await fetch('/pickup-points', {credentials:'include'});
    const payload = await r.json();
    renderPickupPoints(normalizePickupPoints(payload));
  } catch(err) {
    console.error('Failed to load pickup points:', err);
    setPickupInfo(null);
  }
}

function togglePickupSection(show){
  const section = document.getElementById('pickupSection');
  const addressWrap = document.getElementById('coAddressWrap');
  const addressInput = document.getElementById('coAddress');
  const addressErr = document.getElementById('coAddressErr');

  if(section) section.style.display = show ? 'block' : 'none';
  if(addressWrap) addressWrap.style.display = show ? 'none' : 'block';
  if(addressInput) addressInput.required = !show;
  if(addressErr) addressErr.innerText = '';

  if(show){
    setTimeout(function(){ if(pickupMap) pickupMap.invalidateSize(); }, 120);
  }
}

document.getElementById('btnFilter').addEventListener('click', function(){ page=1; fetchProducts();});
document.getElementById('prev').addEventListener('click', function(){ if(page>1){page--; fetchProducts()} });
document.getElementById('next').addEventListener('click', function(){ if(page < totalPages){ page++; fetchProducts(); } });
async function setupHiddenAdminButton(){
  const btn = document.getElementById('btnAdminHidden');
  if(!btn) return;
  btn.addEventListener('click', function(){ location.href = '/admin'; });
  try{
    const r = await fetch('/auth/me', {credentials:'include'});
    if(!r.ok) return;
    const data = await r.json();
    if(data && data.authenticated && data.user && data.user.role === 'admin'){
      btn.style.display = 'inline-flex';
    }
  }catch(e){ console.warn(e); }
}
setupHiddenAdminButton();
fetchProducts();

async function addToCart(e, productId) {
  var btnEl = null;
  try { if (e && e.currentTarget) btnEl = e.currentTarget; else if (e && e.target) btnEl = e.target.closest('button') || e.target; } catch(er) { btnEl = null; }
  const rect = btnEl ? btnEl.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width:40, height:40};

  const cartBtn = document.getElementById('btnCart');
  const cartRect = (cartBtn && cartBtn.getBoundingClientRect) ? cartBtn.getBoundingClientRect() : {left: window.innerWidth-40, top: window.innerHeight-40, width:40, height:40};
  const toX = cartRect.left + (cartRect.width/2);
  const toY = cartRect.top + (cartRect.height/2);

  const el = document.createElement('div');
  el.className = 'fly-cart';
  el.innerText = '🛒';
  el.style.position = 'fixed';
  el.style.left = (rect.left + (rect.width/2) - 25) + 'px';
  el.style.top = (rect.top + (rect.height/2) - 25) + 'px';
  el.style.opacity = '1';
  document.body.appendChild(el);

  requestAnimationFrame(()=>{
    el.style.left = (toX - 25) + 'px';
    el.style.top = (toY - 25) + 'px';
    el.style.transform = 'scale(0.8)';
  });

  const addPromise = fetch('/cart/add', { method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({product_id: productId, quantity: 1}), credentials: 'include' });
  const loadPromise = addPromise.then(r=>{ if(!r.ok) throw new Error('fail'); return loadCart(); }).catch(()=>{});

  const cleanup = ()=>{ el.style.opacity = '0'; setTimeout(()=>{ if(el.parentNode) el.parentNode.removeChild(el); }, 200); };
  el.addEventListener('transitionend', function onEnd(ev){ if(ev.propertyName==='left' || ev.propertyName==='top'){ el.removeEventListener('transitionend', onEnd); cleanup(); if (cartBtn){ cartBtn.classList.add('shake'); cartBtn.addEventListener('animationend', function _rm(){ cartBtn.classList.remove('shake'); cartBtn.removeEventListener('animationend', _rm); }); } } });
  setTimeout(()=>{ cleanup(); if (cartBtn){ cartBtn.classList.add('shake'); setTimeout(()=>{ cartBtn.classList.remove('shake'); }, 700); } }, 900);

  try{ if (btnEl){ btnEl.classList.add('bump'); setTimeout(()=>{ btnEl.classList.remove('bump'); }, 700); } }catch(ex){}

  await addPromise; await loadPromise;
}

document.getElementById('btnCart').addEventListener('click', function(){ document.getElementById('cartModal').style.display=document.getElementById('cartModal').style.display==='none'?'block':'none'; loadCart(); });
document.getElementById('btnLogoutNav').addEventListener('click', async function(){
  await fetch('/auth/logout', {method:'POST', credentials:'include'});
  location.href = '/auth';
});
document.getElementById('closeCart').addEventListener('click', ()=>{ document.getElementById('cartModal').style.display='none';});

async function loadCart(){
  const r = await fetch('/cart', {credentials: 'include'});
  const j = await r.json();
  const cont = document.getElementById('cartContent');
  const totalDiv = document.getElementById('cartTotalDiv');
  cont.innerHTML = '';
  if (totalDiv) totalDiv.innerHTML = '';
  if (!j.items || j.items.length === 0) {
    cont.innerHTML = '<div style="text-align:center;color:#7f8c8d;padding:20px">Корзина пуста</div>';
    return;
  }
  for (let it of j.items) {
    const div = document.createElement('div');
    div.style.marginBottom='12px';
    div.style.paddingBottom='12px';
    div.style.borderBottom='1px solid #e8eef7';
    div.innerHTML = '<div style="font-weight:600;color:#2c3e50">'+escapeHtml(it.product.name)+'</div><div style="color:#7f8c8d;font-size:13px">Цена: '+(it.product.price*it.quantity).toFixed(2)+'  ₽</div><div style="margin-top:6px;font-size:13px">Кол-во: <button onclick="updateCartQty(event.target,'+it.product.id+','+(it.quantity-1)+')">−</button> <span style="margin:0 8px">'+it.quantity+'</span> <button onclick="updateCartQty(event.target,'+it.product.id+','+(it.quantity+1)+')">+</button></div>';
    cont.appendChild(div);
  }
  if (totalDiv) totalDiv.innerHTML = 'Итого: ' + j.total.toFixed(2) + ' ₽';
}

async function updateCartQty(btn, productId, newQty){
  if(newQty<0) return;
  await fetch('/cart/update',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({product_id:productId,quantity:newQty}),credentials:'include'});
  await loadCart();
}

document.getElementById('checkoutBtn').addEventListener('click', async ()=>{
  document.getElementById('checkoutModal').style.display = 'flex';
  document.getElementById('coName').focus();
  const deliveryType = document.querySelector('input[name="coDeliveryType"]:checked')?.value || 'courier';
  togglePickupSection(deliveryType === 'pickup');
  await loadPickupPoints();
});

// Handle delivery type radio buttons
document.getElementById('coDeliveryPickup').addEventListener('change', ()=>{
  togglePickupSection(true);
});
document.getElementById('coDeliveryCourier').addEventListener('change', ()=>{
  togglePickupSection(false);
});

// Handle pickup point selection
document.getElementById('coPickupPoint').addEventListener('change', function(){
  selectPickupPoint(this.value, false);
});

// Phone formatting helper: basic E.164-aware mask with friendly grouping
function formatPhoneInput(el){
  var v = el.value;
  if(!v) return;
  // keep leading +, remove other non-digits
  var hasPlus = v.charAt(0) === '+';
  var digits = v.replace(/[^0-9]/g, '');
  if(hasPlus) digits = digits; // already digits
  // If starts with Russian '7' or country +7 -> format +7 (XXX) XXX-XX-XX
  if(hasPlus && digits.startsWith('7')){
    var d = digits.substring(1); // remove leading 7 for formatting blocks
    var out = '+7';
    if(d.length>0) out += ' (' + d.substring(0,3);
    if(d.length>=3) out += ') ' + d.substring(3,6);
    if(d.length>6) out += '-' + d.substring(6,8);
    if(d.length>8) out += '-' + d.substring(8,10);
    el.value = out;
    return;
  }
  // If country code +1 format as +1 (XXX) XXX-XXXX
  if(hasPlus && digits.startsWith('1')){
    var d1 = digits.substring(1);
    var out1 = '+1';
    if(d1.length>0) out1 += ' (' + d1.substring(0,3);
    if(d1.length>=3) out1 += ') ' + d1.substring(3,6);
    if(d1.length>6) out1 += '-' + d1.substring(6,10);
    el.value = out1;
    return;
  }
  // Generic: put +CC then groups of 3
  if(hasPlus){
    // take first 1-3 digits as country code
    var cc = digits.substring(0,3);
    var rest = digits.substring(cc.length);
    // try smaller country code lengths if rest too long
    // simple approach: first 1-3 as code
    var outg = '+' + cc;
    if(rest.length>0) outg += ' ' + rest.replace(/(\d{3})(?=\d)/g, '$1 ').trim();
    el.value = outg;
    return;
  }
  // If no plus, just group digits
  el.value = digits.replace(/(\d{3})(?=\d)/g, '$1 ').trim();
}

// Attach formatting on input
var phoneEl = document.getElementById('coPhoneFull');
if(phoneEl){
  phoneEl.addEventListener('input', function(e){
    var pos = this.selectionStart;
    formatPhoneInput(this);
    // try to keep caret near end
    this.selectionStart = this.selectionEnd = this.value.length;
  });
}

document.getElementById('coCancel').addEventListener('click', ()=>{
  document.getElementById('checkoutModal').style.display = 'none';
});

// Handle delivery type radio buttons  
var pickupRadio = document.getElementById('coDeliveryPickup');
var courierRadio = document.getElementById('coDeliveryCourier');
if(pickupRadio) pickupRadio.addEventListener('change', ()=>{
  togglePickupSection(true);
});
if(courierRadio) courierRadio.addEventListener('change', ()=>{
  togglePickupSection(false);
});

// Handle pickup point selection - update coordinates
var coPickupPoint = document.getElementById('coPickupPoint');
if(coPickupPoint){
  coPickupPoint.addEventListener('change', function(){
    selectPickupPoint(this.value, false);
  });
}

document.getElementById('checkoutForm').addEventListener('submit', async (e)=>{
  e.preventDefault();
  const name = document.getElementById('coName').value.trim();
  const email = document.getElementById('coEmail').value.trim();
  const address = document.getElementById('coAddress').value.trim();
  const phone = document.getElementById('coPhoneFull').value.trim();
  const deliveryType = document.querySelector('input[name="coDeliveryType"]:checked')?.value;
  
  // Clear previous errors
  document.getElementById('coNameErr').innerText = '';
  document.getElementById('coEmailErr').innerText = '';
  document.getElementById('coAddressErr').innerText = '';
  document.getElementById('coPhoneErr').innerText = '';

  // Validate fields with inline errors
  let hasErr = false;
  if(!name){ document.getElementById('coNameErr').innerText = 'Введите имя'; hasErr = true; }
  const emailRe = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  if(!email || !emailRe.test(email)){ document.getElementById('coEmailErr').innerText = 'Введите корректный email'; hasErr = true; }
  if(deliveryType === 'courier' && !address){ document.getElementById('coAddressErr').innerText = 'Введите адрес'; hasErr = true; }
  // phone validation: must start with + and have 7-15 digits
  const digits = phone.replace(/[^0-9]/g, '');
  if(!(phone.startsWith('+') && digits.length >= 7 && digits.length <= 15)){
    document.getElementById('coPhoneErr').innerText = 'Введите телефон в международном формате, например +7 (999) 123-45-67'; hasErr = true;
  }
  const pickupPointId = parseInt((document.getElementById('coPickupPoint').value || document.getElementById('coPickupPointId').value || '0'), 10);
  if(!deliveryType){ alert('⚠️ Выберите способ доставки'); hasErr = true; }
  if(deliveryType === 'pickup' && !pickupPointId){ alert('⚠️ Выберите пункт выдачи'); hasErr = true; }
  if(hasErr) return;
  
  const btn = document.getElementById('coSubmit');
  btn.disabled = true;
  btn.innerText = '⏳ Обработка...';
  
  try {
    const r = await fetch('/order', {
      method:'POST',
      headers:{'Content-Type':'application/json'},
      body: JSON.stringify({
        name:name,
        email:email,
        address: deliveryType==='courier' ? address : '',
        phone:phone,
        delivery_type:deliveryType,
        pickup_point_id: deliveryType==='pickup' ? pickupPointId : 0,
        pickup_point: deliveryType==='pickup' && selectedPickupPoint ? (selectedPickupPoint.name + ' - ' + selectedPickupPoint.address) : '',
        delivery_lat: deliveryType==='courier' ? (parseFloat(document.getElementById('pickupLat').value||0) || 55.7558) : 0,
        delivery_lng: deliveryType==='courier' ? (parseFloat(document.getElementById('pickupLng').value||0) || 37.6223) : 0
      }),
      credentials:'include'
    });
    
    if (!r.ok) {
      const err = await r.text();
      alert('❌ Ошибка: ' + (err || 'Не удалось оформить заказ'));
      return;
    }
    
    const data = await r.json();
    
    // Show success message
    document.getElementById('checkoutForm').innerHTML = 
      '<div style="text-align:center;padding:20px">' +
      '<div style="font-size:48px;margin-bottom:16px">✅</div>' +
      '<h3 style="color:#2c3e50;margin-bottom:8px">Заказ принят!</h3>' +
      '<p style="color:#7f8c8d;margin-bottom:16px">ID заказа: <strong>' + data.order_id + '</strong></p>' +
      '<p style="color:#7f8c8d;font-size:14px">Данные для доставки отправлены на ' + email + '</p>' +
      '</div>';
    
    setTimeout(()=>{
      document.getElementById('checkoutModal').style.display = 'none';
      document.getElementById('checkoutForm').innerHTML = 
        '<div><label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">👤 Ваше имя</label><input id="coName" type="text" placeholder="Иван Петров" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box"></div><div><label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">✉️ Email</label><input id="coEmail" type="email" placeholder="ivan@example.com" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box"></div><div><label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">🏠 Адрес доставки</label><textarea id="coAddress" placeholder="Город, улица, дом, квартира" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;min-height:80px;resize:vertical;transition:all .3s;box-sizing:border-box;font-family:inherit"></textarea></div><div><label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">📞 Номер телефона</label><input id="coPhone" type="tel" placeholder="+7 (999) 123-45-67" style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box"></div><div style="display:flex;gap:12px;margin-top:20px"><button type="button" id="coCancel" style="flex:1;background:#e8eef7;color:#667eea;border:none;padding:12px;border-radius:8px;font-weight:600;cursor:pointer;transition:all .2s">Отмена</button><button type="submit" id="coSubmit" style="flex:1;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:12px;border-radius:8px;font-weight:600;cursor:pointer;transition:all .2s">✓ Подтвердить</button></div>';
      document.getElementById('coCancel').addEventListener('click', ()=>{ document.getElementById('checkoutModal').style.display = 'none'; });
      document.getElementById('checkoutForm').addEventListener('submit', arguments.callee);
      btn.disabled = false;
      btn.innerText = '✓ Подтвердить';
      location.href='/';
    }, 3000);
    
  } catch(err){ 
    console.error(err);
    alert('❌ Ошибка сети: ' + err.message);
    btn.disabled = false;
    btn.innerText = '✓ Подтвердить';
  }
});
</script>
</body>
</html>`

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
.preview img{width:100%;height:132px;object-fit:cover}
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
.card-image img{width:100%;height:100%;object-fit:cover}
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
  var meta = 'Цена: ' + (p.price || 0) + ' | Остаток: ' + (p.stock || 0);
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
async function uploadFile(file){
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
async function initAdmin(){
  await loadCategories();
  await loadProducts();
}
initAdmin();
</script>
</body>
</html>`
const productHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Товар</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:'Segoe UI',Tahoma,Geneva,Verdana,sans-serif;background:linear-gradient(135deg,#f5f7fa 0%,#c3cfe2 100%);min-height:100vh;color:#2c3e50;padding:20px 0}
.navbar{background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;padding:0 20px;position:sticky;top:0;z-index:100;box-shadow:0 4px 15px rgba(0,0,0,0.15)}
.navbar-content{max-width:1100px;margin:0 auto;display:flex;justify-content:space-between;align-items:center;padding:16px 0}
.navbar-title{font-size:28px;font-weight:700;letter-spacing:1.2px}
.navbar-actions{display:flex;gap:16px;align-items:center}
.back-link{color:#fff;text-decoration:none;font-weight:500;display:flex;align-items:center;gap:8px;transition:opacity .3s;padding:8px 12px;border-radius:6px}
.back-link:hover{opacity:0.85;background:rgba(255,255,255,0.1)}
.container{max-width:1100px;margin:0 auto;padding:32px 20px}
.product-wrap{display:grid;grid-template-columns:1fr 1fr;gap:40px;margin-bottom:32px}
.product-left{border-radius:12px;overflow:hidden;background:#fff;box-shadow:0 4px 15px rgba(0,0,0,0.08)}
.product-left img{width:100%;height:auto;display:block;max-height:600px;object-fit:cover}
.product-left-placeholder{width:100%;height:480px;background:linear-gradient(135deg,#f5f7fa 0%,#c3cfe2 100%);display:flex;align-items:center;justify-content:center;font-size:80px;color:#e8eef7}
.product-right{background:#fff;padding:32px;border-radius:12px;box-shadow:0 4px 15px rgba(0,0,0,0.08);display:flex;flex-direction:column}
.product-title{font-size:32px;font-weight:700;margin-bottom:16px;color:#2c3e50}
.product-price{font-size:40px;font-weight:700;color:linear-gradient(135deg,#667eea 0%,#764ba2 100%);margin-bottom:24px;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);-webkit-background-clip:text;-webkit-text-fill-color:transparent;background-clip:text}
.product-desc{font-size:16px;color:#555;line-height:1.6;margin-bottom:24px}
.specs-grid{display:grid;grid-template-columns:1fr 1fr;gap:16px;margin-bottom:32px}
.spec-item{background:linear-gradient(135deg,#f5f7fa 0%,#e8eef7 100%);padding:16px;border-radius:10px;border:1px solid #e6eef9}
.spec-label{font-size:12px;color:#7f8c8d;text-transform:uppercase;font-weight:600;margin-bottom:6px}
.spec-value{font-size:16px;font-weight:700;color:#2c3e50}
.actions{display:flex;gap:12px;margin-top:auto}
.btn-add{flex:1;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:16px;border-radius:10px;font-size:16px;font-weight:700;cursor:pointer;transition:transform .2s,box-shadow .2s;box-shadow:0 4px 15px rgba(102,126,234,0.3)}
.btn-add:hover{transform:translateY(-3px);box-shadow:0 6px 25px rgba(102,126,234,0.4)}
.btn-add:disabled{opacity:0.9;cursor:not-allowed;transform:none;background:#94a3b8;box-shadow:none}
.flying-cart{position:fixed;width:50px;height:50px;border-radius:50%;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,#f093fb 0%,#f5576c 100%);color:#fff;font-size:22px;z-index:2000;box-shadow:0 4px 15px rgba(245,87,108,0.4)}
.cart-btn{background:linear-gradient(135deg,#f093fb 0%,#f5576c 100%);color:#fff;border:none;padding:10px 16px;border-radius:8px;cursor:pointer;font-weight:600;transition:transform .2s,box-shadow .2s;box-shadow:0 4px 15px rgba(245,87,108,0.2)}
.cart-btn:hover{transform:translateY(-2px);box-shadow:0 6px 20px rgba(245,87,108,0.35)}
.cart-panel{position:fixed;right:24px;bottom:24px;width:360px;max-height:80vh;background:#fff;border-radius:12px;box-shadow:0 12px 40px rgba(0,0,0,0.15);overflow:hidden;padding:20px;z-index:1100;max-width:calc(100% - 40px)}
.cart-content-scroll{overflow-y:auto;max-height:calc(80vh - 210px);padding-right:4px}
.cart-panel h3{margin-bottom:16px;font-size:20px;color:#2c3e50}
.cart-close{background:none;border:none;font-size:18px;cursor:pointer;float:right;color:#999}
.cart-item{padding:12px 0;border-bottom:1px solid #e8eef7;display:flex;gap:12px;align-items:flex-start}
.cart-item:last-child{border-bottom:none}
.cart-item-img{width:60px;height:60px;background:#f8f9fa;border-radius:8px;flex-shrink:0;display:flex;align-items:center;justify-content:center;overflow:hidden}
.cart-item-img img{width:100%;height:100%;object-fit:cover}
.cart-item-info{flex:1}
.cart-item-name{font-weight:600;color:#2c3e50;font-size:14px;margin-bottom:4px}
.cart-item-price{font-size:13px;color:#7f8c8d}
.cart-total{margin-top:16px;padding-top:16px;border-top:2px solid #e8eef7;font-size:18px;font-weight:700;color:#667eea}
.qty-controls{display:flex;gap:6px;align-items:center;margin-top:6px}
.qty-btn{background:#f8f9fa;border:1px solid #e8eef7;padding:4px 8px;border-radius:6px;cursor:pointer;font-weight:600;color:#667eea}
.qty-btn:hover{background:#e8eef7}
.admin-secret{position:fixed;left:16px;bottom:16px;z-index:1200;display:none;align-items:center;justify-content:center;padding:10px 14px;border:none;border-radius:10px;background:linear-gradient(135deg,#0f766e 0%,#0ea5a1 100%);color:#fff;font-size:13px;font-weight:700;cursor:pointer;box-shadow:0 8px 22px rgba(15,118,110,.28);transition:transform .15s ease,box-shadow .15s ease}
.admin-secret:hover{transform:translateY(-1px);box-shadow:0 12px 26px rgba(15,118,110,.35)}

@media(max-width:800px){
  .product-wrap{grid-template-columns:1fr;gap:20px}
  .product-right{padding:20px}
  .specs-grid{grid-template-columns:1fr}
  .navbar-content{flex-direction:column;gap:12px}
}
</style>
</head>
<body>

<nav class="navbar">
  <div class="navbar-content">
    <div class="navbar-title">🛍️ Каталог</div>
    <div class="navbar-actions">
      <a href="/" class="back-link">← Назад в каталог</a>
      <button id="btnLogoutNav" class="cart-btn" style="background:linear-gradient(135deg,#64748b 0%,#334155 100%);padding:8px 12px">Выйти</button>
      <button id="btnCart" class="cart-btn">🛒 Корзина</button>
    </div>
  </div>
</nav>
<button id="btnAdminHidden" class="admin-secret" type="button">Админ-панель</button>

<div class="container">
    <div class="product-wrap">
        <div class="product-left">
            {{if .product.ImageURL}}<img src="{{.product.ImageURL}}" alt="{{.product.Name}}">{{else}}<div class="product-left-placeholder">📦</div>{{end}}
        </div>
        <div class="product-right">
            <h1 class="product-title">{{.product.Name}}</h1>
            <div class="product-price">{{printf "%.2f" .product.Price}} ₽</div>
            <p class="product-desc">{{.product.Description}}</p>
            <div class="specs-grid">
                <div class="spec-item">
                    <div class="spec-label">📦 В наличии</div>
                    <div class="spec-value">{{.product.Stock}} шт</div>
                </div>
                <div class="spec-item">
                    <div class="spec-label">✓ Состояние</div>
                    <div class="spec-value">{{.product.Condition}}</div>
                </div>
                <div class="spec-item">
                    <div class="spec-label">🎨 Цвет</div>
                    <div class="spec-value">{{.product.Color}}</div>
                </div>
                <div class="spec-item">
                    <div class="spec-label">🪡 Материал</div>
                    <div class="spec-value">{{.product.Material}}</div>
                </div>
                <div class="spec-item">
                    <div class="spec-label">🌍 Страна</div>
                    <div class="spec-value">{{.product.Country}}</div>
                </div>
                <div class="spec-item">
                    <div class="spec-label">⭐ Категория</div>
                    <div class="spec-value">Товар</div>
                </div>
            </div>
            <div class="actions">
                <button id="btnAdd" class="btn-add" {{if le .product.Stock 0}}disabled{{end}}>{{if le .product.Stock 0}}Товара нет в наличии{{else}}Добавить в корзину{{end}}</button>
            </div>
        </div>
    </div>
</div>

<!-- Cart panel -->
<div id="cartPanel" class="cart-panel" style="display:none">
    <button id="cartClose" class="cart-close">✕</button>
    <h3>🛒 Корзина</h3>
    <div id="cartContent" class="cart-content-scroll"></div>
    <div class="cart-total" id="cartTotalDiv"></div>
    <button id="checkoutBtn" style="width:100%;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:12px;border-radius:8px;cursor:pointer;font-weight:600;margin-top:12px">Оформить заказ</button>
</div>

<!-- Checkout Modal -->
<div id="checkoutModal" style="display:none;position:fixed;left:0;top:0;right:0;bottom:0;background:rgba(0,0,0,0.5);z-index:2000;align-items:center;justify-content:center;flex-direction:column">
  <div style="background:#fff;border-radius:14px;box-shadow:0 20px 60px rgba(0,0,0,0.3);width:420px;max-width:90%;max-height:90vh;overflow-y:auto;padding:24px;">
    <h2 style="margin-bottom:24px;color:#2c3e50;font-size:24px;text-align:center">📋 Оформление заказа</h2>
    <form id="checkoutForm" style="display:flex;flex-direction:column;gap:16px">
      <div>
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">👤 Ваше имя</label>
        <input id="coName" type="text" placeholder="Иван Петров" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box">
      </div>
      <div>
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">✉️ Email</label>
        <input id="coEmail" type="email" placeholder="ivan@example.com" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box">
      </div>
      <div>
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">🏠 Адрес доставки</label>
        <textarea id="coAddress" placeholder="Город, улица, дом, квартира" required style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;min-height:80px;resize:vertical;transition:all .3s;box-sizing:border-box;font-family:inherit"></textarea>
      </div>
      <div>
        <label style="display:block;font-weight:600;color:#2c3e50;margin-bottom:8px">📞 Номер телефона</label>
          <input id="coPhoneFull" type="tel" placeholder="+7 (999) 123-45-67" style="width:100%;padding:12px;border:2px solid #e8eef7;border-radius:8px;font-size:14px;transition:all .3s;box-sizing:border-box">
      </div>
      <div style="display:flex;gap:12px;margin-top:20px">
        <button type="button" id="coCancel" style="flex:1;background:#e8eef7;color:#667eea;border:none;padding:12px;border-radius:8px;font-weight:600;cursor:pointer;transition:all .2s">Отмена</button>
        <button type="submit" id="coSubmit" style="flex:1;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);color:#fff;border:none;padding:12px;border-radius:8px;font-weight:600;cursor:pointer;transition:all .2s">✓ Подтвердить</button>
      </div>
    </form>
  </div>
</div>

<script>
const productId = {{.product.ID}};
const productStock = {{.product.Stock}};

function escapeHtml(s){ return String(s).replace(/[&<>"'\/]/g, function(c){ return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;','/':'&#x2F;'}[c]; });}
function productOutOfStock(){ return Number(productStock || 0) <= 0; }

function animateFly(fromRect, toX, toY){
	return new Promise((resolve)=>{
		const el = document.createElement('div');
		el.className = 'flying-cart';
		el.innerText = '🛒';
		el.style.position = 'fixed';
		el.style.zIndex = '2000';
		el.style.left = (fromRect.left + (fromRect.width/2) - 25) + 'px';
		el.style.top = (fromRect.top + (fromRect.height/2) - 25) + 'px';
		el.style.opacity = '1';
		el.style.transform = 'scale(1)';
		el.style.transition = 'left .6s cubic-bezier(.2,.9,.2,1), top .6s cubic-bezier(.2,.9,.2,1), transform .6s cubic-bezier(.2,.9,.2,1), opacity .6s ease';
		document.body.appendChild(el);

		requestAnimationFrame(()=>{
			el.style.left = (toX - 25) + 'px';
			el.style.top = (toY - 25) + 'px';
			el.style.transform = 'scale(0.8)';
		});

		const cleanup = ()=>{
			el.style.opacity = '0';
			setTimeout(()=>{ if(el.parentNode) el.parentNode.removeChild(el); resolve(); }, 200);
		};

		el.addEventListener('transitionend', function onEnd(e){
			if (e.propertyName === 'left' || e.propertyName === 'top'){
				el.removeEventListener('transitionend', onEnd);
				cleanup();
			}
		});
		setTimeout(()=>{ if(el.parentNode) { cleanup(); } }, 900);
	});
}

async function addToCart(e, productId){
	var btnEl = null;
	try { if (e && e.currentTarget) btnEl = e.currentTarget; else if (e && e.target) btnEl = e.target.closest('button') || e.target; } catch(er) { btnEl = null; }
	const rect = btnEl ? btnEl.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width:40, height:40};
	const cartBtn = document.getElementById('btnCart');
	const cartRect = (cartBtn && cartBtn.getBoundingClientRect) ? cartBtn.getBoundingClientRect() : {left: window.innerWidth-40, top: window.innerHeight-40, width:40, height:40};
	const toX = cartRect.left + (cartRect.width/2);
	const toY = cartRect.top + (cartRect.height/2);
	const anim = animateFly(rect, toX, toY);
	try{ await fetch('/cart/add', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({product_id: productId, quantity:1}), credentials:'include'}); await loadCart(); }catch(e){ console.warn(e); }
	await anim;
	if (cartBtn){ cartBtn.style.animation='shake .6s ease'; setTimeout(()=>{ cartBtn.style.animation=''; }, 600); }
}

async function loadCart(){
	const r = await fetch('/cart', {credentials: 'include'});
	const j = await r.json();
	const cont = document.getElementById('cartContent');
	cont.innerHTML = '';
	if (!j.items || j.items.length === 0) { cont.innerHTML = '<div style="text-align:center;color:#7f8c8d;padding:20px">Корзина пуста</div>'; document.getElementById('cartTotalDiv').innerHTML=''; return; }
	for (let it of j.items) {
		const img = it.product.image_url ? '<img src="'+it.product.image_url+'">' : '';
		const wrapper = document.createElement('div');
		wrapper.className = 'cart-item';
		wrapper.innerHTML = '<div class="cart-item-img">'+img+'</div><div class="cart-item-info"><div class="cart-item-name">'+escapeHtml(it.product.name)+'</div><div class="cart-item-price">'+Number(it.product.price).toFixed(2)+' ₽ × '+it.quantity+'</div><div class="qty-controls"><button class="qty-btn" onclick="updateQty('+it.product.id+','+(it.quantity-1)+')">−</button><span>'+it.quantity+'</span><button class="qty-btn" onclick="updateQty('+it.product.id+','+(it.quantity+1)+')">+</button></div></div>';
		cont.appendChild(wrapper);
	}
	document.getElementById('cartTotalDiv').innerHTML = 'Итого: <strong>'+j.total.toFixed(2)+' ₽</strong>';
	
	// update add button state
	try{
		const btn = document.getElementById('btnAdd');
		if(btn){
			const found = j.items.find(it => it.product && it.product.id === productId);
			const currentQty = found ? found.quantity : 0;
			btn.disabled = currentQty >= (productStock || 0);
		}
	}catch(e){console.warn(e)}
}

async function updateQty(productId, qty){
    if(qty<0) return;
    const r = await fetch('/cart/update', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({product_id: productId, quantity: qty}), credentials:'include'});
    if (!r.ok){ const txt = await r.text(); alert(txt || 'Ошибка'); return; }
    await loadCart();
}

document.getElementById('btnAdd').addEventListener('click', async function(e){
	if(productOutOfStock()) return;
	try{
		await addToCart(e, productId);
		const btn = e.currentTarget;
		btn.innerText = '✓ Добавлено в корзину';
		btn.style.background = 'linear-gradient(135deg,#6dd5ed 0%,#2193b0 100%)';
		setTimeout(()=>{ btn.innerText = 'Добавить в корзину'; btn.style.background = 'linear-gradient(135deg,#667eea 0%,#764ba2 100%)'; }, 2000);
	}catch(err){ console.warn(err); }
});

document.getElementById('btnCart').addEventListener('click', function(){ 
	const panel = document.getElementById('cartPanel');
	panel.style.display = panel.style.display === 'none' ? 'block' : 'none';
	if(panel.style.display !== 'none') loadCart();
});
document.getElementById('btnLogoutNav').addEventListener('click', async function(){
	await fetch('/auth/logout', {method:'POST', credentials:'include'});
	location.href = '/auth';
});
async function setupHiddenAdminButton(){
	const btn = document.getElementById('btnAdminHidden');
	if(!btn) return;
	btn.addEventListener('click', function(){ location.href = '/admin'; });
	try{
		const r = await fetch('/auth/me', {credentials:'include'});
		if(!r.ok) return;
		const data = await r.json();
		if(data && data.authenticated && data.user && data.user.role === 'admin'){
			btn.style.display = 'inline-flex';
		}
	}catch(e){ console.warn(e); }
}
setupHiddenAdminButton();

document.getElementById('cartClose').addEventListener('click', ()=>{ document.getElementById('cartPanel').style.display='none'; });

document.getElementById('checkoutBtn').addEventListener('click', async ()=>{
  document.getElementById('checkoutModal').style.display = 'flex';
  document.getElementById('coName').focus();
  
  // Load pickup points
  try {
    const r = await fetch('/pickup-points', {credentials: 'include'});
    const data = await r.json();
    const select = document.getElementById('coPickupPoint');
    select.innerHTML = '<option value="">Выберите пункт выдачи</option>';
    if(data && data.points && data.points.length > 0) {
      for(let pt of data.points) {
        const opt = document.createElement('option');
        opt.value = pt.id;
        opt.textContent = pt.name + ' - ' + pt.address;
        opt.dataset.lat = pt.latitude;
        opt.dataset.lng = pt.longitude;
        select.appendChild(opt);
      }
    }
  } catch(err) { console.error('Failed to load pickup points:', err); }
});

document.getElementById('coCancel').addEventListener('click', ()=>{
  document.getElementById('checkoutModal').style.display = 'none';
});

document.getElementById('checkoutForm').addEventListener('submit', async (e)=>{
  e.preventDefault();
  const name = document.getElementById('coName').value.trim();
  const email = document.getElementById('coEmail').value.trim();
  const address = document.getElementById('coAddress').value.trim();
  const phone = document.getElementById('coPhoneFull').value.trim();
  const deliveryType = document.getElementById('coDeliveryType').value;
  
  if(!name || !email || !deliveryType){
    alert('Пожалуйста, заполните все обязательные поля');
    return;
  }
  
  let pickup_point = '';
  let delivery_lat = 0;
  let delivery_lng = 0;
  
  if (deliveryType === 'pickup') {
    pickup_point = document.getElementById('coPickupPoint').value;
    if (!pickup_point) {
      alert('Пожалуйста, выберите точку самовывоза');
      return;
    }
  } else if (deliveryType === 'courier') {
    const latEl = document.getElementById('coDeliveryLat');
    const lngEl = document.getElementById('coDeliveryLng');
    const centerLat = 55.7558;
    const centerLng = 37.6223;
    delivery_lat = parseFloat((latEl && latEl.value) || 0) || centerLat;
    delivery_lng = parseFloat((lngEl && lngEl.value) || 0) || centerLng;
  } catch(err){ 
    console.error(err);
    alert('❌ Ошибка сети: ' + err.message);
    btn.disabled = false;
    btn.innerText = '✓ Подтвердить';
  }
});

loadCart();
if(productOutOfStock()){
	const btn = document.getElementById('btnAdd');
	if(btn){
		btn.disabled = true;
		btn.innerText = 'Товара нет в наличии';
	}
}
</script>

@keyframes shake{0%{transform:translateX(0)}20%{transform:translateX(-6px)}40%{transform:translateX(6px)}60%{transform:translateX(-4px)}80%{transform:translateX(4px)}100%{transform:translateX(0)}}

</body>
</html>`

const cartHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Корзина</title>
<style>
body{font-family:Inter,Arial,Helvetica,sans-serif;margin:20px;background:#f5f7fb;color:#222}
.container{max-width:1100px;margin:0 auto}
h1{margin-bottom:12px}
.list{background:#fff;padding:12px;border-radius:8px;border:1px solid #e6e6e6}
.item{display:flex;gap:12px;align-items:center;padding:8px;border-bottom:1px solid #f3f6f9}
.item:last-child{border-bottom:none}
.item img{width:78px;height:78px;object-fit:cover;border-radius:6px}
.total{margin-top:12px;font-weight:700}
.actions{display:flex;gap:8px;margin-top:12px}
.btn{background:#2b6cb0;color:#fff;border:none;padding:8px 12px;border-radius:8px;cursor:pointer}
.btn.secondary{background:#edf2f7;color:#2b6cb0}
</style>
</head>
<body>
<div class="container">
  <h1>Ваша корзина</h1>
  <div id="cartList" class="list">Загрузка...</div>
  <div class="total" id="cartTotal"></div>
  <div class="actions">
    <button id="goCheckout" class="btn">Оформить заказ</button>
    <button id="backToShop" class="btn secondary">Вернуться в магазин</button>
  </div>
</div>
<script>
async function renderCart(){
  const r = await fetch('/cart', {credentials: 'include'});
  const j = await r.json();
  const list = document.getElementById('cartList');
  const totalEl = document.getElementById('cartTotal');
  list.innerHTML = '';
  if (!j.items || j.items.length === 0){ list.innerHTML = '<div>Корзина пуста</div>'; totalEl.innerText=''; return; }
  for (const it of j.items){
    const div = document.createElement('div'); div.className='item';
    const img = it.product.image_url ? '<img src="'+it.product.image_url+'" alt="">' : '<div style="width:78px;height:78px;background:#f0f2f5;border-radius:6px"></div>';
    div.innerHTML = img + '<div style="flex:1"><div style="font-weight:600">'+(it.product.name||'')+'</div><div style="color:#666">Цена: '+(it.product.price).toFixed(2)+'</div></div><div style="text-align:right"><div>Кол-во: <button onclick="setQty('+it.product.id+','+(it.quantity-1)+')">−</button> <span>'+it.quantity+'</span> <button onclick="setQty('+it.product.id+','+(it.quantity+1)+')">+</button></div><div style="margin-top:6px">Сумма: '+(it.product.price*it.quantity).toFixed(2)+'</div></div>';
    list.appendChild(div);
  }
  totalEl.innerText = 'Итого: ' + j.total.toFixed(2) + ' ₽';
}
async function setQty(productId, qty){ if (qty<0) qty=0; await fetch('/cart/update',{method:'POST',headers:{'Content-Type':'application/json'},body: JSON.stringify({product_id:productId, quantity: qty}),credentials:'include'}); await renderCart(); }
document.getElementById('goCheckout').addEventListener('click', async function(){ const name = prompt('Ваше имя:'); if(!name) return; const email=prompt('Email:'); if(!email) return; const address=prompt('Адрес доставки:'); if(!address) return; const r = await fetch('/order',{method:'POST',headers:{'Content-Type':'application/json'}, body: JSON.stringify({name:name,email:email,address:address,phone:'',delivery_type:'pickup',pickup_point:'',delivery_lat:0,delivery_lng:0}),credentials:'include'}); if (!r.ok){ alert('Ошибка оформления'); return; } const data = await r.json(); alert('Заказ принят, id: '+data.order_id); location.href='/'; });
document.getElementById('backToShop').addEventListener('click', ()=>{ location.href='/' });
renderCart();
</script>
</body>
</html>`
