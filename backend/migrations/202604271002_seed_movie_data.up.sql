-- Seed data for movie table
INSERT INTO movie (id, name, year, language, imdb_rating, genre) VALUES
(gen_random_uuid(), 'Inception', 2010, 'english', 8.8, ARRAY['sci-fi', 'thriller']::genre[]),
(gen_random_uuid(), 'The Shawshank Redemption', 1994, 'english', 9.3, ARRAY['drama']::genre[]),
(gen_random_uuid(), 'Pulp Fiction', 1994, 'english', 8.9, ARRAY['drama', 'thriller']::genre[]),
(gen_random_uuid(), 'The Dark Knight', 2008, 'english', 9.0, ARRAY['drama', 'thriller']::genre[]),
(gen_random_uuid(), 'Forrest Gump', 1994, 'english', 8.8, ARRAY['drama', 'romcom']::genre[]),
(gen_random_uuid(), 'The Matrix', 1999, 'english', 8.7, ARRAY['sci-fi', 'thriller']::genre[]),
(gen_random_uuid(), 'Titanic', 1997, 'english', 7.8, ARRAY['drama', 'romcom']::genre[]),
(gen_random_uuid(), 'Avatar', 2009, 'english', 7.8, ARRAY['sci-fi']::genre[]),
(gen_random_uuid(), 'The Godfather', 1972, 'english', 9.2, ARRAY['drama', 'thriller']::genre[]),
(gen_random_uuid(), 'Schindler''s List', 1993, 'english', 8.9, ARRAY['drama']::genre[]);