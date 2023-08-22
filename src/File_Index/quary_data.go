package File_Index

import (
	"IOT_Storage/src/Block_Chain"
	"github.com/emirpasic/gods/trees/avltree"
	"time"
)

func QueryData(tree *avltree.Tree, deviceId string, startTime time.Time, endTime time.Time) []Block_Chain.DATA {
	if tree.Empty() {
		return []Block_Chain.DATA{}
	}
	floor, found := tree.Floor(TreeKey{DeviceId: deviceId, TimeStamp: startTime, Serial: 1})
	if !found {
		floor = tree.Left()
	}
	if floor.Key.(TreeKey).DeviceId != deviceId || floor.Key.(TreeKey).TimeStamp.Before(startTime) {
		floor = floor.Next()
	}
	if floor == nil || floor.Key.(TreeKey).DeviceId != deviceId {
		return []Block_Chain.DATA{}
	}
	ceiling, found := tree.Ceiling(TreeKey{DeviceId: deviceId, TimeStamp: startTime, Serial: 1})
	if !found {
		ceiling = tree.Right()
	}
	if ceiling.Key.(TreeKey).DeviceId != deviceId || ceiling.Key.(TreeKey).TimeStamp.After(endTime) {
		ceiling = ceiling.Prev()
	}
	if ceiling == nil || ceiling.Key.(TreeKey).DeviceId != deviceId {
		return []Block_Chain.DATA{}
	}
	var data []Block_Chain.DATA
	for {
		data = append(data, floor.Value.(Block_Chain.DATA))
		if floor == ceiling {
			break
		}
		floor = floor.Next()
	}
	return data
}
