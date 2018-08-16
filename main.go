package main

import (
	"flag"
	"log"
)

var (
	configPath string
	singleRun  bool
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "path to the configuration file")
	flag.BoolVar(&singleRun, "single", false, "run only once (to be used when controlled via cron or simiar)")
}

func main() {
	flag.Parse()

	n, err := PrepareNotifiers()
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewChecker(configPath, n)
	if err != nil {
		log.Fatal(err)
	}

	c.Run(singleRun)
}
