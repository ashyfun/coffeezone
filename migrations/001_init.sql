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
