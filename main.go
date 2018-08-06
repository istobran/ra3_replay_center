package main

import (
	"encoding/json"
	_ "ra3_replay_center/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func init() {
	logs.SetLogger(logs.AdapterFile, `{
		"filename": "logs/replay_center_api.log",
		"separate": ["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
		"daily": true,
		"maxdays": 30,
		"rotate": true
	}`)
	mailsource, err := json.Marshal(map[string]interface{}{
		"username":    beego.AppConfig.String("mail_username"),
		"password":    beego.AppConfig.String("mail_password"),
		"fromAddress": beego.AppConfig.String("mail_username"),
		"host":        beego.AppConfig.String("mail_host"),
		"sendTos":     []string{beego.AppConfig.String("mail_sendtos")},
		"subject":     "RA3 Replay Center Error Report",
		"level":       logs.LevelError,
	})
	if err != nil {
		logs.Error(err)
	}
	logs.SetLogger(logs.AdapterMail, string(mailsource))
	logs.EnableFuncCallDepth(true)
	logs.Async(1e3)
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
