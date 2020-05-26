.PHONY: initdb insert

DSN="user=postgres password=admin database=tubekids sslmode=disable"

initdb:
	go run ./scripts/initdb/main.go ${DSN}

insert:
	go run ./scripts/insert/main.go ${DSN} ${VIDEO_ID}
