package exporter

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/zqkgo/goim-enhanced/internal/stat/dao"
)

type Exporter struct {
	hostOnline *prometheus.Desc
	wsOnline   *prometheus.Desc
	tcpOnline  *prometheus.Desc
	roomOnline *prometheus.Desc
	midOnline  *prometheus.Desc

	dao *dao.Dao
}

func NewExporter(d *dao.Dao) *Exporter {
	return &Exporter{
		hostOnline: prometheus.NewDesc(
			"goim_comet_hostOnline",
			"goim comet host online",
			[]string{"host"}, nil),
		roomOnline: prometheus.NewDesc(
			"goim_comet_roomOnline",
			"goim comet room online",
			[]string{"roomID"}, nil),

		dao: d,
	}
}

func (c *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.hostOnline
	ch <- c.roomOnline
}

func (c *Exporter) Collect(ch chan<- prometheus.Metric) {
	m, err := c.dao.GetCometHostOnlines(context.TODO())
	if err != nil {
		panic(err)
	}
	for host, ol := range m {
		ch <- prometheus.MustNewConstMetric(c.hostOnline, prometheus.GaugeValue, float64(ol), host)
	}
}
