CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS financial_operations(
   report_id uuid NOT NULL,
   user_id int REFERENCES users(id),
   cost int NOT NULL,
   state varchar NOT NULL,
   CONSTRAINT financial_operations_pkey PRIMARY KEY (report_id, user_id)
);

