# Sera Ale Job Board API

A RESTful API for companies to post jobs and applicants to apply, built with Go, Gin, GORM, PostgreSQL, Cloudinary, and Clean Architecture.

## Setup

1. Clone the repo
2. Copy `.env` from the example and fill in your credentials
3. Run DB migrations (see `migrations/`)
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Run the server:
   ```bash
   go run cmd/main.go
   ```

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
- `JWT_SECRET`: Secret for JWT signing
- `CLOUDINARY_URL`: Cloudinary API URL

## Tech Stack
- Go (Gin, GORM)
- PostgreSQL (Neon)
- Cloudinary (file uploads)
- JWT (auth)
- Clean Architecture
- Swagger (API docs)

## API Docs

Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

