package controllers

import "strconv"

func parseUintWithDefault(str string, dflt uint64) uint64 {
	if v, err := strconv.ParseUint(str, 10, 64); err == nil {
		return v
	} else {
		return dflt
	}
}
