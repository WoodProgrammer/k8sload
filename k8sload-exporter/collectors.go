package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

type IPerfCollector struct {
	Desc *prometheus.Desc
}

func (c *IPerfCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Desc
}

func (c *IPerfCollector) Collect(ch chan<- prometheus.Metric) {
	var intervals []Intervals
	var testCaseName string

	testCaseName = os.Getenv("LOAD_TEST_NAME")
	if len(testCaseName) == 0 {
		testCaseName = "default_load_test"
	}

	data, err := c.ParseMetricFile("metrics.json")

	if err != nil {
		log.Err(err).Msg("Error while parsing metrics file")
	}
	var mainJson map[string]interface{}
	if err := json.Unmarshal(data, &mainJson); err != nil {
		log.Err(err).Msg("json.Unmarshal:")
	}

	for k, v := range mainJson {
		term := 0
		if k == "intervals" {
			b, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			if err := json.Unmarshal(b, &intervals); err != nil {
				log.Err(err).Msg("json.Unmarshal:")
			}
			for _, s := range intervals {

				for _, t := range s.Streams {
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.Bytes), fmt.Sprintf("bytes_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.BitsPerSecond), fmt.Sprintf("bits_per_second_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.End), fmt.Sprintf("end_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.PMTU), fmt.Sprintf("pmtu_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.RTT), fmt.Sprintf("rtt_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.Retransmits), fmt.Sprintf("retransmits_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.Seconds), fmt.Sprintf("seconds_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.SndCwd), fmt.Sprintf("sndcwd_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.SndWnd), fmt.Sprintf("sndwnd_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.Socket), fmt.Sprintf("sockets_streams_term_%d", term), testCaseName)
					ch <- prometheus.MustNewConstMetric(c.Desc, prometheus.GaugeValue, float64(t.Start), fmt.Sprintf("start_streams_term_%d", term), testCaseName)

					fmt.Println(term)
				}
				term = term + 1
			}
		}
	}
}

func (c *IPerfCollector) ParseMetricFile(metricFile string) ([]byte, error) {
	file, err := os.ReadFile(metricFile)
	if err != nil {
		log.Err(err).Msg("There is and error while reading file")
		return nil, err
	}

	return file, err
}
