package scanner

import (
	"fmt"
	"net"
	"time"
)

var dialTimeout = net.DialTimeout

func ScanPort(host string, port uint16, timeout time.Duration) error {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := dialTimeout("tcp", target, timeout)
	if err != nil {
		return err
	}

	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}
