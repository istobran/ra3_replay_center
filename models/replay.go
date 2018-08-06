package models

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"ra3_replay_center/utils"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

var (
	Replays map[string]*Replay
)

type Sizer interface {
	Size() int64
}

type Replay struct {
	Id              int                      `json:"-" orm:"auto;pk"`                               // 录像Id
	FileHash        string                   `json:"file_hash"`                                     // 文件 hash 值
	FileName        string                   `json:"file_name"`                                     // 文件名
	FileSize        int                      `json:"file_size"`                                     // 文件大小
	SaveName        string                   `json:"save_name"`                                     // 录像保存名称
	NumberOfPlayers int                      `json:"number_of_players"`                             // 玩家数量
	Duration        int                      `json:"duration"`                                      // 游戏时长
	GameVersion     string                   `json:"game_version"`                                  // 游戏版本
	MapName         string                   `json:"map_name"`                                      // 地图名称
	MapPath         string                   `json:"map_path"`                                      // 地图路径
	Players         []map[string]interface{} `json:"players" orm:"-"`                               // 玩家列表
	PlayersJson     string                   `json:"-" orm:"type(json)"`                            // JSON 字符串格式的玩家列表
	Options         map[string]interface{}   `json:"options" orm:"-"`                               // 游戏预设
	OptionsJson     string                   `json:"-" orm:"type(json)"`                            // 游戏预设
	HeaderLen       int                      `json:"-"`                                             // 头部大小
	BodyLen         int                      `json:"-"`                                             // 数据体大小
	FooterLen       int                      `json:"-"`                                             // 底部大小
	CreateTime      time.Time                `json:"create_time" orm:"auto_now_add;type(datetime)"` // 创建时间
	// Uploader        string               `json:"uploader"`          // 上传者名称
	// Email           string               `json:"email"`             // 上传者邮箱
}

var o orm.Ormer

func init() {
	o = orm.NewOrm()
}

func AddReplay(replay *Replay) (Id int) {
	id, err := o.Insert(replay)
	if err != nil {
		logs.Error(err)
		return -1
	}
	return int(id)
}

func GetReplayByHash(hash string) (replay *Replay) {
	replay = &Replay{FileHash: hash}
	err := o.Read(replay, "FileHash")
	if err != nil {
		logs.Notice(err)
		return replay
	}
	var players []map[string]interface{}
	if err := json.Unmarshal([]byte(replay.PlayersJson), &players); err != nil {
		logs.Error(err)
		return replay
	}
	replay.Players = players
	var options map[string]interface{}
	if err := json.Unmarshal([]byte(replay.OptionsJson), &options); err != nil {
		logs.Error(err)
		return replay
	}
	replay.Options = options
	return replay
}

func ResolveReplay(r multipart.File, h *multipart.FileHeader) (replay *Replay, err error) {
	rp := &Replay{}
	if _, err := r.Seek(0, 0); err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	rh, hsize, err := utils.BuildReplayHeader(data)
	if err != nil {
		return nil, err
	}
	rp.HeaderLen = hsize
	rb, bsize, err := utils.BuildReplayBody(data, hsize)
	if err != nil {
		return nil, err
	}
	rp.BodyLen = bsize
	rf, fsize, err := utils.BuildReplayFooter(data, hsize+bsize)
	if err != nil {
		return nil, err
	}
	rp.FooterLen = fsize
	gi := rh.GetGameInfo()
	rp.FileHash = utils.HashFile(r)
	rp.FileName = h.Filename
	rp.FileSize = int(r.(Sizer).Size())
	rp.SaveName = string(rh.Filename)
	rp.NumberOfPlayers = int(rh.NumberOfPlayers)
	rp.Duration = rf.GetDuration()
	rp.GameVersion = string(rh.Vermagic)
	rp.MapName = rh.MatchMapName
	rp.MapPath = gi.M.MapName
	rp.Players = rb.CalcAPM(gi.GetPlayers())
	serializedPlayers, err := json.Marshal(rp.Players)
	if err != nil {
		return nil, err
	}
	rp.PlayersJson = string(serializedPlayers)
	rp.Options = gi.GetOptions()
	serializedOptions, err := json.Marshal(gi.GetOptions())
	if err != nil {
		return nil, err
	}
	rp.OptionsJson = string(serializedOptions)
	return rp, nil
}
