package IOT_Device

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestSendSliceWithSM4(t *testing.T) {
	file, _ := os.OpenFile("test.jpg", os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Read the file into a byte slice
	buffer := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	nodes := []string{
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
		"http://192.168.42.129",
	}
	portForSlice := 10080
	password := "123456"
	SendSliceWithSM4(buffer, nodes, password, portForSlice)
}

//func TestSendSM4Slice(t *testing.T) {
//	file, _ := os.OpenFile("test.jpg", os.O_RDWR|os.O_CREATE, 0755)
//	defer file.Close()
//	// Get the file size
//	stat, err := file.Stat()
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	// Read the file into a byte slice
//	buffer := make([]byte, stat.Size())
//	_, err = bufio.NewReader(file).Read(buffer)
//	if err != nil && err != io.EOF {
//		fmt.Println(err)
//		return
//	}
//
//	nodes := []string{
//		"http://192.168.42.129",
//		"http://192.168.42.129",
//		"http://192.168.42.129",
//		"http://192.168.42.129",
//		"http://192.168.42.129",
//		"http://192.168.42.129",
//		"http://192.168.42.129",
//	}
//	portForSlice := 10080
//	password := "123456"
//	SendSM4Slice(buffer, nodes, password, portForSlice)
//}
