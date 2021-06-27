package main

import (
	"strings"
)

type structure struct {
	action     string
	dirChanges []string
}

func NewStructure(action string) *structure {
	s := structure{action: action}
	s.dirChanges = CalculateDirs(&s)
	return &s
}

func CalculateDirs(s *structure) []string {
	splitAction := strings.Split(s.action, "/")
	if strings.HasSuffix(s.action, "/") {
		return splitAction[:len(splitAction)-1]
	}
	return splitAction
}
