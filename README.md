# FAuthless-Go

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

CREATE TABLE IF NOT EXISTS sessions (
  id VARCHAR(255) PRIMARY KEY NOT NULL,
  username VARCHAR(50) NOT NULL,
  is_revoked BOOL NOT null default false,
  refresh_token VARCHAR(512) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  expires_at TIMESTAMP
 );
```

## Environment variables

Create a `.env` file (example provided in the project). The service expects at least:

      JWT_PORT = "localhost:8001"
      COOKIE_PORT = "localhost:8000"
      JWT_REFRESH_PORT = "localhost:8002"

      #-------------------------------------

      DATABASE_URL="postgres://root:example@localhost:5432/postgres"

      #-------------------------------------

      ISSUER="golang"
      JWT_SECRET_KEY="secretKeyExample"
      TOKEN_DURATION="15" # In minutes...

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
- **PATCH** | `/users/{username}`
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

- JWT + Refresh | `JWT_REFRESH_PORT=8002`

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

## 2.3 JWT + Refresh Token

### Overview

This authentication mode issues two tokens on login: a short-lived **access token** (JWT) used for API requests, and a longer-lived **refresh token** used to obtain new access tokens without re-authenticating. Refresh tokens are stored server-side (sessions table) so they can be revoked.

### Endpoints

- **POST /login** — returns `session_id`, `access_token`, `refresh_token`, and their expiration times.
- **POST /renew** — accepts a refresh token and returns a new access token and its expiration.
- **POST /revoke/{id}** — revoke a session by `id` (sets session as revoked). Returns `204 No Content` on success.

### Login (example)

Request:

```bash
curl -s -X POST http://localhost:8002/login   -H "Content-Type: application/json"   -d '{"username":"alice","password":"plainPassword123"}'
```

Successful response (201 Created):

```json
{
  "session_id": "8a4f2d9e-1a3b-4c2a-9b8f-0a1b2c3d4e5f",
  "access_token": "<JWT_ACCESS_TOKEN>",
  "refresh_token": "<JWT_REFRESH_TOKEN>",
  "access_token_expires_at": "2025-11-30T12:34:56Z",
  "refresh_token_expires_at": "2025-12-01T12:34:56Z"
}
```

### Renew access token

Request:

```bash
curl -s -X POST http://localhost:8002/renew   -H "Content-Type: application/json"   -d '{"refresh_token":"<JWT_REFRESH_TOKEN>"}'
```

Successful response (200 OK):

```json
{
  "access_token": "<NEW_JWT_ACCESS_TOKEN>",
  "access_token_expires_at": "2025-11-30T12:44:56Z"
}
```

### Revoke session

```bash
curl -X POST http://localhost:8002/revoke/8a4f2d9e-1a3b-4c2a-9b8f-0a1b2c3d4e5f
```

---

## 3. Protected Endpoints

### **GET /api/v1/users**

Paginated list\
Supports: - `?page=1` - `?size=25`

### **GET /api/v1/users/{id}**

### **PATCH /api/v1/users/{username}**

Updates user info (must be account owner)

### **DELETE /api/v1/users/{username}**

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

- The project still needs tests case and some coverage about the edge cases...
- If you need any help or are having any issue with it, contact `rafael.cr.carneiro@gmail.com` for assistance.
