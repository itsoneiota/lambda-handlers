package helpers

import "encoding/json"

func IsValidJSONObject(s string) bool {
	var js interface{}
	err := json.Unmarshal([]byte(s), &js)
	if err != nil {
		return false
	}

	switch js.(type) {
	case map[string]interface{}, []interface{}:
		return true
	default:
		return false
	}
}
