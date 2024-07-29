package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	NationalID int `gorm:"unique"`
	FirstName  string
	SecondName string
	SurName    string
	PhoneNum   string `gorm:"unique"`
	Password   string
}
