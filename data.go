package relic

// itemValueMap : 각 아이템의 가치의 정보를 가지는 맵. 키:아이템 id, 값:아이템의 가치를 보이는 정수혐
var itemValueMap = make(map[int]int)

type relicInfo struct {
	// 아이템 다 뽑으면 반복하는거임? 아니면 baseItem만 주는 거임?
	isRepeatType	bool
	// 아이템 리스트
	itemList 		[]int
	// 레어한 아이템 리스트임. 반복 되더라도 이건 다시 안 주는 아이템 리스트.
	rareItemList	[]int
	// 반복 안 할시에 기본으로 보상해줄 아이템 리스트
	baseItemList	[]int
}

// 각 가챠에 속하는 아이템 리스트를 가지는 맵. 키:가챠 id, 값:아이템 리스트
var relicMap = make(map[int]*relicInfo)

// SetItemVal : 
func SetItemVal(itemID int, itemValue int) {
	itemValueMap[itemID] = itemValue
}

// getItemVal : 
func getItemVal(itemID int) int {
	iv := itemValueMap[itemID]
	return iv
}

// AddRelic : 
func AddRelic(gid int, rInfo *relicInfo) {
	relicMap[gid] = rInfo
}

// getRelicInfo : relic id -> item list
func getRelicInfo(rid int) *relicInfo {
	return relicMap[rid]
}