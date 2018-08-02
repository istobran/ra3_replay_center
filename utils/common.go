package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"unicode/utf16"

	"github.com/astaxie/beego"
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

// 求出 SHA1 值
func HashFile(r multipart.File) (hash string) {
	h := sha1.New()
	if _, err := r.Seek(0, 0); err != nil {
		beego.Error(err)
	}
	if _, err := io.Copy(h, r); err != nil {
		beego.Error(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ParseInt(s string) (v int) {
	v, _ = strconv.Atoi(s)
	return
}

func ToUInt16LE(charas []byte) (ch uint16) {
	return uint16(charas[0]) + uint16(charas[1])<<8
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
	return binary.LittleEndian.Uint32(BufNext(buf, 4, counter))
}

func BufNext(buf *bytes.Buffer, size int, counter *int) (data []byte) {
	*counter += size
	return buf.Next(size)
}
