package db

const SchemaSQL = `
CREATE TABLE IF NOT EXISTS spikes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    detected_at DATETIME DEFAULT (datetime('now')),
    pid INTEGER NOT NULL,
    name TEXT NOT NULL,
    cpu_percent REAL,
    memory_rss INTEGER,
    memory_vms INTEGER,
    duration_sec INTEGER DEFAULT 0,
    reason TEXT NOT NULL      
);

CREATE INDEX IF NOT EXISTS idx_pid ON spikes(pid);
CREATE INDEX IF NOT EXISTS idx_detected_at ON spikes(detected_at);
CREATE INDEX IF NOT EXISTS idx_reason ON spikes(reason);

CREATE TABLE IF NOT EXISTS recording_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    started_at DATETIME DEFAULT (datetime('now')),
    ended_at DATETIME,
    cpu_threshold REAL NOT NULL,
    ram_threshold REAL NOT NULL,
    duration_sec INTEGER NOT NULL,
    status TEXT DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS recorded_processes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    recorded_at DATETIME DEFAULT (datetime('now')),
    pid INTEGER NOT NULL,
    name TEXT NOT NULL,
    cpu_percent REAL,
    memory_percent REAL,
    memory_rss INTEGER,
    exe TEXT,
    cmdline TEXT,
    username TEXT,
    FOREIGN KEY (session_id) REFERENCES recording_sessions(id)
);

CREATE INDEX IF NOT EXISTS idx_recorded_session ON recorded_processes(session_id);
CREATE INDEX IF NOT EXISTS idx_recorded_at ON recorded_processes(recorded_at);
CREATE INDEX IF NOT EXISTS idx_recorded_pid ON recorded_processes(pid);

CREATE TABLE IF NOT EXISTS metrics_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    timestamp DATETIME DEFAULT (datetime('now')),
    cpu_percent REAL NOT NULL,
    memory_percent REAL NOT NULL,
    memory_used_mb INTEGER NOT NULL,
    memory_total_mb INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics_history(timestamp);

CREATE TABLE IF NOT EXISTS alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT (datetime('now')),
    type TEXT NOT NULL,
    threshold REAL NOT NULL,
    current_value REAL NOT NULL,
    message TEXT NOT NULL,
    acknowledged INTEGER DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON alerts(created_at);
CREATE INDEX IF NOT EXISTS idx_alerts_type ON alerts(type);
`
