package util

import "encoding/json"

func StructToMap(input interface{}) (result *map[string]interface{}, err error) {
	temp, _ := json.Marshal(input)
	json.Unmarshal(temp, &result)
	return
}
