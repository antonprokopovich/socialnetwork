package mysql

import (
	"database/sql"
	"errors"
	"social-network/internal/models"
)

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(
	firstName, lastName, interests, city string, gender models.Gender, age uint32,
) (int, error) {

	stmt := `INSERT INTO users (created_at, first_name, last_name, age, gender, interests, city)
	VALUES(UTC_TIMESTAMP(), ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, firstName, lastName, age, gender, interests, city)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
	stmt := `SELECT id, created_at, first_name, last_name, age, gender, interests, city FROM users
	WHERE id = ?`

	s := &models.User{}

	row := m.DB.QueryRow(stmt, id)

	err := row.Scan(
		&s.ID, &s.CreatedAt, &s.FirstName, &s.LastName, &s.Age, &s.Gender, &s.Interests, &s.City,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m UserModel) Latest() ([]*models.User, error) {
	stmt := `SELECT id, created_at, first_name, last_name, age, gender, interests, city FROM users
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
			&s.ID, &s.CreatedAt, &s.FirstName, &s.LastName, &s.Age, &s.Gender, &s.Interests, &s.City,
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
