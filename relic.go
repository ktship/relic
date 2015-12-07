package relic

import (
	"sort"
	"fmt"
	"math/rand"
)
// 유물 시스템. (Exclusive)
// 각 가치의 비율과 비례하는 확률 리스트를 구축해서 계산한다.
//  A:1000, B:90, C:15 의 가치를 가진다고 할때, 
//  A-B : 11.1배, B-C : 6배.
// 	A 를 100로 환산하고 그 기준으로 배수를 곱해서 배열을 만든다.
//  A:100, B:A*11.1 + A값, C:B*6 + 값B -> 이것이 확률리스트임.
//  서버데이터로 각 아이템의 가치와 무슨 아이템으로 구성된 가챠인지에 대한 정보가 필요.

//  유저는 각 가챠에서 아이템들의 상태(이미 뽑힘, 뽑을 수 없는 아이템)에 대한 정보가 필요.
// data는 "nextTurn":0이면 첫번째 턴

// Relic : 유물 시스템. (Exclusive)
type Relic struct {
	io 						relicIO
	lazyTurnCount			int
	lazySelectedItemList	[]int
	lazyUID					int
	lazyRelicID				int
	isLoad					bool
}

// New : Relic 
func New(rio relicIO) *Relic {
	return &Relic{
		io		:rio,
		isLoad 	:false,
	}
}

type relicProb struct {
	iid 	int
	prob 	int
}
type byProb []relicProb
func (a byProb) Len() int           { return len(a) }
func (a byProb) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byProb) Less(i, j int) bool { return a[i].prob < a[j].prob }

const ATTR_TURN_COUNT = "tc"
const ATTR_SELECTED_LIST = "sl"

const PROBABILITY_START_VALUE = 100

// StatusReady 아직 안 뽑힘, StatusSelected 뽑힘
const (
	StatusReady = 0
	StatusSelected = 1
)

// GachaRelic : 유물 뽑기 
// 성공하면 유저 저장까지 완료됨.
func (r *Relic)GachaRelic(uid int, relicID int) (int, error) {
	// 유저 데이터 읽기
	r.loadRelicData(uid, relicID)
	// 확률 리스트 생성
	probList, err := r.getRelicProb(uid, relicID)
	if err != nil {
		return 0, err
	}
	lastProb := probList[len(probList)-1].prob
	randV := rand.Intn(lastProb)
	retItemID, err := r.getItemIDFromRandV(probList, randV)
	// 선택된 아이템 리스트에 추가
	r.lazySelectedItemList = append(r.lazySelectedItemList, retItemID)
	// 다 선택된 상태인지를 체크
	r.checkFull()
	err = r.saveRelicData(uid, relicID)
	if err != nil {
		return 0, err
	}
	return retItemID, nil
}

// getRelicProb : 유물 확률 리스트 
func (r *Relic)getRelicProb(uid int, relicID int) (byProb, error) {
	retItems, err := r.getValidItemList(uid, relicID)
	if err != nil {
		return nil, err
	}
	retProb, err := r.calcRelicProb(retItems)
	if err != nil {
		return nil, err
	}
	return retProb, nil
}

func (r *Relic)loadRelicData(uid int, relicID int) error{
	// user 정보
	dat, err := r.io.ReadUserAttr(uid, relicID)
	if err != nil {
		return err
	}
	// 몇번째 턴인가
	turnCount := dat[ATTR_TURN_COUNT].(int)
	// 뽑은 아이템
	strSelectedList := dat[ATTR_SELECTED_LIST].(string)
	selectedItemList, errJSON := strToListInt(strSelectedList)
	if errJSON != nil {
		return errJSON
	}
	
	r.lazyTurnCount = turnCount
	r.lazySelectedItemList = selectedItemList
	r.lazyUID = uid
	r.lazyRelicID = relicID
	r.isLoad = true;
	return nil
}

func (r *Relic)saveRelicData(uid int, relicID int) (error) {
	// 유저 데이터 로드가 끝난 상태에서만 불려져야 함.AddRelicr.
	if r.isLoad != false {
		return fmt.Errorf("Save operation MUST be After load user data")
	}
	
	// Save relic Data
	relicMap := make(map[string]interface{})
	relicMap[ATTR_TURN_COUNT] = r.lazyTurnCount
	var err error
	relicMap[ATTR_SELECTED_LIST], err = listIntToStr(r.lazySelectedItemList)
	if err != nil {
		return err
	}
	err = r.io.WriteUserAttr(uid, relicID, relicMap)
	if err != nil {
		return err
	}
	
	return nil
}
	
func (r *Relic)getValidItemList(uid int, relicID int) ([]int, error) {
	// 유저 데이터 로드가 끝난 상태에서만 불려져야 함.AddRelicr.
	if r.isLoad != false {
		return nil, fmt.Errorf("user data is NOT loaded")
	}
	
	// relic 정보 
	rInfo := getRelicInfo(relicID)
	var newBaseItems []int

	// 몇번째 턴인가
	turnCount := r.lazyTurnCount
	// 반복 가능 타입이면 
	if rInfo.isRepeatType {
		newBaseItems = rInfo.itemList
	} else {
		if turnCount > 0 {		// 여러 턴일때
			newBaseItems = rInfo.baseItemList
		} else {
		newBaseItems = rInfo.itemList
		}
	}
	
	// 이미 뽑은 아이템
	selectedItemList := r.lazySelectedItemList
	
	// 이미 뽑은 아이템 제외
	retItems := minusIntList(newBaseItems, selectedItemList)
	return retItems, nil
}

func (r *Relic)calcRelicProb(retItems []int) (byProb, error) {
	ret := make(byProb, len(retItems))
	for i, v := range retItems {
		val := getItemVal(v)
		if val == 0 {
			return nil, fmt.Errorf("calcRelicProb return nil")
		}
		ret[i] = relicProb {
			iid 	: v,
			prob	: val,
		}
		sort.Sort(ret)
	}
	// 각 value간에 비율
	for i, v := range ret {
		if i+1 >= len(ret) {
			break
		}
		ratio := v.prob / ret[i+1].prob
		if i == 0 {
			ret[0].prob = PROBABILITY_START_VALUE
		}
		ret[i+1].prob = ret[i].prob + int(ret[i].prob * ratio)
	}
	return ret, nil
}

func (r *Relic)getItemIDFromRandV(probList byProb, randV int) (int, error) {
	if randV < 0 {
		return -1, fmt.Errorf("randV must greater than 0")
	}
	if probList[len(probList)-1].prob <= randV {
		return -1, fmt.Errorf("randV must less than last value of probList")
	}
	for _, v := range probList {
		if v.prob > randV {
			return v.iid, nil
		}
	}
	return -1, fmt.Errorf("Invalid Args. randV:%d probList:%v", randV, probList)
}

func (r *Relic)checkFull() bool {
	// relic 정보 아이템 리스트 
	rInfo := getRelicInfo(r.lazyRelicID)
	// 유저가 뽑은 아이템 
	var iList []int
	if rInfo.isRepeatType {
		iList = r.lazySelectedItemList
	} else {
		if r.lazyTurnCount == 0  {
			iList = r.lazySelectedItemList
		} else {
			return false
		}
	}
	excluiveList := minusIntList(rInfo.itemList, iList)
	if len(excluiveList) == 0 {
		r.lazySelectedItemList = nil
		r.lazyTurnCount = r.lazyTurnCount + 1
		return true
	}
	return false
}
