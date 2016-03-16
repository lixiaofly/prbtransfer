package prbtransfer

import (
	"errors"
	"fmt"
	"net"
	"io/ioutil"
)

func TcpUpload(data []byte, tcpServ *net.TCPAddr) (int, error) {
	fmt.Println("tcp send begin")
	conn, err := net.DialTCP("tcp", nil, tcpServ)
	if err != nil {
		return 0, err
	}
	send := 0
	size := len(data)
	if send < size {
		fmt.Println("first send:", data[send:],"len=", len(data))
		i, err := conn.Write(data[send:])
		if err != nil {
			return send, err
		}
		send += i
	}
	fmt.Println("tcp send succeed!")
	result, err := ioutil.ReadAll(conn)
    checkError(err)
    fmt.Println(string(result))
	conn.Close()
	return send, errors.New("tcp send succeed!")
}

func UdpUpload(data []byte, udpServ *net.UDPAddr) (int, error) {
	//udpAddr, err := net.ResolveUDPAddr("udp4", service)
	conn, err := net.DialUDP("udp", nil, udpServ)
	if err != nil {
		return 0,err
	}
	send := 0
	size := len(data)
	if send < size {
		i, err := conn.Write(data[send:])
		if err != nil {
			return send, err
		}
		send += i
	}
	conn.Close()
	return send, errors.New("Send succeed!")
}

func checkError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
	}
}
