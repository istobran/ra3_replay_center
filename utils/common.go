package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"unicode/utf16"

	"github.com/astaxie/beego/logs"
)

// RA3 Replay Constants
var (
	UTF16_TERMINATOR = []byte{0, 0}
	UTF16_CHAR_LEN   = 2
	FILE_HEADER      = "RA3 REPLAY HEADER"
	FILE_FOOTER      = "RA3 REPLAY FOOTER"
	BUFFER_SIZE      = 1536
	MAGIC_SIZE       = 17
	U1_SIZE          = 31
	U2_SIZE          = 20
	MOD_INFO_SIZE    = 22
	DATETIME_SIZE    = 8
)

var (
	ColorMap map[int]string = map[int]string{
		-1: "#000000", // Random
		0:  "#2B2BB3", // Navy
		1:  "#FCE953", // Yellow
		2:  "#00A744", // Green
		3:  "#FD7602", // Orange
		4:  "#8301FC", // Purple
		5:  "#D50000", // Red
		6:  "#04DAFA", // Cyan
	}
)

// 求出 SHA1 值
func HashFile(r multipart.File) (hash string) {
	h := sha1.New()
	if _, err := r.Seek(0, 0); err != nil {
		logs.Error(err)
	}
	if _, err := io.Copy(h, r); err != nil {
		logs.Error(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ParseInt(s string) (v int) {
	v, _ = strconv.Atoi(s)
	return
}

func ToUInt16LE(stream []byte) (ch uint16) {
	return uint16(stream[0]) + uint16(stream[1])<<8
}

func ToUInt32LE(stream []byte) (v uint32) {
	return binary.LittleEndian.Uint32(stream)
}

func ReadUTF16LE(buf *bytes.Buffer, counter *int) (str string) {
	strBuffer := make([]uint16, 0)
	for {
		ch := BufNext(buf, UTF16_CHAR_LEN, counter)
		if bytes.Equal(ch, UTF16_TERMINATOR) {
			return string(utf16.Decode(strBuffer))
		} else {
			// 小端模式解析(Little-Endian)
			strBuffer = append(strBuffer, ToUInt16LE(ch))
		}
	}
}

func ReadUInt32LE(buf *bytes.Buffer, counter *int) (v uint32) {
	return ToUInt32LE(BufNext(buf, 4, counter))
}

func BufNext(buf *bytes.Buffer, size int, counter *int) (data []byte) {
	*counter += size
	return buf.Next(size)
}

func DecodeIP(uid string) (ip string) {
	if uid == "0" {
		return uid
	}
	if len(uid) == 7 {
		uid = uid + "0"
	}
	ipv4Arr := make([]string, 4)
	for i := 0; i < len(uid); i += 2 {
		v, _ := strconv.ParseInt(uid[i:i+2], 16, 0)
		ipv4Arr[i/2] = strconv.FormatInt(v, 10)
	}
	return strings.Join(ipv4Arr, ".")
}

func DecodeColor(cvalue string) (color string) {
	v, _ := strconv.Atoi(cvalue)
	return ColorMap[v]
}
