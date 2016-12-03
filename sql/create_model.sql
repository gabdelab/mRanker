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
  year_ranking integer NOT NULL DEFAULT 0,
  CONSTRAINT unique_album UNIQUE (artist_id, name, year),
  CONSTRAINT unique_ranking UNIQUE (ranking) DEFERRABLE,
  CONSTRAINT unique_year_ranking UNIQUE(year, year_ranking) DEFERRABLE
);

CREATE OR REPLACE FUNCTION insert_album(year int, name text, artist text, ranking int) RETURNS void AS $$
  SET CONSTRAINTS unique_ranking DEFERRED;
  /* Move all rankings up to prepare the insertion */
  UPDATE albums
    SET ranking = ranking + 1
    WHERE ranking >= $4;
  INSERT INTO albums (artist_id, year, name, ranking, year_ranking)
  VALUES (
    (SELECT artist_id FROM artists WHERE name=$3),
    $1,
    $2,
    $4,
    (SELECT coalesce(max(year_ranking), 0) + 1 FROM albums WHERE year=$1)
  );
  /*By default, a new insertion has the highest year_ranking*/
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

CREATE OR REPLACE FUNCTION update_year_ranking(album_id int, old_ranking int, new_ranking int) RETURNS void AS $$
BEGIN
  SET CONSTRAINTS unique_year_ranking DEFERRED;
  IF $3 < $2 THEN
    /* Ranking decreased, increase all other rankings */
    UPDATE albums
      SET year_ranking = year_ranking + 1
      WHERE year_ranking >= $3
      AND year_ranking <= $2
      AND year = (SELECT year FROM albums WHERE albums.album_id=$1);
  ELSIF $2 < $3 THEN
    /* Ranking increased, decrease all other year_rankings */
    UPDATE albums
      SET year_ranking = year_ranking - 1
      WHERE year_ranking >= $2
      AND year_ranking <= $3
      AND year = (SELECT year FROM albums WHERE albums.album_id=$1);
  END IF;
  /* Give the new year_ranking to the album specified by album_id */
  UPDATE albums
    SET year_ranking = $3
    WHERE albums.album_id = $1;
END
$$ LANGUAGE PLPGSQL;

