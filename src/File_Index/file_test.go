package File_Index

import (
	"fmt"
	"testing"
	"time"
)

func TestBuildTraverser(t *testing.T) {
	filepath := string("backup.json")
	tree := BuildTraverser(filepath)

	if tree.Empty() == true {
		fmt.Errorf("Fail to create tree!\n")
	}
	fmt.Println(tree)
	start, _ := time.Parse("2006-01-02 15:04:05", "2023-09-09 13:51:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2023-09-09 13:51:20")
	data := QueryData(tree, "804d0d9378026bcda165ac37e634b5812ba192f6", start, end)
	fmt.Println(data)
}
