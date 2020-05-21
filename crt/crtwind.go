// +build windows

package crt

import (
	"syscall"
	"unsafe"
)

// GetFreeSpace ...
func GetFreeSpace(path string) (int64, error) {
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")

	if err != nil {
		return 0, err
	}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

	//log.Printf("Available  %dmb", lpFreeBytesAvailable/1024/1024.0)
	//log.Printf("Total      %dmb", lpTotalNumberOfBytes/1024/1024.0)
	//log.Printf("Free       %dmb", lpTotalNumberOfFreeBytes/1024/1024.0)

	return lpTotalNumberOfFreeBytes, nil
}
