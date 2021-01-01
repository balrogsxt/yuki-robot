package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

//实现一个RobotCache的redis缓存
type RedisCache struct {
	rdb *redis.Client
	ctx    context.Context
}

func (this *RedisCache) Init(config *RobotConfig) error {
	this.ctx = context.Background()
	conf := config.Cache.Redis
	this.rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d",conf.Host,conf.Port),
		DB:       conf.Index,
		Password: conf.Password,
	})
	err := this.rdb.Ping(this.ctx).Err()
	if err != nil {
		return errors.New(fmt.Sprintf("连接Redis失败: %s",err.Error()))
	}
	return nil
}

func (this *RedisCache) Set(name string, val interface{}, duration ...time.Duration) error {
	var t time.Duration = 0
	if len(duration) > 0 {
		t = duration[0]
	}
	return this.rdb.Set(this.ctx, name, val, t).Err()
}
func (this *RedisCache) Get(name string) (string, error) {
	return this.rdb.Get(this.ctx, name).Result()
}
func (this *RedisCache) Exists(name string) bool {
	flag, err := this.rdb.Exists(this.ctx, name).Result()
	if err != nil {
		return false
	}
	if flag == 0 {
		return false
	}
	return true
}
func (this *RedisCache) GetMap(key, name string) (string, error) {
	return this.rdb.HGet(this.ctx, key, name).Result()
}
func (this *RedisCache) SetMap(key, name string, val interface{}) error {
	return this.rdb.HSet(this.ctx, key, name, val).Err()
}
func (this *RedisCache) ExistsMap(key, name string) bool {
	flag, err := this.rdb.HExists(this.ctx, key, name).Result()
	if err != nil {
		return false
	}
	return flag
}
