//go:build !windows

package resource

import "os"

// openTempArchive opens a temporary archive for reading. openArchive removes it right after; the open file
// stays readable.
func openTempArchive(path string) (*os.File, error) {
	return os.Open(path)
}
