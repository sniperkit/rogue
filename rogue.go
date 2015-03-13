// 2015 Jamie Alquiza
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"./outputs/elasticsearch"
)

var config struct {
	addr     string
	port     string
	queuecap int
	elasticsearch string
}

var (
	messageIncomingQueue = make(chan []byte, config.queuecap)
	sig_chan = make(chan os.Signal)
)

func init() {
	flag.StringVar(&config.addr, "listen-addr", "localhost", "bind address")
	flag.StringVar(&config.port, "listen-port", "6030", "bind port")
	flag.IntVar(&config.queuecap, "queue-cap", 100, "In-flight message queue capacity")
	flag.StringVar(&config.elasticsearch, "elasticsearch", "http://localhost:9200", "ElasticSearch IP")
	flag.Parse()
	// Update vars that depend on flag inputs.
	messageIncomingQueue = make(chan []byte, config.queuecap)
}

// Handles signal events.
// Currently just kills service, will eventually perform graceful shutdown.
func runControl() {
	signal.Notify(sig_chan, syscall.SIGINT)
	<-sig_chan
	log.Printf("Rogue shutting down")
	os.Exit(0)
}

func main() {
	// Start internals.
	go listenTcp()

	// We pass a the message feed queue in addition to
	// flush timeout, message batch count and size threshold values.
	go elasticsearch.Run(config.elasticsearch,
		messageIncomingQueue,
		5, 1000, 100000)

	runControl()
}