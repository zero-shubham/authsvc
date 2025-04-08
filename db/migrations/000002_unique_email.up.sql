
-- Migration: Add unique constraint to email in users table
-- Up migration
ALTER TABLE users ADD CONSTRAINT unique_email UNIQUE (email);