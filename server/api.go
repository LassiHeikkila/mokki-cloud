package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/LassiHeikkila/mokki-cloud/server/auth"
)

var (
	RootPath = "www/index.html"

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
		interval time.Duration,
	) []Measurement = nil
)

const (
	defaultRangeInterval = 30 * time.Minute
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

func Authenticated(w http.ResponseWriter, req *http.Request) bool {
	// header should have X-API-KEY containing valid token
	key := req.Header.Get("X-API-KEY")
	if key == "" {
		return false
	}
	return auth.TokenIsValid(key)
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

func HandleCheckToken(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if Authenticated(w, req) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"ok":false}`))
	}
}

func HandleAuthorization(w http.ResponseWriter, req *http.Request) {
	type authRequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ok":false,"token":""}`))
		return
	}
	var arb authRequestBody
	err = json.Unmarshal(b, &arb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"ok":false,"token":""}`))
		return
	}
	log.Printf(`auth request with username "%s" and password "%s"`, arb.Username, arb.Password)
	if !auth.IsAuthorizedUser(arb.Username, arb.Password) {
		log.Println("credentials nok")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"ok":false,"token":""}`))
		return
	}
	log.Println("credentials ok")
	token, err := auth.GenerateToken(0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ok":false,"token":""}`))
		return
	}
	resp := make(map[string]interface{})
	resp["ok"] = true
	resp["token"] = token
	b, err = json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"ok":false,"token":""}`))
		return
	}
	_, _ = w.Write(b)
}

func HandleLatest(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if !Authenticated(w, req) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
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
	if !Authenticated(w, req) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
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
	interval, err := getDurationFromQueryOrDefault(req.URL.Query(), "interval", defaultRangeInterval)
	if err != nil {
		log.Println("error getting interval from query:", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	data := QueryTimeRange(req.Context(), field, id, start, stop, interval)
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

func getDurationFromQueryOrDefault(values url.Values, key string, defaultDuration time.Duration) (time.Duration, error) {
	value := values.Get(key)
	if value == "" {
		return defaultDuration, nil
	}
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return time.Duration(i) * time.Second, nil
}
