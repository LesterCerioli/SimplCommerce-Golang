package controllers

import (
	"strconv"
)

func parseUint(s string) uint {
	id, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0
	}
	return uint(id)
}

func parseInt(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func hasRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
