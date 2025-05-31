# Reader

**Reader** is a service that runs a [Miniflux](https://miniflux.app/) instance in a single container. I use this to run a local RSS reader.

## Configuration

Reader requires a pretty minimal configuration file.

```yaml
miniflux:
  # these are the credentials used to login as admin to miniflux.
  admin-username: miniflux
  admin-password: your_own_secret_used_to_login_to_miniflux
postgres:
  # miniflux will use this to login to postgres.
  password: your_own_secret_used_as_postgres_password
web:
  # set this to the hostname you intend to serve miniflux on.
  hostname: localhost:8080
```

Copy [reader.example.yaml](reader.example.yaml) to `reader.yaml` and fill in the values.

## Building and running in Docker

### Build the image

```bash
docker build -t reader .
```

### Run the container

```bash
docker run -d \
  --name reader \
  -p 8080:8080 \
  -v $(pwd):/data \
  reader
```

## Developing

### Using docker

Since reader requires a postgres database, the default development environment relies on docker. You can start a container that gives you a shell with the following command:

```bash
./etc/dev-shell
```

Then you can build and run reader until your heart's content.

```bash
make

bin/reader --config-file=reader.yaml
```

### Raw dogging it locally

The only reason that `./etc/dev-shell` exists is to make it easy to have an isolated postgres database. If you already have postgres installed on your machine, you can use that instead. Be aware, though, that `reader` opens both miniflux and postgres as subproceesses, so if you already run postgres on your machine, that might be a problem.

## Author(s)

- [Kelly Norton](https://github.com/kellegous)
