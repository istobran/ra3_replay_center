package models

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	dbname := beego.AppConfig.String("dbname")
	datasource := beego.AppConfig.String("datasource")
	runmode := beego.AppConfig.String("runmode")
	switch runmode {
	case "dev":
		orm.Debug = true
		fallthrough
	default:
		orm.RegisterDriver("postgres", orm.DRPostgres)
		orm.RegisterDataBase(dbname, "postgres", datasource)
		orm.SetMaxIdleConns(dbname, 100)
		orm.SetMaxOpenConns(dbname, 100)
	}
	orm.DefaultTimeLoc = time.FixedZone("Asia/Shanghai", 8*60*60)
	orm.RegisterModel(new(Replay))
	orm.RunSyncdb(dbname, false, orm.Debug) // 自动建表
}
