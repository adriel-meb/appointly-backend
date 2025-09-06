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
# Doctor Appointment Backend - Models Documentation

This document explains the **data models** used in the Doctor Appointment Booking Backend project.  
It covers each table (model), its fields, relationships, and usage within the system.

---

## 1. User Model (`users`)

The **User** model represents every person who interacts with the system.  
A user can have one of three roles:
- **Patient** → can book appointments  
- **Provider (Doctor)** → can offer services and availability  
- **Admin** → manages the system  

### Fields
| Field         | Type      | Description |
|---------------|-----------|-------------|
| `id`          | BIGSERIAL PK | Unique numeric ID for each user (primary key). |
| `name`        | VARCHAR(100) | Full name of the user. |
| `email`       | VARCHAR(150) UNIQUE | Used for login and notifications (must be unique). |
| `password_hash` | TEXT | Encrypted password stored in DB (hashed using bcrypt/Argon2). |
| `role`        | ENUM('patient','provider','admin') | Defines the type of user. Defaults to **patient**. |
| `phone_number` | VARCHAR(20), NULL | Optional phone number for SMS notifications. |
| `created_at`  | TIMESTAMP | When the user account was created. |
| `updated_at`  | TIMESTAMP | Last update timestamp. |
| `deleted_at`  | TIMESTAMP (nullable) | Soft delete column (used for non-permanent deletion). |

### Notes
- All users exist in this table.  
- Providers are also users, but with an additional **Provider profile**.  
- Patients only exist as **users** (they don’t need a provider profile).  

---

## 2. Provider Model (`providers`)

The **Provider** model extends a `User` into a **Doctor/Healthcare provider**.  
It contains extra information about the doctor’s specialty, bio, and timezone.

### Fields
| Field         | Type      | Description |
|---------------|-----------|-------------|
| `id`          | BIGSERIAL PK | Unique provider ID. |
| `user_id`     | BIGINT FK → users.id | Links provider profile to a user account. |
| `specialization` | VARCHAR(100) | Doctor’s specialty (e.g., Cardiologist, Dentist). |
| `bio`         | TEXT | A description of the provider’s background and experience. |
| `timezone`    | VARCHAR(50) | Used to manage appointment times correctly. |
| `created_at`  | TIMESTAMP | When the provider profile was created. |
| `updated_at`  | TIMESTAMP | Last profile update. |

### Notes
- Only users with `role = provider` should have entries in this table.  
- A provider can offer **multiple services** and **availability slots**.  

---

## 3. Service Model (`services`)

The **Service** model describes what a provider offers (e.g., consultation, surgery, therapy session).  

### Fields
| Field         | Type      | Description |
|---------------|-----------|-------------|
| `id`          | BIGSERIAL PK | Unique service ID. |
| `provider_id` | BIGINT FK → providers.id | Which provider offers this service. |
| `title`       | VARCHAR(100) | Service name (e.g., "General Consultation"). |
| `description` | TEXT | Detailed description of the service. |
| `duration_minutes` | INT | Duration of the service (e.g., 30 minutes). |
| `price`       | DECIMAL(10,2) | Cost of the service. |
| `created_at`  | TIMESTAMP | When the service was added. |
| `updated_at`  | TIMESTAMP | When the service was last updated. |

### Notes
- A provider can have many services.  
- Patients choose a service when booking an appointment.  

---

## 4. Availability Model (`availabilities`)

The **Availability** model stores the **time slots** when a provider is available for appointments.  

### Fields
| Field         | Type      | Description |
|---------------|-----------|-------------|
| `id`          | BIGSERIAL PK | Unique availability ID. |
| `provider_id` | BIGINT FK → providers.id | Provider who is available. |
| `start_time`  | TIMESTAMP | Start of availability window. |
| `end_time`    | TIMESTAMP | End of availability window. |
| `is_recurring` | BOOLEAN | If `true`, this availability repeats (e.g., every Monday). |
| `created_at`  | TIMESTAMP | When this availability was created. |
| `updated_at`  | TIMESTAMP | When this availability was last updated. |

### Notes
- Prevent overlapping slots for the same provider.  
- Timezones must be respected (using provider’s timezone).  

---

## 5. Booking Model (`bookings`)

The **Booking** model represents an appointment between a **patient** and a **provider** for a specific service.  

### Fields
| Field         | Type      | Description |
|---------------|-----------|-------------|
| `id`          | BIGSERIAL PK | Unique booking ID. |
| `patient_id`  | BIGINT FK → users.id | The patient making the booking. |
| `provider_id` | BIGINT FK → providers.id | The doctor assigned. |
| `service_id`  | BIGINT FK → services.id | Which service is being booked. |
| `start_time`  | TIMESTAMP | Appointment start time. |
| `end_time`    | TIMESTAMP | Appointment end time. |
| `status`      | ENUM('pending','confirmed','completed','canceled') | Status of the booking. |
| `created_at`  | TIMESTAMP | When the booking was made. |
| `updated_at`  | TIMESTAMP | Last update of booking. |

### Notes
- The system must check availability before creating a booking.  
- A patient cannot book overlapping appointments with the same provider.  

---

## 6. Notification Model (`notifications`)

The **Notification** model stores reminders and messages sent to users.  
Notifications are generated for actions like booking confirmation, cancellation, or reminders.  

### Fields
| Field         | Type      | Description |
|---------------|-----------|-------------|
| `id`          | BIGSERIAL PK | Unique notification ID. |
| `user_id`     | BIGINT FK → users.id | The recipient of the notification. |
| `message`     | TEXT | The notification content. |
| `type`        | ENUM('email','sms') | Notification type. |
| `status`      | ENUM('pending','sent','failed') | Delivery status. |
| `created_at`  | TIMESTAMP | When the notification was created. |
| `sent_at`     | TIMESTAMP NULL | When it was actually sent. |

### Notes
- Notifications can be queued and sent asynchronously.  
- Failed notifications should be retried.  

---

## Relationships Summary

- **User** (base entity)
  - Can be a **Provider** (doctor profile)  
  - Or a **Patient** (just a user with `role=patient`)  
  - Linked to **Notifications**

- **Provider** (doctor profile)
  - Offers multiple **Services**  
  - Has multiple **Availabilities**  
  - Accepts multiple **Bookings**

- **Service**
  - Linked to a **Provider**  
  - Required when creating a **Booking**

- **Availability**
  - Linked to a **Provider**  
  - Defines possible booking slots

- **Booking**
  - Connects **Patient (user)** + **Provider (doctor)** + **Service**

- **Notification**
  - Sent to **Users** when important events happen (e.g., booking confirmed, canceled)

---
