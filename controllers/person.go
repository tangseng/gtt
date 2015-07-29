package controllers

import (
	"gtt/models"
)

type PersonController struct {
	BaseController

	Name string
	LoginName string
	LoginPass string
	Color string
	Group string
}

func (this *PersonController) check() bool {
	if !this.IsAdmin() {
		this.Redirect("/day", 302)
		return false
	}
	return true
}

func (this *PersonController) Get() {
	if !this.check() {
		return
	}
	persons := models.NewPersonOption().ReadAll()
	this.Data["PagePerson"] = true
	this.Data["Group"] = models.GroupConfig
	this.Data["Persons"] = persons
	this.Data["Css"] = []string{"colorpicker"}
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/controllers/person/personCtrl", "colorpicker", "app/directives/color"}
	this.TplNames = "person.html"
}

func (this *PersonController) prefixCheck() bool {
	if !this.check() {
		return false
	}
	this.info = &map[string]string{
		"post" : "所有信息必须填写",
		"group" : "不存在该用户组",
		"user" : "用户已存在",
		"nouser" : "用户不存在",
		"login" : "登录名已存在",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	this.Name = this.GetString("name")
	this.LoginName = this.GetString("loginName")
	this.LoginPass = this.GetString("loginPass")
	this.Color = this.GetString("color")
	this.Group = this.GetString("group")
	if len(this.Name) == 0 || len(this.LoginName) == 0 || len(this.LoginPass) == 0 || len(this.Color) == 0 {
		this.error("post")
		return false
	}
	if _, ok := models.GroupConfig[this.Group]; !ok {
		this.error("group")
		return false
	}
	return true
}

func (this *PersonController) Add() {
	if !this.prefixCheck() {
		return
	}
	personOption := models.NewPersonOption()
	person, err := personOption.ReadByName(this.Name)
	if err == nil && person.Id > 0 {
		this.error("user")
		return
	}
	person, err = personOption.ReadByLoginName(this.LoginName)
	if err == nil && person.Id > 0 {
		this.error("login")
		return
	}
	id, err := personOption.Insert(models.NewPerson(this.Name, this.LoginName, this.LoginPass, this.Color, this.Group))
	if err != nil {
		this.error("err")
		return
	}
	this.success(map[string]interface{} {"id" : id})
}

func (this *PersonController) Update() {
	if !this.prefixCheck() {
		return
	}
	id, _ := this.GetInt("id")
	if id == 0 {
		this.error("nouser")
		return
	}
	personOption := models.NewPersonOption()
	person, err := personOption.ReadById(id)
	if err != nil {
		this.error("nouser")
		return
	}
	person.Name = this.Name
	person.LoginName = this.LoginName
	person.LoginPass = this.LoginPass
	person.Color = this.Color
	person.Group = this.Group
	_, err = personOption.Update(person)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}

func (this *PersonController) Delete() {
	if !this.check() {
		return
	}
	this.info = &map[string]string{
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	id, _ := this.GetInt("id")
	if id == 0 {
		this.error("nouser")
		return
	}
	err := models.NewPersonOption().Delete(id)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}
