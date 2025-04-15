BEGIN;
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
CREATE TABLE IF NOT EXISTS cinemas(
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	private boolean NOT NULL DEFAULT false,
	title text NOT NULL,
	original_title text,
	release_year int,
	age_rating int,
	duration_minutes int,
	poster_url text,
	description text,
	genres  text[],
	actors text[],
	rating float,
	embedding vector(46)
);
CREATE INDEX idx_title_fulltext ON cinemas USING GIN(
	to_tsvector('russian', title), to_tsvector('english', title)
);

CREATE INDEX ON cinemas USING hnsw (embedding vector_cosine_ops);

CREATE TABLE IF NOT EXISTS rated (
	user_id uuid NOT NULL REFERENCES users(id),
	cinema_id uuid NOT NULL REFERENCES cinemas(id),
	upd int NOT NULL
);


CREATE TABLE IF NOT EXISTS saved (
	subject_id uuid NOT NULL,
	cinema_id uuid NOT NULL REFERENCES cinemas(id)
);

CREATE TABLE IF NOT EXISTS blacklisted (
	subject_id uuid NOT NULL,
	cinema_id uuid NOT NULL REFERENCES cinemas(id)
);

COMMIT;
