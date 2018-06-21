package main

import "log"

func main() {
	n, err := PrepareNotifiers()
	if err != nil {
		log.Fatal(err)
	}

	c, err := NewChecker("./config.yaml", n)
	if err != nil {
		log.Fatal(err)
	}

	c.Run()
}
