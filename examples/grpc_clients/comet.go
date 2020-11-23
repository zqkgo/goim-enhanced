package main

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

type CometClient struct {
	client comet.CometClient
	addr   string
}

func cometClients() ([]CometClient, error) {
	dis := naming.New(Conf.DiscoveryConf)
	resolver := dis.Build("goim.comet")
	event := resolver.Watch()
	var comets []CometClient
	select {
	case _, ok := <-event:
		if !ok {
			return nil, fmt.Errorf("watchComet init failed")
		}
		if insInfo, ok := resolver.Fetch(); ok {
			ins := insInfo.Instances[zone]
			for _, in := range ins {
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
				c := comet.NewCometClient(conn)
				comets = append(comets, CometClient{client: c, addr: grpcAddr})
			}
			log.Infof("watchComet init newAddress:%+v", ins)
		}
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("watchComet init instances timeout")
	}
	return comets, nil
}
