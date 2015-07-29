package controllers

import (
	"github.com/astaxie/beego"
	"os/exec"
	"path/filepath"
	"gtt/models"
)

const (
	PAGENUM = 30
)

func RootPath() string {
	file, _ := exec.LookPath("./conf/app.conf")
	path, _ := filepath.Abs(file)
	return path
}

/****** 基类 ******/

type BaseController struct {
	beego.Controller
	info *map[string]string
	user models.Person
}

func (this *BaseController) success(suc interface {}) {
	if index, yes := suc.(string); yes {
		this.Data["json"] = map[string]string{"success" : (*this.info)[index]}
	} else {
		this.Data["json"] = suc
	}
	this.ServeJson()
}

func (this *BaseController) error(err string) {
	errStr := err
	if tmpErrStr, ok := (*this.info)[err]; ok {
		errStr = tmpErrStr
	}
	this.Data["json"] = map[string]string{"error" : errStr}
	this.ServeJson()
}

func (this *BaseController) Prepare() {
	sess := this.StartSession()
	userInfo := sess.Get("user")
	if userInfo != nil {
		this.user = userInfo.(models.Person)
		this.user.LoginPass = "******"
		this.Data["Can"] = this.IsAdmin()
		this.Data["User"] = this.user
	}
	_, err := beego.AppConfig.GetSection("appdb")
	if err == nil {
		this.Data["NeedApp"] = true
	}
}

func (this *BaseController) IsAdmin() bool {
	return this.user.Id == -1 || this.user.Group == "zj"
}


