# chirpy
Chirpy Golang project

## Run locally

Make sure to have:

.env configured:
```
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
```

postgres running (linux setup example):
```
sudo apt install postgresql postgresql-contrib
psql --version
sudo passwd postgres
sudo -u postgres psql
CREATE DATABASE chirpy;
ALTER USER postgres PASSWORD 'postgres';
SELECT version();
exit
psql "postgres://postgres:postgres@localhost:5432/chirpy"
```

goose for migration:
```
go install github.com/pressly/goose/v3/cmd/goose@latest
cd sql/schema
goose postgres postgres://postgres:postgres@localhost:5432/chirpy up
psql chirpy
\dt
```

generate ORM queries:
```
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc version
sqlc generate
```

build and run:
```
go build -o out && ./out
```
