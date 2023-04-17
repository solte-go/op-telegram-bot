package postgresql

import "regexp"

var reDuplicateKey = regexp.MustCompile(`duplicate key value violates unique constraint`)

func duplicateKeyError(err error) bool {
	return reDuplicateKey.MatchString(err.Error())
}
