package Node

import (
	"IOT_Storage/src/Block_Chain"
	"IOT_Storage/src/Controller"
	"IOT_Storage/src/File_Index"
	"encoding/json"
	"github.com/emirpasic/gods/trees/avltree"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var tree *avltree.Tree
var nodeConfig *Config
var table []SearchIndex

type Config struct {
	NodeId                 int
	AddressBook            []string
	PortForPIng            int
	PortForToken           int
	PortForBlock           int
	PortForIndex           int
	PortForGetSlice        int
	PortForQuery           int
	PortForSendSlice       int
	PortForSendIndex       int
	PortForGetIndex        int
	PortForQueryByKeyWords int
}

func CreateConfig() {
	config := new(Config)
	config.NodeId = 0
	config.PortForPIng = 8080
	config.PortForToken = 7080
	config.PortForBlock = 9080
	config.PortForIndex = 9040
	config.PortForGetSlice = 10080
	config.PortForQuery = 8000
	config.PortForSendSlice = 9000
	config.PortForSendIndex = 9060
	config.PortForGetIndex = 9020
	config.PortForQueryByKeyWords = 9080

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
	ReadIndexes()
	pingRouter := Ping()
	go pingRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForPIng))

	tree = File_Index.BuildTraverser("backup.json")
	//
	pipe := make(chan string)
	var urlBooks []string
	//tree := File_Index.BuildTraverser("backup.json")
	for index, nodeAddress := range nodeConfig.AddressBook {
		nodeId := index
		if nodeId == nodeConfig.NodeId {
			continue
		}
		go func(nodeAddress string) {
			trueUrl := nodeAddress + ":" + strconv.Itoa(nodeConfig.PortForPIng+nodeId)
			req := Controller.CreatePingReq(trueUrl)
			for {
				resp := Controller.SendRequest(req)
				if resp == nil {
					log.Printf("Can not get connection with %s\n", trueUrl)
					time.Sleep(time.Second)
					pipe <- ""
					continue
				}
				log.Printf("%s is alive\n", trueUrl)
				pipe <- trueUrl
				break
			}
		}(nodeAddress)
	}
	go func() {
		for {
			select {
			case trueUrl := <-pipe:
				if trueUrl == "" {
					continue
				} else {
					urlBooks = append(urlBooks, trueUrl)
					//log.Println(len(urlBooks))
				}
			}
			if len(urlBooks) == 6 {
				if nodeConfig.NodeId == 0 {

					if Block_Chain.GetPrevBlock() == nil {
						genius := Block_Chain.GeniusBlock()
						log.Println("genius block:")
						log.Println(genius)

						Block_Chain.StoreBlock(*genius)

						File_Index.InsertBlock(genius, tree)
						log.Println("tree:")
						log.Println(tree)

						BroadCastBlock(*genius)
					}

					index := (nodeConfig.NodeId + 1) % 7
					trueUrl := nodeConfig.AddressBook[index] + ":" + strconv.Itoa(nodeConfig.PortForToken+index)
					req, _ := http.NewRequest("GET", trueUrl+"/token", nil)
					Controller.SendRequest(req)
				}
				break
			}
		}
	}()

	//log.Println(tree)

	blockRouter := NodeGetBlock()
	go blockRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForBlock))

	indexRouter := NodeForIndexBroad()
	go indexRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForIndex))

	tokenRouter := NodeGetToken()
	go tokenRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForToken))

	getSliceRouter := NodeGetSlice()
	go getSliceRouter.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForGetSlice))

	queryIndex := NodeGetQuery()
	go queryIndex.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForQuery))

	sendSlice := NodeSendSlice()
	go sendSlice.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForSendSlice))

	getIndex := NodeGetIndex()
	getIndex.Run(":" + strconv.Itoa(nodeConfig.NodeId+nodeConfig.PortForGetIndex))
}
