package relic

type relicIO interface {
	ReadUserRelic(uid int, rid int) (map[string]interface{}, error)
	WriteUserRelic(uid int, rid int, data map[string]interface{}) error
	DelUserRelic(uid int, rid int) error
}
