package deezer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const AlbumAPIFormat = "https://api.deezer.com/album/%d"

type AlbumTrack struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

// AlbumResponse is an intermediate format for getting album data that
// stores the data before putting it in an Album struct
type AlbumResponse struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	CoverURL    string `json:"cover"`
	CoverSmall  string `json:"cover_small"`
	CoverMedium string `json:"cover_medium"`
	CoverBig    string `json:"cover_big"`
	CoverXL     string `json:"cover_xl"`
	Date        string `json:"release_date"`
	Tracks      struct {
		Data []AlbumTrack `json:"data"`
	} `json:tracks"`
}

// Album stores the data for the album of interest
type Album struct {
	ID       int
	Title    string
	Link     string
	CoverURL string
	Covers   Covers
	Date     time.Time
	Tracks   []AlbumTrack
}

// Covers stores the different cover sizes available
type Covers struct {
	Small  string
	Medium string
	Big    string
	XL     string
}

// NewAlbum create an Album from an AlbumResponse
func NewAlbum(response *AlbumResponse) (*Album, error) {
	date, err := time.Parse("2006-01-02", response.Date)
	if err != nil {
		return nil, err
	}
	album := Album{
		ID:       response.ID,
		Title:    response.Title,
		Link:     response.Link,
		CoverURL: response.CoverURL,
		Covers: Covers{
			Small:  response.CoverSmall,
			Medium: response.CoverMedium,
			Big:    response.CoverBig,
			XL:     response.CoverXL,
		},
		Date:   date,
		Tracks: response.Tracks.Data,
	}

	return &album, nil
}

// albumRequest performs the API request for the deezer album
// remember to close the body
func (api *API) albumRequest(ID int) (*http.Response, error) {
	// construct the request
	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf(AlbumAPIFormat, ID),
		nil)

	if err != nil {
		return nil, err
	}

	// send the request
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetAlbum gets the album based on its ID
func (api *API) GetAlbumData(ID int) (*Album, error) {
	// make a request to the public API
	resp, err := api.albumRequest(ID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode the json
	var response AlbumResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	// convert to album
	album, err := NewAlbum(&response)
	if err != nil {
		return nil, err
	}

	return album, nil
}
