.PHONY: startdb initdb insert

DSN="user=postgres password=admin database=tubekids sslmode=disable"

startdb:
	docker run --rm --name tubekids -e POSTGRES_DB=tubekids -e POSTGRES_PASSWORD=admin -p 5432:5432 -d postgres:12.3

initdb:
	go run ./scripts/initdb/main.go ${DSN}

insert:
	go run ./scripts/insert/main.go ${DSN} ${VIDEO_ID}
