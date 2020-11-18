package comet

import (
	"sync"
	"testing"
	"time"
)

func TestStat(t *testing.T) {
	stat := NewStat()
	stat.rstTime = time.Now().Add(-time.Second * 10)
	for i := 0; i < 10000; i++ {
		stat.IncrBroadcastMsgs()
		stat.IncrRoomMsgs()
		stat.IncrMidMsgs()
	}
	stat.GetAndResetMsgs()
	if stat.broadcastMsgs != 0 {
		t.Fatalf("want 0, got %d", stat.broadcastMsgs)
	}
	rid := "foo"
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 500; i++ {
		go func() {
			defer wg.Done()
			stat.IncrRoomOnlines(rid)
		}()
		go func() {
			defer wg.Done()
			stat.DecrRoomOnlines(rid)
		}()
	}
	wg.Wait()
}
