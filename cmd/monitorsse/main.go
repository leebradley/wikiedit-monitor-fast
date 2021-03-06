package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/leebradley/wikiedit-monitor-fast/pkg/monitorsse"
	"github.com/leebradley/wikiedit-monitor-fast/pkg/wiki/recentchanges"
	nats "github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		natsurl  string
		hidebots bool
		wikis    string
	)

	flag.StringVar(&natsurl, "natsurl", nats.DefaultURL, "the url used to connect to nats")
	flag.BoolVar(&hidebots, "hidebots", true, "Whether to hide / ignore bot edits")
	flag.StringVar(&wikis, "wikis", "en", "A comma-delimited list of wikis to listen to")
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	logger := logrus.New()
	logger.Info("Starting wikimedia sse monitor")

	natsconn, err := nats.Connect(natsurl)
	if err != nil {
		logger.WithError(err).Fatal("Could not connect to nats")
	}

	lo := recentchanges.ListenOptions{
		Hidebots: hidebots,
		Wikis:    strings.Split(wikis, ","),
	}

	forward := monitorsse.NewForwarder(natsconn, logger)
	forward.Forward(lo, monitorsse.DefaultForwardSubj)

	done := make(chan struct{})

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			return
		}
	}
}
