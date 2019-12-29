package test

import (
	"fmt"
	"testing"

	"github.com/joshbarrass/deezerdl/pkg/deezer"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

type Configuration struct {
	ArlCookie string `envconfig:"ARL_COOKIE" required:"true"`
	DebugMode bool   `envconfig:"DEBUG_MODE"`
}

func setup(t *testing.T) (Configuration, *deezer.API) {
	var config Configuration
	err := envconfig.Process("", &config)
	assert.Equal(t, nil, err)

	api, err := deezer.NewAPI(config.DebugMode)
	assert.Equal(t, nil, err)
	return config, api
}

func TestGetToken(t *testing.T) {
	config, api := setup(t)

	err := api.CookieLogin(config.ArlCookie)
	assert.Equal(t, nil, err)
	assert.NotEqual(t, "", api.APIToken)
}

func TestGetTrack(t *testing.T) {
	var (
		testID   = 3135553
		testName = "One More Time"
		testMD5  = "43808a3ac856cc117362ab94718603ba"
	)
	config, api := setup(t)

	api.CookieLogin(config.ArlCookie)

	track, err := api.GetSongData(testID)
	assert.Equal(t, err, nil)

	assert.Equal(t, testID, track.ID)
	assert.Equal(t, testName, track.Title)

	if track.MD5 == "" {
		err := track.GetMD5()
		assert.Equal(t, nil, err)
	}

	assert.Equal(t, testMD5, track.MD5)

	u, err := track.GetDownloadURL(deezer.FLAC)
	assert.Equal(t, err, nil)
	fmt.Println(u.String())
}
