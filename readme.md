# 🔒 OWASP Secure Todo App

_A secure, production-ready todo application implementing OWASP best practices for web security._

![Docker](https://img.shields.io/badge/Docker-✓-blue)
![OWASP](https://img.shields.io/badge/OWASP-Compliant-orange)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-336791)

## 🚀 Features

### Security

- ✅ JWT Authentication with refresh tokens
- ✅ Password hashing (bcrypt)
- ✅ Protection against SQL Injection
- ✅ CORS & CSRF mitigation
- ✅ Rate limiting & secure headers

### Functionality

- 📝 CRUD operations for todo items
- 🔐 Role-based access control
- 🧪 Input validation & sanitization
- 🩺 Health check endpoint (`/healthz`)

## 🛠 Tech Stack

**Backend**

- Go (Gin Framework)
- PostgreSQL

**Frontend** _(Optional)_

- React/Next.js

**Infrastructure**

- Docker & Docker Compose

## Create .env file in /api :

```bash
# ========================
# Application Configuration
# ========================
APP_ENV=dev                           # Runtime environment (dev|staging|prod)
APP_PORT=<your-dev-port>                         # Port the application listens on

# ========================
# JWT Configuration
# ========================
# Generate secure secrets with: openssl rand -base64 64
JWT_ACCESS_SECRET=change_me_super_random_64_chars         # 64+ chars for HS256
JWT_REFRESH_SECRET=change_me_even_longer_96_chars         # 96+ chars for HS384/512
JWT_ACCESS_TTL_MIN=15                 # Access token lifetime (minutes)
JWT_REFRESH_TTL_H=168                 # Refresh token lifetime (hours, 7 days)

# ========================
# CSRF Protection
# ========================
# Generate with: openssl rand -base64 32
CSRF_SECRET=32_byte_random_base64_or_hex                  # Exactly 32 bytes

# ========================
# Database Configuration
# ========================
# Format: postgresql://username:password@host:port/database
DB_URL=<your-postgres-db>

# ========================
# Password Hashing
# ========================
BCRYPT_COST=12                        # Hash complexity (4-31, 12 recommended)

# ========================
# Security Headers (Optional)
# ========================
# CSP_DEFAULT_SRC="'self'"
# HSTS_MAX_AGE=31536000
```

## 📂 Project Structure

```bash
.
├── api/               # Go backend
│   ├── handlers/      # Secure endpoints
│   ├── middleware/    # Auth layers
│   └── models/        # DB schema
├── web/               # Frontend
├── docker-compose.yml # Prod/Dev setup
```
