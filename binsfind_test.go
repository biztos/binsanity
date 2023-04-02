// binsfind_test.go - tests for stuff in binsfind.go
package binsanity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biztos/binsanity"
)

// func ValidIdent(s string) bool {
// 	if len(s) == 0 || !unicode.IsLetter(rune(s[0])) {
// 		return false
// 	}
// 	for _, c := range s {
// 		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '_' {
// 			return false
// 		}
// 	}
// 	return true
// }

func TestValidIdentZeroLengthFalse(t *testing.T) {

	assert := assert.New(t)

	assert.False(binsanity.ValidIdent(""))

}

func TestValidIdentNonLetterStartFalse(t *testing.T) {

	assert := assert.New(t)

	assert.False(binsanity.ValidIdent("/this"))

}

func TestValidIdentNonLetterFinishFalse(t *testing.T) {

	assert := assert.New(t)

	assert.False(binsanity.ValidIdent("this/"))

}

func TestValidIdentWeirdButAllowedTrue(t *testing.T) {

	assert := assert.New(t)

	// yay effing ident spec...
	assert.True(binsanity.ValidIdent("ok12ßßมาก"))

}
