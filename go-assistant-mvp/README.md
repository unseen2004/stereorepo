# go-assistant-mvp

An AI-powered personal assistant backend written in Go. It connects your Gmail, Google Calendar, and Notion into a single API gateway — letting an LLM (Gemini) read context, schedule tasks, send notifications, and act on your behalf.

> **Status:** MVP — core integrations working, Kubernetes + Terraform deploy configs included.

---

## What it does

- **AI gateway** — A single HTTP service (`cmd/gateway`) that routes natural-language requests to Gemini and dispatches actions
- **Gmail** — reads and processes incoming emails as task triggers
- **Google Calendar** — creates and queries calendar events
- **Notion** — syncs tasks to/from a Notion database
- **Notifications** — internal notification dispatch layer
- **Location awareness** — optional location context for tasks
- **Persistent state** — PostgreSQL for tasks, Redis for caching/queuing

---

## Stack

| Layer | Tech |
|---|---|
| Language | Go 1.22+ |
| AI | Google Gemini API |
| Integrations | Gmail API, Google Calendar API, Notion API |
| Storage | PostgreSQL + Redis |
| Deploy | Docker Compose / Kubernetes / Terraform |

---

## Quick start

**Prerequisites:** Go 1.22+, Docker

```bash
git clone https://github.com/unseen2004/go-assistant-mvp.git
cd go-assistant-mvp

# Copy and fill in your credentials
cp .env.example .env

# Start Postgres + Redis + app
make docker-up

# Or run locally (Postgres/Redis must already be running)
make run
```

Check it's alive:
```bash
make health   # curl http://localhost:8080/health
```

---

## Environment variables

```env
GEMINI_API_KEY=
NOTION_API_KEY=
NOTION_DATABASE_ID=
GMAIL_CLIENT_ID=
GMAIL_CLIENT_SECRET=
GOOGLE_CALENDAR_ID=
POSTGRES_URL=postgres://goassistant:secret@localhost:5433/goassistant?sslmode=disable
REDIS_URL=redis://localhost:6379
```

See `.env.example` for the full template.

---

## Project layout

```
cmd/gateway/        # HTTP entry point
internal/
  ai/               # Gemini client & prompt handling
  integrations/
    gmail/          # Gmail OAuth + message reading
    gcal/           # Google Calendar API
    notion/         # Notion database sync
  tasks/            # Task storage & scheduling
  notifications/    # Notification dispatch
  location/         # Location context
deploy/
  k8s/              # Kubernetes manifests
  terraform/        # Infrastructure as code
scripts/            # e2e test script
```

---

## Make targets

```
make build          Build binary → bin/gateway
make run            Run locally
make test           Run unit tests
make docker-up      Start with Docker Compose
make docker-down    Stop containers
make e2e-test       Run end-to-end test scenario
make k8s-deploy     Deploy to Kubernetes
make tf-apply       Apply Terraform infrastructure
```

---

## License

MIT
