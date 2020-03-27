package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"../../util"
	"encoding/json"
	"log"
	"fmt"
)

type BlueBook struct {
	AddressBook map[string]string
	BluebookAddress string
}


//Finds the network address of an email server
func (blueBook *BlueBook) FindURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	emailAddress := vars["emailAddress"]

	//extract the domain of the email address as a substring
	if domain, ok := util.GetDomain(emailAddress); ok {
		
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
	fmt.Println(domain + " is available at " + networkAddress)
}


func (blueBook *BlueBook) HandleRequests() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/address/{emailAddress}", blueBook.FindURL).Methods("GET")
	fmt.Println("Bluebook Server is running at: " + blueBook.BluebookAddress)
	log.Fatal(http.ListenAndServe(blueBook.BluebookAddress, router));

}





func main() {
	blueBook := BlueBook{make(map[string]string), ":9000"}
	blueBook.AddMapping("here.com", "http://mta:7000")
	blueBook.AddMapping("there.com", "http://mta:8000")
	
	blueBook.HandleRequests()
}

