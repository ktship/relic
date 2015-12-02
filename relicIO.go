package relic

//  유저는 각 가챠에서 아이템들의 상태(이미 뽑힘, 뽑을 수 없는 아이템)에 대한 정보가 필요.
// data는 "nextTurn":0 1이면 두번이상임. "itemID":status
type relicIO interface {
	ReadUserAttr(uid int, rid int) (map[string]interface{}, error)
	WriteUserAttr(uid int, rid int, data map[string]interface{}) error
	DelUserAttr(uid int, rid int) error
}
