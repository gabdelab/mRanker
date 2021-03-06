CREATE SEQUENCE artist_id_seq;
CREATE SEQUENCE album_id_seq;

CREATE TABLE artists (
  artist_id integer PRIMARY KEY NOT NULL DEFAULT NEXTVAL('artist_id_seq'),
  name varchar(40) UNIQUE NOT NULL
);

CREATE TABLE albums (
  album_id integer PRIMARY KEY NOT NULL DEFAULT NEXTVAL('album_id_seq'),
  name varchar(100) NOT NULL,
  artist_id integer NOT NULL REFERENCES artists (artist_id),
  year integer NOT NULL,
  ranking integer NOT NULL,
  CONSTRAINT unique_album UNIQUE (artist_id, name, year),
  CONSTRAINT unique_ranking UNIQUE (ranking) DEFERRABLE,
);

CREATE OR REPLACE FUNCTION insert_album(year int, name text, artist text, ranking int) RETURNS void AS $$
  SET CONSTRAINTS unique_ranking DEFERRED;
  /* Move all rankings up to prepare the insertion */
  UPDATE albums
    SET ranking = ranking + 1
    WHERE ranking >= $4;
  INSERT INTO albums (artist_id, year, name, ranking)
  VALUES (
    (SELECT artist_id FROM artists WHERE name=$3),
    $1,
    $2,
    $4
  );
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION delete_album(album_id int) RETURNS void AS $$
  WITH rank AS
    (DELETE FROM albums WHERE album_id=$1 RETURNING ranking)
  UPDATE albums
    SET ranking = ranking - 1
    WHERE ranking > (SELECT ranking FROM rank);
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION update_ranking(album_id int, old_ranking int, new_ranking int) RETURNS void AS $$
BEGIN
  SET CONSTRAINTS unique_ranking DEFERRED;
  IF $3 < $2 THEN
    /* Ranking decreased, increase all other rankings */
    UPDATE albums
      SET ranking = ranking + 1
      WHERE ranking >= $3
      AND ranking <= $2;
  ELSIF $2 < $3 THEN
    /* Ranking increased, decrease all other rankings */
    UPDATE albums
      SET ranking = ranking - 1
      WHERE ranking >= $2
      AND ranking <= $3;
  END IF;
  /* Give the new ranking to the album specified by album_id */
  UPDATE albums
    SET ranking = $3
    WHERE albums.album_id = $1;
END
$$ LANGUAGE PLPGSQL;

