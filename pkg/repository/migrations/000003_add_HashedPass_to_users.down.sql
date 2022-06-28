BEGIN;
  ALTER TABLE user 
  DROP COLUMN hashedPassword; 
COMMIT;