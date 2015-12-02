package relic

type relicIO interface {
	ReadUserAttr(uid int, rid int) (map[string]interface{}, error)
	WriteUserAttr(uid int, rid int, data map[string]interface{}) error
	DelUserAttr(uid int, rid int) error
}
