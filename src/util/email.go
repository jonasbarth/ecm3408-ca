package util

import (
	"github.com/google/uuid"
)

type Email struct {
	Source string `json:"source"`
	Destination string `json:"destination"`
	Body string `json:"body"`
	UUID string `json:"uuid"`
}

//Sets the UUID of the email
func (email *Email) SetUUID() bool {
	id, err := uuid.NewUUID()
	if err != nil {
		return false
	}
	email.UUID = id.String()
	return true
}
