package resource

import (
	"os"
	"syscall"
)

// openTempArchive opens a temporary archive for reading. Windows cannot remove an open file, so the system
// removes it once this handle closes, which includes the process ending in any way.
func openTempArchive(path string) (*os.File, error) {
	name, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	const deleteOnClose = 0x04000000 // FILE_FLAG_DELETE_ON_CLOSE
	h, err := syscall.CreateFile(name, syscall.GENERIC_READ, syscall.FILE_SHARE_READ|syscall.FILE_SHARE_DELETE, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL|deleteOnClose, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	return os.NewFile(uintptr(h), path), nil
}
