package utils

import (
	"strings"

	"github.com/astaxie/beego/logs"
)

type M struct {
	MapId   int
	MapName string
}

type GameInfo struct {
	M    M      // Map Info
	MC   int    // Map CRC
	MS   int    // Map Size
	SD   int    // Seed
	GSID int    // GameSpy id
	GT   int    // unknown
	PC   int    // Post Commentator
	RU   string // Game Options
	S    string // Player Detail
}

// type PlayerDetail struct {
// 	Id       int          `json:"-" orm:"auto;pk"` // 玩家Id
// 	Name     string       `json:"name"`            // 玩家名
// 	Uid      string       `json:"uid"`             // 玩家 ip 地址的源字串
// 	Ip       string       `json:"p"`               // ip 地址
// 	Port     int          `json:"port"`            // 端口号
// 	Flag     string       `json:"flag"`            // TT|FT
// 	Color    string       `json:"color"`           // 颜色
// 	Faction  int          `json:"faction"`         // 阵营
// 	Position int          `json:"position"`        // 位置
// 	Team     int          `json:"team"`            // 队伍
// 	Handicap int          `json:"handicap"`        // 暂时未知
// 	Clan     string       `json:"clan"`            // 战队名
// 	Mode     int          `json:"mode"`            // AI 模式
// 	Human    bool         `json:"human"`           // 是否人类玩家
// 	Apm      int          `json:"apm"`             // actions per minute 平均每分钟操作数
// 	Rp       *interface{} `json:"-" orm:"rel(fk)"` // 对应的replay
// }

// type GameOption struct {
// 	Id                  int          `json:"-" orm:"auto;pk"`       // 录像配置Id
// 	InitialCameraPlayer int          `json:"initial_camera_player"` // 初始视角所在玩家
// 	GameSpeed           int          `json:"game_speed"`            // 游戏速度
// 	InitialResources    int          `json:"initial_resources"`     // 初始资金
// 	BroadcastGame       bool         `json:"broadcast_game"`        // 允许广播
// 	AllowCommentary     bool         `json:"allow_commentary"`      // 允许评论
// 	TapeDelay           int          `json:"tape_delay"`            // 启动延迟
// 	RandomCrates        bool         `json:"random_crates"`         // 随机生成箱子
// 	EnableVoIP          bool         `json:"enable_voip"`           // 允许语音
// 	Rp                  *interface{} `json:"-" orm:"reverse(one)"`  // 对应的replay
// }

func (g *GameInfo) GetPlayers() (players []map[string]interface{}) {
	players = make([]map[string]interface{}, 0)
	playerItems := strings.Split(g.S, ":")
	playerItems = playerItems[:len(playerItems)-1]
	for _, v := range playerItems {
		p := make(map[string]interface{})
		logs.Info("player", v)
		switch string(v[0]) {
		case "H": // Human
			pData := strings.Split(v, ",")
			p["Name"] = pData[0][1:]
			p["Uid"] = pData[1]
			p["Ip"] = DecodeIP(pData[1])
			p["Port"] = ParseInt(pData[2])
			p["Flag"] = pData[3]
			p["Color"] = DecodeColor(pData[4])
			p["Faction"] = ParseInt(pData[5])
			p["Position"] = ParseInt(pData[6])
			p["Team"] = ParseInt(pData[7]) + 1
			p["Handicap"] = ParseInt(pData[8])
			if len(pData) > 11 {
				p["Clan"] = pData[11]
			}
			p["Human"] = true
			players = append(players, p)
		case "C": // Computer
			pData := strings.Split(v, ",")
			p["Name"] = pData[0][1:]
			p["Color"] = DecodeColor(pData[1])
			p["Faction"] = ParseInt(pData[2])
			p["Position"] = ParseInt(pData[3])
			p["Team"] = ParseInt(pData[4]) + 1
			p["Handicap"] = ParseInt(pData[5])
			p["Mode"] = ParseInt(pData[6])
			p["Human"] = false
			players = append(players, p)
		case "X": // Closed
			continue
		}
	}
	return players
}

func (g *GameInfo) GetOptions() (opt map[string]interface{}) {
	opt = make(map[string]interface{})
	arr := strings.Split(g.RU, " ")
	opt["InitialCameraPlayer"] = ParseInt(arr[0])
	opt["GameSpeed"] = ParseInt(arr[1])
	opt["InitialResources"] = ParseInt(arr[2])
	opt["BroadcastGame"] = ParseInt(arr[3]) == 1
	opt["AllowCommentary"] = ParseInt(arr[4]) == 1
	opt["TapeDelay"] = ParseInt(arr[5])
	opt["RandomCrates"] = ParseInt(arr[6]) == 1
	opt["EnableVoIP"] = ParseInt(arr[7]) == 1
	return opt
}
