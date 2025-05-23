package types

import (
	"database/sql"
	"time"
)

type User struct {
	Name               string    		`json:"name,omitempty"`
	DOB                string    		`json:"dob,omitempty" gorm:"type:date"`
	ID                 int64     		`json:"id,omitempty" gorm:"primaryKey"`
	Password           string    		`json:"password,omitempty"`
	StudentCardURI     string    		`json:"studentCardURI,omitempty"`
	LivePictureURI     string    		`json:"livePictureURI,omitempty"`
	CreatedAt          time.Time 		`json:"createdAt,omitempty" gorm:"created_at;autoCreateTime"`
	UpdatedAt          time.Time 		`json:"updatedAt,omitempty" gorm:"updated_at;autoUpdateTime"`
	RegistrationNumber string    		`json:"registrationNumber,omitempty" gorm:"unique"`
	Department         string    		`json:"department,omitempty"`
	Semester           int8      		`json:"semester,omitempty" gorm:"default:1"`
	Email              string    		`json:"email,omitempty" gorm:"unique"`
	EmailVerifiedAt	   sql.NullTime		`json:"emailVerified,omitempty" gorm:"default:null"`
	CardVerifiedAt	   sql.NullTime		`json:"cardVerified,omitempty" gorm:"default:null"`
	SelfieVerifiedAt   sql.NullTime		`json:"selfieVerified,omitempty" gorm:"default:null"`
	OTP				   sql.NullInt16	`json:"otp,omitempty" gorm:"default:null"false`
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
	VehicleID		*int64			`json:"vehicle_id,omitempty"`
	Trunk			bool			`json:"trunk,omitempty"`
	Persons 		uint8			`json:"persons" validate:"required"`
	VehicleType		string			`json:"vehicle_type" validate:"required"`
	FromLat			float64			`json:"from_lat" validate:"required"`
	FromLong		float64			`json:"from_long" validate:"required"`
	ToLat			float64			`json:"to_lat" validate:"required"`
	ToLong			float64			`json:"to_long" validate:"required"`
	ToName			*string			`json:"to_name" validate:"required"`
	FromName		*string			`json:"from_name" validate:"required"`
	Status			string			`json:"status" validate:"required"`
	OriginatorRole	string			`json:"originator_role" validate:"required"`
	RequestID		int64			`json:"request_id"`
	Until			time.Time		`json:"until" validate:"required"`
	CreatedAt 		time.Time 		`json:"created_at" gorm:"created_at;autoCreateTime"`
	UpdatedAt 		time.Time 		`json:"updated_at" gorm:"updated_at;autoUpdateTime"`
	Vehicle			Vehicle			
	User			User
}

// make RequestID and ProposalID composite key
type Room struct {
	ID              int64     		`json:"id" gorm:"primaryKey"`
	RequestID		int64 			`json:"request_id" validate:"required" gorm:"uniqueIndex:idx_request_proposal"`
	ProposalID		int64 			`json:"proposal_id" validate:"required" gorm:"uniqueIndex:idx_request_proposal"`
	CreatedAt 		time.Time 		`json:"created_at" gorm:"created_at;autoCreateTime"`
	UpdatedAt 		time.Time 		`json:"updated_at" gorm:"updated_at;autoUpdateTime"`
	Proposal		Request			
	Request			Request			
}

type Message struct {
	ID              int64     		`json:"id" gorm:"primaryKey"`
	RoomID			int64 			`json:"room_id" validate:"required"`
	UserID			int64 			`json:"user_id" validate:"required"`
	Message			string			`json:"message" validate:"required"`
	CreatedAt 		time.Time 		`json:"created_at" gorm:"created_at;autoCreateTime"`
	UpdatedAt 		time.Time 		`json:"updated_at" gorm:"updated_at;autoUpdateTime"`
	Room			Room			`json:"room" gorm:"foreignKey:RoomID;references:ID"`
	User			User			`json:"user" gorm:"foreignKey:UserID;references:ID"`
}

type Ride struct {
	DriverID          int64     `json:"driverId"`
	PassengerID       int64     `json:"passengerId"`
	ID                int64     `json:"id" gorm:"primaryKey"`
	VehicleID         int64     `json:"vehicle_id"`
	RequestID         int64     `json:"request_id"`
	ProposalID		  int64     `json:"proposal_id"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Status		   	  string    `json:"status" gorm:"default:'reaching_passenger'"`
	Passenger 	   	  User     
	Driver 		   	  User     
	Vehicle 	   	  Vehicle
	Request			  Request
	Proposal		  Request
	CreatedAt         time.Time `json:"createdAt" gorm:"created_at;autoCreateTime"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"updated_at;autoUpdateTime"`
	// DistanceCovered   int       `json:"distanceCovered"`
	// PaymentPercentage uint8     `json:"paymentPercentage"`
}

type Feedbacks struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	RideID      int64  `json:"ride_id"`
	UserID      int64  `json:"user_Id"`
	Stars       uint8  `json:"stars"`
	Description string `json:"description"`
	Ride 		Ride   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	User 		User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
