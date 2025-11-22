# Cocktails API

A RESTful API for managing cocktails, ingredients, and users built with TypeScript and Express.

---

## Overview

Node.js backend application with TypeScript that manages cocktails, ingredients, and users. Uses Prisma ORM for database management with PostgreSQL and includes comprehensive testing with Jest.

Created for https://github.com/Solvro/rekrutacja/blob/main/backend.md

---

## Features

- RESTful API endpoints for cocktails, ingredients, and user management
- JWT authentication and authorization
- Prisma ORM with PostgreSQL database
- Image upload with multer
- Comprehensive test suite (unit + integration)
- TypeScript for type safety

---

## Project Structure

```
cocktails-api/
├── src/
│   ├── handlers/          # Route handlers
│   ├── modules/           # Auth, middleware, types
│   ├── db.ts             # Prisma client
│   ├── router.ts         # Route definitions
│   └── server.ts         # Express setup
├── tests/
│   ├── integration/      # API endpoint tests
│   └── unit/             # Module unit tests
├── prisma/
│   ├── schema.prisma     # Database schema
│   └── migrations/       # Schema migrations
└── package.json
```

---

## Installation

Install dependencies:
```bash
npm install
```

---

## Configuration

Copy example environment file:
```bash
cp example.env .env
```

Configure database connection and other settings in `.env`

---

## Running

Development server:
```bash
npm run dev
```

Build for production:
```bash
npm run build
```

---

## Testing

Run all tests:
```bash
npm test
```

Watch mode:
```bash
npm run test-watch
```

---

## License

MIT License
