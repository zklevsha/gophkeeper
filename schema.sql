CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		email VARCHAR(50) UNIQUE NOT NULL,
		password BYTEA NOT NULL);
