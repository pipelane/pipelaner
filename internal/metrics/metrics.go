/*
 * Copyright (c) 2024 Alexey Khokhlov
 */

package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(TotalMessagesCount)
	prometheus.MustRegister(TotalTransformationError)
	prometheus.MustRegister(BufferCapacity)
	prometheus.MustRegister(BufferLength)
}

var TotalMessagesCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "message_total",
		Help: "Total number of messages.",
	},
	[]string{"type", "name"},
)

var TotalTransformationError = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "transformations_error_total",
		Help: "Total number of errors.",
	},
	[]string{"type", "name"},
)

var BufferCapacity = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "buffer_capacity",
		Help: "Buffer capacity.",
	},
	[]string{"type", "name"},
)

var BufferLength = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "buffer_length",
		Help: "Buffer length.",
	},
	[]string{"type", "name"},
)
