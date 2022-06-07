package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "stonechat_processed_requests_total",
		Help: "The total number of processed requests",
	})
)
