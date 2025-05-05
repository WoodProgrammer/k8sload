package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

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
	data, err := c.ParseMetricFile("metrics.json")

	if err != nil {
		log.Err(err).Msg("Error while parsing metrics file")
	}
	for _, v := range data.Intervals {
		for _, i := range v.Streams {
			var iperfMetric bytes.Buffer
			enc := gob.NewEncoder(&iperfMetric)
			err := enc.Encode(i)
			if err != nil {
				log.Err(err).Msg("encode error:")
			}

			fields := reflect.VisibleFields(reflect.TypeOf(i))

			b, err := json.Marshal(i)
			if err != nil {
				panic(err)
			}

			for _, f := range fields {
				var parsed map[string]interface{}
				if err := json.Unmarshal(b, &parsed); err != nil {
					log.Err(err).Msg("json.Unmarshal:")
				}
				fmt.Println("oc ", f.Name)
				fmt.Println(parsed["rtt"])

			}

		}
	}
}

func (c *IPerfCollector) ParseMetricFile(metricFile string) (Metric, error) {
	var metric Metric
	file, err := os.ReadFile(metricFile)
	if err != nil {
		log.Err(err).Msg("There is and error while reading file")
		return metric, err
	}
	err = json.Unmarshal(file, &metric)
	if err != nil {
		log.Err(err).Msg("Error while unmarshaling")
		return metric, err
	}

	return metric, err
}
