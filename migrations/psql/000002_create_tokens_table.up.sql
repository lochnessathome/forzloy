CREATE TABLE IF NOT EXISTS tokens(
   user_id int UNIQUE REFERENCES users(id),
   access_token varchar NOT NULL,
   access_token_exires_at timestamp without time zone NOT NULL,
   access_token_issued_at timestamp without time zone NOT NULL
);

