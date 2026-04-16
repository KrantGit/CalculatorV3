CREATE TABLE IF NOT EXISTS calculations (
                                            id SERIAL PRIMARY KEY,
                                            expression TEXT NOT NULL,
                                            result NUMERIC NOT NULL,
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_calculations_created_at ON calculations(created_at);