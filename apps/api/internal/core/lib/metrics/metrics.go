package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// EventBus metrics
	EventsPublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_published_total",
			Help: "Total number of events published to Redis Streams",
		},
		[]string{"event_name", "result"},
	)

	EventsPublishedDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "events_published_duration_seconds",
			Help:    "Duration of event publishing in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_name"},
	)

	// Consumer metrics
	EventsConsumedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_consumed_total",
			Help: "Total number of events consumed from Redis Streams",
		},
		[]string{"event_name", "result"},
	)

	HandlerDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "event_handler_duration_seconds",
			Help:    "Duration of event handler execution in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"event_name"},
	)

	EventsRetriedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "events_retried_total",
			Help: "Total number of events retried",
		},
		[]string{"event_name"},
	)

	// Outbox metrics
	OutboxEventsPending = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "outbox_events_pending",
			Help: "Number of events pending publication in outbox",
		},
	)

	OutboxEventsPublishedTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "outbox_events_published_total",
			Help: "Total number of events published from outbox",
		},
		[]string{"result"},
	)

	OutboxPublishDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "outbox_publish_duration_seconds",
			Help:    "Duration of outbox batch publish in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)

	OutboxRetriesTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "outbox_retries_total",
			Help: "Total number of outbox event retries",
		},
	)

	OutboxDLQTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "outbox_dlq_total",
			Help: "Total number of events moved to dead letter queue",
		},
	)
)
