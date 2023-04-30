// nolint
package metrics

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var workerMetrics *Worker

type Worker struct {
	dataNetFlow prometheus.Gauge
}

func (sp *Worker) Register() (string, chi.Router) {
	routes := chi.NewRouter()

	routes.Handle("/", promhttp.Handler())
	return "/metrics", routes
}

func NewWorker() *Worker {
	var (
		dataNetFlow = promauto.NewGauge(prometheus.GaugeOpts{
			Name: "data_net_flow",
			Help: "Data Net Flow",
		})
	)
	workerMetrics = &Worker{dataNetFlow: dataNetFlow}
	return workerMetrics
}

func (sp *Worker) CountNetFlow(number int) {
	sp.dataNetFlow.Set(float64(number))
}
