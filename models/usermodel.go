package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	NationalID int    `gorm:"unique"`
	Email      string `gorm:"unique"`
	FirstName  string
	SecondName string
	SurName    string
	PhoneNum   string `gorm:"unique"`
	Password   string
	ResetToken string

	ResetTokenExpiry time.Time
}
