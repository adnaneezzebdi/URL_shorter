CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE short_urls (
  id            UUID DEFAULT gen_random_uuid() PRIMARY KEY,
  code          VARCHAR(16) NOT NULL UNIQUE,
  original_url  TEXT NOT NULL UNIQUE,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  access_count  BIGINT DEFAULT 0 NOT NULL,
  last_access_at TIMESTAMPTZ NULL
);