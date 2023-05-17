package database

import (
	"fmt"
	

	"github.com/LucaWilliams4831/liquiditysniperbot/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	
	HOST := "localhost"
	DBUSER := "postgres"
	PASSWORD := "postgres"
	DBNAME := "bdjuno"
	PORT := "5432"
	
	config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, DBUSER, PASSWORD, DBNAME)
	fmt.Println(config)
	connection, err := gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}
	
	DB = connection

	connection.AutoMigrate(&models.Admin{},&models.Account{}) 
}
