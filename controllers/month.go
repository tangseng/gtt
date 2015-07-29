package controllers

import (
	"gtt/models"
	"time"
)

type MonthController struct {
	BaseController
}

func (this *MonthController) Get() {
	y, m, d := time.Now().Date()
	persons := make([]models.Person, 0)
	admin := this.IsAdmin()
	if admin {
		persons = models.NewPersonOption().ReadAll()
	} else {
		persons = models.NewPersonOption().ReadByGroup(this.user.Group)
	}
	outputPersons := make([]models.Person, len(persons))
	for index, person := range persons {
		person.LoginPass = "******"
		outputPersons[index] = person
	}
	cy, cm, cd := y, int(m), d
	pm, _ := this.GetInt("m")
	if pm > 0 {
		if pm > cm {
			cy -= 1
		}
		cm = pm
		cd = 1
	}
	this.Data["Y"] = cy
	this.Data["M"] = cm
	this.Data["D"] = cd
	this.Data["PageMonth"] = true
	this.Data["Group"] = models.GroupConfig
	this.Data["Persons"] = outputPersons
	this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/directives/plan", "app/services/month", "app/services/time", "app/controllers/month/monthCtrl"}
	this.TplNames = "month.html"
}

func (this *MonthController) Plan() {
	this.info = &map[string]string{
		"diff" : "所请求的时间大于一个月的时间",
	}
	startTime, _ := this.GetInt64("start")
	endTime, _ := this.GetInt64("end")
	if endTime - startTime > 32 * 24 * 3600 {
		this.error("diff")
		return
	}
	isAdmin := this.IsAdmin()
	planOption := models.NewPlanOption()
	var plans []models.Plan
	if isAdmin {
		plans = planOption.ReadInMonth(startTime, endTime)
	} else {
		persons := models.NewPersonOption().ReadByGroup(this.user.Group)
		uids := make([]int, 0)
		for _, person := range persons {
			uids = append(uids, person.Id)
		}
		plans = planOption.ReadByUidsInMonth(uids, startTime, endTime)
	}
	this.success(plans)
}
