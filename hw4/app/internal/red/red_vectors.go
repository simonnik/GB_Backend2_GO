package red

import "github.com/prometheus/client_golang/prometheus"

var (
	durationReq = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "backend2_duration_seconds",
			Help:       "Summary of request duration in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{"URI", "METHOD"},
	)
	errorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "backend2_errors_total",
			Help: "Total number of errors",
		},
		[]string{"URI", "METHOD", "CODE"},
	)
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "backend2_request_total",
			Help: "Total number of requests",
		},
		[]string{"URI", "METHOD"},
	)
)

var (
	durationDBFunc = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "backend2_db_query_duration_seconds",
			Help:       "Summary of database method request duration in seconds",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{"FUNC"},
	)

	errorsDB = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "backend2_db_query_errors",
			Help: "Amount of errors on database requests",
		},
		[]string{"FUNC"},
	)

	requestsDB = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "backend2_db_query_counter",
			Help: "Total amount of database requests",
		},
		[]string{"URI"},
	)
)

func init() {
	prometheus.MustRegister(durationReq)
	prometheus.MustRegister(errorsTotal)
	prometheus.MustRegister(requestsTotal)

	prometheus.MustRegister(durationDBFunc)
	prometheus.MustRegister(errorsDB)
	prometheus.MustRegister(requestsDB)
}
