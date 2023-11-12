package mysql

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"social-network/internal/models"
	"strings"
)

type UserModel struct {
	DB *sql.DB

	RequestsModel FriendRequestModel
}

func (m *UserModel) Insert(
	firstName, lastName, interests, city, email, password string, gender models.Gender, age uint32,
) (int, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}

	stmt := `INSERT INTO users (created_at, first_name, last_name, age, gender, interests, city, email, hashed_password)
	VALUES(UTC_TIMESTAMP(), ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, firstName, lastName, age, gender, interests, city, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return 0, models.ErrDuplicateEmail
			}
		}

		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// TODO query friends and friend requests
func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := `
	SELECT 
	       id, 
		   created_at, 
		   first_name, 
		   last_name, 
		   age, 
		   gender, 
		   interests, 
		   city, 
		   email, 
		   hashed_password 
	FROM 
	       users
	WHERE 
	       id = ?`

	user := &models.User{}

	row := m.DB.QueryRow(stmt, id)

	err := row.Scan(
		&user.ID,
		&user.CreatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Gender,
		&user.Interests,
		&user.City,
		&user.Email,
		&user.HashedPassword,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	// Querying friend requests
	friendRequests, err := m.RequestsModel.ListAllForRecipient(id)
	if err != nil {
		return nil, err
	}

	user.FiendRequests = friendRequests

	// TODO friends list

	// Querying friend list
	/*stmt = `
	SELECT U.first_name AS friend_name
	FROM friendships AS F
	JOIN users AS U ON F.user_1_id = U.id
	WHERE F.user_2_id = ?
		UNION
	SELECT U.first_name AS friend_name
	FROM friendships AS F
	JOIN users AS U ON F.user_2_id = U.id
	WHERE F.user_1_id = ?;`

	row = m.DB.QueryRow(stmt, id)
	err = row.Scan(
		&user.Friends,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}*/

	// TODO

	return user, nil
}

func (m UserModel) Latest() ([]*models.User, error) {
	stmt := `
		SELECT 
			 id, 
			 created_at, 
			 first_name, 
			 last_name, 
			 age, 
			 gender, 
			 interests, 
			 city, 
			 email,
			 hashed_password 
		FROM 
			 users
		ORDER BY created_at DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		s := &models.User{}

		err = rows.Scan(
			&s.ID, &s.CreatedAt, &s.FirstName, &s.LastName, &s.Age, &s.Gender, &s.Interests, &s.City, &s.Email, &s.HashedPassword,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Authenticate verifies a user exists with the provided email address and password.
// If the user exists the relevant user ID is returned.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var (
		id       int
		hashedPw []byte
	)

	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`
	row := m.DB.QueryRow(stmt, email)

	if err := row.Scan(&id, &hashedPw); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	if err := bcrypt.CompareHashAndPassword(hashedPw, []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}
