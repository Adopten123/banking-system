-- name: CreatePost :one
INSERT INTO posts (author_id, type_id, content, media_attachments, related_asset_ticker, status)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetPostByID :one
SELECT *
FROM posts
WHERE id = $1
  AND deleted_at IS NULL LIMIT 1;

-- name: ListFeedPosts :many
SELECT p.id,
       p.content,
       p.media_attachments,
       p.related_asset_ticker,
       p.status,
       p.likes_count,
       p.comments_count,
       p.is_pinned,
       p.created_at,
       u.user_id    AS author_id,
       u.username   AS author_username,
       u.avatar_url AS author_avatar
FROM posts p
         JOIN users_profile u ON p.author_id = u.user_id
WHERE p.deleted_at IS NULL
  AND p.status = 'published'
ORDER BY p.is_pinned DESC, p.created_at DESC LIMIT $1
OFFSET $2;

-- name: IncrementPostLikes :exec
UPDATE posts
SET likes_count = likes_count + 1
WHERE id = $1;

-- name: DecrementPostLikes :exec
UPDATE posts
SET likes_count = GREATEST(likes_count - 1, 0)
WHERE id = $1;

-- name: AddPostLike :exec
INSERT INTO post_likes (post_id, user_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING;

-- name: RemovePostLike :exec
DELETE
FROM post_likes
WHERE post_id = $1
  AND user_id = $2;