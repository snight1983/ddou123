// +build !windows

package ddlib

import (
	"syscall"
)

//GetFreeSpace ...
func GetFreeSpace(path string) (int64, error) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return 0, err
	}
	return fs.Bfree * int64(fs.Bsize), nil
}
