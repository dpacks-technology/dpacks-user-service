# DPacks User Platform Microservice
This is a microservice that provides the user platform for the DPacks project. It is responsible for managing users, their roles, and their permissions.

## Development Branch: Dev

## Installation (Server - Docker)
1. Clone the repository
2. Run `docker build -t dpacks-user-platform .`
3. Run `docker run -p 4001:4001 dpacks-user-platform`
4. The service should now be running on `localhost:4001`

## Installation (Local)
1. Clone the repository
2. Run `go mod download`
3. Run `go run main.go`
4. The service should now be running on `localhost:4001`

## Technologies
- Go, Gin
- Docker
- PostgreSQL

## Copyright
© 2024 DPacks. All Rights Reserved.
