package util


type User struct {
	Inbox []*Email;
	Outbox []*Email;
	EmailAddress string;
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


//Pops the last email found in the outbox
func (user *User) PopOutbox() *Email {
	var email *Email
	user.Outbox, email = remove(user.Outbox, len(user.Outbox) - 1)
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