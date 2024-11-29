package handler

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Request struct {
	Entity    string `json:"entity"`
	Operation string `json:"operation"`
	Payload json.RawMessage `json:"payload"`
}

func GlobalHandler(request Request, db *gorm.DB) interface{} {

	entityMap := make(map[string]map[string]interface{})

	entityMap["users"] = map[string]interface{} {
		"create": func() interface{} { return UserCreate(db, request.Payload) },
		"read": func() interface{} { return UserRead(db, request.Payload) },
		"update": func() interface{} { return UserUpdate(db, request.Payload) },
		"delete": func() interface{} { return UserDelete(db, request.Payload) },
	}

	entityMap["auth"] = map[string]interface{} {
		"login": func() interface{} { return Login(db, request.Payload) },
	}

	return entityMap[request.Entity][request.Operation].(func() interface{})()
}
