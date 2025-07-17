BEGIN;

-- Add outdoor and shade fields to bookings table
ALTER TABLE public.bookings 
ADD COLUMN IF NOT EXISTS is_outdoor BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS has_shade BOOLEAN DEFAULT FALSE;

-- Force an error if columns don't exist to catch silent failures
DO $$
BEGIN
  PERFORM is_outdoor FROM public.bookings LIMIT 1;
  EXCEPTION WHEN undefined_column THEN
    RAISE EXCEPTION 'Column is_outdoor was not created successfully';
END $$;

COMMIT;