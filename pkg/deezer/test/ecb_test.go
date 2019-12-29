package test

import (
	"testing"

	"github.com/joshbarrass/deezloader/pkg/deezer"
	"github.com/stretchr/testify/assert"
)

func TestECB(t *testing.T) {
	var (
		testKey    = []byte("aaaaaaaaaaaaaaaa")
		testCT     = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		testResult = []byte("\xa2\xa1h\x84\xdf;.7\xd8c\xddh\xcc\xf0\x8e\x1c\xa2\xa1h\x84\xdf;.7\xd8c\xddh\xcc\xf0\x8e\x1c")
	)

	result, err := deezer.ECB(testKey, testCT)
	assert.Equal(t, nil, err, "An error should not occur")
	assert.Equal(t, testResult, result, "The result should match the expected result from a different language")
}
