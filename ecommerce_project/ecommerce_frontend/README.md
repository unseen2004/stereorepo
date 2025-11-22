# E-commerce Frontend

React-based frontend for e-commerce platform with Vite build system.

---

## Overview

Modern React application providing user interface for the e-commerce platform. Features product browsing, cart management, and checkout functionality with responsive design.

---

## Building

Install dependencies:
```bash
npm install
```

---

## Configuration

Create `.env` file with API endpoint:
```
VITE_API_URL=http://localhost:3000
```

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

Preview production build:
```bash
npm run preview
```

---

## Linting

Run ESLint:
```bash
npm run lint
```

---

## Features

- React 18 with hooks
- React Router for navigation
- Axios for API communication
- React Hook Form for form handling
- React Hot Toast for notifications
- Responsive design
- Fast refresh with Vite HMR

---

## Project Structure

```
ecommerce_frontend/
├── src/
│   ├── components/   # React components
│   ├── pages/        # Page components
│   ├── services/     # API service layer
│   └── App.jsx       # Root component
├── public/           # Static assets
└── index.html        # HTML entry point
```

---

## Tech Stack

- React 18
- Vite
- React Router DOM
- Axios
- React Hook Form
- React Hot Toast
- ESLint

---

## Development

Vite provides fast HMR and optimized builds. The dev server runs on `http://localhost:5173` by default.
