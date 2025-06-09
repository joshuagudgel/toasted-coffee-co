CREATE TABLE IF NOT EXISTS menu_items (
    id SERIAL PRIMARY KEY,
    value VARCHAR(100) NOT NULL,
    label VARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default menu items
INSERT INTO menu_items (value, label, type) VALUES
('french_toast', 'French Toast', 'coffee_flavor'),
('dirty_vanilla_chai', 'Dirty Vanilla Chai', 'coffee_flavor'),
('mexican_mocha', 'Mexican Mocha', 'coffee_flavor'),
('cinnamon_brown_sugar', 'Cinnamon Brown Sugar', 'coffee_flavor'),
('horchata', 'Horchata (made w/ rice milk)', 'coffee_flavor');

INSERT INTO menu_items (value, label, type) VALUES
('whole', 'Whole Milk', 'milk_option'),
('half_and_half', 'Half & Half', 'milk_option'),
('oat', 'Oat Milk', 'milk_option'),
('almond', 'Almond Milk', 'milk_option'),
('rice', 'Rice Milk', 'milk_option');