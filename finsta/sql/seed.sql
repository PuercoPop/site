-- -*- sql-product: postgres -*-

-- Seeds for local development

INSERT INTO
  finsta.users (email, password)
VALUES
  ('honcho@puercopop.com', crypt('topsecret', gen_salt('bf', 8)));
