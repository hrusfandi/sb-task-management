# Task Management REST API

A REST API for task management built with Go, featuring JWT authentication and PostgreSQL.

## Tech Stack

- **Go** with Chi Router
- **PostgreSQL** with GORM
- **JWT** Authentication
- **Docker** & Docker Compose

## Quick Start

### Prerequisites
- Docker & Docker Compose installed

### Running the Application

1. Clone the repository:
```bash
git clone https://github.com/hrusfandi/sb-task-management.git
cd sb-task-management
```

2. Start the services:
```bash
docker compose up -d --build
```

The API will be available at `http://localhost:8080`

### Stopping the Application
```bash
# Stop services
docker compose down

# Stop and remove all data
docker compose down -v
```

## API Endpoints

Base URL: `http://localhost:8080/api`

### Authentication
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/register` | Register new user | No |
| POST | `/login` | Login user | No |

### Tasks
| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/tasks` | List all tasks | Yes |
| GET | `/tasks/{id}` | Get task details | Yes |
| POST | `/tasks` | Create new task | Yes |
| PUT | `/tasks/{id}` | Update task | Yes |
| DELETE | `/tasks/{id}` | Delete task | Yes |

### Authentication Header
```
Authorization: Bearer <jwt-token>
```

## API Request Examples

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Create Task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive documentation for the REST API",
    "status": "pending"
  }'
```

### 4. List Tasks
```bash
curl -X GET "http://localhost:8080/api/tasks?page=1&limit=10&status=pending&sort_by=created_at&order=desc" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Query Parameters:**
- `page` - Page number (default: 1)
- `limit` - Items per page (default: 10, max: 100)
- `status` - Filter by status (pending, in_progress, completed)
- `sort_by` - Sort field (created_at, updated_at, title, status)
- `order` - Sort order (asc, desc)

### 5. Get Task Details
```bash
curl -X GET http://localhost:8080/api/tasks/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 6. Update Task
```bash
curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Complete project documentation",
    "description": "Write comprehensive documentation for the REST API with examples",
    "status": "in_progress"
  }'
```

### 7. Delete Task
```bash
curl -X DELETE http://localhost:8080/api/tasks/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Logging

```bash
# View logs
docker compose logs -f app
```

## Project Structure

```
├── main.go                 # Application entry point
├── docker-compose.yml      # Docker configuration
├── Dockerfile              # Multi-stage Docker build
├── docker-entrypoint.sh    # Container startup script
├── config/                 # Configuration
├── database/               # Database connection
├── migrations/             # Database migrations
├── models/                 # Data models
├── handlers/               # Request handlers
├── middleware/             # Auth & logging middleware
├── utils/                  # Utilities (JWT, validation, etc.)
└── routes/                 # API routes
```

## Environment Variables

The application uses these environment variables (configured in docker-compose.yml):
- `DB_USER` - PostgreSQL user (default: postgres)
- `DB_PASSWORD` - PostgreSQL password
- `DB_NAME` - Database name (default: task_management)
- `JWT_SECRET` - Secret key for JWT tokens
- `PORT` - Application port (default: 8080)