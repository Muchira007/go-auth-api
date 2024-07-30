package initializers

import "github.com/muchira007/jambo-green-go/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{}, &models.Sale{}, &models.Customer{})
}
