package main

import (
	//"Email"
	"github.com/gorilla/mux"
	"net/http"
	//"fmt"
	"encoding/json"
)

type MSA struct {
	Users map[string] *User;
	domain string;
}


func (msa *MSA) CreateUser(user *User) {
	if !msa.exists(user.EmailAddress) {
		msa.Users[user.EmailAddress] = user;
	}
}


func (msa *MSA) AddEmailToOutbox(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//user := vars["user"]
	decoder := json.NewDecoder(r.Body)
	
	var email Email;
	
	if err := decoder.Decode(&email); err == nil {
		if _, ok := msa.Users[email.Source]; ok {
			w.WriteHeader(http.StatusCreated)
			email.SetUUID()
			msa.Users[email.Source].AddEmailToOutbox(email)
		} else {
		//the user does not exist on this MSA serverelse {
			w.WriteHeader(http.StatusBadRequest)
		}
	
	  } else {
		w.WriteHeader( http.StatusBadRequest )
	  }
	
}


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
		w.WriteHeader(http.StatusNotFound)
	}
}


func (msa *MSA) ListInBox(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusNotFound)
	}
}


func (msa *MSA) DeleteEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	uuid := vars["uuid"]

	if msa.exists(user) {

		if msa.Users[user].ExistsInInbox(uuid) {
			



		} else {
			w.WriteHeader(http.StatusNotFound)
		}

	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}


func (msa *MSA) exists(emailAddress string) bool {
	_, ok := msa.Users[emailAddress]
	return ok;
}



