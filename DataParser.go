package gosocket

import (
	bytes2 "bytes"
	"encoding/binary"
)

const (
	STATE_FLAG = iota
	STATE_VERSION
	STATE_TYPE
	STATE_LEN
)

type DataParser struct {
	state int
}

func (dp *DataParser) AddData(bytes *[]byte, n int) bool {
	count := 0
	for i := 0; i < n; i++ {

		b := (*bytes)[i]
		switch dp.state {
		case STATE_FLAG:
			{
				count = 0

				if n < 3 {
					return false
				}
				flag := string((*bytes)[0:3])
				if flag == "skt" {
					dp.state = STATE_VERSION
				} else {
					return false
				}
				continue
			}
		case STATE_VERSION:
			{
				if count == 0 {

				}
				count++

				if count == 2 {
					var version int16
					binary.Read(buf, binary.BigEndian, &x)
				}

			}
		}
	}

	return true
}
