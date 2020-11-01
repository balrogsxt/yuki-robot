package cache

import (
	"context"
	"github.com/balrogsxt/xtbot-go/db"
	"github.com/balrogsxt/xtbot-go/util/entity"
	"time"
)

//实现一个redis缓存
type RedisCache struct {
	handle *db.RedisHandle
	ctx    context.Context
}

func (this *RedisCache) Init(config entity.UserConfig) error {
	this.ctx = context.Background()
	this.handle = &db.RedisHandle{
		Ctx: this.ctx,
	}
	if err := this.handle.ConnectRDB(config.Redis); err != nil {
		return err
	}
	return nil
}

func (this *RedisCache) Set(name string, val interface{}, duration ...time.Duration) error {
	var t time.Duration = 0
	if len(duration) > 0 {
		t = duration[0]
	}
	return this.handle.Rdb.Set(this.handle.Ctx, name, val, t).Err()
}
func (this *RedisCache) Get(name string) (string, error) {
	return this.handle.Rdb.Get(this.ctx, name).Result()
}
func (this *RedisCache) Exists(name string) bool {
	flag, err := this.handle.Rdb.Exists(this.ctx, name).Result()
	if err != nil {
		return false
	}
	if flag == 0 {
		return false
	}
	return true
}
func (this *RedisCache) GetMap(key, name string) (string, error) {
	return this.handle.Rdb.HGet(this.ctx, key, name).Result()
}
func (this *RedisCache) SetMap(key, name string, val interface{}) error {
	return this.handle.Rdb.HSet(this.ctx, key, name, val).Err()
}
func (this *RedisCache) ExistsMap(key, name string) bool {
	flag, err := this.handle.Rdb.HExists(this.ctx, key, name).Result()
	if err != nil {
		return false
	}
	return flag
}
