-- Down migration (to revert if needed)
ALTER TABLE users DROP CONSTRAINT unique_email;