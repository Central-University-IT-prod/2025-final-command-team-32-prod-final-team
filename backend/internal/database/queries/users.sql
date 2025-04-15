-- name: UserExists :one
SELECT EXISTS(SELECT(1) FROM users WHERE login = $1);

-- name: GetUserByLogin :one
SELECT * FROM users WHERE login = $1;

-- name: CreateUser :exec
INSERT INTO users(login, password, embedding) VALUES ($1, $2, $3);

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetVector :one
SELECT embedding FROM users WHERE login = $1;

-- name: UploadVector :exec
UPDATE users
SET embedding = $1
WHERE id = $2;

-- name: SaveFilmToFavorites :exec
INSERT INTO saved (subject_id, cinema_id) VALUES ($1, $2);

-- name: DeleteFilmFromFavorites :exec
DELETE FROM saved WHERE subject_id = $1 AND cinema_id = $2;

-- name: CheckIfFavorite :one
SELECT EXISTS(SELECT(1) FROM saved WHERE subject_id = $1 AND cinema_id = $2);

-- name: SaveFilmToBlacklist :exec
INSERT INTO blacklisted (subject_id, cinema_id) VALUES ($1, $2);

-- name: DeleteFilmFromBlacklist :exec
DELETE FROM blacklisted WHERE subject_id = $1 AND cinema_id = $2;

-- name: CheckIfBlacklisted :one
SELECT EXISTS(SELECT(1) FROM blacklisted WHERE subject_id = $1 AND cinema_id = $2);

-- name: SaveRate :exec
INSERT INTO rated (user_id, cinema_id, upd) VALUES ($1, $2, $3);

-- name: RateExists :one
SELECT EXISTS(SELECT(1) FROM rated WHERE user_id = $1 AND cinema_id = $2);

-- name: UpdateRate :exec
UPDATE rated SET upd = $3 WHERE user_id = $1 AND cinema_id = $2;

-- name: GetRate :one
SELECT upd FROM rated WHERE user_id = $1 AND cinema_id = $2 LIMIT 1;

-- name: GetFavorites :many
SELECT * FROM saved WHERE subject_id = $1 LIMIT $2;

-- name: MarkAsViewed :exec
INSERT INTO viewed (subject_id, cinema_id) VALUES ($1, $2);

-- name: MarkAsViewedBulk :copyfrom
INSERT INTO viewed (subject_id, cinema_id) VALUES ($1, $2);

-- name: IsFavorite :one
SELECT EXISTS(SELECT(1) FROM saved WHERE subject_id = $1 AND cinema_id = $2);

-- name: IsBlacklisted :one
SELECT EXISTS(SELECT(1) FROM blacklisted WHERE subject_id = $1 AND cinema_id = $2);

-- name: UserSearch :many
SELECT * FROM users 
ORDER BY levenshtein(lower(users.login), lower($1))
LIMIT 10; 



