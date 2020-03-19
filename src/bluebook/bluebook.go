package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"encoding/json"
	"log"
	"fmt"
)

type BlueBook struct {
	AddressBook map[string]string
}


//Finds the network address of an email server
func (blueBook *BlueBook) FindURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	emailAddress := vars["emailAddress"]

	//extract the domain of the email address as a substring
	if domain, ok := getDomain(emailAddress); ok {
		
		//ensure there exists a mapping between the domain and a network address
		if networkAddress, ok := blueBook.AddressBook[domain]; ok {
			w.WriteHeader(http.StatusOK)

			if enc, err := json.Marshal(networkAddress); err == nil {
				w.Write([]byte(enc))
	
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			//the domain does not exist on the blue book server
			w.WriteHeader(http.StatusNotFound)
			fmt.Println("Domain doesnt exist:")
			fmt.Println(domain)
		}

	} else {
		//the email address is of the wrong format
		w.WriteHeader(http.StatusNotFound)
		fmt.Println("Email address has wrong format:")
		fmt.Println(emailAddress)
	}

}

//Adds a mapping from a domain to a network address
func (blueBook *BlueBook) AddMapping(domain string, networkAddress string) {
	blueBook.AddressBook[domain] = networkAddress
}


func (blueBook *BlueBook) HandleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/findURL/{emailAddress}", blueBook.FindURL).Methods("GET")
	log.Fatal(http.ListenAndServe(":8887", router));

}

//Retrieves the domain of an email address
//Returns the domain and true is the email address is legal (legal = contains @ somewhere)
//Returns an empty string and false if the email address is illegal (illegal = does not contain @)
func getDomain(emailAddress string) (string, bool) {
	index := strings.Index(emailAddress, "@")


	if index != -1 {
		domain := emailAddress[index + 1:len(emailAddress)]
		return domain, true
	}
	return "", false
	
}



func main() {
	blueBook := BlueBook{make(map[string]string)}
	blueBook.AddMapping("here.com", "http://localhost:8888")
	blueBook.HandleRequests()
}

