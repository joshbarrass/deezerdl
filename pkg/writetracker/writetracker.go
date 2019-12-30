package writetracker

import (
	"fmt"
	"strings"

	humanize "github.com/dustin/go-humanize"
)

// https://progolang.com/how-to-download-files-in-go/

type WriteTracker struct {
	bytes  uint64
	Format string
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
	fmt.Printf(tracker.Format, humanize.Bytes(tracker.bytes))
}

// NewWriteTracker returns a pointer to a new write tracker with a
// custom format. Use an empty string for a default "downloaded"
// formatter. %s is used for the bytes.
func NewWriteTracker(format string) *WriteTracker {
	if format == "" {
		format = "\rDownloaded %s"
	}
	return &WriteTracker{
		Format: format,
	}
}
