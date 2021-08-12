# PW_Less
## A mad simple passwordless login service written in go

### About
This is a service that facilitates a passwordless login flow. This service would sit in a backend and be something that a front end would reach out to with an email address to generate a "magic link" that a user would click to complete the login. 

The service caches a given email with a randomly generated token and sends an email to that address with a link consisting of the email and token in the query parameters (e.g. `https://redirect.url?email=test@emai.com&token=abc123`). The service handles a subsequent GET request to check if the user has initiated the login flow (exists in the cache) and that the token is valid before returning a 200 with the user data or an error back to a frontend.

This service users [OneLogin](https://www.onelogin.com) as the user store via the [OneLogin Go SDK](https://github.com/onelogin/onelogin-go-sdk) which leans on the [Users API](https://developers.onelogin.com/api-docs/2/users/list-users). You are of course welcome to bring your own storage (e.g. postgres or mongodb).

A more in-depth description can be found in the [accompanying article](https://d-caponi1.medium.com/writing-a-passwordless-service-with-go-and-docker-594d693689dc).

### Endpoints
`POST /users {email: "email@test.com"}`

`GET /users?email=test@emai.com&token=abc123`

### Getting Started
1. Clone this repository
2. `docker-compose up`

### Flow Diagram
![passwordless login flow](https://miro.medium.com/max/700/1*ye-0kb3JFb9bYxzfW6RpdQ.png "passwordless login flow")