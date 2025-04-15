-- name: UserIsPrivileged :one
SELECT privileged FROM users WHERE login = $1;
