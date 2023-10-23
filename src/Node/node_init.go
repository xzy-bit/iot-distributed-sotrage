package Node

import (
	"IOT_Storage/src/Database"
	"encoding/json"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
)

var tree *avltree.Tree
var nodeConfig *Config
var table []SearchIndex

type Config struct {
	NodeId                  int
	AddressBook             []string
	PortForPIng             int
	PortForToken            int
	PortForBlock            int
	PortForIndexBroad       int
	PortForGetSlice         int
	PortForQuery            int
	PortForSendSlice        int
	PortForSendIndex        int
	PortForGetIndexFromUser int
	PortForQueryByKeyWords  int
	PortIndexPageForUser    int
}

func CreateConfig() {
	config := new(Config)
	config.NodeId = 0
	config.PortForPIng = 8080
	config.PortForToken = 7080
	config.PortForBlock = 9080
	config.PortForIndexBroad = 8060
	config.PortForGetSlice = 10080
	config.PortForQuery = 8000
	config.PortForSendSlice = 9000
	config.PortForSendIndex = 9060
	config.PortForGetIndexFromUser = 8040
	config.PortForQueryByKeyWords = 8020
	config.PortIndexPageForUser = 10000

	address := []string{
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
	}
	config.AddressBook = address
	data, _ := json.Marshal(config)
	os.WriteFile("config.json", data, 0666)
}

func ReadConfig() *Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
		log.Fatal("Read config error!")
	}
	config := new(Config)
	err = json.Unmarshal(data, config)
	return config
}

func NodeInit() {
	nodeConfig = ReadConfig()
	//ReadIndexes()
	//pingRouter := Ping()
	//go pingRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForPIng))
	//
	//tree = File_Index.BuildTraverser("backup.json")
	////
	//pipe := make(chan string)
	//var urlBooks []string
	////tree := File_Index.BuildTraverser("backup.json")
	//for index, nodeAddress := range nodeConfig.AddressBook {
	//	nodeId := index
	//	if nodeId == nodeConfig.NodeId {
	//		continue
	//	}
	//	go func(nodeAddress string) {
	//		trueUrl := nodeAddress + ":" + strconv.Itoa(nodeConfig.PortForPIng+nodeId)
	//		req := Controller.CreatePingReq(trueUrl)
	//		for {
	//			resp := Controller.SendRequest(req)
	//			if resp == nil {
	//				//log.Printf("Can not get connection with %s\n", trueUrl)
	//				time.Sleep(time.Second)
	//				pipe <- ""
	//				continue
	//			}
	//			//log.Printf("%s is alive\n", trueUrl)
	//			pipe <- trueUrl
	//			break
	//		}
	//	}(nodeAddress)
	//}
	//go func() {
	//	for {
	//		select {
	//		case trueUrl := <-pipe:
	//			if trueUrl == "" {
	//				continue
	//			} else {
	//				urlBooks = append(urlBooks, trueUrl)
	//				//log.Println(len(urlBooks))
	//			}
	//		}
	//		if len(urlBooks) == 6 {
	//			if nodeConfig.NodeId == 0 {
	//
	//				if Block_Chain.GetPrevBlock() == nil {
	//					genius := Block_Chain.GeniusBlock()
	//					log.Println("genius block:")
	//					log.Println(genius)
	//
	//					Block_Chain.StoreBlock(*genius)
	//
	//					File_Index.InsertBlock(genius, tree)
	//					log.Println("tree:")
	//					log.Println(tree)
	//
	//					BroadCastBlock(*genius)
	//				}
	//
	//				index := (nodeConfig.NodeId + 1) % 7
	//				trueUrl := nodeConfig.AddressBook[index] + ":" + strconv.Itoa(nodeConfig.PortForToken+index)
	//				req, _ := http.NewRequest("GET", trueUrl+"/token", nil)
	//				Controller.SendRequest(req)
	//			}
	//			break
	//		}
	//	}
	//}()
	//
	//log.Println(tree)

	blockRouter := NodeGetBlock()
	go blockRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForBlock))

	indexBroadRouter := NodeForIndexBroad()
	go indexBroadRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForIndexBroad))

	tokenRouter := NodeGetToken()
	go tokenRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForToken))

	router_for_send_slice := gin.Default()
	v5 := router_for_send_slice.Group("")
	SendSlice(v5)
	SendSliceSm4(v5)
	go router_for_send_slice.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForSendSlice))

	getSliceRouter := NodeGetSlice()
	go getSliceRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForGetSlice))

	router_for_node_query := gin.Default()
	v4 := router_for_node_query.Group("")
	NodeSendIndex(v4)
	NodeSendSplitMat(v4)
	go router_for_node_query.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForSendIndex))

	router_for_keywords_query := gin.Default()
	v3 := router_for_keywords_query.Group("")
	NodeForKeyWordsQuery(v3)
	NodeForKeyWordsQueryWithSplitMat(v3)
	go router_for_keywords_query.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForQueryByKeyWords))

	router_for_indexes := gin.Default()
	v2 := router_for_indexes.Group("")
	NodeGetIndex(v2)
	NodeGetIndexWithSplitMat(v2)
	NodeGetSpltMat(v2)
	go router_for_indexes.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForGetIndexFromUser))

	db := Database.ConnectDB()
	router_for_front_page := gin.Default()
	router_for_front_page.LoadHTMLGlob("templates/*")
	v1 := router_for_front_page.Group("")
	NodeGetQuery(v1)
	NodeIndexPageForUser(v1)
	NodeUploadPageForUser(v1)
	NodeQueryDataForUSer(v1)
	NodeLoginPage(v1, db)
	NodeSearchPageForUser(v1)
	NodeSearchServerForUser(v1)
	NodeDownload(v1)
	router_for_front_page.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForQuery))

	//queryIndex := NodeGetQuery()
	//go queryIndex.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForQuery))
	//
	//indexPage := NodeIndexPageForUser()
	//indexPage.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortIndexPageForUser))

}
