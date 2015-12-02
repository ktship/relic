package relic
import (
	"testing"
	"strconv"
	"runtime"
	"fmt"
	"math/rand"
	"github.com/ktship/testio"
	gc "gopkg.in/check.v1"
)

const isExpertMode = false

func (s *TableSuite) SetUpSuite(c *gc.C) {}
func (s *TableSuite) SetUpTest(c *gc.C) {}
func (s *TableSuite) TearDownTest(c *gc.C) {}
func (s *TableSuite) TearDownSuite(c *gc.C) {}
func Test(t *testing.T) { gc.TestingT(t) }

type TableSuite struct {
	client *client
}
var _ = gc.Suite(&TableSuite {
	client : newClient(1, 0),
})
var _ = gc.Suite(&TableSuite {
	client : newClient(2, 0),
})

/*
가챠(유물) 시스템 테스트
	1. 각 아이템들의 확률 확인
	  1 확률 리스트 기본 : 아이템이 하나이상이라야 함
	  2 확률 리스트 기본 : 아이템이 하나일때는 확률은 그냥 1임.
	  모든 아이템에 대해,
	  3 확률 리스트 기본 : 아이템의 확률은 무조건 0보다 커야 함
	  4 확률 리스트 기본 : 뒤의 확률은 무조건 앞 확률보다 커야함
	2. 뽑힌 아이템이 정당한지 확인
	 1. itemID가 이미 뽑힌 상태여야 함.
	 2. prob 리스트에는 itemID 값이 들어 있어야함.
*/

type client struct {
	uid			int
	relics		map[int]map[int]int
	// 유물뽑기 작업을 연속으로 실행시에 평균 딜레이 시간 
	delayCall	int
}

func newClient(uid int, dc int) *client {
	return &client {
		uid 	: uid,
		relics 	: make(map[int]map[int]int),
		delayCall : dc,
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
}

// 50, 40, 7.5, 2.5
func (s *TableSuite) Test001_DynamoDBIO(c *gc.C) {
	tio := testio.New()
	relic := New(tio)

	rid := 0
	probList := relic.GetRelicProb(s.client.uid, rid)
	// 확률 리스트에 대한 검사
	s.checkProbList(c, rid, probList)
	iid, err := relic.GachaRelic(s.client.uid, rid)
	c.Assert(err, gc.IsNil)
	// 뽑힌 아이템의 정당성 검사
	s.checkGachaRelic(c, tio, s.client.uid, rid, iid, probList)
	probList = relic.GetRelicProb(s.client.uid, rid)
	// 뽑힌 아이템을 제외한 확률 리스트인지 확인
	s.checkPostGachaRelic(c, tio, s.client.uid, rid, iid, probList)
}

func (s *TableSuite)checkProbList(c *gc.C, rid int, probL []relicProb) {
	// 1 확률 리스트 기본 : 아이템이 하나이상이라야 함
	if len(probL) < 1 {
		c.Fatalf("아이템이 하나이상이라야 함 rid:%d", rid)
	}
	// 2 확률 리스트 기본 : 아이템이 하나일때는 확률은 그냥 1임.
	if len(probL) == 1 {
		if probL[0].prob != 1 {
			c.Fatalf("아이템이 하나라야 함 rid:%d", rid)
		}
	}
	// 모든 아이템에 대해,
	for i,v := range probL {
		// 3 확률 리스트 기본 : 아이템의 확률은 무조건 0보다 커야 함
		if v.prob <= 0 {
			c.Fatalf("아이템의 확률은 무조건 0보다 커야 함 rid:%d i:%d", rid, i)
		}
		// 3 확률 리스트 기본 : 뒤의 확률은 무조건 앞 확률보다 커야함
		if i+1 < len(probL) {
			if probL[i].prob >= probL[i+1].prob {
				c.Fatalf("앞의 확률은 무조건 뒤 확률보다 커야함. rid:%d, i:%d, v:%f, i+1:%d, v+1:%f", rid, i, probL[i].prob, i+1, probL[i+1].prob)
			}
		}
	}
}

func (s *TableSuite)checkGachaRelic(c *gc.C, io relicIO, uid int, relicID int, itemID int, probL []relicProb) {
	// 1. itemID가 뽑힌 상태여야 함.
	dat, err := io.ReadUserAttr(uid, relicID)
	c.Assert(err, gc.IsNil)
	strItemID := strconv.Itoa(itemID)
	status, ok := dat[strItemID]
	// itemID는 유저 데이터에 
	c.Assert(ok, gc.Equals, true)
	// itemID는 Exception 상태가 아니라 반드시 뽑힌 
	c.Assert(status, gc.Equals, StatusSelected)

	// 2. probL에는 itemID 값이 들어 있어야함.
	isIn := false
	for i,v := range probL {
		if v.iid == itemID {
			isIn = true
		}
	}
	c.Assert(isIn, gc.Equals, true)
}

func (s *TableSuite)checkPostGachaRelic(c *gc.C, io relicIO, uid int, relicID int, itemID int, probL []relicProb) {
	// 1. probL에는 itemID 값이 없어야함
	isIn := false
	for i,v := range probL {
		if v.iid == itemID {
			isIn = true
		}
	}
	c.Assert(isIn, gc.Equals, false)
}

