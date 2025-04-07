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
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	// User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Role struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	UserID int64  `json:"userId"`
	Role   string `json:"role"`
	// User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
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
