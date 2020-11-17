package comet

import (
	"sync"
	"sync/atomic"
)

type stat struct {
	hostOnline    int64
	tcpOnline     int64
	wsOnline      int64
	roomOnlines   map[string]int64
	midOnlines    map[int64]int64
	allMsgs       uint64
	broadcastMsgs uint64
	roomMsgs      uint64
	midMsgs       uint64
	mu            sync.RWMutex
}

var DefaultStat = NewStat()

func NewStat() *stat {
	return &stat{
		roomOnlines: make(map[string]int64),
		midOnlines:  make(map[int64]int64),
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

func (s *stat) IncrAllMsgs() {
	atomic.AddUint64(&s.allMsgs, 1)
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
