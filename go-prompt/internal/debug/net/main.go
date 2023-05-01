package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

// netLoggerService to output message through tcp conn
type netLoggerService struct {
	net.Listener
}

func newNetLoggerService(port int) *netLoggerService {
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		panic("fail to listen, err: " + err.Error())
	}
	return &netLoggerService{
		Listener: listen,
	}
}

func (s *netLoggerService) recvLogger(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err: ", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println(recvStr)
		conn.Write([]byte(recvStr))
	}
}

var port = flag.Int("port", 9090, "specify port to listen (default value 9090)")

func main() {
	flag.Parse()
	service := newNetLoggerService(*port)
	for {
		conn, err := service.Accept()
		if err != nil {
			fmt.Println("fail to accpect conn, err: ", err)
			continue
		}
		go service.recvLogger(conn)
	}
}
