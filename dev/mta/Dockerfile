FROM golang:latest
MAINTAINER 123446
RUN mkdir /app/
ADD . /app/

WORKDIR /app
RUN go get github.com/gorilla/mux
RUN go get github.com/google/uuid
EXPOSE 7000
EXPOSE 8000
CMD [ "go", "run", "/app/mta/src/mta.go" ]