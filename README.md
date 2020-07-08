# archibe-backend

Archibe web app for summarizing and browsing YouTube Live archives.

Demo page: https://archibe.vercel.app

## Develop

### Start database

```bash
$ make startdb
```

### Stop database

```bash
$ make stopdb
```

### Init database

```bash
$ make initdb
```

### Insert YouTube archive data

```bash
$ CHANNEL_ID=<channel-id> make insert
```

### Start Server

```bash
$ make serve
```

### Test

```bash
$ make test
```

### Backup

```bash
$ pg_dump archibe -h 127.0.0.1 -U admin > example.dump
```

### Restore

```bash
# psql -c 'create database archibe'
$ psql archibe < example.dump
```

## License

MIT
