package mysql

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"social-network/internal/models"
)

type FriendRequestModel struct {
	DB *sql.DB
}

func (m *FriendRequestModel) Insert(senderUserID, recipientUserID int) error {
	stmt := `INSERT INTO friend_requests (created_at, sender_user_id, recipient_user_id)
	VALUES(UTC_TIMESTAMP(), ?, ?)`

	if _, err := m.DB.Exec(stmt, senderUserID, recipientUserID); err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 {
				return models.ErrDuplicateFriendRequest
			}
		}

		return err
	}

	return nil
}

// TODO test
func (m *FriendRequestModel) ListAllForRecipient(recipient int) ([]*models.FriendRequest, error) {
	stmt := `
	SELECT 
		   created_at, 
		   sender_user_id, 
		   recipient_user_id
	FROM 
	       friend_requests
	WHERE 
	       recipient_user_id = ?`

	var allFriendReq []*models.FriendRequest

	rows, err := m.DB.Query(stmt, recipient)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		req := &models.FriendRequest{}

		err = rows.Scan(
			&req.CreatedAt, &req.SenderUserID, &req.RecipientUserID,
		)
		if err != nil {
			return nil, err
		}
		allFriendReq = append(allFriendReq, req)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allFriendReq, nil
}

func (m *FriendRequestModel) Get(sender, recipient int) (*models.FriendRequest, error) {
	stmt := `
	SELECT 
		   created_at, 
		   sender_user_id, 
		   recipient_user_id
	FROM 
	       friend_requests
	WHERE 
	       sender_user_id = ? AND recipient_user_id = ?`

	friendReq := &models.FriendRequest{}

	row := m.DB.QueryRow(stmt, sender, recipient)

	err := row.Scan(
		&friendReq.CreatedAt,
		&friendReq.SenderUserID,
		&friendReq.RecipientUserID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return friendReq, nil
}
