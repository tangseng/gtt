package controllers

import (
	"gtt/models"
	"strconv"
)

type TBController struct {
	BaseController

	App string
	Title string
	Desc string
	Why string
	Fix string
	Custom string
	Date int64
}

func (this *TBController) check() bool {
	return true
}

func (this *TBController) Get() {
	if !this.check() {
		return
	}
	this.Data["PageNum"] = PAGENUM
	this.Data["PageTb"] = true
	this.Data["Css"] = []string{"colorpicker"}
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/directives/date", "app/services/time", "app/controllers/tb/tbCtrl", "colorpicker", "app/directives/color"}
	this.TplNames = "tb.html"
}

func (this *TBController) Search() {
	if !this.check() {
		return
	}
	app := this.GetString("app")
	custom := this.GetString("custom")
	key := this.GetString("key")
	offset, _ := this.GetInt("offset")
	columns := make(map[string]string)
	if len(app) > 0 {
		columns["app"] = app
	}
	if len(custom) > 0 {
		columns["custom"] = custom
	}
	tbs, pps := this.get(key, columns, offset)
	this.success(map[string]interface {}{"tbs" : tbs, "pps": pps})
}

func (this *TBController) get(key string, columns map[string]string, offset int) ([]models.Tb, map[string]string) {
	tbs := make([]models.Tb, 0)
	if len(key) > 0 || len(columns) > 0 {
		tbs = models.NewTbOption().Search(key, columns, offset, PAGENUM)
	} else {
		tbs = models.NewTbOption().ReadMore(offset, PAGENUM)
	}
	pids := make([]int, 0)
	pidsCheck := make(map[string]bool)
	for _, tb := range tbs {
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
	outPersons["-1"] = "管理员"
	return tbs, outPersons
}

func (this *TBController) prefixCheck() bool {
	if !this.check() {
		return false
	}
	this.info = &map[string]string{
		"post" : "标题与解决方案必须填写",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	this.App = this.GetString("app")
	this.Title = this.GetString("title")
	this.Desc = this.GetString("desc")
	this.Why = this.GetString("why")
	this.Fix = this.GetString("fix")
	this.Custom = this.GetString("custom")
	this.Date, _ = this.GetInt64("date")
	if len(this.Title) == 0 || len(this.Fix) == 0 {
		this.error("post")
		return false
	}
	return true
}

func (this *TBController) Add() {
	if !this.prefixCheck() {
		return
	}
	id, err := models.NewTbOption().Insert(models.NewTb(this.user.Id, this.App, this.Title, this.Desc, this.Custom, this.Why, this.Fix, this.Date))
	if err != nil {
		this.error("err")
		return
	}
	this.success(map[string]interface{} {"id" : id})
}

func (this *TBController) Update() {
	id, _ := this.GetInt("id")
	if id == 0 {
		this.Add()
		return
	}
	tbOption := models.NewTbOption()
	tb, err := tbOption.Read(id)
	if err != nil {
		this.Add()
		return
	}
	if !this.prefixCheck() {
		return
	}
	tb.App = this.App
	tb.Title = this.Title
	tb.Desc = this.Desc
	tb.Why = this.Why
	tb.Fix = this.Fix
	tb.Custom = this.Custom
	tb.Date = this.Date
	_, err = tbOption.Update(tb)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}

func (this *TBController) Delete() {
	if !this.check() {
		return
	}
	this.info = &map[string]string{
		"id" : "不存在该解决方案",
		"auth" : "你没有权限删除该问题解决方案",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	id, _ := this.GetInt("id")
	if id == 0 {
		this.error("id")
		return
	}
	tb, err := models.NewTbOption().Read(id)
	if err != nil {
		this.error("id")
		return
	}
	if !this.IsAdmin() && tb.Uid != this.user.Id {
		this.error("auth")
		return
	}
	this.success("ok")
}
