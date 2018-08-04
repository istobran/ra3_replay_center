package utils

type ReplayFooter struct {
	FooterStr     string
	FinalTimeCode uint32
	Data          []byte
	FooterLength  uint32
}

func (rf *ReplayFooter) GetDuration() (dur int) {
	return int(rf.FinalTimeCode) / 15
}
