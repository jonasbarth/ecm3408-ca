package util

//User representation on an email server

type User struct {
	Inbox []*Email; //The user's inbox
	Outbox []*Email; //The user's outbox
	EmailAddress string; //Email address of the user
}


//Deletes the email with the specified UUID from the inbox
//Returns true if the email exists and is deleted
//Returns false if the email does not exist
func (user *User) DeleteFromInbox(uuid string) bool {

	for index, email := range user.Inbox {
	
		if email.UUID == uuid {
			user.Outbox, _ = remove(user.Inbox, index)
			return true
		}
	}
	return false
}


//Deletes the email with the specified UUID from the outbox
//Returns true if the email exists and is deleted
//Returns false if the email does not exist
func (user *User) DeleteFromOutbox(uuid string) bool {

	for index, email := range user.Outbox {
		
		if email.UUID == uuid {
			user.Outbox, _ = remove(user.Outbox, index)
			return true
		}
	}
	return false
}


//Adds an email to the users outbox
func (user *User) AddEmailToOutbox(email *Email) {
	user.Outbox = append(user.Outbox, email);
}

//Adds an email to the users inbox
func (user *User) AddEmailToInbox(email *Email) {
	user.Inbox = append(user.Inbox, email)
}


//Pops the most recent email found in the outbox
func (user *User) PopOutbox() *Email {
	var email *Email
	user.Outbox, email = remove(user.Outbox, len(user.Outbox) - 1)
	return email
}


//Peeks at the most recent email in the outbox
func (user *User) PeekOutbox() *Email {
	if len(user.Outbox) == 0 {
		return nil
	}
	email := user.Outbox[len(user.Outbox)-1]
	return email
}


//Peeks at the most recent email in the inbox
func (user *User) PeekInbox() *Email {
	if len(user.Inbox) == 0 {
		return nil
	}
	email := user.Inbox[len(user.Inbox)-1]
	return email
}


//Removes an element from an Email slice at the given index
//Returns the new slice and the removed element if found
//Returns the original box and nil if the box is empty
func remove(box []*Email, index int) ([]*Email, *Email) {
	if len(box) == 0 {
		return box, nil
	}
    box[len(box)-1], box[index] = box[index], box[len(box)-1]
    return box[:len(box)-1], box[index]
}