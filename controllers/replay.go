package controllers

import (
	"github.com/astaxie/beego"
	"io"
	"os"
)

type ReplayController struct {
	beego.Controller
}

// @Title upload replay
// @Description upload a ra3 replay
// @router /upload [post]
func (c *ReplayController) Post() {
	files, err := c.GetFiles("files[]")
	if err != nil {
		beego.Error("getfiles err ", err)
	}
	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			beego.Error("file open err ", err)
		}
		//create destination file making sure the path is writeable.
		dst, err := os.Create("replays/" + files[i].Filename)
		defer dst.Close()
		if err != nil {
			beego.Error("file path is not writeable ", err)
		}
		if _, err := io.Copy(dst, file); err != nil {
			beego.Error("file: " + files[i].Filename + " failed to save")
		}
		beego.Informational("file: " + files[i].Filename + " upload successfully")
	}
	c.Data["json"] = &map[string]interface{}{"errcode": 0, "message": "upload successfully"}
	c.ServeJSON()
}