package util

import (
	"github.com/google/uuid"
	"strings"
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




//Retrieves the domain of an email address
//Returns the domain and true is the email address is legal (legal = contains @ somewhere)
//Returns an empty string and false if the email address is illegal (illegal = does not contain @)
func GetDomain(emailAddress string) (string, bool) {
	index := strings.Index(emailAddress, "@")


	if index != -1 {
		domain := emailAddress[index + 1:len(emailAddress)]
		return domain, true
	}
	return "", false
	
}
