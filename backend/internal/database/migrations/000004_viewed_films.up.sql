CREATE TABLE IF NOT EXISTS viewed (
	subject_id uuid NOT NULL,
	cinema_id uuid NOT NULL REFERENCES cinemas(id)
);
