package main

import (
	"bet-server/src/server"
	"fmt"
)

func main() {
	var instance = server.NewServer()
	instance.Run()
	fmt.Println("Hello, world.")
}
