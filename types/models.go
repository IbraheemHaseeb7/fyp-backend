package types

import (
	"database/sql"
	"time"
)

type User struct {
	Name               string    		`json:"name"`
	DOB                string    		`json:"dob" gorm:"type:date"`
	ID                 int64     		`json:"id" gorm:"primaryKey"`
	Password           string    		`json:"password"`
	StudentCardURI     string    		`json:"studentCardURI"`
	LivePictureURI     string    		`json:"livePictureURI"`
	CreatedAt          time.Time 		`json:"createdAt" gorm:"created_at;autoCreateTime"`
	UpdatedAt          time.Time 		`json:"updatedAt" gorm:"updated_at;autoUpdateTime"`
	RegistrationNumber string    		`json:"registrationNumber" gorm:"unique"`
	Department         string    		`json:"department"`
	Semester           int8      		`json:"semester" gorm:"default:1"`
	Email              string    		`json:"email" gorm:"unique"`
	EmailVerifiedAt	   sql.NullTime		`json:"emailVerified" gorm:"default:null"`
	CardVerifiedAt	   sql.NullTime		`json:"cardVerified" gorm:"default:null"`
	SelfieVerifiedAt   sql.NullTime		`json:"selfieVerified" gorm:"default:null"`
	OTP				   sql.NullInt16	`json:"otp" gorm:"default:null"false`
}

type Vehicle struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	Type      string    `json:"type"`
	Make      string    `json:"make"`
	Model     string    `json:"model"`
	Year      int       `json:"year"`
	VIN       string    `json:"vin"`
	UserID    int64     `json:"userId"`
	FrontURI  string    `json:"frontUri"`
	BackURI   string    `json:"backUri"`
	InsideURI string    `json:"insideUri"`
	CreatedAt time.Time `json:"createdAt" gorm:"created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"updated_at;autoUpdateTime"`
	// User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// type Role struct {
// 	ID     int64  `json:"id" gorm:"primaryKey"`
// 	UserID int64  `json:"userId"`
// 	Role   string `json:"role"`
// 	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
// }

type Request struct {
	ID              int64     		`json:"id" gorm:"primaryKey"`
	UserID			int64 			`json:"user_id" validate:"required"`
	VehicleID		int64			`json:"vehicle_id"`
	Trunk			bool			`json:"trunk" validate:"required"`
	Persons 		uint8			`json:"persons" validate:"required"`
	VehicleType		string			`json:"vehicle_type" validate:"required"`
	FromLat			float64			`json:"from_lat" validate:"required"`
	FromLong		float64			`json:"from_long" validate:"required"`
	ToLat			float64			`json:"to_lat" validate:"required"`
	ToLong			float64			`json:"to_long" validate:"required"`
	Status			string			`json:"status" validate:"required"`
	OriginatorRole	string			`json:"originator_role" validate:"required"`
	RequestID		int64			`json:"request_id"`
	CreatedAt 		time.Time 		`json:"created_at" gorm:"created_at;autoCreateTime"`
	UpdatedAt 		time.Time 		`json:"updated_at" gorm:"updated_at;autoUpdateTime"`
}

type Ride struct {
	DriverID          int64     `json:"driverId"`
	PassengerID       int64     `json:"passengerId"`
	ID                int64     `json:"id" gorm:"primaryKey"`
	DistanceCovered   int       `json:"distanceCovered"`
	PaymentPercentage uint8     `json:"paymentPercentage"`
	StartTime         time.Time `json:"startTime"`
	EndTime           time.Time `json:"endTime"`
	// VehicleID int64 `json:"vehicleId"`
	// Vehicle Vehicle `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Feedbacks struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	RideID      int64  `json:"rideId"`
	Stars       uint8  `json:"stars"`
	Description string `json:"description"`
	// Ride Ride `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
