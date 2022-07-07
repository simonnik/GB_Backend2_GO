package red

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

// MeasurableDBExec executes db.Exec method and collects metrics
// (total queries, errors and duration)
func MeasurableDBExec(db *sql.DB, query string, values ...string) error {
	funcName := "exec"
	requestsDB.WithLabelValues(funcName).Inc()
	timer := prometheus.NewTimer(durationDBFunc.WithLabelValues(funcName))
	defer timer.ObserveDuration()

	_, err := db.Exec(query, values)
	if err != nil {
		errorsDB.WithLabelValues(funcName).Inc()
	}
	return err
}

// MeasurableDBQuery executes db.Query method and collects metrics
// (total queries, errors and duration)
func MeasurableDBQuery(db *sql.DB, query string) (*sql.Rows, error) {
	funcName := "query"
	requestsDB.WithLabelValues(funcName).Inc()
	timer := prometheus.NewTimer(durationDBFunc.WithLabelValues(funcName))
	defer timer.ObserveDuration()

	rr, err := db.Query(query)
	if err != nil {
		errorsDB.WithLabelValues(funcName).Inc()
	}
	return rr, err
}

// MeasurableDBQueryRow executes db.QueryRow method and collects metrics
// (total queries, errors and duration)
func MeasurableDBQueryRow(db *sql.DB, query string, id string) *sql.Row {
	funcName := "query_row"
	requestsDB.WithLabelValues(funcName).Inc()
	timer := prometheus.NewTimer(durationDBFunc.WithLabelValues(funcName))
	defer timer.ObserveDuration()

	row := db.QueryRow(query, id)

	return row
}

// ErrorDBQueryRow counts errors on sql.Row.Scan method
func ErrorDBQueryRow() {
	funcName := "query_row"
	errorsDB.WithLabelValues(funcName).Inc()
	return
}
