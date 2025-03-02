CREATE TABLE IF NOT EXISTS status (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

INSERT OR IGNORE INTO status (name) VALUES ('pending'), ('done');

CREATE TABLE IF NOT EXISTS tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    status INT NOT NULL REFERENCES status(id),
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);