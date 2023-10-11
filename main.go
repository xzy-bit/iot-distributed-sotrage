package main

import "IOT_Storage/src/User"

func main() {
	//iot device process
	//file, _ := os.Open("private.pem")
	//if file == nil {
	//	Identity_Verify.GenerateKey(false)
	//}
	//IOT_Device.IotInit()
	//nodes := []string{
	//	"http://192.168.42.129:10080",
	//	"http://192.168.42.129:10081",
	//	"http://192.168.42.129:10082",
	//	"http://192.168.42.129:10083",
	//	"http://192.168.42.129:10084",
	//	"http://192.168.42.129:10085",
	//	"http://192.168.42.129:10086",
	//}
	//IOT_Device.SendSliceToNode(nodes)

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
	
	portForSendSlice := 9000
	nodeToQuery := "http://192.168.42.129:10000"
	startTime := "2023-10-07 15:25:10"
	endTime := "2023-10-07 15:25:10"
	User.QueryData(nodeToQuery, startTime, endTime, portForSendSlice)
	}
