# Auth-Go

[![language](https://img.shields.io/badge/language-Go-00ADD8?labelColor=2F2F2F)](https://go.dev/doc/)
[![version](https://img.shields.io/badge/version-1.25-9C27B0?labelColor=2F2F2F)](https://go.dev/doc/install)

An API with Authentication using JWT (with and without refresh token) and cookie based

## Overview

This repository contains a small HTTP API that manages users backed by PostgreSQL. It was built using **Go 1.25** and assumes a PostgreSQL instance is available, if not... I've provided a Docker command to run one locally.

This is an area of ​​study aimed at testing different ways to create an authentication system.

## Technologies used

- Go 1.25
- PostgreSQL (run with Docker)
- chi (router)
- jackc/pgx driver
- gorilla/schema for query decoding
- godotenv for environment variables

## Requirements

- Go 1.25 or later installed on your machine
- Docker (for running PostgreSQL locally) or an accessible Postgres instance
- Make sure the `DATABASE_URL` environment variable points to a reachable Postgres database

## Database Schema

Use the schema below to initialize the database:

```sql
CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(50) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  session_token VARCHAR(255),
  csrf_token VARCHAR(255),
  age INT not null
 );
```

## Environment variables

Create a `.env` file (example provided in the project). The service expects at least:

    ```
      JWT_PORT = "localhost:8001"
      COOKIE_PORT = "localhost:8000"
      JWT_REFRESH_PORT = "localhost:8002"

      #-------------------------------------

      DATABASE_URL="postgres://root:example@localhost:5432/postgres"

      #-------------------------------------

      ISSUER="golang"
      JWT_SECRET_KEY="secretKeyExample"
      TOKEN_DURATION="15" # In minutes...
    ```

## How to use

1. Clone the repo:

   ```bash
   git clone <repo-url>
   cd go-full-crud
   ```

2. Create a `.env` file (or use the provided `.env.example`) and set `DATABASE_URL`, `JWT_PORT`, `COOKIE_PORT` and `JWT_REFRESH_PORT` as needed.

3. Start PostgreSQL with Docker:

   ```bash
   docker run --name go-postgres -e POSTGRES_PASSWORD=example -e POSTGRES_USER=root -e POSTGRES_DB=golang-database -p 5432:5432 -d postgres:15
   ```

   Adjust user/password/db name to match your `DATABASE_URL` if necessary.

4. Apply the database schema (run the SQL above) using `psql` or a GUI tool.

5. Run the service:

   ```bash
   # from the project root
   go run cmd/cookie-based/main.go         # This one for the Cookie token with csrf protection.
   go run cmd/jwt-based/main.go            # This one for the JWT but without the refresh token, only the expiration time
   go run cmd/jwt-refresh-based/main.go    # This one for the JWT with refresh token.
   ```
