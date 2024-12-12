# Go Dentist API

A simple REST API to manage dentist information, built with **Go** and using **AWS DynamoDB** as the database.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Running the Project](#running-the-project)
- [Endpoints](#endpoints)
- [Contributing](#contributing)
- [License](#license)

## Overview

This project is a REST API that allows:
- Creating dentist records.
- Retrieving a list of all dentists.
- Fetching details of a specific dentist by ID.

The API uses **DynamoDB** for data storage and follows best practices for backend development in **Go**.

## Features

- Create dentists (`POST /dentist`).
- List all dentists (`GET /dentists`).
- Retrieve a dentist by ID (`GET /dentist/{id}`).

## Project Structure

```plaintext
go-dentist-api/
│
├── cmd/
│   └── main.go                # Application entry point
│
├── internal/
│   ├── config/                # DynamoDB configuration
│   │   └── dynamodb.go
│   ├── handlers/              # Endpoint handlers
│   │   └── dentist.go
│   ├── models/                # Data models and validations
│   │   └── dentist.go
│   └── router/                # Route configuration
│       └── router.go
├── test/
│
├── go.mod                     # Project dependencies
├── go.sum                     # Dependency version details
└── README.md                  # Project documentation
```

## Prerequisites

Before setting up the project, ensure you have the following tools installed:

- **Go 1.19+**: [Download Go](https://golang.org/dl/)
- **AWS CLI**: [Install AWS CLI](https://aws.amazon.com/cli/)
- **Docker**: [Install Docker](https://www.docker.com/) (to run DynamoDB locally)

## Setup

Follow these steps to set up the project:

### 1. Clone the repository

Run the following commands to clone the project and navigate into the project directory:

```bash
git clone https://github.com/your-username/go-dentist-api.git
cd go-dentist-api
```

### 2. Install dependencies
Run the following command to install project dependencies:

```bash
go mod tidy
```

3. Configure DynamoDB Local
Start DynamoDB locally using Docker:

docker run -d -p 8000:8000 amazon/dynamodb-local
4. Verify DynamoDB
Ensure DynamoDB is running correctly:

aws dynamodb list-tables --endpoint-url http://localhost:8000
Running the Project

Start the server using:

go run cmd/main.go
The server will be available at: http://localhost:8080

Endpoints

Create a Dentist
POST /dentist
Request Body:
{
  "ID": "1",
  "Name": "Dr. John Smith",
  "Email": "john.smith@example.com",
  "Phone": "123456789",
  "CRO": "12345",
  "Country": "USA",
  "Specialty": "Orthodontics",
  "CreatedAt": "2024-11-27T12:00:00Z",
  "UpdatedAt": "2024-11-27T12:00:00Z"
}
Response:
{
  "ID": "1",
  "Name": "Dr. John Smith",
  "Email": "john.smith@example.com",
  "Phone": "123456789",
  "CRO": "12345",
  "Country": "USA",
  "Specialty": "Orthodontics",
  "CreatedAt": "2024-11-27T12:00:00Z",
  "UpdatedAt": "2024-11-27T12:00:00Z"
}
List All Dentists
GET /dentists
Response:
[
  {
    "ID": "1",
    "Name": "Dr. John Smith",
    "Email": "john.smith@example.com",
    "Phone": "123456789",
    "CRO": "12345",
    "Country": "USA",
    "Specialty": "Orthodontics",
    "CreatedAt": "2024-11-27T12:00:00Z",
    "UpdatedAt": "2024-11-27T12:00:00Z"
  }
]
Retrieve Dentist by ID
GET /dentist/{id}
Response:
{
  "ID": "1",
  "Name": "Dr. John Smith",
  "Email": "john.smith@example.com",
  "Phone": "123456789",
  "CRO": "12345",
  "Country": "USA",
  "Specialty": "Orthodontics",
  "CreatedAt": "2024-11-27T12:00:00Z",
  "UpdatedAt": "2024-11-27T12:00:00Z"
}
Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

License

This project is licensed under the MIT License. See the LICENSE file for more details.


### Steps to Finalize:
1. Replace `https://github.com/your-username/go-dentist-api.git` with the correct repository URL.
2. Add additional endpoints or features to the documentation if applicable.
3. If using a license, include the corresponding `LICENSE` file in the project.

Let me know if you'd like further adjustments! 😊