package flag

import (
	"fmt"
	"strconv"
	"strings"
)

type sliceInt []uint16

func (s *sliceInt) String() string {
	return fmt.Sprintf("%d", *s)
}

func (s *sliceInt) Set(value string) error {
	for _, v := range strings.Split(value, ",") {
		num, err := strconv.ParseUint(v, 10, 16)
		if err != nil {
			return fmt.Errorf("invalid port number: %w", err)
		}
		*s = append(*s, uint16(num))
	}
	return nil
}
