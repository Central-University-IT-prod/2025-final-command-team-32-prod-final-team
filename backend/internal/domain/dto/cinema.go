package dto

import "github.com/google/uuid"

type CinemaView struct {
	Id              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	ReleaseYear     *int32    `json:"year,omitempty"`
	AgeRating       *int32    `json:"age_rating,omitempty"`
	DurationMinutes *int32    `json:"duration_minutes,omitempty"`
	PosterURL       *string   `json:"poster_url,omitempty"`
	Description     *string   `json:"description,omitempty"`
	Genres          []string  `json:"genres,omitempty"`
	Actors          []string  `json:"actors,omitempty"`
	Rating          *float64  `json:"rating,omitempty"`
	UserRating      *int32    `json:"user_rating,omitempty"`
}

type CinemaCreate struct {
	Name            string   `json:"name" validate:"required"`
	ReleaseYear     *int32   `json:"year,omitempty"`
	AgeRating       *int32   `json:"age_rating,omitempty"`
	DurationMinutes *int32   `json:"duration_minutes,omitempty"`
	PosterURL       string   `json:"poster_url,omitempty"`
	Description     string   `json:"description,omitempty"`
	Genres          []string `json:"genres,omitempty"`
	Actors          []string `json:"actors,omitempty"`
	Rating          float32  `json:"rating,omitempty"`
	UserRating      int      `json:"user_rating,omitempty"`
	Private         bool     `json:"-"`
	OriginalTitle   *string  `json:"original_title,omitempty"`
}

type CinemaCreateResponse struct {
	Id uuid.UUID `json:"id"`
}

type CinemaUpdate struct {
	Title           string    `json:"name"`
	OriginalTitle   *string   `json:"original_title,omitempty"`
	ReleaseYear     *int32    `json:"year,omitempty"`
	AgeRating       *int32    `json:"age_rating,omitempty"`
	DurationMinutes *int32    `json:"duration_minutes,omitempty"`
	PosterUrl       *string   `json:"poster_url,omitempty"`
	Description     *string   `json:"description,omitempty"`
	Genres          []string  `json:"genres,omitempty"`
	Actors          []string  `json:"actors,omitempty"`
	Rating          *float64  `json:"rating,omitempty"`
	ID              uuid.UUID `json:"-"`
}
