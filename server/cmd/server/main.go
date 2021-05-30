package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/LassiHeikkila/mokki-cloud/server"
)

const applicationVersion = "0.0.1"

func main() {
	var (
		help    = flag.Bool("help", false, "Print help text")
		version = flag.Bool("version", false, "Print version")

		cert = flag.String("cert", "", "Path to certificate file")
		key  = flag.String("key", "", "Path to private TLS key file")

		influxDBConfigFile = flag.String("influxDBConfig", "influxdb.json", "Path to config JSON containing InfluxDB parameters")
	)
	flag.Parse()

	if *version {
		fmt.Println(applicationVersion)
		return
	}

	if *help {
		flag.Usage()
		return
	}

	var influxConfig InfluxDBConfig
	err := loadConfig(*influxDBConfigFile, &influxConfig)
	if err != nil {
		log.Println("error loading influxdb config:", err)
		return
	}

	q := server.NewQuerier(
		influxConfig.Address,
		influxConfig.AuthToken,
		influxConfig.Organization,
	)
	defer q.Close()

	bucket := influxConfig.Bucket
	measurement := influxConfig.Measurement

	server.QueryLatest = func(
		ctx context.Context,
		field string,
		id string,
	) server.Measurement {
		return q.QueryLastValue(ctx, bucket, field, id, measurement)
	}
	server.QueryTimeRange = func(
		ctx context.Context,
		field string,
		id string,
		start time.Time,
		stop time.Time,
	) []server.Measurement {
		return q.QueryBetweenTimes(ctx, bucket, field, id, measurement, start, stop)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", server.HandleRoot)
	r.HandleFunc("/api/data/{field}/{id}/latest", server.HandleLatest)
	r.HandleFunc("/api/data/{field}/{id}/range", server.HandleRange)

	const dir = "."
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	log.Fatal(http.ListenAndServeTLS(":443", *cert, *key, r))
}
