# BookMyShow MVP Clone (Go + Svelte)
![GitHub License](https://img.shields.io/github/license/AdityaShankar1/bmsclone-go?style=flat-square&color=blue)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/AdityaShankar1/bmsclone-go/ci.yml?style=flat-square&label=build)
![GitHub last commit](https://img.shields.io/github/last-commit/AdityaShankar1/bmsclone-go?style=flat-square)

A high-performance, robust movie booking application with integrated financial tools and multiplex chaining.

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![Svelte](https://img.shields.io/badge/Svelte-FF3E00?style=flat-square&logo=svelte&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat-square&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=flat-square&logo=redis&logoColor=white)
![Bootstrap](https://img.shields.io/badge/Bootstrap-7952B3?style=flat-square&logo=bootstrap&logoColor=white)

<img width="1470" height="956" alt="Screenshot 2026-03-30 at 3 38 32 PM" src="https://github.com/user-attachments/assets/901cf57c-2c13-4432-84fb-152341b6213b" />

## 🚀 Key Features

![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=flat-square&logo=github-actions&logoColor=white)
![GCP](https://img.shields.io/badge/GCP-4285F4?style=flat-square&logo=google-cloud&logoColor=white)

### 🎬 Cinema & Booking
- **Multiplex Chains**: Support for PVR-INOX, Cinepolis, and Independent theatres.
- **Dual Flow**: Book by movie or explore specific theatres/chains.
- **Seat Locking**: Real-time Redis-based seat locking (5-minute TTL) to prevent double booking.
- **Smart Formats**: Dynamic filtering for IMAX, 4DX, and 2D based on theatre capabilities.

### Walkthrough:

<img width="1122" height="655" alt="image" src="https://github.com/user-attachments/assets/5efbfb08-f64a-400b-b413-fe3ff7fa2728" />
<img width="1123" height="684" alt="image" src="https://github.com/user-attachments/assets/f2d94029-406a-4785-b95b-40a9815259ae" />
<img width="1130" height="477" alt="image" src="https://github.com/user-attachments/assets/7e0b16f5-e0cf-4237-a34d-6328c1fa8978" />
<img width="1125" height="729" alt="image" src="https://github.com/user-attachments/assets/ef721c7e-7560-462f-803d-0ece040bd44d" />
<img width="1466" height="917" alt="image" src="https://github.com/user-attachments/assets/f3d2d99a-a3d9-4c3e-88ce-1f1d5f2d3096" />
<img width="1020" height="493" alt="image" src="https://github.com/user-attachments/assets/1346daea-e26a-4ceb-8d02-0240c742fd8a" />
<img width="1020" height="621" alt="image" src="https://github.com/user-attachments/assets/43e663b4-a8e3-4199-8de8-f5c8e83a16c9" />
<img width="1012" height="475" alt="image" src="https://github.com/user-attachments/assets/e63ce6b0-b61e-417b-b184-1e0275c04cd9" />


### 💰 Financial Ecosystem (USP)
- **Unified Wallet**: Manage payments, investments, and refunds in one place.
- **BMSCash Refunds**: Automatic 70% refund to wallet for eligible cancellations (>2h before showtime).
- **Inflation Protection**: Invest ₹500+ to "freeze" ticket rates and unlock a permanent 20% discount.
- **Transaction Logs**: Full transparency for all payments, investments, and refunds.

<img width="1126" height="837" alt="image" src="https://github.com/user-attachments/assets/0809aa24-274a-4892-8089-d63e1af84dcc" />


### ⭐ Engagement
- **Star Ratings**: Rate your movie experience after booking.
- **Modern UI**: responsive Svelte frontend with Bootstrap aesthetics.

## 🛠️ Tech Stack
- **Backend**: Go (Gin, GORM, Redis-Go)
- **Database**: PostgreSQL (Persistence), Redis (Locking)
- **Frontend**: Svelte 5 + Bootstrap 5
- **Communication**: RESTful API / JSON

## 📊 Architecture Diagram

```mermaid
graph TD
    User((User))
    subgraph "Frontend (Svelte)"
        Home[Home View]
        Cinema[Cinema View]
        Payments[Payment Gateways]
        WalletView[Wallet/Transactions]
    end
    
    subgraph "Backend (Go)"
        API[Gin API Layer]
        Auth[Handlers/Middleware]
        Lock[Redis SETNX Lock]
    end
    
    subgraph "Storage"
        PG[(PostgreSQL)]
        R[(Redis)]
    end

    User --> Home
    Home --> API
    API --> Lock
    Lock --> R
    API --> PG
    Payments --> API
    WalletView --> API
```

## ⚙️ Setup Instructions

1. **Prerequisites**: Go 1.25+, Docker, Node.js.
2. **Run Databases**:
   ```bash
   docker run --name bms-db -e POSTGRES_PASSWORD=password123 -p 5432:5432 -d postgres:alpine
   docker run --name bms-redis -p 6379:6379 -d redis:alpine
   ```
3. **Start Backend**:
   ```bash
   cd backend
   go run main.go
   ```
4. **Start Frontend**:
   ```bash
   cd frontend
   npm run dev
    ```
   
## 🐳 Docker Deployment

The entire stack can be launched with a single command:

```bash
docker compose up --build
```

- **Frontend**: http://localhost:5173
- **Backend**: http://localhost:8080
- **Postgres**: localhost:5432
- **Redis**: localhost:6379

## 🧪 CI/CD
Automated builds and tests are handled via **GitHub Actions** (`.github/workflows/ci.yml`), ensuring that every push to `main` is buildable and robust.

## A Note on the deployment:
If you're reading this after mid-Apr26, you'll likely hit an error with the URL as my GCP trial would expire by then.
