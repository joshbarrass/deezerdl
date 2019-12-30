package deezer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBlowfishKey(t *testing.T) {
	var (
		testResult = "e219??zng,s67i72"
	)
	result := testTrack.GetBlowfishKey()
	assert.Equal(t, testResult, string(result))
}
