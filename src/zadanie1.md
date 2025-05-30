a. Budowanie aplikacji:

docker build -t aplikacja-pogodowa-michal-krocz .

b. Uruchomienie kontenera:

docker run -p 8080:8080 aplikacja-pogodowa-michal-krocz

c. Sprawdzenie logów

docker logs <id kontenera>

d. Rozmiar obrazu

docker images

