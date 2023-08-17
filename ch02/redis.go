package main

import (
	"time"

	"github.com/go-redis/redis"
)

// 定义常量
const (

	// 生成短链key
	URLIDKEY = "next.url.id"

	// 短链映射长链
	ShortlinkKey = "shortlink:%s:url"

	// 短链hash与地址进行映射
	URLHashKey = "urlhash:%s:url"

	// 短链详情
	ShortlinkDetailKey = "shortlink:%s:detail"
)

// redis 客户端数据结构
type RedisCli struct {
	Cli *redis.Client
}

type URLDetail struct {
	URL                 string        `json:"url"`
	CreateAt            string        `json:"create_at"`
	ExpirationInMinutes time.Duration `json:"expiration_in_minutes"`
}

// 初始化redis客户端
func NewRedisCli(addr string, passwd string, db int) *RedisCli {
	c := redis.NewClient(opt * redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}

	return &RedisCli{Cli: c}
}
