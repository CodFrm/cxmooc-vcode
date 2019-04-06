package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"runtime"
	"time"
)

type Packgee struct {
	Conn    net.Conn
	Version [4]byte
	Len     int32
	Tag     [8]byte
	Time    int64
	Data    []byte
}

var ch chan *Packgee

func InitHandel() error {
	ch = make(chan *Packgee, 100)
	runtime.GOMAXPROCS(4)
	for i := 0; i < 10; i++ {
		go handelPackOcr()
	}
	return nil
}

//处理客户端连接请求
func HandelRequest(conn net.Conn) {
	defer conn.Close()
	var tempPack = make([]byte, 0, 4096)
	for {
		buf := make([]byte, 4096)
		readlen, err := conn.Read(buf)
		if err != nil {
			println("Client error: ", err.Error())
			break
		}
		tempPack = append(tempPack, buf[:readlen]...)
		if tempPack[0] == 'v' && len(tempPack) >= 8 {
			pack, err := unpack(bytes.NewReader(tempPack))
			if err == nil {
				pack.Conn = conn
				tempPack = tempPack[pack.Len:]
				ch <- pack
			} else {
				println("unpack err: ", err.Error())
			}
		} else if len(tempPack) > 1048576 {
			//大于1m的数据关闭客户链接
			break
		}
	}
}

func unpack(io io.Reader) (*Packgee, error) {
	var err error
	ret := Packgee{}
	err = binary.Read(io, binary.BigEndian, &ret.Version)
	err = binary.Read(io, binary.BigEndian, &ret.Len)
	err = binary.Read(io, binary.BigEndian, &ret.Tag)
	err = binary.Read(io, binary.BigEndian, &ret.Time)
	ret.Data = make([]byte, ret.Len-24)
	err = binary.Read(io, binary.BigEndian, &ret.Data)
	return &ret, err
}

func pack(pack *Packgee) ([]byte, error) {
	pack.Len = int32(24 + len(pack.Data))
	pack.Time = time.Now().Unix()
	var ret = make([]byte, pack.Len)
	var err error
	io := new(bytes.Buffer)
	err = binary.Write(io, binary.BigEndian, &pack.Version)
	err = binary.Write(io, binary.BigEndian, &pack.Len)
	err = binary.Write(io, binary.BigEndian, &pack.Tag)
	err = binary.Write(io, binary.BigEndian, &pack.Time)
	err = binary.Write(io, binary.BigEndian, &pack.Data)
	_, err = io.Read(ret)
	return ret, err
}

//处理ocr识别包请求
func handelPackOcr() {
	for recvPack := range ch {
		ret := ocr(&recvPack.Data)
		recvPack.Data = []byte(ret)
		send, err := pack(recvPack)
		if err != nil {
			println("pack err: ", err.Error())
			recvPack.Conn.Close()
			continue
		}
		sendLen, err := recvPack.Conn.Write(send)
		if err != nil || int32(sendLen) != recvPack.Len {
			println("send err: ", sendLen, recvPack.Len)
			recvPack.Conn.Close()
			continue
		}
	}
}
