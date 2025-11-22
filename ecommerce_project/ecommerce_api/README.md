# E-commerce API

RESTful API backend for e-commerce platform built with Node.js and Express.

---

## Overview

Backend API providing authentication, product management, shopping cart, and order processing functionality. Uses PostgreSQL database with Sequelize ORM.

---

## Building

Install dependencies:
```bash
npm install
```

---

## Configuration

Create `.env` file with configuration:
```bash
cp .env.test .env
```

Configure database connection and JWT secret in `.env`

---

## Running

Development mode:
```bash
npm run dev
```

Production mode:
```bash
npm start
```

---

## Testing

Run test suite:
```bash
npm test
```

Watch mode:
```bash
npm run test:watch
```

---

## Features

- User authentication with JWT
- Password hashing with bcryptjs
- Product CRUD operations
- Shopping cart management
- Order processing
- Input validation with express-validator
- CORS support
- PostgreSQL database with Sequelize

---

## API Structure

```
ecommerce_api/
├── config/          # Database and app configuration
├── controllers/     # Route controllers
├── middleware/      # Authentication and validation
├── models/          # Sequelize models
├── routes/          # API route definitions
├── tests/           # Test suite
├── app.js          # Express app setup
└── server.js       # Server entry point
```

---

## Tech Stack

- Node.js
- Express
- PostgreSQL
- Sequelize ORM
- JWT authentication
- Jest for testing
