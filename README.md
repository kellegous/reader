# Reader

**Reader** is a [Miniflux](https://miniflux.app/)-based RSS reader. It starts and manages Miniflux and PostgreSQL, and serves Reader's web UI in front of Miniflux.

## Configuration

Reader reads `reader.yaml` by default. Start by copying the example:

```bash
cp reader.example.yaml reader.yaml
```

The minimal required configuration is:

```yaml
miniflux:
  # these are the credentials used to login as admin to miniflux.
  admin-username: miniflux
  admin-password: your_own_secret_used_to_login_to_miniflux
postgres:
  # miniflux will use this to login to postgres.
  password: your_own_secret_used_as_postgres_password
web:
  # The public host and port used to reach Reader.
  hostname: localhost:8080
```

`postgres.data-dir`, `postgres.database`, and `postgres.username` are optional; they default to `db`, `reader`, and `reader`, respectively. Relative data directories are resolved relative to the configuration file.

`web.addr` controls the listening address and defaults to `:4040`. `web.hostname` must be the public host and port (without a scheme) that Miniflux should use for its URLs. When running in Docker with the port mapping below, leave it as `localhost:8080`.

Optional settings:

```yaml
miniflux:
  # Automatically sign requests in as this Miniflux user.
  auto-login-as: your-miniflux-username
ollama:
  # Defaults: http://localhost:11434 and gemma3:27b
  url: http://localhost:11434
  model: gemma3:27b
```

## Building and running in Docker

### Build the image

```bash
docker build -t reader .
```

### Run the container

```bash
docker run -ti --rm \
  --name reader \
  -p 8080:4040 \
  -v $(pwd):/data \
  reader
```

Open <http://localhost:8080>. The mounted directory holds `reader.yaml` and, by default, the PostgreSQL data directory at `db/`.

## Developing

### Using Docker

The development shell supplies PostgreSQL, Miniflux, Bun, and the Go toolchain. It exposes Reader on port 4040 by default:

```bash
./etc/dev-shell
```

Then you can build and run reader until your heart's content.

```bash
make develop
```

`make develop` starts Reader with the Vite development server and hot module reloading. Ensure `reader.yaml` exists first. Pass `--port` to `./etc/dev-shell` to use a different host port.

### Running locally

If Miniflux and a compatible PostgreSQL installation are available locally, build and run the server directly:

```bash
make

bin/reader server --config-file=reader.yaml
```

Reader starts PostgreSQL and Miniflux as subprocesses, so their binaries must be available on `PATH`. Use `bin/reader server --help` to see server options, including logging and debug flags.

### Checks and formatting

```bash
make test
make lint
make fmt
```

Run all three with:

```bash
make validate
```

## Updating Miniflux

To bump the bundled Miniflux version in the setup scripts and `go.mod`:

```bash
./etc/update-miniflux 2.3.3
```

## Author(s)

- [Kelly Norton](https://github.com/kellegous)
