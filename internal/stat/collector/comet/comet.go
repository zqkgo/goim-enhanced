package comet

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/golang/glog"
	comet "github.com/zqkgo/goim-enhanced/api/comet/grpc"
	"github.com/zqkgo/goim-enhanced/internal/stat/collector"
)

type CometCollector struct {
	opts collector.Options

	mu      sync.RWMutex
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
	cc.ctx, cc.cancel = context.WithCancel(context.Background())
	err := cc.watchComet()
	if err != nil {
		return err
	}
	go cc.aggregate()
	return nil
}

func (cc *CometCollector) Stop() {
	cc.cancel()
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

func (cc *CometCollector) aggregate() {
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-cc.ctx.Done():
			log.Infof("CometCollector.aggregate() stop due to comet collector shutdown")
			return
		case <-t.C:
			var (
				wsOnlines   int64
				tcpOnlines  int64
				roomOnlines = make(map[string]int64)
				midOnlines  = make(map[int64]int64)
			)
			cc.mu.RLock()
			for _, c := range cc.clients {
				wsOnlines += c.wsOnlines
				tcpOnlines += c.tcpOnlines
				c.mu.RLock()
				for rid, ol := range c.roomOnlines {
					roomOnlines[rid] += ol
				}
				for mid, ol := range c.midOnlines {
					midOnlines[mid] += ol
				}
				c.mu.RUnlock()
			}
			cc.mu.RUnlock()
			if err := cc.opts.Dao.SetWSOnline(context.TODO(), wsOnlines); err != nil {
				log.Errorf("CometCollector.aggregate(), failed to set ws online, online: %d, err: %v", wsOnlines, err)
			}
			if err := cc.opts.Dao.SetTCPOnline(context.TODO(), tcpOnlines); err != nil {
				log.Errorf("CometCollector.aggregate(), failed to set tcp online, online: %d, err: %v", tcpOnlines, err)
			}
			for rid, ol := range roomOnlines {
				if err := cc.opts.Dao.SetRoomOnline(context.TODO(), rid, ol); err != nil {
					log.Errorf("CometCollector.aggregate(), failed to set room online, online: %d, rid: %s, err: %v", ol, rid, err)
				}
			}
			for mid, ol := range midOnlines {
				if err := cc.opts.Dao.SetMidOnline(context.TODO(), mid, ol); err != nil {
					log.Errorf("CometCollector.aggregate(), failed to set mid online, online: %d, mid: %d, err: %v", ol, mid, err)
				}
			}
		}
	}
}

func (cc *CometCollector) collect(cmt *cometClient) {
	cmt.ctx, cmt.cancel = context.WithCancel(context.Background())
	for {
		select {
		case <-cc.ctx.Done():
			log.Infof("CometCollector.collect() stop due to comet collector shutdown")
			return
		case <-cmt.ctx.Done():
			log.Infof("CometCollector.collect() stop due to comet offline, node: %s", cmt.addr)
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
			// memory cache
			cmt.wsOnlines = onlineRsp.WsOnline
			cmt.tcpOnlines = onlineRsp.TcpOnline
			cmt.mu.Lock()
			cmt.roomOnlines = onlineRsp.RoomOnlines
			cmt.midOnlines = onlineRsp.MidOnlines
			cmt.mu.Unlock()
		}
	}
}
