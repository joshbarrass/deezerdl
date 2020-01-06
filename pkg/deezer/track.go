package deezer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Format int

const (
	FLAC    Format = 9
	MP3_320        = 3
	MP3_256        = 5
)

const (
	downloadHostFormat = "e-cdns-proxy-%c.dzcdn.net"
	downloadPathFormat = "/mobile/1/%s"
)

type Track struct {
	ID           int     `json:"SNG_ID,string"`
	Title        string  `json:"SNG_TITLE"`
	TrackNumber  int     `json:"TRACK_NUMBER,string"`
	Gain         float32 `json:"GAIN,string"`
	MD5          string  `json:"MD5_ORIGIN"`
	MediaVersion int     `json:"MEDIA_VERSION,string"`
	api          *API
}

var NoMD5Error = errors.New("no MD5 hash -- try authenticating")

// GetDownloadURL gets the download url (as a *url.URL) for a given
// format
func (track *Track) GetDownloadURL(format Format) (*url.URL, error) {
	if len(track.MD5) == 0 {
		if err := track.GetMD5(); err != nil {
			return nil, NoMD5Error
		}
	}
	path, err := MakeURLPath(track, format)
	if err != nil {
		return nil, err
	}
	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf(downloadHostFormat, track.MD5[0]),
		Path:   fmt.Sprintf(downloadPathFormat, path),
	}
	return &u, nil
}

// GetMD5 uses an alternative API to get the MD5 of the track
func (track *Track) GetMD5() error {
	resp, err := track.api.MobileApiRequest(getSongMobileMethod,
		strings.NewReader(fmt.Sprintf(`{"SNG_ID":%d}`, track.ID)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if track.api.DebugMode {
		DumpResponse(resp, "GetMD5")
	}

	// decode result key into a struct from the body
	var data struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { // uses the body directly
		return err
	}
	// and then get the MD5
	var results struct {
		MD5 string `json:"MD5_ORIGIN"`
	}
	if err := json.Unmarshal(data.Results, &results); err != nil {
		return err
	}

	if results.MD5 == "" {
		return NoMD5Error
	}
	track.MD5 = results.MD5

	return nil

}

// GetSongData gets a track
func (api *API) GetSongData(ID int) (*Track, error) {
	// make the request
	body := strings.NewReader(fmt.Sprintf(`{"SNG_ID":%d}`, ID))
	resp, err := api.ApiRequest(getSongMethod, body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if api.DebugMode {
		DumpResponse(resp, "GetSongData")
	}

	// decode results key
	var data struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// decode track from results
	var track Track
	if err := json.Unmarshal(data.Results, &track); err != nil {
		return nil, err
	}
	track.api = api

	return &track, nil
}
