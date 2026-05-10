CREATE TYPE genre AS ENUM (
  'sci-fi',
  'drama',
  'romcom',
  'horror',
  'horror-comedy',
  'thriller'
);

CREATE TYPE language AS ENUM (
  'english'
);

CREATE TABLE movie (
  id uuid PRIMARY KEY,
  name text NOT NULL,
  year smallint NOT NULL,
  language language,
  imdb_rating numeric(3,1),
  genre genre[] NOT NULL
);