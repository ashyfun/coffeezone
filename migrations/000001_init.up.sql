CREATE TABLE IF NOT EXISTS cz_cafe_locations (
    id serial PRIMARY KEY,
    address_name varchar(255) NOT NULL,
    longitude float8 NOT NULL,
    latitude float8 NOT NULL,
    CONSTRAINT location_unique UNIQUE (longitude, latitude)
);

CREATE TABLE IF NOT EXISTS cz_cafes (
    code char(24) PRIMARY KEY,
    title varchar(255) NOT NULL,
    location_id integer REFERENCES cz_cafe_locations (id) ON DELETE SET NULL,
    updated_at timestamp DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS cz_topics (
    id serial PRIMARY KEY,
    feature varchar(255) NOT NULL,
    CONSTRAINT feature_unique UNIQUE (feature)
);

CREATE TABLE IF NOT EXISTS cz_cafes_topics (
    cafe_code char(24) REFERENCES cz_cafes (code),
    topic_id integer REFERENCES cz_topics (id),
    PRIMARY KEY (cafe_code, topic_id)
);
