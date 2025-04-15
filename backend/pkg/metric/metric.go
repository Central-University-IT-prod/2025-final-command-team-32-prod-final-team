package metric

import "github.com/prometheus/client_golang/prometheus"

type PromMetrics struct {
	TotalRegistered    prometheus.Counter
	TotalSeen          prometheus.Counter
	TotalLikedFilms    prometheus.Counter
	TotalDislikedFilms prometheus.Counter
	TotalCouches       prometheus.Counter
}
