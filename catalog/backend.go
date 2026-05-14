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
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

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
	var cat Category
	if err := DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	ids := []uint{cat.ID}
	for i := 0; i < len(ids); i++ {
		var children []Category
		DB.Where("parent_id = ?", ids[i]).Find(&children)
		for _, child := range children {
			ids = append(ids, child.ID)
		}
	}

	if err := DB.Model(&Product{}).Where("category_id IN ?", ids).Update("category_id", nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot detach products"})
		return
	}
	if err := DB.Exec("DELETE FROM product_categories WHERE category_id IN ?", ids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot detach product categories"})
		return
	}
	if err := DB.Where("id IN ?", ids).Delete(&Category{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ==================== PRODUCT HANDLERS ====================
func normalizeCategoryIDs(ids []uint, fallback *uint) []uint {
	seen := map[uint]bool{}
	out := []uint{}
	add := func(id uint) {
		if id == 0 || seen[id] {
			return
		}
		var count int64
		DB.Model(&Category{}).Where("id = ?", id).Count(&count)
		if count == 0 {
			return
		}
		seen[id] = true
		out = append(out, id)
	}
	if fallback != nil {
		add(*fallback)
	}
	for _, id := range ids {
		add(id)
	}
	return out
}

func categoryWithDescendantIDs(root uint) []uint {
	if root == 0 {
		return []uint{}
	}
	ids := []uint{root}
	for i := 0; i < len(ids); i++ {
		var children []Category
		DB.Where("parent_id = ?", ids[i]).Find(&children)
		for _, child := range children {
			ids = append(ids, child.ID)
		}
	}
	return ids
}

func setProductCategories(p *Product, ids []uint) error {
	ids = normalizeCategoryIDs(ids, p.CategoryID)
	if len(ids) > 0 {
		first := ids[0]
		p.CategoryID = &first
	} else {
		p.CategoryID = nil
	}
	categories := []Category{}
	if len(ids) > 0 {
		if err := DB.Where("id IN ?", ids).Find(&categories).Error; err != nil {
			return err
		}
	}
	if err := DB.Save(p).Error; err != nil {
		return err
	}
	return DB.Model(p).Association("Categories").Replace(categories)
}

func enrichProductCategories(p *Product) {
	if len(p.Categories) == 0 {
		DB.Model(p).Association("Categories").Find(&p.Categories)
	}
	ids := []uint{}
	seen := map[uint]bool{}
	for _, category := range p.Categories {
		if category.ID == 0 || seen[category.ID] {
			continue
		}
		seen[category.ID] = true
		ids = append(ids, category.ID)
	}
	if p.CategoryID != nil && !seen[*p.CategoryID] {
		ids = append([]uint{*p.CategoryID}, ids...)
	}
	p.CategoryIDs = ids
}

func productSearchTerms(search string) []string {
	search = strings.TrimSpace(search)
	seen := map[string]bool{}
	terms := []string{}
	add := func(term string) {
		term = strings.TrimSpace(term)
		if term == "" || seen[term] {
			return
		}
		seen[term] = true
		terms = append(terms, term)
		runes := []rune(term)
		if len(runes) == 0 {
			return
		}
		title := string(unicode.ToUpper(runes[0])) + string(runes[1:])
		lower := string(unicode.ToLower(runes[0])) + string(runes[1:])
		if title != term && !seen[title] {
			seen[title] = true
			terms = append(terms, title)
		}
		if lower != term && !seen[lower] {
			seen[lower] = true
			terms = append(terms, lower)
		}
	}
	add(search)
	low := strings.ToLower(search)
	if strings.Contains(low, "смартфон") || strings.Contains(low, "телефон") {
		add("телефон")
		add("смартфон")
		add("ксяоми")
		add("xiaomi")
		add("сасунг")
		add("samsung")
		add("айфон")
		add("iphone")
	}
	if strings.Contains(low, "samsung") || strings.Contains(low, "самсунг") || strings.Contains(low, "сасунг") {
		add("Samsung")
		add("samsung")
		add("самсунг")
		add("сасунг")
	}
	if strings.Contains(low, "xiaomi") || strings.Contains(low, "ксяоми") || strings.Contains(low, "сяоми") {
		add("Xiaomi")
		add("xiaomi")
		add("ксяоми")
		add("сяоми")
	}
	if strings.Contains(low, "apple") || strings.Contains(low, "iphone") || strings.Contains(low, "айфон") {
		add("Apple")
		add("apple")
		add("iPhone")
		add("айфон")
	}
	if strings.Contains(low, "huawei") || strings.Contains(low, "хуавей") {
		add("Huawei")
		add("huawei")
		add("хуавей")
	}
	if strings.Contains(low, "honor") || strings.Contains(low, "хонор") {
		add("Honor")
		add("honor")
		add("хонор")
	}
	if strings.Contains(low, "часы") {
		add("часы")
		add("смарт часы")
		add("SmartWatch")
		add("watch")
	}
	if strings.Contains(low, "ноутбук") {
		add("ноутбук")
		add("laptop")
	}
	return terms
}

func applyProductSearch(query *gorm.DB, search string, fields []string) *gorm.DB {
	variants := productSearchTerms(search)
	clauses := make([]string, 0, len(variants))
	args := make([]interface{}, 0, len(variants)*len(fields))
	for _, variant := range variants {
		like := "%" + variant + "%"
		fieldClauses := make([]string, 0, len(fields))
		for _, field := range fields {
			fieldClauses = append(fieldClauses, field+" LIKE ?")
			args = append(args, like)
		}
		clauses = append(clauses, "("+strings.Join(fieldClauses, " OR ")+")")
	}
	if len(clauses) == 0 {
		return query
	}
	return query.Where(strings.Join(clauses, " OR "), args...)
}

func inferProductBrand(p Product) string {
	if strings.TrimSpace(p.Brand) != "" {
		return strings.TrimSpace(p.Brand)
	}
	text := strings.ToLower(strings.Join([]string{p.Name, p.Description, p.Material, p.Country}, " "))
	switch {
	case strings.Contains(text, "xiaomi") || strings.Contains(text, "ксяоми") || strings.Contains(text, "сяоми"):
		return "Xiaomi"
	case strings.Contains(text, "samsung") || strings.Contains(text, "самсунг") || strings.Contains(text, "сасунг"):
		return "Samsung"
	case strings.Contains(text, "iphone") || strings.Contains(text, "apple") || strings.Contains(text, "айфон"):
		return "Apple"
	case strings.Contains(text, "huawei") || strings.Contains(text, "хуавей"):
		return "Huawei"
	case strings.Contains(text, "honor") || strings.Contains(text, "хонор"):
		return "Honor"
	case strings.Contains(text, "smartwatch") || strings.Contains(text, "smart watch"):
		return "SmartWatch"
	case strings.Contains(text, "asus"):
		return "ASUS"
	case strings.Contains(text, "lenovo"):
		return "Lenovo"
	case strings.Contains(text, "acer"):
		return "Acer"
	case strings.Contains(text, "dell"):
		return "Dell"
	case strings.Contains(text, "hp"):
		return "HP"
	default:
		return ""
	}
}

func enrichProductBrand(p *Product) {
	if strings.TrimSpace(p.Brand) == "" {
		p.Brand = inferProductBrand(*p)
	}
}

func enrichProduct(p *Product) {
	enrichProductBrand(p)
	enrichProductCategories(p)
}

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
		Brand:       strings.TrimSpace(input.Brand),
		Color:       input.Color,
		Condition:   input.Condition,
		Country:     input.Country,
		Material:    input.Material,
	}
	if err := DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create product"})
		return
	}
	if err := setProductCategories(&p, input.CategoryIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save product categories"})
		return
	}
	enrichProduct(&p)
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
	brand := c.Query("brand")

	query := DB.Model(&Product{})

	if search != "" {
		query = applyProductSearch(query, search, []string{"name", "description", "brand", "color", "condition", "country", "material"})
	}
	if brand != "" {
		query = applyProductSearch(query, brand, []string{"brand", "name", "description"})
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
			ids := categoryWithDescendantIDs(uint(v))
			query = query.Where("category_id IN ? OR id IN (SELECT product_id FROM product_categories WHERE category_id IN ?)", ids, ids)
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
	query.Preload("Categories").Offset(offset).Limit(perPage).Find(&products)
	for i := range products {
		enrichProduct(&products[i])
	}

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
	if err := DB.Preload("Categories").First(&p, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	enrichProduct(&p)
	c.JSON(http.StatusOK, p)
}

func GetBrands(c *gin.Context) {
	var products []Product
	DB.Select([]string{"brand", "name", "description", "material", "country"}).Find(&products)
	seen := map[string]bool{}
	brands := []string{}
	for _, product := range products {
		brand := inferProductBrand(product)
		if brand == "" || seen[brand] {
			continue
		}
		seen[brand] = true
		brands = append(brands, brand)
	}
	sort.Strings(brands)
	c.JSON(http.StatusOK, gin.H{"items": brands})
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
	if v, ok := input["brand"].(string); ok {
		prod.Brand = strings.TrimSpace(v)
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
	if err := DB.Save(&prod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save product"})
		return
	}
	if v, ok := input["category_ids"]; ok {
		ids := []uint{}
		if raw, ok := v.([]interface{}); ok {
			for _, item := range raw {
				switch t := item.(type) {
				case float64:
					ids = append(ids, uint(t))
				case string:
					if parsed, err := strconv.Atoi(t); err == nil {
						ids = append(ids, uint(parsed))
					}
				}
			}
		}
		if err := setProductCategories(&prod, ids); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot save product categories"})
			return
		}
	}
	enrichProduct(&prod)
	c.JSON(http.StatusOK, prod)
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	res := DB.Delete(&Product{}, id)
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	DB.Where("product_id = ?", id).Delete(&CartItem{})
	DB.Exec("DELETE FROM product_categories WHERE product_id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ==================== CART & ORDER HANDLERS ====================
func getOrCreateCartID(c *gin.Context) string {
	if user, err := currentUserFromCookie(c); err == nil {
		cartID := fmt.Sprintf("user-%d", user.ID)
		c.SetCookie(CartCookieName, cartID, 3600*24*30, "/", "", false, false)
		return cartID
	}
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
			DB.Delete(&it)
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

	var userID *uint
	if user, err := currentUserFromCookie(c); err == nil {
		userID = &user.ID
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
		UserID:        userID,
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

func GetOrderStatus(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := DB.Preload("Items").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.UserID != nil {
		user, err := currentUserFromCookie(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "auth required"})
			return
		}
		if normalizeUserRole(user.Role) != "admin" && *order.UserID != user.ID {
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}
	}
	c.JSON(http.StatusOK, order)
}

func UserListOrders(c *gin.Context) {
	user, err := currentUserFromCookie(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"items": []Order{}})
		return
	}
	var orders []Order
	DB.Preload("Items").
		Where("user_id = ?", user.ID).
		Order("created_at desc, id desc").
		Find(&orders)
	c.JSON(http.StatusOK, gin.H{"items": orders})
}

// ==================== ADMIN ORDERS ====================
func validOrderStatus(status string) bool {
	switch status {
	case "pending", "processing", "shipped", "delivered", "cancelled":
		return true
	default:
		return false
	}
}

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
	DB.Preload("Items").Order("created_at desc").Offset(offset).Limit(perPage).Find(&orders)
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
	status := strings.TrimSpace(input.Status)
	if !validOrderStatus(status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}
	var order Order
	if err := DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	order.Status = status
	if err := DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update order"})
		return
	}
	DB.Preload("Items").First(&order, id)
	c.JSON(http.StatusOK, order)
}

func AdminDeleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order Order
	if err := DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	tx := DB.Begin()
	if err := tx.Where("order_id = ?", order.ID).Delete(&OrderItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete order items"})
		return
	}
	if err := tx.Delete(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete order"})
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}

	_ = os.Remove(fmt.Sprintf("./orders/order-%d.json", order.ID))
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
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

	userToken, err := createToken(strconv.FormatUint(uint64(user.ID), 10), "user", time.Hour*24*30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "РѕС€РёР±РєР° С‚РѕРєРµРЅР°"})
		return
	}
	c.SetCookie(UserCookieName, userToken, 3600*24*30, "/", "", false, true)
	c.SetCookie(AdminCookieName, "", -1, "/", "", false, true)
	c.SetCookie(CartCookieName, fmt.Sprintf("user-%d", user.ID), 3600*24*30, "/", "", false, false)

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
	c.SetCookie(CartCookieName, fmt.Sprintf("user-%d", user.ID), 3600*24*30, "/", "", false, false)
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
	c.SetCookie(CartCookieName, "", -1, "/", "", false, false)
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

// ==================== ADMIN USERS ====================
func AdminListUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "50")
	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 50
	}
	offset := (page - 1) * perPage

	query := DB.Model(&User{})
	search := strings.ToLower(strings.TrimSpace(c.Query("search")))
	if search != "" {
		query = query.Where("LOWER(username) LIKE ?", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var users []User
	query.Order("created_at desc, id desc").Offset(offset).Limit(perPage).Find(&users)
	out := make([]gin.H, 0, len(users))
	for i := range users {
		item := userPublic(&users[i])
		item["password_hash"] = users[i].PasswordHash
		item["created_at"] = users[i].CreatedAt
		out = append(out, item)
	}
	c.JSON(http.StatusOK, gin.H{
		"items":      out,
		"page":       page,
		"per_page":   perPage,
		"total":      total,
		"total_page": (total + int64(perPage) - 1) / int64(perPage),
	})
}

func AdminUpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user User
	if err := DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var input struct {
		Username string `json:"username"`
		Role     string `json:"role"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := strings.ToLower(strings.TrimSpace(input.Username))
	if username != "" && username != user.Username {
		if len(username) < 3 || len(username) > 32 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username must be 3-32 characters"})
			return
		}
		if strings.ContainsAny(username, " \t\r\n") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username must not contain spaces"})
			return
		}
		if user.Username == ADMIN_USERNAME {
			c.JSON(http.StatusBadRequest, gin.H{"error": "built-in admin username cannot be changed"})
			return
		}
		var existing User
		err := DB.Where("username = ? AND id <> ?", username, user.ID).First(&existing).Error
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
			return
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		user.Username = username
	}

	role := strings.ToLower(strings.TrimSpace(input.Role))
	if role != "" {
		if role != "user" && role != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
		if user.Username == ADMIN_USERNAME && role != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "built-in admin must keep admin role"})
			return
		}
		user.Role = role
	}

	password := strings.TrimSpace(input.Password)
	if password != "" {
		if len(password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 6 characters"})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "password hash error"})
			return
		}
		user.PasswordHash = string(hash)
	}

	if err := DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": userPublic(&user)})
}

func AdminDeleteUser(c *gin.Context) {
	id := c.Param("id")
	current, err := currentUserFromCookie(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth required"})
		return
	}

	var user User
	if err := DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if user.Username == ADMIN_USERNAME {
		c.JSON(http.StatusBadRequest, gin.H{"error": "built-in admin cannot be deleted"})
		return
	}
	if user.ID == current.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "current user cannot be deleted"})
		return
	}

	tx := DB.Begin()
	if err := tx.Model(&Order{}).Where("user_id = ?", user.ID).Update("user_id", nil).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot detach user orders"})
		return
	}
	if err := tx.Where("cart_id = ?", fmt.Sprintf("user-%d", user.ID)).Delete(&CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete user cart"})
		return
	}
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot delete user"})
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// ==================== PICKUP POINTS ====================
func GetPickupPoints(c *gin.Context) {
	var points []PickupPoint
	query := DB.Order("id asc")
	if c.Query("all") != "1" {
		query = query.Where("city IN ?", []string{"Moscow", "Moscow Oblast", "Москва", "Московская область"})
	}
	query.Find(&points)
	if len(points) == 0 {
		for _, point := range defaultPickupPoints() {
			DB.FirstOrCreate(&point, PickupPoint{Name: point.Name})
		}
		query.Find(&points)
	}
	c.JSON(http.StatusOK, gin.H{"points": points})
}
