BEGIN;

INSERT INTO users (username, password, age)
SELECT
  ('user' || i::text) AS username,
  md5('password' || i::text) AS password,
  (floor(random()*60) + 18)::int AS age
FROM generate_series(1, 500000) AS s(i);

COMMIT;

-- output example:
-- id   |   username    |   password    |   age     | tokens | tokens
-- i+1  |   user1       |   7c6a180b... |   18..77  |   null | null