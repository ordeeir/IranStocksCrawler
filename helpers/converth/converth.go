package converth

import (
	"strconv"
	"strings"
)

func HumanReadableToBytes(text string) int64 {

	text = strings.ToUpper(text)
	var res int64 = 0
	var m int64 = 1
	if strings.Contains(text, "B") {
		text = strings.Replace(text, "B", "", 1)
	}
	if strings.Contains(text, "K") {
		text = strings.Replace(text, "K", "", 1)
		m = 1024
	}
	if strings.Contains(text, "M") {
		text = strings.Replace(text, "M", "", 1)
		m = 1024 * 1024
	}
	if strings.Contains(text, "G") {
		text = strings.Replace(text, "G", "", 1)
		m = 1024 * 1024 * 1024
	}
	if strings.Contains(text, "T") {
		text = strings.Replace(text, "T", "", 1)
		m = 1024 * 1024 * 1024 * 1024
	}
	if result, err := strconv.ParseInt(text, 10, 64); err == nil {
		res = m * result
	}

	return res
}
