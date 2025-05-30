DO $$
BEGIN
    -- Check if the column already exists
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'bookings' AND column_name = 'archived'
    ) THEN
        -- Add the archived column with default value of FALSE
        ALTER TABLE bookings ADD COLUMN archived BOOLEAN NOT NULL DEFAULT FALSE;
        
        -- Set all existing bookings to not archived
        UPDATE bookings SET archived = FALSE;
        
        RAISE NOTICE 'Added archived column to bookings table';
    ELSE
        RAISE NOTICE 'Column archived already exists in bookings table';
    END IF;
END $$;