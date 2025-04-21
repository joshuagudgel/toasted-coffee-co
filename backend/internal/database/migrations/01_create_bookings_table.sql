CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    date DATE NOT NULL,
    time VARCHAR(10) NOT NULL,
    people INTEGER NOT NULL,
    location VARCHAR(255) NOT NULL,
    notes TEXT,
    coffee_flavors VARCHAR[] NOT NULL,
    milk_options VARCHAR[] NOT NULL,
    package VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);