package main

import (
	//"Email"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"../../util"
	"encoding/json"
	"fmt"
)

type MSA struct {
	Users map[string] *util.User;
	Domain string;
	NetworkAddress string;
}

func (msa *MSA) HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/outbox", msa.AddEmailToOutbox).Methods("POST")
	router.HandleFunc("/inbox", msa.AddEmailToInbox).Methods("POST")
	router.HandleFunc("/outbox/{user}", msa.ListOutbox).Methods("GET")
	router.HandleFunc("/inbox/{user}", msa.ListInbox).Methods("GET")
	router.HandleFunc("/outbox/{user}/{uuid}", msa.DeleteEmailFromOutbox).Methods("DELETE")
	router.HandleFunc("/inbox/{user}/{uuid}", msa.DeleteEmailFromInbox).Methods("DELETE")
	router.HandleFunc("/peekOutbox/{user}", msa.PeekOutbox).Methods("GET")
	router.HandleFunc("/peekInbox/{user}", msa.PeekInbox).Methods("GET")
	router.HandleFunc("/users", msa.ListUsers).Methods("GET")
	fmt.Println("MSA Service is running at " + msa.NetworkAddress)
	log.Fatal(http.ListenAndServe(msa.NetworkAddress, router));
}

//Creates a new user on the email server
//Users email address has to be unique
func (msa *MSA) CreateUser(emailAddress string) bool {
	domain, _ := util.GetDomain(emailAddress)

	//If the user does not already exist on this server and the domain of the email address matches that of this email server
	if (!msa.exists(emailAddress) && domain == msa.Domain) {
		fmt.Printf("Creating user %s on email server %s\n", emailAddress, msa.Domain)
		user := util.User{make([]*util.Email, 0), make([]*util.Email, 0), emailAddress}
		msa.Users[emailAddress] = &user
		return true
	
	//Case where the user already exists and the domain of the email addresses is a match for this domain
	} else if (msa.exists(emailAddress) && domain == msa.Domain) {
		return true
	}
	return false
}



//Adds an email to the user outbox of the email source address
func (msa *MSA) AddEmailToOutbox(w http.ResponseWriter, r *http.Request) {
	

	decoder := json.NewDecoder(r.Body)
	var email util.Email;
	
	//Decode the JSON into an email struct
	if err := decoder.Decode(&email); err == nil {
		
		//Create a user on this email server if necessary
		if ok := msa.CreateUser(email.Source); ok {
			
			//Make sure the user exists on this erver
			if _, ok := msa.Users[email.Source]; ok {
				w.WriteHeader(http.StatusCreated)
				email.SetUUID()
				msa.Users[email.Source].AddEmailToOutbox(&email)
				fmt.Println(msa.NetworkAddress + " : adding email to outbox")
	
			} else {
				//the user does not exist on this MSA server
				fmt.Printf("User does not exist on this server\n")
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			//The domain of the email does not correspond to this email server
			fmt.Printf("The domain of %s does not correspond to this domain %s\n", email.Source, msa.Domain)
			w.WriteHeader(http.StatusBadRequest)
		}

	} else {
		//the JSON cannot be decoded
		fmt.Printf("JSON cannot be decoded with error %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}


//Adds an email to the user inbox of the email source address
func (msa *MSA) AddEmailToInbox(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var email util.Email;
	
	if err := decoder.Decode(&email); err == nil {

		//If the user exists on this server
		if msa.exists(email.Destination) {
			w.WriteHeader(http.StatusCreated)
			email.SetUUID()
			msa.Users[email.Destination].AddEmailToInbox(&email)
			fmt.Println(msa.NetworkAddress + " : adding email to " + email.Destination + "inbox")

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

	//make sure the user exists on this server
	if msa.exists(user) {
		
		//The email was successfully deleted from the server
		if ok := msa.Users[user].DeleteFromOutbox(uuid); ok {
			fmt.Printf("Email with uuid %s successfully deleted from %s inbox\n", uuid, user)
			w.WriteHeader(http.StatusAccepted)
		} else {
			//no email has the specified UUID
			fmt.Printf("Email with uuid %s does not exist in %s inbox\n", uuid, user)
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


//Returns the most recent email from the specified user's inbox
func (msa *MSA) PeekInbox(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]

	//make sure user exists on this server
	if msa.exists(user) {

		w.WriteHeader(http.StatusOK)

		//get the latest email from the outbox
		var email *util.Email
		email = msa.Users[user].PeekInbox()

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

//Lists all the users of this email server so that the MTA can utilise
//the list for polling all of their outboxes
func (msa *MSA) ListUsers(w http.ResponseWriter, r *http.Request) {

	
	if enc, err := json.Marshal(msa.getUsers()); err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(enc))

	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	
}


//Helper function to determine whether a user exists on this server
//Returns true if user exists, false if user does not exist
func (msa *MSA) exists(emailAddress string) bool {
	_, ok := msa.Users[emailAddress]
	return ok;
}

//Helper function to get all users on this email server
func (msa *MSA) getUsers() []string {
	
	users := []string{}
	for user := range msa.Users {
		users = append(users, user)
	}
	return users
}

func main() {
	msa := MSA{make(map[string]*util.User), "here.com", ":7001"}

	go msa.HandleRequests()

	msa2 := MSA{make(map[string]*util.User), "there.com", ":8001"}

	msa2.HandleRequests()
}



