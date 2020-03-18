package main

import (
	"github.com/google/uuid"
)

type Email struct {
	Source string `json:"source"`
	Destination string `json:"destination"`
	Body string `json:"body"`
	UUID string `json:"uuid"`
}


func (email *Email) SetUUID() bool {
	id, err := uuid.NewUUID()
	if err != nil {
		return false
	}
	email.UUID = id
	return true
}

/*
func New(source string, destination string, body string) Email {
	email := Email {source, destination, body};
	return email;
} */