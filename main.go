package main

import (
	"IOT_Storage/src/Node"
)

func main() {
	//iot device process
	//file, _ := os.Open("private.pem")
	//if file == nil {
	//	Identity_Verify.GenerateKey(false)
	//}
	//IOT_Device.IotInit()
	//nodes := []string{
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//	"http://192.168.42.129",
	//}
	//portForSlice := 10080
	//portForIndex := 8040
	//IOT_Device.SendSliceAndIndexToNode(nodes, portForSlice, portForIndex)

	//node process
	//router := Node.Ping()
	//go router.Run(":8080")
	//
	//routerReceivePublicKey := Controller.ReceivePublicKey()
	//go routerReceivePublicKey.Run(":8090")
	//routerChallenge := Node.Challenge()
	//routerChallenge.Run(":8081")
	//routerGetSlice := Node.ServerGetSlice()
	//routerGetSlice.Run(":8081")

	//user process
	//userIsAlive := User.Ping()
	//go userIsAlive.Run(":8080")
	//userReceiveKeys := User.ReceiveKeys()
	//userReceiveKeys.Run(":8090")
	//User.SignForRandom("http://192.168.42.129:8081")
	//P2P_Net.P2pPing()
	//SearchableEncrypt.GenerateSk()

	//query := []string{
	//	"精神科",
	//	"食欲不振",
	//	"记忆力衰退",
	//}
	//
	//queryCompare := []string{
	//	"精神科",
	//	"心律不齐",
	//	"记忆力衰退",
	//}
	//
	//User.QueryByKeyWords(query)
	//fmt.Println()
	//User.QueryByKeyWords(queryCompare)
	//portForSendSlice := 9000
	//nodeToQuery := "http://192.168.42.129:8000"
	//startTime := "2023-10-11 14:32:57"
	//endTime := "2023-10-11 14:32:57"
	//
	//User.QueryData(nodeToQuery, startTime, endTime, portForSendSlice)
	//
	//Node.CreateConfig()
	Node.NodeInit()
}
