package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	RootPath = "index.html"

	QueryLatest func(
		ctx context.Context,
		field string,
		id string,
	) Measurement = nil

	QueryTimeRange func(
		ctx context.Context,
		field string,
		id string,
		start time.Time,
		stop time.Time,
	) []Measurement = nil
)

func HandleRoot(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	b, err := ioutil.ReadFile(RootPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("error reading %s", RootPath), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		HandleRoot(w, req)
		return
	}
	if strings.HasSuffix(req.URL.Path, "/latest") {
		HandleLatest(w, req)
		return
	}
	if strings.HasSuffix(req.URL.Path, "/range") {
		HandleRange(w, req)
		return
	}
}

func HandleLatest(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	id, err := getSensorIDFromPath(req.URL.Path)
	if err != nil {
		log.Printf("error getting sensor id from request path (%s): %s\n", req.URL.Path, err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	field, err := getFieldFromPath(req.URL.Path)
	if err != nil {
		log.Printf("error getting field from request path (%s): %s\n", req.URL.Path, err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	data := QueryLatest(req.Context(), field, id)
	if data == nil {
		log.Println("nil data returned from query")
		http.Error(w, "no data found for given parameters", http.StatusNotFound)
		return
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Println("error marshalling data:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func HandleRange(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	id, err := getSensorIDFromPath(req.URL.Path)
	if err != nil {
		log.Printf("error getting sensor id from request path (%s): %s\n", req.URL.Path, err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	field, err := getFieldFromPath(req.URL.Path)
	if err != nil {
		log.Printf("error getting field from request path (%s): %s\n", req.URL.Path, err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	start, err := getTimeFromQuery(req.URL.Query(), "from")
	if err != nil {
		log.Println("error getting start time from query:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	stop, err := getTimeFromQuery(req.URL.Query(), "to")
	if err != nil {
		log.Println("error getting stop time from query:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	data := QueryTimeRange(req.Context(), field, id, start, stop)
	if data == nil {
		log.Println("nil data returned from query")
		http.Error(w, "no data found for given parameters", http.StatusNotFound)
		return
	}
	b, err := json.Marshal(data)
	if err != nil {
		log.Println("error marshalling data:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(b)
}

func getFieldFromPath(path string) (string, error) {
	// /api/data/{field}/{id}/latest
	// field is third item, but split counts the empty value before the first /
	segments := strings.Split(path, "/")
	if len(segments) != 6 {
		return "", fmt.Errorf("malformed path: %v", segments)
	}
	return segments[3], nil
}

func getSensorIDFromPath(path string) (string, error) {
	// /api/data/{field}/{id}/latest
	// id is fourth item, but split counts the empty value before the first /
	segments := strings.Split(path, "/")
	if len(segments) != 6 {
		return "", fmt.Errorf("malformed path: %v", segments)
	}
	return segments[4], nil
}

func getTimeFromQuery(values url.Values, key string) (time.Time, error) {
	value := values.Get(key)
	if value == "" {
		return time.Time{}, fmt.Errorf("requested key %s not present", key)
	}
	return time.Parse(time.RFC3339Nano, value)
}
