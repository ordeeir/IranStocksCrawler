package stringsh

import (
	"math"
	"strconv"
	"strings"
)

func TextPersianize(text string) string {

	text = strings.ReplaceAll(text, "ي", "ی")
	text = strings.ReplaceAll(text, "ك", "ک")
	return text
}

func ToInt(text string) int64 {
	if strings.Contains(text, ".") {
		x1, _ := strconv.ParseFloat(text, 64)
		x2 := int(math.Round(x1))
		return int64(x2)

	} else {
		x1, err := strconv.Atoi(text)
		if err != nil {
			return 0
		}
		return int64(x1)

	}

}

func ToString(num int64) string {
	str := strconv.Itoa(int(num))
	return str
}
