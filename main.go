package main

import (
	"encoding/json"
	"os"
	_ "ra3_replay_center/routers"
	"ra3_replay_center/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

const (
	LOGPATH = "logs/"
	REPPATH = "replays/"
)

func init() {
	if _, err := os.Stat(LOGPATH); os.IsNotExist(err) {
		os.Mkdir(LOGPATH, os.ModePerm)
	}
	if _, err := os.Stat(REPPATH); os.IsNotExist(err) {
		os.Mkdir(REPPATH, os.ModePerm)
	}
	logs.SetLogger(logs.AdapterFile, `{
		"filename": "logs/replay_center_api.log",
		"separate": ["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"],
		"daily": true,
		"maxdays": 30,
		"rotate": true
	}`)
	// mailsource, err := json.Marshal(map[string]interface{}{
	// 	"username":    beego.AppConfig.String("mail_username"),
	// 	"password":    beego.AppConfig.String("mail_password"),
	// 	"fromAddress": beego.AppConfig.String("mail_username"),
	// 	"host":        beego.AppConfig.String("mail_host"),
	// 	"sendTos":     []string{beego.AppConfig.String("mail_sendtos")},
	// 	"subject":     "RA3 Replay Center Error Report",
	// 	"level":       logs.LevelError,
	// })
	// if err != nil {
	// 	logs.Error(err)
	// }
	// logs.SetLogger(logs.AdapterMail, string(mailsource))
	qcConfig, err := json.Marshal(map[string]interface{}{
		"policyId":  beego.AppConfig.String("policy_id"),
		"region":    beego.AppConfig.String("region"),
		"secretId":  beego.AppConfig.String("secret_id"),
		"secretKey": beego.AppConfig.String("secret_key"),
		"level":     logs.LevelError,
	})
	if err != nil {
		logs.Error(err)
	}
	logs.SetLogger(utils.AdapterQCloud, string(qcConfig))
	logs.EnableFuncCallDepth(true)
	logs.Async()
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
