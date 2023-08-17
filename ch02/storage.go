package main

type Storage interface {
	// 长链生成短链
	Shorten(url string, expirationInMinutes int64) (string, error)
	// 短链详情
	ShortlinkInfo(shortUrl string) (interface{}, error)
	// 短链还原长链
	Unshortend(shortUrl string) (string, error)
}
