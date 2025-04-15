BEGIN;
CREATE TABLE IF NOT EXISTS couches (
	id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
	name text NOT NULL,
	author text NOT NULL,
    embedding vector(46)
);

CREATE TABLE IF NOT EXISTS couch_sitters (
	couch_id uuid NOT NULL REFERENCES couches(id),
    user_name text NOT NULL
);
COMMIT;
