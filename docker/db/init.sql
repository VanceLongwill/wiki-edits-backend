CREATE TABLE edits (
  id SERIAL PRIMARY KEY,
  lang_code VARCHAR(2),
  byte_change INTEGER,
  modified_at TIMESTAMPZ NOT NULL DEFAULT NOW()
);
