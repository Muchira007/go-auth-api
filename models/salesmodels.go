package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	Name        string
	Gender      string
	PhoneNumber string
	CustomerID  uint // Ensure this matches your database schema
	Latitude    float64
	Longitude   float64
	Country     string
	County      string
	Subcounty   string
	Village     string
	// Ward        string
	Sale []Sale
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
	NationalID      uint
	Latitude        float64
	Longitude       float64
}
