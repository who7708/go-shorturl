package main

func main() {
	a := App{}
	e := getEnv()
	a.Initialize(e)
	a.Run(":8000")
}

// export APP_REDIS_ADDR=192.168.1.3
// export APP_REDIS_PASSWD=123456
// export APP_REDIS_DB=1
