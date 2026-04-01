# Distributed Task Queue

A job queue system built in Go that spawns workers, handles retries, and tracks job status. Built with Redis for queueing and PostgreSQL for persistent job storage.

## Tech Stack

- **Go** — core application
- **Redis** — job queue (LPUSH/BRPOP)
- **PostgreSQL** — persistent job storage and status tracking
- **Docker** — running Redis and PostgreSQL locally

## Project Structure

```
distributed-task-queue/
├── api/                  # HTTP handlers (POST /jobs, GET /jobs/:id)
├── db/
│   ├── db.go             # Postgres connection
│   └── jobs.go           # InsertJob, UpdateJobStatus, GetJob
├── migrations/
│   └── 001_create_jobs.sql
├── models/
│   └── job.go            # Job struct and status constants
├── queue/
│   └── queue.go          # Enqueue + SaveJob
├── worker/
│   └── worker.go         # StartWorker, job processing, retry logic
├── go.mod
├── go.sum
└── main.go
```

## How It Works

```
POST /jobs
    ↓
Enqueue()
    ├── Save full job to PostgreSQL  (status = pending)
    └── Push job ID to Redis list    (LPUSH queue:default)

Worker (running in background)
    ├── BRPOP queue:default          (blocks until job arrives)
    ├── Fetch full job from Postgres
    ├── Mark as running
    ├── Process job
    └── Mark as completed / failed
            └── If failed + attempts < max_retries → re-enqueue (RPUSH)
```

## Job Lifecycle

```
pending → running → completed
                 ↘ failed (if attempts >= max_retries)
                 ↘ pending (re-enqueued for retry)
```

## Prerequisites

- Go 1.21+
- Docker Desktop

## Setup

**1. Clone the repo**

```bash
git clone https://github.com/krutinewalkar/distributed-task-queue
cd distributed-task-queue
```

**2. Start PostgreSQL and Redis with Docker**

```bash
docker run --name postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_USER=taskqueue \
  -e POSTGRES_DB=taskqueue \
  -p 5432:5432 \
  -d postgres:latest

docker run --name redis \
  -p 6379:6379 \
  -d redis:latest
```

**3. Run the database migration**

```bash
docker exec -i postgres psql -U taskqueue -d taskqueue < migrations/001_create_jobs.sql
```

**4. Install dependencies**

```bash
go mod tidy
```

**5. Run**

```bash
go run main.go
```

You should see:

```
Postgres connected
Redis connected
Worker started for queue: default
```

## Job Schema

| Column      | Type         | Description                          |
|-------------|--------------|--------------------------------------|
| id          | VARCHAR(36)  | UUID, primary key                    |
| queue       | VARCHAR(100) | Queue name, default = "default"      |
| payload     | TEXT         | JSON string with job data            |
| status      | VARCHAR(20)  | pending / running / completed / failed |
| attempts    | INTEGER      | How many times the job has been tried |
| max_retries | INTEGER      | Max attempts before marking failed   |
| error       | TEXT         | Error message if job failed          |
| created_at  | TIMESTAMP    | When the job was created             |
| updated_at  | TIMESTAMP    | Last status update                   |

## Retry Logic

- Every job gets `max_retries = 3` by default
- On failure, `attempts` is incremented
- If `attempts < max_retries`, job is re-enqueued with `RPUSH` (goes to back of queue)
- If `attempts >= max_retries`, job is marked `failed` permanently

## Verifying Jobs

**Check all jobs in Postgres:**

```bash
docker exec -it postgres psql -U taskqueue -d taskqueue -c "SELECT id, status, attempts FROM jobs;"
```

**Check Redis queue:**

```bash
docker exec -it redis redis-cli LRANGE queue:default 0 -1
```

## Roadmap

- [ ] HTTP API — `POST /jobs` and `GET /jobs/:id`
- [ ] Multiple named queues
- [ ] Delayed job scheduling
- [ ] Worker concurrency (multiple goroutines)
- [ ] Dead letter queue for permanently failed jobs
