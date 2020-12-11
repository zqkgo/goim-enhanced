package dao

import (
	"context"
	"fmt"

	"github.com/garyburd/redigo/redis"
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

func (d *Dao) GetCometHostOnlines(c context.Context) (map[string]int64, error) {
	conn := d.redis.Get()
	defer conn.Close()
	k := d.cometHostOnlineKey()
	m, err := redis.Int64Map(conn.Do("HGETALL", k))
	if err != nil {
		log.Errorf("conn.Do(HGETALL %s) err(%v)", k, err)
		return nil, err
	}
	return m, nil
}

func (d *Dao) SetWSOnline(c context.Context, ol int64) error {
	conn := d.redis.Get()
	defer conn.Close()
	k := d.cometWSOnlineKey()
	_, err := conn.Do("SET", k, ol)
	if err != nil {
		log.Errorf("conn.Do(SET %s,%d) err(%v)", k, ol, err)
		return err
	}
	return nil
}

func (d *Dao) SetTCPOnline(c context.Context, ol int64) error {
	conn := d.redis.Get()
	defer conn.Close()
	k := d.cometTCPOnlineKey()
	_, err := conn.Do("SET", k, ol)
	if err != nil {
		log.Errorf("conn.Do(SET %s,%d) err(%v)", k, ol, err)
		return err
	}
	return nil
}

func (d *Dao) SetRoomOnline(c context.Context, rid string, ol int64) error {
	conn := d.redis.Get()
	defer conn.Close()
	k := d.cometRoomOnlineKey(rid)
	_, err := conn.Do("SET", k, ol)
	if err != nil {
		log.Errorf("conn.Do(SET %s,%d) err(%v)", k, ol, err)
		return err
	}
	// TODO: add key to all room ids set
	return nil
}

func (d *Dao) SetMidOnline(c context.Context, mid int64, ol int64) error {
	conn := d.redis.Get()
	defer conn.Close()
	k := d.cometMidOnlineKey(mid)
	_, err := conn.Do("SET", k, ol)
	if err != nil {
		log.Errorf("conn.Do(SET %s,%d) err(%v)", k, ol, err)
		return err
	}
	// TODO: add key to all mids set
	return nil
}

func (d *Dao) cometHostOnlineKey() string {
	return "cmt:online:hosts"
}

func (d *Dao) cometWSOnlineKey() string {
	return "cmt:online:ws"
}

func (d *Dao) cometTCPOnlineKey() string {
	return "cmt:online:tcp"
}

func (d *Dao) cometRoomOnlineKey(rid string) string {
	return fmt.Sprintf("cmt:online:rid:%s", rid)
}

func (d *Dao) cometMidOnlineKey(mid int64) string {
	return fmt.Sprintf("cmt:online:mid:%d", mid)
}
