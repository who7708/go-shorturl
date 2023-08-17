package main

func main() {
	a := App{}
	e := getEnv()
	a.Initialize(e)
	a.Run(":8000")
}
