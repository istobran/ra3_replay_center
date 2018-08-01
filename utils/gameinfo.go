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
	MS   int    // Map Size
	SD   int    // Seed
	GSID int    // GameSpy id
	GT   int    // unknown
	PC   int    // Post Commentator
	RU   string // Game Options
	S    string // Player Detail
}

type PlayerDetail struct {
	Name     string `json:"name"`
	IP       string `json:"ip"`
	Port     int    `json:"port"`
	Flag     string `json:"flag"` // TT|FT
	Color    string `json:"color"`
	Army     int    `json:"army"`
	Position int    `json:"position"`
	Team     int    `json:"team"`
	Handicap int    `json:"handicap"`
	Clan     string `json:"clan"` // 战队名
	Mode     int    `json:"mode"` // AI 模式
	Human    bool   `json:"human"`
}

type GameOption struct {
	InitialCameraPlayer int  `json:"initial_camera_player"` // 初始视角所在玩家
	GameSpeed           int  `json:"game_speed"`            // 游戏速度
	InitialResources    int  `json:"initial_resources"`     // 初始资金
	BroadcastGame       bool `json:"broadcast_game"`        // 允许广播
	AllowCommentary     bool `json:"allow_commentary"`      // 允许评论
	TapeDelay           int  `json:"tape_delay"`            // 启动延迟
	RandomCrates        bool `json:"random_crates"`         // 随机生成箱子
	EnableVoIP          bool `json:"enable_voip"`           // 允许语音
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

func (g *GameInfo) GetOptions() (opt GameOption) {
	opt = GameOption{}
	arr := strings.Split(g.RU, " ")
	opt.InitialCameraPlayer = ParseInt(arr[0])
	opt.GameSpeed = ParseInt(arr[1])
	opt.InitialResources = ParseInt(arr[2])
	opt.BroadcastGame = ParseInt(arr[3]) == 1
	opt.AllowCommentary = ParseInt(arr[4]) == 1
	opt.TapeDelay = ParseInt(arr[5])
	opt.RandomCrates = ParseInt(arr[6]) == 1
	opt.EnableVoIP = ParseInt(arr[7]) == 1
	return opt
}

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
