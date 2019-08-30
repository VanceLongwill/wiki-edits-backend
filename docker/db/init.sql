CREATE TABLE edits (
  id SERIAL PRIMARY KEY,
  lang_code VARCHAR(2),
  bytes_changed INTEGER,
  modified_at TIMESTAMP NOT NULL DEFAULT NOW()
);
