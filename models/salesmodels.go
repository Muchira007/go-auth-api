package models

import (
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	gorm.Model
	Name        string
	Gender      string
	PhoneNumber string
	NationalID  string
	Geolocation string
	Country     string
	County      string
	Subcounty   string
	Village     string
	Sale        []Sale
}

type Sale struct {
	gorm.Model
	ProductID       uint
	Product         Product
	CustomerID      uint
	Customer        Customer
	DateOfSale      time.Time
	SerialNumber    string
	PaymentOption   string
	StatusOfAccount string
	Quantity        int
	Total           float64
}
