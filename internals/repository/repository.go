package repository

import (
	"time"

	"github.com/imrcht/bed-n-breakfast/internals/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomId(start_date, end_date time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRoomsByDates(start_date, end_date time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)

	GetUserById(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, testPassword string) (models.User, string, error)
	InsertUser(user models.User) (models.User, error)
}
