-- name: CreateUserProfile :one
INSERT INTO users_profile (user_id, username, display_name, avatar_url, bio, is_verified, is_staff, is_private)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetUserProfile :one
SELECT *
FROM users_profile
WHERE user_id = $1 LIMIT 1;

-- name: UpdateUserProfile :one
UPDATE users_profile
SET username     = COALESCE(sqlc.narg(username), username),
    display_name = COALESCE(sqlc.narg(display_name), display_name),
    avatar_url   = COALESCE(sqlc.narg(avatar_url), avatar_url),
    bio          = COALESCE(sqlc.narg(bio), bio),
    is_private   = COALESCE(sqlc.narg(is_private), is_private),
    updated_at   = NOW()
WHERE user_id = $1 RETURNING *;

-- name: FollowUser :exec
INSERT INTO follows (follower_id, following_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: UnfollowUser :exec
DELETE
FROM follows
WHERE follower_id = $1
  AND following_id = $2;