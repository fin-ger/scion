// Copyright 2017 ETH Zurich
// Copyright 2018 ETH Zurich, Anapaya Systems
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package prom contains some utility functions for dealing with prometheus
// metrics.
package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Common label values.
const (
	// LabelResult is the label for result classifications.
	LabelResult = "result"
	// LabelStatus for latency status classifications, possible values are prefixed with Status*.
	LabelStatus = "status"
	// LabelOperation is the label for the name of an executed operation.
	LabelOperation = "op"
	// LabelSrc is the label for the src of a request.
	LabelSrc = "src"
)

// Common result values.
const (
	// Success is no error.
	Success = "ok_success"
	// ErrCrypto is used for crypto related errors.
	ErrCrypto = "err_crypto"
	// ErrDB is used for db related errors.
	ErrDB = "err_db"
	// ErrInternal is an internal error.
	ErrInternal = "err_internal"
	// ErrInvalidReq is an invalid request.
	ErrInvalidReq = "err_invalid_request"
	// ErrNotClassified is an error that is not further classified.
	ErrNotClassified = "err_not_classified"
	// ErrParse failed to parse request.
	ErrParse = "err_parse"
	// ErrProcess is an error during processing e.g. parsing failed.
	ErrProcess = "err_process"
	// ErrTimeout is a timeout error.
	ErrTimeout = "err_timeout"
	// ErrValidate is used for validation related errors.
	ErrValidate = "err_validate"
	// ErrVerify is used for validation related errors.
	ErrVerify = "err_verify"
	// ErrReply is used for errors when sending the reply.
	ErrReply = "err_reply"
)

// FIXME(roosd): remove when moving messenger to new metrics style.
const (
	StatusOk      = "ok"
	StatusErr     = "err"
	StatusTimeout = "err_timeout"
)

var (
	// DefaultLatencyBuckets 10ms, 20ms, 40ms, ... 5.12s, 10.24s.
	DefaultLatencyBuckets = []float64{0.01, 0.02, 0.04, 0.08, 0.16, 0.32, 0.64,
		1.28, 2.56, 5.12, 10.24}
)

// ExportElementID exports the element ID as configured in the config file.
func ExportElementID(id string) {
	NewGaugeVec("scion", "", "elem_id",
		"The element ID from the config file", []string{"cfg"}).WithLabelValues(id).Set(1)
}

// FIXME(roosd): remove.
func CopyLabels(labels prometheus.Labels) prometheus.Labels {
	l := make(prometheus.Labels)
	for k, v := range labels {
		l[k] = v
	}
	return l
}

// NewCounter creates a new prometheus counter that is registered with the default registry.
func NewCounter(namespace, subsystem, name, help string) prometheus.Counter {
	return promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
	)
}

// NewCounterVec creates a new prometheus counter vec that is registered with the default registry.
func NewCounterVec(namespace, subsystem, name, help string,
	labelNames []string) *prometheus.CounterVec {

	return promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		labelNames,
	)
}

// NewGauge creates a new prometheus gauge that is registered with the default registry.
func NewGauge(namespace, subsystem, name, help string) prometheus.Gauge {
	return promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
	)
}

// NewGaugeVec creates a new prometheus gauge vec that is registered with the default registry.
func NewGaugeVec(namespace, subsystem, name, help string,
	labelNames []string) *prometheus.GaugeVec {

	return promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
		},
		labelNames,
	)
}

// NewHistogramVec creates a new prometheus histogram vec
// that is registered with the default registry.
func NewHistogramVec(namespace, subsystem, name, help string,
	labelNames []string, buckets []float64) *prometheus.HistogramVec {

	return promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      name,
			Help:      help,
			Buckets:   buckets,
		},
		labelNames,
	)
}
