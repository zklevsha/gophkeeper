CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		email VARCHAR(50) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL);

CREATE TABLE IF NOT EXISTS private_types (
		id serial PRIMARY KEY, 
		name VARCHAR(50) UNIQUE NOT NULL);

CREATE TABLE IF NOT EXISTS private_data (
	id serial PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	user_id integer REFERENCES users (id),
	type_id integer REFERENCES private_types(id),
	khash_base64 TEXT,
	data_base64 TEXT,
	UNIQUE (id, name));
	