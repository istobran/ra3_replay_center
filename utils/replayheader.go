package utils

type ReplayHeader struct {
	StrMagic 						FixedLenASCII
	Hnumber1						byte
	Vermajor						uint32
	Verminor						uint32
	Buildmajor					uint32				
	Buildminor					uint32

	Hnumber2						byte
	Zero1								byte

	MatchTitle					string
	MatchDescription		string
	MatchMapName				string
	MatchMapId					string

	NumberOfPlayers			byte
	PlayerMap						[]Player

	Offset							uint32
	StrReplLength				uint32
	StrReplMagic				FixedLenASCII

	ModInfo							FixedLenASCII

	Timestamp						uint32

	Unknown1						[]byte

	HeaderLen						uint32
	Header							FixedLenASCII
	ReplaySaver					byte
	Zero2								uint32
	Zero3								uint32
	FilenameLength			uint32
	Filename						FixedLenUTF16
	DateTime						DateTime
	VermagicLen					uint32
	Vermagic						FixedLenASCII
	MagicHash						uint32
	Zero4								byte

	Unknown2						[]uint32
}

// func (rh *ReplayHeader) ReadPlayers() {
// 	return
// }

// func (rh *ReplayHeader) ReadMap() {
// 	return
// }

// func (rh *ReplayHeader) ReadMisc() {
// 	return
// }

// func (rh *ReplayHeader) ReadOptions() {
// 	return
// }