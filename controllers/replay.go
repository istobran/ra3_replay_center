package controllers

import (
	"io"
	"os"
	"ra3_replay_center/models"

	"github.com/astaxie/beego"
)

type ReplayController struct {
	beego.Controller
}

type CommonResponse struct {
	Errcode int         `json:"errcode"`
	Errmsg  string      `json:"errmsg"`
	Data    interface{} `json:"data"`
}

// TODO:
// 需要做文件大小检测（InsertFilter）
// 需要做 hash 重复上传检测

// @Title upload replay
// @Description upload a ra3 replay
// @router /upload [post]
func (this *ReplayController) Post() {
	files, err := this.GetFiles("files[]")
	if err != nil {
		beego.Error("getfiles err ", err)
	}
	rps := make([]models.Replay, len(files))
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			beego.Error("file open err:", err)
			break
		}
		rp, err := models.ResolveReplay(file, files[i])
		if err != nil {
			beego.Error(files[i].Filename, "resolve err:", err)
			break
		}
		beego.Informational("builded replay model:", rp)
		dst, err := os.Create("replays/" + rp.FileHash + ".ra3replay")
		defer dst.Close()
		if err != nil {
			beego.Error("file path is not writeable:", err)
			break
		}
		if _, err = file.Seek(0, 0); err != nil {
			beego.Error("file: " + files[i].Filename + " failed to seek")
			break
		}
		if _, err := io.Copy(dst, file); err != nil {
			beego.Error("file: " + files[i].Filename + " failed to save")
			break
		}
		rps[i] = *rp
		beego.Informational("file: " + files[i].Filename + " upload successfully")
	}
	this.Data["json"] = &CommonResponse{
		Errcode: 0,
		Errmsg:  "resolve finished",
		Data:    &rps,
	}
	this.ServeJSON(true)
	return
}
