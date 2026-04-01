CREATE TABLE IF NOT EXISTS jobs (
    id          VARCHAR(36) PRIMARY KEY,
    queue       VARCHAR(100) NOT NULL DEFAULT 'default',
    payload     TEXT NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'pending',
    attempts    INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 3,
    error       TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_queue  ON jobs(queue);