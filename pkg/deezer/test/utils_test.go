package test

import (
	"testing"

	"github.com/joshbarrass/deezerdl/pkg/deezer"
	"github.com/stretchr/testify/assert"
)

var testTrack = deezer.Track{
	ID:           3135553,
	MD5:          "43808a3ac856cc117362ab94718603ba",
	MediaVersion: 7,
}

func TestMakeURLPath(t *testing.T) {
	var (
		testResult = "95e759c702e1ba9e43abf00d32049e257055f9973f05a51bebc99699b8ecb441d0a8070405e2f08e950e93aed73ab2c73a8d69e4e56700e7655f5281603f676d4b6c4a14ee471e77d44b485e7f27c7d9"
	)
	result, err := deezer.MakeURLPath(&testTrack, deezer.FLAC)
	assert.Equal(t, nil, err, "An error should not occur")
	assert.Equal(t, testResult, result, "The result should match the expected result")
}
