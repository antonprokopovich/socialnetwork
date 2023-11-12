package models

import (
	"errors"
	"time"
)

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

var (
	ErrNoRecord               = errors.New("models: no matching record found")
	ErrInvalidCredentials     = errors.New("models: invalid credentials")
	ErrDuplicateEmail         = errors.New("models: duplicate email")
	ErrDuplicateFriendRequest = errors.New("models: duplicate friend request")
)

type Gender string

type User struct {
	ID             int64     `json:"id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	FirstName      string    `json:"first_name,omitempty"`
	LastName       string    `json:"last_name,omitempty"`
	Age            uint32    `json:"age,omitempty"`
	Gender         Gender    `json:"gender,omitempty"`
	Interests      string    `json:"interests,omitempty"`
	City           string    `json:"city,omitempty"`
	Email          string    `json:"email,omitempty"`
	HashedPassword string    `json:"hashed_password,omitempty"`
	// Friends - список ID пользователей, которых данный пользователь добавил в друзья.
	Friends []int `json:"friends,omitempty"`
	// FiendRequests - список пользователей, которые отправили данному пользователю запрос на добавление в друзья.
	FiendRequests []*FriendRequest `json:"friend_requests,omitempty"`
}

type Friendship struct {
	ID        int64     `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	User1ID   int       `json:"user_1_id,omitempty"`
	User2ID   int       `json:"user_2_id,omitempty"`
}

type FriendRequest struct {
	CreatedAt       time.Time `json:"created_at"`
	SenderUserID    int       `json:"user_1_id,omitempty"`
	RecipientUserID int       `json:"user_2_id,omitempty"`
}
