package utils

import (
	"encoding/json"
	"fmt"
)

func ToJSON(v any) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %s", err)
	}
	return string(data)
}
