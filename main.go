package main

import (
	"flag"

	log "github.com/liudanking/goutil/logutil"
)

func main() {
	flag.Parse()

	log.Info("ok")
}
