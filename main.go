package main

import (
	"flag"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "path to the configuration file")
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

	c.Run()
}
