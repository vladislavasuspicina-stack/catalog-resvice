package main

import (
	"errors"
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
	UserCookieName  = "user_token"
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
	Brand       string    `json:"brand"`
	Color       string    `json:"color"`
	Condition   string    `json:"condition"`
	Country     string    `json:"country"`
	Material    string    `json:"material"`
	CreatedAt   time.Time `json:"created_at"`
}

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64" json:"username"`
	Role         string    `gorm:"size:16;default:user" json:"role"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type CartItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CartID    string    `json:"cart_id" gorm:"index"`
	ProductID uint      `json:"product_id" gorm:"index"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type Order struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	CustomerName  string      `json:"customer_name"`
	Email         string      `json:"email"`
	Phone         string      `json:"phone"`
	Address       string      `json:"address"`
	Total         float64     `json:"total"`
	Status        string      `json:"status"`
	DeliveryType  string      `json:"delivery_type"`  // "pickup" or "courier"
	PickupPoint   string      `json:"pickup_point"`   // human-readable pickup location
	DeliveryLat   float64     `json:"delivery_lat"`   // latitude for courier
	DeliveryLng   float64     `json:"delivery_lng"`   // longitude for courier
	DeliveryPrice float64     `json:"delivery_price"` // delivery price
	CreatedAt     time.Time   `json:"created_at"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type PickupPoint struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `json:"name"`
	City         string    `json:"city"`
	Address      string    `json:"address"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	WorkingHours string    `json:"working_hours"`
	Phone        string    `json:"phone"`
	Details      string    `json:"details"`
	CreatedAt    time.Time `json:"created_at"`
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
		panic("DB connection error: " + err.Error())
	}
	if err := db.AutoMigrate(&Category{}, &Product{}, &User{}, &CartItem{}, &Order{}, &OrderItem{}, &PickupPoint{}); err != nil {
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

	// Keep admin credentials available in the regular user auth flow as well.
	var adminUser User
	err = DB.Where("username = ?", ADMIN_USERNAME).First(&adminUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := DB.Create(&User{
			Username:     ADMIN_USERNAME,
			Role:         "admin",
			PasswordHash: string(adminHash),
		}).Error; err != nil {
			panic("cannot create admin user: " + err.Error())
		}
	} else if err == nil {
		if err := DB.Model(&adminUser).Updates(map[string]interface{}{
			"password_hash": string(adminHash),
			"role":          "admin",
		}).Error; err != nil {
			panic("cannot update admin user data: " + err.Error())
		}
	} else {
		panic("cannot load admin user: " + err.Error())
	}

	// Backfill role for old users created before role support.
	DB.Model(&User{}).Where("role IS NULL OR role = ''").Update("role", "user")

	defaults := defaultPickupPoints()
	for _, point := range defaults {
		var existing PickupPoint
		err := DB.Where("name = ?", point.Name).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			DB.Create(&point)
			continue
		}
		if err != nil {
			continue
		}
		DB.Model(&existing).Updates(map[string]interface{}{
			"city":          point.City,
			"address":       point.Address,
			"latitude":      point.Latitude,
			"longitude":     point.Longitude,
			"working_hours": point.WorkingHours,
			"phone":         point.Phone,
			"details":       point.Details,
		})
	}

	// Remove obsolete seeded points from previous versions.
	DB.Where("name IN ?", []string{
		"Moscow Center Pickup",
		"Moscow Central Pickup",
		"Moscow North Pickup",
		"Moscow South Pickup",
		"Moscow West Pickup",
		"Moscow East Pickup",
		"Moscow North-West Pickup",
		"Moscow North-East Pickup",
		"Moscow South-East Pickup",
		"Zelenograd Pickup",
		"Khimki Pickup",
		"Mytishchi Pickup",
		"Korolyov Pickup",
		"Balashikha Pickup",
		"Lyubertsy Pickup",
		"Podolsk Pickup",
		"Odintsovo Pickup",
		"Krasnogorsk Pickup",
		"Domodedovo Pickup",
		"Saint Petersburg Nevsky Pickup",
		"Kazan Center Pickup",
		"Yekaterinburg Downtown Pickup",
		"Novosibirsk Pickup Hub",
		"Krasnoyarsk Central Pickup",
		"Vladivostok Pickup Point",
		"Sochi City Pickup",
		"Nizhny Novgorod Pickup",
		"Rostov-on-Don Pickup",
	}).Delete(&PickupPoint{})

	DB.Model(&PickupPoint{}).
		Where("working_hours IS NULL OR working_hours = '' OR phone IS NULL OR phone = '' OR details IS NULL OR details = ''").
		Updates(map[string]interface{}{
			"working_hours": "10:00-22:00, daily",
			"phone":         "+7 (800) 555-35-35",
			"details":       "Pickup and return counter. Bring passport or order code.",
		})
	DB.Model(&PickupPoint{}).
		Where("city IS NULL OR city = ''").
		Update("city", "Moscow")
}

func defaultPickupPoints() []PickupPoint {
	return []PickupPoint{
		{
			Name:         "Пункт выдачи Тверская",
			City:         "Москва",
			Address:      "Москва, Тверская улица, 7",
			Latitude:     55.7571,
			Longitude:    37.6156,
			WorkingHours: "ежедневно 09:00-22:00",
			Phone:        "+7 (495) 100-10-10",
			Details:      "Центральный пункт выдачи, быстрая стойка получения.",
		},
		{
			Name:         "Пункт выдачи Алтуфьево",
			City:         "Москва",
			Address:      "Москва, Алтуфьевское шоссе, 86",
			Latitude:     55.8798,
			Longitude:    37.5869,
			WorkingHours: "ежедневно 10:00-22:00",
			Phone:        "+7 (495) 100-10-11",
			Details:      "Рядом со станцией метро Алтуфьево.",
		},
		{
			Name:         "Пункт выдачи Каширский",
			City:         "Москва",
			Address:      "Москва, Каширское шоссе, 61",
			Latitude:     55.6101,
			Longitude:    37.7174,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (495) 100-10-12",
			Details:      "Удобен для южных районов, рядом есть парковка.",
		},
		{
			Name:         "Пункт выдачи Рублевский",
			City:         "Москва",
			Address:      "Москва, Рублевское шоссе, 52А",
			Latitude:     55.7584,
			Longitude:    37.4082,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (495) 100-10-13",
			Details:      "Пункт выдачи и возврата на западе города.",
		},
		{
			Name:         "Пункт выдачи Новокосино",
			City:         "Москва",
			Address:      "Москва, Новокосинская улица, 32",
			Latitude:     55.7448,
			Longitude:    37.8632,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (495) 100-10-14",
			Details:      "Пункт выдачи для восточных районов Москвы.",
		},
		{
			Name:         "Пункт выдачи Митино",
			City:         "Москва",
			Address:      "Москва, Митинская улица, 40",
			Latitude:     55.8450,
			Longitude:    37.3622,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (495) 100-10-15",
			Details:      "Рядом со станцией метро Митино.",
		},
		{
			Name:         "Пункт выдачи ВДНХ",
			City:         "Москва",
			Address:      "Москва, проспект Мира, 211",
			Latitude:     55.8211,
			Longitude:    37.6380,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (495) 100-10-16",
			Details:      "Пункт выдачи рядом с ВДНХ.",
		},
		{
			Name:         "Пункт выдачи Люблино",
			City:         "Москва",
			Address:      "Москва, Люблинская улица, 153",
			Latitude:     55.6776,
			Longitude:    37.7612,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (495) 100-10-17",
			Details:      "Пункт выдачи для юго-восточных районов.",
		},
		{
			Name:         "Пункт выдачи Зеленоград",
			City:         "Москва",
			Address:      "Зеленоград, Панфиловский проспект, 6А",
			Latitude:     55.9825,
			Longitude:    37.1815,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (495) 100-10-18",
			Details:      "Удобный подъезд на автомобиле.",
		},
		{
			Name:         "Пункт выдачи Химки",
			City:         "Московская область",
			Address:      "Химки, Московская улица, 14",
			Latitude:     55.8892,
			Longitude:    37.4457,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-10",
			Details:      "Северный пункт выдачи в Московской области.",
		},
		{
			Name:         "Пункт выдачи Мытищи",
			City:         "Московская область",
			Address:      "Мытищи, улица Мира, 32/2",
			Latitude:     55.9104,
			Longitude:    37.7360,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-11",
			Details:      "Пункт выдачи на северо-востоке области.",
		},
		{
			Name:         "Пункт выдачи Королев",
			City:         "Московская область",
			Address:      "Королев, проспект Космонавтов, 34Б",
			Latitude:     55.9162,
			Longitude:    37.8541,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-12",
			Details:      "Рядом с железнодорожной станцией.",
		},
		{
			Name:         "Пункт выдачи Балашиха",
			City:         "Московская область",
			Address:      "Балашиха, Советская улица, 9",
			Latitude:     55.7963,
			Longitude:    37.9382,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-13",
			Details:      "Восточный пункт выдачи в Московской области.",
		},
		{
			Name:         "Пункт выдачи Люберцы",
			City:         "Московская область",
			Address:      "Люберцы, Октябрьский проспект, 112",
			Latitude:     55.6765,
			Longitude:    37.8981,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-14",
			Details:      "Пункт выдачи на юго-востоке области.",
		},
		{
			Name:         "Пункт выдачи Подольск",
			City:         "Московская область",
			Address:      "Подольск, улица Кирова, 39",
			Latitude:     55.4311,
			Longitude:    37.5455,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-15",
			Details:      "Южный пункт выдачи в Московской области.",
		},
		{
			Name:         "Пункт выдачи Одинцово",
			City:         "Московская область",
			Address:      "Одинцово, Можайское шоссе, 58А",
			Latitude:     55.6780,
			Longitude:    37.2777,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-16",
			Details:      "Западный пункт выдачи в Московской области.",
		},
		{
			Name:         "Пункт выдачи Красногорск",
			City:         "Московская область",
			Address:      "Красногорск, улица Ленина, 38Б",
			Latitude:     55.8317,
			Longitude:    37.3299,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-17",
			Details:      "Северо-западный пункт выдачи в области.",
		},
		{
			Name:         "Пункт выдачи Домодедово",
			City:         "Московская область",
			Address:      "Домодедово, Каширское шоссе, 3А",
			Latitude:     55.4430,
			Longitude:    37.7470,
			WorkingHours: "ежедневно 10:00-21:00",
			Phone:        "+7 (498) 200-20-18",
			Details:      "Удобный пункт выдачи рядом с трассой к аэропорту.",
		},
	}
}
