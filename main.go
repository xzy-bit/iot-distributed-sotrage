package main

import (
	"IOT_Storage/src/User"
)

func main() {
	//iot device process
	//file, _ := os.Open("private.pem")
	//if file == nil {
	//	Identity_Verify.GenerateKey(false)
	//}
	//IOT_Device.IotInit()

	////node process
	//router := Node.Ping()
	//go router.Run(":8080")
	//
	//routerReceivePublicKey := Controller.ReceivePublicKey()
	//go routerReceivePublicKey.Run(":8090")
	//
	//routerChallenge := Node.Challenge()
	//routerChallenge.Run(":8081")

	////user process
	////userIsAlive := User.Ping()
	////go userIsAlive.Run(":8080")
	////userReceiveKeys := User.ReceiveKeys()
	////userReceiveKeys.Run(":8090")
	User.SignForRandom("http://localhost:8081")
}
