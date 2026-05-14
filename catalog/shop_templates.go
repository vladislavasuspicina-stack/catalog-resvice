package main

const shopIndexHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>Shop Go</title>
<style>
:root{
  --blue:#005bff;
  --blue-dark:#0048d8;
  --pink:#f91155;
  --green:#00b75a;
  --teal:#12bdbd;
  --ink:#001a34;
  --muted:#667b90;
  --soft:#f1f4f9;
  --line:#dfe7f1;
  --card:#fff;
  --shadow:0 12px 30px rgba(0,26,52,.08);
}
*{box-sizing:border-box}
html{scroll-behavior:smooth}
body{margin:0;background:var(--soft);color:var(--ink);font-family:Inter,Segoe UI,Arial,sans-serif;font-size:16px}
button,input,select,textarea{font:inherit}
button{cursor:pointer}
a{color:inherit}
@keyframes softRise{from{opacity:0;transform:translateY(12px)}to{opacity:1;transform:translateY(0)}}
@keyframes accountMenuIn{from{opacity:0;transform:translate(-50%,-6px)}to{opacity:1;transform:translate(-50%,0)}}
@keyframes popIn{from{opacity:0;transform:translate(-50%,-48%) scale(.98)}to{opacity:1;transform:translate(-50%,-50%) scale(1)}}
@keyframes toastIn{from{opacity:0;transform:translate(-50%,10px)}to{opacity:1;transform:translate(-50%,0)}}
.page-bg{min-height:100vh}
.market-header{position:sticky;top:0;z-index:70;background:rgba(241,244,249,.96);backdrop-filter:saturate(140%) blur(10px)}
.header-card{max-width:1410px;margin:0 auto;background:#fff;border-radius:0 0 28px 28px;padding:10px 22px 14px;box-shadow:0 8px 24px rgba(0,26,52,.05)}
.header-main{display:grid;grid-template-columns:auto auto minmax(280px,1fr) auto;gap:14px;align-items:center}
.logo{font-size:42px;line-height:1;font-weight:900;letter-spacing:1px;text-decoration:none;color:var(--blue);transition:transform .18s ease}
.logo:hover{transform:translateY(-1px)}
.catalog-button{height:46px;border:0;border-radius:11px;background:var(--blue);color:#fff;padding:0 18px;display:flex;align-items:center;gap:10px;font-weight:800;transition:transform .16s ease,background .16s ease,box-shadow .16s ease}
.catalog-button:hover{transform:translateY(-1px);box-shadow:0 9px 18px rgba(0,91,255,.2)}
.catalog-button.open{background:var(--blue-dark)}
.catalog-button.open .grid-mark{transform:rotate(45deg)}
.grid-mark{width:18px;height:18px;display:grid;grid-template-columns:repeat(2,7px);gap:4px;transition:transform .2s ease}
.grid-mark span{display:block;border-radius:2px;background:#fff}
.search-box{height:46px;border:2px solid var(--blue);border-radius:12px;display:grid;grid-template-columns:1fr auto;align-items:center;overflow:hidden;background:#fff}
.search-box input{border:0;outline:0;min-width:0;padding:0 12px;color:var(--ink)}
.search-box input::placeholder{color:#8798aa}
.search-submit{height:100%;min-width:66px;border:0;background:var(--blue);color:#fff;font-weight:900;transition:background .16s ease}
.search-submit:hover{background:var(--blue-dark)}
.header-actions{display:flex;gap:18px;align-items:center;justify-content:flex-end}
.nav-link,.nav-button{position:relative;min-width:66px;border:0;background:transparent;color:#728397;text-decoration:none;display:flex;flex-direction:column;align-items:center;gap:5px;font-size:12px;line-height:1.15}
.nav-link:hover,.nav-button:hover{color:var(--blue)}
.nav-link.active{color:var(--blue)}
.nav-icon{width:28px;height:24px;position:relative;color:#9aaaba;display:block}
.nav-icon:before,.nav-icon:after{content:"";position:absolute;display:block}
.icon-user:before{left:9px;top:2px;width:9px;height:9px;border:2px solid currentColor;border-radius:50%}
.icon-user:after{left:5px;bottom:1px;width:18px;height:10px;border:2px solid currentColor;border-top:0;border-radius:0 0 10px 10px}
.icon-orders:before{left:5px;top:4px;width:17px;height:17px;border:2px solid currentColor;border-radius:5px;transform:rotate(45deg)}
.icon-favorite:before{left:4px;top:6px;width:10px;height:16px;background:currentColor;border-radius:10px 10px 0 0;transform:rotate(-45deg);transform-origin:100% 100%}
.icon-favorite:after{left:14px;top:6px;width:10px;height:16px;background:currentColor;border-radius:10px 10px 0 0;transform:rotate(45deg);transform-origin:0 100%}
.icon-cart:before{left:3px;top:4px;width:20px;height:15px;border:3px solid currentColor;border-radius:2px}
.icon-cart:after{left:8px;top:8px;width:10px;height:7px;border:2px solid currentColor;border-top:0}
.nav-link:hover .nav-icon,.nav-button:hover .nav-icon,.nav-link.active .nav-icon{color:var(--blue)}
.counter{position:absolute;top:-6px;right:8px;min-width:21px;height:18px;padding:0 6px;border-radius:999px;background:var(--pink);color:#fff;font-size:12px;font-weight:900;display:none;align-items:center;justify-content:center}
.subnav{display:flex;justify-content:flex-start;gap:20px;margin-top:12px;color:#66778a;font-size:14px;white-space:nowrap}
.subnav-left,.subnav-right{display:flex;align-items:center;gap:16px;min-width:0}
.subnav a{text-decoration:none}
.subnav a:hover{color:var(--blue)}
.fresh{color:#00a38f;font-weight:800}
.pickup{color:var(--blue);font-weight:700}
.account-wrap{position:relative;display:flex;justify-content:center}
.account-menu{position:absolute;left:50%;top:48px;width:min(274px,calc(100vw - 28px));background:#fff;border:1px solid var(--line);border-radius:14px;box-shadow:0 16px 36px rgba(0,26,52,.16);padding:10px;display:none;transform:translateX(-50%)}
.account-menu.open{display:block;animation:accountMenuIn .18s ease both}
.account-menu:before{content:"";position:absolute;left:50%;top:-9px;transform:translateX(-50%);border-left:9px solid transparent;border-right:9px solid transparent;border-bottom:9px solid #fff}
.menu-item{display:flex;align-items:center;justify-content:space-between;gap:10px;padding:10px 12px;border-radius:10px;text-decoration:none;color:#1b2838}
.menu-item:hover{background:#f5f8fc}
.menu-item small{display:block;color:#667b90;font-size:14px;margin-top:2px}
.menu-badge{min-width:32px;height:22px;border-radius:999px;background:var(--pink);color:#fff;display:grid;place-items:center;font-size:12px;font-weight:900}
.menu-separator{height:1px;background:var(--line);margin:8px -10px}
.account-title{font-weight:900;color:#16314c;padding:0 14px 10px}
.catalog-panel{position:fixed;inset:82px 0 0;z-index:65;background:rgba(255,255,255,.98);opacity:0;visibility:hidden;pointer-events:none;overflow:hidden;transform:translateY(-10px);transition:opacity .2s ease,transform .22s ease,visibility .22s}
.catalog-panel.open{opacity:1;visibility:visible;pointer-events:auto;transform:translateY(0)}
.catalog-inner{height:100%;max-width:1410px;margin:0 auto;padding:28px 22px 42px}
.catalog-content{height:100%;overflow:auto}
.catalog-content h1{font-size:36px;margin:0 0 22px}
.catalog-columns{display:grid;grid-template-columns:repeat(3,minmax(220px,1fr));gap:16px 18px}
.catalog-group{border:1px solid #e6edf6;border-radius:8px;background:#fff;padding:18px;box-shadow:0 1px 0 rgba(0,26,52,.03);transition:transform .18s ease,box-shadow .18s ease,border-color .18s ease}
.catalog-group:hover{transform:translateY(-2px);border-color:#cfe0f3;box-shadow:0 12px 28px rgba(0,26,52,.08)}
.catalog-group h3{margin:0 0 12px;font-size:17px}
.catalog-group a{display:block;text-decoration:none;color:#5f7187;margin:0 0 10px;line-height:1.35}
.catalog-group a:last-child{margin-bottom:0}
.catalog-group a:hover{color:var(--blue)}
.shell{max-width:1410px;margin:0 auto;padding:28px 22px 60px}
.page-view{display:none}
.page-view.active{display:block;animation:softRise .22s ease both}
.promo-banner{height:clamp(150px,22vw,300px);border-radius:24px;background:#080305;margin:0 0 14px;overflow:hidden;box-shadow:0 18px 38px rgba(0,26,52,.12);animation:softRise .32s ease both}
.promo-banner img{width:100%;height:100%;display:block;object-fit:cover;object-position:center}
.section-line{display:flex;align-items:flex-end;justify-content:space-between;gap:16px;margin:22px 0 18px}
.section-line h1,.section-line h2{margin:0;font-size:34px;line-height:1.05}
.section-tabs{display:flex;align-items:center;gap:12px;color:#6b7f94;font-size:30px;font-weight:900}
.section-tabs a{text-decoration:none;color:#6b7f94}
.section-tabs a.active{color:#000}
.tiny-count{font-size:13px;background:#e8eef6;border-radius:999px;padding:3px 7px;color:#60758b;vertical-align:middle}
.home-tools{display:flex;align-items:center;justify-content:space-between;gap:14px;margin-bottom:10px}
.pill{height:38px;border:0;border-radius:999px;background:#fff;color:#304256;padding:0 14px;font-weight:700;white-space:nowrap;transition:transform .16s ease,background .16s ease,color .16s ease}
.pill:hover{transform:translateY(-1px)}
.pill.active{background:#001a34;color:#fff}
.brand-row{display:flex;gap:8px;overflow:auto;margin:0 0 16px;padding-bottom:2px}
.brand-pill{height:34px;border:1px solid var(--line);border-radius:999px;background:#fff;color:#52677c;padding:0 13px;font-weight:800;white-space:nowrap;transition:transform .16s ease,border-color .16s ease,background .16s ease}
.brand-pill:hover{transform:translateY(-1px);border-color:#b9cce0}
.brand-pill.active{border-color:var(--blue);background:#eef6ff;color:var(--blue)}
.select{height:44px;border:1px solid var(--line);border-radius:12px;background:#fff;padding:0 42px 0 14px;color:#1b2838;outline:none}
.product-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(190px,1fr));gap:18px}
.product-card{background:#fff;border-radius:8px;overflow:hidden;position:relative;min-width:0;box-shadow:0 1px 0 rgba(0,26,52,.04);display:flex;flex-direction:column;animation:softRise .24s ease both;transition:transform .18s ease,box-shadow .18s ease}
.product-card:hover{transform:translateY(-3px);box-shadow:var(--shadow)}
.product-media{position:relative;display:block;aspect-ratio:1/1.22;background:#fff;overflow:hidden;text-decoration:none}
.product-media img{width:100%;height:100%;object-fit:contain;display:block;padding:8px;background:#fff;transition:transform .22s ease}
.product-card:hover .product-media img{transform:scale(1.035)}
.media-placeholder{height:100%;display:grid;place-items:center;background:linear-gradient(135deg,#e9f1ff,#f6f9ff)}
.media-placeholder span{font-weight:900;color:#8aa0b8;font-size:34px;letter-spacing:.04em}
.heart{position:absolute;right:10px;top:10px;width:34px;height:34px;border:0;border-radius:50%;background:rgba(255,255,255,.94);font-size:22px;line-height:1;color:#001a34;box-shadow:0 4px 12px rgba(0,26,52,.12);transition:transform .16s ease,color .16s ease}
.heart:hover{transform:scale(1.08)}
.heart.active{color:var(--pink)}
.sale-label{position:absolute;left:8px;bottom:8px;background:var(--pink);color:#fff;border-radius:6px;padding:5px 9px;font-weight:900;font-size:13px}
.out-label{position:absolute;left:8px;bottom:8px;background:#001a34;color:#fff;border-radius:6px;padding:5px 9px;font-weight:900;font-size:13px}
.card-body{padding:10px 9px 12px;display:flex;flex-direction:column;gap:6px;min-height:184px}
.price-row{display:flex;align-items:baseline;gap:7px;flex-wrap:wrap}
.price{font-size:21px;font-weight:900;color:var(--pink)}
.old-price{font-size:14px;color:#8a99a8;text-decoration:line-through;font-weight:700}
.discount{font-size:13px;color:var(--pink);font-weight:900}
.product-name{display:block;text-decoration:none;color:#001a34;line-height:1.28;min-height:41px;overflow:hidden}
.meta-row{display:flex;gap:6px;align-items:center;color:#667b90;font-size:13px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.star{color:#f5a000;font-weight:900}
.delivery-button{height:36px;border:0;border-radius:8px;background:var(--blue);color:#fff;font-weight:900;width:100%;margin-top:auto;transition:transform .16s ease,background .16s ease}
.delivery-button:hover:not(:disabled){transform:translateY(-1px);background:var(--blue-dark)}
.delivery-button:disabled{background:#d7e0ea;color:#7f91a4}
.empty-state{background:#fff;border-radius:22px;padding:46px 22px;text-align:center;color:#667b90}
.empty-illustration{width:118px;height:88px;border-radius:28px;background:linear-gradient(135deg,#f8ce9a,#d98b4d);margin:0 auto 18px;position:relative}
.empty-illustration:before{content:"";position:absolute;left:22px;right:22px;top:-16px;height:32px;background:#f7bd77;border-radius:9px;transform:skewX(-18deg)}
.empty-state h2{margin:0 0 8px;color:#000;font-size:25px}
.blue-button,.light-button,.dark-button{height:46px;border:0;border-radius:11px;padding:0 18px;font-weight:900;text-decoration:none;display:inline-flex;align-items:center;justify-content:center;gap:8px;transition:transform .16s ease,background .16s ease,box-shadow .16s ease}
.blue-button:hover,.light-button:hover,.dark-button:hover{transform:translateY(-1px)}
.blue-button{background:var(--blue);color:#fff}
.blue-button:hover{background:var(--blue-dark)}
.light-button{background:#eef6ff;color:var(--blue)}
.dark-button{background:#001a34;color:#fff}
.favorites-layout{display:grid;grid-template-columns:1fr;gap:24px;align-items:start}
.filter-panel{background:#fff;border-radius:22px;padding:22px;position:sticky;top:118px}
.filter-panel h3{margin:0 0 18px;font-size:18px}
.filter-list{display:grid;gap:14px}
.toggle-row{display:flex;justify-content:space-between;align-items:center;gap:10px;font-weight:800;margin-top:24px}
.switch{width:44px;height:28px;border-radius:999px;background:#d9e2ec;position:relative}
.switch:before{content:"";position:absolute;width:24px;height:24px;border-radius:50%;background:#fff;left:2px;top:2px;box-shadow:0 2px 5px rgba(0,26,52,.14)}
.cart-title{display:flex;align-items:center;gap:4px;margin:0 0 28px;font-size:32px}
.cart-layout{display:grid;grid-template-columns:minmax(0,1fr) 380px;gap:24px;align-items:start}
.cart-main{display:grid;gap:10px}
.sale-strip,.select-strip,.cart-list{background:#fff;border-radius:22px}
.sale-strip{height:66px;padding:0 18px;display:flex;align-items:center;justify-content:space-between;gap:14px}
.sale-left{display:flex;align-items:center;gap:12px}
.fire-dot{width:42px;height:42px;border-radius:12px;background:#fff1f6;color:var(--pink);display:grid;place-items:center;font-weight:900}
.countdown{background:#fff1f6;color:var(--pink);border-radius:999px;padding:7px 12px;font-weight:900;font-size:13px}
.select-strip{height:58px;padding:0 18px;display:flex;align-items:center;justify-content:space-between}
.checkline{display:flex;align-items:center;gap:10px;font-weight:800}
.check{width:22px;height:22px;border-radius:7px;background:var(--blue);display:grid;place-items:center;color:#fff;font-size:15px;font-weight:900}
.cart-list{padding:16px}
.cart-list h3{margin:0 0 10px;font-size:17px;background:#f4f7fb;border-radius:14px;padding:14px 16px}
.cart-row{display:grid;grid-template-columns:24px 92px minmax(0,1fr) 150px 132px;gap:12px;align-items:start;padding:18px 0;border-bottom:1px solid #eef2f7;animation:softRise .22s ease both}
.cart-row:last-child{border-bottom:0}
.cart-row img,.cart-picture{width:92px;height:118px;border-radius:10px;object-fit:contain;background:#fff;display:grid;place-items:center;color:#8fa0b2;font-weight:900}
.cart-name{font-size:17px;line-height:1.28;margin-bottom:10px}
.cart-actions{display:flex;gap:8px;flex-wrap:wrap}
.mini-action{height:30px;border:0;border-radius:9px;background:#f3f6fa;color:#001a34;padding:0 11px;font-weight:800;transition:transform .16s ease,background .16s ease}
.mini-action:hover{transform:translateY(-1px);background:#e7eef7}
.cart-price{font-size:18px;font-weight:900;color:var(--pink)}
.cart-old{font-size:13px;color:#91a0af;text-decoration:line-through;font-weight:700;margin-top:3px}
.quantity-box{height:36px;border-radius:10px;background:#f3f6fa;display:grid;grid-template-columns:36px 1fr 36px;align-items:center;text-align:center;font-weight:900}
.quantity-box button{height:36px;border:0;background:transparent;font-size:22px;font-weight:900;color:#001a34}
.limit-note{font-size:12px;color:#f04438;font-weight:800;margin-top:8px;text-align:center}
.summary-panel{background:#fff;border-radius:22px;position:sticky;top:118px;overflow:hidden}
.summary-top{padding:24px;border-bottom:1px solid var(--line)}
.summary-top .blue-button{width:100%;margin-bottom:16px}
.summary-note{color:#667b90;font-size:14px;line-height:1.35}
.summary-body{padding:22px 24px}
.summary-head{display:flex;justify-content:space-between;gap:10px;align-items:center;margin-bottom:18px}
.summary-head h2{margin:0;font-size:22px}
.summary-line{display:flex;justify-content:space-between;gap:12px;margin:12px 0;color:#001a34}
.summary-line.discount{color:var(--pink);font-weight:800}
.summary-total{display:flex;justify-content:space-between;align-items:flex-start;gap:12px;padding-top:18px;margin-top:18px;border-top:1px solid var(--line);font-weight:900}
.summary-total strong{font-size:24px;color:var(--green)}
.orders-layout{display:grid;grid-template-columns:216px 1fr;gap:16px;align-items:start}
.profile-panel{background:#fff;border-radius:22px;padding:20px;position:sticky;top:118px}
.avatar{width:92px;height:92px;border-radius:50%;background:linear-gradient(135deg,#dbeafe,#7dd3fc);display:grid;place-items:center;font-size:36px;font-weight:900;color:var(--blue);margin-bottom:14px}
.profile-name{font-size:21px;line-height:1.12;font-weight:900;margin-bottom:6px}
.profile-link{color:var(--blue);font-weight:800;text-decoration:none;font-size:14px}
.orders-head{background:#fff;border-radius:22px;padding:26px}
.orders-head h1{margin:0;font-size:34px}
.tab-buttons{display:flex;gap:10px;margin-top:18px}
.tab-button{height:42px;border:0;border-radius:13px;background:#f3f6fa;padding:0 18px;font-weight:800}
.tab-button.active{background:#000;color:#fff}
.orders-empty{min-height:340px;display:grid;place-items:center;text-align:center}
.orders-list{display:grid;gap:12px;margin-top:16px}
.order-card{background:#fff;border-radius:18px;padding:18px;border:1px solid #edf2f7;animation:softRise .24s ease both;transition:transform .18s ease,box-shadow .18s ease}
.order-card:hover{transform:translateY(-2px);box-shadow:0 12px 28px rgba(0,26,52,.08)}
.order-card-head{display:flex;justify-content:space-between;gap:16px;margin-bottom:12px}
.status-chip{border-radius:999px;background:#fff1f6;color:var(--pink);padding:6px 10px;font-weight:900;font-size:13px}
.recommend-block{background:#fff;border-radius:22px;padding:22px;margin-top:16px}
.recommend-block h2{margin:0 0 16px;font-size:28px}
.modal-cover{position:fixed;inset:0;background:rgba(0,26,52,.42);z-index:120;display:none}
.modal-cover[style*="block"]{animation:softRise .16s ease both}
.checkout-modal{position:fixed;left:50%;top:50%;transform:translate(-50%,-50%);width:min(760px,94vw);max-height:92vh;overflow:auto;background:#fff;border-radius:22px;z-index:130;display:none;box-shadow:0 28px 80px rgba(0,26,52,.28)}
.checkout-modal[style*="block"]{animation:popIn .2s ease both}
.modal-head{display:flex;justify-content:space-between;gap:14px;align-items:center;padding:20px 22px;border-bottom:1px solid var(--line)}
.modal-head h2{margin:0;font-size:24px}
.close{width:40px;height:40px;border:0;border-radius:12px;background:#f2f6fb;font-weight:900;font-size:22px}
.modal-body{padding:22px}
.checkout-grid{display:grid;grid-template-columns:1fr 1fr;gap:14px}
.field{display:grid;gap:7px}
.field label{font-size:13px;font-weight:900;color:#5d7085}
.field input,.field textarea,.field select{width:100%;border:1px solid var(--line);border-radius:12px;padding:12px 13px;outline:none;background:#fff}
.field input:focus,.field textarea:focus,.field select:focus{border-color:var(--blue);box-shadow:0 0 0 3px rgba(0,91,255,.12)}
.delivery-choice{display:grid;grid-template-columns:1fr 1fr;gap:10px;margin:16px 0}
.delivery-choice label{border:1px solid var(--line);border-radius:16px;padding:14px;cursor:pointer;font-weight:800}
.pickup-info{margin-top:8px;background:#f4f7fb;border-radius:14px;padding:12px;color:#52677c;font-size:14px;line-height:1.45}
.success{display:none;text-align:center;padding:38px 24px}
.success-mark{width:76px;height:76px;border-radius:50%;margin:0 auto 16px;background:var(--green);color:#fff;display:grid;place-items:center;font-size:44px;font-weight:900}
.toast{position:fixed;left:50%;bottom:24px;transform:translateX(-50%);z-index:150;background:#001a34;color:#fff;border-radius:999px;padding:12px 18px;display:none;font-weight:900;box-shadow:0 12px 30px rgba(0,26,52,.25)}
.toast[style*="block"]{animation:toastIn .18s ease both}
@media(prefers-reduced-motion:reduce){
  *,*:before,*:after{animation:none!important;transition:none!important;scroll-behavior:auto!important}
}
.muted{color:#667b90}
@media(max-width:1100px){
  .header-main{grid-template-columns:auto auto 1fr}.header-actions{grid-column:1/-1;justify-content:space-between}.subnav{overflow:auto}.cart-layout,.orders-layout,.favorites-layout{grid-template-columns:1fr}.summary-panel,.profile-panel,.filter-panel{position:static}.cart-row{grid-template-columns:24px 78px 1fr}.cart-price,.quantity-box{grid-column:3}.catalog-columns{grid-template-columns:repeat(2,minmax(220px,1fr))}
}
@media(max-width:760px){
  body{font-size:15px}.header-card{padding:8px 10px 12px;border-radius:0 0 18px 18px}.header-main{grid-template-columns:1fr auto}.logo{font-size:34px}.catalog-button{height:42px;padding:0 12px}.search-box{grid-column:1/-1;order:3}.header-actions{gap:6px;overflow:auto}.nav-link,.nav-button{min-width:58px}.nav-icon{width:25px;height:22px}.subnav-right{display:none}.shell{padding:18px 10px 50px}.promo-banner{margin:0 0 18px;border-radius:18px}.section-line{align-items:flex-start;flex-direction:column}.section-line h1,.section-line h2,.section-tabs{font-size:27px}.home-tools{align-items:flex-start;flex-direction:column}.product-grid{grid-template-columns:repeat(2,minmax(0,1fr));gap:10px}.card-body{min-height:170px}.price{font-size:18px}.old-price,.discount,.meta-row{font-size:12px}.catalog-panel{inset:64px 0 0}.catalog-inner{padding:14px 10px 28px}.catalog-content h1{font-size:28px}.catalog-columns{grid-template-columns:1fr}.catalog-group{padding:15px}.checkout-grid,.delivery-choice{grid-template-columns:1fr}.cart-row{gap:9px}.cart-row img,.cart-picture{width:78px;height:98px}.account-menu{width:min(274px,calc(100vw - 20px))}
}
</style>
</head>
<body>
<div class="page-bg">
<header class="market-header">
  <div class="header-card">
    <div class="header-main">
      <a class="logo" href="/" aria-label="Shop Go">SHOPGO</a>
      <button class="catalog-button" id="catalogToggle" type="button">
        <span class="grid-mark" aria-hidden="true"><span></span><span></span><span></span><span></span></span>
        <span>Каталог</span>
      </button>
      <div class="search-box">
        <input id="searchInput" placeholder="Искать на Shop Go" autocomplete="off">
        <button class="search-submit" id="searchBtn" type="button">Найти</button>
      </div>
      <nav class="header-actions" aria-label="Основная навигация">
        <div class="account-wrap">
          <button class="nav-button" id="accountBtn" type="button"><span class="nav-icon icon-user" aria-hidden="true"></span><span id="accountName">Войти</span></button>
          <div class="account-menu" id="accountMenu">
            <a class="menu-item" href="/auth" id="loginMenuItem"><span><b>Личный кабинет</b><small>Войти или создать аккаунт</small></span></a>
            <a class="menu-item" href="/orders"><span><b>Заказы</b><small>Покупки и возвраты</small></span></a>
            <a class="menu-item" href="/favorites"><span><b>Избранное</b><small>Сохранённые товары</small></span><span class="menu-badge" id="menuFavBadge" style="display:none"></span></a>
            <a class="menu-item" href="/cart/view"><span><b>Корзина</b><small>Оформление заказа</small></span><span class="menu-badge" id="menuCartBadge" style="display:none"></span></a>
            <a class="menu-item" href="/auth?next=/admin" id="adminMenuItem"><span><b>Админ-панель</b><small>Вход: admin / 12345</small></span></a>
            <div class="menu-separator"></div>
            <a class="menu-item" href="#" id="logoutMenuItem" style="display:none"><span><b>Выйти</b></span></a>
          </div>
        </div>
        <a class="nav-link" data-route="orders" href="/orders"><span class="counter" id="ordersBadge"></span><span class="nav-icon icon-orders" aria-hidden="true"></span><span>Заказы</span></a>
        <a class="nav-link" data-route="favorites" href="/favorites"><span class="counter" id="favBadge"></span><span class="nav-icon icon-favorite" aria-hidden="true"></span><span>Избранное</span></a>
        <a class="nav-link" data-route="cart" href="/cart/view"><span class="counter" id="cartBadge"></span><span class="nav-icon icon-cart" aria-hidden="true"></span><span>Корзина</span></a>
      </nav>
    </div>
    <div class="subnav">
      <div class="subnav-left">
        <a class="fresh" href="/">Новинки гаджетов</a>
        <a href="/?search=смартфон">Смартфоны</a>
        <a href="/?search=ноутбук">Ноутбуки</a>
        <a href="/?search=компьютер">Компьютеры</a>
        <a href="/?search=планшет">Планшеты</a>
        <a href="/?search=часы">Смарт-часы</a>
        <a href="/?search=комплектующие">Комплектующие</a>
        <a href="/?search=аксессуары">Аксессуары</a>
      </div>
    </div>
  </div>
</header>

<section class="catalog-panel" id="catalogPanel" aria-hidden="true">
  <div class="catalog-inner">
    <div class="catalog-content">
      <h1>Электроника</h1>
      <div class="catalog-columns">
        <div class="catalog-group"><h3>Телефоны и смарт-часы</h3><a href="/?search=смартфон">Смартфоны</a><a href="/?search=часы">Смарт-часы</a><a href="/?brand=Samsung">Samsung</a><a href="/?brand=Xiaomi">Xiaomi</a><a href="/?brand=Apple">Apple</a><a href="/?brand=Huawei">Huawei</a></div>
        <div class="catalog-group"><h3>Аксессуары и звук</h3><a href="/?search=AirPods">AirPods</a><a href="/?search=Buds">Samsung Galaxy Buds</a><a href="/?search=наушники">Наушники</a><a href="/?search=JBL">Портативные колонки</a><a href="/?search=SmartTag">Метки для поиска вещей</a></div>
        <div class="catalog-group"><h3>Зарядка и питание</h3><a href="/?search=MagSafe">MagSafe</a><a href="/?search=заряд">Зарядные устройства</a><a href="/?search=Power Bank">Внешние аккумуляторы</a><a href="/?brand=Anker">Anker</a><a href="/?brand=Belkin">Belkin</a></div>
        <div class="catalog-group"><h3>Компьютерные аксессуары</h3><a href="/?search=мышь">Мыши</a><a href="/?search=клавиатура">Клавиатуры</a><a href="/?search=SSD">Внешние SSD</a><a href="/?brand=Logitech">Logitech</a><a href="/?brand=Keychron">Keychron</a><a href="/?brand=SanDisk">SanDisk</a></div>
      </div>
    </div>
  </div>
</section>

<main class="shell">
  <section class="page-view" id="homePage">
    <div class="promo-banner"><img src="/static/banners/phone-sale-banner.png" alt="Скидки на телефоны до 80 процентов"></div>
    <div class="home-tools">
      <span class="muted" id="resultInfo">Загрузка товаров...</span>
      <select class="select" id="sortSelect">
        <option value="id_asc">Сначала популярные</option>
        <option value="newest">Сначала новые</option>
        <option value="price_asc">Дешевле</option>
        <option value="price_desc">Дороже</option>
        <option value="name_asc">По названию</option>
      </select>
    </div>
    <div class="brand-row" id="brandPills"></div>
    <input type="hidden" id="categoryFilter" value="">
    <input type="hidden" id="brandFilter" value="">
    <div class="product-grid" id="productsGrid"></div>
  </section>

  <section class="page-view" id="favoritesPage">
    <div class="section-line">
      <div class="section-tabs">
        <a class="active" href="/favorites">Избранное <span class="tiny-count" id="favoritesTitleCount">0</span></a>
      </div>
    </div>
    <div class="favorites-layout">
      <section>
        <select class="select" id="favoriteSort">
          <option value="new">Сначала новые</option>
          <option value="price_asc">Сначала дешевле</option>
          <option value="price_desc">Сначала дороже</option>
        </select>
        <div class="product-grid" id="favoritesGrid" style="margin-top:12px"></div>
      </section>
    </div>
  </section>

  <section class="page-view" id="cartPage">
    <h1 class="cart-title">Корзина<span class="tiny-count" id="cartTitleCount">0</span></h1>
    <div class="cart-layout">
      <section class="cart-main">
        <div class="sale-strip">
          <div class="sale-left"><span class="fire-dot">%</span><div><b>Не упустите распродажу</b><br><span class="muted" id="cartSaleHint">Товары скоро подорожают</span></div></div>
          <span class="countdown">Осталось 2 дня</span>
        </div>
        <div class="select-strip">
          <div class="checkline"><span class="check">✓</span><span>Выбрать все</span></div>
          <div><button class="mini-action" type="button">Поделиться</button></div>
        </div>
        <div class="cart-list">
          <h3>Доступны для заказа</h3>
          <div id="cartPageItems"></div>
        </div>
      </section>
      <aside>
        <div class="summary-panel">
          <div class="summary-top">
            <button class="blue-button" id="checkoutOpen" type="button">Перейти к оформлению</button>
            <div class="summary-note">Доступные способы и время доставки можно выбрать при оформлении заказа</div>
          </div>
          <div class="summary-body">
            <div class="summary-head"><h2>Ваша корзина</h2><span class="muted" id="summaryCount">0 товаров</span></div>
            <div class="summary-line"><span>Товары</span><b id="summarySubtotal">0 ₽</b></div>
            <div class="summary-line discount"><span>Скидка</span><b id="summaryDiscount">0 ₽</b></div>
            <div class="summary-total"><span>Итого<br><small class="muted">К оплате при получении</small></span><strong id="summaryTotal">0 ₽</strong></div>
          </div>
        </div>
      </aside>
    </div>
  </section>

  <section class="page-view" id="ordersPage">
    <div class="orders-layout">
      <aside class="profile-panel">
        <div class="avatar" id="profileAvatar">S</div>
        <div class="profile-name" id="profileName">Гость</div>
      </aside>
      <section>
        <div class="orders-head">
          <h1>Заказы · Покупки · Сервис гаджетов</h1>
          <div class="tab-buttons"><button class="tab-button active" type="button" data-order-tab="active">Актуальные</button><button class="tab-button" type="button" data-order-tab="finished">Завершенные</button></div>
        </div>
        <div id="ordersContent"></div>
        <div class="recommend-block">
          <h2>Подобрали по вашим интересам</h2>
          <div class="product-grid" id="ordersRecommendGrid"></div>
        </div>
      </section>
    </div>
  </section>
</main>

<div class="modal-cover" id="modalCover"></div>
<section class="checkout-modal" id="checkoutModal">
  <div class="modal-head"><h2>Оформление заказа</h2><button class="close" id="closeCheckout" type="button">×</button></div>
  <div id="checkoutFormWrap">
    <form class="modal-body" id="checkoutForm">
      <div class="checkout-grid">
        <div class="field"><label for="coName">Имя</label><input id="coName" required placeholder="Иван Петров"></div>
        <div class="field"><label for="coEmail">Email</label><input id="coEmail" required type="email" placeholder="ivan@example.com"></div>
        <div class="field"><label for="coPhone">Телефон</label><input id="coPhone" placeholder="+7 999 123-45-67"></div>
        <div class="field"><label for="coComment">Комментарий</label><input id="coComment" placeholder="Позвонить заранее"></div>
      </div>
      <div class="delivery-choice">
        <label><input type="radio" name="delivery" value="pickup" checked> Самовывоз бесплатно</label>
        <label><input type="radio" name="delivery" value="courier"> Курьерская доставка</label>
      </div>
      <div id="pickupBox">
        <div class="field">
          <label for="pickupSelect">Пункт выдачи</label>
          <select id="pickupSelect"></select>
          <div class="pickup-info" id="pickupInfo">Загрузка пунктов выдачи...</div>
        </div>
      </div>
      <div id="courierBox" style="display:none">
        <div class="field"><label for="coAddress">Адрес доставки</label><textarea id="coAddress" rows="3" placeholder="Город, улица, дом, квартира"></textarea></div>
      </div>
      <button class="blue-button" id="submitOrder" type="submit" style="width:100%;margin-top:14px">Завершить заказ</button>
    </form>
  </div>
  <div class="success" id="checkoutSuccess">
    <div class="success-mark">✓</div>
    <h2>Заказ оформлен</h2>
    <p class="muted" id="successText"></p>
    <a class="blue-button" href="/orders">К заказам</a>
  </div>
</section>

<div class="toast" id="toast"></div>
</div>

<script>
const qs = id => document.getElementById(id);
const money = value => new Intl.NumberFormat('ru-RU',{style:'currency',currency:'RUB',maximumFractionDigits:0}).format(Number(value||0));
function formatDate(value){
  if(!value) return '';
  const d = new Date(value);
  return Number.isNaN(d.getTime()) ? String(value) : d.toLocaleString('ru-RU');
}
const esc = value => String(value ?? '').replace(/[&<>"']/g, c => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));
const store = {
  get(key, fallback){try{return JSON.parse(localStorage.getItem(key)) ?? fallback}catch(e){return fallback}},
  set(key, value){localStorage.setItem(key, JSON.stringify(value))}
};
let state = {
  products: [],
  categories: [],
  brands: [],
  cart: {items:[], total:0},
	favorites: [],
	orders: [],
	orderTab: 'active',
	pickupPoints: [],
  selectedPickupId: 0,
  user: null,
  page: 1,
  perPage: 80,
  total: 0,
  totalPages: 1
};

function storagePrefix(){
  return state.user && state.user.id ? 'shopgo:user:' + state.user.id : 'shopgo:guest';
}
function storageKey(name){
  return storagePrefix() + ':' + name;
}
function loadScopedState(){
  state.favorites = store.get(storageKey('favorites'), []);
  state.orders = store.get(storageKey('orders'), []);
}
function saveFavorites(){
  store.set(storageKey('favorites'), state.favorites);
}
function saveOrders(){
  store.set(storageKey('orders'), state.orders);
}

function currentRoute(){
  const path = location.pathname;
  if(path.startsWith('/favorites')) return 'favorites';
  if(path.startsWith('/cart')) return 'cart';
  if(path.startsWith('/orders')) return 'orders';
  return 'home';
}
function toast(text){
  const el = qs('toast');
  el.textContent = text;
  el.style.display = 'block';
  clearTimeout(window.toastTimer);
  window.toastTimer = setTimeout(()=>el.style.display='none', 1900);
}
async function api(url, options){
  const res = await fetch(url, Object.assign({credentials:'include'}, options || {}));
  const text = await res.text();
  let payload = {};
  try{payload = text ? JSON.parse(text) : {}}catch(e){payload = {error:text}}
  if(!res.ok) throw new Error(payload.error || text || 'Ошибка запроса');
  return payload;
}
function oldPrice(p){
  const ratio = 1.12 + ((Number(p.id)||1) % 5) * 0.08;
  return Math.round(Number(p.price || 0) * ratio);
}
function discountText(p){
  const old = oldPrice(p);
  const price = Number(p.price || 0);
  if(!old || old <= price) return '';
  return '-' + Math.round((old - price) / old * 100) + '%';
}
function ratingFor(p){
  return (4.6 + ((Number(p.id)||1) % 4) * 0.1).toFixed(1);
}
function isFavorite(id){return state.favorites.map(Number).includes(Number(id))}
function productImage(p, className){
  if(p.image_url) return '<img src="'+esc(p.image_url)+'" alt="'+esc(p.name)+'">';
  return '<div class="'+(className || 'media-placeholder')+'"><span>SHOP</span></div>';
}
function productCard(p){
  const stock = Number(p.stock || 0);
  const fav = isFavorite(p.id);
  const old = oldPrice(p);
  const meta = p.brand ? esc(p.brand) + (p.condition ? ' · ' + esc(p.condition) : '') : esc(p.condition || 'Оригинал');
  const label = stock > 0 ? '<span class="sale-label">Распродажа</span>' : '<span class="out-label">Нет в наличии</span>';
  return '<article class="product-card">'+
    '<a class="product-media" href="/product/'+p.id+'">'+productImage(p)+label+'</a>'+
    '<button class="heart '+(fav?'active':'')+'" type="button" data-fav="'+p.id+'" aria-label="В избранное">♥</button>'+
    '<div class="card-body">'+
      '<div class="price-row"><span class="price">'+money(p.price)+'</span><span class="old-price">'+money(old)+'</span><span class="discount">'+discountText(p)+'</span></div>'+
      '<a class="product-name" href="/product/'+p.id+'">'+esc(p.name)+'</a>'+
      '<div class="meta-row"><span class="star">★</span><b>'+ratingFor(p)+'</b><span>проверенный товар</span></div>'+
      '<div class="meta-row">'+meta+(p.country ? ' · '+esc(p.country) : '')+'</div>'+
      '<button class="delivery-button" type="button" data-add="'+p.id+'" '+(stock<=0?'disabled':'')+'>'+(stock>0 ? 'Купить' : 'Нет в наличии')+'</button>'+
    '</div>'+
  '</article>';
}
function emptyProducts(text, action){
  return '<div class="empty-state" style="grid-column:1/-1"><div class="empty-illustration"></div><h2>'+esc(text)+'</h2>'+(action || '')+'</div>';
}
async function loadUser(){
  const data = await api('/auth/me');
  state.user = data.authenticated ? data.user : null;
  loadScopedState();
  const name = state.user ? state.user.username : 'Войти';
  qs('accountName').textContent = name;
  qs('profileName').textContent = state.user ? state.user.username : 'Гость';
  qs('profileAvatar').textContent = (state.user ? state.user.username : 'S').slice(0,1).toUpperCase();
  qs('loginMenuItem').href = state.user ? '/orders' : '/auth';
  qs('loginMenuItem').querySelector('small').textContent = state.user ? 'Профиль и заказы' : 'Войти или создать аккаунт';
  qs('logoutMenuItem').style.display = state.user ? 'flex' : 'none';
  const isAdmin = state.user && state.user.role === 'admin';
  qs('adminMenuItem').style.display = isAdmin ? 'flex' : 'none';
  qs('adminMenuItem').href = isAdmin ? '/admin' : '#';
  qs('adminMenuItem').querySelector('small').textContent = 'Товары, категории, заказы';
}
async function loadCategories(){
  try{state.categories = await api('/categories')}catch(e){state.categories = []}
}
async function loadBrands(){
  try{
    const data = await api('/brands');
    state.brands = data.items || [];
  }catch(e){
    state.brands = [];
  }
  renderBrandPills();
}
async function loadProducts(){
  const params = new URLSearchParams({page:state.page, per_page:state.perPage, sort:qs('sortSelect').value});
  const search = qs('searchInput').value.trim();
  const category = qs('categoryFilter').value;
  const brand = qs('brandFilter').value;
  if(search) params.set('search', search);
  if(category) params.set('category', category);
  if(brand) params.set('brand', brand);
  const data = await api('/products?'+params.toString());
  state.products = data.items || [];
  state.total = Number(data.total || 0);
  state.totalPages = Number(data.total_page || 1);
  await refreshFavoriteProducts();
  renderHome();
  renderFavorites();
  renderRecommendations();
}
async function loadCart(){
  state.cart = await api('/cart');
  renderCart();
  updateBadges();
}
async function refreshStoredOrders(){
  if(!state.orders.length) return;
  const refreshed = [];
  for(const saved of state.orders){
    const id = saved.order_id || saved.id;
    if(!id){
      refreshed.push(saved);
      continue;
    }
    try{
      const latest = await api('/order/' + encodeURIComponent(id));
      refreshed.push(Object.assign({}, saved, latest, {
        order_id: latest.id || id,
        created_at: formatDate(latest.created_at || saved.created_at),
        items: latest.items && latest.items.length ? latest.items : (saved.items || [])
      }));
    }catch(e){
      refreshed.push(saved);
    }
  }
  state.orders = refreshed;
  saveOrders();
}
async function loadUserOrders(){
  if(!state.user) return;
  try{
    const data = await api('/orders/my');
    if(data && Array.isArray(data.items)){
      state.orders = data.items.map(order => Object.assign({order_id: order.id}, order));
      saveOrders();
    }
  }catch(e){}
}
async function refreshFavoriteProducts(){
  if(!state.favorites.length) return;
  const seen = new Set();
  const productsById = new Map(state.products.map(p => [Number(p.id), p]));
  const nextFavorites = [];
  for(const rawId of state.favorites){
    const id = Number(rawId);
    if(!id || seen.has(id)) continue;
    seen.add(id);
    if(productsById.has(id)){
      nextFavorites.push(id);
      continue;
    }
    try{
      const product = await api('/products/' + encodeURIComponent(id));
      state.products.push(product);
      productsById.set(id, product);
      nextFavorites.push(id);
    }catch(e){}
  }
  if(nextFavorites.length !== state.favorites.length || nextFavorites.some((id, idx) => Number(state.favorites[idx]) !== id)){
    state.favorites = nextFavorites;
    saveFavorites();
  }
}
function setBadge(id, value){
  const el = qs(id);
  if(!el) return;
  const n = Number(value || 0);
  el.textContent = n > 99 ? '99+' : String(n);
  el.style.display = n > 0 ? (el.classList.contains('menu-badge') ? 'grid' : 'flex') : 'none';
}
function updateBadges(){
  const count = (state.cart.items || []).reduce((sum,it)=>sum+Number(it.quantity||0),0);
  setBadge('cartBadge', count);
  setBadge('menuCartBadge', count);
  qs('cartTitleCount').textContent = count;
  qs('summaryCount').textContent = count + (count === 1 ? ' товар' : ' товаров');
  setBadge('favBadge', state.favorites.length);
  setBadge('menuFavBadge', state.favorites.length);
  qs('favoritesTitleCount').textContent = state.favorites.length;
  setBadge('ordersBadge', state.orders.filter(order => !isFinishedOrder(order.status)).length);
}
function renderHome(){
  const grid = qs('productsGrid');
  if(!grid) return;
  grid.innerHTML = state.products.length ? state.products.map(productCard).join('') : emptyProducts('Товары не найдены','<a class="blue-button" href="/">Сбросить поиск</a>');
  qs('resultInfo').textContent = 'Найдено товаров: ' + state.total;
}
function renderFavorites(){
  const grid = qs('favoritesGrid');
  if(!grid) return;
  let items = state.products.filter(p => isFavorite(p.id));
  const search = currentRoute() === 'favorites' ? qs('searchInput').value.trim().toLowerCase() : '';
  if(search) items = items.filter(p => String(p.name || '').toLowerCase().includes(search));
  const sort = qs('favoriteSort').value;
  if(sort === 'price_asc') items.sort((a,b)=>Number(a.price)-Number(b.price));
  if(sort === 'price_desc') items.sort((a,b)=>Number(b.price)-Number(a.price));
  grid.innerHTML = items.length ? items.map(productCard).join('') : emptyProducts('В избранном пока пусто','<a class="blue-button" href="/">К покупкам</a>');
}
function renderRecommendations(){
  const grids = [qs('ordersRecommendGrid')].filter(Boolean);
  grids.forEach(grid => {
    grid.innerHTML = state.products.length ? state.products.slice(0,8).map(productCard).join('') : emptyProducts('Скоро появятся рекомендации');
  });
}
function renderBrandPills(){
  const box = qs('brandPills');
  if(!box) return;
  if(!state.brands.length){
    box.innerHTML = '';
    return;
  }
  const current = qs('brandFilter').value || '';
  const buttons = ['<button class="brand-pill '+(!current?'active':'')+'" data-brand="" type="button">Все бренды</button>'];
  state.brands.forEach(brand => {
    buttons.push('<button class="brand-pill '+(brand === current ? 'active' : '')+'" data-brand="'+esc(brand)+'" type="button">'+esc(brand)+'</button>');
  });
  box.innerHTML = buttons.join('');
}
function renderCart(){
  const items = state.cart.items || [];
  const count = items.reduce((sum,it)=>sum+Number(it.quantity||0),0);
  const subtotal = Number(state.cart.total || 0);
  const visualDiscount = Math.round(subtotal * 0.12);
  const total = Math.max(0, subtotal - visualDiscount);
  qs('summarySubtotal').textContent = money(subtotal);
  qs('summaryDiscount').textContent = '- ' + money(visualDiscount);
  qs('summaryTotal').textContent = money(total);
  qs('cartSaleHint').textContent = count ? count + ' товаров скоро подорожают' : 'Добавьте товары, чтобы поймать скидку';
  qs('checkoutOpen').disabled = count === 0;
  updateBadges();
  const box = qs('cartPageItems');
  if(!box) return;
  if(!items.length){
    box.innerHTML = emptyProducts('Корзина пустая','<a class="blue-button" href="/">К покупкам</a>');
    return;
  }
  box.innerHTML = items.map(it => {
    const p = it.product || {};
    const qty = Number(it.quantity || 0);
    return '<div class="cart-row">'+
      '<span class="check">✓</span>'+
      (p.image_url ? '<img src="'+esc(p.image_url)+'" alt="'+esc(p.name)+'">' : '<div class="cart-picture">SHOP</div>')+
      '<div><div class="cart-name">'+esc(p.name)+'</div><div class="cart-actions"><button class="mini-action" data-fav="'+p.id+'" type="button">'+(isFavorite(p.id)?'В избранном':'В избранное')+'</button><button class="mini-action" data-remove="'+p.id+'" type="button">Удалить</button><button class="mini-action" data-add="'+p.id+'" type="button">Купить еще</button></div><div class="muted" style="margin-top:8px">'+esc(p.color || '')+' '+esc(p.material || '')+'</div></div>'+
      '<div><div class="cart-price">'+money(p.price)+'</div><div class="cart-old">'+money(oldPrice(p))+'</div></div>'+
      '<div><div class="quantity-box"><button data-qty="'+p.id+'" data-value="'+(qty-1)+'" type="button">−</button><span>'+qty+'</span><button data-qty="'+p.id+'" data-value="'+(qty+1)+'" type="button">+</button></div>'+(qty > 1 ? '<div class="limit-note">Количество ограничено</div>' : '')+'</div>'+
    '</div>';
  }).join('');
}
function renderOrders(){
  const box = qs('ordersContent');
  if(!box) return;
  syncOrderTabs();
  const filtered = state.orders.filter(order => state.orderTab === 'finished' ? isFinishedOrder(order.status) : !isFinishedOrder(order.status));
  if(!filtered.length){
    const title = state.orderTab === 'finished' ? 'Завершенных заказов пока нет' : 'У вас пока нет актуальных заказов';
    const text = state.orderTab === 'finished' ? 'Когда заказ будет доставлен или отменен, он появится здесь.' : 'Когда появятся, будут отображаться здесь.';
    box.innerHTML = '<div class="orders-empty"><div class="empty-state" style="width:100%;background:transparent"><div class="empty-illustration"></div><h2>'+title+'</h2><p>'+text+'</p><a class="blue-button" href="/">К покупкам</a></div></div>';
    return;
  }
  box.innerHTML = '<div class="orders-list">'+filtered.map(order => {
    const orderId = order.order_id || order.id;
    const items = (order.items || []).map(it => esc((it.product || {}).name || it.product_name || 'Товар')).join(', ');
    const place = order.delivery_type === 'pickup' ? order.pickup_point : order.address;
    return '<article class="order-card"><div class="order-card-head"><div><b>Заказ №'+orderId+'</b><br><span class="muted">'+esc(formatDate(order.created_at))+'</span></div><span class="status-chip">'+statusLabel(order.status)+'</span></div><div>'+items+'</div>'+(place ? '<p class="muted">'+esc(place)+'</p>' : '')+'<h3>'+money(order.total)+'</h3></article>';
  }).join('')+'</div>';
}
function isFinishedOrder(status){
  return status === 'delivered' || status === 'completed' || status === 'cancelled';
}
function syncOrderTabs(){
  document.querySelectorAll('[data-order-tab]').forEach(btn => btn.classList.toggle('active', btn.dataset.orderTab === state.orderTab));
}
function setOrderTab(tab){
  state.orderTab = tab === 'finished' ? 'finished' : 'active';
  renderOrders();
}
function statusLabel(status){
  if(status === 'pending') return 'Заказ принят, ожидайте оператора';
  if(status === 'processing') return 'В обработке';
  if(status === 'shipped') return 'В пути';
  if(status === 'delivered' || status === 'completed') return 'Заказ доставлен';
  if(status === 'cancelled') return 'Заказ отменён';
  return 'Заказ принят, ожидайте оператора';
}
function showRoute(){
  const route = currentRoute();
  document.querySelectorAll('.page-view').forEach(v=>v.classList.remove('active'));
  const page = qs(route + 'Page');
  if(page) page.classList.add('active');
  document.querySelectorAll('.nav-link').forEach(a=>a.classList.toggle('active', a.dataset.route === route));
  qs('searchInput').placeholder = route === 'favorites' ? 'Найти в избранном' : 'Искать на Shop Go';
  if(route === 'orders') renderOrders();
  if(route === 'cart') renderCart();
  if(route === 'favorites') renderFavorites();
}
function setCategory(category){
  qs('categoryFilter').value = category || '';
  document.querySelectorAll('[data-category]').forEach(btn=>btn.classList.toggle('active', String(btn.dataset.category || '') === String(category || '')));
  state.page = 1;
  loadProducts().catch(err=>toast(err.message));
}
function setBrand(brand){
  qs('brandFilter').value = brand || '';
  renderBrandPills();
  state.page = 1;
  loadProducts().catch(err=>toast(err.message));
}
function submitSearch(){
  const route = currentRoute();
  const term = qs('searchInput').value.trim();
  if(route === 'favorites'){renderFavorites(); return}
  if(route !== 'home'){
    location.href = term ? '/?search=' + encodeURIComponent(term) : '/';
    return;
  }
  state.page = 1;
  loadProducts().catch(err=>toast(err.message));
}
function applyCatalogSelection(category, query, brand){
  const cat = category || '';
  const term = query || '';
  const brandValue = brand || '';
  closeCatalog();
  if(currentRoute() === 'home'){
    qs('searchInput').value = term;
    qs('categoryFilter').value = cat;
    qs('brandFilter').value = brandValue;
    document.querySelectorAll('[data-category]').forEach(btn=>btn.classList.toggle('active', String(btn.dataset.category || '') === String(cat)));
    renderBrandPills();
    state.page = 1;
    loadProducts().catch(err=>toast(err.message));
    return;
  }
  const params = new URLSearchParams();
  if(term) params.set('search', term);
  if(cat) params.set('category', cat);
  if(brandValue) params.set('brand', brandValue);
  location.href = params.toString() ? '/?' + params.toString() : '/';
}
async function addToCart(id){
  await api('/cart/add',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({product_id:Number(id), quantity:1})});
  await loadCart();
  toast('Товар добавлен в корзину');
}
async function setQty(id, qty){
  await api('/cart/update',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({product_id:Number(id), quantity:Number(qty)})});
  await loadCart();
}
function toggleFavorite(id){
  const n = Number(id);
  state.favorites = isFavorite(n) ? state.favorites.filter(x => Number(x) !== n) : [n].concat(state.favorites);
  saveFavorites();
  updateBadges();
  renderHome();
  renderFavorites();
  renderCart();
}
function openCatalog(){
  const panel = qs('catalogPanel');
  const open = !panel.classList.contains('open');
  panel.classList.toggle('open', open);
  qs('catalogToggle').classList.toggle('open', open);
  panel.setAttribute('aria-hidden', open ? 'false' : 'true');
}
function closeCatalog(){
  qs('catalogPanel').classList.remove('open');
  qs('catalogToggle').classList.remove('open');
  qs('catalogPanel').setAttribute('aria-hidden','true');
}
function openCheckout(){
  const items = state.cart.items || [];
  if(!items.length){toast('Корзина пустая');return}
  qs('modalCover').style.display = 'block';
  qs('checkoutModal').style.display = 'block';
  qs('checkoutFormWrap').style.display = 'block';
  qs('checkoutSuccess').style.display = 'none';
  if(state.user && !qs('coName').value) qs('coName').value = state.user.username;
  loadPickupPoints().catch(err=>toast(err.message));
}
function closeCheckout(){
  qs('modalCover').style.display = 'none';
  qs('checkoutModal').style.display = 'none';
}
async function loadPickupPoints(){
  const data = await api('/pickup-points');
  state.pickupPoints = data.points || [];
  if(!state.selectedPickupId && state.pickupPoints.length) state.selectedPickupId = state.pickupPoints[0].id;
  renderPickupPoints();
}
function renderPickupPoints(){
  const select = qs('pickupSelect');
  if(!state.pickupPoints.length){
    select.innerHTML = '<option value="">Пункты выдачи не найдены</option>';
    qs('pickupInfo').textContent = 'Пункты выдачи пока недоступны.';
    return;
  }
  select.innerHTML = state.pickupPoints.map(p => '<option value="'+p.id+'" '+(Number(p.id)===Number(state.selectedPickupId)?'selected':'')+'>'+esc(p.name)+' · '+esc(p.address)+'</option>').join('');
  const point = state.pickupPoints.find(p => Number(p.id) === Number(state.selectedPickupId));
  qs('pickupInfo').innerHTML = point ? '<b>'+esc(point.name)+'</b><br>'+esc(point.address)+'<br>'+esc(point.working_hours || 'Время работы не указано')+' · '+esc(point.phone || 'телефон не указан') : 'Выберите пункт выдачи.';
}
async function submitOrder(e){
  e.preventDefault();
  const delivery = document.querySelector('input[name="delivery"]:checked').value;
  const name = qs('coName').value.trim();
  const email = qs('coEmail').value.trim();
  const address = qs('coAddress').value.trim();
  if(!name || !email){toast('Заполните имя и email'); return}
  if(delivery === 'courier' && !address){toast('Введите адрес доставки'); return}
  if(delivery === 'pickup' && !state.selectedPickupId){toast('Выберите пункт выдачи'); return}
  const snapshot = JSON.parse(JSON.stringify(state.cart));
  const btn = qs('submitOrder');
  btn.disabled = true;
  btn.textContent = 'Оформляем...';
  try{
    const order = await api('/order',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({
      name, email, phone: qs('coPhone').value.trim(), address: delivery === 'courier' ? address : '',
      delivery_type: delivery,
      pickup_point_id: delivery === 'pickup' ? Number(state.selectedPickupId) : 0,
      delivery_lat: delivery === 'courier' ? 55.7558 : 0,
      delivery_lng: delivery === 'courier' ? 37.6223 : 0
    })});
    const saved = Object.assign({}, order, {items:snapshot.items || [], created_at:new Date().toLocaleString('ru-RU')});
    state.orders = [saved].concat(state.orders).slice(0,20);
    saveOrders();
    await loadCart();
    renderOrders();
    qs('checkoutFormWrap').style.display = 'none';
    qs('checkoutSuccess').style.display = 'block';
    qs('successText').innerHTML = 'Ваш заказ принят, ожидайте звонка оператора.<br>Номер заказа: <b>№'+order.order_id+'</b><br>Сумма: <b>'+money(order.total)+'</b>';
  }catch(err){
    toast(err.message);
  }finally{
    btn.disabled = false;
    btn.textContent = 'Завершить заказ';
  }
}

document.addEventListener('click', e => {
  const add = e.target.closest('[data-add]');
  if(add){addToCart(add.dataset.add).catch(err=>toast(err.message)); return}
  const fav = e.target.closest('[data-fav]');
  if(fav){toggleFavorite(fav.dataset.fav); return}
  const qty = e.target.closest('[data-qty]');
  if(qty){setQty(qty.dataset.qty, qty.dataset.value).catch(err=>toast(err.message)); return}
  const rem = e.target.closest('[data-remove]');
  if(rem){setQty(rem.dataset.remove, 0).catch(err=>toast(err.message)); return}
  const category = e.target.closest('[data-category]');
  if(category){setCategory(category.dataset.category); return}
  const brand = e.target.closest('[data-brand]');
  if(brand){setBrand(brand.dataset.brand); return}
  const orderTab = e.target.closest('[data-order-tab]');
  if(orderTab){setOrderTab(orderTab.dataset.orderTab); return}
  if(!e.target.closest('.account-wrap')) qs('accountMenu').classList.remove('open');
});
qs('catalogToggle').addEventListener('click', openCatalog);
qs('accountBtn').addEventListener('click', () => qs('accountMenu').classList.toggle('open'));
qs('searchBtn').addEventListener('click', submitSearch);
qs('searchInput').addEventListener('keydown', e => {
  if(e.key === 'Enter'){
    submitSearch();
  }
});
qs('sortSelect').addEventListener('change', () => {state.page=1; loadProducts().catch(err=>toast(err.message))});
qs('favoriteSort').addEventListener('change', renderFavorites);
qs('checkoutOpen').addEventListener('click', openCheckout);
qs('closeCheckout').addEventListener('click', closeCheckout);
qs('modalCover').addEventListener('click', closeCheckout);
qs('checkoutForm').addEventListener('submit', submitOrder);
qs('pickupSelect').addEventListener('change', e => {state.selectedPickupId = Number(e.target.value); renderPickupPoints()});
document.querySelectorAll('input[name="delivery"]').forEach(r => r.addEventListener('change', () => {
  const delivery = document.querySelector('input[name="delivery"]:checked').value;
  qs('pickupBox').style.display = delivery === 'pickup' ? 'block' : 'none';
  qs('courierBox').style.display = delivery === 'courier' ? 'block' : 'none';
}));
qs('logoutMenuItem').addEventListener('click', async e => {
  e.preventDefault();
  await api('/auth/logout',{method:'POST'});
  location.href = '/auth';
});
window.addEventListener('keydown', e => {if(e.key === 'Escape'){closeCatalog(); closeCheckout(); qs('accountMenu').classList.remove('open')}});

(async function init(){
  const params = new URLSearchParams(location.search);
  qs('searchInput').value = params.get('search') || '';
  qs('categoryFilter').value = params.get('category') || '';
  qs('brandFilter').value = params.get('brand') || '';
  document.querySelectorAll('[data-category]').forEach(btn=>btn.classList.toggle('active', String(btn.dataset.category || '') === String(qs('categoryFilter').value || '')));
  showRoute();
  await Promise.all([loadUser(), loadCategories(), loadBrands(), loadCart()]);
  await loadUserOrders();
  await refreshStoredOrders();
  await loadProducts();
  renderOrders();
  updateBadges();
})().catch(err=>toast(err.message));
</script>
</body>
</html>`

const shopProductHTML = `<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.product.Name}}</title>
<style>
:root{--blue:#005bff;--pink:#f91155;--green:#00b75a;--ink:#001a34;--muted:#667b90;--soft:#f1f4f9;--line:#dfe7f1}
*{box-sizing:border-box}body{margin:0;background:var(--soft);color:var(--ink);font-family:Inter,Segoe UI,Arial,sans-serif}button,input{font:inherit}button{cursor:pointer}
.top{position:sticky;top:0;background:#fff;border-radius:0 0 28px 28px;box-shadow:0 8px 24px rgba(0,26,52,.05);z-index:20}.bar{max-width:1410px;margin:0 auto;padding:10px 22px;display:grid;grid-template-columns:auto auto 1fr auto;gap:14px;align-items:center}.logo{font-size:42px;line-height:1;font-weight:900;color:var(--blue);text-decoration:none}.catalog{height:46px;border:0;border-radius:11px;background:var(--blue);color:#fff;padding:0 18px;font-weight:900}.search{height:46px;border:2px solid var(--blue);border-radius:12px;display:flex;overflow:hidden}.search input{border:0;outline:0;flex:1;padding:0 14px}.search a{background:var(--blue);color:#fff;display:grid;place-items:center;text-decoration:none;padding:0 20px;font-weight:900}.nav{display:flex;gap:16px}.nav a{color:#7a8b9e;text-decoration:none;font-weight:800;font-size:13px;text-align:center}.shell{max-width:1410px;margin:0 auto;padding:28px 22px}.product{display:grid;grid-template-columns:minmax(320px,560px) 1fr 330px;gap:28px;align-items:start}.gallery,.info,.buybox{background:#fff;border-radius:22px}.gallery{padding:18px}.image{aspect-ratio:1/1.08;border-radius:16px;background:#fff;display:grid;place-items:center;overflow:hidden}.image img{width:100%;height:100%;object-fit:contain;padding:10px;background:#fff}.empty{font-size:38px;color:#8aa0b8;font-weight:900}.info{padding:26px}.crumb{color:var(--muted);margin-bottom:12px}.title{margin:0 0 16px;font-size:34px;line-height:1.08}.rating{display:flex;gap:10px;color:var(--muted);font-size:14px;margin-bottom:18px}.star{color:#f5a000;font-weight:900}.desc{color:#3f5267;line-height:1.55;white-space:pre-wrap}.specs{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:10px;margin-top:20px}.spec{background:#f5f8fc;border-radius:14px;padding:12px}.spec span{display:block;color:var(--muted);font-size:13px;margin-bottom:5px}.buybox{padding:24px;position:sticky;top:100px}.price{font-size:34px;font-weight:900;color:var(--pink);margin-bottom:4px}.old{color:#8b9aaa;text-decoration:line-through;font-weight:800}.green{color:var(--green);font-weight:900;margin:12px 0}.blue-button,.light-button{height:48px;border:0;border-radius:12px;padding:0 18px;font-weight:900;text-decoration:none;display:flex;align-items:center;justify-content:center;width:100%;margin-top:10px}.blue-button{background:var(--blue);color:#fff}.blue-button:disabled{background:#d7e0ea;color:#7f91a4}.light-button{background:#eef6ff;color:var(--blue)}.toast{position:fixed;left:50%;bottom:24px;transform:translateX(-50%);background:#001a34;color:#fff;border-radius:999px;padding:12px 18px;display:none;font-weight:900}
@media(max-width:1000px){.product{grid-template-columns:1fr}.buybox{position:static}.bar{grid-template-columns:1fr auto}.search{grid-column:1/-1}.nav{grid-column:1/-1;justify-content:space-around}.logo{font-size:34px}.specs{grid-template-columns:1fr}}
</style>
</head>
<body>
<header class="top"><div class="bar"><a class="logo" href="/">SHOPGO</a><a class="catalog" href="/">Каталог</a><div class="search"><input placeholder="Искать на Shop Go"><a href="/">Найти</a></div><nav class="nav"><a href="/orders">Заказы</a><a href="/favorites">Избранное</a><a href="/cart/view">Корзина</a></nav></div></header>
<main class="shell">
  <section class="product">
    <div class="gallery"><div class="image">{{if .product.ImageURL}}<img src="{{.product.ImageURL}}" alt="{{.product.Name}}">{{else}}<div class="empty">SHOP</div>{{end}}</div></div>
    <div class="info">
      <div class="crumb">Каталог / товар</div>
      <h1 class="title">{{.product.Name}}</h1>
      <div class="rating"><span class="star">★ 4.9</span><span>проверенный товар</span><span>оригинальный товар</span></div>
      <p class="desc">{{.product.Description}}</p>
      <div class="specs">
        <div class="spec"><span>Бренд</span><b>{{.product.Brand}}</b></div>
        <div class="spec"><span>В наличии</span><b>{{.product.Stock}} шт</b></div>
        <div class="spec"><span>Состояние</span><b>{{.product.Condition}}</b></div>
        <div class="spec"><span>Цвет</span><b>{{.product.Color}}</b></div>
        <div class="spec"><span>Материал</span><b>{{.product.Material}}</b></div>
        <div class="spec"><span>Страна</span><b>{{.product.Country}}</b></div>
      </div>
    </div>
    <aside class="buybox">
      <div class="price">{{printf "%.0f" .product.Price}} ₽</div>
      <div class="old">Старая цена выше</div>
      <div class="green">Доставка завтра в пункт выдачи</div>
      <button class="blue-button" id="addBtn" {{if le .product.Stock 0}}disabled{{end}}>{{if le .product.Stock 0}}Нет в наличии{{else}}Добавить в корзину{{end}}</button>
      <button class="light-button" id="favBtn" type="button">В избранное</button>
      <a class="light-button" href="/">Продолжить покупки</a>
    </aside>
  </section>
</main>
<div class="toast" id="toast"></div>
<script>
const productId = {{.product.ID}};
let currentUser = null;
const qs = id => document.getElementById(id);
function toast(text){const el=qs('toast');el.textContent=text;el.style.display='block';clearTimeout(window.t);window.t=setTimeout(()=>el.style.display='none',1800)}
async function api(url, options){const r=await fetch(url,Object.assign({credentials:'include'},options||{}));const text=await r.text();let p={};try{p=text?JSON.parse(text):{}}catch(e){p={error:text}}if(!r.ok)throw new Error(p.error||text||'Ошибка');return p}
function favoriteKey(){return currentUser && currentUser.id ? 'shopgo:user:'+currentUser.id+':favorites' : 'shopgo:guest:favorites'}
async function loadCurrentUser(){try{const data=await api('/auth/me');currentUser=data.authenticated?data.user:null}catch(e){currentUser=null}}
qs('addBtn').addEventListener('click',async()=>{try{await api('/cart/add',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({product_id:productId,quantity:1})});toast('Товар добавлен в корзину')}catch(e){toast(e.message)}})
qs('favBtn').addEventListener('click',async()=>{if(currentUser===null)await loadCurrentUser();let fav=[];try{fav=JSON.parse(localStorage.getItem(favoriteKey()))||[]}catch(e){};fav=fav.map(Number).includes(Number(productId))?fav.filter(x=>Number(x)!==Number(productId)):[productId].concat(fav);localStorage.setItem(favoriteKey(),JSON.stringify(fav));toast('Избранное обновлено')})
loadCurrentUser();
</script>
</body>
</html>`
