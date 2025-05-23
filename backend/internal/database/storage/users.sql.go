// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

const checkIfBlacklisted = `-- name: CheckIfBlacklisted :one
SELECT EXISTS(SELECT(1) FROM blacklisted WHERE subject_id = $1 AND cinema_id = $2)
`

type CheckIfBlacklistedParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) CheckIfBlacklisted(ctx context.Context, db DBTX, arg CheckIfBlacklistedParams) (bool, error) {
	row := db.QueryRow(ctx, checkIfBlacklisted, arg.SubjectID, arg.CinemaID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkIfFavorite = `-- name: CheckIfFavorite :one
SELECT EXISTS(SELECT(1) FROM saved WHERE subject_id = $1 AND cinema_id = $2)
`

type CheckIfFavoriteParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) CheckIfFavorite(ctx context.Context, db DBTX, arg CheckIfFavoriteParams) (bool, error) {
	row := db.QueryRow(ctx, checkIfFavorite, arg.SubjectID, arg.CinemaID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users(login, password, embedding) VALUES ($1, $2, $3)
`

type CreateUserParams struct {
	Login     string
	Password  *string
	Embedding pgvector.Vector
}

func (q *Queries) CreateUser(ctx context.Context, db DBTX, arg CreateUserParams) error {
	_, err := db.Exec(ctx, createUser, arg.Login, arg.Password, arg.Embedding)
	return err
}

const deleteFilmFromBlacklist = `-- name: DeleteFilmFromBlacklist :exec
DELETE FROM blacklisted WHERE subject_id = $1 AND cinema_id = $2
`

type DeleteFilmFromBlacklistParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) DeleteFilmFromBlacklist(ctx context.Context, db DBTX, arg DeleteFilmFromBlacklistParams) error {
	_, err := db.Exec(ctx, deleteFilmFromBlacklist, arg.SubjectID, arg.CinemaID)
	return err
}

const deleteFilmFromFavorites = `-- name: DeleteFilmFromFavorites :exec
DELETE FROM saved WHERE subject_id = $1 AND cinema_id = $2
`

type DeleteFilmFromFavoritesParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) DeleteFilmFromFavorites(ctx context.Context, db DBTX, arg DeleteFilmFromFavoritesParams) error {
	_, err := db.Exec(ctx, deleteFilmFromFavorites, arg.SubjectID, arg.CinemaID)
	return err
}

const getFavorites = `-- name: GetFavorites :many
SELECT subject_id, cinema_id FROM saved WHERE subject_id = $1 LIMIT $2
`

type GetFavoritesParams struct {
	SubjectID uuid.UUID
	Limit     int32
}

func (q *Queries) GetFavorites(ctx context.Context, db DBTX, arg GetFavoritesParams) ([]Saved, error) {
	rows, err := db.Query(ctx, getFavorites, arg.SubjectID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Saved
	for rows.Next() {
		var i Saved
		if err := rows.Scan(&i.SubjectID, &i.CinemaID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRate = `-- name: GetRate :one
SELECT upd FROM rated WHERE user_id = $1 AND cinema_id = $2 LIMIT 1
`

type GetRateParams struct {
	UserID   uuid.UUID
	CinemaID uuid.UUID
}

func (q *Queries) GetRate(ctx context.Context, db DBTX, arg GetRateParams) (int32, error) {
	row := db.QueryRow(ctx, getRate, arg.UserID, arg.CinemaID)
	var upd int32
	err := row.Scan(&upd)
	return upd, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, login, password, privileged, embedding, provider, yandex_id, access_token FROM users WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, db DBTX, id uuid.UUID) (User, error) {
	row := db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.Password,
		&i.Privileged,
		&i.Embedding,
		&i.Provider,
		&i.YandexID,
		&i.AccessToken,
	)
	return i, err
}

const getUserByLogin = `-- name: GetUserByLogin :one
SELECT id, login, password, privileged, embedding, provider, yandex_id, access_token FROM users WHERE login = $1
`

func (q *Queries) GetUserByLogin(ctx context.Context, db DBTX, login string) (User, error) {
	row := db.QueryRow(ctx, getUserByLogin, login)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Login,
		&i.Password,
		&i.Privileged,
		&i.Embedding,
		&i.Provider,
		&i.YandexID,
		&i.AccessToken,
	)
	return i, err
}

const getVector = `-- name: GetVector :one
SELECT embedding FROM users WHERE login = $1
`

func (q *Queries) GetVector(ctx context.Context, db DBTX, login string) (pgvector.Vector, error) {
	row := db.QueryRow(ctx, getVector, login)
	var embedding pgvector.Vector
	err := row.Scan(&embedding)
	return embedding, err
}

const isBlacklisted = `-- name: IsBlacklisted :one
SELECT EXISTS(SELECT(1) FROM blacklisted WHERE subject_id = $1 AND cinema_id = $2)
`

type IsBlacklistedParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) IsBlacklisted(ctx context.Context, db DBTX, arg IsBlacklistedParams) (bool, error) {
	row := db.QueryRow(ctx, isBlacklisted, arg.SubjectID, arg.CinemaID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isFavorite = `-- name: IsFavorite :one
SELECT EXISTS(SELECT(1) FROM saved WHERE subject_id = $1 AND cinema_id = $2)
`

type IsFavoriteParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) IsFavorite(ctx context.Context, db DBTX, arg IsFavoriteParams) (bool, error) {
	row := db.QueryRow(ctx, isFavorite, arg.SubjectID, arg.CinemaID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const markAsViewed = `-- name: MarkAsViewed :exec
INSERT INTO viewed (subject_id, cinema_id) VALUES ($1, $2)
`

type MarkAsViewedParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) MarkAsViewed(ctx context.Context, db DBTX, arg MarkAsViewedParams) error {
	_, err := db.Exec(ctx, markAsViewed, arg.SubjectID, arg.CinemaID)
	return err
}

type MarkAsViewedBulkParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

const rateExists = `-- name: RateExists :one
SELECT EXISTS(SELECT(1) FROM rated WHERE user_id = $1 AND cinema_id = $2)
`

type RateExistsParams struct {
	UserID   uuid.UUID
	CinemaID uuid.UUID
}

func (q *Queries) RateExists(ctx context.Context, db DBTX, arg RateExistsParams) (bool, error) {
	row := db.QueryRow(ctx, rateExists, arg.UserID, arg.CinemaID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const saveFilmToBlacklist = `-- name: SaveFilmToBlacklist :exec
INSERT INTO blacklisted (subject_id, cinema_id) VALUES ($1, $2)
`

type SaveFilmToBlacklistParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) SaveFilmToBlacklist(ctx context.Context, db DBTX, arg SaveFilmToBlacklistParams) error {
	_, err := db.Exec(ctx, saveFilmToBlacklist, arg.SubjectID, arg.CinemaID)
	return err
}

const saveFilmToFavorites = `-- name: SaveFilmToFavorites :exec
INSERT INTO saved (subject_id, cinema_id) VALUES ($1, $2)
`

type SaveFilmToFavoritesParams struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

func (q *Queries) SaveFilmToFavorites(ctx context.Context, db DBTX, arg SaveFilmToFavoritesParams) error {
	_, err := db.Exec(ctx, saveFilmToFavorites, arg.SubjectID, arg.CinemaID)
	return err
}

const saveRate = `-- name: SaveRate :exec
INSERT INTO rated (user_id, cinema_id, upd) VALUES ($1, $2, $3)
`

type SaveRateParams struct {
	UserID   uuid.UUID
	CinemaID uuid.UUID
	Upd      int32
}

func (q *Queries) SaveRate(ctx context.Context, db DBTX, arg SaveRateParams) error {
	_, err := db.Exec(ctx, saveRate, arg.UserID, arg.CinemaID, arg.Upd)
	return err
}

const updateRate = `-- name: UpdateRate :exec
UPDATE rated SET upd = $3 WHERE user_id = $1 AND cinema_id = $2
`

type UpdateRateParams struct {
	UserID   uuid.UUID
	CinemaID uuid.UUID
	Upd      int32
}

func (q *Queries) UpdateRate(ctx context.Context, db DBTX, arg UpdateRateParams) error {
	_, err := db.Exec(ctx, updateRate, arg.UserID, arg.CinemaID, arg.Upd)
	return err
}

const uploadVector = `-- name: UploadVector :exec
UPDATE users
SET embedding = $1
WHERE id = $2
`

type UploadVectorParams struct {
	Embedding pgvector.Vector
	ID        uuid.UUID
}

func (q *Queries) UploadVector(ctx context.Context, db DBTX, arg UploadVectorParams) error {
	_, err := db.Exec(ctx, uploadVector, arg.Embedding, arg.ID)
	return err
}

const userExists = `-- name: UserExists :one
SELECT EXISTS(SELECT(1) FROM users WHERE login = $1)
`

func (q *Queries) UserExists(ctx context.Context, db DBTX, login string) (bool, error) {
	row := db.QueryRow(ctx, userExists, login)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const userSearch = `-- name: UserSearch :many
SELECT id, login, password, privileged, embedding, provider, yandex_id, access_token FROM users 
ORDER BY levenshtein(lower(users.login), lower($1))
LIMIT 10
`

func (q *Queries) UserSearch(ctx context.Context, db DBTX, lower string) ([]User, error) {
	rows, err := db.Query(ctx, userSearch, lower)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Login,
			&i.Password,
			&i.Privileged,
			&i.Embedding,
			&i.Provider,
			&i.YandexID,
			&i.AccessToken,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
