package main

import (
	//"Email"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"../util"
	"encoding/json"
	"fmt"
)

type MSA struct {
	Users map[string] *util.User;
	domain string;
	NetworkAddress string;
}

func (msa *MSA) HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/toOutbox/{user}", msa.AddEmailToOutbox).Methods("POST")
	router.HandleFunc("/toInbox/{user}", msa.AddEmailToInbox).Methods("POST")
	router.HandleFunc("/listOutbox/{user}", msa.ListOutbox).Methods("GET")
	router.HandleFunc("/listInbox/{user}", msa.ListInbox).Methods("GET")
	router.HandleFunc("/deleteOutbox/{user}/{uuid}", msa.DeleteEmailFromOutbox).Methods("DELETE")
	router.HandleFunc("/deleteInbox/{user}/{uuid}", msa.DeleteEmailFromInbox).Methods("DELETE")
	router.HandleFunc("/peekOutbox/{user}", msa.PeekOutbox).Methods("GET")
	fmt.Println("MSA Service is running at " + msa.NetworkAddress)
	log.Fatal(http.ListenAndServe(msa.NetworkAddress, router));
}

//Creates a new user on the email server
//Users email address has to be unique
func (msa *MSA) CreateUser(user *util.User) {
	if !msa.exists(user.EmailAddress) {
		msa.Users[user.EmailAddress] = user;
	}
}


//Adds an email to the user outbox of the email source address
func (msa *MSA) AddEmailToOutbox(w http.ResponseWriter, r *http.Request) {
	fmt.Println(msa.NetworkAddress + " : adding email to outbox")
	decoder := json.NewDecoder(r.Body)
	var email util.Email;
	
	if err := decoder.Decode(&email); err == nil {

		if _, ok := msa.Users[email.Source]; ok {
			w.WriteHeader(http.StatusCreated)
			email.SetUUID()
			msa.Users[email.Source].AddEmailToOutbox(&email)

		} else {
			//the user does not exist on this MSA server
			w.WriteHeader(http.StatusNotFound)
		}
	
	  } else {
		//the JSON cannot be decoded
		w.WriteHeader(http.StatusBadRequest)
	  }
}


//Adds an email to the user inbox of the email source address
func (msa *MSA) AddEmailToInbox(w http.ResponseWriter, r *http.Request) {
	fmt.Println(msa.NetworkAddress + " : adding email to inbox")

	decoder := json.NewDecoder(r.Body)
	var email util.Email;
	
	if err := decoder.Decode(&email); err == nil {

		//If the user exists on this server
		if _, ok := msa.Users[email.Destination]; ok {
			w.WriteHeader(http.StatusCreated)
			email.SetUUID()
			msa.Users[email.Destination].AddEmailToInbox(&email)
			fmt.Println("Added email to " + email.Destination + " inbox")

		} else {
			//the user does not exist on this MSA server
			w.WriteHeader(http.StatusNotFound)
		}
	
	  } else {
		//the JSON cannot be decoded
		w.WriteHeader(http.StatusBadRequest)
	  }
}

//Lists all emails in the specified user's outbox
func (msa *MSA) ListOutbox(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	user := vars["user"]

	/*Check if the user exists*/
	if msa.exists(user) {
		w.WriteHeader(http.StatusOK)

		/*List the outbox*/
		for _, email := range msa.Users[user].Outbox {
			if enc, err := json.Marshal(email); err == nil {
				w.Write([]byte(enc))

			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	} else {
		//user does not exist on this server
		w.WriteHeader(http.StatusNotFound)
	}
}


//Lists all emails in the specified user's inbox
func (msa *MSA) ListInbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	/*Check if the user exists*/
	if msa.exists(user) {
		w.WriteHeader(http.StatusOK)
		
		/*List the outbox*/
		for _, email := range msa.Users[user].Inbox {
			if enc, err := json.Marshal(email); err == nil {
				w.Write([]byte(enc))

			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	} else {
		//user does not exist on this server
		w.WriteHeader(http.StatusNotFound)
	}
}

//Deletes a specific email from the user's outbox
//The email is identified using its UUID
func (msa *MSA) DeleteEmailFromOutbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	uuid := vars["uuid"]

	if msa.exists(user) {

		if ok := msa.Users[user].DeleteFromOutbox(uuid); ok {
			w.WriteHeader(http.StatusAccepted)
		} else {
			//no email has the specified UUID
			w.WriteHeader(http.StatusAccepted)
		}

	} else {
		//User does not exist on this server
		w.WriteHeader(http.StatusNotFound)
	}
}

//Deletes a specific email from the user's inbox
//The email is identified using its UUID
func (msa *MSA) DeleteEmailFromInbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	uuid := vars["uuid"]

	//Make sure the user exists on this server
	if msa.exists(user) {

		if ok := msa.Users[user].DeleteFromInbox(uuid); ok {
			w.WriteHeader(http.StatusAccepted)
		} else {
			//no email with the UUID can be found
			w.WriteHeader(http.StatusAccepted)
		}

	} else {
		//user does not exist on this server
		w.WriteHeader(http.StatusNotFound)
	}
}

//Returns the most recent email from the specified user's outbox
func (msa *MSA) PeekOutbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	//make sure user exists on this server
	if msa.exists(user) {

		w.WriteHeader(http.StatusOK)

		//get the latest email from the outbox
		var email *util.Email
		email = msa.Users[user].PeekOutbox()

		if enc, err := json.Marshal(email); err == nil {
			w.Write([]byte(enc))

		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

	} else {
		//user does not exist on this server
		w.WriteHeader(http.StatusNotFound)
	}
}


//Helper method to determine whether a user exists on this server
//Returns true if user exists, false if user does not exist
func (msa *MSA) exists(emailAddress string) bool {
	_, ok := msa.Users[emailAddress]
	return ok;
}


func main() {
	user := util.User{make([]*util.Email, 0), make([]*util.Email, 0), "fred@here.com"}
	msa := MSA{make(map[string]*util.User), "here.com", ":7001"}
	msa.CreateUser(&user)
	go msa.HandleRequests()

	msa2 := MSA{make(map[string]*util.User), "there.com", ":8001"}
	user2 := util.User{make([]*util.Email, 0), make([]*util.Email, 0), "fred@there.com"}
	msa2.CreateUser(&user2)
	msa2.HandleRequests()


}



