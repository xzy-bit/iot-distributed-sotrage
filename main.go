package main

import (
	"IOT_Storage/src/IOT_Device"
	"IOT_Storage/src/Identity_Verify"
	"os"
)

func main() {
	//iot device process
	file, _ := os.Open("private.pem")
	if file == nil {
		Identity_Verify.GenerateKey(false)
	}
	IOT_Device.IotInit()

	////server process
	//router := Node.Ping()
	//go router.Run(":8080")

	//routerReceivePublicKey := Controller.ReceivePublicKey()
	//routerReceivePublicKey.Run(":8090")

	////user process
	//userIsAlive := User.Ping()
	//go userIsAlive.Run(":8080")
	//userReceiveKeys := User.ReceiveKeys()
	//userReceiveKeys.Run(":8090")
}
