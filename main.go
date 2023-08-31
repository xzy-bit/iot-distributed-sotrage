package main

import (
	"IOT_Storage/src/Node"
	"time"
)

func main() {
	router := Node.Ping()
	go router.Run(":8080")
	for {
		time.Sleep(time.Second * 5)
		go Node.Send()
	}
	router_1 := Node.Login()
	router_1.Run(":8090")
}
