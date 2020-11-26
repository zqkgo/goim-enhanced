package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/golang/glog"
	"github.com/zqkgo/goim-enhanced/internal/stat/collector"
	"github.com/zqkgo/goim-enhanced/internal/stat/collector/comet"
	"github.com/zqkgo/goim-enhanced/internal/stat/conf"
	"github.com/zqkgo/goim-enhanced/internal/stat/dao"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic("failed to init config, err: " + err.Error())
	}
	dao := dao.NewDao(conf.Conf)
	cc := comet.NewCometCollector()
	err := cc.Init(
		comet.CometServiceName("goim.comet"),
		collector.Interval(conf.Conf.Collector.Itvl),
		collector.DiscoveryConf(conf.Conf.Discovery),
		collector.Dao(dao),
	)
	if err != nil {
		panic("failed to init comet collectors, err: " + err.Error())
	}
	collector.FireCollectors(cc)
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("goim-stat get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Infof("goim-stat exit")
			log.Flush()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
