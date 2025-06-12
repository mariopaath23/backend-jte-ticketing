# Go JWT-Auth Backend for JTE-Ticketing

This is the backend server for the JTE-Ticketing application, built with Go, MongoDB, and using JWT for authentication.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup & Installation](#setup--installation)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)

## Features

- User Registration
- User Login with JWT-based authentication
- Protected routes using JWT middleware
- MongoDB integration

## Project Structure


/backend
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   │   └── jwt.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   └── user_handler.go
│   ├── middleware/
│   │   └── auth.go
│   └── models/
│       └── user.go
├── go.mod
├── go.sum
└── .env


## Prerequisites

- [Go](https://go.dev/doc/install) (version 1.18 or newer)
- [MongoDB](https://www.mongodb.com/try/download/community) installed and running.
- A tool to make API requests, like [Postman](https://www.postman.com/downloads/) or `curl`.

## Setup & Installation

1.  **Clone the repository and navigate to the `backend` directory.**

2.  **Create the module:**
    Initialize the Go module. Replace `your-username/your-repo-name` with your actual GitHub username and repository name.
    ```bash
    go mod init [github.com/your-username/your-repo-name/backend](https://github.com/your-username/your-repo-name/backend)
    ```

3.  **Install dependencies:**
    This command will download all the necessary libraries defined in `go.mod`.
    ```bash
    go get .
    ```

4.  **Set up environment variables:**
    Create a `.env` file in the root of the `/backend` directory and add the following variables. Replace the placeholder values with your actual configuration.
    ```env
    MONGO_URI=mongodb://localhost:27017
    MONGO_DATABASE=jte_ticketing
    JWT_SECRET_KEY=your_super_secret_key
    API_PORT=8080
    ```

## Running the Application

To start the server, run the following command from the `/backend` directory:

```bash
go run cmd/server/main.go

You should see a message in your console indicating that the server is running:
Server starting on port 8080...
API Endpoints

All endpoints are prefixed with /api.

Method
	

Endpoint
	

Description
	

Authentication

POST
	

/api/register
	

Register a new user.
	

None

POST
	

/api/login
	

Log in an existing user.
	

None

GET
	

/api/protected
	

Example protected route.
	

JWT Token

Request Body for /api/register and /api/login:

{
    "email": "test@example.com",
    "password": "password123"
}

