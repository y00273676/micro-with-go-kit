package fuse

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/afex/hystrix-go/plugins"

	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"log"
	"net"
	"net/http"
)

func Init() {
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(net.JoinHostPort("", "81"), hystrixStreamHandler)
	c, err := plugins.InitializeStatsdCollector(&plugins.StatsdCollectorConfig{
		StatsdAddr: "localhost:8125",
		Prefix:     "myapp.hystrix",
	})
	if err != nil {
		log.Fatalf("could not initialize statsd client: %v", err)
	}

	metricCollector.Registry.Register(c.NewStatsdCollector)
}
