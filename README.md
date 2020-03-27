# ECM3408 CA - Email Service Microservice Implementation

This CA builds a dockerised basic email service using a microservice architecture.



## Architecture

This application mirrors the email servers in the provided CA spec. It consists out of two email servers, each with one Mail Transfer Agent and one Mail Submission Agent. 
Additionally there is a single Bluebook server. 


Email Server A:
- domain name: here.com
- address: localhost:4001

Email Server B:
- domain name: there.com
- address: localhost:5001



## Setup

Run the provided docker-compose.yaml file.

```bash
docker-compose build
docker-compose up
```

The application is first being built and then started.

## Interacting with the service

All microservices are now running. To add an email to every of the 4 user's outbox run the provided shell script.

```bash
sh start.sh
```

