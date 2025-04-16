package scanner

import (
	"fmt"
	"net"
	"time"
)

func ScanPort(host string, port uint16, timeout time.Duration) error {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}
