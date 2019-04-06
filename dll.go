package main

import (
	"syscall"
	"unsafe"
)

var dll syscall.Handle
var initFun uintptr
var ocrFunc uintptr

func InitDll() error {
	var err error
	dll, err = syscall.LoadLibrary("ocr.dll")
	if err != nil {
		return err
	}
	initFun, err = syscall.GetProcAddress(dll, "init")
	if err != nil {
		return err
	}
	syscall.Syscall(initFun, 0, 0, 0, 0)
	ocrFunc, err = syscall.GetProcAddress(dll, "ocr")
	if err != nil {
		return err
	}
	return nil
}

func ocr(img *[]byte) string {
	p := *((*int32)(unsafe.Pointer(img)))
	ret, _, _ := syscall.Syscall(ocrFunc, 2, uintptr(p), uintptr(len(*img)), 0)
	return prttostr(ret)
}

func prttostr(vcode uintptr) string {
	if vcode <= 0 {
		return ""
	}
	var vbyte []byte
	for {
		sbyte := *((*byte)(unsafe.Pointer(vcode)))
		if sbyte == 0 {
			break
		}
		vbyte = append(vbyte, sbyte)
		vcode += 1
	}
	return string(vbyte)
}
