package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Definimos las metricas que vamos a utilizar en la aplicacion, estas metricas se pueden usar en cualquier parte de la aplicacion para medir el rendimiento y el comportamiento de las diferentes operaciones.
var (
	HttpRequestsCounter = promauto.NewCounter( // un contador
		prometheus.CounterOpts{
			Name: "runners_app_http_requests", // nombre de la metrica
			Help: "Número total de peticiones HTTP", // descripcion de la metrica
		},
	)

	GetRunnerHttpResponsesCounter = promauto.NewCounterVec( // un contador con etiquetas. Las etiquetas nos permiten clasificar las métricas en diferentes dimensiones. En este caso, vamos a clasificar las respuestas HTTP del endpoint get runner por su código de estado (status).
		prometheus.CounterOpts{
			Name: "runners_app_get_runner_http_responses",
			Help: "Número total de respuestas HTTP para el endpoint get runner",
		},
		[]string{"estado"}, // creamos una etiqueta llamada estado para clasificar las respuestas HTTP por su código de estado
	)

	GetAllRunnersTimer = promauto.NewHistogram( // un histograma. Un histograma nos permite medir la distribución de los tiempos de ejecución de una operación. En este caso, vamos a medir la duración de la operación get all runners.
		prometheus.HistogramOpts{
			Name: "runners_app_get_all_runners_duration",
			Help: "Duración de la operación get all runners en segundos",
		},
	)
)
