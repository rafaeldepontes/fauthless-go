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

# Endpoints Overview

## Public

- **POST /register**
- **POST /login** (auth-type dependent: Cookie, JWT, JWT+Refresh)

## Protected (all require authentication)

Base path: `/api/v1`

- **GET** | `/users/{id}`
- **GET** | `/users?page=$&size=$`
- **PATCH** | `/users`
- **DELETE** | `/users/{username}`

---

## 1. Registration

### **POST /register**

**Body**

```json
{
  "username": "alice",
  "hashed_password": "plainPassword123",
  "age": 30
}
```

**Example**

```bash
curl -i -X POST http://localhost:8000/register   -H "Content-Type: application/json"   -d '{"username":"alice","hashed_password":"plainPassword123","age":30}'
```

---

## 2. Authentication Methods

Your server can run in one of three modes:

- Cookie-based | `COOKIE_PORT=8000`

- JWT-based | `JWT_PORT=8001`

- JWT + Refresh | `JWT_REFRESH_PORT=8002 (WIP)`

---

## 2.1 Cookie-Based Authentication

### Login

**POST /login**

Creates: - `session_token` cookie (HttpOnly) - `crsf_token` cookie

**Example**

```bash
curl -i -c cookies.txt -X POST http://localhost:8000/login   -H "Content-Type: application/json"   -d '{"username":"alice","password":"plainPassword123"}'
```

Read CSRF token:

```bash
CRSF=$(awk '/csrf_token/ {print $7}' cookies.txt)
```

### Access Protected Routes

```bash
curl -b cookies.txt http://localhost:8000/api/v1/users
```

### Mutating Requests (Require CSRF Header)

```bash
curl -b cookies.txt -X PATCH http://localhost:8000/api/v1/users   -H "Content-Type: application/json"   -H "X-CSRF-Token: $CRSF"   -d '{"age": 31}'
```

---

## 2.2 JWT-Based Authentication

### Login

**POST /login**

Response:

```json
{
  "token": "<JWT>"
}
```

Example:

```bash
TOKEN=$(curl -s -X POST http://localhost:8001/login   -H "Content-Type: application/json"   -d '{"username":"alice","password":"plainPassword123"}' | jq -r '.token')
```

### Access Protected Routes

```bash
curl -H "Authorization: Bearer $TOKEN"   http://localhost:8001/api/v1/users
```

---

## 2.3 JWT + Refresh Token (WIP)

## `WIP`

---

## 3. Protected Endpoints

### **GET /api/v1/users**

Paginated list\
Supports: - `?page=1` - `?size=25`

### **GET /api/v1/users/{id}**

### **PATCH /api/v1/users**

Updates user info (must be account owner)

### **DELETE /api/v1/users**

Deletes account (must be account owner)

---

## 4. Quick Reference

```bash
# Register
curl -X POST http://localhost:8000/register   -H "Content-Type: application/json"   -d '{"username":"alice","hashed_password":"pwd","age":30}'

# Cookie Login
curl -i -c cookies.txt -X POST http://localhost:8000/login   -H "Content-Type: application/json"   -d '{"username":"alice","password":"pwd"}'

# JWT Login
TOKEN=$(curl -s -X POST http://localhost:8001/login   -H "Content-Type: application/json"   -d '{"username":"alice","password":"pwd"}' | jq -r '.token')
```

---

## 5. Notes

- JWT refresh login endpoint still needs implementation
