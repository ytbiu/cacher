package cac

import (
	"sync"
	"time"
)

type cac struct {
	sync.RWMutex
	key         interface{}
	val         interface{}
	expire      time.Duration
	expireCb    func()
	expireTimer *time.Timer
}

func newCac(key string, val interface{}, expire time.Duration, deleteFn func(key string), cbs ...func()) *cac {
	c := &cac{
		key:    key,
		val:    val,
		expire: expire,
	}

	c.expireTimer = time.AfterFunc(expire, func() {
		deleteFn(key)
		for _, cb := range cbs {
			cb()
		}
	})
	return c
}

func (c *cac) reset(expire time.Duration) {
	c.expire = expire
	c.expireTimer.Reset(expire)
}

func (c *cac) value() interface{} {
	return c.val
}

func (c *cac) put(val interface{}) {
	c.val = val
}
