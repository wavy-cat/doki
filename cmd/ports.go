package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

type portsFlag []uint16

func (p *portsFlag) String() string {
	if p == nil || len(*p) == 0 {
		return ""
	}
	values := make([]string, 0, len(*p))
	for _, v := range *p {
		values = append(values, strconv.FormatUint(uint64(v), 10))
	}
	return strings.Join(values, ",")
}

func (p *portsFlag) Set(value string) error {
	if value == "" {
		return fmt.Errorf("ports value is empty")
	}
	for _, v := range strings.Split(value, ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			return fmt.Errorf("invalid port number: empty")
		}
		num, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return fmt.Errorf("invalid port number: %w", err)
		}
		*p = append(*p, uint16(num))
	}
	return nil
}

func (p *portsFlag) Type() string {
	return "ports"
}
