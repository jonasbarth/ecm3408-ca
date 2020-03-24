#!/bin/bash

echo Starting bluebook server
go run ./bluebook/bluebook.go &

echo Starting msa
go run ./mailsubmissionagent/msa.go &

echo Starting mta
go run ./mailtransferagent/mta.go

