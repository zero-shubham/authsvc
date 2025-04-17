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