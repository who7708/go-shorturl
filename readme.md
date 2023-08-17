# 下载依赖
go mod tidy


# 编译
go build


# 运行
./app

# 访问

http://localhost:8000/api/shorten
请求
```json
{
    "url": "https://www.baidu.com",
    "expiration_in_minutes": 10
}
```

响应
```json
{
    "shortlink": "D"
}
```

http://localhost:8000/api/info?shortlink=xxx

返回
```json
{
    "url":"https://www.baidu.com",
    "create_at":"2023-08-17 21:17:12.34971 +0800 CST m=+62.079104201",
    "expiration_in_minutes":1
}
```

访问
http://localhost:8000/D

