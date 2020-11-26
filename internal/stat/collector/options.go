package collector

import (
	"context"
	"time"

	"github.com/bilibili/discovery/naming"
	"github.com/zqkgo/goim-enhanced/internal/stat/dao"
	xtime "github.com/zqkgo/goim-enhanced/pkg/time"
)

const (
	DefaultItvl = xtime.Duration(3 * time.Second)
)

type Option func(o *Options)

type Options struct {
	Discovery *naming.Config
	Itvl      xtime.Duration
	Dao       *dao.Dao
	Context   context.Context
}

func Interval(itvl xtime.Duration) Option {
	return func(o *Options) {
		o.Itvl = itvl
	}
}

func DiscoveryConf(dis *naming.Config) Option {
	return func(o *Options) {
		o.Discovery = dis
	}
}

func Dao(d *dao.Dao) Option {
	return func(o *Options) {
		o.Dao = d
	}
}
