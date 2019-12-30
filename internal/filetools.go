package internal

import "os"

// https://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-exists

// FileExists checks whether a file/directory exists
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
