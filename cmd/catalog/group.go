package main

import (
	"errors"
)

func checkGroupsInput(strs []string) error {
	hasPublic := false
	hasOther := false

	if strs == nil || len(strs) < 1 {
		return errors.New("no groups given")
	}

	for _, str := range strs {
		if str == "public" {
			hasPublic = true
		} else {
			hasOther = true
		}
	}

	if hasPublic && hasOther {
		return errors.New("public and groups cannot be mixed")
	}

	return nil
}
