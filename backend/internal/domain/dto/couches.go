package dto

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

type CreateCouch struct {
	Name    string   `json:"name"`
	Sitters []string `json:"users"`

	AuthorName string    `json:"-"`
	AuthorID   uuid.UUID `json:"-"`

	Embedding pgvector.Vector `json:"-"`
}

type CouchView struct {
	Id        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Sitters   []string        `json:"users"`
	Author    string          `json:"author"`
	Embedding pgvector.Vector `json:"-"`
}

type UpdateCouch struct {
	Name    *string   `json:"name"`
	Sitters *[]string `json:"users"`

	AuthorName string `json:"-"`
}
