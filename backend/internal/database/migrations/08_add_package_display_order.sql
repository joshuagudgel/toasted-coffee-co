ALTER TABLE packages ADD COLUMN display_order INTEGER DEFAULT 0;

-- Update any existing packages with sequential order
WITH ordered_packages AS (
  SELECT id, ROW_NUMBER() OVER (ORDER BY name) AS row_num
  FROM packages
)
UPDATE packages
SET display_order = op.row_num - 1
FROM ordered_packages op
WHERE packages.id = op.id;