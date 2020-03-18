package main

type User struct {
	Inbox []Email];
	Outbox []Email;
	EmailAddress string;
}

/*
func (user *User) ExistsInInbox(string uuid) bool {
	
	for _, email := range user.Inbox {
		if email.UUID == uuid {
			return true
		}
	}
	return false  
}


func (user *User) DeleteFromInbox(string uuid) bool {

	for _, email := range user.Inbox {
		if email.UUID == uuid {

			return true
		}
	}
	return false
}


func (user *User) AddEmailToOutbox(email *Email) {
	user.Outbox = append(msa.Users[email.Source].Outbox, email);
} */