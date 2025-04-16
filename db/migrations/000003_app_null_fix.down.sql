-- Migration: make org_id nullable and app_grp_id not nullable in apps table
-- Down migration
ALTER TABLE apps ALTER COLUMN org_id SET NOT NULL;
ALTER TABLE apps ALTER COLUMN app_grp_id DROP NOT NULL;