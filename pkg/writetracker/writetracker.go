package writetracker

import (
	"fmt"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// https://progolang.com/how-to-download-files-in-go/

type WriteTracker struct {
	bytes uint64
}

// Write takes the input bytes and increments the total number of
// bytes written
func (tracker *WriteTracker) Write(b []byte) (int, error) {
	n := len(b)
	tracker.bytes += uint64(n)
	tracker.ShowProgress()
	return n, nil
}

// ShowProgress displays the current number of bytes written
func (tracker *WriteTracker) ShowProgress() {
	// reset line
	fmt.Printf("\r%s", strings.Repeat(" ", 50))

	// print current progress
	fmt.Printf("\rDownloaded %s", humanize.Bytes(tracker.bytes))
}
