package controllers

import (
	"gtt/models"
)

type GGController struct {
	BaseController

	Title string
	Desc string
	Class string
	Date int64
}

func (this *GGController) Get() {
	this.Data["PageNum"] = PAGENUM
	this.Data["PageGg"] = true
	if this.IsAdmin() {
		this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/directives/date", "app/services/time", "app/controllers/gg/ggAdminCtrl"}
		this.TplNames = "gg.admin.html"
	} else {
		this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/services/time", "app/controllers/gg/ggCtrl"}
		this.TplNames = "gg.html"
	}
}

func (this *GGController) Ajax() {
	offset, _ := this.GetInt("offset")
	gg := models.NewGgOption().ReadMore(offset, PAGENUM)
	this.success(gg)
}

func (this *GGController) prefixCheck() bool {
	this.info = &map[string]string{
		"auth" : "没有权限发布公告",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	if !this.IsAdmin() {
		this.error("auth")
		return false
	}
	this.Title = this.GetString("title")
	this.Desc = this.GetString("desc")
	this.Class = this.GetString("class")
	this.Date, _ = this.GetInt64("date")
	return true
}

func (this *GGController) Add() {
	if !this.prefixCheck() {
		return
	}
	id, err := models.NewGgOption().Insert(models.NewGg(this.Title, this.Desc, this.Class, this.Date))
	if err != nil {
		this.error("err")
		return
	}
	this.success(map[string]interface{} {"id" : id})
}

func (this *GGController) Update() {
	if !this.prefixCheck() {
		return
	}
	ggOption := models.NewGgOption()
	id, _ := this.GetInt("id")
	gg, err := ggOption.Read(id)
	if err != nil {
		this.error("err")
		return
	}
	gg.Title = this.Title
	gg.Desc = this.Desc
	gg.Class = this.Class
	gg.Date = this.Date
	_, err = ggOption.Update(gg)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}

func (this *GGController) Delete() {
	if !this.prefixCheck() {
		return
	}
	id, _ := this.GetInt("id")
	err := models.NewGgOption().Delete(id)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}
