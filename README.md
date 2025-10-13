# Dental SaaS API

A comprehensive SaaS platform for dental clinic management built with **Go** and **DynamoDB**, following **Hexagonal Architecture** principles.

## 🏗️ Architecture Overview

This project implements a **Hexagonal Architecture** (Ports and Adapters) pattern, ensuring clean separation of concerns and high testability:

```
internal/
├── core/                   # Business Logic Layer
│   ├── domain/            # Domain entities and business rules
│   └── ports/             # Interfaces (repositories & services)
│       ├── repositories/  # Data access interfaces
│       └── services/      # Business service interfaces
├── application/           # Application Layer
│   └── usecases/         # Use case implementations
└── adapters/             # Infrastructure Layer
    ├── primary/          # Input adapters (HTTP handlers)
    │   └── http/
    │       ├── handlers/ # HTTP request handlers
    │       └── routes/   # Route definitions
    └── secondary/        # Output adapters (databases, external APIs)
        └── database/     # Database implementations
```

## 🚀 Technologies Used

- **Go 1.22**: Modern, efficient programming language
- **Gorilla Mux**: HTTP router and URL matcher
- **AWS DynamoDB**: NoSQL database for scalable data storage
- **DynamoDB Local**: Local development database
- **Swagger/OpenAPI**: API documentation and testing
- **Docker & Docker Compose**: Containerization and orchestration
- **UUID**: Unique identifier generation
- **Hexagonal Architecture**: Clean architecture pattern

## 📦 Features

### Currently Implemented
- ✅ **Dentist Management**: Complete CRUD operations for dentists
- ✅ **Domain Models**: Patient, Appointment, Procedure, Clinic entities
- ✅ **Use Cases**: Business logic implementation
- ✅ **Mock Repositories**: In-memory data storage for development
- ✅ **API Documentation**: Swagger/OpenAPI integration
- ✅ **Docker Support**: Full containerization

### In Development
- 🔄 **HTTP Routes**: Patient, Appointment, Procedure, Clinic endpoints
- 🔄 **Database Integration**: DynamoDB implementation
- 🔄 **Authentication**: JWT-based authentication system

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.22+ (for local development)

### Running with Docker

1. **Clone the repository:**
```bash
git clone <repository-url>
cd go-dentist-api
```

2. **Start the application:**
```bash
docker-compose up --build
```

3. **Access the services:**
- **API**: http://localhost:8080
- **Swagger Documentation**: http://localhost:8080/swagger/
- **DynamoDB Local**: http://localhost:8000

### Local Development

1. **Start DynamoDB Local:**
```bash
docker-compose up dynamodb-local -d
```

2. **Run the application:**
```bash
go run cmd/main.go
```

3. **Build the application:**
```bash
go build -o main cmd/main.go
./main
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Available Endpoints

#### Dentists (Fully Implemented)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/dentists` | List all dentists |
| `POST` | `/dentists` | Create a new dentist |
| `GET` | `/dentists/{id}` | Get dentist by ID |
| `PUT` | `/dentists/{id}` | Update dentist |
| `DELETE` | `/dentists/{id}` | Delete dentist |
| `GET` | `/dentists/search?cro={cro}` | Search dentist by CRO |
| `GET` | `/dentists/search?name={name}` | Search dentists by name |
| `GET` | `/dentists/clinic/{clinicId}` | Get dentists by clinic |

#### Query Parameters for Dentists
- `clinic_id`: Filter by clinic ID
- `active`: Filter by active status (true/false)
- `specialty`: Filter by specialty

#### Example Requests

**Create Dentist:**
```bash
curl -X POST "http://localhost:8080/api/v1/dentists" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dr. John Smith",
    "cro": "CRO-SP-12345",
    "email": "john.smith@clinic.com",
    "phone": "+1 555-0123",
    "specialty": "Orthodontics",
    "clinic_id": "clinic-123"
  }'
```

**List Dentists:**
```bash
curl -X GET "http://localhost:8080/api/v1/dentists"
```

**Filter by Clinic:**
```bash
curl -X GET "http://localhost:8080/api/v1/dentists?clinic_id=clinic-123"
```

#### Response Format

```json
{
  "id": "uuid-generated-id",
  "clinic_id": "clinic-123",
  "name": "Dr. John Smith",
  "cro": "CRO-SP-12345",
  "specialty": "Orthodontics",
  "phone": "+1 555-0123",
  "email": "john.smith@clinic.com",
  "active": true,
  "created_at": "2025-01-13T22:03:13.421363Z",
  "updated_at": "2025-01-13T22:03:13.421363Z"
}
```

### Planned Endpoints (Domain Models Ready)

#### Patients
- `GET /patients` - List patients
- `POST /patients` - Create patient
- `GET /patients/{id}` - Get patient by ID
- `PUT /patients/{id}` - Update patient
- `DELETE /patients/{id}` - Delete patient

#### Appointments
- `GET /appointments` - List appointments
- `POST /appointments` - Create appointment
- `GET /appointments/{id}` - Get appointment by ID
- `PUT /appointments/{id}` - Update appointment
- `DELETE /appointments/{id}` - Delete appointment

#### Procedures
- `GET /procedures` - List procedures
- `POST /procedures` - Create procedure
- `GET /procedures/{id}` - Get procedure by ID
- `PUT /procedures/{id}` - Update procedure
- `DELETE /procedures/{id}` - Delete procedure

#### Clinics
- `GET /clinics` - List clinics
- `POST /clinics` - Create clinic
- `GET /clinics/{id}` - Get clinic by ID
- `PUT /clinics/{id}` - Update clinic
- `DELETE /clinics/{id}` - Delete clinic

## 🏗️ Project Structure

```
go-dentist-api/
├── cmd/
│   └── main.go                    # Application entry point
├── internal/
│   ├── core/
│   │   ├── domain/               # Domain entities
│   │   │   ├── appointment.go
│   │   │   ├── clinic.go
│   │   │   ├── dentist.go
│   │   │   ├── patient.go
│   │   │   └── procedure.go
│   │   └── ports/
│   │       ├── repositories/     # Repository interfaces
│   │       └── services/         # Service interfaces
│   ├── application/
│   │   └── usecases/            # Business logic implementation
│   └── adapters/
│       ├── primary/
│       │   └── http/
│       │       ├── handlers/    # HTTP handlers
│       │       └── routes/      # Route definitions
│       └── secondary/
│           └── database/        # Database implementations
├── docs/                        # Swagger documentation
├── docker-compose.yml          # Docker configuration
├── Dockerfile                  # Container definition
└── README.md                   # This file
```

## 🧪 Testing

### Manual API Testing

Use the provided test script:
```bash
chmod +x test_api.sh
./test_api.sh
```

### Using cURL

Test the dentist endpoints:
```bash
# List all dentists
curl -X GET "http://localhost:8080/api/v1/dentists"

# Create a dentist
curl -X POST "http://localhost:8080/api/v1/dentists" \
  -H "Content-Type: application/json" \
  -d '{"name": "Dr. Test", "cro": "CRO-TEST-001", "email": "test@clinic.com", "phone": "+1 555-0001", "specialty": "General", "clinic_id": "clinic-123"}'
```

### Using Swagger UI

Visit http://localhost:8080/swagger/ to interact with the API through the Swagger interface.

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DYNAMODB_ENDPOINT` | DynamoDB endpoint URL | `http://localhost:8000` |
| `PORT` | Application port | `8080` |

### DynamoDB Tables

The application uses the following DynamoDB tables:
- `Dentists` - Dentist information
- `Patients` - Patient records
- `Appointments` - Appointment scheduling
- `Procedures` - Medical procedures
- `Clinics` - Clinic information

## 🚧 Development Roadmap

### Phase 1: Core API (Current)
- ✅ Hexagonal architecture implementation
- ✅ Domain models and use cases
- ✅ Dentist endpoints (complete)
- 🔄 Patient, Appointment, Procedure, Clinic endpoints

### Phase 2: Database Integration
- 🔄 DynamoDB repository implementations
- 🔄 Database migrations and setup
- 🔄 Data validation and constraints

### Phase 3: Authentication & Security
- 🔄 JWT authentication
- 🔄 Role-based access control
- 🔄 API rate limiting
- 🔄 Input validation and sanitization

### Phase 4: Advanced Features
- 🔄 File upload (patient documents, X-rays)
- 🔄 Email notifications
- 🔄 Appointment reminders
- 🔄 Reporting and analytics

### Phase 5: Production Ready
- 🔄 Comprehensive testing suite
- 🔄 Logging and monitoring
- 🔄 Performance optimization
- 🔄 CI/CD pipeline

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Maintain the hexagonal architecture pattern
- Write comprehensive tests for new features
- Update documentation for API changes
- Use meaningful commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

If you have any questions or need help with the project, please:
1. Check the [Swagger documentation](http://localhost:8080/swagger/)
2. Review the existing issues
3. Create a new issue with detailed information

---

*Built with ❤️ using Go and modern software architecture principles*
