package internal

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/joshbarrass/deezerdl/pkg/deezer"
	"github.com/joshbarrass/deezerdl/pkg/writetracker"
	"github.com/sirupsen/logrus"
)

// https://progolang.com/how-to-download-files-in-go/

func DownloadFile(url, outPath string) error {
	outFile, err := os.Create(outPath + ".part")
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Get the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	wt := writetracker.NewWriteTracker("")
	_, err = io.Copy(outFile, io.TeeReader(resp.Body, wt))
	if err != nil {
		return err
	}

	// move to new line because of how ShowProgress works
	fmt.Println("")

	// rename part file
	err = os.Rename(outPath+".part", outPath)
	if err != nil {
		return err
	}

	return nil
}

// Download reads arguments from docopt options to work out what to
// download
func Download(opts docopt.Opts, config *Configuration) {
	// get the ID
	if _, ok := opts["<ID>"]; !ok {
		logrus.Fatal("missing ID")
	}
	ID, err := opts.Int("<ID>")
	if err != nil {
		logrus.Fatalf("failed to parse arguments: %s", err)
	}

	// get format
	var formatString string
	_, ok := opts["--format"]
	if ok {
		// exists, so use that
		formatString, err = opts.String("--format")
	}
	if !ok || err != nil || formatString == "" {
		// does not exist or failed, use config
		if config.DefaultFormat != "" {
			formatString = config.DefaultFormat
			err = nil
			fmt.Printf("Using format from config: %s\n", formatString)
		} else {
			logrus.Fatal("no format specified and no format in your config")
		}
	} else {
		fmt.Printf("Using format: %s\n", formatString)
	}
	if err != nil {
		logrus.Fatalf("failed to get format: %s", err)
	}
	format := FormatStringToFormat(formatString)

	// make API
	api, err := deezer.NewAPI(false)
	if err != nil {
		logrus.Fatalf("failed to create api: %s", err)
	}

	// log in
	if err := api.CookieLogin(config.ARLCookie); err != nil {
		logrus.Fatalf("failed to log in: %s", err)
	}

	// check track
	if track, err := opts.Bool("track"); err != nil {
		logrus.Fatalf("failed to parse arguments: %s", err)
	} else if track {
		if err := downloadTrack(format, ID, api); err != nil {
			logrus.Fatalf("failed to download track: %s", err)
		}
	}

}

// downloadTrack is for downloading an individual track
func downloadTrack(format deezer.Format, ID int, api *deezer.API) error {
	// get track info
	fmt.Println("\nGetting track info...")
	track, err := api.GetSongData(ID)
	if err != nil {
		return err
	}
	fmt.Println("Got track info")
	fmt.Println("")

	// get the download URL
	downloadUrl, err := track.GetDownloadURL(format)
	if err != nil {
		return err
	}

	filename := track.Title + ".flac"
	encFilename := filename + ".enc"

	// download file
	if err := DownloadFile(downloadUrl.String(), encFilename); err != nil {
		return err
	}
	defer os.Remove(encFilename)

	// decrypt song
	fmt.Println("\nDecrypting...")
	fmt.Println("")
	key := track.GetBlowfishKey()
	err = deezer.DecryptSongFile(key, encFilename, filename)
	if err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}

func FormatStringToFormat(formatString string) deezer.Format {
	var format deezer.Format
	switch formatString {
	case "FLAC":
		format = deezer.FLAC
	case "MP3_320":
		format = deezer.MP3_320
	case "MP3_256":
		format = deezer.MP3_256
	default:
		logrus.Fatalf("invalid format: %s", formatString)
	}
	return format
}
