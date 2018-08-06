package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["ra3_replay_center/controllers:ReplayController"] = append(beego.GlobalControllerRouter["ra3_replay_center/controllers:ReplayController"],
		beego.ControllerComments{
			Method: "QueryHash",
			Router: `/queryhash`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["ra3_replay_center/controllers:ReplayController"] = append(beego.GlobalControllerRouter["ra3_replay_center/controllers:ReplayController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/upload`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

}
