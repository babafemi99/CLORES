package database

import (
	models "clores-local/Models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

const DNS = "root:aina4orosun@tcp(127.0.0.1:3306)/cloredb"

func Connect() {
	DB, err = gorm.Open(mysql.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(&models.User{}, &models.Product{}, &models.Link{}, &models.Order{}, models.OrderItem{})
}
