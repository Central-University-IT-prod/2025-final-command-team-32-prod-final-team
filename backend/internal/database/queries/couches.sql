-- name: CouchExists :one
SELECT EXISTS(SELECT(1) FROM couches WHERE id = $1);

-- name: CreateCouch :one
INSERT INTO couches (
    name, author, embedding
) VALUES (
    $1, $2, $3
) RETURNING id;

-- name: CreateSitter :copyfrom
INSERT INTO couch_sitters(couch_id, user_name) VALUES ($1, $2);

-- name: GetCouch :one
SELECT 
    couches.id, couches.name,
    couches.author, couches.embedding
FROM couches
WHERE id = $1;

-- name: GetCouchVector :one
SELECT embedding FROM couches WHERE id = $1;

-- name: UpdateCouch :exec
UPDATE couches SET name = $1 WHERE id = $2;

-- name: SetCouchVector :exec
UPDATE couches SET embedding = $1 WHERE id = $2;

-- name: GetCouches :many
SELECT DISTINCT 
    couch_sitters.couch_id,
    couches.name,
    couches.author
FROM couches INNER JOIN couch_sitters ON 
    couch_sitters.couch_id = couches.id
WHERE couch_sitters.user_name = $1;

-- name: GetSitters :many
SELECT couch_sitters.user_name FROM couch_sitters WHERE couch_id = $1;

-- name: ClearSitters :exec
DELETE FROM couch_sitters WHERE couch_id = $1;


-- name: GetCouchFeed :many
SELECT * FROM cinemas
WHERE private = False AND
NOT EXISTS(SELECT(1) FROM viewed WHERE viewed.cinema_id = cinemas.id AND viewed.subject_id = $3)
ORDER BY cinemas.embedding <=> $1 LIMIT $2;

-- name: MarkCouchAsViewedBulk :copyfrom
INSERT INTO viewed (subject_id, cinema_id) VALUES ($1, $2);
