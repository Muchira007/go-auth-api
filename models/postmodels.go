package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Description string `gorm:"unique"`
	Price       float64
	Quantity    int
}
