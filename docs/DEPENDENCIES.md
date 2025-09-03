
# Doctor Appointment Backend - Dependencies Guide

This document lists all the **dependencies** for the Doctor Appointment Backend project, explains their purpose, and includes instructions for installation.

---

## 1. Web Framework — Gin
- **Package:** `github.com/gin-gonic/gin`
- **Purpose:** Handles routing, HTTP requests, and middleware.
- **Why?** Fast, minimal web framework with built-in JSON parsing, routing groups, middleware, and error handling.

**Install:**
```bash
go get github.com/gin-gonic/gin

## 2. ORM — GORM

Package: gorm.io/gorm

Purpose: Object Relational Mapping for PostgreSQL.

Why? Simplifies database operations, migrations, and managing relationships.

## 4. JWT Authentication

Package: github.com/golang-jwt/jwt/v5

Purpose: Generate and verify JSON Web Tokens for authentication.


## 5. Password Hashing

Package: golang.org/x/crypto

Purpose: Securely hash and verify passwords.


## 6. Configuration Management — Viper

Package: github.com/spf13/viper

Purpose: Load config files (YAML, JSON, ENV) and environment variables.


9. Logger (Optional but Recommended)

Package: go.uber.org/zap

Purpose: Structured, production-grade logging.


## 10. Testing — Testify

Package: github.com/stretchr/testify

Purpose: Unit testing and assertions.




## 11. Migrations — Goose or Golang-Migrate

Purpose: Manage database schema changes across environments.



