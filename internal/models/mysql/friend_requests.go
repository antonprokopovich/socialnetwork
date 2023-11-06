package mysql

import (
	"database/sql"
	"errors"
	"social-network/internal/models"
)

type FriendRequestModel struct {
	DB *sql.DB
}

func (m *FriendRequestModel) Insert(senderUserID, recipientUserID int64) error {
	stmt := `INSERT INTO friend_requests (created_at, sender_user_id, recipient_user_id)
	VALUES(UTC_TIMESTAMP(), ?, ?)`

	if _, err := m.DB.Exec(stmt, senderUserID, recipientUserID); err != nil {
		return err
	}

	return nil
}

func (m *FriendRequestModel) ListAllForRecipient(recipient int) (*models.FriendRequest, error) {
	// TODO
	return nil, nil
}

func (m *FriendRequestModel) Get(id int) (*models.FriendRequest, error) {
	stmt := `
	SELECT 
	       id, 
		   created_at, 
		   sender_user_id, 
		   recipient_user_id
	FROM 
	       friend_requests
	WHERE 
	       id = ?`

	friendReq := &models.FriendRequest{}

	row := m.DB.QueryRow(stmt, id)

	err := row.Scan(
		&friendReq.ID,
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
