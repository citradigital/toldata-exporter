package exporter

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/citradigital/toldata"
	"github.com/gogo/protobuf/proto"

	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	Bus *toldata.Bus
}

var (
	labelNames = []string{"name"}

	endPoints []string

	upDesc = prometheus.NewDesc(
		prometheus.BuildFQName(NAMESPACE, "services", "up"),
		"Service Up", labelNames, nil,
	)
)

func init() {
	//	register("toldata", NewUpCollector())
}

func NewUpCollector(bus *toldata.Bus) (Collector, error) {
	c := Collector{
		Bus: bus,
	}

	b, err := ioutil.ReadFile("endpoints.txt")
	if err != nil {
		fmt.Print(err)
	}

	endPoints = strings.Split(string(b), "\n")

	prometheus.MustRegister(c)
	return c, nil
}

func (c *Collector) healthCheck(functionName string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reqRaw, err := proto.Marshal(nil)

	result, err := c.Bus.Connection.RequestWithContext(ctx, functionName, reqRaw)
	if err != nil {
		return err
	}

	if result.Data[0] == 0 {
		// 0 means no error
		p := &toldata.ToldataHealthCheckInfo{}
		err = proto.Unmarshal(result.Data[1:], p)
		if err != nil {
			return err
		}
		return nil
	} else {
		var pErr toldata.ErrorMessage
		err = proto.Unmarshal(result.Data[1:], &pErr)
		if err == nil {
			return errors.New(pErr.ErrorMessage)
		} else {
			return err
		}
	}
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upDesc
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {

	for _, endPoint := range endPoints {
		if endPoint == "" {
			continue
		}
		value := 0.0
		var labels = []string{endPoint}

		err := c.healthCheck(endPoint + "/ToldataHealthCheck")
		if err == nil {
			value = 1
		} else {
			log.Println(err)
		}

		ch <- prometheus.MustNewConstMetric(upDesc, prometheus.GaugeValue, value, labels...)
	}
}
