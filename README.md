
# Doctor Appointment Backend (Go)

A backend-only doctor appointment booking system built in **Go (Golang)**. 
This project is designed to help you learn backend development with Go while building a real-world, scalable API.

## Documentation
- [Introduction](docs/INTRODUCTION.md)
- [Database Schema](docs/DATABASE_SCHEMA.md)
- [API Endpoints](docs/API_ENDPOINTS.md)
- [Workflows](docs/WORKFLOWS.md)
- [Development Roadmap](docs/ROADMAP.md)


## 🚀 Overview
This backend provides APIs for:
- User authentication (patients and doctors)
- Doctor profile and specialization management
- Services and availability scheduling
- Appointment booking and cancellation
- Notifications via email/SMS (background jobs)
- Optional payment integration (Stripe/PayPal)

The architecture follows a **modular, clean structure**, using modern Go libraries and tools.

## 🛠 Tech Stack
- **Language:** Go
- **Web Framework:** Gin
- **Database:** PostgreSQL (GORM ORM)
- **Cache / Locking:** Redis
- **Background Jobs:** Go worker pools
- **Authentication:** JWT
- **Notifications:** Twilio (SMS), SendGrid (Email)
- **Deployment:** Docker & Kubernetes
- **Monitoring:** Prometheus + Grafana

## 📂 Project Structure
```
appointment-app/
├── cmd/
│   └── server/
│       └── main.go             # Entry point
├── internal/
│   ├── auth/                   # JWT, middleware
│   ├── booking/                # Booking engine logic
│   ├── notifications/          # Background jobs for SMS/email
│   ├── payments/               # Stripe/PayPal integration
│   └── users/                  # Providers and patients
├── pkg/
│   ├── logger/                 # Logging utilities
│   └── config/                  # App configuration
├── api/
│   └── http/                   # REST API handlers
├── scripts/                    # DB migrations, seeds
├── go.mod
└── go.sum
```

## 🔗 Key API Endpoints
| Endpoint | Method | Purpose |
|-----------|--------|---------|
| `/auth/register` | POST | Register new user |
| `/auth/login` | POST | Authenticate user, return JWT |
| `/providers` | POST | Create a new provider (admin only) |
| `/providers` | GET | List all providers |
| `/providers/{id}` | GET | Get provider details |
| `/providers/{id}/services` | POST | Add service for provider |
| `/providers/{id}/availability` | POST | Set provider availability |
| `/bookings` | POST | Create appointment booking |
| `/bookings/{id}` | GET | Get booking details |
| `/bookings/{id}/cancel` | PUT | Cancel booking |

## 🗄 Database Schema (Summary)
- **USERS:** Stores patients and providers basic data
- **PROVIDERS:** Specialization, bio, timezone
- **SERVICES:** Provider’s services with duration and price
- **AVAILABILITIES:** Available time slots
- **BOOKINGS:** Appointment records linking patient, provider, and service
- **NOTIFICATIONS:** Stores scheduled or sent notifications

## 📦 Installation & Setup
1. Clone the repo:
   ```bash
   git clone https://github.com/your-username/doctor-appointment-backend.git
   cd doctor-appointment-backend
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Configure environment variables (`.env`):
   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/appointment_db
   REDIS_URL=redis://localhost:6379
   JWT_SECRET=your_jwt_secret
   ```
4. Run database migrations:
   ```bash
   go run scripts/migrate.go
   ```
5. Start the server:
   ```bash
   go run cmd/server/main.go
   ```

## 📅 Project Timeline (2 Months)
- **Phase 1 (Week 1-2):** Setup project, authentication, users module
- **Phase 2 (Week 3-4):** Providers and services
- **Phase 3 (Week 5):** Availability scheduling
- **Phase 4 (Week 6):** Booking system with Redis lock
- **Phase 5 (Week 7):** Notifications
- **Phase 6 (Week 8):** Payments (optional), testing, and deployment

## 🤝 Contributing
Contributions are welcome! Please fork the repository and submit a pull request.

## 📜 License
This project is licensed under the MIT License.
