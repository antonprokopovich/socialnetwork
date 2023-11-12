package mysql

import (
	"database/sql"
	"social-network/internal/models"
)

type FriendshipModel struct {
	DB *sql.DB
}

func (m *FriendshipModel) Insert(firstUserID, secondUserID int) error {
	stmt := `INSERT INTO friendships (created_at, user_1_id, user_2_id)
	VALUES(UTC_TIMESTAMP(), ?, ?)`

	if _, err := m.DB.Exec(stmt, firstUserID, secondUserID); err != nil {
		return err
	}

	return nil
}

func (m *FriendshipModel) Delete(firstUserID, secondUserID int) error {
	stmt := `DELETE FROM friendships WHERE (user_1_id = ? AND user_2_id = ?) OR (user_2_id = ? AND user_1_id = ?)`

	if _, err := m.DB.Exec(stmt, firstUserID, secondUserID, firstUserID, secondUserID); err != nil {
		return err
	}

	return nil
}

func (m *FriendshipModel) ListAllForUser(id int) ([]*models.Friendship, error) {
	stmt := `
	SELECT 
		   created_at, 
		   user_1_id, 
		   user_2_id
	FROM 
	       friendships
	WHERE 
	       user_1_id = ? OR user_2_id = ?`

	var friends []*models.Friendship

	rows, err := m.DB.Query(stmt, id, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		req := &models.Friendship{}

		err = rows.Scan(
			&req.CreatedAt, &req.User1ID, &req.User2ID,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

/*func (m *FriendshipModel) Get(id int) (*models.Friendship, error) {
	stmt := `
	SELECT
		   created_at,
		   user_1_id,
		   user_2_id
	FROM
	       friendships
	WHERE
	       user_1_id = ? AND user_2_id = ?`

	friendship := &models.Friendship{}

	row := m.DB.QueryRow(stmt, id)

	err := row.Scan(
		&friendship.ID,
		&friendship.CreatedAt,
		&friendship.User1ID,
		&friendship.User2ID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return friendship, nil
}
*/
