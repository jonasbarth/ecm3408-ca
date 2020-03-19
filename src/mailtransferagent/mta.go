package main

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"fmt"
	"../util"
	"github.com/gorilla/mux"
	"log"
)

type MTA struct {
	OutEmails []*util.Email
	BlueBookURL string
}

//Posts the email into a users inbox via the MSA microservice
func (mta *MTA) PostEmail(w http.ResponseWriter, r *http.Request) {

	//get body of the eamil
	decoder := json.NewDecoder(r.Body)
	var email util.Email;

	//make sure the email follows the correct format
	if err := decoder.Decode(&email); err == nil {
		
		//Marshal the email object into JSON
		if enc, err := json.Marshal(email); err == nil {
			
			//Create the URL 
			MSAPostURL := mta.getURL(email.Destination) + "/toInbox/" + email.Destination

			//Create and make the POST request
			if req, err1 := http.NewRequest("POST", MSAPostURL, bytes.NewBuffer(enc)); err1 == nil {
				
				client := &http.Client {}
				//Get the response
				if resp, err2 := client.Do(req); err2 == nil {

					if _, err3 := ioutil.ReadAll(resp.Body); err3 == nil {
						w.WriteHeader(http.StatusOK)
					}
				} else {
					//POST request failed
					w.WriteHeader(http.StatusInternalServerError)
				}

			} else {
				//POST request failed
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			//The JSON cannot be marshalled
			w.WriteHeader(http.StatusInternalServerError)
		}

	} else {
		//the JSON cannot be decoded
		w.WriteHeader(http.StatusBadRequest)
	}
}

//Polls its own MSA microservice every x seconds and reads and deletes the latest email
func (mta *MTA) PollMSA(emailAddress string) {

	var email *util.Email
	//Get the email from outbox of the MSA to which the email belongs
	if body, ok := mta.popOutbox(emailAddress); ok {

		//unmarshal the email
		if err := json.Unmarshal(body, email); err == nil {
			//deliver the email to the MTA of the destination address
		} else {
			fmt.Println("Cannot unmarshal JSON")
		}
	} else {
		fmt.Println("Failed to pop from outbox")
	}

	

}


func (mta *MTA) getURL(email string) string {

	url := mta.BlueBookURL + "/" + email
	
	if req, err1 := http.NewRequest("GET", url, nil); err1 == nil {
				
		client := &http.Client {}
		//Get the response
		if resp, err2 := client.Do(req); err2 == nil {

			if body, err3 := ioutil.ReadAll(resp.Body); err3 == nil {
				fmt.Println(string(body))
				return string(body)[1:len(string(body))-1]
			}
		} else {
			//GET request failed
			return ""
		}

	} else {
		//GET request failed
		return ""
	}
	return ""
}


func (mta *MTA) popOutbox(emailAddress string) ([]byte, bool) {
	//Create the URL
	MSAGetURL := fmt.Sprintf("%s%s%s", mta.getURL(emailAddress), "/popOutbox/", emailAddress)
	fmt.Println(MSAGetURL)
	//poll the MSA for new emails ever x seconds
	if req, err := http.NewRequest("DELETE", MSAGetURL, nil); err == nil {
				
		client := &http.Client {}
		//Get the response
		if resp, err1 := client.Do(req); err1 == nil {

			if body, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
				fmt.Println(string(body))
				return body, true
				//Find the URL of the destination address
				
			}
		} else {
			//DELETE request failed
			fmt.Printf("DELETE failed with %s\n", err1)
		}

	} else {
		//DELETE request failed
		fmt.Printf("DELETE failed with %s\n", err)
	}
	return nil, false
}


func HandleRequests(mta *MTA) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/postEmail", mta.PostEmail).Methods("POST")
	log.Fatal(http.ListenAndServe(":8887", router));
}


func main() {
	mta := MTA{make([]*util.Email, 0), "http://localhost:8887/findURL"}
	HandleRequests(&mta)
}


