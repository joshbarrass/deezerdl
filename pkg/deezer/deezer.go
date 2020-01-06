package deezer

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/sirupsen/logrus"
)

const (
	getTokenMethod      = "deezer.getUserData"
	getSongMethod       = "song.getData"
	getSongMobileMethod = "song_getData"
)

var apiUrl = url.URL{
	Scheme: "https",
	Host:   "www.deezer.com",
	Path:   "/ajax/gw-light.php",
}

var mobileApiUrl = url.URL{
	Scheme: "https",
	Host:   "api.deezer.com",
	Path:   "/1.0/gateway.php",
}

var deezerUrl = url.URL{
	Scheme: "http",
	Host:   ".deezer.com",
	Path:   "/",
}

type API struct {
	APIToken  string
	client    *http.Client
	DebugMode bool
}

// NewAPI creates a new API with a http Client with cookie jar
func NewAPI(debugMode bool) (*API, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Jar: cookieJar,
	}
	api := API{
		client:    &client,
		DebugMode: debugMode,
	}
	return &api, nil
}

// ApiRequest performs an API request
func (api *API) ApiRequest(method string, body io.Reader) (*http.Response, error) {
	// add the required parameters to the URL
	u := apiUrl
	q := url.Values{
		"api_version": {"1.0"},
		"input":       {"3"},
		"method":      {method},
	}

	// add the token only if we know the token
	if method == getTokenMethod {
		q.Set("api_token", "null")
	} else {
		q.Set("api_token", api.APIToken)
	}
	u.RawQuery = q.Encode()

	// construct the request
	req, err := http.NewRequest(http.MethodPost,
		u.String(),
		body)
	if err != nil {
		return nil, err
	}

	// change user agent
	req.Header.Add("User-Agent", "PostmanRuntime/7.21.0")

	// send
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// MobileApiRequest performs a mobile API request
func (api *API) MobileApiRequest(method string, body io.Reader) (*http.Response, error) {
	// add the required parameters to the URL
	u := mobileApiUrl
	q := url.Values{
		"api_key": {"4VCYIJUCDLOUELGD1V8WBVYBNVDYOXEWSLLZDONGBBDFVXTZJRXPR29JRLQFO6ZE"},
		"input":   {"3"},
		"output":  {"3"},
		"method":  {method},
	}

	// get the current sid from the cookie jar
	var sid string
	for _, cookie := range api.client.Jar.Cookies(&deezerUrl) {
		if cookie.Name == "sid" {
			sid = cookie.Value
			break
		}
	}
	q.Set("sid", sid)

	u.RawQuery = q.Encode()

	// construct the request
	req, err := http.NewRequest(http.MethodPost,
		u.String(),
		body)
	if err != nil {
		return nil, err
	}

	// change user agent
	req.Header.Add("User-Agent", "PostmanRuntime/7.21.0")

	// send
	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// CookieLogin allows the user to log in using their arl cookie taken
// from a browser
func (api *API) CookieLogin(arl string) error {
	// add the cookie to the jar
	cookie := http.Cookie{
		Name:   "arl",
		Value:  arl,
		Domain: ".deezer.com",
		Path:   "/",
	}
	api.client.Jar.SetCookies(&deezerUrl, []*http.Cookie{&cookie})

	// get a session
	err := api.getSession()
	if err != nil {
		return err
	}

	// try to get the token
	api.APIToken, err = api.getToken()
	if err != nil {
		return err
	}

	return nil
}

// GetToken gets the user's API token
func (api *API) getToken() (string, error) {
	// make the request
	resp, err := api.ApiRequest(getTokenMethod, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if api.DebugMode {
		DumpResponse(resp, "GetToken")
	}

	// decode result key into a struct from the body
	var data struct {
		Results json.RawMessage `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil { // uses the body directly
		return "", err
	}
	// and then the checkForm key (the token)
	var results struct {
		Token string `json:"checkForm"`
	}
	if err := json.Unmarshal(data.Results, &results); err != nil {
		return "", err
	}

	if api.DebugMode {
		logrus.WithFields(logrus.Fields{
			"token": results.Token,
		}).Info("got api token")
	}
	return results.Token, nil
}

// GetSession makes a request to the base URL to get any required
// cookies
func (api *API) getSession() error {
	// construct the request
	req, err := http.NewRequest(http.MethodPost,
		"https://www.deezer.com",
		nil)
	if err != nil {
		return err
	}

	// send
	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	if api.DebugMode {
		DumpResponse(resp, "GetSession")
	}
	return nil
}
