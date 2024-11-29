package main

import (
	"micro/database"
	"micro/utils"
	"micro/broker"

	"github.com/joho/godotenv"
)

func main() {
	// getting env variables
	err := godotenv.Load(".env")
	utils.ErrorHandler(err)

	// connecting to database
	db := utils.ConnectToDB()

	// running migrations
	db.AutoMigrate(&database.User{})
	database.RunMigrations(db)

	// setting up rabbit mq
	broker.RabbitMQ(db)
}

