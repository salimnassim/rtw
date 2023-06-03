# rtw

Provides a JSON API for interacting with rTorrent over XML-RPC.

The server exposes a simple index on the `/` route.

## API routes

`GET /api/methods`
Retrieves all system methods from rTorrent.

`GET /api/view/{view}`
Retrieves all torrents in the view.

`POST /api/load`
Uploads torrent metadata file (.torrent) as a multipart file upload. The form key should be `file`.

`GET /api/torrent/{info_hash}/{action}`
Action can be: `stop`, `start`, `files`, `peers`, `trackers`

## Practical examples

List all unregistered torrents

`curl 127.0.0.1:8080/api/view/main | jq -r '.torrents[] | select(.message | ascii_downcase | contains("unregistered torrent")) | .hash'`

### Environment variables

- `BIND_ADDRESS`: server IP:port (e.g. 0.0.0.0:8080)
- `URL`: rTorrent XML-RPC endpoint (e.g. https://hostname/rpc2)
- `BASIC_USERNAME`: rTorrent XML-RPC basic auth username (optional)
- `BASIC_PASSWORD`: rTorrent XML-RPC basic auth password (optional)
- `CORS_ORIGIN`: *
- `CORS_AGE`: 86400