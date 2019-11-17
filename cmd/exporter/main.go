package main

import (
	"context"
	"exporter"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/citradigital/toldata"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	natsURL := os.Getenv("NATS_URL")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println(fmt.Sprintf("NATS_URL: %s", natsURL))
	client, err := toldata.NewBus(ctx, toldata.ServiceConfiguration{URL: natsURL})
	if err != nil {
		log.Fatalln("unable-to-contact-toldata")
	}
	defer client.Close()

	exporter.NewUpCollector(client)
	flag.Parse()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
