package File_Index

import (
	"fmt"
	"testing"
	"time"
)

func TestBuildTraverser(t *testing.T) {
	filepath := string("backup.json")
	tree := BuildTraverser(filepath)
	fmt.Println(tree)
	if tree.Empty() == true {
		fmt.Errorf("Fail to create tree!\n")
	}
	//fmt.Println(tree)
	start, _ := time.Parse("2006-01-02 15:04:05", "2023-10-22 14:14:00")
	end, _ := time.Parse("2006-01-02 15:04:05", "2023-10-22 14:14:00")
	data := QueryData(tree, "6866974dc2b54eb8c5363f174f6f0c0e8d6a69ba3d9035957dceff69992402f2", start, end)
	fmt.Println(data)
}
