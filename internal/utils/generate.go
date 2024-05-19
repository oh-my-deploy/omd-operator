package utils

import (
	"github.com/google/uuid"
)

func RandomStringLists(idx int) []string {
	r := make([]string, 0)
	for i := 0; i < idx; i++ {
		r = append(r, uuid.New().String())
	}
	return r
}
