package main

import (
	"fmt"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ==================== CONFIG ====================
const (
	DBPath          = "./catalog.db"
	UploadDir       = "./static/uploads"
	AdminCookieName = "admin_token"
	CartCookieName  = "cart_id"
	ADMIN_USERNAME  = "admin"
	ADMIN_PASSWORD  = "12345"
	JWT_SECRET      = "replace_this_with_strong_secret"
)

// ==================== MODELS ====================
type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
	Products  []Product `json:"-" gorm:"foreignKey:CategoryID"`
}

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  *uint     `json:"category_id"`
	ImageURL    *string   `json:"image_url"`
	Color       string    `json:"color"`
	Condition   string    `json:"condition"`
	Country     string    `json:"country"`
	Material    string    `json:"material"`
	CreatedAt   time.Time `json:"created_at"`
}

type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    string    `json:"cart_id" gorm:"index"`
	ProductID uint      `json:"product_id" gorm:"index"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	CustomerName string      `json:"customer_name"`
	Email        string      `json:"email"`
	Address      string      `json:"address"`
	Total        float64     `json:"total"`
	Status       string      `json:"status"`
	CreatedAt    time.Time   `json:"created_at"`
	Items        []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	OrderID     uint    `json:"order_id" gorm:"index"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

// ==================== DB ====================
var DB *gorm.DB

func InitDB() {
	if err := os.MkdirAll(UploadDir, 0o755); err != nil {
		panic("cannot create upload dir: " + err.Error())
	}
	db, err := gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		panic("Ошибка подключения к БД: " + err.Error())
	}
	if err := db.AutoMigrate(&Category{}, &Product{}, &CartItem{}, &Order{}, &OrderItem{}); err != nil {
		panic("AutoMigrate error: " + err.Error())
	}
	DB = db
	fmt.Println("DB connected:", DBPath)
	ensureAdminData()
}

// ==================== ADMIN SETUP ====================
var adminHash []byte

func ensureAdminData() {
	h, err := bcrypt.GenerateFromPassword([]byte(ADMIN_PASSWORD), bcrypt.DefaultCost)
	if err != nil {
		panic("bcrypt error: " + err.Error())
	}
	adminHash = h
}
