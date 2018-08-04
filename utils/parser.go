package utils

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"unicode/utf16"
)

type FixedLenASCII string
type FixedLenUTF16 string

type DateTime struct {
	Year        uint16
	Month       uint16
	Weekday     uint16 // (0-6 Sun-Sat)
	Day         uint16
	Hour        uint16
	Minute      uint16
	Second      uint16
	Millisecond uint16
}

type Player struct {
	PlayerId   uint32
	PlayerName string
	TeamNumber byte
}

func BuildReplayHeader(bytesData []byte) (rh *ReplayHeader, size int, err error) {
	headBuffer := bytesData[:BUFFER_SIZE]
	if string(headBuffer[:MAGIC_SIZE]) != FILE_HEADER {
		return nil, 0, errors.New("File is not Red Alert Replay!!")
	}
	HEADER_SIZE := 0
	rh = &ReplayHeader{}
	// 利用反射遍历结构体
	buffer := bytes.NewBuffer(headBuffer)
	v := reflect.ValueOf(rh).Elem() // 参考 https://blog.golang.org/laws-of-reflection
	fieldCount := v.NumField()
	for i := 0; i < fieldCount; i++ {
		f := v.Field(i)
		fieldName := v.Type().Field(i).Name
		switch f.Kind() {
		case reflect.Uint8:
			data := uint64(BufNext(buffer, f.Type().Align(), &HEADER_SIZE)[0])
			f.SetUint(data)
		case reflect.Uint32:
			data := uint64(ReadUInt32LE(buffer, &HEADER_SIZE))
			f.SetUint(data)
		case reflect.String:
			switch f.Type() {
			case reflect.TypeOf(FixedLenASCII(0)):
				len := 0
				switch fieldName {
				case "HeaderStr":
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
				data := string(BufNext(buffer, len, &HEADER_SIZE))
				f.Set(reflect.ValueOf(FixedLenASCII(data)))
			case reflect.TypeOf(FixedLenUTF16(0)):
				len := 0
				switch fieldName {
				case "Filename":
					len = int(rh.FilenameLength) * 2
				}
				ch := BufNext(buffer, len, &HEADER_SIZE)
				strBuffer := make([]uint16, len/2)
				for j := 0; j < len; j += 2 { // 小端读取
					strBuffer[j/2] = ToUInt16LE([]byte{ch[j], ch[j+1]})
				}
				f.Set(reflect.ValueOf(FixedLenUTF16(utf16.Decode(strBuffer))))
			default:
				line := ReadUTF16LE(buffer, &HEADER_SIZE)
				f.Set(reflect.ValueOf(line))
			}
		case reflect.Slice:
			switch fieldName {
			case "PlayerMap":
				players := make([]Player, 0)
				for j := 0; j < int(rh.NumberOfPlayers)+1; j++ {
					player := Player{}
					player.PlayerId = ReadUInt32LE(buffer, &HEADER_SIZE)
					player.PlayerName = ReadUTF16LE(buffer, &HEADER_SIZE)
					if rh.Hnumber1 == 0x05 {
						player.TeamNumber = uint8(BufNext(buffer, 1, &HEADER_SIZE)[0])
					}
					players = append(players, player)
				}
				f.Set(reflect.ValueOf(players))
			case "Unknown1":
				f.Set(reflect.ValueOf(BufNext(buffer, U1_SIZE, &HEADER_SIZE)))
			case "Unknown2":
				unknown2 := make([]uint32, 0)
				for j := 0; j < U2_SIZE; j++ {
					unknown2 = append(unknown2, ReadUInt32LE(buffer, &HEADER_SIZE))
				}
				f.Set(reflect.ValueOf(unknown2))
			}
		case reflect.Struct:
			switch fieldName {
			case "DateTime":
				len := DATETIME_SIZE * 2
				dtBuffer, dt := BufNext(buffer, len, &HEADER_SIZE), &DateTime{}
				dtKeys := reflect.ValueOf(dt).Elem()
				for j := 0; j < len; j += 2 {
					dtKeys.Field(j / 2).SetUint(uint64(ToUInt16LE([]byte{dtBuffer[j], dtBuffer[j+1]})))
				}
				f.Set(reflect.ValueOf(*dt))
			}
		}
		fmt.Println("result:", f.Type(), fieldName, f)
	}
	return rh, HEADER_SIZE, nil
}

func BuildReplayBody(bytesData []byte, offset int) (rb *ReplayBody, size int, err error) {
	BODY_SIZE := 0
	rb = &ReplayBody{}
	buffer := bytes.NewBuffer(bytesData[offset:])
	chunks := make([]BodyChunk, 0)
	for {
		timecode := ReadUInt32LE(buffer, &BODY_SIZE)
		if timecode == 0x7FFFFFFF { // timecode terminator
			break
		}
		chunktype := BufNext(buffer, 1, &BODY_SIZE)[0]
		chunksize := ReadUInt32LE(buffer, &BODY_SIZE)
		data := BufNext(buffer, int(chunksize), &BODY_SIZE)
		zero := ReadUInt32LE(buffer, &BODY_SIZE)
		chunks = append(chunks, BodyChunk{
			timecode, chunktype, chunksize, data, zero,
		})
	}
	rb.Size = len(chunks)
	rb.Chunks = chunks
	return rb, BODY_SIZE, nil
}

func BuildReplayFooter(bytesData []byte, offset int) (rf *ReplayFooter, size int, err error) {
	FOOTER_SIZE := 0
	footBuffer := bytesData[offset:]
	if string(footBuffer[:MAGIC_SIZE]) != FILE_FOOTER {
		return nil, 0, errors.New("File has crashed, Pls reupload replay!!")
	}
	buffer := bytes.NewBuffer(footBuffer)
	footerstr := string(BufNext(buffer, MAGIC_SIZE, &FOOTER_SIZE))
	finaltimecode := ReadUInt32LE(buffer, &FOOTER_SIZE)
	footerlength := ToUInt32LE(footBuffer[len(footBuffer)-4 : len(footBuffer)])
	data := BufNext(buffer, int(footerlength)-MAGIC_SIZE-12, &FOOTER_SIZE)
	rf = &ReplayFooter{
		footerstr, finaltimecode, data, footerlength,
	}
	return rf, FOOTER_SIZE, nil
}
