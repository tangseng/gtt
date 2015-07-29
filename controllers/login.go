package controllers

import (
	"gtt/models"
	"github.com/astaxie/beego"
	"errors"
)

type LoginController struct {
	BaseController
}

func (this *LoginController) check() bool {
	if this.user.Id != 0 {
		this.Redirect("/day", 302)
		return false
	}
	return true
}

func (this *LoginController) Login() {
	if !this.check() {
		return
	}
	this.Data["PageLogin"] = true
	this.Data["BM"] = beego.AppConfig.String("bm")
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/controllers/login/loginCtrl"}
	this.TplNames = "login.html"
}

func (this *LoginController) DoLogin() {
	if !this.check() {
		return
	}
	this.info = &map[string]string{
		"success" : "登录成功",
		"login" : "登录失败",
	}
	name := this.GetString("name")
	pass := this.GetString("pass")
	person, err := CheckAdmin(name, pass)
	if err != nil {
		person, err = models.NewPersonOption().Read(name, pass)
		if err != nil {
			this.error("login")
			return
		}
	}
	sess := this.StartSession()
	err = sess.Set("user", *person)
	if err != nil {
		this.error("login")
		return
	}
	this.success(map[string]int{"id" : person.Id})
}

func (this *LoginController) DoLoginOut() {
	sess := this.StartSession()
	sess.Set("user", nil)
	this.Redirect("/login", 302)
}

func (this *LoginController) Md5() {
	val := this.GetString("val")
	this.Ctx.Output.Body([]byte(models.MD5(val, beego.AppConfig.String("bm"))))
}

func CheckAdmin(name, pass string) (*models.Person, error) {
	bm := beego.AppConfig.String("bm")
	admin := beego.AppConfig.String("admin")
	adminPass := beego.AppConfig.String("pass")
	if name == admin && models.MD5(pass, bm) == adminPass {
		person := models.NewPerson("管理员", "admin", "***", "", "")
		person.Id = -1
		return person, nil
	}
	return nil, errors.New("不是管理员")
}
