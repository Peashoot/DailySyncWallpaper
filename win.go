package main

import (
	"syscall"
	"unsafe"
)

const (
	// SpiSerDeskWallpaper 设置桌面背景
	SpiSerDeskWallpaper = 0x0014

	// SpiFUpdateInifile Writes the new system-wide parameter setting to the user profile.
	SpiFUpdateInifile = 1

	// SpifSendWininiChange Broadcasts the WM_SETTINGCHANGE message after updating the user profile.
	SpifSendWininiChange = 2

	// FALSE 0
	FALSE = 0
	// TRUE 1
	TRUE = 1
)

var (
	// Library
	libuser32   uintptr
	libkernel32 uintptr

	// Functions
	systemParametersInfo uintptr
	getVersion           uintptr
)

func init() {
	// Library
	libuser32 = MustLoadLibrary("user32.dll")
	libkernel32 = MustLoadLibrary("kernel32.dll")
	// Functions
	systemParametersInfo = MustGetProcAddress(libuser32, "SystemParametersInfoW")
	getVersion = MustGetProcAddress(libkernel32, "GetVersion")
}

// MustLoadLibrary 加载动态库
func MustLoadLibrary(name string) uintptr {
	lib, err := syscall.LoadLibrary(name)
	if err != nil {
		panic(err)
	}

	return uintptr(lib)
}

// MustGetProcAddress 获取动态库地址
func MustGetProcAddress(lib uintptr, name string) uintptr {
	addr, err := syscall.GetProcAddress(syscall.Handle(lib), name)
	if err != nil {
		panic(err)
	}

	return uintptr(addr)
}

// SystemParametersInfo 通过调用Win32 API函数SystemParametersInfo 设置桌面壁纸
func SystemParametersInfo(uiAction, uiParam uint32, pvParam unsafe.Pointer, fWinIni uint32) bool {
	ret, _, _ := syscall.Syscall6(systemParametersInfo, 4,
		uintptr(uiAction),
		uintptr(uiParam),
		uintptr(pvParam),
		uintptr(fWinIni),
		0,
		0)

	return ret != 0
}

// GetVersion 获取windows系统版本
func GetVersion() int64 {
	ret, _, _ := syscall.Syscall(getVersion, 0,
		0,
		0,
		0)
	return int64(ret)
}
