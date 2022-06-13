package main

import (
	"gorm.io/gorm"
)

type Dog struct {
	gorm.Model
	Name             string
	NumberOfMentions int
	BarkMentions     string
}

var allDogs []Dog

var negativeWords = []string{"no", "never", "don't", "rarely", "seldom", "won't"}

var barkWords = []string{"bark", "yap", "noise", "noisy"}
