CREATE TABLE
    IF NOT EXISTS clients (
        id SERIAL PRIMARY KEY,
        client_id VARCHAR(100) NOT NULL UNIQUE,
        algorithm VARCHAR(20) NOT NULL DEFAULT 'fixed_window',
        "limit" INTEGER NOT NULL DEFAULT 10,
        window_seconds INTEGER NOT NULL DEFAULT 60,
        api_key VARCHAR(255) NOT NULL UNIQUE,
        deleted_at TIMESTAMPTZ
    );

CREATE INDEX IF NOT EXISTS idx_clients_deleted_at ON clients (deleted_at);
