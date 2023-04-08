package prompt

import (
	"net"
	"strconv"
)

var (
	conn       net.Conn
	debug_mode = false
)

func init() {
	if debug_mode {
		var err error
		conn, err = net.Dial("tcp", "127.0.0.1:9999")
		if err != nil {
			panic(err)
		}
	}
}

func Log(data ...any) {
	if debug_mode {
		res := ""
		for _, data := range data {
			switch v := data.(type) {
			case int:
				res += strconv.Itoa(v)
			case uint16:
				res += strconv.Itoa(int(v))
			case string:
				res += v
			}
		}
		conn.Write([]byte(res))
	}
}
