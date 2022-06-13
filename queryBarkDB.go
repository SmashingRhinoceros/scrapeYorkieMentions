package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("newBarkData.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	var dogs []Dog

	db.Order("number_of_mentions desc, name").Where("number_of_mentions >= ?", 5).Find(&dogs)
	fmt.Println(dogs)

	for _, d := range dogs {
		fmt.Println(d.Name, d.NumberOfMentions)

	}

}
