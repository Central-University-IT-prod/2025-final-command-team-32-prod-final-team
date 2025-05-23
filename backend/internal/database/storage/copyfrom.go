// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: copyfrom.go

package storage

import (
	"context"
)

// iteratorForCreateSitter implements pgx.CopyFromSource.
type iteratorForCreateSitter struct {
	rows                 []CreateSitterParams
	skippedFirstNextCall bool
}

func (r *iteratorForCreateSitter) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForCreateSitter) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].CouchID,
		r.rows[0].UserName,
	}, nil
}

func (r iteratorForCreateSitter) Err() error {
	return nil
}

func (q *Queries) CreateSitter(ctx context.Context, db DBTX, arg []CreateSitterParams) (int64, error) {
	return db.CopyFrom(ctx, []string{"couch_sitters"}, []string{"couch_id", "user_name"}, &iteratorForCreateSitter{rows: arg})
}

// iteratorForMarkAsViewedBulk implements pgx.CopyFromSource.
type iteratorForMarkAsViewedBulk struct {
	rows                 []MarkAsViewedBulkParams
	skippedFirstNextCall bool
}

func (r *iteratorForMarkAsViewedBulk) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForMarkAsViewedBulk) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].SubjectID,
		r.rows[0].CinemaID,
	}, nil
}

func (r iteratorForMarkAsViewedBulk) Err() error {
	return nil
}

func (q *Queries) MarkAsViewedBulk(ctx context.Context, db DBTX, arg []MarkAsViewedBulkParams) (int64, error) {
	return db.CopyFrom(ctx, []string{"viewed"}, []string{"subject_id", "cinema_id"}, &iteratorForMarkAsViewedBulk{rows: arg})
}

// iteratorForMarkCouchAsViewedBulk implements pgx.CopyFromSource.
type iteratorForMarkCouchAsViewedBulk struct {
	rows                 []MarkCouchAsViewedBulkParams
	skippedFirstNextCall bool
}

func (r *iteratorForMarkCouchAsViewedBulk) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForMarkCouchAsViewedBulk) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].SubjectID,
		r.rows[0].CinemaID,
	}, nil
}

func (r iteratorForMarkCouchAsViewedBulk) Err() error {
	return nil
}

func (q *Queries) MarkCouchAsViewedBulk(ctx context.Context, db DBTX, arg []MarkCouchAsViewedBulkParams) (int64, error) {
	return db.CopyFrom(ctx, []string{"viewed"}, []string{"subject_id", "cinema_id"}, &iteratorForMarkCouchAsViewedBulk{rows: arg})
}
