CREATE TABLE IF NOT EXISTS menu_items (
    id SERIAL PRIMARY KEY,
    value VARCHAR(100) NOT NULL,
    label VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);