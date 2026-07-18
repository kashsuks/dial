CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task TEXT NOT NULL,
    project TEXT,
    tags TEXT,
    started_at DATETIME NOT NULL,
    ended_at DATETIME,
    paused_at DATETIME,
    paused_seconds INTEGER DEFAULT 0,
    source TEXT DEFAULT 'manual',
    note TEXT
);

CREATE INDEX IF NOT EXISTS idx_sessions_started ON sessions(started_at);
CREATE INDEX IF NOT EXISTS idx_sessions_running ON sessions(ended_at) WHERE ended_at IS NULL;
