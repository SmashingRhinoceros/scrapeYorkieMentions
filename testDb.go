package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cat struct {
	gorm.Model
	Name string
	Age  int
}

func main() {
	fmt.Println("hi")
	db, err := gorm.Open(sqlite.Open("test2.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Cat{})
	db.Create(&Cat{Name: "Pud", Age: 16})
	db.Create(&Cat{Name: "Mako", Age: 19})

	var cat Cat
	pud := db.First(&cat, "Name = ?", "Pud")
	fmt.Println(pud)

}
