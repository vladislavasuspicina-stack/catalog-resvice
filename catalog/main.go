package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

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
	r.GET("/brands", GetBrands)

	// Cart & Order
	r.POST("/cart/add", AddToCart)
	r.POST("/cart/update", UpdateCartQuantity)
	r.POST("/cart/remove", RemoveFromCart)
	r.GET("/cart", GetCart)
	r.POST("/order", CreateOrderHandler)
	r.GET("/order/:id", GetOrderStatus)
	r.GET("/orders/my", UserListOrders)
	r.GET("/pickup-points", GetPickupPoints)

	// Customer auth
	r.POST("/auth/register", UserRegister)
	r.POST("/auth/login", UserLogin)
	r.POST("/auth/logout", UserLogout)
	r.GET("/auth/me", UserMe)

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
		admin.DELETE("/orders/:id", AdminDeleteOrder)
		// users
		admin.GET("/users", AdminListUsers)
		admin.PUT("/users/:id", AdminUpdateUser)
		admin.DELETE("/users/:id", AdminDeleteUser)
	}

	// Web UI
	t := template.Must(template.New("index").Parse(shopIndexHTML))
	template.Must(t.New("admin").Parse(adminHTML))
	template.Must(t.New("product").Parse(shopProductHTML))
	template.Must(t.New("auth").Parse(authHTML))
	r.SetHTMLTemplate(t)

	r.GET("/auth", func(c *gin.Context) {
		next := sanitizeNextPath(c.DefaultQuery("next", "/"))
		if _, err := currentUserFromCookie(c); err == nil {
			c.Redirect(http.StatusFound, next)
			return
		}
		c.HTML(http.StatusOK, "auth", gin.H{"next": next})
	})

	renderShopPage := func(c *gin.Context) {
		var categories []Category
		var products []Product
		DB.Find(&categories)
		DB.Preload("Categories").Order("id asc").Limit(12).Find(&products)
		for i := range products {
			enrichProduct(&products[i])
		}
		c.HTML(http.StatusOK, "index", gin.H{"categories": categories, "products": products})
	}
	r.GET("/", renderShopPage)
	r.GET("/shop", renderShopPage)
	r.GET("/favorites", renderShopPage)
	r.GET("/orders", renderShopPage)
	r.GET("/catalog", renderShopPage)
	r.GET("/cart/view", renderShopPage)

	r.GET("/product/:id", func(c *gin.Context) {
		id := c.Param("id")
		var p Product
		if err := DB.Preload("Categories").First(&p, id).Error; err != nil {
			c.String(http.StatusNotFound, "product not found")
			return
		}
		enrichProduct(&p)
		c.HTML(http.StatusOK, "product", gin.H{"product": p})
	})

	r.GET("/admin", func(c *gin.Context) {
		user, err := currentUserFromCookie(c)
		if err != nil {
			c.Redirect(http.StatusFound, "/auth?next="+url.QueryEscape("/admin"))
			return
		}
		if normalizeUserRole(user.Role) != "admin" {
			c.String(http.StatusForbidden, "доступ только для админа")
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

func sanitizeNextPath(next string) string {
	next = strings.TrimSpace(next)
	if next == "" || !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return "/"
	}
	if next == "/" || next == "/shop" || next == "/cart/view" || next == "/favorites" || next == "/orders" || next == "/catalog" || next == "/admin" || strings.HasPrefix(next, "/product/") {
		return next
	}
	return "/"
}
