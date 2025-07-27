CREATE TABLE IF NOT EXISTS users(
   id serial PRIMARY KEY,
   login varchar UNIQUE NOT NULL,
   password_hash varchar NOT NULL
);
