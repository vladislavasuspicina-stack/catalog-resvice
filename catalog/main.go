package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	r := gin.Default()
	r.Static("/static", "./static")

	// Public API
	r.GET("/products", GetProducts)
	r.GET("/products/:id", GetProduct)
	r.GET("/categories", GetCategories)

	// Cart & Order
	r.POST("/cart/add", AddToCart)
	r.POST("/cart/update", UpdateCartQuantity)
	r.POST("/cart/remove", RemoveFromCart)
	r.GET("/cart", GetCart)
	r.POST("/order", CreateOrderHandler)

	// Admin auth
	r.POST("/admin/login", AdminLogin)
	r.POST("/admin/logout", AdminLogout)

	// Admin API (protected)
	admin := r.Group("/admin/api", AdminRequired())
	{
		// categories
		admin.POST("/categories", CreateCategory)
		admin.GET("/categories", GetCategories)
		admin.PUT("/categories/:id", UpdateCategory)
		admin.DELETE("/categories/:id", DeleteCategory)
		// products
		admin.POST("/products", CreateProduct)
		admin.PUT("/products/:id", UpdateProduct)
		admin.DELETE("/products/:id", DeleteProduct)
		// upload
		admin.POST("/upload", UploadImageHandler)
		// orders
		admin.GET("/orders", AdminListOrders)
		admin.GET("/orders/:id", AdminGetOrder)
		admin.PUT("/orders/:id/status", AdminUpdateOrderStatus)
	}

	// Web UI
	t := template.Must(template.New("index").Parse(indexHTML))
	template.Must(t.New("admin").Parse(adminHTML))
	template.Must(t.New("product").Parse(productHTML))
	template.Must(t.New("cart").Parse(cartHTML))
	r.SetHTMLTemplate(t)

	r.GET("/", func(c *gin.Context) {
		var categories []Category
		var products []Product
		DB.Find(&categories)
		DB.Order("id asc").Limit(12).Find(&products)
		c.HTML(http.StatusOK, "index", gin.H{"categories": categories, "products": products})
	})

	r.GET("/product/:id", func(c *gin.Context) {
		id := c.Param("id")
		var p Product
		if err := DB.First(&p, id).Error; err != nil {
			c.String(http.StatusNotFound, "Товар не найден")
			return
		}
		c.HTML(http.StatusOK, "product", gin.H{"product": p})
	})

	// Cart page (separate view)
	r.GET("/cart/view", func(c *gin.Context) {
		c.HTML(http.StatusOK, "cart", nil)
	})

	r.GET("/admin", func(c *gin.Context) {
		_, err := c.Cookie(AdminCookieName)
		if err != nil {
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.String(http.StatusOK, adminLoginPage())
			return
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		if err := t.ExecuteTemplate(c.Writer, "admin", nil); err != nil {
			c.String(http.StatusInternalServerError, "template error")
		}
	})

	fmt.Println("Server running on :8080")
	r.Run(":8080")
}

// ==================== SIMPLE LOGIN PAGE ====================
func adminLoginPage() string {
	return `<!doctype html><html><head><meta charset="utf-8"><title>Admin login</title></head><body>
<h2>Вход в админку</h2>
<form id="f">
<input id="u" placeholder="username"><br><br>
<input id="p" placeholder="password" type="password"><br><br>
<button type="button" onclick="login()">Войти</button>
</form>
<script>
async function login(){
  const u=document.getElementById('u').value;
  const p=document.getElementById('p').value;
  const r=await fetch('/admin/login',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({username:u,password:p})});
  if(r.ok){ location.href='/admin'; } else { alert('auth failed'); }
}
</script>
</body></html>`
}
