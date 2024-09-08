package metrics

import (
	"reflect"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestMustRegisterGauge(t *testing.T) {
	type args struct {
		subsystem  string
		name       string
		help       string
		labelNames []string
	}
	tests := []struct {
		name   string
		args   args
		labels prometheus.Labels
		want   *prometheus.Desc
	}{
		{
			name: "Test register new gauge vec",
			args: args{
				subsystem:  "example",
				name:       "example_metrcis",
				help:       "help",
				labelNames: []string{"name"},
			},
			labels: prometheus.Labels{"name": "example"},
			want:   prometheus.NewDesc("example", "help", []string{"name"}, prometheus.Labels{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MustRegisterGauge(
				tt.args.subsystem,
				tt.args.name,
				tt.args.help,
				tt.args.labelNames...)
			defer prometheus.Unregister(got)

			if !reflect.DeepEqual(got.With(tt.labels).Desc(), tt.want) {
				return
			}
		})
	}
}

func TestMustRegisterCounter(t *testing.T) {
	type args struct {
		subsystem  string
		name       string
		help       string
		labelNames []string
	}
	tests := []struct {
		name string
		args args
		want *prometheus.Desc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MustRegisterCounter(
				tt.args.subsystem,
				tt.args.name,
				tt.args.help,
				tt.args.labelNames...)
			defer prometheus.Unregister(got)
			require.Equal(t, tt.want, got, "MustRegisterCounter() = %v, want %v", got, tt.want)
		})
	}
}

func TestMustRegisterHistogram(t *testing.T) {
	type args struct {
		subsystem  string
		name       string
		help       string
		buckets    []float64
		labelNames []string
	}
	tests := []struct {
		name string
		args args
		want *prometheus.Desc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MustRegisterHistogram(
				tt.args.subsystem,
				tt.args.name,
				tt.args.help,
				tt.args.buckets,
				tt.args.labelNames...)
			defer prometheus.Unregister(got)
			require.Equal(t, tt.want, got, "MustRegisterHistogram() = %v, want %v", got, tt.want)
		})
	}
}
