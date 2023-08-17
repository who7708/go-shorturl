package main

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/mattheath/base62"
)

// 定义常量
const (

	// 全局自增ID
	URLIDKEY = "next.url.id"

	// 短链映射长链
	ShortlinkKey = "shortlink:%s:url"

	// 长链hash与短链进行映射
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
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})
	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}

	return &RedisCli{Cli: c}
}

func toSha1(url string) string {
	// panic("unimplemented")
	sha := sha1.New()
	return string(sha.Sum([]byte(url)))
}

// 长链生成短链
// Shorten(url string, expirationInMinutes int64) (string, error)
// // 短链详情
// ShortlinkInfo(shortUrl string) (interface{}, error)
// // 短链还原长链
// Unshortend(shortUrl string) (string, error)

// 长链生成短链
func (r *RedisCli) Shorten(url string, expirationInMinutes int64) (string, error) {
	h := toSha1(url)

	d, err := r.Cli.Get(fmt.Sprintf(URLHashKey, h)).Result()

	if err == redis.Nil {

	} else if err != nil {
		return "", err
	} else {
		if d == "{}" {
			// 过期了，nothing to do
		} else {
			return d, nil
		}
	}

	// increase the global counter
	err = r.Cli.Incr(URLIDKEY).Err()

	if err != nil {
		return "", err
	}

	id, err := r.Cli.Get(URLIDKEY).Int64()
	if err != nil {
		return "", err
	}

	shortUrl := base62.EncodeInt64(id)

	err = r.Cli.Set(fmt.Sprintf(ShortlinkKey, shortUrl), url,
		time.Minute*time.Duration(expirationInMinutes)).Err()

	if err != nil {
		return "", err
	}

	detail, err := json.Marshal(
		&URLDetail{
			URL:                 url,
			CreateAt:            time.Now().String(),
			ExpirationInMinutes: time.Duration(expirationInMinutes),
		},
	)

	if err != nil {
		return "", nil
	}

	err = r.Cli.Set(fmt.Sprintf(ShortlinkDetailKey, shortUrl), detail,
		time.Minute*time.Duration(expirationInMinutes)).Err()

	if err != nil {
		return "", err
	}

	return shortUrl, nil

}

// 短链详情
func (r *RedisCli) ShortlinkInfo(shortUrl string) (interface{}, error) {
	d, err := r.Cli.Get(fmt.Sprintf(ShortlinkDetailKey, shortUrl)).Result()

	if err == redis.Nil {
		return "", StatusError{404, errors.New("短链不存在")}
	} else if err != nil {
		return "", nil
	}
	return d, nil
}

func (r *RedisCli) Unshortend(shortUrl string) (string, error) {
	return "", nil
}
