package main

import "log"

func main() {
	err := parseFlags()
	if err != nil {
		log.Fatalln(err)
	}

	bootstrap()
}
