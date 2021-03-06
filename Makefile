.PHONY: serve startdb stopdb initdb insert test

DSN_DEVELOPMENT="user=admin password=admin database=archibe sslmode=disable"
DSN_TEST="user=admin password=admin database=archibe-test port=5433 sslmode=disable"

DSN=${DSN_DEVELOPMENT}
ADDRESS=:4000

serve:
	go run main.go ${ADDRESS} ${DSN}

startdb:
	docker run --rm \
		--name archibe \
		-v archibe-postgtesql-data:/var/lib/postgresql/data \
		-e POSTGRES_DB=archibe \
		-e POSTGRES_USER=admin \
		-e POSTGRES_PASSWORD=admin \
		-p 5432:5432 \
		-d postgres:12.3

stopdb:
	docker rm archibe -f

initdb:
	go run ./scripts/initdb/main.go ${DSN}

insert:
	go run ./scripts/insert/main.go ${DSN} ${CHANNEL_ID}

test:
	(docker rm archibe-test -f || true) > /dev/null 2>&1
	docker run --rm --name archibe-test -e POSTGRES_DB=archibe-test -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -p 5433:5432 -d postgres:12.3
	# wait few seconds for start database
	sleep 2
	go run ./scripts/initdb/main.go ${DSN_TEST}
	go run ./scripts/testdata/main.go ${DSN_TEST}
	DSN=${DSN_TEST} go test -v ./...
	docker rm archibe-test -f
