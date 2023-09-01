package IOT_Device

import (
	"IOT_Storage/src/Controller"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"io"
	"log"
	"os"
)

func GenerateIotId(file *os.File) (id string) {
	h := sha3.NewLegacyKeccak256()
	_, err := io.Copy(h, file)
	if err == nil {
		fmt.Println("Generating userId from public key...")
		hash := h.Sum(nil)
		return hex.EncodeToString(hash[12:])
	} else {
		log.Fatal("Failed to generate iot id !")
		return "fatal error!"
	}
}

func IotInit() {

	publicFile, _ := os.Open("public.pem")
	userId := GenerateIotId(publicFile)
	publicFile.Close()
	fmt.Println(userId)

	privateFile, _ := os.Open("private.pem")
	sendPrivateToServer := Controller.CreateSendFileReq(privateFile, userId+".pem", "http://192.168.42.129:8090/receive")
	privateFile.Close()

	publicFile, _ = os.Open("public.pem")
	sendPublicToUser := Controller.CreateSendFileReq(publicFile, "public.pem", "http://localhost:8090/receive")

	resp := Controller.SendRequest(sendPublicToUser)
	if resp.StatusCode != 200 {
		log.Fatal(resp.StatusCode)
	}

	privateFile, _ = os.Open("private.pem")
	sendPrivateToUser := Controller.CreateSendFileReq(privateFile, "private.pem", "http://localhost:8090/receive")

	resp = Controller.SendRequest(sendPrivateToServer)
	if resp.StatusCode != 200 {
		log.Fatal(resp.StatusCode)
	}

	resp = Controller.SendRequest(sendPrivateToUser)
	if resp.StatusCode != 200 {
		log.Fatal(resp.StatusCode)
	}

}
