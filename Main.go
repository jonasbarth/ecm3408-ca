package main

import (
	//"MSA"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	//"fmt"
)


func handleRequests(msa MSA) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/sendEmail/{user}", msa.AddEmailToOutbox).Methods("POST")
	router.HandleFunc("/sendEmail/{user}", msa.ListOutbox).Methods("GET")
	log.Fatal(http.ListenAndServe(":8888", router));
}




func main() {

	user := User{make([]Email, 0), make([]Email, 0), "fred@here.com"}
	msa := MSA{make(map[string]*User), "here.com"}
	msa.CreateUser(&user)
	handleRequests(msa);
}