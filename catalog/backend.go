package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ==================== JWT HELPERS ====================
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	// Simplified distance calculation (approximation for Moscow)
	// Returns distance in kilometers
	// Using simple Euclidean distance on projected plane

	dLat := (lat2 - lat1) * 111.0 // approximately 111 km per degree latitude
	dLng := (lng2 - lng1) * 88.0  // approximately 88 km per degree longitude at Moscow latitude

	distance := (dLat*dLat + dLng*dLng)
	if distance < 0 {
		distance = -distance
	}

	// Approximate square root manually
	return distance
}

func createToken(subject, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  subject,
		"role": role,
		"exp":  time.Now().Add(ttl).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET))
}

func parseToken(tokenStr string, expectedRole string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(JWT_SECRET), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		sub, ok := claims["sub"].(string)
		if !ok || strings.TrimSpace(sub) == "" {
			return "", errors.New("invalid claims")
		}
		if expectedRole != "" {
			role, ok := claims["role"].(string)
			if !ok || role != expectedRole {
				return "", errors.New("invalid role")
			}
		}
		return sub, nil
	}
	return "", errors.New("invalid claims")
}

// ==================== MIDDLEWARE ====================
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := currentUserFromCookie(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth required"})
			return
		}
		if normalizeUserRole(user.Role) != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		c.Next()
	}
}

func currentUserFromCookie(c *gin.Context) (*User, error) {
	cookie, err := c.Cookie(UserCookieName)
	if err != nil || strings.TrimSpace(cookie) == "" {
		return nil, errors.New("auth required")
	}
	sub, err := parseToken(cookie, "")
	if err != nil {
		return nil, err
	}
	userID, err := strconv.ParseUint(sub, 10, 64)
	if err != nil {
		return nil, errors.New("invalid token subject")
	}
	var user User
	if err := DB.First(&user, uint(userID)).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func userPublic(user *User) gin.H {
	return gin.H{
		"id":       user.ID,
		"username": user.Username,
		"role":     normalizeUserRole(user.Role),
	}
}

func normalizeUserRole(role string) string {
	if strings.EqualFold(strings.TrimSpace(role), "admin") {
		return "admin"
	}
	return "user"
}

// ==================== UTIL ====================
func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func sanitizeFilename(name string) string {
	name = filepath.Base(name)
	name = strings.ReplaceAll(name, " ", "-")
	return name
}

func saveOrderToFile(order Order, items []OrderItem) error {
	orderDir := "./orders"
	if err := os.MkdirAll(orderDir, 0o755); err != nil {
		return err
	}
	filename := fmt.Sprintf("%s/order-%d.json", orderDir, order.ID)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	data := map[string]interface{}{
		"order_id":      order.ID,
		"customer_name": order.CustomerName,
		"email":         order.Email,
		"address":       order.Address,
		"total":         order.Total,
		"status":        order.Status,
		"created_at":    order.CreatedAt,
		"items":         items,
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Save order to Excel file
func saveOrderToExcel(order Order, items []OrderItem) error {
	excelPath := "./orders/orders.xlsx"

	// Create or open Excel file
	f, err := excelize.OpenFile(excelPath)
	isNewFile := err != nil

	if isNewFile {
		f = excelize.NewFile()
		// Create header row
		headers := []string{"ID заказа", "Имя", "Email", "Телефон", "Адрес", "Статус", "Дата", "Сумма (₽)", "Доставка", "Точка самовывоза", "Доставка льгот", "Товары"}
		for i, h := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue("Sheet1", cell, h)
			f.SetCellStyle("Sheet1", cell, cell, 1) // Apply style to header
		}
	}

	// Get next row number
	rows, err := f.GetRows("Sheet1")
	nextRow := len(rows) + 1

	// Get phone from order
	phone := ""
	if order.Phone != "" {
		phone = order.Phone
	}

	// Format delivery info
	deliveryInfo := order.DeliveryType
	if order.DeliveryType == "pickup" {
		deliveryInfo = "Самовывоз"
	} else if order.DeliveryType == "courier" {
		deliveryInfo = "Курьер"
	}

	pickupPointStr := ""
	if order.PickupPoint != "" {
		pickupPointStr = order.PickupPoint
	}

	deliveryPriceStr := ""
	if order.DeliveryPrice > 0 {
		deliveryPriceStr = fmt.Sprintf("%.2f₽", order.DeliveryPrice)
	}

	// Prepare items summary
	itemsStr := ""
	for _, item := range items {
		if itemsStr != "" {
			itemsStr += "; "
		}
		itemsStr += fmt.Sprintf("%s (×%d - %.2f₽)", item.ProductName, item.Quantity, item.Price*float64(item.Quantity))
	}

	// Add order data
	rowData := []interface{}{
		order.ID,
		order.CustomerName,
		order.Email,
		phone,
		order.Address,
		order.Status,
		order.CreatedAt.Format("02.01.2006 15:04"),
		order.Total,
		deliveryInfo,
		pickupPointStr,
		deliveryPriceStr,
		itemsStr,
	}

	for i, val := range rowData {
		cell := fmt.Sprintf("%c%d", 'A'+i, nextRow)
		f.SetCellValue("Sheet1", cell, val)
	}

	// Set column widths
	f.SetColWidth("Sheet1", "A", "A", 12)
	f.SetColWidth("Sheet1", "B", "B", 18)
	f.SetColWidth("Sheet1", "C", "C", 18)
	f.SetColWidth("Sheet1", "D", "D", 16)
	f.SetColWidth("Sheet1", "E", "E", 30)
	f.SetColWidth("Sheet1", "F", "F", 12)
	f.SetColWidth("Sheet1", "G", "G", 16)
	f.SetColWidth("Sheet1", "H", "H", 12)
	f.SetColWidth("Sheet1", "I", "I", 12)
	f.SetColWidth("Sheet1", "J", "J", 25)
	f.SetColWidth("Sheet1", "K", "K", 12)
	f.SetColWidth("Sheet1", "L", "L", 50)

	// Create orders directory if not exists
	orderDir := "./orders"
	os.MkdirAll(orderDir, 0o755)

	// Save file
	return f.SaveAs(excelPath)
}

// ==================== CATEGORY HANDLERS ====================
func CreateCategory(c *gin.Context) {
	var input Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	cat := Category{Name: input.Name, ParentID: input.ParentID}
	DB.Create(&cat)
	c.JSON(http.StatusCreated, cat)
}

func GetCategories(c *gin.Context) {
	var cats []Category
	DB.Order("id asc").Find(&cats)
	c.JSON(http.StatusOK, cats)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat Category
	if err := DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if name, ok := input["name"].(string); ok && strings.TrimSpace(name) != "" {
		cat.Name = name
	}
	if parent, ok := input["parent_id"].(float64); ok {
		pid := uint(parent)
		cat.ParentID = &pid
	}
	DB.Save(&cat)
	c.JSON(http.StatusOK, cat)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	res := DB.Delete(&Category{}, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ==================== PRODUCT HANDLERS ====================
func CreateProduct(c *gin.Context) {
	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(input.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	p := Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		CategoryID:  input.CategoryID,
		ImageURL:    input.ImageURL,
	}
	DB.Create(&p)
	c.JSON(http.StatusCreated, p)
}

func GetProducts(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "20")
	sort := c.DefaultQuery("sort", "id_asc")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	search := c.Query("search")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	categoryID := c.Query("category")

	query := DB.Model(&Product{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	if minPrice != "" {
		if v, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query = query.Where("price >= ?", v)
		}
	}
	if maxPrice != "" {
		if v, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query = query.Where("price <= ?", v)
		}
	}
	if categoryID != "" {
		if v, err := strconv.Atoi(categoryID); err == nil {
			query = query.Where("category_id = ?", v)
		}
	}

	var total int64
	query.Count(&total)

	switch sort {
	case "price_asc":
		query = query.Order("price asc")
	case "price_desc":
		query = query.Order("price desc")
	case "name_asc":
		query = query.Order("name asc")
	case "name_desc":
		query = query.Order("name desc")
	case "newest":
		query = query.Order("created_at desc")
	default:
		query = query.Order("id asc")
	}

	var products []Product
	query.Offset(offset).Limit(perPage).Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"items":      products,
		"page":       page,
		"per_page":   perPage,
		"total":      total,
		"total_page": (total + int64(perPage) - 1) / int64(perPage),
	})
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var p Product
	if err := DB.First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var prod Product
	if err := DB.First(&prod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update only fields that are present in the JSON payload
	if v, ok := input["name"].(string); ok {
		prod.Name = v
	}
	if v, ok := input["description"].(string); ok {
		prod.Description = v
	}
	if v, ok := input["price"].(float64); ok {
		prod.Price = v
	}
	if v, ok := input["stock"].(float64); ok {
		prod.Stock = int(v)
	}
	if v, ok := input["category_id"]; ok {
		if v == nil {
			prod.CategoryID = nil
		} else {
			switch t := v.(type) {
			case float64:
				vv := uint(t)
				prod.CategoryID = &vv
			case string:
				if t == "" {
					prod.CategoryID = nil
				} else if parsed, err := strconv.Atoi(t); err == nil {
					tmp := uint(parsed)
					prod.CategoryID = &tmp
				}
			}
		}
	}
	// image_url: distinguish between absent, null (delete) and a string value
	if v, present := input["image_url"]; present {
		if v == nil {
			// explicit null -> remove image
			prod.ImageURL = nil
		} else if s, ok := v.(string); ok {
			prod.ImageURL = &s
		}
	}
	if v, ok := input["color"].(string); ok {
		prod.Color = v
	}
	if v, ok := input["condition"].(string); ok {
		prod.Condition = v
	}
	if v, ok := input["country"].(string); ok {
		prod.Country = v
	}
	if v, ok := input["material"].(string); ok {
		prod.Material = v
	}
	DB.Save(&prod)
	c.JSON(http.StatusOK, prod)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	res := DB.Delete(&Product{}, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ==================== CART & ORDER HANDLERS ====================
func getOrCreateCartID(c *gin.Context) string {
	cartID, err := c.Cookie(CartCookieName)
	if err == nil && cartID != "" {
		return cartID
	}
	newID := "cart-" + randomHex(12)
	c.SetCookie(CartCookieName, newID, 3600*24*30, "/", "", false, false)
	return newID
}

type addCartReq struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

func AddToCart(c *gin.Context) {
	var req addCartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	var prod Product
	if err := DB.First(&prod, req.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		return
	}
	// check stock and existing quantity
	cartID := getOrCreateCartID(c)
	var existing CartItem
	err := DB.Where("cart_id = ? AND product_id = ?", cartID, req.ProductID).First(&existing).Error
	existingQty := 0
	if err == nil {
		existingQty = existing.Quantity
	}
	if existingQty+req.Quantity > prod.Stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock"})
		return
	}
	var item CartItem
	err = DB.Where("cart_id = ? AND product_id = ?", cartID, req.ProductID).First(&item).Error
	if err == nil {
		item.Quantity += req.Quantity
		DB.Save(&item)
	} else {
		item = CartItem{CartID: cartID, ProductID: req.ProductID, Quantity: req.Quantity}
		DB.Create(&item)
	}
	c.JSON(http.StatusOK, gin.H{"message": "added"})
}

// UpdateCartQuantity sets a new quantity for a product in the cart (or removes it if quantity <= 0)
func UpdateCartQuantity(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cartID := getOrCreateCartID(c)
	var prod Product
	if err := DB.First(&prod, req.ProductID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		return
	}
	if req.Quantity < 0 {
		req.Quantity = 0
	}
	if req.Quantity > prod.Stock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock"})
		return
	}
	var item CartItem
	err := DB.Where("cart_id = ? AND product_id = ?", cartID, req.ProductID).First(&item).Error
	if err != nil {
		// if not found and quantity > 0 create
		if req.Quantity > 0 {
			item = CartItem{CartID: cartID, ProductID: req.ProductID, Quantity: req.Quantity}
			if err := DB.Create(&item).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot add item"})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "updated"})
		return
	}
	if req.Quantity == 0 {
		DB.Delete(&item)
		c.JSON(http.StatusOK, gin.H{"message": "removed"})
		return
	}
	item.Quantity = req.Quantity
	DB.Save(&item)
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func RemoveFromCart(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cartID := getOrCreateCartID(c)
	res := DB.Where("cart_id = ? AND product_id = ?", cartID, req.ProductID).Delete(&CartItem{})
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "removed"})
}

func GetCart(c *gin.Context) {
	cartID := getOrCreateCartID(c)
	var items []CartItem
	DB.Where("cart_id = ?", cartID).Find(&items)
	type ItemOut struct {
		Product  Product `json:"product"`
		Quantity int     `json:"quantity"`
	}
	out := []ItemOut{}
	var total float64
	for _, it := range items {
		var p Product
		if err := DB.First(&p, it.ProductID).Error; err != nil {
			continue
		}
		out = append(out, ItemOut{Product: p, Quantity: it.Quantity})
		total += float64(it.Quantity) * p.Price
	}
	c.JSON(http.StatusOK, gin.H{"items": out, "total": total})
}

func clearCart(tx *gorm.DB, cartID string) error {
	return tx.Where("cart_id = ?", cartID).Delete(&CartItem{}).Error
}

func CreateOrderHandler(c *gin.Context) {
	var req struct {
		Name          string  `json:"name"`
		Email         string  `json:"email"`
		Address       string  `json:"address"`
		Phone         string  `json:"phone"`
		DeliveryType  string  `json:"delivery_type"`   // "pickup" or "courier"
		PickupPointID uint    `json:"pickup_point_id"` // preferred for pickup
		PickupPoint   string  `json:"pickup_point"`    // только для pickup
		DeliveryLat   float64 `json:"delivery_lat"`    // только для courier
		DeliveryLng   float64 `json:"delivery_lng"`    // только для courier
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and email required"})
		return
	}
	if strings.TrimSpace(req.DeliveryType) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delivery_type required (pickup or courier)"})
		return
	}
	if req.DeliveryType == "courier" && (req.DeliveryLat == 0 || req.DeliveryLng == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "delivery coordinates required for courier delivery"})
		return
	}
	if req.DeliveryType == "pickup" {
		if req.PickupPointID > 0 {
			var pickupPoint PickupPoint
			if err := DB.First(&pickupPoint, req.PickupPointID).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "pickup point not found"})
				return
			}
			req.PickupPoint = strings.TrimSpace(pickupPoint.Name + " - " + pickupPoint.Address)
		}
		if strings.TrimSpace(req.PickupPoint) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "pickup_point_id or pickup_point required for pickup delivery"})
			return
		}
	}

	cartID := getOrCreateCartID(c)
	var items []CartItem
	DB.Where("cart_id = ?", cartID).Find(&items)
	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart empty"})
		return
	}

	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var total float64
	var deliveryPrice float64 = 0

	// Calculate delivery price
	if req.DeliveryType == "courier" {
		// базовая цена доставки 200 рублей + 50 рублей за каждый км от центра (красная площадь)
		// координаты красной площади: 55.7558, 37.6223
		centerLat := 55.7558
		centerLng := 37.6223

		distance := calculateDistance(centerLat, centerLng, req.DeliveryLat, req.DeliveryLng)
		deliveryPrice = 200 + (distance * 50)
	}
	// для самовывоза доставка бесплатна

	order := Order{
		CustomerName:  req.Name,
		Email:         req.Email,
		Phone:         req.Phone,
		Address:       req.Address,
		Status:        "pending",
		DeliveryType:  req.DeliveryType,
		PickupPoint:   req.PickupPoint,
		DeliveryLat:   req.DeliveryLat,
		DeliveryLng:   req.DeliveryLng,
		DeliveryPrice: deliveryPrice,
		CreatedAt:     time.Now(),
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create order"})
		return
	}

	var orderItems []OrderItem
	for _, it := range items {
		var p Product
		if err := tx.First(&p, it.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			return
		}
		if p.Stock < it.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("not enough stock for %s", p.Name)})
			return
		}
		oi := OrderItem{
			OrderID:     order.ID,
			ProductID:   p.ID,
			ProductName: p.Name,
			Price:       p.Price,
			Quantity:    it.Quantity,
		}
		if err := tx.Create(&oi).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create order item"})
			return
		}
		orderItems = append(orderItems, oi)
		p.Stock -= it.Quantity
		tx.Save(&p)
		total += float64(it.Quantity) * p.Price
	}

	total += deliveryPrice
	order.Total = total
	tx.Save(&order)
	if err := clearCart(tx, cartID); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot clear cart"})
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}

	// Save to file
	if err := saveOrderToFile(order, orderItems); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to save order file: %v\n", err)
	}

	// Save to Excel
	if err := saveOrderToExcel(order, orderItems); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to save order to Excel: %v\n", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"order_id":       order.ID,
		"status":         order.Status,
		"delivery_type":  order.DeliveryType,
		"delivery_price": order.DeliveryPrice,
		"total":          order.Total,
	})
}

// ==================== ADMIN ORDERS ====================
func AdminListOrders(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "20")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	var total int64
	DB.Model(&Order{}).Count(&total)

	var orders []Order
	DB.Order("created_at desc").Offset(offset).Limit(perPage).Find(&orders)
	c.JSON(http.StatusOK, gin.H{
		"items":      orders,
		"page":       page,
		"per_page":   perPage,
		"total":      total,
		"total_page": (total + int64(perPage) - 1) / int64(perPage),
	})
}

func AdminGetOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := DB.Preload("Items").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func AdminUpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var order Order
	if err := DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	order.Status = input.Status
	DB.Save(&order)
	c.JSON(http.StatusOK, order)
}

// ==================== IMAGE UPLOAD ====================
func UploadImageHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		ext = ".bin"
	}
	fn := fmt.Sprintf("%d-%s%s", time.Now().Unix(), randomHex(6), ext)
	fn = sanitizeFilename(fn)
	dst := filepath.Join(UploadDir, fn)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}
	url := "/static/uploads/" + fn
	c.JSON(http.StatusOK, gin.H{"url": url})
}

// ==================== AUTH ====================
func AdminLogin(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Username != ADMIN_USERNAME {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword(adminHash, []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	token, err := createToken(input.Username, "admin", time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token error"})
		return
	}
	c.SetCookie(AdminCookieName, token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func AdminLogout(c *gin.Context) {
	c.SetCookie(AdminCookieName, "", -1, "/", "", false, false)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func UserRegister(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := strings.ToLower(strings.TrimSpace(input.Username))
	password := strings.TrimSpace(input.Password)

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "логин и пароль обязательны"})
		return
	}
	if len(username) < 3 || len(username) > 32 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "логин должен быть от 3 до 32 символов"})
		return
	}
	if strings.ContainsAny(username, " \t\r\n") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "логин не должен содержать пробелы"})
		return
	}
	if len(password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пароль должен быть не короче 6 символов"})
		return
	}

	var existing User
	if err := DB.Where("username = ?", username).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "пользователь с таким логином уже существует"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка базы данных"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка шифрования пароля"})
		return
	}

	user := User{
		Username:     username,
		Role:         "user",
		PasswordHash: string(hash),
	}
	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "registered", "user": userPublic(&user)})
}

func UserLogin(c *gin.Context) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := strings.ToLower(strings.TrimSpace(input.Username))
	password := strings.TrimSpace(input.Password)
	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "логин и пароль обязательны"})
		return
	}

	var user User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный логин или пароль"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный логин или пароль"})
		return
	}

	role := normalizeUserRole(user.Role)
	userToken, err := createToken(strconv.FormatUint(uint64(user.ID), 10), role, time.Hour*24*30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка токена"})
		return
	}
	adminToken := ""
	if user.Username == ADMIN_USERNAME {
		adminToken, err = createToken(user.Username, "admin", time.Hour*24)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "РѕС€РёР±РєР° С‚РѕРєРµРЅР°"})
			return
		}
	}
	c.SetCookie(UserCookieName, userToken, 3600*24*30, "/", "", false, true)
	if adminToken != "" {
		c.SetCookie(AdminCookieName, adminToken, 3600*24, "/", "", false, true)
	} else {
		c.SetCookie(AdminCookieName, "", -1, "/", "", false, true)
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "user": userPublic(&user)})
}

func UserLogout(c *gin.Context) {
	c.SetCookie(UserCookieName, "", -1, "/", "", false, true)
	c.SetCookie(AdminCookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func UserMe(c *gin.Context) {
	user, err := currentUserFromCookie(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"authenticated": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"authenticated": true, "user": userPublic(user)})
}

// ==================== PICKUP POINTS ====================
func GetPickupPoints(c *gin.Context) {
	var points []PickupPoint
	query := DB.Order("id asc")
	if c.Query("all") != "1" {
		query = query.Where("city IN ?", []string{"Moscow", "Moscow Oblast", "Москва", "Московская область"})
	}
	query.Find(&points)
	c.JSON(http.StatusOK, points)
}
