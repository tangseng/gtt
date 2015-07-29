package controllers

import (
	"gtt/models"
	"strconv"
	"github.com/astaxie/beego"
)

type AppController struct {
	BaseController
}

func (this *AppController) check() bool {
	return true
}

func (this *AppController) Get() {
	if !this.check() {
		return
	}
	outApps, err := this.app()
	if err != nil {
		this.Redirect("/login", 302)
		return
	}
	this.Data["Apps"] = outApps
	this.Data["PageApp"] = true
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/directives/date", "app/services/time", "app/controllers/app/appCtrl"}
	this.TplNames = "app.html"
}

func (this *AppController) Search() {
	if !this.check() {
		return
	}
	appid, _ := this.GetInt("appid")
	page, _ := this.GetInt("page")
	offset, _ := this.GetInt("offset")
	logs, pps := this.get(appid, page, offset)
	this.success(map[string]interface {}{"logs" : logs, "pps": pps})
}

func (this *AppController) get(appid, page, offset int) ([]models.App, map[string]string) {
	logs := models.NewAppOption().ReadMore(appid, offset, page)
	pids := make([]int, 0)
	pidsCheck := make(map[string]bool)
	for _, tb := range logs {
		if tb.Uid != -1 {
			if _, ok := pidsCheck[strconv.Itoa(tb.Uid)]; !ok {
				pidsCheck[strconv.Itoa(tb.Uid)] = true;
				pids = append(pids, tb.Uid)
			}
		}
	}
	outPersons := make(map[string]string)
	if len(pids) > 0 {
		persons := models.NewPersonOption().ReadByIds(pids)
		for _, person := range persons {
			outPersons[strconv.Itoa(person.Id)] = person.Name
		}
	}
	if len(outPersons) > 0 {
		outPersons["-1"] = "管理员"
	}
	return logs, outPersons
}

func (this *AppController) Add() {
	if !this.check() {
		return
	}
	this.info = &map[string]string{
		"post" : "内容必须填写",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	appid, _ := this.GetInt("appid")
	desc := this.GetString("desc")
	if len(desc) == 0 {
		this.error("post")
		return
	}
	id, err := models.NewAppOption().Insert(models.NewApp(this.user.Id, appid, desc))
	if err != nil {
		this.error("err")
		return
	}
	this.success(map[string]interface{} {"id" : id})
}

func (this *AppController) GetApp() {
	if !this.check() {
		return
	}
	this.info = &map[string]string{
		"err" : "查询失败",
	}
	apps, err := this.app()
	if err != nil {
		this.error("err")
		return
	}
	this.success(apps)
}


func (this *AppController) app() ([]map[string]string, error) {
	appdb, err := beego.AppConfig.GetSection("appdb")
	if err != nil {
		return nil, err
	}
	db := models.NewDB(appdb["host"], appdb["user"], appdb["pass"])
	dbDo := models.NewDBDo(db, appdb["database"])
	apps, err := dbDo.ShowData(appdb["table"], 0)
	outApps := make([]map[string]string, 0)
	if err == nil {
		for _, app := range apps {
			outApps = append(outApps, map[string]string{
				"id" : app["id"].(string),
				"name" : app["name"].(string),
				"create_time" : app["create_time"].(string),
			})
		}
	}
	return outApps, nil
}
