// Copyright © 2018 Barthelemy Vessemont
// This program is free software under the terms of the GNU GPLv3

package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var loglevel string
var pollingPeriod string
var promURL string
var statsdAddr string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "prom2statsd_streamer",
	Short: "Periodically scrapes a subset of prometheus metrics and send them to a statsd target",
	Long: `Periodically scrapes a subset of prometheus metrics, format them, 
and send them to a specified statsd target`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.prom2statsd_streamer.yaml)")
	rootCmd.PersistentFlags().StringVarP(&promURL, "prometheusUrl", "s", "http://127.0.0.1:9090", "prometheus source api host:port")
	rootCmd.PersistentFlags().StringVarP(&statsdAddr, "statsdTarget", "d", "http://127.0.0.1:8125", "statsd destination host:port")
	rootCmd.PersistentFlags().StringVar(&pollingPeriod, "scrapePeriod", "60s", "metrics scraping interval")
	rootCmd.PersistentFlags().StringVarP(&loglevel, "loglevel", "l", "info", "log level")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".prom2statsd_streamer" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".prom2statsd_streamer")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}

	// Init logger
	log.SetOutput(os.Stdout)
	lvl, err := log.ParseLevel(loglevel)
	if err != nil {
		log.Warning("Log level not recognized, fallback to default level (INFO)")
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)
	log.Info("Logger initialized")
}
