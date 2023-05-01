package debug

import (
	"net"
	"os"
	"strconv"
)

// netLogger to output message through tcp conn
type netLogger struct {
	conn net.Conn
}

var (
	netlogger       *netLogger
	envEnableNetLog = "GO_PROMPT_ENABLE_NET_LOG"
)

func init() {
	if e := os.Getenv(envEnableNetLog); e == "true" || e == "1" {
		conn, err := net.Dial("tcp", "127.0.0.1:9999")
		if err != nil {
			panic(err)
		}
		netlogger = &netLogger{
			conn: conn,
		}
	}
}

// NetLog to output message
func NetLog(data ...any) {
	if netlogger != nil {
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
		netlogger.conn.Write([]byte(res))
	}
}
