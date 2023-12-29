package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/imrcht/bed-n-breakfast/internals/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgressDBRepo) AllUsers() bool {
	return true
}

// * InsertReservation: inserts a reservation into the database
func (m *postgressDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// * Context is used to set a timeout for the query to maintain the transaction atomicity
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var resId int

	// * `returning id` is used to return the id of the inserted row and this makes the `insert statement` a `query`
	query := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, query,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Time{},
		time.Time{},
	).Scan(&resId)

	if err != nil {
		m.App.ErrorLog.Println(err)
		return 0, err
	}

	return resId, nil
}

// * InsertReservation: inserts a reservation into the database
func (m *postgressDBRepo) InsertRoomRestriction(roomRestriction models.RoomRestriction) error {
	// * Context is used to set a timeout for the query to maintain the transaction atomicity
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// * `returning id` is used to return the id of the inserted row and this makes the `insert statement` a `query`
	query := `insert into room_restrictions (start_date, end_date, reservation_id, room_id, restriction_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, query,
		roomRestriction.StartDate,
		roomRestriction.EndDate,
		roomRestriction.ReservationId,
		roomRestriction.RoomID,
		roomRestriction.RestrictionID,
		time.Time{},
		time.Time{},
	)

	if err != nil {
		m.App.ErrorLog.Println(err)
		return err
	}

	return nil
}

// * SearchAvailabilityByDates: returns true if room is available, and false if not available for a single room
func (m *postgressDBRepo) SearchAvailabilityByDatesByRoomId(start_date, end_date time.Time, roomId int) (bool, error) {
	// * Context is used to set a timeout for the query to maintain the transaction atomicity
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int

	query := `select count(id) from room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date`
	err := m.DB.QueryRowContext(ctx, query, roomId, start_date, end_date).Scan(&numRows)

	if err != nil {
		m.App.ErrorLog.Println(err)
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// * SearchAvailabilityForAllRoomsByDates: returns a slice of available rooms, if any, for a given date range
func (m *postgressDBRepo) SearchAvailabilityForAllRoomsByDates(start_date, end_date time.Time) ([]models.Room, error) {
	// * Context is used to set a timeout for the query to maintain the transaction atomicity
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// query := `select * from rooms inner join room_restrictions on rooms.id=room_restrictions.room_id where $1 < end_date and $2 > start_date`
	query := `select r.id, r.room_name from rooms r where r.id not in (select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`
	rows, err := m.DB.QueryContext(ctx, query, start_date, end_date)

	var availableRooms []models.Room
	if err != nil {
		m.App.ErrorLog.Println(err)
		return availableRooms, err
	}

	for rows.Next() {
		var room models.Room
		err = rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			m.App.ErrorLog.Println(err)
			return availableRooms, err
		}

		availableRooms = append(availableRooms, room)
	}

	if err = rows.Err(); err != nil {
		m.App.ErrorLog.Println(err)
		return availableRooms, err
	}

	return availableRooms, nil
}

// * SearchAvailabilityForAllRoomsByDates: returns a slice of available rooms, if any, for a given date range
func (m *postgressDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		m.App.ErrorLog.Println(err)
		return room, err
	}

	return room, nil
}

// * GetUserByEmail: returns a user by email
func (m *postgressDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1`
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		m.App.ErrorLog.Println(err)
		return user, err
	}

	return user, nil
}

// * UpdateUser: updates a user in the database
func (m *postgressDBRepo) UpdateUser(user models.User) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5 where id = $6`
	_, err := m.DB.ExecContext(ctx, query, user.FirstName, user.LastName, user.Email, user.AccessLevel, time.Now(), user.ID)

	if err != nil {
		m.App.ErrorLog.Println(err)
		return err
	}

	return nil
}

// * Authenticate: returns user id and hashed password if email and password are correct
func (m *postgressDBRepo) Authenticate(email, testPassword string) (models.User, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user models.User
	var hashedPassword string

	query := `select id, first_name, last_name, email, access_level, password from users where email = $1`
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.AccessLevel, &hashedPassword)

	if err != nil {
		return user, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return user, "", errors.New("incorrect password")
	} else if err != nil {
		return user, "", err
	}

	return user, hashedPassword, nil
}

// * InsertUser: inserts a user into the database
func (m *postgressDBRepo) InsertUser(user models.User) (models.User, error) {
	// * Context is used to set a timeout for the query to maintain the transaction atomicity
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// * `returning id` is used to return the id of the inserted row and this makes the `insert statement` a `query`
	query := `insert into users (first_name, last_name, email, password, access_level, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7) returning id`

	err := m.DB.QueryRowContext(ctx, query, user.FirstName, user.LastName, user.Email, user.Password, user.AccessLevel, time.Now(), time.Now()).Scan(&user.ID)

	if err != nil {
		m.App.ErrorLog.Println(err)
		return models.User{}, err
	}

	return user, nil
}
