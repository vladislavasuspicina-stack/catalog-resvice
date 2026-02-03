package main

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"golang.org/x/text/encoding/charmap"
	"gorm.io/gorm"
)

func decodeWindows1251(b []byte) string {
	dec := charmap.Windows1251.NewDecoder()
	out, _ := dec.Bytes(b)
	return string(out)
}

type Product struct {
	ID   uint
	Name string
}

func main() {
	db, err := gorm.Open(sqlite.Open("./catalog.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	// Update names
	updates := map[uint]string{
		2: "Poco X3 Pro",
		4: "Redmi 10",
		5: "Самсунг Галакси S23",
		6: "Айфон 15 Про",
		7: "Сяоми Ми 13",
	}
	for id, name := range updates {
		if err := db.Model(&Product{}).Where("id = ?", id).Update("name", name).Error; err != nil {
			log.Printf("update error for id %d: %v", id, err)
		}
	}
	var products []Product
	if err := db.Model(&Product{}).Find(&products).Error; err != nil {
		log.Fatalf("query: %v", err)
	}
	for _, p := range products {
		fmt.Printf("%d: %s\n", p.ID, p.Name)
		fmt.Printf("bytes: % x\n", []byte(p.Name))
	}
}
