CREATE TABLE status (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

INSERT INTO status (name) VALUES ('pending'), ('done');

CREATE TABLE IF NOT EXISTS schema_migrations (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    status INT NOT NULL REFERENCES status(id)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);
