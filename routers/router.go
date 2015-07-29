package routers

import (
	"gtt/controllers"
	"github.com/astaxie/beego"
)

func init() {
	personCC := &controllers.PersonController{}
	personNS := beego.NewNamespace("/person",
		beego.NSRouter("/", personCC, "get:Get"),
		beego.NSRouter("/add", personCC, "post:Add"),
		beego.NSRouter("/update", personCC, "post:Update"),
		beego.NSRouter("/delete", personCC, "post:Delete"),
	)
	beego.AddNamespace(personNS)

	appCC := &controllers.AppController{}
	appNS := beego.NewNamespace("/app",
		beego.NSRouter("/", appCC, "get:Get"),
		beego.NSRouter("/add", appCC, "post:Add"),
		beego.NSRouter("/search", appCC, "get:Search"),
		beego.NSRouter("/getApp", appCC, "get:GetApp"),
	)
	beego.AddNamespace(appNS)

	jxCC := &controllers.JXController{}
	jxNS := beego.NewNamespace("/jx",
		beego.NSRouter("/", jxCC, "get:Get"),
		beego.NSRouter("/getDay", jxCC, "get:GetDay"),
		beego.NSRouter("/getMjx", jxCC, "get:GetMjx"),
		beego.NSRouter("/addMjx", jxCC, "post:AddMjx"),
		beego.NSRouter("/updateMjx", jxCC, "post:UpdateMjx"),
		beego.NSRouter("/deleteMjx", jxCC, "post:DeleteMjx"),
		beego.NSRouter("/score", jxCC, "post:Score"),
		beego.NSRouter("/xx", jxCC, "get:XX"),
		beego.NSRouter("/xxs", jxCC, "get:XXs"),
	)
	beego.AddNamespace(jxNS)

	tbCC := &controllers.TBController{}
	tbNS := beego.NewNamespace("/tb",
		beego.NSRouter("/", tbCC, "get:Get"),
		beego.NSRouter("/add", tbCC, "post:Add"),
		beego.NSRouter("/update", tbCC, "post:Update"),
		beego.NSRouter("/delete", tbCC, "post:Delete"),
		beego.NSRouter("/search", tbCC, "get:Search"),
	)
	beego.AddNamespace(tbNS)

	ggCC := &controllers.GGController{}
	ggNS := beego.NewNamespace("/gg",
		beego.NSRouter("/", ggCC, "get:Get"),
		beego.NSRouter("/ajax", ggCC, "get:Ajax"),
		beego.NSRouter("/add", ggCC, "post:Add"),
		beego.NSRouter("/update", ggCC, "post:Update"),
		beego.NSRouter("/delete", ggCC, "post:Delete"),
	)
	beego.AddNamespace(ggNS)

	monthCC := &controllers.MonthController{}
	monthNS := beego.NewNamespace("/month",
		beego.NSRouter("/", monthCC, "get:Get"),
		beego.NSRouter("/plan", monthCC, "get:Plan"),
	)
	beego.AddNamespace(monthNS)

	dayCC := &controllers.DayController{}
	dayNS := beego.NewNamespace("/day",
		beego.NSRouter("/", dayCC, "get:Get"),
		beego.NSRouter("/add", dayCC, "post:Add"),
		beego.NSRouter("/update", dayCC, "post:Update"),
		beego.NSRouter("/delete", dayCC, "post:Delete"),
		beego.NSRouter("/search", dayCC, "get:Search"),
	)
	beego.AddNamespace(dayNS)

	planCC := &controllers.PlanController{}
	planNS := beego.NewNamespace("/plan",
		beego.NSRouter("/", planCC, "get:Get"),
		beego.NSRouter("/complete", planCC, "get:GetComplete"),
		beego.NSRouter("/add", planCC, "post:Add"),
		beego.NSRouter("/update", planCC, "post:Update"),
		beego.NSRouter("/delete", planCC, "post:Delete"),
		beego.NSRouter("/getPlan", planCC, "get:GetPlan"),
	)
	beego.AddNamespace(planNS)

	loginCC := &controllers.LoginController{}
	loginNS := beego.NewNamespace("/login",
		beego.NSRouter("/", loginCC, "get:Login"),
		beego.NSRouter("/login", loginCC, "post:DoLogin"),
		beego.NSRouter("/loginOut", loginCC, "get:DoLoginOut"),
		beego.NSRouter("/md5", loginCC, "get:Md5"),
	)
	beego.AddNamespace(loginNS)
}
