package storage

import "strings"

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}

	if strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return true
	}

	return false
}
