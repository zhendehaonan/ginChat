package utils

import (
	"context"
	"fmt"
	"ginchat/config"
	"github.com/go-redis/redis/v8"
)

// Redis
type UserRedis struct {
	*redis.Client
}

// 创建操作redis的缓存对象
func NewUserRedisBasic(ctx context.Context) *UserRedis {
	return &UserRedis{config.NewRedisClient(ctx)}
}

const (
	PublishKey = "websocket"
)

// Publish 发布消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	red := NewUserRedisBasic(ctx)
	var err error
	err = red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅Redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	red := NewUserRedisBasic(ctx)
	sub := red.PSubscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return msg.Payload, err
}
