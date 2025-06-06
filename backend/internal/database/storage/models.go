// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package storage

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type Blacklisted struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

type Cinema struct {
	ID              uuid.UUID
	Private         bool
	Title           string
	OriginalTitle   *string
	ReleaseYear     *int32
	AgeRating       *int32
	DurationMinutes *int32
	PosterUrl       *string
	Description     *string
	Genres          []string
	Actors          []string
	Rating          *float64
	Embedding       pgvector.Vector
}

type Couch struct {
	ID        uuid.UUID
	Name      string
	Author    string
	Embedding pgvector.Vector
}

type CouchSitter struct {
	CouchID  uuid.UUID
	UserName string
}

type Rated struct {
	UserID   uuid.UUID
	CinemaID uuid.UUID
	Upd      int32
}

type Saved struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}

type User struct {
	ID          uuid.UUID
	Login       string
	Password    *string
	Privileged  bool
	Embedding   pgvector.Vector
	Provider    string
	YandexID    *string
	AccessToken *string
}

type Viewed struct {
	SubjectID uuid.UUID
	CinemaID  uuid.UUID
}
