package database

import "time"

type User struct {
	Name		   string
	Dob			   time.Time		`gorm:"type:date"`
	Id			   int64			`gorm:"primaryKey"`
	StudentCardUri string
	LivePictureUri string
	CreatedAt	   time.Time		`gorm:"autoCreateTime"`
	UpdatedAt	   time.Time		`gorm:"autoUpdateTime"`
}
