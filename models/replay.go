package models

import (
	"errors"
	"strconv"
	"time"
	"mime/multipart"
	"fmt"
)

var (
	Replays map[string]*Replay
)

type Player struct {
	Name				string
	Position		int
	Color				int
	Team				int
}

type Replay struct {
	Id					string
	Name				string
	FileName		string
	MapName			string
	GameVersion	string
	// Players			Player[]
	// PlayTime		Date
	// CreateTime	Date
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

func ResolveReplay(r multipart.File) {
	fmt.Println(r)
}