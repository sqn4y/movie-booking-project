CREATE TABLE IF NOT EXISTS genre
(
    id   BIGSERIAL PRIMARY KEY,
    name VARCHAR(55) NOT NULL UNIQUE
);


CREATE TABLE IF NOT EXISTS movies
(
    id           BIGSERIAL PRIMARY KEY,
    title        VARCHAR(255) NOT NULL,
    director     VARCHAR(255) NOT NULL,
    duration     INTEGER      NOT NULL,
    description  TEXT,
    image_url    TEXT         NOT NULL DEFAULT '',
    age_rating   INTEGER      NOT NULL,
    release_date TIMESTAMP,
    created_at   TIMESTAMP,
    updated_at   TIMESTAMP
);


CREATE TABLE IF NOT EXISTS movie_genre
(
    movie_id BIGINT NOT NULL,
    genre_id BIGINT NOT NULL,
    PRIMARY KEY (movie_id, genre_id),
    CONSTRAINT fk_movie_genre_movie FOREIGN KEY (movie_id)
        REFERENCES movies (id) ON DELETE CASCADE,
    CONSTRAINT fk_movie_genre_genre FOREIGN KEY (genre_id)
        REFERENCES genre (id) ON DELETE RESTRICT
);


CREATE TABLE IF NOT EXISTS bookings
(
    id         BIGSERIAL PRIMARY KEY,
    user_id    UUID        NOT NULL,
    movie_id   BIGINT      NOT NULL,
    seats      TEXT[]      NOT NULL,
    status     VARCHAR(50) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_bookings_movie FOREIGN KEY (movie_id)
        REFERENCES movies (id) ON DELETE RESTRICT
);
CREATE INDEX idx_bookings_user_id ON bookings (user_id);
CREATE INDEX idx_bookings_movie_id ON bookings (movie_id);
CREATE INDEX idx_bookings_status ON bookings (status);
CREATE INDEX idx_movie_genre_movie_id ON movie_genre (movie_id);
CREATE INDEX idx_movie_genre_genre_id ON movie_genre (genre_id);
