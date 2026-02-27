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
			Name:         "Moscow Central Pickup",
			City:         "Moscow",
			Address:      "Tverskaya St, 7, Moscow",
			Latitude:     55.7571,
			Longitude:    37.6156,
			WorkingHours: "09:00-22:00, daily",
			Phone:        "+7 (495) 100-10-10",
			Details:      "Main city center pickup, fast issue counter.",
		},
		{
			Name:         "Moscow North Pickup",
			City:         "Moscow",
			Address:      "Altufyevskoye Hwy, 86, Moscow",
			Latitude:     55.8798,
			Longitude:    37.5869,
			WorkingHours: "10:00-22:00, daily",
			Phone:        "+7 (495) 100-10-11",
			Details:      "Near Altufyevo metro.",
		},
		{
			Name:         "Moscow South Pickup",
			City:         "Moscow",
			Address:      "Kashirskoye Hwy, 61, Moscow",
			Latitude:     55.6101,
			Longitude:    37.7174,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (495) 100-10-12",
			Details:      "Convenient for south districts, parking available.",
		},
		{
			Name:         "Moscow West Pickup",
			City:         "Moscow",
			Address:      "Rublyovskoye Hwy, 52A, Moscow",
			Latitude:     55.7584,
			Longitude:    37.4082,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (495) 100-10-13",
			Details:      "West side pickup and return desk.",
		},
		{
			Name:         "Moscow East Pickup",
			City:         "Moscow",
			Address:      "Novokosinskaya St, 32, Moscow",
			Latitude:     55.7448,
			Longitude:    37.8632,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (495) 100-10-14",
			Details:      "East district pickup point.",
		},
		{
			Name:         "Moscow North-West Pickup",
			City:         "Moscow",
			Address:      "Mitinskaya St, 40, Moscow",
			Latitude:     55.8450,
			Longitude:    37.3622,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (495) 100-10-15",
			Details:      "Near Mitino metro.",
		},
		{
			Name:         "Moscow North-East Pickup",
			City:         "Moscow",
			Address:      "Mira Ave, 211, Moscow",
			Latitude:     55.8211,
			Longitude:    37.6380,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (495) 100-10-16",
			Details:      "Near VDNH, quick pickup line.",
		},
		{
			Name:         "Moscow South-East Pickup",
			City:         "Moscow",
			Address:      "Lyublinskaya St, 153, Moscow",
			Latitude:     55.6776,
			Longitude:    37.7612,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (495) 100-10-17",
			Details:      "South-east pickup point.",
		},
		{
			Name:         "Zelenograd Pickup",
			City:         "Moscow",
			Address:      "Panfilovsky Ave, 6A, Zelenograd",
			Latitude:     55.9825,
			Longitude:    37.1815,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (495) 100-10-18",
			Details:      "North-west Moscow area, easy car access.",
		},
		{
			Name:         "Khimki Pickup",
			City:         "Moscow Oblast",
			Address:      "Moskovskaya St, 14, Khimki",
			Latitude:     55.8892,
			Longitude:    37.4457,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-10",
			Details:      "Moscow region north hub.",
		},
		{
			Name:         "Mytishchi Pickup",
			City:         "Moscow Oblast",
			Address:      "Mira St, 32/2, Mytishchi",
			Latitude:     55.9104,
			Longitude:    37.7360,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-11",
			Details:      "Moscow region north-east hub.",
		},
		{
			Name:         "Korolyov Pickup",
			City:         "Moscow Oblast",
			Address:      "Kosmonavtov Ave, 34B, Korolyov",
			Latitude:     55.9162,
			Longitude:    37.8541,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-12",
			Details:      "Near railway station, fast pickup counter.",
		},
		{
			Name:         "Balashikha Pickup",
			City:         "Moscow Oblast",
			Address:      "Sovetskaya St, 9, Balashikha",
			Latitude:     55.7963,
			Longitude:    37.9382,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-13",
			Details:      "East Moscow region pickup point.",
		},
		{
			Name:         "Lyubertsy Pickup",
			City:         "Moscow Oblast",
			Address:      "Oktyabrsky Ave, 112, Lyubertsy",
			Latitude:     55.6765,
			Longitude:    37.8981,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-14",
			Details:      "South-east Moscow region hub.",
		},
		{
			Name:         "Podolsk Pickup",
			City:         "Moscow Oblast",
			Address:      "Kirova St, 39, Podolsk",
			Latitude:     55.4311,
			Longitude:    37.5455,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-15",
			Details:      "South Moscow region pickup hub.",
		},
		{
			Name:         "Odintsovo Pickup",
			City:         "Moscow Oblast",
			Address:      "Mozhayskoye Hwy, 58A, Odintsovo",
			Latitude:     55.6780,
			Longitude:    37.2777,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-16",
			Details:      "West Moscow region pickup point.",
		},
		{
			Name:         "Krasnogorsk Pickup",
			City:         "Moscow Oblast",
			Address:      "Lenina St, 38B, Krasnogorsk",
			Latitude:     55.8317,
			Longitude:    37.3299,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-17",
			Details:      "North-west Moscow region pickup point.",
		},
		{
			Name:         "Domodedovo Pickup",
			City:         "Moscow Oblast",
			Address:      "Kashirskoye Hwy, 3A, Domodedovo",
			Latitude:     55.4430,
			Longitude:    37.7470,
			WorkingHours: "10:00-21:00, daily",
			Phone:        "+7 (498) 200-20-18",
			Details:      "Near airport route, convenient parking.",
		},
	}
}
