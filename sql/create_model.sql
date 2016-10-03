CREATE SEQUENCE album_id_seq;
CREATE SEQUENCE artist_id_seq;

CREATE TABLE artists (
  artist_id integer PRIMARY KEY NOT NULL DEFAULT NEXTVAL('artist_id_seq'),
  name varchar(40) UNIQUE
);

CREATE TABLE albums (
  album_id integer PRIMARY KEY NOT NULL DEFAULT NEXTVAL('album_id_seq'),
  name varchar(100),
  artist_id integer NOT NULL REFERENCES artists (artist_id),
  year integer,
  CONSTRAINT albums_constraint UNIQUE (artist_id, name, year)
);

CREATE OR REPLACE FUNCTION insert_album(year int, name text, artist text) RETURNS integer AS $$
  INSERT INTO albums (artist_id, year, name) VALUES (
    (SELECT artist_id FROM artists WHERE name=$3),
    $1,
    $2
  ) RETURNING album_id;
$$ LANGUAGE SQL;
