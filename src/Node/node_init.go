package Node

import (
	"IOT_Storage/src/Controller"
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
	NodeId      int
	AddressBook []string
	Port        string
}

func CreateConfig() {
	config := new(Config)
	config.NodeId = 0
	config.Port = ":8082"
	address := []string{
		"http://192.168.42.129:8082",
		"http://192.168.42.129:8083",
		"http://192.168.42.129:8084",
		"http://192.168.42.129:8085",
		"http://192.168.42.129:8086",
		"http://192.168.42.129:8087",
		"http://192.168.42.129:8088",
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
	config := ReadConfig()

	//pingRouter := Ping()
	//go pingRouter.Run(config.Port)

	sig := make(chan bool)
	count := 1
	//tree := File_Index.BuildTraverser("backup.json")
	for nodeId, nodeAddress := range config.AddressBook {
		if nodeId == config.NodeId {
			continue
		}
		go func(nodeAddress string) {
			req := Controller.CreatePingReq(nodeAddress)
			for {
				resp := Controller.SendRequest(req)
				if resp.StatusCode != 200 {
					log.Printf("Can not get connection with %s\n", nodeAddress)
					time.Sleep(time.Second)
					sig <- false
					continue
				}
				sig <- true
				break
			}
		}(nodeAddress)
	}
	go func() {
		for {
			select {
			case <-sig:
				count++
				log.Printf("%d nodes connected\n", count)
			}
			if count == 7 {
				break
			}
		}
	}()
	log.Println(count)
}
