```markdown
# Housing Lottery Web App

A web application for managing housing/condominium lotteries per sub-city, built with **Golang**, **Fiber**, **Vue.js**, and **PostgreSQL**. The system allows authorities to upload applicants, run lotteries per sub-city, and announce winners in a transparent and auditable way.

---

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Folder Structure](#folder-structure)
- [Database Schema](#database-schema)
- [Setup Instructions](#setup-instructions)
- [API Endpoints](#api-endpoints)
- [Future Enhancements](#future-enhancements)

---

## Project Overview

- Authorities upload a list of applicants via Excel.
- Each sub-city has its own lottery draw.
- The system randomly selects winners and stores them with position order.
- Admins can view, manage, and audit lotteries and applicants.

### Key Features

- Admin registration and authentication
- Sub-city management
- Excel-based applicant upload
- Lottery engine with random selection
- Audit logging
- Clean architecture separation for maintainability

---

## Architecture

We follow **Clean Architecture** principles:

- **Domain**: Core entities and business logic (`internal/domain`)
- **Usecase/Service**: Application logic (`internal/usecase`)
- **Infrastructure**: Database repositories (`internal/infrastructure`)
- **Delivery/Handler**: HTTP handlers using Fiber (`internal/delivery/http`)
- **Middleware**: Auth and request handling (`internal/delivery/middleware`)
- **Utils**: Excel parsing helpers (`internal/utils`)

### Stack

- **Backend**: Golang (Fiber, pgx)
- **Frontend**: Vue.js
- **Database**: PostgreSQL
- **Authentication**: Bcrypt + JWT (planned)
- **File Parsing**: Excelize

---

## Folder Structure

```

housing-lottery/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── domain/
│   ├── usecase/
│   ├── infrastructure/
│   ├── delivery/
│   ├── repository/
│   └── utils/
├── migrations/
├── pkg/
├── go.mod
└── README.md

```

---

## Database Schema

### Tables

- `admins`
- `subcities`
- `upload_batches`
- `applicants`
- `lotteries`
- `lottery_winners`
- `audit_logs`

### Sample Admin Entries

```

id                                   | email             | password_hash      | created_at
-------------------------------------+-------------------+--------------------+-------------------------
65e4038c-12de-4fa8-8e97-cc1463aef47b  | [admin@test.com](mailto:admin@test.com)    | <hashed_password>  | 2026-02-16 09:52:02
d79322fa-4567-4c7f-90ca-c0f1198fefcb  | [adminn@test.com](mailto:adminn@test.com)   | <hashed_password>  | 2026-02-16 09:58:19

````

---

## Setup Instructions

### 1. Clone the repository

```bash
git clone https://github.com/FANIMAN/housing-lottery.git
cd housing-lottery
````

### 2. Create PostgreSQL database

```sql
CREATE DATABASE housing_lottery;
\c housing_lottery
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

### 3. Create `.env` file

```env
DATABASE_URL=postgres://postgres:YOUR_PASSWORD@localhost:5432/housing_lottery?sslmode=disable
```

### 4. Install dependencies

```bash
go mod tidy
```

### 5. Run the server

```bash
go run cmd/server/main.go
```

---

## API Endpoints

### POST /admin/register

Request:

```json
{
  "email": "admin@test.com",
  "password": "123456"
}
```

Response:

```json
{
  "message": "admin created"
}
```

More endpoints (login, subcity CRUD, upload, lottery) are in progress.

---

## Future Enhancements

* Admin login + JWT authentication
* Subcity CRUD endpoints
* Excel applicant upload
* Lottery engine with seed-based random selection
* Winner announcement with position ordering
* Audit logs
* Vue.js dashboard frontend

---

## Notes

* Clean Architecture for maintainability
* Repository pattern for DB operations
* Fiber v2 for HTTP
* bcrypt for password hashing
* Excelize for Excel parsing

```

---

