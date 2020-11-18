package comet

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zqkgo/goim-enhanced/api/logic/grpc"
)

type stat struct {
	hostOnline    int64
	tcpOnline     int64
	wsOnline      int64
	roomOnlines   map[string]int64
	midOnlines    map[int64]int64
	broadcastMsgs uint64
	roomMsgs      uint64
	midMsgs       uint64
	mu            sync.RWMutex
	rstTime       time.Time
}

type MsgStat struct {
	MsgType grpc.PushMsg_Type
	Count   uint64
	Speed   float64
}

var DefaultStat = NewStat()

func NewStat() *stat {
	return &stat{
		roomOnlines: make(map[string]int64),
		midOnlines:  make(map[int64]int64),
		rstTime:     time.Now(),
	}
}

func (s *stat) IncrHostOnline() {
	atomic.AddInt64(&s.hostOnline, 1)
}

func (s *stat) DecrHostOnline() {
	atomic.AddInt64(&s.hostOnline, -1)
}

func (s *stat) IncrTCPOnline() {
	atomic.AddInt64(&s.tcpOnline, 1)
}

func (s *stat) DecrTCPOnline() {
	atomic.AddInt64(&s.tcpOnline, -1)
}

func (s *stat) IncrWsOnline() {
	atomic.AddInt64(&s.wsOnline, 1)
}

func (s *stat) DecrWsOnline() {
	atomic.AddInt64(&s.wsOnline, -1)
}

func (s *stat) IncrMidOnlines(mid int64) {
	s.mu.Lock()
	s.midOnlines[mid]++
	s.mu.Unlock()
}

func (s *stat) DecrMidOnlines(mid int64) {
	s.mu.Lock()
	s.midOnlines[mid]--
	s.mu.Unlock()
}

func (s *stat) IncrRoomOnlines(rid string) {
	s.mu.Lock()
	s.roomOnlines[rid]++
	s.mu.Unlock()
}

func (s *stat) DecrRoomOnlines(rid string) {
	s.mu.Lock()
	s.roomOnlines[rid]--
	s.mu.Unlock()
}

func (s *stat) IncrBroadcastMsgs() {
	atomic.AddUint64(&s.broadcastMsgs, 1)
}

func (s *stat) IncrRoomMsgs() {
	atomic.AddUint64(&s.roomMsgs, 1)
}

func (s *stat) IncrMidMsgs() {
	atomic.AddUint64(&s.midMsgs, 1)
}

func (s *stat) GetOnlines() (hostOnline, tcpOnline, wsOnline int64, roomOnlines map[string]int64, midOnlines map[int64]int64) {
	s.mu.RLock()
	roomOnlines = make(map[string]int64)
	for rid, online := range s.roomOnlines {
		roomOnlines[rid] = online
	}
	midOnlines = make(map[int64]int64)
	for mid, online := range s.midOnlines {
		midOnlines[mid] = online
	}
	s.mu.RUnlock()
	hostOnline, tcpOnline, wsOnline = s.hostOnline, s.tcpOnline, s.wsOnline
	return
}

func (s *stat) GetAndResetMsgs() []MsgStat {
	var (
		broadcast = MsgStat{
			MsgType: grpc.PushMsg_BROADCAST,
			Count:   atomic.LoadUint64(&s.broadcastMsgs),
		}
		room = MsgStat{
			MsgType: grpc.PushMsg_ROOM,
			Count:   atomic.LoadUint64(&s.roomMsgs),
		}
		mid = MsgStat{
			MsgType: grpc.PushMsg_PUSH,
			Count:   atomic.LoadUint64(&s.midMsgs),
		}
	)
	now := time.Now()
	if sec := now.Sub(s.rstTime).Seconds(); sec > 0 {
		broadcast.Speed = s.calSpd(broadcast.Count, sec)
		room.Speed = s.calSpd(room.Count, sec)
		mid.Speed = s.calSpd(mid.Count, sec)
	}
	// reset
	s.rstMsgs()
	s.rstTime = now
	return []MsgStat{broadcast, room, mid}
}

// round to .2f
func (s *stat) calSpd(cnt uint64, dur float64) float64 {
	t := float64(cnt) / dur
	spd := math.Round(t*100) / 100
	return spd
}

func (s *stat) rstMsgs() {
	atomic.StoreUint64(&s.broadcastMsgs, 0)
	atomic.StoreUint64(&s.roomMsgs, 0)
	atomic.StoreUint64(&s.midMsgs, 0)
}
