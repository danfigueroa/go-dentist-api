# Dental Clinic API

A comprehensive REST API to manage a dental clinic, including dentists, patients, procedures, and appointments. Built with **Go** and using **AWS DynamoDB** as the database.

## Table of Contents

-   [Overview](#overview)
-   [Features](#features)
-   [Project Structure](#project-structure)
-   [Prerequisites](#prerequisites)
-   [Setup](#setup)
-   [Running the Project](#running-the-project)
-   [Endpoints](#endpoints)
-   [Contributing](#contributing)
-   [License](#license)

## Overview

This project is a REST API that allows managing all aspects of a dental clinic, including:

-   Dentists (create, read, update, delete)
-   Patients (create, read, update, delete)
-   Procedures (create, read, update, delete)
-   Appointments (create, read, update, delete)

The API uses **DynamoDB** for data storage and follows best practices for backend development in **Go**.

## Features

### Dentists

-   Create dentist (`POST /dentist`)
-   List all dentists (`GET /dentist`)
-   Retrieve a dentist by ID (`GET /dentist/{id}`)
-   Search dentists by name (`GET /dentist/name/{name}`)
-   Search dentist by CRO (`GET /dentist/cro/{cro}`)
-   Update a dentist (`PUT /dentist/{id}`)
-   Delete a dentist (`DELETE /dentist/{id}`)

### Patients

-   Create patient (`POST /patient`)
-   List all patients (`GET /patient`)
-   Retrieve a patient by ID (`GET /patient/{id}`)
-   Search patients by name (`GET /patient/name/{name}`)
-   Update a patient (`PUT /patient/{id}`)
-   Delete a patient (`DELETE /patient/{id}`)

### Procedures

-   Create procedure (`POST /procedure`)
-   List all procedures (`GET /procedure`)
-   Retrieve a procedure by ID (`GET /procedure/{id}`)
-   Search procedures by name (`GET /procedure/name/{name}`)
-   Update a procedure (`PUT /procedure/{id}`)
-   Delete a procedure (`DELETE /procedure/{id}`)

### Appointments

-   Create appointment (`POST /appointment`)
-   List all appointments (`GET /appointment`)
-   Retrieve an appointment by ID (`GET /appointment/{id}`)
-   Update an appointment (`PUT /appointment/{id}`)
-   Delete an appointment (`DELETE /appointment/{id}`)

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
│   │   ├── dentist.go
│   │   ├── patient.go
│   │   ├── procedure.go
│   │   └── appointment.go
│   ├── models/                # Data models and validations
│   │   ├── dentist.go
│   │   ├── patient.go
│   │   ├── procedure.go
│   │   └── appointment.go
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

-   **Go 1.19+**: [Download Go](https://golang.org/dl/)
-   **AWS CLI**: [Install AWS CLI](https://aws.amazon.com/cli/)
-   **Docker**: [Install Docker](https://www.docker.com/) (to run DynamoDB locally)

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

### 3. Running with Docker Compose (Recommended)

The easiest way to run the project is using Docker Compose, which will start both the local DynamoDB and the API with a single command:

```bash
docker-compose up -d
```

This will build the API image and start both services. The API will be available at http://localhost:8080 and the local DynamoDB at http://localhost:8000.

To stop the services:

```bash
docker-compose down
```

### 4. Running Manually (Alternative)

If you prefer to run the components separately:

#### 4.1. Configure DynamoDB Local

Start DynamoDB locally using Docker:

```bash
docker run -d -p 8000:8000 amazon/dynamodb-local
```

#### 4.2. Verify DynamoDB

Ensure DynamoDB is running correctly:

```bash
aws dynamodb list-tables --endpoint-url http://localhost:8000
```

#### 4.3. Running the API

Start the server using:

```bash
go run cmd/main.go
```

The server will be available at: http://localhost:8080

## Endpoints

### Dentists

**Create a Dentist**

```
POST /dentist
```

Request Body:

```json
{
    "Name": "Dr. John Smith",
    "Email": "john.smith@example.com",
    "Phone": "123456789",
    "CRO": "12345",
    "Country": "USA",
    "Specialty": "Orthodontics"
}
```

**List All Dentists**

```
GET /dentist
```

Response Example:

```json
[
    {
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "Name": "Dr. John Smith",
        "Email": "john.smith@example.com",
        "Phone": "123456789",
        "CRO": "12345",
        "Country": "USA",
        "Specialty": "Orthodontics",
        "CreatedAt": "2023-01-01T10:00:00Z",
        "UpdatedAt": "2023-01-01T10:00:00Z"
    },
    {
        "ID": "223e4567-e89b-12d3-a456-426614174001",
        "Name": "Dr. Sarah Johnson",
        "Email": "sarah.johnson@example.com",
        "Phone": "987654321",
        "CRO": "54321",
        "Country": "USA",
        "Specialty": "Periodontics",
        "CreatedAt": "2023-01-02T11:00:00Z",
        "UpdatedAt": "2023-01-02T11:00:00Z"
    }
]
```

**Retrieve Dentist by ID**

```
GET /dentist/{id}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "Name": "Dr. John Smith",
    "Email": "john.smith@example.com",
    "Phone": "123456789",
    "CRO": "12345",
    "Country": "USA",
    "Specialty": "Orthodontics",
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-01T10:00:00Z"
}
```

**Search Dentists by Name**

```
GET /dentist/name/{name}
```

Response Example:

```json
[
    {
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "Name": "Dr. John Smith",
        "Email": "john.smith@example.com",
        "Phone": "123456789",
        "CRO": "12345",
        "Country": "USA",
        "Specialty": "Orthodontics",
        "CreatedAt": "2023-01-01T10:00:00Z",
        "UpdatedAt": "2023-01-01T10:00:00Z"
    }
]
```

**Search Dentist by CRO**

```
GET /dentist/cro/{cro}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "Name": "Dr. John Smith",
    "Email": "john.smith@example.com",
    "Phone": "123456789",
    "CRO": "12345",
    "Country": "USA",
    "Specialty": "Orthodontics",
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-01T10:00:00Z"
}
```

**Update Dentist**

```
PUT /dentist/{id}
```

Request Example:

```json
{
    "Name": "Dr. John Smith Jr.",
    "Email": "john.smith.jr@example.com",
    "Phone": "123456789",
    "CRO": "12345",
    "Country": "USA",
    "Specialty": "Orthodontics and Pediatric Dentistry"
}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "Name": "Dr. John Smith Jr.",
    "Email": "john.smith.jr@example.com",
    "Phone": "123456789",
    "CRO": "12345",
    "Country": "USA",
    "Specialty": "Orthodontics and Pediatric Dentistry",
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-03T15:30:00Z"
}
```

**Delete Dentist**

```
DELETE /dentist/{id}
```

Response Example:

```json
{
    "message": "Dentist deleted successfully"
}
```

### Patients

**Create a Patient**

```
POST /patient
```

Request Body:

```json
{
    "Name": "Jane Doe",
    "Email": "jane.doe@example.com",
    "Phone": "987654321",
    "DateOfBirth": "1990-01-01",
    "MedicalNotes": "No allergies"
}
```

**List All Patients**

```
GET /patient
```

Response Example:

```json
[
    {
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "Name": "Jane Doe",
        "Email": "jane.doe@example.com",
        "Phone": "987654321",
        "DateOfBirth": "1990-01-01",
        "MedicalNotes": "No allergies",
        "CreatedAt": "2023-01-01T10:00:00Z",
        "UpdatedAt": "2023-01-01T10:00:00Z"
    },
    {
        "ID": "223e4567-e89b-12d3-a456-426614174001",
        "Name": "John Smith",
        "Email": "john.smith@example.com",
        "Phone": "123456789",
        "DateOfBirth": "1985-05-15",
        "MedicalNotes": "Allergic to penicillin",
        "CreatedAt": "2023-01-02T11:00:00Z",
        "UpdatedAt": "2023-01-02T11:00:00Z"
    }
]
```

**Retrieve Patient by ID**

```
GET /patient/{id}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "Name": "Jane Doe",
    "Email": "jane.doe@example.com",
    "Phone": "987654321",
    "DateOfBirth": "1990-01-01",
    "MedicalNotes": "No allergies",
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-01T10:00:00Z"
}
```

**Search Patients by Name**

```
GET /patient/name/{name}
```

Response Example:

```json
[
    {
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "Name": "Jane Doe",
        "Email": "jane.doe@example.com",
        "Phone": "987654321",
        "DateOfBirth": "1990-01-01",
        "MedicalNotes": "No allergies",
        "CreatedAt": "2023-01-01T10:00:00Z",
        "UpdatedAt": "2023-01-01T10:00:00Z"
    }
]
```

**Update Patient**

```
PUT /patient/{id}
```

Request Example:

```json
{
    "Name": "Jane Doe Smith",
    "Email": "jane.smith@example.com",
    "Phone": "987654321",
    "DateOfBirth": "1990-01-01",
    "MedicalNotes": "No allergies, recent dental surgery"
}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "Name": "Jane Doe Smith",
    "Email": "jane.smith@example.com",
    "Phone": "987654321",
    "DateOfBirth": "1990-01-01",
    "MedicalNotes": "No allergies, recent dental surgery",
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-03T15:30:00Z"
}
```

**Delete Patient**

```
DELETE /patient/{id}
```

Response Example:

```json
{
    "message": "Patient deleted successfully"
}
```

### Procedures

**Create a Procedure**

```
POST /procedure
```

Request Body:

```json
{
    "Name": "Teeth Cleaning",
    "Description": "Standard teeth cleaning procedure",
    "Price": 100.0,
    "Duration": 30
}
```

**List All Procedures**

```
GET /procedure
```

Response Example:

```json
[
    {
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "Name": "Teeth Cleaning",
        "Description": "Standard teeth cleaning procedure",
        "Price": 100.0,
        "Duration": 30,
        "CreatedAt": "2023-01-01T10:00:00Z",
        "UpdatedAt": "2023-01-01T10:00:00Z"
    },
    {
        "ID": "223e4567-e89b-12d3-a456-426614174001",
        "Name": "Root Canal",
        "Description": "Endodontic treatment for infected pulp",
        "Price": 800.0,
        "Duration": 90,
        "CreatedAt": "2023-01-02T11:00:00Z",
        "UpdatedAt": "2023-01-02T11:00:00Z"
    }
]
```

**Retrieve Procedure by ID**

```
GET /procedure/{id}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "Name": "Teeth Cleaning",
    "Description": "Standard teeth cleaning procedure",
    "Price": 100.0,
    "Duration": 30,
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-01T10:00:00Z"
}
```

**Search Procedures by Name**

```
GET /procedure/name/{name}
```

Response Example:

```json
[
    {
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "Name": "Teeth Cleaning",
        "Description": "Standard teeth cleaning procedure",
        "Price": 100.0,
        "Duration": 30,
        "CreatedAt": "2023-01-01T10:00:00Z",
        "UpdatedAt": "2023-01-01T10:00:00Z"
    }
]
```

**Update Procedure**

```
PUT /procedure/{id}
```

Request Example:

```json
{
    "Name": "Advanced Teeth Cleaning",
    "Description": "Deep cleaning procedure with polishing",
    "Price": 120.0,
    "Duration": 45
}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "Name": "Advanced Teeth Cleaning",
    "Description": "Deep cleaning procedure with polishing",
    "Price": 120.0,
    "Duration": 45,
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-03T15:30:00Z"
}
```

**Delete Procedure**

```
DELETE /procedure/{id}
```

Response Example:

```json
{
    "message": "Procedure deleted successfully"
}
```

### Appointments

**Create an Appointment**

```
POST /appointment
```

Request Body:

```json
{
    "DentistID": "1",
    "PatientID": "2",
    "DateTime": "2024-12-01T10:00:00Z",
    "Notes": "Regular checkup"
}
```

**List All Appointments**

```
GET /appointment
```

Response Example:

```json
[
    {
        "ID": "123e4567-e89b-12d3-a456-426614174000",
        "DentistID": "123e4567-e89b-12d3-a456-426614174000",
        "PatientID": "223e4567-e89b-12d3-a456-426614174001",
        "ProcedureID": "323e4567-e89b-12d3-a456-426614174002",
        "DateTime": "2024-12-01T10:00:00Z",
        "Notes": "Regular checkup",
        "Status": "Scheduled",
        "CreatedAt": "2023-01-01T10:00:00Z",
        "UpdatedAt": "2023-01-01T10:00:00Z"
    },
    {
        "ID": "223e4567-e89b-12d3-a456-426614174001",
        "DentistID": "123e4567-e89b-12d3-a456-426614174000",
        "PatientID": "223e4567-e89b-12d3-a456-426614174001",
        "ProcedureID": "423e4567-e89b-12d3-a456-426614174003",
        "DateTime": "2024-12-02T14:00:00Z",
        "Notes": "Follow-up treatment",
        "Status": "Scheduled",
        "CreatedAt": "2023-01-02T11:00:00Z",
        "UpdatedAt": "2023-01-02T11:00:00Z"
    }
]
```

**Retrieve Appointment by ID**

```
GET /appointment/{id}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "DentistID": "123e4567-e89b-12d3-a456-426614174000",
    "PatientID": "223e4567-e89b-12d3-a456-426614174001",
    "ProcedureID": "323e4567-e89b-12d3-a456-426614174002",
    "DateTime": "2024-12-01T10:00:00Z",
    "Notes": "Regular checkup",
    "Status": "Scheduled",
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-01T10:00:00Z"
}
```

**Update Appointment**

```
PUT /appointment/{id}
```

Request Example:

```json
{
    "DentistID": "123e4567-e89b-12d3-a456-426614174000",
    "PatientID": "223e4567-e89b-12d3-a456-426614174001",
    "ProcedureID": "323e4567-e89b-12d3-a456-426614174002",
    "DateTime": "2024-12-05T11:30:00Z",
    "Notes": "Rescheduled regular checkup",
    "Status": "Rescheduled"
}
```

Response Example:

```json
{
    "ID": "123e4567-e89b-12d3-a456-426614174000",
    "DentistID": "123e4567-e89b-12d3-a456-426614174000",
    "PatientID": "223e4567-e89b-12d3-a456-426614174001",
    "ProcedureID": "323e4567-e89b-12d3-a456-426614174002",
    "DateTime": "2024-12-05T11:30:00Z",
    "Notes": "Rescheduled regular checkup",
    "Status": "Rescheduled",
    "CreatedAt": "2023-01-01T10:00:00Z",
    "UpdatedAt": "2023-01-03T15:30:00Z"
}
```

**Delete Appointment**

```
DELETE /appointment/{id}
```

Response Example:

```json
{
    "message": "Appointment deleted successfully"
}
```

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
