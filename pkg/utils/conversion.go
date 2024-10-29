package utils

import (
	"strconv"
)

func StringToUint(s string) (uint, error) {
	u64, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(u64), nil
}
