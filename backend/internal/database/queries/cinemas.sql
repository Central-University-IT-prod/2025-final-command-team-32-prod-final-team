-- name: FilmExists :one
SELECT EXISTS(SELECT(1) FROM cinemas WHERE id = $1);

-- name: GetCinemaByID :one
SELECT * FROM cinemas WHERE id = $1;

-- name: GetUserModelByLogin :one
SELECT * FROM users WHERE login = $1;

-- name: GetFeed :many
SELECT * FROM cinemas
WHERE private = False AND
NOT EXISTS(SELECT(1) FROM viewed WHERE viewed.cinema_id = cinemas.id AND viewed.subject_id = $3)
ORDER BY cinemas.embedding <=> $1 LIMIT $2;


-- name: CreateCinema :one
INSERT INTO cinemas (
    private, title, original_title,
    release_year, age_rating, duration_minutes,
    poster_url, description, genres,
    actors, rating, embedding
) VALUES (
        $1, $2, $3, $4, $5, $6,
        $7, $8, $9, $10, $11, $12
) RETURNING id;


-- name: GetTopRated :many
SELECT * FROM cinemas
WHERE private = False
ORDER BY cinemas.rating DESC OFFSET $2 LIMIT $1;

-- name: SetCinemaPic :exec
UPDATE cinemas 
SET poster_url = $2
WHERE id = $1;

-- name: GetCinemaPic :one
SELECT poster_url FROM cinemas
WHERE id = $1;

-- name: SearchFilm :many
SELECT * FROM cinemas
WHERE private = False 
ORDER BY levenshtein(lower(cinemas.title), lower($1)) ASC
LIMIT 10;

-- name: SearchFilmWitTags :many
SELECT * FROM cinemas
WHERE private = False 
AND genres @> $2
ORDER BY levenshtein(lower(cinemas.title), lower($1)) ASC
LIMIT 10;

-- name: SearchTagsOnly :many
SELECT * FROM cinemas
WHERE genres @> $1
ORDER BY rating LIMIT 50;

-- name: UpdateCinema :exec
UPDATE cinemas SET 
    title = $1, original_title = $2,
    release_year = $3, age_rating = $4,
    duration_minutes = $5, poster_url = $6,
    description = $7, genres = $8,
    actors = $9, rating = $10
WHERE id = $11;

-- name: DeleteCinema :exec
DELETE FROM cinemas WHERE id = $1;

-- name: GetAllCinemas :many
SELECT * FROM cinemas;
