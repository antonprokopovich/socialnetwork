package models

import (
	"errors"
	"time"
)

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Gender string

type User struct {
	ID        int64     `json:"ID,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Age       uint32    `json:"age,omitempty"`
	Gender    Gender    `json:"gender,omitempty"`
	Interests string    `json:"interests,omitempty"`
	City      string    `json:"city,omitempty"`
}
