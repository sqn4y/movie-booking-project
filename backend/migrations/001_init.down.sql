
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP INDEX IF EXISTS idx_bookings_movie_id;
DROP INDEX IF EXISTS idx_bookings_status;
DROP INDEX IF EXISTS idx_movie_genre_movie_id;
DROP INDEX IF EXISTS idx_movie_genre_genre_id;

DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS movie_genre;
DROP TABLE IF EXISTS movies;
DROP TABLE IF EXISTS genre;
