package handler

import (
	"encoding/json"
	"micro/database"
	"micro/utils"

	"gorm.io/gorm"
)

type LimitAndOffset struct {
	Limit int
	Offset int
}

func UserCreate(db *gorm.DB, payload json.RawMessage) interface{} {
	// storing and mapping payload in user struct/model
	var user database.User
	err := json.Unmarshal([]byte(payload), &user)
	utils.ErrorHandler(err)

	// storing in database
	res := db.Create(&user)
	utils.ErrorHandler(res.Error)
	return "Successfully added a new user"
}

func UserRead(db *gorm.DB, payload json.RawMessage) interface{} {
	// extracting limit and offset for paginated data
	var limAndOff LimitAndOffset
	err := json.Unmarshal([]byte(payload), &limAndOff)
	utils.ErrorHandler(err)

	// fetching from database
	var users []database.User
	res := db.Limit(limAndOff.Limit).Offset(limAndOff.Offset).Find(&users)
	utils.ErrorHandler(res.Error)

	// creating a json string for response
	jsonData, err := json.Marshal(users)
	utils.ErrorHandler(err)

	return string(jsonData)	
}

func UserUpdate(db *gorm.DB, payload json.RawMessage) interface{} {
	// storing and mapping payload in user struct/model
	var user database.User
	err := json.Unmarshal([]byte(payload), &user)
	utils.ErrorHandler(err)

	// storing in database
	res := db.Save(&user)
	utils.ErrorHandler(res.Error)

	return string(payload)
}

func UserDelete(db *gorm.DB, payload json.RawMessage) interface{} {
	type DeleteRequest struct {
		Id int16
	}
	var req DeleteRequest
	err := json.Unmarshal([]byte(payload), &req)
	utils.ErrorHandler(err)
	db.Delete(&database.User{}, req.Id)
	return "Successfully deleted a new user"
}
