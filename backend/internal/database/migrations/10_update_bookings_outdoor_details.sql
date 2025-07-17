-- Update existing records to have default values
UPDATE bookings
SET is_outdoor = FALSE, has_shade = FALSE
WHERE is_outdoor IS NULL OR has_shade IS NULL;