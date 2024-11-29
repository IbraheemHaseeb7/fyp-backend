package handler

import (
	"encoding/json"
	"micro/database"
	"micro/utils"
			
	"gorm.io/gorm"
)

func Login(db *gorm.DB, payload json.RawMessage) interface{} {
	// extracting limit and offset for paginated data
	type LoginBody struct  {
		Email string
		Password string
	}
	var loginBody LoginBody
	err := json.Unmarshal([]byte(payload), &loginBody)
	utils.ErrorHandler(err)

	// fetching from database
	var users []database.User
	res := db.Find(&users, "email = ?", loginBody.Email)
	utils.ErrorHandler(res.Error)

	// creating a json string for response
	jsonData, err := json.Marshal(users)
	utils.ErrorHandler(err)

	return string(jsonData)	
}
