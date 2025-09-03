# Database Schema

The database schema consists of six core tables:

---

## 1. USERS Table
| Column        | Type                     | Description                  |
|---------------|--------------------------|------------------------------|
| id            | BIGSERIAL PK             | Unique user ID                |
| name          | VARCHAR(100)             | Full name                      |
| email         | VARCHAR(150) UNIQUE      | Login and notifications email |
| password_hash | TEXT                      | Hashed password (bcrypt/Argon2)|
| role          | ENUM('patient','provider','admin') | User type (patient/provider/admin) |
| phone_number  | VARCHAR(20)              | Optional phone for SMS alerts |
| created_at    | TIMESTAMP                | Account creation date          |
| updated_at    | TIMESTAMP                | Last update date                |

---

## 2. PROVIDERS Table
| Column        | Type                     | Description                  |
|---------------|--------------------------|------------------------------|
| id            | BIGSERIAL PK             | Provider unique ID            |
| user_id       | BIGINT FK → USERS.id      | Link to user account          |
| specialization| VARCHAR(100)             | e.g., Cardiologist, Dentist  |
| bio           | TEXT                      | Provider profile description  |
| timezone      | VARCHAR(50)              | For appointment scheduling    |
| created_at    | TIMESTAMP                | Record creation timestamp     |
| updated_at    | TIMESTAMP                | Last update timestamp         |

---

## 3. SERVICES Table
| Column           | Type             | Description            |
|------------------|------------------|------------------------|
| id               | BIGSERIAL PK     | Unique service ID      |
| provider_id      | BIGINT FK → PROVIDERS.id | Provider offering the service |
| title            | VARCHAR(100)     | Name of the service    |
| description      | TEXT             | Detailed description   |
| duration_minutes | INT              | Duration in minutes    |
| price            | DECIMAL(10,2)    | Service cost           |
| created_at       | TIMESTAMP        | Creation timestamp     |
| updated_at       | TIMESTAMP        | Last update timestamp  |

---

## 4. AVAILABILITIES Table
| Column        | Type             | Description            |
|---------------|------------------|------------------------|
| id            | BIGSERIAL PK     | Unique availability ID |
| provider_id   | BIGINT FK → PROVIDERS.id | Provider for this slot      |
| start_time    | TIMESTAMP        | Start of availability  |
| end_time      | TIMESTAMP        | End of availability    |
| is_recurring  | BOOLEAN          | Whether slot repeats   |
| created_at    | TIMESTAMP        | Creation timestamp     |
| updated_at    | TIMESTAMP        | Last update timestamp  |

---

## 5. BOOKINGS Table
| Column        | Type                     | Description                   |
|---------------|--------------------------|-------------------------------|
| id            | BIGSERIAL PK             | Booking unique ID             |
| patient_id    | BIGINT FK → USERS.id      | Patient booking the appointment|
| provider_id   | BIGINT FK → PROVIDERS.id  | Doctor assigned               |
| service_id    | BIGINT FK → SERVICES.id   | Service being booked          |
| start_time    | TIMESTAMP                | Appointment start time         |
| end_time      | TIMESTAMP                | Appointment end time           |
| status        | ENUM('pending','confirmed','completed','canceled') | Booking status |
| created_at    | TIMESTAMP                | Booking creation timestamp     |
| updated_at    | TIMESTAMP                | Last update timestamp          |

---

## 6. NOTIFICATIONS Table
| Column     | Type                    | Description                    |
|------------|-------------------------|--------------------------------|
| id         | BIGSERIAL PK            | Notification unique ID          |
| user_id    | BIGINT FK → USERS.id     | Recipient of notification      |
| message    | TEXT                     | Notification content           |
| type       | ENUM('email','sms')     | Notification type              |
| status     | ENUM('pending','sent','failed') | Delivery status            |
| created_at | TIMESTAMP               | Creation timestamp             |
| sent_at    | TIMESTAMP NULL          | Time when notification sent    |

---

## Relationships Overview
