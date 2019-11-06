package rules

import (
	"strconv"
)

func unsafeAtoi(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return val
}
