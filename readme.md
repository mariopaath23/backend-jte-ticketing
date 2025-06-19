# Go JWT-Auth Backend for JTE-Ticketing

This is the backend server for the JTE-Ticketing application, built with Go, MongoDB, and using JWT for authentication.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup & Installation](#setup--installation)
- [Database Migration & Seeding](#database-migration--seeding)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)

## Features

- User Registration with default 'student' role
- User Login with JWT-based authentication
- JWT tokens automatically expire after 120 minutes
- User Logout (clears authentication cookie)
- Login session logging (timestamp, user agent)
- Protected routes using JWT middleware
- MongoDB integration with migrations and seeding

## API Endpoints

All endpoints are prefixed with `/api`.

| Method | Endpoint          | Description                       | Authentication |
| :----- | :---------------- | :-------------------------------- | :------------- |
| `POST` | `/api/register`   | Register a new user.              | None           |
| `POST` | `/api/login`      | Log in an existing user.          | None           |
| `POST` | `/api/logout`     | Log out the current user.         | JWT Token      |
| `GET`  | `/api/protected`  | Example protected route.          | JWT Token      |
| `GET`  | `/api/login-logs` | Get login history for the user.   | JWT Token      |
