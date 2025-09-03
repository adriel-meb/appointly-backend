# API Endpoints

## Authentication
| Endpoint        | Method | Related Table | Purpose          |
|-----------------|--------|---------------|-----------------|
| /auth/register  | POST   | USERS         | Register new user |
| /auth/login     | POST   | USERS         | Authenticate user |

**Example Registration Payload:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "role": "patient"
}

| Endpoint        | Method | Related Table                       | Purpose                      |
| --------------- | ------ | ----------------------------------- | ---------------------------- |
| /providers      | POST   | PROVIDERS                           | Create provider (admin only) |
| /providers      | GET    | PROVIDERS                           | List all providers           |
| /providers/{id} | GET    | PROVIDERS, SERVICES, AVAILABILITIES | Get provider details         |


| Endpoint                 | Method | Related Table | Purpose       |
| ------------------------ | ------ | ------------- | ------------- |
| /providers/{id}/services | POST   | SERVICES      | Add service   |
| /providers/{id}/services | GET    | SERVICES      | List services |
{
  "title": "Consultation",
  "description": "General checkup",
  "duration_minutes": 30,
  "price": 50.00
}


| Endpoint                     | Method | Related Table  | Purpose                 |
| ---------------------------- | ------ | -------------- | ----------------------- |
| /providers/{id}/availability | POST   | AVAILABILITIES | Add availability slots  |
| /providers/{id}/availability | GET    | AVAILABILITIES | List availability slots |


| Endpoint              | Method | Related Table            | Purpose              |
| --------------------- | ------ | ------------------------ | -------------------- |
| /bookings             | POST   | BOOKINGS, AVAILABILITIES | Create appointment   |
| /bookings/{id}        | GET    | BOOKINGS                 | View booking details |
| /bookings/{id}/cancel | PUT    | BOOKINGS, NOTIFICATIONS  | Cancel appointment   |
{
  "patient_id": 1,
  "provider_id": 2,
  "service_id": 5,
  "start_time": "2025-09-05T09:30:00Z",
  "end_time": "2025-09-05T10:00:00Z"
}


---

### **`docs/WORKFLOWS.md`**
```markdown
# Example Workflows

## Patient Booking Flow
1. **Admin creates a provider**
   - Add entry to `USERS` and `PROVIDERS` tables.

2. **Provider adds services**
   - Insert into `SERVICES` table.

3. **Provider sets availability**
   - Insert into `AVAILABILITIES` table.

4. **Patient browses providers**
   - `GET /providers` and `GET /providers/{id}/services`.

5. **Patient books an appointment**
   - `POST /bookings`.

6. **System sends notifications**
   - Insert entry into `NOTIFICATIONS` table.

7. **Patient cancels appointment**
   - `PUT /bookings/{id}/cancel`.

---

## Notifications Workflow
- Booking creation triggers:
  - Confirmation email/SMS
- 24 hours before appointment:
  - Reminder notification
- On cancellation:
  - Cancellation email/SMS


