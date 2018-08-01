package models

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"ra3_replay_center/utils"
	"strconv"
	"time"
)

var (
	Replays map[string]*Replay
)

type Sizer interface {
	Size() int64
}

type Replay struct {
	Id              string               // 录像Id
	FileHash        string               // 文件 hash 值
	FileName        string               // 文件名
	FileSize        int                  // 文件大小
	NumberOfPlayers int                  // 玩家数量
	Duration        string               // 游戏时长
	GameVersion     string               // 游戏版本
	MapName         string               // 地图名称
	MapPath         string               // 地图路径
	Players         []utils.PlayerDetail // 玩家列表
	Options         utils.GameOption     // 游戏预设
	// Header          utils.ReplayHeader   // 文件头部
	// Body            utils.ReplayBody   // 文件数据
	// Footer          utils.ReplayFooter // 文件底部
}

func init() {
	Replays = make(map[string]*Replay)
}

func AddReplay(replay Replay) (Id string) {
	replay.Id = "replay_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	Replays[replay.Id] = &replay
	return replay.Id
}

func GetReplay(Id string) (replay *Replay, err error) {
	if v, ok := Replays[Id]; ok {
		return v, nil
	}
	return nil, errors.New("Replay Not Exist")
}

func GetReplayList() map[string]*Replay {
	return Replays
}

func DeleteReplay(Id string) {
	delete(Replays, Id)
}

func ResolveReplay(r multipart.File, h *multipart.FileHeader) (replay *Replay, err error) {
	rp := &Replay{}
	rh, err := utils.BuildReplayHeader(r)
	if err != nil {
		return nil, err
	}
	gi := rh.GetGameInfo()
	rp.FileHash = HashFile(r)
	rp.FileName = h.Filename
	rp.FileSize = int(r.(Sizer).Size())
	rp.NumberOfPlayers = int(rh.NumberOfPlayers)
	// rp.Duration =
	rp.GameVersion = string(rh.Vermagic)
	rp.MapName = rh.MatchMapName
	rp.MapPath = gi.M.MapName
	rp.Players = gi.GetPlayers()
	rp.Options = gi.GetOptions()
	return rp, nil
}

// 求出 SHA1 值
func HashFile(r multipart.File) (hash string) {
	h := sha1.New()
	if _, err := io.Copy(h, r); err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
