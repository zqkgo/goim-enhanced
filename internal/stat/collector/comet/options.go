package comet

import (
	"context"

	"github.com/zqkgo/goim-enhanced/internal/stat/collector"
)

type cometServiceNameKey struct{}

func CometServiceName(name string) collector.Option {
	return func(o *collector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		context.WithValue(o.Context, cometServiceNameKey{}, name)
	}
}
