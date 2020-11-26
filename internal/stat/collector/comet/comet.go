package comet

import (
	"context"
	"fmt"

	log "github.com/golang/glog"
	comet "github.com/zqkgo/goim-enhanced/api/comet/grpc"
	"github.com/zqkgo/goim-enhanced/internal/stat/collector"
)

type CometCollector struct {
	opts    collector.Options
	clients map[string]*cometClient

	ctx    context.Context
	cancel context.CancelFunc
}

func NewCometCollector() *CometCollector {
	return &CometCollector{}
}

func (cc *CometCollector) Init(opts ...collector.Option) error {
	for _, opt := range opts {
		opt(&cc.opts)
	}
	if cc.opts.Dao == nil {
		return fmt.Errorf("Dao must be present")
	}
	if cc.opts.Itvl == 0 {
		cc.opts.Itvl = collector.DefaultItvl
	}
	return nil
}

func (cc *CometCollector) Collect() error {
	err := cc.watchComet()
	if err != nil {
		return err
	}
	return nil
}

func (cc *CometCollector) Stop() {

}

func (cc *CometCollector) ReCollect(opts ...collector.Option) error {
	return nil
}

func (cc *CometCollector) Options() collector.Options {
	return cc.opts
}

func (cc *CometCollector) String() string {
	return "comet collector"
}

func (cc *CometCollector) collect(cmt *cometClient) {
	cmt.ctx, cmt.cancel = context.WithCancel(context.Background())
	for {
		select {
		case <-cmt.ctx.Done():
			log.Infof("collect stop, node: %s", cmt.addr)
			return
		case <-cmt.ticker.C:
			log.Infof("collecting comet stats")
			onlineRsp, err := cmt.client.Onlines(context.Background(), &comet.OnlinesReq{})
			if err != nil {
				log.Errorf("CometCollector.collect(), failed to get onlines, addr: %s, err: %v", cmt.addr, err)
			}
			err = cc.opts.Dao.SetCometHostOnline(context.TODO(), cmt.addr, onlineRsp.HostOnline)
			if err != nil {
				log.Errorf("CometCollector.collect(), failed to set host online, host: %s, online: %d", cmt.addr, onlineRsp.HostOnline)
			}
		}
	}
}
