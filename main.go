package main

import (
	"C"
	"net"
	"os"
)

func main() {
	if err := InitDll(); err != nil {
		println("Identification module failed to load: ", err.Error())
		return
	}
	if err := InitHandel(); err != nil {
		println("Handel init error: ", err.Error())
		return
	}
	//开启服务
	server, err := net.Listen("tcp", ":5208")
	if err != nil {
		println("Server start error: ", err.Error())
		return
	}
	defer server.Close()
	println("Service start...")
	for {
		conn, err := server.Accept()
		if err != nil {
			println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go HandelRequest(conn)
	}
	println("Service stop")
}
