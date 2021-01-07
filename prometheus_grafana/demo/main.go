package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var port string

var httpRequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_request_count",
		Help: "http request count",
	},
	[]string{"endpoint"},
)

func init() {
	prometheus.MustRegister(httpRequestCount)
}

func main() {
	flag.StringVar(&port, "port", ":10240", "port of http server")
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/api", handler)

	log.Printf("listen on: %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	httpRequestCount.WithLabelValues(path).Inc()

	w.Write([]byte("hello zy"))
}
