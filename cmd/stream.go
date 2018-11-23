// Copyright © 2018 Barthelemy Vessemont
// This program is free software under the terms of the GNU GPLv3

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func convertMetricName(s *model.Sample) (string, error) {
	fmt.Println(s.Metric.String())
	return s.Metric.String(), nil
	return "", fmt.Errorf("Impossible to convert prometheus metric name to statsd format")
}

// streamCmd represents the stream command
var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Start the main scrape'n stream periodic job",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting scraping and streaming statsd metrics")

		log.Info("Initializing tickers")
		scrapePeriod, err := time.ParseDuration(pollingPeriod)
		if err != nil {
			log.Warning("Impossible to parse consulPeriod value, fallback to 120s")
			scrapePeriod = 120 * time.Second
		}
		if scrapePeriod < 60*time.Second {
			log.Warning("Scraping and streaming metric more than once a minute is not allowed, fallback to 60s")
			// scrapePeriod = 60 * time.Second
		}
		log.Info("Metrics scraping interval: ", scrapePeriod.String())
		scrapeMetricsTicker := time.NewTicker(scrapePeriod)

		log.Info("Initializing prometheus API client")
		promClient, err := api.NewClient(api.Config{Address: promURL, RoundTripper: api.DefaultRoundTripper})
		if err != nil {
			log.Fatal("Prometheus client init failed: ", err.Error())
		}
		q := prometheus.NewAPI(promClient)

		/* sd, err := statsd.NewBufferedClient(statsdAddr, "espoke", 300*time.Millisecond, 0)
		if err != nil {
			log.Fatal("Statsd forwarder init failed: ", err.Error())
		} */

		for {
			select {
			case <-scrapeMetricsTicker.C:
				log.Debug("Scraping metrics from Prometheus source")

				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				resp, err := q.Query(ctx, "es_node_search_latency{quantile=\"0_9\"}", time.Now())
				if err != nil {
					log.Error("Prometheus api query failed: ", err.Error())
					continue
				}

				switch {

				case resp.Type() == model.ValScalar:
					scalarVal := resp.(*model.Scalar)
					fmt.Println(scalarVal.String())

				case resp.Type() == model.ValVector:
					vectorVal := resp.(model.Vector)

					for _, elem := range vectorVal {
						convertMetricName(elem)
						// sd.SetInt(elem.Metric.String(), elem.Value, 1)
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(streamCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// streamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// streamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
