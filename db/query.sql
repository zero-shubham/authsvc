-- name: CreateAppGroup :one
INSERT INTO app_groups (org_id, name, scopes)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateApp :one
INSERT INTO apps (org_id, app_grp_id)
VALUES ($1, $2)
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (email, password, app_group_id, org_id)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: GetUserAndAppGrpByEmail :one
SELECT users.*, app_groups.scopes
FROM users
JOIN app_groups ON users.app_group_id = app_groups.id WHERE email = $1;

-- name: UpdateAppGroup :one
UPDATE app_groups
SET name = $2,
    scopes = $3,
    org_id = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING id, org_id, name, scopes, created_at, updated_at;

-- name: GetAppGroupByID :one
SELECT id, org_id, name, scopes, created_at, updated_at
FROM app_groups
WHERE id = $1; 

-- name: GetAppGrpByAppID :one
SELECT * from app_groups JOIN apps on app_groups.id = apps.app_grp_id WHERE apps.id = $1;

-- name: GetAppGroupByName :one
SELECT id, org_id, name, scopes, created_at, updated_at
FROM app_groups
WHERE name = $1; 

-- name: DbCleanUp :exec
DROP TABLE IF EXISTS apps, app_groups, users CASCADE;
