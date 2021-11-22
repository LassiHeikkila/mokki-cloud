package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"

	"github.com/LassiHeikkila/mokki-cloud/server"
	"github.com/LassiHeikkila/mokki-cloud/server/auth"
)

const applicationVersion = "0.0.1"

var (
	httpsPort = 443
	httpPort  = 80
)

func main() {
	var (
		help    = flag.Bool("help", false, "Print help text")
		version = flag.Bool("version", false, "Print version")

		dev          = flag.Bool("dev", false, "Run in development mode (plain HTTP)")
		httpsPortArg = flag.Int("httpsPort", 443, "HTTPS port")
		httpPortArg  = flag.Int("httpPort", 80, "HTTP port")

		cert = flag.String("cert", "", "Path to certificate file")
		key  = flag.String("key", "", "Path to private TLS key file")

		influxDBConfigFile = flag.String("influxDBConfig", "influxdb.json", "Path to config JSON containing InfluxDB parameters")

		authDB = flag.String("authdb", "auth.db", "Path to authentication database")
	)
	flag.Parse()

	httpsPort = *httpsPortArg
	httpPort = *httpPortArg

	if *version {
		fmt.Println(applicationVersion)
		return
	}

	if *help {
		flag.Usage()
		return
	}

	if *authDB != "" {
		db, err := sql.Open("sqlite3", *authDB)
		if err != nil {
			log.Println("failed to open:", *authDB)
			return
		}
		defer db.Close()

		auth.RegisterDatabase(db)
		if err := auth.InitializeDatabase(); err != nil {
			log.Println("failed to initialize database", err)
			return
		}
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
		interval time.Duration,
	) []server.Measurement {
		return q.QueryBetweenTimes(ctx, bucket, field, id, measurement, start, stop, interval)
	}

	r := mux.NewRouter()
	r.HandleFunc("/", server.HandleRoot)
	r.HandleFunc("/api/authorize", server.HandleAuthorization)
	r.HandleFunc("/api/checkToken", server.HandleCheckToken)
	r.HandleFunc("/api/data/{field}/{id}/latest", server.HandleLatest)
	r.HandleFunc("/api/data/{field}/{id}/range", server.HandleRange)

	const dir = "static"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	if !*dev {
		go func() {
			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), http.HandlerFunc(redirectTLS)))
		}()

		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%d", httpsPort), *cert, *key, r))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), r))
	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("https://%s:%d/%s", r.Host, httpsPort, r.RequestURI), http.StatusMovedPermanently)
}
