package comet

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/bilibili/discovery/naming"
	log "github.com/golang/glog"
	comet "github.com/zqkgo/goim-enhanced/api/comet/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	// grpc options
	grpcKeepAliveTime    = time.Duration(10) * time.Second
	grpcKeepAliveTimeout = time.Duration(3) * time.Second
	grpcBackoffMaxDelay  = time.Duration(3) * time.Second
	grpcMaxSendMsgSize   = 1 << 24
	grpcMaxCallMsgSize   = 1 << 24
)

const (
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
)

type cometClient struct {
	client comet.CometClient
	addr   string
	ticker *time.Ticker

	ctx    context.Context
	cancel context.CancelFunc
}

func (cc *CometCollector) watchComet() error {
	dis := naming.New(cc.opts.Discovery)
	resolver := dis.Build(cc.cometServiceName())
	event := resolver.Watch()
	select {
	case _, ok := <-event:
		if !ok {
			return fmt.Errorf("watchComet init failed")
		}
		if insInfo, ok := resolver.Fetch(); ok {
			if err := cc.updateClients(insInfo.Instances); err != nil {
				return err
			}

		}
	case <-time.After(10 * time.Second):
		return fmt.Errorf("watchComet init instances timeout")
	}
	go func() {
		for {
			if _, ok := <-event; !ok {
				log.Info("watchComet exit")
				return
			}
			ins, ok := resolver.Fetch()
			if ok {
				if err := cc.updateClients(ins.Instances); err != nil {
					log.Errorf("watchComet newAddress(%+v) error(%+v)", ins, err)
					continue
				}
				log.Infof("watchComet change newAddress:%+v", ins)
			}
		}
	}()
	return nil
}

func (cc *CometCollector) updateClients(insMap map[string][]*naming.Instance) error {
	ins := insMap[cc.opts.Discovery.Zone]
	if len(ins) == 0 {
		return fmt.Errorf("intances is empty")
	}
	clients := make(map[string]*cometClient)
	for _, in := range ins {
		if old, ok := cc.clients[in.Hostname]; ok {
			clients[in.Hostname] = old
			continue
		}
		// new comet, connect and collect
		client, err := cc.connect(in)
		if err != nil {
			return err
		}
		clients[in.Hostname] = client
		go cc.collect(client)
		log.Infof("new comet node added: %+v", in)
	}
	// stop the dropped comets
	for k, old := range cc.clients {
		if _, ok := clients[k]; !ok {
			old.cancel()
			log.Infof("drop comet: %s", k)
		}
	}
	cc.clients = clients
	return nil
}

func (cc *CometCollector) cometServiceName() string {
	defaultName := "goim.comet"
	if cc.opts.Context == nil {
		return defaultName
	}
	v := cc.opts.Context.Value(cometServiceNameKey{})
	if v == nil {
		return defaultName
	}
	return v.(string)
}

func (cc *CometCollector) connect(in *naming.Instance) (*cometClient, error) {
	var grpcAddr string
	for _, addrs := range in.Addrs {
		u, err := url.Parse(addrs)
		if err == nil && u.Scheme == "grpc" {
			grpcAddr = u.Host
		}
	}
	if grpcAddr == "" {
		return nil, fmt.Errorf("invalid grpc address:%v", in.Addrs)
	}
	cmt := &cometClient{
		addr:   grpcAddr,
		ticker: time.NewTicker(time.Duration(cc.opts.Itvl)),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second))
	defer cancel()
	conn, err := grpc.DialContext(ctx, grpcAddr,
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
		}...,
	)
	if err != nil {
		return nil, err
	}
	cmt.client = comet.NewCometClient(conn)
	return cmt, nil
}
