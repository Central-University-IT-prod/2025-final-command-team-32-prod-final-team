package wrapper

import (
	"github.com/pgvector/pgvector-go"
	"solution/internal/database/storage"
	"solution/internal/domain/dto"
)

func CinemaWithView(c *storage.Cinema) *dto.CinemaView {
	return &dto.CinemaView{
		Id:              c.ID,
		Name:            c.Title,
		ReleaseYear:     c.ReleaseYear,
		AgeRating:       c.AgeRating,
		DurationMinutes: c.DurationMinutes,
		PosterURL:       c.PosterUrl,
		Description:     c.Description,
		Genres:          c.Genres,
		Actors:          c.Actors,
		Rating:          c.Rating,
		UserRating:      nil,
	}
}

func FeedWitView(c *storage.Cinema) *dto.CinemaView {
	return &dto.CinemaView{
		Id:              c.ID,
		Name:            c.Title,
		ReleaseYear:     c.ReleaseYear,
		AgeRating:       c.AgeRating,
		DurationMinutes: c.DurationMinutes,
		PosterURL:       c.PosterUrl,
		Description:     c.Description,
		Genres:          c.Genres,
		Actors:          c.Actors,
		Rating:          c.Rating,
		UserRating:      nil,
	}
}

func ToCreateCinema(c *dto.CinemaCreate, posterURL *string, description *string, rating *float64, vector pgvector.Vector) storage.CreateCinemaParams {
	params := storage.CreateCinemaParams{
		Private:         c.Private,
		Title:           c.Name,
		OriginalTitle:   c.OriginalTitle,
		ReleaseYear:     c.ReleaseYear,
		AgeRating:       c.AgeRating,
		DurationMinutes: c.DurationMinutes,
		PosterUrl:       posterURL,
		Description:     description,
		Genres:          c.Genres,
		Actors:          c.Actors,
		Rating:          rating,
		Embedding:       vector,
	}
	return params
}

func ToUpdateCinema(param dto.CinemaUpdate) storage.UpdateCinemaParams {
	return storage.UpdateCinemaParams{
		Title:           param.Title,
		OriginalTitle:   param.OriginalTitle,
		ReleaseYear:     param.ReleaseYear,
		AgeRating:       param.AgeRating,
		DurationMinutes: param.DurationMinutes,
		PosterUrl:       param.PosterUrl,
		Description:     param.Description,
		Genres:          param.Genres,
		Actors:          param.Actors,
		Rating:          param.Rating,
		ID:              param.ID,
	}
}
