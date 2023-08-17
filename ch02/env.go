package main

import (
	"log"
	"os"
	"strconv"
)

type Env struct {
	S Storage
}

func getEnv() *Env {
	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "192.168.1.3:6379"
	}

	passwd := os.Getenv("APP_REDIS_PASSWD")
	if passwd == "" {
		passwd = ""
	}

	dbS := os.Getenv("APP_REDIS_DB")
	if dbS == "" {
		dbS = "0"
	}

	dbIndex, err := strconv.Atoi(dbS)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("connect to redis (addr: %s ,passwd: %s ,db: %d)", addr, passwd, dbIndex)

	r := NewRedisCli(addr, passwd, dbIndex)

	return &Env{S: r}
}
