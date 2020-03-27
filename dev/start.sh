#!/bin/bash

curl -v -X POST -d "{\"source\": \"fred@here.com\",\"destination\" : \"bob@there.com\", \"body\": \"This email is from fred to bob\"}" localhost:4001/outbox
curl -v -X POST -d "{\"source\": \"wilma@here.com\",\"destination\" : \"betty@there.com\", \"body\": \"This email is from wilma to betty\"}" localhost:4001/outbox



curl -v -X POST -d "{\"source\": \"bob@there.com\",\"destination\" : \"fred@here.com\", \"body\": \"This email is from bob to fred\"}" localhost:5001/outbox
curl -v -X POST -d "{\"source\": \"betty@there.com\",\"destination\" : \"wilma@here.com\", \"body\": \"This email is from betty to wilma\"}" localhost:5001/outbox
