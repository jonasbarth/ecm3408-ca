package main

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"fmt"
	"../../util"
	"github.com/gorilla/mux"
	"log"
	"time"
	"math/rand"
	"errors"
)

type MTA struct {
	OutEmails []*util.Email
	BlueBookURL string
	URL string
	MSAURL string
}

//Posts the email into a users inbox via the MSA microservice
//MTA and MSA have the same network address
func (mta *MTA) SendEmail(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Received post request")

	//get body of the eamil
	decoder := json.NewDecoder(r.Body)
	var email util.Email;

	//make sure the email follows the correct format
	if err := decoder.Decode(&email); err == nil {
		
		//Marshal the email object into JSON
		if enc, err := json.Marshal(email); err == nil {
			
			//Create the URL 
			MSAPostURL := mta.MSAURL + "/inbox"

			//Create and make the POST request
			if req, err1 := http.NewRequest("POST", MSAPostURL, bytes.NewBuffer(enc)); err1 == nil {
				
				client := &http.Client {}
				//Get the response
				if resp, err2 := client.Do(req); err2 == nil {

					w.WriteHeader(resp.StatusCode)
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
//Then sends the email to the MTA of the destination address
func (mta *MTA) PollMSA() {


	for ok := true; ok; ok = true {
		
		//MSA is polled every 3-4 seconds to avoid all MTAs using the network at the same time
		min := 3000
		max := 4000
		time.Sleep(mta.getRandom(min, max) * time.Millisecond)

		//get all users from the MSA
		users := mta.getUsers()

		for _, emailAddress := range users {		

			//Get the email from outbox of the MSA to which the email belongs. Also make sure we actually get an email before making further requests
			if email, err := mta.peekOutbox(emailAddress); err == nil {
				fmt.Printf("Peeking outbox %s\n", email.Destination)
				fmt.Println(email)
				//Delete the email from the outbox 
				if ok, err1 := mta.deleteOutbox(email); ok  {
					
					//Marshal the email and send it to the correct MTA
					if enc, err2 := json.Marshal(&email); err2 == nil {
						
						//Get the network address of the correct MTA service
						if MTAPostURL, err3 := mta.getURL(email.Destination); err3 == nil {

							MTAPostURL = MTAPostURL + "/send/" + email.Destination
							
							//Create and make the POST request
							if req, err4 := http.NewRequest("POST", MTAPostURL, bytes.NewBuffer(enc)); err4 == nil {
								
								client := &http.Client {}

								//Get the response and implement error handling
								if resp, err5 := client.Do(req); err5 == nil {
									fmt.Printf("Response Code %s\n", resp.Status)
									
								} else {
									//POST request failed
									fmt.Printf("POST request to %s failed with error %s\n", email.Destination, err5)				
								}

							} else {
								//POST request failed
								fmt.Printf("POST request to %s failed with error %s\n", email.Destination, err4)
									
							}
						} else {
							//Could not get the URL for the specified email addresss
							fmt.Printf("Bluebook server could not find network address of %s with error %s\n", email.Destination, err3)
						} 
					
						
					} else {
						//Marshalling failed
						fmt.Printf("Cannot marshal JSON with error %s\n", err2)
						
					}
				} else {
					//Could not delete the email from the user's outbox
					fmt.Printf("Deleting email from outbox failed with error %s\n", err1)
				}

							
			} else {
				//Could not get the email from the user's outbox
				fmt.Printf("Failed to get email from outbox with error %s\n", err)
				
			}

		}
		
	}
}



//Gets the URL of the email server belonging to the email
//Returns the URL as a string if found on the bluebook server
//Returns an empty string if the URL does not exist on the bluebook server
func (mta *MTA) getURL(email string) (string, error) {

	url := mta.BlueBookURL + "/" + email
	
	if req, err1 := http.NewRequest("GET", url, nil); err1 == nil {
				
		client := &http.Client {}
		//Get the response
		if resp, err2 := client.Do(req); err2 == nil {

			if body, err3 := ioutil.ReadAll(resp.Body); err3 == nil {
				return string(body)[1:len(string(body))-1], nil
			} else {
				fmt.Printf("GET request to %s failed with %s\n", url, err3)
				return "", err3
			}
		} else {
			fmt.Printf("GET request to %s failed with %s\n", url, err2)
			//GET request failed
			return "", err2
		}

	} else {
		fmt.Printf("GET request to %s failed with %s\n", url, err1)
		//GET request failed
		return "", err1
	}
}



//Pops the outbox of the user as specified by the email address
//Returns the email and true if an email exists in the user's outbox
//Returns nil and false if the outbox is empty
func (mta *MTA) peekOutbox(emailAddress string) (*util.Email, error) {


	MSAGetURL := mta.MSAURL + "/peekOutbox/" + emailAddress

	var email util.Email
	

	if req, err := http.NewRequest("GET", MSAGetURL, nil); err == nil {
				
		client := &http.Client {}
		//Get the response
		if resp, err1 := client.Do(req); err1 == nil {

			if body, err2 := ioutil.ReadAll(resp.Body); err2 == nil {

				if err3 := json.Unmarshal(body, &email); err3 == nil {
					
					if email.Source != "" {
						return &email, nil
					} else {
						return nil, errors.New("Outbox of user " + emailAddress + " is empty")
					}
					
				} else {
					fmt.Printf("Could not unmarshal email with error %s\n", err3)
					return nil, err3
				}
				
				
			} else {
				fmt.Printf("GET failed with %s\n", err2)
				return nil, err2
			}
		} else {
			
			//DELETE request failed
			fmt.Printf("DELETE failed with %s\n", err1)
			return nil, err1
		}

	} else {
		
		//DELETE request failed
		fmt.Printf("DELETE failed with %s\n", err)
		return nil, err
	}

}

//Deletes an email from this email server
//Returns true, nil if delete successful
//Returns false, error if unsuccessful
func (mta *MTA) deleteOutbox(email *util.Email) (bool, error) {

	MSADeleteURL := mta.MSAURL + "/outbox/" + email.Source + "/" + email.UUID

	//poll the MSA for new emails ever x seconds
	if req, err := http.NewRequest("DELETE", MSADeleteURL, nil); err == nil {
				
		client := &http.Client {}
		//Get the response
		if resp, err1 := client.Do(req); err1 == nil {

			if _, err2 := ioutil.ReadAll(resp.Body); err2 == nil {
			
				return true, nil
			} else {
				fmt.Printf("DELETE failed with %s\n", err2)
				return false, err2
			}
		} else {
			//DELETE request failed
			fmt.Printf("DELETE failed with %s\n", err1)
			return false, err1
		}

	} else {
		//DELETE request failed
		fmt.Printf("DELETE failed with %s\n", err)
		return false, err
	}

}

//Gets a list of users from the MSA of this email server
//The list is used to poll the users' inboxes
//Returns an empty list if no users exist on the server
func (mta *MTA) getUsers() []string {

	MSAGetURL := mta.MSAURL + "/users"
	var users []string

	if req, err := http.NewRequest("GET", MSAGetURL, nil); err == nil {
				
		client := &http.Client {}

		//Get the response
		if resp, err1 := client.Do(req); err1 == nil {
	

			if body, err2 := ioutil.ReadAll(resp.Body); err2 == nil {

				if err3 := json.Unmarshal(body, &users); err3 == nil {
					return users

				} else {
					fmt.Printf("Could not unmarshal JSON with error %s\n", err3)
				}
			} else {
				fmt.Printf("Could not read response body with error %\n", err2)
			}
		} else {
			fmt.Printf("GET request failed with error %s\n", err1)
		}
	} else {
		fmt.Printf("GET request failed with error %s\n", err)
	}
	return users

}

//Produces a random integer in the range of [min, max)
//Used for the outbox polling intervals
func (mta *MTA) getRandom(min int, max int) time.Duration {
	rand.Seed(time.Now().UnixNano())
 	return time.Duration(rand.Intn(max - min + 1) + min)
}


func (mta *MTA) HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/send/{user}", mta.SendEmail).Methods("POST")
	
	fmt.Println("MTA Service is running at " + mta.URL)
	log.Fatal(http.ListenAndServe(mta.URL, router));
}


func main() {

	bluebookURL := "http://bluebook:9000/address"
	msaAddress := "http://msa:7001"
	mta := MTA{make([]*util.Email, 0), bluebookURL, ":7000", msaAddress}
	go mta.HandleRequests()

	go mta.PollMSA()

	msa2Address := "http://msa:8001"
	mta2 := MTA{make([]*util.Email, 0), bluebookURL, ":8000", msa2Address}

	go mta2.PollMSA()

	mta2.HandleRequests()
}


