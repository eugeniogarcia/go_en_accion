package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// definimos el exporter de métricas de Prometheus. En este caso, exponemos las métricas en el endpoint /metrics del puerto 9000
func InitPrometheus() {
	http.Handle("/metrics", promhttp.Handler()) // endpoint en el que expondremos las métricas
	http.ListenAndServe(":9000", nil)
}
