package cache

import (
	"context"

	"github.com/allegro/bigcache"
	"github.com/pkg/errors"
)

type bigCache struct {
	ins *bigcache.BigCache
}

func NewBigCacheAdaptor(instance *bigcache.BigCache) *bigCache {
	ins := &bigCache{
		ins: instance,
	}
	return ins
}

func (c *bigCache) Get(_ context.Context, key string) ([]byte, error) {
	if c == nil {
		return nil, errors.New("free xcache is disabled")
	}
	value, err := c.ins.Get(key)
	return value, err
}

func (c *bigCache) Set(_ context.Context, key string, value []byte) error {
	if c == nil {
		return errors.New("free xcache is disabled")
	}
	err := c.ins.Set(key, value)
	return err
}

func (c *bigCache) Delete(_ context.Context, key string) error {
	if c == nil {
		return errors.New("inmem xcache is disabled")
	}
	err := c.ins.Delete(key)
	return err
}

func (c *bigCache) Clear(_ context.Context) error {
	if c == nil {
		return errors.New("inmem xcache is disabled")
	}
	err := c.ins.Reset()
	return err
}
