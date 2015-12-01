package relic
import (
	"testing"
	"strconv"
	"runtime"
	"fmt"
	"math/rand"
	"github.com/ktship/testio"
)

const isExpertMode = false

var client0 *client

/*
가챠(유물) 시스템 테스트
	데이터
	 유물 0번, 반복 OK, 아이템 : 0, 2, 6, 7, 레어 아이템 : 0, 2
	 유물 1번, 반복 NO, 아이템 : 0, 3, 6, 기본 아이템 : 7, 8, 9
기본 테스트
	1. 클라이언트 0, 1로 0.1 정도의 랜덤 간격으로 계속 뽑음. (기본 모드)
	   유물 0은 처음은 모두 뽑히고, 이후로는 2, 3만 뽑히게 되는지 확인. 
	   각 뽑기 할때마다 각 아이템들의 확률이 맞는지를 확인. 3번씩 반복.
	2. 5초 랜덤 간격으로 뽑음. (클라이언트0과 동시에 함.고급 모드)
	   각 뽑기 할때마다 각 아이템들의 확률이 맞는지를 확인. 1분간 확인.
 */

type client struct {
	uid			int
	relics		map[int]map[int]int
}

func newClient(uid int) *client {
	return &client {
		uid 	: uid,
		relics 	: make(map[int]map[int]int),
	}
}

func (c *client)setItemStatus(relicID int, itemID int, status int) {
	if _, ok := c.relics[relicID]; !ok {
		c.relics[relicID] = make(map[int]int)
	}
	c.relics[relicID][itemID] = status
}

func init() {
	fmt.Printf("Running On %s, %s, %s, %d-bit \n", runtime.Compiler, runtime.GOARCH, runtime.GOOS, strconv.IntSize)
	// 랜덤 함수의 값을 고정시킴.
	rand.Seed(0)

	// item 가치 테이블
	// 0:1000, 1:900, 2:800, 3:500, 4:300, 5:200, 6:100, 7:50, 8:20, 9:10
	SetItemVal(0, 1000)
	SetItemVal(1, 900)
	SetItemVal(2, 800)
	SetItemVal(3, 500)
	SetItemVal(4, 300)
	SetItemVal(5, 200)
	SetItemVal(6, 150)
	SetItemVal(7, 50)
	SetItemVal(8, 20)
	SetItemVal(9, 10)

	// 유물 구성
	// 유물 rid 0번, 반복 OK
	// 아이템 : 0, 1, 2, 3
	// 레어 아이템 : 0, 1
	relic0 := relicInfo {
		isRepeatType 	: true,
		itemList 		: []int{0, 1, 2, 3},
		rareItemList 	: []int{0, 1},
		baseItemList	: []int{},
	}
	AddRelic(0, &relic0)
	// 유물 rid 1번, 반복 NO
	// 아이템 : 0, 3, 6
	// 기본 아이템 : 7, 8, 9
	relic1 := relicInfo {
		isRepeatType 	: false,
		itemList 		: []int{0, 3, 6},
		rareItemList 	: []int{},
		baseItemList	: []int{7, 8, 9},
	}
	AddRelic(0, &relic1)
	
	client0 = newClient(111)
}

// 50, 40, 7.5, 2.5
func Test000_relic0(t *testing.T) { 
	tio := testio.New()
	relic := New(tio)

	rid := 0
	probList := relic.GetRelicProb(client0.uid, rid)
	checkProbList(t, rid, probList)
	iid, err := relic.GachaRelic(client0.uid, rid)
	if err != nil {
		t.Fatalf("err occured:%s", err)
	}
	checkGachaRelic(t, rid, iid, probList)
}

func checkProbList(t *testing.T, rid int, probL []relicProb) {
	// 1 확률 리스트 기본 : 아이템이 하나이상이라야 함
	if len(probL) < 1 {
		t.Fatalf("아이템이 하나이상이라야 함 rid:%d", rid)
	}
	// 2 확률 리스트 기본 : 아이템이 하나일때는 확률은 그냥 1임.
	if len(probL) == 1 {
		if probL[0].prob != 1 {
			t.Fatalf("아이템이 하나라야 함 rid:%d", rid)
		}
	}
	// 모든 아이템에 대해,
	for i,v := range probL {
		// 3 확률 리스트 기본 : 아이템의 확률은 무조건 0보다 커야 함
		if v <= 0 {
			t.Fatalf("아이템의 확률은 무조건 0보다 커야 함 rid:%d i:%d", rid, i)
		}
		// 3 확률 리스트 기본 : 뒤의 확률은 무조건 앞 확률보다 커야함
		if i+1 < len(probL) {
			if probL[i].prob >= probL[i+1].prob {
				t.Fatalf("앞의 확률은 무조건 뒤 확률보다 커야함. rid:%d, i:%d, v:%f, i+1:%d, v+1:%f", rid, i, probL[i].prob, i+1, probL[i+1].prob)
			}
		}
	}
}

func checkGachaRelic(t *testing.T, relicID int, itemID int, probL []relicProb) {
	// 1. itemID가 이미 뽑힌 상태여야 함.

	// 2. itemID가 probL에는 없는 상태여아 함.( 뽑히기 전 타이밍이므로)

	// 3. itemID는 Exception 상태가 아니어야 함.

	// 4. relicProb List 에는 itemID 값이 들어 있어야함.
}



