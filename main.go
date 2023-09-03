package main

import (
	"IOT_Storage/src/IOT_Device"
	"IOT_Storage/src/User"
)

func main() {
	//iot device process
	//file, _ := os.Open("private.pem")
	//if file == nil {
	//	Identity_Verify.GenerateKey(false)
	//}
	//IOT_Device.IotInit()
	nodes := []string{
		"http://192.168.42.129:8082",
	}
	IOT_Device.SendSliceToNode(nodes)

	//node process
	//router := Node.Ping()
	//go router.Run(":8080")
	//
	//routerReceivePublicKey := Controller.ReceivePublicKey()
	//go routerReceivePublicKey.Run(":8090")
	//
	//routerChallenge := Node.Challenge()
	//routerChallenge.Run(":8081")
	//routerGetSlice := Node.GetSlice()
	//routerGetSlice.Run(":8082")

	//user process
	//userIsAlive := User.Ping()
	//go userIsAlive.Run(":8080")
	//userReceiveKeys := User.ReceiveKeys()
	//userReceiveKeys.Run(":8090")
	User.SignForRandom("http://192.168.42.129:8081")
}
