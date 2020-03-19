package main

import (
	"user"

)






func main() {
	/*
	blueBook := BlueBook{make(map[string]string)}
	blueBook.AddMapping("here.com", "http://localhost:8889")
	go blueBook.HandleRequests()

	*/
	
	user := User{make([]*Email, 0), make([]*Email, 0), "fred@here.com"}
	msa := MSA{make(map[string]*User), "here.com", "http://localhost:8888"}
	msa.CreateUser(&user)
	/*
	go msa.HandleRequests()
	mta := MTA{make([]*Email, 0), "http://localhost:8889", "http://localhost:8888"}
	mta.PollMSA(user.EmailAddress) */
	
}