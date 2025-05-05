package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var iperf3Metrics = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "iperf3_metrics",
		Help: "Metrics collected from iperf3 client metric output",
	},
	[]string{"metric"},
)

func init() {
	prometheus.MustRegister(iperf3Metrics)
}

func NewIperf3Collector() *IPerfCollector {
	return &IPerfCollector{
		Desc: prometheus.NewDesc("iperf3_metric",
			"Iperf3 v2 metrics",
			[]string{"metric"},
			nil,
		),
	}
}

func main() {
	router := gin.Default()
	iperf3Collector := NewIperf3Collector()
	prometheus.MustRegister(iperf3Collector)
	router.GET("/metrics", PrometheusHandler())
	router.Run(fmt.Sprintf("localhost:%s", "9100"))
}

func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
