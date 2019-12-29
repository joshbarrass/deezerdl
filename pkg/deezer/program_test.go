package deezer

import (
	"fmt"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

type Configuration struct {
	ArlCookie string `envconfig:"ARL_COOKIE" required:"true"`
	DebugMode bool   `envconfig:"DEBUG_MODE"`
}

func testSetup(t *testing.T) (Configuration, *API) {
	var config Configuration
	if err := envconfig.Process("", &config); err != nil {
		t.Skip("Unable to load config -- skipping")
	}

	api, err := NewAPI(config.DebugMode)
	assert.Equal(t, nil, err)
	return config, api
}

// TestProgram runs any tests that work on the real api and hence
// require a working ARL cookie -- these will be skipped if the config
// cannot be loaded.
func TestProgram(t *testing.T) {
	t.Run("Cookie Login", func(t *testing.T) {
		config, api := testSetup(t)

		err := api.CookieLogin(config.ArlCookie)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, "", api.APIToken)
	})

	t.Run("Get Download Link", func(t *testing.T) {
		var (
			testID   = 3135553
			testName = "One More Time"
			testMD5  = "43808a3ac856cc117362ab94718603ba"
		)
		config, api := testSetup(t)

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

		u, err := track.GetDownloadURL(FLAC)
		assert.Equal(t, err, nil)
		fmt.Println(u.String())
	})
}
