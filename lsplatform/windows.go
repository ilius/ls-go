//go:build windows

package lsplatform

import (
	"io"
	"os"
	"syscall"
	"time"
	"unsafe"

	colorable "github.com/ilius/go-colorable"
)

var (
	advapi32                       = syscall.NewLazyDLL("advapi32.dll")
	procGetFileSecurity            = advapi32.NewProc("GetFileSecurityW")
	procGetSecurityDescriptorOwner = advapi32.NewProc("GetSecurityDescriptorOwner")
)

const (
	OWNER_SECURITY_INFORMATION = 0x00000001
	ERROR_INSUFFICIENT_BUFFER  = "The data area passed to a system call is too small."
)

func (*LocalPlatform) OwnerAndGroupNames(info FileInfo) (*OwnerGroup, error) {
	path := info.PathAbs()
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return nil, &PlatformError{
			Operation: "syscall.UTF16PtrFromString",
			Path:      path,
			Msg:       err.Error(),
		}
	}

	var lenNeeded uint32
	r1, _, err := procGetFileSecurity.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		OWNER_SECURITY_INFORMATION,
		0,                                   // pSecurityDescriptor
		0,                                   // nLength
		uintptr(unsafe.Pointer(&lenNeeded)), // lpnLengthNeeded
	)
	if r1 == 0 && err != nil { // lenNeeded == 0
		if err.Error() != ERROR_INSUFFICIENT_BUFFER {
			return nil, &PlatformError{
				Operation: "GetFileSecurity",
				Path:      path,
				Msg:       err.Error(),
			}
		}
	}
	if lenNeeded == 0 {
		return nil, &PlatformError{
			Operation: "OwnerAndGroupNames",
			Path:      path,
			Msg:       "unexpected lenNeeded = 0",
		}
	}
	pSecDesc := make([]byte, lenNeeded)
	var lenNeeded2 uint32
	r1, _, err = procGetFileSecurity.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		OWNER_SECURITY_INFORMATION,
		uintptr(unsafe.Pointer(&pSecDesc[0])), // pSecurityDescriptor
		uintptr(lenNeeded),                    // nLength
		uintptr(unsafe.Pointer(&lenNeeded2)),  // lpnLengthNeeded
	)
	if r1 == 0 && err != nil {
		return nil, &PlatformError{
			Operation: "GetFileSecurity (2)",
			Path:      path,
			Msg:       err.Error(),
		}
	}
	var ownerDefaulted uint32
	var sid *syscall.SID
	// The security identifier (SID) structure is a variable-length structure used
	// to uniquely identify users or groups.
	r1, _, err = procGetSecurityDescriptorOwner.Call(
		uintptr(unsafe.Pointer(&pSecDesc[0])),
		uintptr(unsafe.Pointer(&sid)),
		uintptr(unsafe.Pointer(&ownerDefaulted)),
	)
	// seems like ownerDefaulted (the last argument) is always zero
	// but you have to pass it or it will panic
	if r1 == 0 && err != nil {
		return nil, &PlatformError{
			Operation: "GetSecurityDescriptorOwner",
			Path:      path,
			Msg:       err.Error(),
		}
	}
	// LookupAccount(system string) (account string, domain string, accType uint32, err error)
	// what is accType exactly?
	// I have seen 1 (normal user/admin), 4 (root user?), and 5 (system user?)
	account, domain, _, err := sid.LookupAccount("")
	if err != nil {
		return nil, &PlatformError{
			Operation: "LookupAccount",
			Path:      path,
			Msg:       err.Error(),
		}
	}
	// log.Printf("account=%#v, domain=%#v", account, domain)
	// account="Administrators", domain="BUILTIN", accType=4
	// account="SYSTEM", domain="NT AUTHORITY", accType=5
	// account="TrustedInstaller", domain="NT SERVICE", accType=5
	// account = fmt.Sprintf("%v (%v)", account, accType)
	return &OwnerGroup{
		account,
		domain,
	}, nil
}

// RootUserName returns name of root user (the main admin)
func (*LocalPlatform) RootUserName() string {
	return "Administrators"
}

func (*LocalPlatform) SystemUserNames() []string {
	return []string{
		"SYSTEM",
		"TrustedInstaller",
	}
}

// UserName returns name of current user
func (*LocalPlatform) UserName() string {
	return os.Getenv("USERNAME")
}

func (f *LocalPlatform) OwnerAndGroupIDs(info FileInfo) (*OwnerGroup, error) {
	return f.OwnerAndGroupNames(info)
}

func (*LocalPlatform) DeviceNumbers(_ FileInfo) (string, error) {
	return "", nil
}

func (*LocalPlatform) NumberOfHardLinks(info FileInfo) (uint64, error) {
	if info.IsDir() {
		return 0, nil
	}
	path := info.PathAbs()
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, &PlatformError{
			Operation: "syscall.UTF16PtrFromString",
			Path:      path,
			Msg:       err.Error(),
		}
	}
	handle, err := syscall.CreateFile(pathPtr, 0, 0, nil, syscall.OPEN_EXISTING, 0, 0)
	if err != nil {
		return 0, &PlatformError{
			Operation: "syscall.CreateFile",
			Path:      path,
			Msg:       err.Error(),
		}
	}

	var fi syscall.ByHandleFileInformation
	if err = syscall.GetFileInformationByHandle(handle, &fi); err != nil {
		syscall.CloseHandle(handle)
		return 0, &PlatformError{
			Operation: "syscall.GetFileInformationByHandle",
			Path:      path,
			Msg:       err.Error(),
		}
	}
	return uint64(fi.NumberOfLinks), nil
}

func (*LocalPlatform) FileInode(info FileInfo) (uint64, error) {
	path := info.PathAbs()
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, &PlatformError{
			Operation: "syscall.UTF16PtrFromString",
			Path:      path,
			Msg:       err.Error(),
		}
	}
	attrs := uint32(0)
	if info.IsDir() {
		attrs = 0x02000000
	}
	handle, err := syscall.CreateFile(
		pathPtr,               // name
		0,                     // access
		0,                     // mode
		nil,                   // sa *SecurityAttributes
		syscall.OPEN_EXISTING, // createmode
		attrs,                 // attrs
		0,                     // templatefile
	)
	if err != nil {
		return 0, &PlatformError{
			Operation: "syscall.CreateFile",
			Path:      path,
			Msg:       err.Error(),
		}
	}

	var fi syscall.ByHandleFileInformation
	if err = syscall.GetFileInformationByHandle(handle, &fi); err != nil {
		syscall.CloseHandle(handle)
		return 0, &PlatformError{
			Operation: "syscall.GetFileInformationByHandle",
			Path:      path,
			Msg:       err.Error(),
		}
	}

	return uint64(fi.FileIndexHigh)<<32 | uint64(fi.FileIndexLow), nil
}

func (*LocalPlatform) FileCTime(info FileInfo) *time.Time {
	data := info.Sys().(*syscall.Win32FileAttributeData)
	_time := time.Unix(0, data.LastWriteTime.Nanoseconds())
	return &_time
}

func (*LocalPlatform) FileATime(info FileInfo) *time.Time {
	data := info.Sys().(*syscall.Win32FileAttributeData)
	_time := time.Unix(0, data.LastAccessTime.Nanoseconds())
	return &_time
}

// FileBlocks returns number of 1024-byte blocks occupied by a file
func (*LocalPlatform) FileBlocks(_ FileInfo) int64 {
	// FIXME
	// data := info.Sys().(*syscall.Win32FileAttributeData)
	// data.FileSizeHigh is always zero
	// data.FileSizeLow is same as file.Size(), not "size on disk"
	return 0
}

func (*LocalPlatform) EmptyFileInfoSys() any {
	return &syscall.Win32FileAttributeData{}
}

func (*LocalPlatform) OutputAndError(colors bool) (io.Writer, io.Writer) {
	if !colors {
		return os.Stdout, os.Stderr
	}
	if os.Getenv("TERM") != "" {
		// "xterm" or "xterm-256color" for Unix
		// "xterm" for Git-Bash on Windows
		// empty for Windows cmd
		return os.Stdout, os.Stderr
	}
	return colorable.NewColorableStdout(), colorable.NewColorableStderr()
}
