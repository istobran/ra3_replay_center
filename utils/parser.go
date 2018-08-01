package utils

import (
	"mime/multipart"
	"errors"
	"reflect"
	"encoding/binary"
	"bytes"
	"fmt"
	"unicode/utf16"
)

var (		// RA3 Replay Constants
	UTF16_CHAR_LEN		= 		2
	UTF16_TERMINATOR	= 		[]byte{0, 0}
	FILE_HEADER				=			"RA3 REPLAY HEADER"
	BUFFER_SIZE				=			1536
	MAGIC_SIZE				=			17
	U1_SIZE						=			31
	U2_SIZE						=			20
	MOD_INFO_SIZE			=			22
	DATETIME_SIZE			=			8
)

type FixedLenASCII string
type FixedLenUTF16 string

type DateTime struct {
	Year				uint16
	Month				uint16
	Weekday			uint16		// (0-6 Sun-Sat)
	Day					uint16
	Hour				uint16
	Minute			uint16
	Second			uint16
	Millisecond	uint16
}

type Player struct {
	PlayerId			uint32
	PlayerName		string
	TeamNumber		byte
}

func BuildReplayHeader(f multipart.File) (rh *ReplayHeader, err error) {
	headBuffer := make([]byte, BUFFER_SIZE)
	_, err = f.Read(headBuffer)
	if err != nil {
		return nil, err
	}
	if string(headBuffer[:MAGIC_SIZE]) != FILE_HEADER {
		return nil, errors.New("File is not Red Alert Replay!!")
	}
	return Parse(headBuffer), nil
}

func Parse(bytesBuffer []byte) (rh *ReplayHeader) {
	rh = &ReplayHeader{}
	// 利用反射遍历结构体
	buffer := bytes.NewBuffer(bytesBuffer)
	v := reflect.ValueOf(rh).Elem()			// 参考 https://blog.golang.org/laws-of-reflection
	fieldCount := v.NumField()
	for i := 0; i < fieldCount; i++ {
		f := v.Field(i)
		fieldName := v.Type().Field(i).Name
		switch f.Kind() {
			case reflect.Uint8:
				data := uint64(buffer.Next(f.Type().Align())[0])
				f.SetUint(data)
			case reflect.Uint32:
				buf := buffer.Next(f.Type().Align())
				data := uint64(binary.LittleEndian.Uint32(buf))
				f.SetUint(data)
			case reflect.String:
				switch f.Type() {
					case reflect.TypeOf(FixedLenASCII(0)):
						len := 0
						switch fieldName {
							case "StrMagic":
								len = MAGIC_SIZE
							case "StrReplMagic":
								len = int(rh.StrReplLength)
							case "ModInfo":
								len = MOD_INFO_SIZE
							case "Header":
								len = int(rh.HeaderLen)
							case "Vermagic":
								len = int(rh.VermagicLen)
						}
						data := string(buffer.Next(len))
						f.Set(reflect.ValueOf(FixedLenASCII(data)))
					case reflect.TypeOf(FixedLenUTF16(0)):
						len := 0
						switch fieldName {
							case "Filename":
								len = int(rh.FilenameLength)*2
						}
						ch := buffer.Next(len)
						strBuffer := make([]uint16, len/2)
						for j := 0; j < len; j += 2 {	// 小端读取
							strBuffer[j/2] = ToUInt16LE([]byte{ch[j], ch[j+1]})
						}
						f.Set(reflect.ValueOf(FixedLenUTF16(utf16.Decode(strBuffer))))
					default:
						line := ReadUTF16LE(buffer);
						f.Set(reflect.ValueOf(line));
				}
			case reflect.Slice:
				switch fieldName {
					case "PlayerMap":
						players := make([]Player, 0)
						for j := 0; j < int(rh.NumberOfPlayers) + 1; j++ {
							player := Player{}
							player.PlayerId = binary.LittleEndian.Uint32(buffer.Next(4))
							player.PlayerName = ReadUTF16LE(buffer)
							if rh.Hnumber1 == 0x05 {
								player.TeamNumber = uint8(buffer.Next(1)[0])
							}
							players = append(players, player)
						}
						f.Set(reflect.ValueOf(players))
					case "Unknown1":
						f.Set(reflect.ValueOf(buffer.Next(U1_SIZE)))
					case "Unknown2":
						unknown2 := make([]uint32, 0)
						for j := 0; j < U2_SIZE; j++ {
							unknown2 = append(unknown2, binary.LittleEndian.Uint32(buffer.Next(4)))
						}
						f.Set(reflect.ValueOf(unknown2))
				}
			case reflect.Struct:
				switch fieldName {
					case "DateTime":
						len := DATETIME_SIZE*2
						dtBuffer, dt := buffer.Next(len), &DateTime{}
						dtKeys := reflect.ValueOf(dt).Elem()
						for j := 0; j < len; j += 2 {
							dtKeys.Field(j/2).SetUint(uint64(ToUInt16LE([]byte{dtBuffer[j], dtBuffer[j+1]})))
						}
						f.Set(reflect.ValueOf(*dt))
				}
		}
		fmt.Println("result:", f.Type(), fieldName, f)
	}
	return rh
}

func ToUInt16LE(charas []byte) (ch uint16) {
	return uint16(charas[0]) + uint16(charas[1]) << 8
}

func ReadUTF16LE(buf *bytes.Buffer) (str string) {
	strBuffer := make([]uint16, 0)
	for {
		ch := buf.Next(UTF16_CHAR_LEN)		
		if bytes.Equal(ch, UTF16_TERMINATOR) {
			return string(utf16.Decode(strBuffer))
		} else {
			// 小端模式解析(Little-Endian)
			strBuffer = append(strBuffer, ToUInt16LE(ch))
		}
	}
}