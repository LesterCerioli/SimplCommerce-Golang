package controllers

import (
	"strconv"
)

func parseInt(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func hasRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}
