package controllers

import (
	"fmt"
	"io"
	"os"
	"ra3_replay_center/models"
	"ra3_replay_center/utils"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

type ReplayController struct {
	beego.Controller
}

const (
	SUCCESS = 0
	FAILURE = 1
)

type CommonResponse struct {
	Errcode int         `json:"errcode"`
	Errmsg  string      `json:"errmsg"`
	Data    interface{} `json:"data"`
}

// @Title upload replay
// @Description upload a ra3 replay
// @router /upload [post]
func (this *ReplayController) Post() {
	response := &CommonResponse{
		Errcode: SUCCESS,
		Errmsg:  "resolve finished",
	}
	files, err := this.GetFiles("files[]")
	if err != nil {
		ProcErr(response, "getfiles err:"+err.Error())
	} else {
		rps := make([]models.Replay, len(files))
		for i, _ := range files {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				ProcErr(response, "file open err:"+err.Error())
				break
			}
			// 不考虑 SHA1 碰撞的情况，因为几率过小，基本不可能
			tmp := models.GetReplayByHash(utils.HashFile(file))
			if tmp.Id > 0 { // 表示录像已存在，不需再次解析
				rps[i] = *tmp
			} else { // 录像不存在，需要进行解析
				rp, err := models.ResolveReplay(file, files[i])
				if err != nil {
					ProcErr(response, fmt.Sprintf("file %s resolve err: %s", files[i].Filename, err.Error()))
					break
				}
				logs.Info("builded replay model:", rp)
				dst, err := os.Create("replays/" + rp.FileHash + ".ra3replay")
				defer dst.Close()
				if err != nil {
					ProcErr(response, "file path is not writeable:"+err.Error())
					break
				}
				if _, err = file.Seek(0, 0); err != nil {
					ProcErr(response, fmt.Sprintf("file: %s failed to seek", files[i].Filename))
					break
				}
				if _, err := io.Copy(dst, file); err != nil {
					ProcErr(response, fmt.Sprintf("file: %s failed to save", files[i].Filename))
					break
				}
				models.AddReplay(rp)
				rps[i] = *rp
			}
			logs.Info("file: " + files[i].Filename + " upload successfully")
		}
		response.Data = &rps
	}
	this.Data["json"] = response
	this.ServeJSON(true)
	return
}

func ProcErr(resp *CommonResponse, msg string) {
	resp.Errcode = FAILURE
	resp.Errmsg = msg
	logs.Error(msg)
}
