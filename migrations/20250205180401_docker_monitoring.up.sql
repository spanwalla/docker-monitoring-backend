CREATE TABLE pingers(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pingers_name_hash ON pingers USING HASH (name);

CREATE TABLE reports(
    id SERIAL PRIMARY KEY,
    pinger_id INT NOT NULL REFERENCES pingers(id),
    content JSONB NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reports_pinger_id_created_at ON reports (pinger_id, created_at DESC);