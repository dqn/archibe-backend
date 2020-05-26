# tubekids-backend

## Development

### Start database locally using Docker.

```bash
$ docker run --name tubekids -e POSTGRES_DB=tubekids -e POSTGRES_PASSWORD=admin -p 5432:5432 -d postgres:12.3
```

### Init Database

```bash
$ make initdb
```

### Insert data

```bash
$ VIDEO_ID=<VIDEO_ID> make insert
```
