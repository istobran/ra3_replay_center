package utils

import (
	"fmt"
	"strconv"
	"strings"
)

type M struct {
	MapId   int
	MapName string
}

type GameInfo struct {
	M    M
	MC   int    // Map CRC
	MS   int    // Map FileSize
	SD   int    // Seed
	GSID int    // GameSpy id
	GT   int    // unknown
	PC   int    // Post Commentator
	RU   string // Global Config
	S    string // Player Config
}

type PlayerDetail struct {
	Name     string
	IP       string
	Port     int
	Flag     string // TT|FT
	Color    string
	Army     int
	Position int
	Team     int
	Handicap int
	Clan     string // 战队名
	Mode     int    // AI 模式
	Human    bool
}

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

func (g *GameInfo) GetPlayers() (players []PlayerDetail) {
	players = make([]PlayerDetail, 0)
	playerItems := strings.Split(g.S, ":")
	playerItems = playerItems[:len(playerItems)-1]
	for _, v := range playerItems {
		p := PlayerDetail{}
		fmt.Println("player", v)
		switch string(v[0]) {
		case "H": // Human
			pData := strings.Split(v, ",")
			p.Name = pData[0][1:]
			p.IP = transformIP(pData[1])
			p.Port = ParseInt(pData[2])
			p.Flag = pData[3]
			p.Color = transformColor(pData[4])
			p.Army = ParseInt(pData[5])
			p.Position = ParseInt(pData[6])
			p.Team = ParseInt(pData[7]) + 1
			p.Handicap = ParseInt(pData[8])
			if len(pData) > 11 {
				p.Clan = pData[11]
			}
			p.Human = true
			players = append(players, p)
		case "C": // Computer
			pData := strings.Split(v, ",")
			p.Name = pData[0][1:]
			p.Color = transformColor(pData[1])
			p.Army = ParseInt(pData[2])
			p.Position = ParseInt(pData[3])
			p.Team = ParseInt(pData[4]) + 1
			p.Handicap = ParseInt(pData[5])
			p.Mode = ParseInt(pData[6])
			p.Human = false
			players = append(players, p)
		case "X": // Closed
			continue
		}
	}
	return players
}

// func (rh *ReplayHeader) ReadMap() {
// 	return
// }

// func (rh *ReplayHeader) ReadMisc() {
// 	return
// }

// func (rh *ReplayHeader) ReadOptions() {
// 	return
// }

func transformColor(cvalue string) (color string) {
	v, _ := strconv.Atoi(cvalue)
	return ColorMap[v]
}

func transformIP(uid string) (ip string) {
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
