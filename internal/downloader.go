package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/joshbarrass/deezerdl/pkg/deezer"
	"github.com/joshbarrass/deezerdl/pkg/writetracker"
	"github.com/sirupsen/logrus"
)

// https://progolang.com/how-to-download-files-in-go/

func DownloadFile(url, outPath string) error {
	// Get the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("bad status code: %d", resp.StatusCode))
	}
	defer resp.Body.Close()

	// make the file on disk to be written to
	outFile, err := os.Create(outPath + ".part")
	if err != nil {
		return err
	}
	defer outFile.Close()

	// write to the file
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
			return
		}
	}

	// check album
	if album, err := opts.Bool("album"); err != nil {
		logrus.Fatalf("failed to parse arguments: %s", err)
	} else if album {
		if err := downloadAlbum(format, ID, api); err != nil {
			logrus.Fatalf("failed to download album: %s", err)
			return
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

	filename := CalculateFilename(track, format)
	fmt.Printf("Downloading %s\n", filename)
	fmt.Println("")

	encFilename := filename + ".enc"

	// download file
	if err := DownloadFile(downloadUrl.String(), encFilename); err != nil {
		return err
	}
	defer os.Remove(encFilename)

	// decrypt song
	fmt.Println("Decrypting...")
	fmt.Println("")
	key := track.GetBlowfishKey()
	err = deezer.DecryptSongFile(key, encFilename, filename)
	if err != nil {
		return err
	}

	fmt.Println("Done!")
	return nil
}

// downloadAlbum downloads all tracks in an album
func downloadAlbum(format deezer.Format, ID int, api *deezer.API) error {
	// get album info
	fmt.Println("\nGetting album info...")
	album, err := api.GetAlbumData(ID)
	if err != nil {
		return err
	}
	fmt.Println("Got album info")
	fmt.Println("")

	// get tracks
	tracks, err := album.GetTracks()
	if err != nil {
		return err
	}

	// make new dir for the album and CD to it
	if err := os.Mkdir(album.Title, configDirPerms); err != nil {
		logrus.Warnf("couldn't make new dir -- dir possibly exists? error: %s", err)
	}
	if err := os.Chdir(album.Title); err != nil {
		return err
	}

	// download all tracks
	for index, track := range tracks {
		if err := downloadTrack(format, track.ID, api); err != nil {
			return err
		}
		// rename to have track number at front
		oldFilename := CalculateFilename(track, format)
		newFilename := fmt.Sprintf("%02d - %s", index+1, oldFilename)
		if err := os.Rename(oldFilename, newFilename); err != nil {
			logrus.Warnf("couldn't rename file: %s", err)
		}
	}

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

func FormatExtension(format deezer.Format) string {
	switch format {
	case deezer.FLAC:
		return ".flac"
	case deezer.MP3_320, deezer.MP3_256:
		return ".mp3"
	default:
		logrus.Fatalf("invalid format: %d", format)
	}
	return ""
}

func CalculateFilename(track *deezer.Track, format deezer.Format) string {
	filename := track.Title + FormatExtension(format)

	return escapeFilename(filename)
}

// escapeFilename removes illegal characters from filenames
// https://stackoverflow.com/a/31976060
func escapeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, "/", "-")
	switch runtime.GOOS {
	case "windows":
		filename = strings.ReplaceAll(filename, "<", "-")
		filename = strings.ReplaceAll(filename, ">", "-")
		filename = strings.ReplaceAll(filename, ":", "-")
		filename = strings.ReplaceAll(filename, "\"", "-")
		filename = strings.ReplaceAll(filename, "\\", "-")
		filename = strings.ReplaceAll(filename, "|", "-")
		filename = strings.ReplaceAll(filename, "?", "-")
		filename = strings.ReplaceAll(filename, "*", "-")
	}

	return filename
}
