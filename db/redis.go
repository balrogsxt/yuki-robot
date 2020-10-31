package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/balrogsxt/xtbot-go/util/entity"
	"github.com/go-redis/redis/v8"
)

//这里初始化redis数据库连接
type RedisHandle struct {
	Rdb *redis.Client
	Ctx context.Context
}

func (this *RedisHandle) ConnectRDB(config entity.RedisConfig) error {
	this.Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		DB:       config.Index,
		Password: config.Password,
	})
	err := this.Rdb.Ping(this.Ctx).Err()
	if err != nil {
		return errors.New("连接Redis失败:" + err.Error())
	}
	return nil
}
