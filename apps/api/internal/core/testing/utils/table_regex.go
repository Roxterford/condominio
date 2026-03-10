package utils

import "fmt"

func GetMockTableRegex(table string) string {
	return fmt.Sprintf("[`\"]?%s[`\"]?", table)
}
