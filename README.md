# ✈️ Travel Planner

A full-stack, production-grade travel planning application built with **Go (Gin)** and **Vue 3 (Vite)**. Browse destinations, leave reviews, and build custom itineraries — all backed by MongoDB with enterprise-level reliability.

---

## 🏗️ Architecture

```
┌──────────────┐      ┌──────────────────┐      ┌─────────┐
│   Frontend   │─────▶│   Backend API    │─────▶│ MongoDB │
│  Vue 3/Vite  │ :80  │   Go / Gin       │ :5000│         │
│    Nginx     │      │                  │      │  Mongo 7│
└──────────────┘      └──────────────────┘      └─────────┘
```

| Layer    | Stack                                       |
|----------|---------------------------------------------|
| Frontend | Vue 3, Vite, Pinia, Vue Router, Axios        |
| Backend  | Go 1.24, Gin, MongoDB Driver, JWT, bcrypt    |
| Database | MongoDB 7 with connection pooling             |
| Deploy   | Docker, Docker Compose, Nginx                |

---

## ✨ Features

### Core
- 🌍 **Destination Catalog** — Browse, search & filter destinations by name, country, or rating
- ⭐ **Reviews & Ratings** — Submit reviews with automatic average-rating recalculation
- 📋 **Itineraries** — Create, edit & manage personal travel itineraries with multiple destinations
- 🔐 **Authentication** — JWT-based auth with bcrypt password hashing (cost 12)
- 👤 **Role-Based Access** — Admin dashboard for managing destinations and moderating reviews

### Production-Grade Infrastructure
- ⚡ **100K Request Capacity** — MongoDB connection pool (200 max / 20 min), tuned HTTP server timeouts
- 🛡️ **Rate Limiting** — Per-IP token-bucket limiter (100 req/s default, configurable)
- 📊 **Structured Logging** — JSON logs via zerolog with request ID tracing
- 🔗 **Request Tracing** — ULID-based `X-Request-ID` on every request
- ⏱️ **Request Timeouts** — Configurable per-request deadline (default 30s)
- 🗜️ **GZIP Compression** — Response compression reducing bandwidth 60-80%
- 🏥 **Health Probes** — `/health/live` (liveness) and `/health/ready` (DB readiness)
- 🔄 **Graceful Shutdown** — SIGINT/SIGTERM handling with 15s connection drain
- 🛡️ **Security Headers** — X-Content-Type-Options, X-Frame-Options, CSP, XSS Protection
- 🔒 **Input Sanitization** — Regex escaping in search filters to prevent ReDoS attacks
- 🔁 **Frontend Retry** — Exponential backoff (3 retries) on network/429/5xx errors

---

## 📁 Project Structure

```
Travel Planner/
├── backend/
│   ├── config/          # MongoDB connection & pooling
│   ├── errors/          # Typed AppError for clean error handling
│   ├── handlers/        # HTTP request handlers (auth, destinations, reviews, itineraries)
│   ├── logger/          # Structured zerolog setup
│   ├── middleware/       # Auth, rate limiting, request ID, timeout, logging
│   ├── models/          # Data models & input validation structs
│   ├── repository/      # MongoDB data access layer
│   ├── routes/          # Route definitions & handler wiring
│   ├── seed/            # Database seeding (admin user + sample destinations)
│   ├── service/         # Business logic layer
│   ├── Dockerfile       # Multi-stage build (Go → Alpine)
│   ├── main.go          # Application entry point
│   └── .env.example     # Environment variable template
├── frontend/
│   ├── src/
│   │   ├── api/         # Axios client with retry & timeout
│   │   ├── components/  # Reusable Vue components
│   │   ├── router/      # Vue Router with auth guards
│   │   ├── stores/      # Pinia state management
│   │   └── views/       # Page components (Home, Login, Admin, etc.)
│   ├── Dockerfile       # Multi-stage build (Node → Nginx)
│   ├── nginx.conf       # Production Nginx config
│   └── vite.config.js   # Vite config with code splitting
├── docker-compose.yml   # Full stack orchestration
└── README.md
```

---

## 🚀 Getting Started

### Prerequisites

| Tool     | Version  |
|----------|----------|
| Go       | 1.24+    |
| Node.js  | 22.12+   |
| MongoDB  | 7.0+     |
| Docker   | 24+      |

### Option 1: Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/Aniketssg1/Travel-Planner.git
cd Travel-Planner

# Configure environment
cp backend/.env.example backend/.env
# Edit backend/.env with your JWT_SECRET (min 32 chars)

# Start everything
docker-compose up --build
```

The app will be available at:
- **Frontend:** http://localhost
- **API:** http://localhost:5000
- **MongoDB:** localhost:27017

### Option 2: Local Development

**Backend:**
```bash
cd backend
cp .env.example .env
# Edit .env with your MongoDB URI and JWT secret
go run .
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:5000

---

## ⚙️ Configuration

All configuration is via environment variables. See [`backend/.env.example`](backend/.env.example) for the full list.

| Variable | Default | Description |
|----------|---------|-------------|
| `MONGO_URI` | *required* | MongoDB connection string |
| `DB_NAME` | *required* | Database name |
| `JWT_SECRET` | *required* | Signing secret (min 32 chars) |
| `PORT` | `5000` | Server port |
| `GIN_MODE` | `release` | `debug` or `release` |
| `ALLOWED_ORIGINS` | `http://localhost:3000` | CORS origins (comma-separated) |
| `MONGO_MAX_POOL_SIZE` | `200` | Max MongoDB connections |
| `MONGO_MIN_POOL_SIZE` | `20` | Warm MongoDB connections |
| `RATE_LIMIT_RPS` | `100` | Max requests/second per IP |
| `REQUEST_TIMEOUT_SECONDS` | `30` | Per-request timeout |
| `LOG_LEVEL` | `info` | Log level (trace/debug/info/warn/error/fatal) |

---

## 📡 API Reference

### Authentication
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/register` | — | Register new user |
| POST | `/api/auth/login` | — | Login & get JWT |

### Destinations
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/destinations` | — | List with search, filter & pagination |
| GET | `/api/destinations/:id` | — | Get destination + reviews |
| POST | `/api/destinations` | Admin | Create destination |
| PUT | `/api/destinations/:id` | Admin | Update destination |
| DELETE | `/api/destinations/:id` | Admin | Delete destination + reviews |

**Query Parameters:** `name`, `country`, `minRating`, `page`, `limit` (max 50)

### Reviews
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/destinations/:id/reviews` | User | Submit review (1-5 stars) |
| GET | `/api/reviews` | Admin | List all reviews |
| DELETE | `/api/reviews/:id` | Admin | Delete review |

### Itineraries
| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/itineraries` | User | Get user's itineraries |
| POST | `/api/itineraries` | User | Create itinerary |
| PUT | `/api/itineraries/:id` | User | Update itinerary |
| DELETE | `/api/itineraries/:id` | User | Delete itinerary |

### Health
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health/live` | Liveness probe (always 200) |
| GET | `/health/ready` | Readiness probe (pings MongoDB) |

---

## 🔒 Default Admin Account

> **Email:** `admin@travel.com`
> **Password:** `Admin@123`

⚠️ Change these credentials immediately in production.

---

## 🛡️ Production Middleware Pipeline

```
Request → Recovery → RequestID → RateLimit → Logger → Timeout → GZIP → CORS → SecurityHeaders → Handler
```

---

## 🐳 Docker Compose Services

| Service | Image | Port | Health Check |
|---------|-------|------|-------------|
| `mongo` | mongo:7 | 27017 | `mongosh ping` |
| `backend` | Go multi-stage | 5000 | `/health/live` |
| `frontend` | Node → Nginx | 80 | `/health` |

---
