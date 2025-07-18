package db

import (
	"fmt"
	"os"

	"github.com/IbraheemHaseeb7/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DB_STRING")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := db.AutoMigrate(&types.User{}, &types.Vehicle{}, &types.Request{}, &types.Room{}, &types.Message{}, &types.Ride{}, &types.Feedback{}); err != nil {
		fmt.Println(err.Error())
	}

	DB = db
}
