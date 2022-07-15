package main

import (
	"context"
	"fmt"
)

type Code struct {
	Code int
}

func (c *RedisClient) Create(code int) (*Code, error) {
	mkey := newRedisKey(code)
	id := Code{
		Code: code,
	}
	err := c.Set(context.Background(), mkey, code, c.TTL).Err()
	if err != nil {
		return nil, fmt.Errorf("redis: set key %q: %w", mkey, err)
	}
	return &id, nil
}

func (c *RedisClient) Check(code int) (bool, error) {
	mkey := newRedisKey(code)
	data, err := c.GetRecord(mkey)
	if err != nil {
		return false, fmt.Errorf("redis: get record by key %q: %w", mkey, err)
	} else if data == nil {
		// add here custom err handling
		return false, nil
	}
	return true, nil
}

func (c *RedisClient) Delete(code int) error {
	mkey := newRedisKey(code)
	err := c.Del(context.Background(), mkey).Err()
	if err != nil {
		return fmt.Errorf("redis: trying to delete value by key %q: %w", mkey, err)
	}
	return nil
}

func newRedisKey(code int) string {
	return fmt.Sprintf("sms.Code.%d", code)
}
