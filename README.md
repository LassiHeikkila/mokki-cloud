# Mokki-cloud
This project contains the backend and a website for visualizing environmental data from an InfluxDB database.

Frontend is implemented with React.

Backend is a custom Go web server.

## Running the backend in dev mode

Dev mode allows running with plain HTTP.

Running the backend requires access to an InfluxDB instance containing appropriate data.

Access credentials are given via a config file.
Path to the config file can be given with `-influxDBConfig <file>` argument.
Default is `influxdb.json` in the current working directory.

The config file content should look something like:

```json
{
    "address": "https://my-influxdb-instance.website.cloud",
    "org": "myuser@website.cloud",
    "token": "secret",
    "bucket": "Mokki data",
    "measurement": "ruuvidata"
}
```

Then to run the server, either compile and execute (in `server` directory):

```console
go build -o mokki-server ./cmd/server
./mokki-server -dev -influxDBConfig ../config.json -httpPort 8080
```

or run directly (in `server` directory):

```console
go run ./cmd/server -dev -influxDBConfig ../<path to config.json> -httpPort 8080 
```