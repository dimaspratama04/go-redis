SET @exist := (
  SELECT COUNT(*)
  FROM INFORMATION_SCHEMA.COLUMNS
  WHERE TABLE_NAME = 'products'
    AND COLUMN_NAME = 'category'
    AND TABLE_SCHEMA = DATABASE()
);

SET @stmt := IF(@exist = 0,
  'ALTER TABLE products ADD COLUMN category VARCHAR(100) NOT NULL DEFAULT "general";',
  'SELECT "Column already exists" as info;'
);

PREPARE stmt FROM @stmt;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- Update data existing
UPDATE products
SET category = 'electronics'
WHERE category IS NULL OR category = '';
