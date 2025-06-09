-- Only add default items if table is empty
DO $$ 
BEGIN
    -- Check if menu_items table has any rows
    IF (SELECT COUNT(*) FROM menu_items) = 0 THEN
        -- Insert coffee flavors
        INSERT INTO menu_items (value, label, type) VALUES
        ('french_toast', 'French Toast', 'coffee_flavor'),
        ('dirty_vanilla_chai', 'Dirty Vanilla Chai', 'coffee_flavor'),
        ('mexican_mocha', 'Mexican Mocha', 'coffee_flavor'),
        ('cinnamon_brown_sugar', 'Cinnamon Brown Sugar', 'coffee_flavor'),
        ('horchata', 'Horchata (made w/ rice milk)', 'coffee_flavor');

        -- Insert milk options
        INSERT INTO menu_items (value, label, type) VALUES
        ('whole', 'Whole Milk', 'milk_option'),
        ('half_and_half', 'Half & Half', 'milk_option'),
        ('oat', 'Oat Milk', 'milk_option'),
        ('almond', 'Almond Milk', 'milk_option'),
        ('rice', 'Rice Milk', 'milk_option');
    END IF;
END $$;