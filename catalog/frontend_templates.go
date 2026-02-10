package main

// ==================== TEMPLATES (inline) ====================
const indexHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Каталог</title>
<style>
body{font-family:Arial,Helvetica,sans-serif;margin:16px;background:#f3f6f9}
.container{max-width:1100px;margin:0 auto}
.header{display:flex;align-items:center;justify-content:space-between}
h1{margin:0}
.controls{display:flex;gap:8px;align-items:center;margin-top:12px}
select,input{padding:8px;border-radius:6px;border:1px solid #ddd}
button{padding:8px 12px;border-radius:6px;border:none;background:#2b6cb0;color:#fff;cursor:pointer}
.grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(220px,1fr));gap:12px;margin-top:16px}
.card{background:#fff;padding:12px;border-radius:8px;border:1px solid #e6e6e6}
.card img{width:100%;height:140px;object-fit:cover;border-radius:6px}
.meta{color:#666;font-size:13px;margin-top:8px;display:-webkit-box;-webkit-line-clamp:3;-webkit-box-orient:vertical;overflow:hidden;text-overflow:ellipsis;max-height:4.8em}
.pager{margin-top:12px;display:flex;gap:8px;align-items:center}
/* flying cart and small animation */
.fly-cart{position:fixed;width:44px;height:44px;border-radius:50%;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,#ffd166,#f08a5d);color:#fff;font-size:20px;z-index:2000;transition:left .6s cubic-bezier(.2,.9,.2,1), top .6s cubic-bezier(.2,.9,.2,1), transform .6s cubic-bezier(.2,.9,.2,1), opacity .6s ease}
.bump{animation:bump .38s ease}
@keyframes bump{0%{transform:scale(1)}50%{transform:scale(1.12)}100%{transform:scale(1)}}

/* cart button one-time shake */
.cart-btn.shake{animation:shake .6s ease}
@keyframes shake{0%{transform:translateX(0)}20%{transform:translateX(-6px)}40%{transform:translateX(6px)}60%{transform:translateX(-4px)}80%{transform:translateX(4px)}100%{transform:translateX(0)}}
</style>
</head>
<body>
<div class="container">
  <div class="header">
    <h1>Каталог товаров</h1>
    <div style="display:flex;gap:8px;align-items:center"><a href="/admin">Админ-панель</a><button id="btnCart" class="cart-btn" style="margin-left:12px">Корзина</button></div>
  </div>

  <div class="controls">
    <input type="text" id="search" placeholder="Поиск по названию">
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
    <button id="btnFilter">Применить</button>
  </div>

  <div id="productsArea" class="grid" aria-live="polite">
    {{range .products}}
    <div class="card">
      <a href="/product/{{.ID}}" style="text-decoration:none;color:inherit;display:block">
      {{if .ImageURL}}
        <img src="{{.ImageURL}}" alt="{{.Name}}">
      {{end}}
      <h3>{{.Name}}</h3>
      <div class="meta">{{.Description}}</div>
			<div class="meta">Цена: {{printf "%.2f" .Price}} ₽ • В наличии: {{.Stock}}</div>
			<div class="meta">Цвет: {{.Color}} • Состояние: {{.Condition}}</div>
			<div class="meta">Страна: {{.Country}} • Материал: {{.Material}}</div>
      </a>
      <div style="margin-top:8px">
        <button onclick="addToCart(event, {{.ID}})">Добавить в корзину</button>
      </div>
    </div>
    {{else}}
    <div class="card">Товаров нет</div>
    {{end}}
  </div>

  <div class="pager">
    <button id="prev">Назад</button>
    <div id="pageInfo"></div>
    <button id="next">Вперёд</button>
  </div>
</div>

<div id="cartModal" style="display:none;position:fixed;right:20px;bottom:20px;width:320px;background:#fff;padding:12px;border-radius:8px;box-shadow:0 8px 24px rgba(0,0,0,0.12)">
  <h3>Корзина</h3>
  <div id="cartContent"></div>
  <div style="margin-top:8px">
    <button id="checkoutBtn">Оформить заказ</button>
    <button id="closeCart">Закрыть</button>
  </div>
</div>

<script>
let page = 1, per_page = 6, totalPages = 1;
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
    area.innerHTML = '<div class="card">Товары не найдены</div>';
  } else {
			for (var i = 0; i < data.items.length; i++) {
      var p = data.items[i];
      var div = document.createElement('div');
      div.className = 'card';
      var imghtml = '';
      if (p.image_url) imghtml = '<img src="' + p.image_url + '" alt="">';
			div.innerHTML = '<a href="/product/' + p.id + '" style="text-decoration:none;color:inherit">' + imghtml + '<h3>' + escapeHtml(p.name) + '</h3>' +
										'<div class="meta">' + escapeHtml(p.description || '') + '</div>' +
										'<div class="meta">Цена: ' + Number(p.price).toFixed(2) + ' ₽ • В наличии: ' + (p.stock||0) + '</div>' +
										'<div class="meta">Цвет: ' + escapeHtml(p.color || '') + ' • Состояние: ' + escapeHtml(p.condition || '') + '</div>' +
										'<div class="meta">Страна: ' + escapeHtml(p.country || '') + ' • Материал: ' + escapeHtml(p.material || '') + '</div>' +
										'</a><div style="margin-top:8px"><button onclick="addToCart(event, ' + p.id + ')">Добавить в корзину</button></div>'; 
      area.appendChild(div);
    }
  }
  document.getElementById('pageInfo').innerText = 'Страница ' + data.page + ' / ' + data.total_page + ' (Всего: ' + data.total + ')';
  totalPages = data.total_page;
  document.getElementById('prev').disabled = page <= 1;
  document.getElementById('next').disabled = page >= totalPages;
}
function escapeHtml(s){ return String(s).replace(/[&<>"'\/]/g, function(c){ return {'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;','/':'&#x2F;'}[c]; });}
document.getElementById('btnFilter').addEventListener('click', function(){ page=1; fetchProducts();});
document.getElementById('prev').addEventListener('click', function(){ if(page>1){page--; fetchProducts()} });
document.getElementById('next').addEventListener('click', function(){ if(page < totalPages){ page++; fetchProducts(); } });
fetchProducts();

async function addToCart(e, productId) {
  // determine source rect (button that was clicked)
  var btnEl = null;
  try { if (e && e.currentTarget) btnEl = e.currentTarget; else if (e && e.target) btnEl = e.target.closest('button') || e.target; } catch(er) { btnEl = null; }
  const rect = btnEl ? btnEl.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width:40, height:40};

  // target: the cart button
  const cartBtn = document.getElementById('btnCart');
  const cartRect = (cartBtn && cartBtn.getBoundingClientRect) ? cartBtn.getBoundingClientRect() : {left: window.innerWidth-40, top: window.innerHeight-40, width:40, height:40};
  const toX = cartRect.left + (cartRect.width/2);
  const toY = cartRect.top + (cartRect.height/2);

  // create fly element
  const el = document.createElement('div');
  el.className = 'fly-cart';
  el.innerText = '🛒';
  el.style.position = 'fixed';
  el.style.left = (rect.left + (rect.width/2) - 22) + 'px';
  el.style.top = (rect.top + (rect.height/2) - 22) + 'px';
  el.style.opacity = '1';
  document.body.appendChild(el);

  // animate to cart button
  requestAnimationFrame(()=>{
    el.style.left = (toX - 22) + 'px';
    el.style.top = (toY - 22) + 'px';
    el.style.transform = 'scale(0.8)';
  });

  // add to cart and refresh cart contents in background
  const addPromise = fetch('/cart/add', { method: 'POST', headers: {'Content-Type':'application/json'}, body: JSON.stringify({product_id: productId, quantity: 1}) });
  const loadPromise = addPromise.then(r=>{ if(!r.ok) throw new Error('fail'); return loadCart(); }).catch(()=>{});

  // on animation end -> cleanup and shake cart btn once
  const cleanup = ()=>{ el.style.opacity = '0'; setTimeout(()=>{ if(el.parentNode) el.parentNode.removeChild(el); }, 200); };
  el.addEventListener('transitionend', function onEnd(ev){ if(ev.propertyName==='left' || ev.propertyName==='top'){ el.removeEventListener('transitionend', onEnd); cleanup(); if (cartBtn){ cartBtn.classList.add('shake'); cartBtn.addEventListener('animationend', function _rm(){ cartBtn.classList.remove('shake'); cartBtn.removeEventListener('animationend', _rm); }); } } });
  setTimeout(()=>{ cleanup(); if (cartBtn){ cartBtn.classList.add('shake'); setTimeout(()=>{ cartBtn.classList.remove('shake'); }, 700); } }, 900);

  // small feedback on source button
  try{ if (btnEl){ btnEl.classList.add('bump'); setTimeout(()=>{ btnEl.classList.remove('bump'); }, 700); } }catch(ex){}

  await addPromise; await loadPromise;
} 


document.getElementById('btnCart').addEventListener('click', function(){ location.href = '/cart/view'; });
document.getElementById('closeCart').addEventListener('click', ()=>{ document.getElementById('cartModal').style.display='none';});

async function loadCart(){
  const r = await fetch('/cart');
  const j = await r.json();
  const cont = document.getElementById('cartContent');
  cont.innerHTML = '';
  if (!j.items || j.items.length === 0) { cont.innerHTML = '<div>Корзина пуста</div>'; return; }
  for (let it of j.items) {
    const div = document.createElement('div');
    div.style.marginBottom='8px';
    div.innerHTML = '<b>'+escapeHtml(it.product.name)+'</b><div>Кол-во: '+it.quantity+' • Цена: '+(it.product.price*it.quantity).toFixed(2)+'</div><div style="margin-top:4px"><button onclick="removeFromCart('+it.product.id+')">Удалить</button></div>';
    cont.appendChild(div);
  }
  cont.insertAdjacentHTML('beforeend','<div style="margin-top:8px"><b>Итого: '+j.total.toFixed(2)+'</b></div>');
}

async function removeFromCart(productId){
  await fetch('/cart/remove', {
    method:'POST',
    headers: {'Content-Type':'application/json'},
    body: JSON.stringify({product_id: productId})
  });
  await loadCart();
}

document.getElementById('checkoutBtn').addEventListener('click', async ()=>{
  const name = prompt('Ваше имя:');
  if(!name) return;
  const email = prompt('Email:');
  if(!email) return;
  const address = prompt('Адрес доставки:');
  if(!address) return;
  const r = await fetch('/order', {
    method:'POST',
    headers:{'Content-Type':'application/json'},
    body: JSON.stringify({name:name,email:email,address:address})
  });
  if (!r.ok) {
    alert('Ошибка оформления');
    return;
  }
  const data = await r.json();
  alert('Заказ принят, id: ' + data.order_id);
  document.getElementById('cartModal').style.display='none';
  fetchProducts();
});
</script>
</body>
</html>`

const adminHTML = `<!doctype html>
<html lang="ru"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Admin</title>
<style>
  body{font-family:Inter,Arial,Helvetica,sans-serif;margin:16px;background:#f7fafc;color:#0b374f}
  .container{max-width:1100px;margin:0 auto}
  h1{display:flex;justify-content:space-between;align-items:center}
  input,select,textarea{padding:8px;border-radius:8px;border:1px solid #e6eef9;background:#fff}
  button{background:#2b6cb0;color:#fff;border:none;padding:8px 10px;border-radius:8px;cursor:pointer}
  .btn.secondary{background:#edf2f7;color:#2b6cb0}
  .grid{display:grid;grid-template-columns:1fr 1fr;gap:12px}
  .card{background:#fff;padding:12px;border-radius:10px;border:1px solid #e6e6e6}
  .card:hover{box-shadow:0 8px 24px rgba(11,20,40,0.06);transform:translateY(-3px)}
  .small{font-size:13px;color:#5a6b7a}

  /* File input beautify */
  .file-input{display:inline-flex;align-items:center;gap:8px;padding:8px 10px;border-radius:8px;background:#fff;border:1px dashed #cfe3f6;color:#0b4a77;cursor:pointer}
  .file-input input{display:none}
  .img-preview{width:100%;min-height:120px;border-radius:8px;border:1px dashed #e6eef9;background:#fbfdff;display:flex;align-items:center;justify-content:center;overflow:hidden}
  .img-thumb{width:100%;height:120px;display:flex;align-items:center;justify-content:center}
  .img-thumb img{max-width:100%;max-height:100%;object-fit:cover;border-radius:8px}
  .remove-btn{background:transparent;border:none;color:#ef4444;cursor:pointer;padding:6px 8px;border-radius:6px}

  /* Modal nicer */
  #editModal > div{box-shadow:0 18px 50px rgba(11,20,40,0.12)}
  #editModal .btn{min-width:110px}

  @media(max-width:800px){ .grid{grid-template-columns:1fr} }
</style>
</head><body><div class="container"><h1>Admin панель <button id="btnLogout">Выйти</button></h1>
<section><h2>Создать категорию</h2><div><input id="catName" placeholder="Название"><button id="btnCreateCat">Создать</button></div></section>
<section><h2>Создать товар</h2>
<div style="display:flex;gap:12px;flex-wrap:wrap;align-items:flex-start">
  <div style="flex:1;min-width:320px">
    <input id="pName" placeholder="Название" style="width:100%;margin-bottom:8px">
    <input id="pPrice" placeholder="Цена" style="width:48%;margin-right:4%;margin-bottom:8px">
    <input id="pStock" placeholder="Кол-во" style="width:48%;margin-bottom:8px">
    <select id="pCategory" style="width:100%;margin-bottom:8px"><option value="">Категория</option></select>
    <div style="margin-bottom:8px">
      <label class="file-input">
        <input id="pImage" type="file" accept="image/*">
        <span>Выбрать изображение</span>
      </label>
      <button type="button" id="pImageRemove" class="remove-btn" style="display:none;margin-left:8px">Удалить</button>
    </div>
    <input id="pColor" placeholder="Цвет" style="width:48%;margin-right:4%;margin-bottom:8px">
    <input id="pCondition" placeholder="Состояние (новый/бу)" style="width:48%;margin-bottom:8px">
    <input id="pCountry" placeholder="Страна" style="width:48%;margin-right:4%;margin-bottom:8px">
    <input id="pMaterial" placeholder="Материал" style="width:48%;margin-bottom:8px">
    <textarea id="pDesc" placeholder="Описание" style="width:100%;min-height:80px;margin-top:8px"></textarea>
  </div>
  <div style="width:220px">
    <div id="pImagePreview" class="img-preview">Нет изображения</div>
  </div>
</div>
<div style="margin-top:10px;text-align:right"><button id="btnCreateProd">Создать товар</button></div>
</section>
<section><h2>Товары</h2><div id="prodList" class="grid"></div></section></div>

<!-- Edit product modal -->
<div id="editModal" style="display:none;position:fixed;left:0;top:0;right:0;bottom:0;background:rgba(0,0,0,0.4);align-items:center;justify-content:center;z-index:2000">
	<div style="width:760px;max-width:96%;margin:0 auto;background:#fff;border-radius:10px;padding:16px;position:relative;">
		<button id="editClose" style="position:absolute;right:12px;top:12px;border:none;background:transparent;font-size:20px;cursor:pointer">✕</button>
		<h2>Редактировать товар</h2>
		<div style="display:flex;gap:12px;align-items:flex-start;flex-wrap:wrap">
			<div style="flex:1 1 320px">
				<input id="eName" placeholder="Название" style="width:100%;margin-bottom:8px">
				<input id="ePrice" placeholder="Цена" style="width:100%;margin-bottom:8px">
				<input id="eStock" placeholder="Кол-во" style="width:100%;margin-bottom:8px">
				<select id="eCategory" style="width:100%;margin-bottom:8px"><option value="">Категория</option></select>
				<label class="file-input" style="display:inline-block"><input id="eImage" type="file" accept="image/*"><span>Сменить изображение</span></label>
			</div>
			<div style="flex:1 1 320px">
				<div id="eImagePreview" class="img-thumb" style="margin-bottom:8px">Нет изображения</div>
				<div class="img-actions"><button type="button" id="eImageReplaceBtn" class="primary-btn">Загрузить новую</button> <button type="button" id="eImageRemoveBtn" class="remove-btn" style="display:none">Удалить</button></div>
			<input id="eColor" placeholder="Цвет" style="width:100%;margin-bottom:8px">
				<input id="eCondition" placeholder="Состояние" style="width:100%;margin-bottom:8px">
				<input id="eCountry" placeholder="Страна" style="width:100%;margin-bottom:8px">
				<input id="eMaterial" placeholder="Материал" style="width:100%;margin-bottom:8px">
			</div>
		</div>
		<textarea id="eDesc" placeholder="Описание" style="width:100%;height:80px;margin-top:8px"></textarea>
		<div style="margin-top:12px;text-align:right">
			<button id="saveEdit" class="btn">Сохранить</button>
			<button id="cancelEdit" class="btn secondary" style="margin-left:8px">Отмена</button>
		</div>
	</div>
</div>

<script>
async function api(path, method, body, isJson){
  if (!method) method='GET'; if (isJson===undefined) isJson=true;
  var opts={method:method, headers:{}};
  if (isJson && body) { opts.headers['Content-Type']='application/json'; opts.body = JSON.stringify(body); }
  if (!isJson && body) { opts.body = body; }
  var r = await fetch(path, opts);
  if (!r.ok) { var txt = await r.text(); throw new Error(txt); }
  return r.json();
}
async function loadCategories(){
  var cats = await api('/categories','GET'); var sel = document.getElementById('pCategory'); sel.innerHTML = '<option value=\"\">Категория</option>';
  for(var i=0;i<cats.length;i++){ var c=cats[i]; sel.insertAdjacentHTML('beforeend','<option value=\"'+c.id+'\">'+c.name+'</option>'); }
  return cats;
}
async function loadProducts(){
  var data = await api('/products?page=1&per_page=100','GET'); var items = data.items; var list = document.getElementById('prodList'); list.innerHTML='';
  for(var i=0;i<items.length;i++){ var p=items[i]; var img = p.image_url?'<img src=\"'+p.image_url+'\" style=\"width:100%;height:120px;object-fit:cover;border-radius:6px\">':''; var html='<div class=\"card\" onclick=\"editProduct('+p.id+')\" style=\"cursor:pointer\">'+img+'<h3>'+p.name+'</h3><div class=\"small\">Цена: '+p.price+' • В наличии: '+p.stock+'</div><div class=\"small\">'+(p.description||'')+'</div><div style=\"margin-top:8px;text-align:right\"><button onclick=\"event.stopPropagation(); editProduct('+p.id+')\">Редактировать</button> <button onclick=\"event.stopPropagation(); delProduct('+p.id+')\">Удалить</button></div></div>'; list.insertAdjacentHTML('beforeend',html); }
}
document.getElementById('btnCreateCat').addEventListener('click', async function(){ var name=document.getElementById('catName').value; await api('/admin/api/categories','POST',{name:name}); document.getElementById('catName').value=''; await loadCategories(); await loadProducts(); });
async function uploadFile(file){ var fd = new FormData(); fd.append('file', file); var r = await fetch('/admin/api/upload', {method:'POST', body: fd}); if(!r.ok) throw new Error('upload failed'); return (await r.json()).url; }

// Create product: preview and remove handlers
var pImageInput = document.getElementById('pImage');
var pImagePreview = document.getElementById('pImagePreview');
var pImageRemoveBtn = document.getElementById('pImageRemove');
if (pImageInput) {
  pImageInput.addEventListener('change', function(){
    var f = this.files[0];
    if (f){
      var r = new FileReader(); r.onload = function(e){ pImagePreview.innerHTML = '<div class="img-thumb"><img src="'+e.target.result+'"></div>'; pImageRemoveBtn.style.display='inline-block'; }
      r.readAsDataURL(f);
    } else { pImagePreview.innerHTML = 'Нет изображения'; pImageRemoveBtn.style.display='none'; }
  });
}
if (pImageRemoveBtn){ pImageRemoveBtn.addEventListener('click', function(){ pImageInput.value=''; pImagePreview.innerHTML = 'Нет изображения'; this.style.display='none'; }); }

document.getElementById('btnCreateProd').addEventListener('click', async function(){ try{ 
		var name=document.getElementById('pName').value; 
		var price=parseFloat(document.getElementById('pPrice').value||0); 
		var stock=parseInt(document.getElementById('pStock').value||0); 
		var cat=document.getElementById('pCategory').value||null; 
		var desc=document.getElementById('pDesc').value; 
		var color=document.getElementById('pColor').value||'';
		var condition=document.getElementById('pCondition').value||'';
		var country=document.getElementById('pCountry').value||'';
		var material=document.getElementById('pMaterial').value||'';
		var file=document.getElementById('pImage').files[0]; 
		var image_url=null; if(file) image_url=await uploadFile(file);
		await api('/admin/api/products','POST',{
			name:name,price:price,stock:stock,category_id:cat?parseInt(cat):null,description:desc,image_url:image_url,
			color: color, condition: condition, country: country, material: material
		}); 
		document.getElementById('pName').value=''; document.getElementById('pPrice').value=''; document.getElementById('pStock').value=''; document.getElementById('pDesc').value=''; document.getElementById('pImage').value=''; document.getElementById('pColor').value=''; document.getElementById('pCondition').value=''; document.getElementById('pCountry').value=''; document.getElementById('pMaterial').value='';
		pImagePreview.innerHTML = 'Нет изображения'; pImageRemoveBtn.style.display='none';
		await loadProducts(); }catch(e){ alert(e.message); } });

// Edit modal: enhanced preview/remove logic
var eRemoveImage = false;
var eSelectedFile = null;
var eImageInput = document.getElementById('eImage');
var eImagePreview = document.getElementById('eImagePreview');
var eImageRemoveBtn = document.getElementById('eImageRemoveBtn');
if (eImageInput){
  eImageInput.addEventListener('change', function(){
    var f = this.files[0];
    eSelectedFile = f || null;
    eRemoveImage = false;
    if (f){ var r = new FileReader(); r.onload = function(e){ eImagePreview.innerHTML = '<div class="img-thumb"><img src="'+e.target.result+'"></div>'; eImageRemoveBtn.style.display='inline-block'; } ; r.readAsDataURL(f); }
    else { eImagePreview.innerHTML = ''; eImageRemoveBtn.style.display='none'; }
  });
}
if (eImageRemoveBtn){ eImageRemoveBtn.addEventListener('click', function(){ eRemoveImage = true; eSelectedFile = null; if (eImageInput) eImageInput.value=''; eImagePreview.innerHTML = '<div class="small">Изображение будет удалено</div>'; this.style.display='none'; }); }
var eImageReplaceBtn = document.getElementById('eImageReplaceBtn'); if (eImageReplaceBtn){ eImageReplaceBtn.addEventListener('click', function(){ if (eImageInput) eImageInput.click(); }); }

// Update editProduct() loader to use preview/remove state
async function editProduct(id){
	try{
		const prod = await api('/products/'+id,'GET');
		await loadCategories();
		const src = document.getElementById('pCategory');
		const dest = document.getElementById('eCategory');
		dest.innerHTML = src.innerHTML;
		document.getElementById('eName').value = prod.name || '';
		document.getElementById('ePrice').value = prod.price || '';
		document.getElementById('eStock').value = prod.stock || '';
		document.getElementById('eCategory').value = prod.category_id || '';
		document.getElementById('eColor').value = prod.color || '';
		document.getElementById('eCondition').value = prod.condition || '';
		document.getElementById('eCountry').value = prod.country || '';
		document.getElementById('eMaterial').value = prod.material || '';
		document.getElementById('eDesc').value = prod.description || '';
		document.getElementById('eImage').value = '';
		// preview and remove
		eRemoveImage = false; eSelectedFile = null;
		if (prod.image_url) {
			eImagePreview.innerHTML = '<div class="img-thumb"><img src="'+prod.image_url+'" style="display:block;margin-bottom:6px"></div>';
			eImageRemoveBtn.style.display = 'inline-block';
		} else {
			eImagePreview.innerHTML = '';
			eImageRemoveBtn.style.display = 'none';
		}
		document.getElementById('editModal').dataset.editing = id;
		document.getElementById('editModal').style.display = 'flex';
	}catch(e){ alert('Ошибка загрузки товара: '+e.message); }
}

// Save edit: respect selected file or explicit removal
document.getElementById('saveEdit').addEventListener('click', async function(){
	try{
		const id = document.getElementById('editModal').dataset.editing;
		if(!id) return alert('Нет товара для сохранения');
		const payload = {};
		payload.name = document.getElementById('eName').value;
		payload.price = parseFloat(document.getElementById('ePrice').value||0);
		payload.stock = parseInt(document.getElementById('eStock').value||0);
		const cat = document.getElementById('eCategory').value; payload.category_id = cat?parseInt(cat):null;
		payload.color = document.getElementById('eColor').value||'';
		payload.condition = document.getElementById('eCondition').value||'';
		payload.country = document.getElementById('eCountry').value||'';
		payload.material = document.getElementById('eMaterial').value||'';
		payload.description = document.getElementById('eDesc').value||'';
		if (eSelectedFile){
			const url = await uploadFile(eSelectedFile);
			payload.image_url = url;
		} else if (eRemoveImage){
			payload.image_url = null;
		}
		await api('/admin/api/products/'+id,'PUT',payload);
		closeEditModal();
		await loadProducts();
	}catch(e){ alert('Ошибка сохранения: '+e.message); }
});
async function delProduct(id){ if(!confirm('Удалить?')) return; await api('/admin/api/products/'+id,'DELETE'); await loadProducts(); }

	document.getElementById('btnLogout').addEventListener('click', async function(){ await fetch('/admin/logout',{method:'POST'}); location.href='/'; });

async function delProduct(id){ if(!confirm('Удалить?')) return; await api('/admin/api/products/'+id,'DELETE'); await loadProducts(); }
async function editProduct(id){
	try{
		const prod = await api('/products/'+id,'GET');
		await loadCategories();
		const src = document.getElementById('pCategory');
		const dest = document.getElementById('eCategory');
		dest.innerHTML = src.innerHTML;
		document.getElementById('eName').value = prod.name || '';
		document.getElementById('ePrice').value = prod.price || '';
		document.getElementById('eStock').value = prod.stock || '';
		document.getElementById('eCategory').value = prod.category_id || '';
		document.getElementById('eColor').value = prod.color || '';
		document.getElementById('eCondition').value = prod.condition || '';
		document.getElementById('eCountry').value = prod.country || '';
		document.getElementById('eMaterial').value = prod.material || '';
		document.getElementById('eDesc').value = prod.description || '';
		document.getElementById('eImage').value = '';
		// preview and remove
		if (prod.image_url) {
			document.getElementById('eImagePreview').innerHTML = '<div class="img-thumb"><img src="'+prod.image_url+'" style="display:block;margin-bottom:6px"></div>';
			eRemoveImage = false;
			if (typeof eImageRemoveBtn !== 'undefined' && eImageRemoveBtn) eImageRemoveBtn.style.display = 'inline-block';
		} else {
			document.getElementById('eImagePreview').innerHTML = '';
			eRemoveImage = false;
			if (typeof eImageRemoveBtn !== 'undefined' && eImageRemoveBtn) eImageRemoveBtn.style.display = 'none';
		}
		document.getElementById('editModal').dataset.editing = id;
		document.getElementById('editModal').style.display = 'flex';
	}catch(e){ alert('Ошибка загрузки товара: '+e.message); }
}

document.getElementById('editClose').addEventListener('click', closeEditModal);
document.getElementById('cancelEdit').addEventListener('click', closeEditModal);
function closeEditModal(){ const m = document.getElementById('editModal'); m.style.display='none'; delete m.dataset.editing; }

async function initAdmin(){ await loadCategories(); await loadProducts(); } initAdmin();
</script></body></html>`

const productHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Товар</title>
<style>
body{font-family:Inter,Arial,Helvetica,sans-serif;margin:20px;background:#f5f7fb;color:#222}
.container{max-width:1100px;margin:0 auto}
.back{color:#2b6cb0;text-decoration:none;font-weight:600}
.product-header{display:flex;align-items:center;justify-content:space-between}
.product-wrap{display:flex;gap:28px;align-items:flex-start;margin-top:16px;flex-wrap:wrap}
.left{flex:0 0 48%;max-width:48%}
.left img{width:100%;height:auto;max-height:520px;object-fit:cover;border-radius:10px;border:1px solid #e6e6e6}
.right{flex:1;background:#fff;padding:20px;border-radius:10px;border:1px solid #e6e6e6;min-width:260px}
.title{font-size:24px;margin:0 0 8px}
.price{color:#2b6cb0;font-weight:700;font-size:20px;margin:8px 0}
.meta{color:#555;margin:6px 0}
.specs{display:grid;grid-template-columns:1fr 1fr;gap:8px;margin-top:12px}
.spec{background:#fafafa;padding:10px;border-radius:8px;border:1px solid #f0f0f0}
.actions{margin-top:16px;display:flex;gap:12px}
.btn{background:#2b6cb0;color:#fff;padding:10px 14px;border-radius:8px;border:none;cursor:pointer}
.btn.secondary{background:#e6eef9;color:#0b4a77}

/* flying cart + feedback */
.fly-cart{position:fixed;width:44px;height:44px;border-radius:50%;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,#ffd166,#f08a5d);color:#fff;font-size:20px;z-index:2000;transition:left .6s cubic-bezier(.2,.9,.2,1), top .6s cubic-bezier(.2,.9,.2,1), transform .6s cubic-bezier(.2,.9,.2,1), opacity .6s ease}
.bump{animation:bump .38s ease}
@keyframes bump{0%{transform:scale(1)}50%{transform:scale(1.12)}100%{transform:scale(1)}}

/* cart button one-time shake */
.cart-btn{background:#2b6cb0;color:#fff;border:none;padding:8px 12px;border-radius:8px;cursor:pointer}
.cart-btn.shake{animation:shake .6s ease}
@keyframes shake{0%{transform:translateX(0)}20%{transform:translateX(-6px)}40%{transform:translateX(6px)}60%{transform:translateX(-4px)}80%{transform:translateX(4px)}100%{transform:translateX(0)}}

.cart-panel{position:fixed;right:24px;bottom:24px;top:auto;width:360px;max-height:60vh;background:#fff;border-radius:10px;border:1px solid #e6e6e6;box-shadow:0 12px 40px rgba(11,20,40,0.12);overflow:auto;padding:12px;transform:translateX(120%);transition:transform .38s cubic-bezier(.2,.9,.2,1);z-index:1200}
.cart-panel.open{transform:translateX(0)}
.cart-panel h3{margin:8px 0 10px}
.cart-close{position:absolute;right:10px;top:8px;background:transparent;border:none;font-size:18px;cursor:pointer}

/* flying cart icon */
.fly-cart{position:fixed;width:44px;height:44px;border-radius:50%;display:flex;align-items:center;justify-content:center;background:linear-gradient(135deg,#ffd166,#f08a5d);color:#fff;font-size:20px;z-index:2000;transition:transform .6s ease,opacity .6s ease}

/* small bump animation for button */
.bump{animation:bump .38s ease}
@keyframes bump{0%{transform:scale(1)}50%{transform:scale(1.12)}100%{transform:scale(1)}}

/* cart button one-time shake */
.cart-btn.shake{animation:shake .6s ease}
@keyframes shake{0%{transform:translateX(0)}20%{transform:translateX(-6px)}40%{transform:translateX(6px)}60%{transform:translateX(-4px)}80%{transform:translateX(4px)}100%{transform:translateX(0)}}

@media(max-width:800px){.product-wrap{flex-direction:column}.left,.right{flex:1;max-width:100%}}
</style>
</head>
<body>
<div class="container">
    <div class="product-header"><a class="back" href="/">← Назад к каталогу</a><button id="btnCart" class="cart-btn">Корзина</button></div>
    <div class="product-wrap">
        <div class="left">
            {{if .product.ImageURL}}<img src="{{.product.ImageURL}}" alt="{{.product.Name}}">{{else}}<div style="width:100%;height:380px;background:#f0f2f5;border-radius:10px;display:flex;align-items:center;justify-content:center;color:#999">Нет изображения</div>{{end}}
        </div>
        <div class="right">
            <h1 class="title">{{.product.Name}}</h1>
            <div class="price">{{printf "%.2f" .product.Price}} ₽</div>
            <div class="meta">{{.product.Description}}</div>
            <div class="specs">
                <div class="spec"><strong>В наличии</strong><div>{{.product.Stock}}</div></div>
                <div class="spec"><strong>Состояние</strong><div>{{.product.Condition}}</div></div>
                <div class="spec"><strong>Цвет</strong><div>{{.product.Color}}</div></div>
                <div class="spec"><strong>Материал</strong><div>{{.product.Material}}</div></div>
                <div class="spec"><strong>Страна</strong><div>{{.product.Country}}</div></div>
            </div>
			<div class="actions">
				<button id="btnAdd" class="btn">Добавить в корзину</button>
			</div>
        </div>
    </div>

    <!-- slide-in cart panel (left) -->
    <aside id="cartPanel" class="cart-panel" aria-hidden="true">
        <button id="cartClose" class="cart-close">✕</button>
        <h3>Корзина</h3>
        <div id="cartContent"></div>
        <div style="margin-top:8px">
            <button id="checkoutBtn" class="btn">Оформить заказ</button>
        </div>
    </aside>
</div>

<script>
const productId = {{.product.ID}};
const productStock = {{.product.Stock}};

async function addToCartSilent(id){
    await fetch('/cart/add', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({product_id: id, quantity:1})});
}

function animateFly(fromRect, toX, toY){
	return new Promise((resolve)=>{
		const el = document.createElement('div');
		el.className = 'fly-cart';
		el.innerText = '🛒';
		// ensure fixed positioning
		el.style.position = 'fixed';
		el.style.zIndex = '2000';
		// start near button center
		el.style.left = (fromRect.left + (fromRect.width/2) - 22) + 'px';
		el.style.top = (fromRect.top + (fromRect.height/2) - 22) + 'px';
		el.style.opacity = '1';
		el.style.transform = 'scale(1)';
		el.style.transition = 'left .6s cubic-bezier(.2,.9,.2,1), top .6s cubic-bezier(.2,.9,.2,1), transform .6s cubic-bezier(.2,.9,.2,1), opacity .6s ease';
		document.body.appendChild(el);

		// target position: place icon center at (toX,toY)
		// run in next frame to trigger transition
		requestAnimationFrame(()=>{
			el.style.left = (toX - 22) + 'px';
			el.style.top = (toY - 22) + 'px';
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
		// safety fallback
		setTimeout(()=>{ if(el.parentNode) { cleanup(); } }, 900);
	});
}

// Local addToCart for product page: animate to header cart button and shake it once
async function addToCart(e, productId){
	var btnEl = null;
	try { if (e && e.currentTarget) btnEl = e.currentTarget; else if (e && e.target) btnEl = e.target.closest('button') || e.target; } catch(er) { btnEl = null; }
	const rect = btnEl ? btnEl.getBoundingClientRect() : {left: window.innerWidth/2, top: window.innerHeight/2, width:40, height:40};
	const cartBtn = document.getElementById('btnCart');
	const cartRect = (cartBtn && cartBtn.getBoundingClientRect) ? cartBtn.getBoundingClientRect() : {left: window.innerWidth-40, top: window.innerHeight-40, width:40, height:40};
	const toX = cartRect.left + (cartRect.width/2);
	const toY = cartRect.top + (cartRect.height/2);
	const anim = animateFly(rect, toX, toY);
	// server add
	try{ await fetch('/cart/add', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({product_id: productId, quantity:1})}); await loadCart(); }catch(e){ console.warn(e); }
	await anim;
	if (cartBtn){ cartBtn.classList.add('shake'); cartBtn.addEventListener('animationend', function _rm(){ cartBtn.classList.remove('shake'); cartBtn.removeEventListener('animationend', _rm); }); }
}

async function loadCart(){
	const r = await fetch('/cart');
	const j = await r.json();
	const cont = document.getElementById('cartContent');
	cont.innerHTML = '';
	if (!j.items || j.items.length === 0) { cont.innerHTML = '<div>Корзина пуста</div>'; return; }
	for (let it of j.items) {
		const wrapper = document.createElement('div');
		wrapper.style.display = 'flex';
		wrapper.style.gap = '10px';
		wrapper.style.alignItems = 'center';
		wrapper.style.marginBottom = '10px';
		const imgHtml = it.product.image_url ? '<img src="'+it.product.image_url+'" style="width:58px;height:58px;object-fit:cover;border-radius:6px;flex:0 0 58px">' : '<div style="width:58px;height:58px;background:#f0f2f5;border-radius:6px;flex:0 0 58px"></div>';
		const info = '<div style="flex:1"><div style="font-weight:600">'+(it.product.name||'')+'</div><div style="font-size:13px;color:#666">Цена: '+(it.product.price*it.quantity).toFixed(2)+'</div></div>';
		const controls = '<div style="display:flex;flex-direction:column;align-items:center;gap:6px"><div style="display:flex;align-items:center;gap:6px"><button class="qty-minus" data-id="'+it.product.id+'">−</button><span class="qty-value" data-id="'+it.product.id+'">'+it.quantity+'</span><button class="qty-plus" data-id="'+it.product.id+'">+</button></div><div style="font-size:12px;color:#666">доступно: '+(it.product.stock||0)+'</div></div>';
		wrapper.innerHTML = imgHtml + info + controls;
		cont.appendChild(wrapper);
		// attach handlers
		const minus = wrapper.querySelector('.qty-minus');
		const plus = wrapper.querySelector('.qty-plus');
		const qtySpan = wrapper.querySelector('.qty-value');
		if (minus) minus.addEventListener('click', async ()=>{
			let q = parseInt(qtySpan.innerText,10) - 1;
			if (q < 0) q = 0;
			await setCartQuantity(it.product.id, q);
			await loadCart();
		});
		if (plus) {
			if (it.quantity >= (it.product.stock||0)) plus.disabled = true;
			plus.addEventListener('click', async ()=>{
				let q = parseInt(qtySpan.innerText,10) + 1;
				if (q > (it.product.stock||0)) { alert('Превышен доступный запас'); return; }
				await setCartQuantity(it.product.id, q);
				await loadCart();
			});
		}
	}
	cont.insertAdjacentHTML('beforeend','<div style="margin-top:8px"><b>Итого: '+j.total.toFixed(2)+'</b></div>');
	// update add button state based on cart quantity
	try{
		const btn = document.getElementById('btnAdd');
		if(btn){
			const found = j.items.find(it => it.product && it.product.id === productId);
			const currentQty = found ? found.quantity : 0;
			btn.disabled = currentQty >= (productStock || 0);
		}
	}catch(e){console.warn(e)}
}

async function setCartQuantity(productId, quantity){
    const r = await fetch('/cart/update', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({product_id: productId, quantity: quantity})});
    if (!r.ok){
        const txt = await r.text();
        alert(txt || 'Ошибка обновления корзины');
    }
}

document.getElementById('btnAdd').addEventListener('click', async function(e){
	try{
		await addToCart(e, productId);
		const btn = e.currentTarget;
		btn.innerText = 'Добавлено';
		btn.classList.add('bump');
		setTimeout(()=>{ btn.classList.remove('bump'); btn.innerText = 'Добавить в корзину'; }, 900);
	}catch(err){ console.warn(err); }
});

document.getElementById('cartClose').addEventListener('click', ()=>{ const panel=document.getElementById('cartPanel'); panel.classList.remove('open'); panel.setAttribute('aria-hidden','true'); });

// Header cart button on product page
var headerCartBtn = document.getElementById('btnCart'); if (headerCartBtn){ headerCartBtn.addEventListener('click', function(){ location.href = '/cart/view'; }); }

document.getElementById('checkoutBtn').addEventListener('click', async ()=>{
    const name = prompt('Ваше имя:'); if(!name) return;
    const email = prompt('Email:'); if(!email) return;
    const address = prompt('Адрес доставки:'); if(!address) return;
    const r = await fetch('/order', {method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify({name:name,email:email,address:address})});
    if (!r.ok){ alert('Ошибка оформления'); return; }
    const data = await r.json(); alert('Заказ принят, id: ' + data.order_id);
    const panel=document.getElementById('cartPanel'); panel.classList.remove('open'); panel.setAttribute('aria-hidden','true');
});
</script>
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
  const r = await fetch('/cart');
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
async function setQty(productId, qty){ if (qty<0) qty=0; await fetch('/cart/update',{method:'POST',headers:{'Content-Type':'application/json'},body: JSON.stringify({product_id:productId, quantity: qty})}); await renderCart(); }
document.getElementById('goCheckout').addEventListener('click', async function(){ const name = prompt('Ваше имя:'); if(!name) return; const email=prompt('Email:'); if(!email) return; const address=prompt('Адрес доставки:'); if(!address) return; const r = await fetch('/order',{method:'POST',headers:{'Content-Type':'application/json'}, body: JSON.stringify({name:name,email:email,address:address})}); if (!r.ok){ alert('Ошибка оформления'); return; } const data = await r.json(); alert('Заказ принят, id: '+data.order_id); location.href='/'; });
document.getElementById('backToShop').addEventListener('click', ()=>{ location.href='/' });
renderCart();
</script>
</body>
</html>`
