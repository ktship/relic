package relic
import (
		"encoding/json"
)

func minusIntList(a []int, b []int) []int {
	retList := make([]int, len(a))
	for _, va := range a {
		isIn := false
		for _, vb := range b {
			if va == vb {
				isIn = true
				break	
			}
		}
		if isIn == false {
			retList = append(retList, va)
		}
	}
	return retList
}

func strToListInt(str string) ([]int, error) {
	byteStr := []byte(str)
	var items []int
	if errJSON := json.Unmarshal(byteStr, &items); errJSON != nil {
		return nil, errJSON
	}
	return items, nil
}

func listIntToStr(listInt []int) (string, error) {
	byteArray, err := json.Marshal(listInt)
	if err != nil {
		return "", err
	}
	s := string(byteArray[:])
	return s, nil
}