# Go Tasker Backend

REST API for a task management system built in Go. Features todos with subtasks, categories, comments, file attachments, background email notifications, and automated cron jobs.

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23+ |
| Web framework | Echo v4 |
| Database | PostgreSQL (pgx v5, Tern migrations) |
| Job queue | Redis + Asynq |
| Auth | Clerk |
| Email | Resend |
| File storage | AWS S3 |
| Observability | New Relic APM, ZeroLog |
| Config | Koanf |

## Project Structure

```
backend/
├── cmd/
│   ├── tasker/               # HTTP API server + background job worker
│   └── cron/                 # CLI for running scheduled jobs
├── internal/
│   ├── config/               # Config loading (Koanf, env vars)
│   ├── server/               # Server init & lifecycle
│   ├── router/
│   │   ├── router.go         # Root router, middleware, system routes
│   │   └── v1/               # v1 API route registration
│   ├── handler/              # HTTP request handlers
│   ├── service/              # Business logic
│   ├── repository/           # Database queries
│   ├── model/                # Domain models (todo, category, comment)
│   ├── middleware/           # Auth, rate limiting, tracing, request ID, CORS
│   ├── cron/                 # Scheduled job definitions & registry
│   ├── database/
│   │   └── migrations/       # SQL migration files (Tern)
│   ├── lib/
│   │   ├── aws/              # S3 client
│   │   ├── email/            # Email templating & Resend client
│   │   └── job/              # Asynq task definitions & handlers
│   ├── logger/               # Structured logging with New Relic integration
│   ├── sqlerr/               # SQL error → domain error mapping
│   └── validation/           # Input validation helpers
└── templates/
    └── emails/               # HTML email templates
```

## Prerequisites

- Go 1.23+
- PostgreSQL
- Redis
- [Clerk](https://clerk.com) account (authentication)
- [Resend](https://resend.com) API key (email)
- AWS S3 bucket (file attachments)
- [Task](https://taskfile.dev) (optional, for Taskfile commands)

## Configuration

Copy `.env.sample` to `.env` and fill in the values.

```bash
cp .env.sample .env
```

### Environment Variables

**Primary**
```
TASKER_PRIMARY.ENV                                    # local | development | production
```

**Server**
```
TASKER_SERVER.PORT                                    # HTTP port (default: 8080)
TASKER_SERVER.READ_TIMEOUT                            # seconds (default: 30)
TASKER_SERVER.WRITE_TIMEOUT                           # seconds (default: 30)
TASKER_SERVER.IDLE_TIMEOUT                            # seconds (default: 60)
TASKER_SERVER.CORS_ALLOWED_ORIGINS                    # comma-separated origins
```

**Database**
```
TASKER_DATABASE.HOST
TASKER_DATABASE.PORT                                  # default: 5432
TASKER_DATABASE.USER
TASKER_DATABASE.PASSWORD
TASKER_DATABASE.NAME
TASKER_DATABASE.SSL_MODE                              # disable | require | prefer
TASKER_DATABASE.MAX_OPEN_CONNS                        # default: 25
TASKER_DATABASE.MAX_IDLE_CONNS                        # default: 25
TASKER_DATABASE.CONN_MAX_LIFETIME                     # seconds (default: 300)
TASKER_DATABASE.CONN_MAX_IDLE_TIME                    # seconds (default: 300)
```

**Auth**
```
TASKER_AUTH.SECRET_KEY
```

**Integrations**
```
TASKER_INTEGRATION.RESEND_API_KEY
```

**Redis**
```
TASKER_REDIS.ADDRESS                                  # e.g. localhost:6379
```

**AWS**
```
TASKER_AWS.ACCESS_KEY_ID
TASKER_AWS.SECRET_ACCESS_KEY
TASKER_AWS.REGION
TASKER_AWS.UPLOAD_BUCKET
TASKER_AWS.ENDPOINT_URL                               # optional: custom S3-compatible endpoint
```

**Observability**
```
TASKER_OBSERVABILITY.SERVICE_NAME
TASKER_OBSERVABILITY.ENVIRONMENT
TASKER_OBSERVABILITY.LOGGING.LEVEL                    # debug | info | warn | error
TASKER_OBSERVABILITY.LOGGING.FORMAT                   # console | json
TASKER_OBSERVABILITY.LOGGING.SLOW_QUERY_THRESHOLD     # e.g. 100ms
TASKER_OBSERVABILITY.NEW_RELIC.LICENSE_KEY
TASKER_OBSERVABILITY.NEW_RELIC.APP_LOG_FORWARDING_ENABLED
TASKER_OBSERVABILITY.NEW_RELIC.DISTRIBUTED_TRACING_ENABLED
TASKER_OBSERVABILITY.NEW_RELIC.DEBUG_LOGGING
TASKER_OBSERVABILITY.HEALTH_CHECKS.ENABLED
TASKER_OBSERVABILITY.HEALTH_CHECKS.INTERVAL           # e.g. 30s
TASKER_OBSERVABILITY.HEALTH_CHECKS.TIMEOUT            # e.g. 5s
TASKER_OBSERVABILITY.HEALTH_CHECKS.CHECKS             # e.g. database,redis
```

**Cron**
```
TASKER_CRON.ARCHIVE_DAYS_THRESHOLD                    # days before archiving completed todos (default: 30)
TASKER_CRON.BATCH_SIZE                                # records per cron batch (default: 100)
TASKER_CRON.REMINDER_HOURS                            # hours before due date to send reminders (default: 24)
TASKER_CRON.MAX_TODOS_PER_USER_NOTIFICATION           # max todos per notification email (default: 10)
```

## Running

### API Server

```bash
go run ./cmd/tasker
# or
task run
```

> On non-local environments, database migrations run automatically on startup. In `local`, run them manually (see below).

### Database Migrations

```bash
# Apply all pending migrations
TASKER_DB_DSN="postgres://user:pass@localhost:5432/tasker" task migrations:up

# Create a new migration file
task migrations:new name=<migration_name>
```

## API Endpoints

All `/api/v1/*` routes require a valid Clerk session token.

### System

| Method | Path | Description |
|--------|------|-------------|
| GET | `/status` | Health check |
| GET | `/docs` | OpenAPI documentation UI |

### Todos

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/todos` | Create a todo |
| GET | `/api/v1/todos` | List todos (paginated, filterable) |
| GET | `/api/v1/todos/stats` | Todo statistics for the current user |
| GET | `/api/v1/todos/:id` | Get a todo by ID |
| PATCH | `/api/v1/todos/:id` | Update a todo |
| DELETE | `/api/v1/todos/:id` | Delete a todo |

### Attachments

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/todos/:id/attachments` | Upload a file attachment (S3) |
| DELETE | `/api/v1/todos/:id/attachments/:attachmentId` | Delete an attachment |
| GET | `/api/v1/todos/:id/attachments/:attachmentId/download` | Get a presigned S3 download URL |

### Comments

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/todos/:id/comments` | Add a comment to a todo |
| GET | `/api/v1/todos/:id/comments` | List comments on a todo |
| PATCH | `/api/v1/comments/:id` | Update a comment |
| DELETE | `/api/v1/comments/:id` | Delete a comment |

### Categories

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/categories` | Create a category |
| GET | `/api/v1/categories` | List categories |
| PATCH | `/api/v1/categories/:id` | Update a category |
| DELETE | `/api/v1/categories/:id` | Delete a category |

## Background Jobs

The API server (`cmd/tasker`) also runs the Asynq worker that processes queued tasks. Cron jobs (`cmd/cron`) only **enqueue** tasks to Redis — the server must be running for emails and other tasks to be processed.

### Cron Jobs

| Job | Description |
|-----|-------------|
| `due-date-reminders` | Enqueues reminder emails for todos due within the configured hours |
| `overdue-notifications` | Enqueues notification emails for overdue todos |
| `weekly-reports` | Enqueues weekly productivity report emails per user |
| `auto-archive` | Archives completed todos older than the configured threshold |

```bash
# List all available jobs
go run ./cmd/cron list

# Run a specific job
go run ./cmd/cron due-date-reminders
go run ./cmd/cron overdue-notifications
go run ./cmd/cron weekly-reports
go run ./cmd/cron auto-archive
```

### Job Queues

Tasks are processed with three priority queues: `critical` (60%), `default` (30%), `low` (10%).

## Development Tasks

```bash
task run                        # start the API server
task migrations:up              # apply pending migrations
task migrations:new name=<name> # create a new migration
task tidy                       # go fmt + go mod tidy + go mod verify
task help                       # list all available tasks
```
