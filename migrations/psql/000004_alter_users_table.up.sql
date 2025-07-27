ALTER TABLE users ADD balance int DEFAULT 0 NOT NULL;

ALTER TABLE users ADD CONSTRAINT balance_never_negative CHECK (balance >= 0);
