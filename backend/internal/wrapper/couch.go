package wrapper

import (
	"github.com/google/uuid"
	"solution/internal/database/storage"
	"solution/internal/domain/dto"
)

func SittersConvert(sitters []string, couchID uuid.UUID) []storage.CreateSitterParams {
	res := make([]storage.CreateSitterParams, 0)
	for _, sitter := range sitters {
		res = append(res, storage.CreateSitterParams{
			CouchID:  couchID,
			UserName: sitter,
		})
	}
	return res
}

func ToCreateCouch(param dto.CreateCouch) storage.CreateCouchParams {
	return storage.CreateCouchParams{
		Name:      param.Name,
		Author:    param.AuthorName,
		Embedding: param.Embedding,
	}
}

func ToCouchView(couch storage.Couch, sitters []string) *dto.CouchView {
	return &dto.CouchView{
		Id:        couch.ID,
		Name:      couch.Name,
		Sitters:   sitters,
		Author:    couch.Author,
		Embedding: couch.Embedding,
	}
}

func ToCouchesView(couch storage.GetCouchesRow, sitters []string) dto.CouchView {
	return dto.CouchView{
		Id:      couch.CouchID,
		Name:    couch.Name,
		Sitters: sitters,
		Author:  couch.Author,
	}
}

func CouchViewForDatabase(couchId uuid.UUID, ids []uuid.UUID) []storage.MarkCouchAsViewedBulkParams {
	res := make([]storage.MarkCouchAsViewedBulkParams, 0)
	for _, id := range ids {
		res = append(res, storage.MarkCouchAsViewedBulkParams{
			SubjectID: couchId,
			CinemaID:  id,
		})
	}
	return res
}
