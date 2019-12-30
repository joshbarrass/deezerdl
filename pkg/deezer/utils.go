package deezer

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

var combineChar = []byte("\xa4")
var deezerKey = []byte("jo6aey6haid2Teih")

// DumpResponse dumps a response with logrus
func DumpResponse(resp *http.Response, message string) {
	data, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return
	}
	logrus.WithFields(logrus.Fields{
		"response": string(data),
	}).Info(message)
}

// MD5Hash hashes the input data and returns it as a string
func MD5Hash(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

// MakeURLPath generates the path of the download URL
func MakeURLPath(track *Track, format Format) (string, error) {
	// generate MD5 data
	chars := []string{
		track.MD5,
		strconv.Itoa(int(format)),
		strconv.Itoa(track.ID),
		strconv.Itoa(track.MediaVersion),
	}
	md5Data := strings.Join(chars, string(combineChar))
	hash := []byte(MD5Hash([]byte(md5Data)))

	// generate and return hex of encrypted data
	encData := append(hash, combineChar...)
	encData = append(encData, md5Data...)
	encData = append(encData, combineChar...)
	ecb, err := ECB(deezerKey, encData)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", ecb), nil
}
