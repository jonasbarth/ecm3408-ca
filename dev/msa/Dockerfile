FROM golang:latest
MAINTAINER 123446
RUN mkdir /app/
ADD . /app/

WORKDIR /app
RUN go get github.com/gorilla/mux
RUN go get github.com/google/uuid
EXPOSE 7001
EXPOSE 8001
CMD [ "go", "run", "/app/msa/src/msa.go" ]