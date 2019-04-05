package main

import (
	"syscall"
	"unsafe"
)

var ocrDll *syscall.Proc

func InitDll() error {
	dll := syscall.MustLoadDLL("ocr.dll")
	init := dll.MustFindProc("init")
	init.Call()
	ocrDll = dll.MustFindProc("ocr")
	return nil
}

func ocr(img []byte) (string, error) {
	p := *((*int32)(unsafe.Pointer(&img)))
	ret, _, err := ocrDll.Call(uintptr(p), uintptr(len(img)), 0)
	if err != nil {
		return "", err
	}
	return prttostr(ret), nil
}

func prttostr(vcode uintptr) string {
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
