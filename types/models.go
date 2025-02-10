package types

import "time"

type User struct {
	Name               string    `json:"name"`
	DOB                string    `json:"dob" gorm:"type:date"`
	ID                 int64     `json:"id" gorm:"primaryKey"`
	Password           string    `json:"password"`
	StudentCardURI     string    `json:"studentCardURI"`
	LivePictureURI     string    `json:"livePictureURI"`
	CreatedAt          time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	RegistrationNumber string    `json:"registrationNumber" gorm:"unique"`
	Department         string    `json:"department"`
	Semester           int8      `json:"semester" gorm:"default:1"`
	Email              string    `json:"email" gorm:"unique"`
}

type Vehicle struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	Type   string `json:"type"`
	Make   string `json:"make"`
	Model  string `json:"model"`
	Year   int    `json:"year"`
	VIN    string `json:"vin"`
	UserID int64  `json:"user_id"`
}
