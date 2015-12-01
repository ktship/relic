package relic

// Relic : 유물 시스템. (Exclusive)
//  A:1000, B:90, C:15 의 가치를 가진다고 가정.
//  가치의 역순으로 나열해서 퍼센트를 계산. A:15%, B:90%, C:1000% 의 확률을 가짐. (X)
//  수정! : 역순은 그냥 운 좋을때만 맞을듯.
//  1. 아이템 순서대로 각각의 가치의 배율을 뽑음.
//  A-B : 11.1배, B-C : 6배.
// 	2. A 를 1로 환산하고 그 기준으로 배수를 곱해서 배열을 만든다.
//  A:1, B:11.1, C:66.6 -> 이것이 확률리스트임.
//  뽑힌 아이템은 제외하고 다시 퍼센트를 계산하는 방식임.
//  서버데이터로 각 아이템의 가치와 무슨 아이템으로 구성된 가챠인지에 대한 정보가 필요.
//  유저는 각 가챠에서 아이템들의 상태(이미 뽑힘, 뽑을 수 없는 아이템)에 대한 정보가 필요.
type Relic struct {
	io 		relicIO
}

// New : Relic 
func New(rio relicIO) *Relic {
	return &Relic{
		io	:rio,
	}
}

type relicProb struct {
	iid 	int
	prob 	float64
}


// StatusReady 아직 안 뽑힘, StatusSelected 뽑힘, StatusException 레어라서 한번뽑으면 예외
const (
	StatusReady = 0
	StatusSelected = 1
	StatusException = 2
)

// GachaRelic : 유물 뽑기 
func (r *Relic)GachaRelic(uid int, rid int) (int, error) {
	return 0, nil
}

// GachaRelic : 유물 뽑기
func (r *Relic)GetRelicProb(uid int, rid int) []relicProb {
	return 0, nil
}
