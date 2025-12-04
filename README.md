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
- Docker

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

1. Run PostgreSQL using Docker:
   ```bash
   docker run -d \
     --name chirpy-db \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=chirpy \
     -p 5432:5432 \
     postgres
   ```
   This starts a PostgreSQL container named `chirpy-db` exposing port 5432. The database `chirpy` is automatically created.

   > **Note:** If you need to access the database shell (psql) without installing it locally, you can use:
   > ```bash
   > docker exec -it chirpy-db psql -U postgres -d chirpy
   > ```

2. Run migrations:
   ```bash
   go run github.com/pressly/goose/v3/cmd/goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable" up
   ```

4. (Optional) Generate SQLC code:
   ```bash
   go run github.com/sqlc-dev/sqlc/cmd/sqlc generate
   ```

### Running the Application

Run the server directly:

```bash
go run ./cmd/server
```

Or build and run:

```bash
go build -o out ./cmd/server && ./out
```

The server will start on port `8080`.

## Project Structure

```
chirpy/
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── auth/            # Authentication (JWT, passwords)
│   ├── config/          # Application configuration
│   ├── database/        # Database queries (sqlc generated)
│   └── handler/         # HTTP handlers
├── sql/
│   ├── queries/         # SQLC query definitions
│   └── schema/          # Database migrations
└── bruno-collection/    # API testing collection
```

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
