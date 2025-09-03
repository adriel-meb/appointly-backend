# Doctor Appointment Backend - Introduction

## Overview
This project provides a **scalable backend system** for a doctor appointment booking application.  
It is built with **Go** and **PostgreSQL**, focusing solely on backend functionalities.

### Core Features
- User authentication (patients, providers, admins)
- Provider management
- Services and pricing management
- Availability scheduling
- Appointment booking
- Notifications via email/SMS
- Optional payment integration

### Architecture
- REST API with **Gin**
- PostgreSQL with **GORM** ORM
- JWT-based authentication
- Background workers for notifications
- Optional Docker/Kubernetes for deployment
