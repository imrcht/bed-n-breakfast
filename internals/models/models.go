package models

import (
	"time"
)

type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// Users: is the user model
type Users struct {
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
type Rooms struct {
	ID        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions: is the restriction model
type Restrictions struct {
	ID              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Reservations struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomId    int
	Room      Rooms
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RoomRestrictions: is the reservation model
type RoomRestrictions struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	ReservationId int
	Reservation   Reservations
	RoomID        int
	Room          Rooms
	RestrictionID int
	Restriction   Restrictions
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
