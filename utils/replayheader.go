package utils

import (
	"regexp"
	"strings"
)

type ReplayHeader struct {
	HeaderStr        FixedLenASCII
	Hnumber1         byte
	Vermajor         uint32
	Verminor         uint32
	Buildmajor       uint32
	Buildminor       uint32
	Hnumber2         byte
	Zero1            byte
	MatchTitle       string
	MatchDescription string
	MatchMapName     string
	MatchMapId       string
	NumberOfPlayers  byte
	PlayerMap        []Player
	Offset           uint32
	StrReplLength    uint32
	StrReplMagic     FixedLenASCII
	ModInfo          FixedLenASCII
	Timestamp        uint32
	Unknown1         []byte
	HeaderLen        uint32
	Header           FixedLenASCII
	ReplaySaver      byte
	Zero2            uint32
	Zero3            uint32
	FilenameLength   uint32
	Filename         FixedLenUTF16
	DateTime         DateTime
	VermagicLen      uint32
	Vermagic         FixedLenASCII
	MagicHash        uint32
	Zero4            byte
	Unknown2         []uint32
}

func (rh *ReplayHeader) GetGameInfo() (g *GameInfo) {
	arr := strings.Split(string(rh.Header), ";")
	arr = arr[:len(arr)-1]
	for i, v := range arr {
		strl := strings.Split(v, "=")
		arr[i] = strl[1]
	}
	regM, _ := regexp.Compile("^(\\d+)(\\S+)$")
	pM := regM.FindStringSubmatch(arr[0])
	g = &GameInfo{
		M{ParseInt(pM[1]), pM[2]}, // M
		ParseInt(arr[1]),          // MC
		ParseInt(arr[2]),          // MS
		ParseInt(arr[3]),          // SD
		ParseInt(arr[4]),          // GSID
		ParseInt(arr[5]),          // PT
		ParseInt(arr[6]),          // PC
		arr[7],                    // RU
		arr[8],                    // S
	}
	return g
}
