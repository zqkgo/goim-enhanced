package dao

import (
	"context"

	log "github.com/golang/glog"
)

func (d *Dao) SetCometHostOnline(c context.Context, host string, ol int64) error {
	conn := d.redis.Get()
	defer conn.Close()
	k := d.cometHostOnlineKey()
	_, err := conn.Do("HSET", k, host, ol)
	if err != nil {
		log.Errorf("conn.Do(HSET %s,%s,%d) err(%v)", k, host, ol, err)
		return err
	}
	return nil
}

// TODO: save other stats

func (d *Dao) cometHostOnlineKey() string {
	return "cmt:online:hosts"
}
