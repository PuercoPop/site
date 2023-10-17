-- -*- sql-product: postgres -*-
BEGIN;
CREATE TABLE IF NOT EXISTS public.versions (
  version INTEGER PRIMARY KEY NOT NULL,
  checksum BYTEA NOT NULL
);
COMMIT;
