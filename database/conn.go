package database

import (
	"fmt"
	"os"

	"github.com/LucaWilliams4831/uniswap-pancakeswap-tradingbot/liquiditysniperbot/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	HOST := os.Getenv("HOST")
	DBUSER := os.Getenv("DBUSER")
	PASSWORD := os.Getenv("PASSWORD")
	DBNAME := os.Getenv("DBNAME")
	PORT := os.Getenv("DB_PORT")
	
	config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", HOST, PORT, DBUSER, PASSWORD, DBNAME)
	
	connection, err := gorm.Open(postgres.Open(config), &gorm.Config{})

	if err != nil {
		panic("could not connect to the database")
	}
	fmt.Println("Database connected...")
	DB = connection

	connection.AutoMigrate(&models.Admin{},&models.Account{}) 
}
