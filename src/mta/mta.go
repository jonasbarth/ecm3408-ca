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
	"time"
	"os"
)

type MTA struct {
	OutEmails []*util.Email
	BlueBookURL string
	URL string
	MSAURL string
}

//Posts the email into a users inbox via the MSA microservice
//MTA and MSA have the same network address
func (mta *MTA) PostEmail(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Received post request")

	//get body of the eamil
	decoder := json.NewDecoder(r.Body)
	var email util.Email;

	//make sure the email follows the correct format
	if err := decoder.Decode(&email); err == nil {
		
		//Marshal the email object into JSON
		if enc, err := json.Marshal(email); err == nil {
			
			//Create the URL 
			MSAPostURL := mta.MSAURL + "/toInbox/" + email.Destination

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
//Then sends the email to the MTA of the destination address
func (mta *MTA) PollMSA(emailAddress string) {


	for ok := true; ok; ok = true {

		time.Sleep(5000 * time.Millisecond)

		//Get the email from outbox of the MSA to which the email belongs
		if email, err := mta.peekOutbox(emailAddress); err == nil {

			fmt.Println("Email successfully retrieved from")
			//Delete the email from the outbox 
			if ok, err0 := mta.deleteOutbox(email); ok  {
				
				//Marshal the email and send it to the correct MTA
				if enc, err1 := json.Marshal(&email); err1 == nil {
					
					//Get the network address of the correct MTA service
					MTAPostURL := mta.getURL(email.Destination) + "/postEmail/" + email.Destination

					fmt.Println("Posting email to MTA at " + MTAPostURL)
					//Create and make the POST request
					if req, err2 := http.NewRequest("POST", MTAPostURL, bytes.NewBuffer(enc)); err2 == nil {
						
						client := &http.Client {}

						//Get the response and implement error handling
						if resp, err3 := client.Do(req); err3 == nil {
							fmt.Printf("Reponse %s\n", resp.Status)
							break
							
						} else {
							//POST request failed
							fmt.Printf("POST request to %s failed with error %s\n", email.Destination, err3)
							break
							
						}

					} else {
						//POST request failed
						fmt.Printf("POST request to %s failed with error %s\n", email.Destination, err2)
						break
						
					}
				} else {
					//Marshalling failed
					fmt.Printf("Cannot marshal JSON with error %s\n", err1)
					break
				}
			} else {
				fmt.Printf("Deleting email from outbox failed with error %s\n", err0)
			}

				
				
			
		} else {mta.getURL(email.Destination)
			fmt.Println("Failed to pop from outbox")
			break
		}
	}

		

}



//Gets the URL of the email server belonging to the email
//Returns the URL as a string if found on the bluebook server
//Returns an empty string if the URL does not exist on the bluebook server
func (mta *MTA) getURL(email string) string {

	url := mta.BlueBookURL + "/" + email
	
	if req, err1 := http.NewRequest("GET", url, nil); err1 == nil {
				
		client := &http.Client {}
		//Get the response
		if resp, err2 := client.Do(req); err2 == nil {

			if body, err3 := ioutil.ReadAll(resp.Body); err3 == nil {
				fmt.Println(email + " can be found at " + string(body))
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
					return &email, nil
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


func (mta *MTA) deleteOutbox(email *util.Email) (bool, error) {


	MSADeleteURL := mta.MSAURL + "/deleteOutbox/" + email.Source + "/" + email.UUID
	

	fmt.Println(MSADeleteURL)

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


func (mta *MTA) HandleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/postEmail/{user}", mta.PostEmail).Methods("POST")
	
	fmt.Println("MTA Service is running at " + mta.URL)
	log.Fatal(http.ListenAndServe(mta.URL, router));
}


func main() {

	jsonFile, err := os.Open("../resources/init.json")
    // if we os.Open returns an error then handle it
    if err != nil {
        fmt.Println(err)
	}
	
	defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	
	fmt.Println(result["mta"])


	bluebookURL := "http://localhost:9000/findURL"
	msaAddress := "http://localhost:7001"
	mta := MTA{make([]*util.Email, 0), bluebookURL, ":7000", msaAddress}
	go mta.HandleRequests()

	mta.getURL("fred@here.com")
	mta.getURL("fred@there.com")

	go mta.PollMSA("fred@here.com")

	msa2Address := "http://localhost:8001"
	mta2 := MTA{make([]*util.Email, 0), bluebookURL, ":8000", msa2Address}
	mta2.HandleRequests()
}


