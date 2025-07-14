# 🔐 auth-service

A lightweight, modular authentication microservice written in Go with PostgreSQL, designed for multi-tenant SaaS platforms.

---

## ✨ Features

- Superadmin bootstrapping via environment variables
- JWT-based authentication
- Role-based access control (Superadmin, Admin, Member, Guest)
- Multi-tenant user support (each user is scoped to a TenantID)
- Secure password hashing with bcrypt
- Pagination-ready endpoints
- Built with [Gin](https://github.com/gin-gonic/gin), [GORM](https://gorm.io/), and Go modules

---

## 🧱 Tech Stack

- Language: [Go 1.23+](https://golang.org)
- Framework: [Gin](https://github.com/gin-gonic/gin)
- ORM: [GORM](https://gorm.io)
- Database: PostgreSQL
- Auth: JWT, bcrypt

---

## 🗂 Project Structure

