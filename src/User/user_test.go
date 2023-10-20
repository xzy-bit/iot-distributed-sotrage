package User

import "testing"

func TestQueryDataWithSM4(t *testing.T) {
	node := "http://192.168.42.129:8000"

	startTime := "2023-10-19 18:19:20"
	endTime := "2023-10-19 18:19:20"
	QueryDataWithSM4(node, startTime, endTime, 111, "123456")
}
