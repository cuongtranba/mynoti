CREATE TABLE comic_tracking (
    id          SERIAL PRIMARY KEY,
    url         TEXT NOT NULL,
    name        VARCHAR(255),
    description TEXT,
    html        TEXT,
    cron_spec   VARCHAR(100),
    last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);