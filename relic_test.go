package relic
import (
	"testing"
	"strconv"
	"runtime"
	"fmt"
	"math/rand"
)

const isExpertMode = false

/*
가챠(유물) 시스템 테스트
	데이터
	 유물 0번, 반복 OK, 아이템 : 0, 1, 2, 3, 레어 아이템 : 0, 1
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
	relic0		map[itemID]itemStatus
	relic1		map[itemID]itemStatus
}

func newClient(uid int) *client {
	return &client {
		uid 	: uid,
		relic0 	: make(map[itemID]itemStatus),
		relic1 	: make(map[itemID]itemStatus),
	}
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
	SetItemVal(6, 100)
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
}

func Test000_relic0(t *testing.T) { 
	client0 := newClient(123)
	
}
