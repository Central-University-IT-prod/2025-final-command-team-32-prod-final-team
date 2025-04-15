BEGIN;
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;
CREATE TABLE IF NOT EXISTS users(
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	login text NOT NULL UNIQUE,
	password text,
	privileged boolean NOT NULL DEFAULT false, 
	embedding vector(46),
	provider text NOT NULL DEFAULT 'local',
	yandex_id text,
	access_token text
);
CREATE INDEX idx_login_fulltext ON users USING GIN(
	to_tsvector('russian', login), to_tsvector('english', login)
);
CREATE INDEX ON users USING hnsw (embedding vector_cosine_ops);

COMMIT;
