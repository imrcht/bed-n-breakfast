package models

import (
	"time"
)

// Users: is the user model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Rooms: is the room model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions: is the restriction model
type Restriction struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// RoomRestrictions: is the reservation model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	ReservationId int
	RoomID        int
	RestrictionID int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Reservation   Reservation
	Room          Room
	Restriction   Restriction
}
