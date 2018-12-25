package common

import "strconv"

func String2int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		Log("String2int Error", err)
	}
	return i
}
