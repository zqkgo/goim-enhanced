package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	cometRPC "github.com/zqkgo/goim-enhanced/api/comet/grpc"
)

func main() {
	flag.Parse()
	if err := Init(); err != nil {
		panic(err)
	}
	cometMessages()
	cometOnlines()
}

func cometMessages() {
	comets, err := cometClients()
	if err != nil {
		panic(err)
	}
	for _, comet := range comets {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		reply, err := comet.client.Messages(ctx, &cometRPC.MessagesReq{})
		if err != nil {
			panic(err)
		}
		if reply.MsgStats != nil {
			for _, msgStat := range reply.MsgStats {
				fmt.Printf("addr: %s, type: %d, count: %d, speed: %f\n", comet.addr, msgStat.MsgType, msgStat.Count, msgStat.Speed)
			}
		}
	}
}

func cometOnlines() {
	comets, err := cometClients()
	if err != nil {
		panic(err)
	}
	for _, comet := range comets {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		reply, err := comet.client.Onlines(ctx, &cometRPC.OnlinesReq{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("addr: %s, host: %d, tcp: %d, websocket: %d\n", comet.addr, reply.HostOnline, reply.TcpOnline, reply.WsOnline)
		for mid, online := range reply.MidOnlines {
			fmt.Printf("addr: %s, mid: %d, online: %d\n", comet.addr, mid, online)
		}
		for rid, online := range reply.RoomOnlines {
			fmt.Printf("addr: %s, mid: %s, online: %d\n", comet.addr, rid, online)
		}
	}
}
