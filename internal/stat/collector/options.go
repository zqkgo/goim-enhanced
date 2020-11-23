package collector

import (
	"context"
	"time"
)

type Option func(o *Options)

type Options struct {
	Itvl         time.Duration
	MaxLockRetry uint64
	Context      context.Context
}

func Interval(itvl time.Duration) Option {
	return func(o *Options) {
		o.Itvl = itvl
	}
}

func MaxLockRetry(retry uint64) Option {
	return func(o *Options) {
		o.MaxLockRetry = retry
	}
}
