package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string `gorm:"unique"`
	Description string
	Price       float64
	Quantity    int
	Color       string
	ImageData   []byte // Field to store image binary data
}
