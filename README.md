# tubekids-backend

## Development

### Start database locally using Docker.

```bash
$ make startdb
```

### Init database

```bash
$ make initdb
```

### Fetch chat data

```bash
$ CHANNEL_ID=<channel-id> VIDEO_ID=<video-id> make insert
```

### Insert test data

```bash
$ make testdata
```

### Test

```bash
$ make test
```
