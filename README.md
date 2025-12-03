# Chirpy

Chirpy is a social media API built with Go. It allows users to create accounts, post "chirps" (messages), and interact with other users.

## Features

- **User Management**: Create users, login, and authentication via JWT.
- **Chirps**: Create, read, and delete chirps.
- **Metrics**: Admin endpoint to view server hits.
- **Database**: Persistent storage using PostgreSQL.
- **API Testing**: Includes a [Bruno](https://www.usebruno.com/) collection for easy API testing.

## Tech Stack

- **Language**: Go
- **Database**: PostgreSQL
- **ORM/Query Builder**: SQLC
- **Migrations**: Goose
- **Router**: Standard library `net/http` with `ServeMux`
- **Environment**: `godotenv`

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL
- [Goose](https://github.com/pressly/goose) (for migrations)
- [SQLC](https://sqlc.dev/) (for generating Go code from SQL)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/KleitonBarone/chirpy.git
   cd chirpy
   ```

2. Create a `.env` file in the root directory:
   ```env
   DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
   PLATFORM="dev"
   JWT_SECRET="your_secret_key"
   ```

### Database Setup

1. Install PostgreSQL (Linux example):
   ```bash
   sudo apt install postgresql postgresql-contrib
   sudo -u postgres psql
   # Inside psql:
   CREATE DATABASE chirpy;
   ALTER USER postgres PASSWORD 'postgres';
   \q
   ```

2. Install Goose:
   ```bash
   go install github.com/pressly/goose/v3/cmd/goose@latest
   ```

3. Run migrations:
   ```bash
   cd sql/schema
   goose postgres "postgres://postgres:postgres@localhost:5432/chirpy" up
   cd ../..
   ```

4. (Optional) Generate SQLC code:
   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   sqlc generate
   ```

### Running the Application

Build and run the server:

```bash
go build -o out && ./out
```

The server will start on port `8080`.

### Running Tests

```bash
go test -v ./...
```

## API Endpoints

### Health
- `GET /api/healthz`: Check API health.

### Users
- `POST /api/users`: Create a new user.
- `POST /api/login`: Login and receive a JWT.

### Chirps
- `POST /api/chirps`: Create a new chirp.
- `GET /api/chirps`: Get all chirps.
- `GET /api/chirps/{chirpID}`: Get a specific chirp.

### Admin
- `GET /admin/metrics`: View server hit count.
- `POST /admin/reset`: Reset server hit count.

## API Testing with Bruno

This project includes a [Bruno](https://www.usebruno.com/) collection in the `bruno-collection` directory. You can import this collection into Bruno to easily test the API endpoints.
