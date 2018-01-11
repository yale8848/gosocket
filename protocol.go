package gosocket

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

const (
	STATE_FLAG = iota
	STATE_VERSION
	STATE_TYPE
	STATE_LEN
	STATE_CRC32
	STATE_DATA
)

const ENCODE_FLAG = "skt"
const HEART_BEAT = 1

type Protocol struct {
	state   int
	count   uint32
	success bool
	Version uint16
	Reserve uint16
	dataLen uint32
	crc32   uint32
	data    []byte
}

func (p *Protocol) IsHeartBeat() bool {
	return HEART_BEAT == p.Reserve
}
func convertUint16(v uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, v)
	return b
}
func convertUint32(v uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, v)
	return b
}
func (dp *Protocol) String() string {
	return string(dp.data)
}
func (dp *Protocol) GetData() []byte {
	return dp.data
}
func (dp *Protocol) Encode(byteData []byte) []byte {

	flag := []byte(ENCODE_FLAG)
	buff := bytes.NewBuffer(flag)

	version := convertUint16(dp.Version)
	reserve := convertUint16(dp.Reserve)

	buff.Write(version)
	buff.Write(reserve)
	var dataLenN uint32
	if byteData != nil && len(byteData) > 0 {
		dataLenN = uint32(len(byteData))
	}
	dataLen := convertUint32(dataLenN)
	buff.Write(dataLen)

	if byteData != nil && len(byteData) > 0 {
		crc32 := convertUint32(crc32.ChecksumIEEE(byteData))
		buff.Write(crc32)
		buff.Write(byteData)
	}
	return buff.Bytes()

}

func (dp *Protocol) Decode(byteData []byte) bool {
	buff := make([]byte, 4)
	n := len(byteData)

	init := func() {
		dp.state = STATE_FLAG
		dp.count = 0
		dp.success = false
	}
	getShort := func(b byte, state int) (uint16, bool) {
		buff[dp.count] = b
		dp.count++

		if dp.count == 2 {
			var version uint16
			version = binary.BigEndian.Uint16(buff[0:2])
			dp.state = state
			dp.count = 0
			return version, true
		}
		return 0, false
	}
	getInt := func(b byte, state int) (uint32, bool) {
		buff[dp.count] = b
		dp.count++

		if dp.count == 4 {
			var len uint32
			len = binary.BigEndian.Uint32(buff[0:4])
			dp.state = state
			dp.count = 0
			return len, true
		}
		return 0, false
	}
	for i := 0; i < n; i++ {

		b := (byteData)[i]
		switch dp.state {
		case STATE_FLAG:
			{
				dp.count = 0

				if n < len(ENCODE_FLAG) {
					init()
					return true
				}
				flag := string((byteData)[0:3])
				if flag == ENCODE_FLAG {
					dp.state = STATE_VERSION
					i += (len(ENCODE_FLAG) - 1)
				} else {
					init()
					return true
				}
				break
			}
		case STATE_VERSION:
			{

				v, finish := getShort(b, STATE_TYPE)
				if finish {
					dp.Version = v
				}
				break

			}
		case STATE_TYPE:
			{
				t, finish := getShort(b, STATE_LEN)
				if finish {
					dp.Reserve = t
				}
				break
			}

		case STATE_LEN:
			{
				len, finish := getInt(b, STATE_CRC32)
				if finish {

					if len == 0 {
						init()
						dp.data = nil
						dp.success = true
						return true
					}
					dp.dataLen = len
					dp.data = make([]byte, len)
				}
				break
			}
		case STATE_CRC32:
			{
				crc32, finish := getInt(b, STATE_DATA)
				if finish {
					dp.crc32 = crc32

				}
				break
			}

		case STATE_DATA:
			{
				dp.data[dp.count] = b
				dp.count++

				if dp.count == dp.dataLen {
					init()
					if crc32.ChecksumIEEE(dp.data) == dp.crc32 {
						dp.success = true
					}
					return true
				}

				break
			}
		}
	}

	return false
}
