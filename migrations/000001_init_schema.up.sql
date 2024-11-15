CREATE TABLE comic_tracking (
    id          SERIAL PRIMARY KEY,
    url         TEXT NOT NULL,
    name        VARCHAR(255),
    description TEXT,
    html        TEXT,
    last_checked TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);