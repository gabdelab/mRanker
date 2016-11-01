psql mrankerdb -c "DELETE FROM albums WHERE album_id >=1;"
psql mrankerdb -c "DELETE FROM artists WHERE artist_id >=1;"

# Insert some artists
curl -d "name=Radiohead" localhost:8080/artist/
curl -d "name=Phoenix" localhost:8080/artist/
curl -d "name=Modest Mouse" localhost:8080/artist/
curl -d "name=Television" localhost:8080/artist/

# Insert some albums
curl -d "artist=Radiohead&name=Kid A&year=2000&rank=1" localhost:8080/album/
curl -d "artist=Radiohead&name=OK Computer&year=1997&rank=1" localhost:8080/album/
curl -d "artist=Radiohead&name=Amnesiac&year=2001&rank=1" localhost:8080/album/
curl -d "artist=Television&name=Marquee Moon&year=1977&rank=1" localhost:8080/album/
curl -d "artist=Television&name=Adventure&year=1979&rank=4" localhost:8080/album/
curl -d "artist=Phoenix&name=Alphabetical&year=2000&rank=3" localhost:8080/album/
curl -d "artist=Modest Mouse&name=The Moon and Antarctica&year=2000&rank=4" localhost:8080/album/

# Perform some updates
curl -d "artist=Radiohead&name=Amnesiac&year=2001&rank=1" localhost:8080/album/
curl -d "artist=Modest Mouse&name=The Moon and Antarctica&year=2000&rank=7" localhost:8080/album/

curl localhost:8080
