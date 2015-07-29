package controllers

import (
	"gtt/models"
	"time"
)

type DayController struct {
	BaseController
	Today time.Time
	Y int
	M int
	D int
}

func (this *DayController) prefixInit() {
	this.Today = time.Now()
	y, m, d := this.Today.Date()
	this.Y = y
	this.M = int(m)
	this.D = d
}

func (this *DayController) Get() {
	this.prefixInit()
	dayOption := models.NewDayOption()
	day, _ := dayOption.ReadDay(this.user.Id, this.Y, this.M, this.D)
	this.Data["PageDay"] = true
	this.Data["Day"] = day
	this.Data["Y"] = this.Y
	this.Data["M"] = this.M
	this.Data["D"] = this.D
	js := []string{"app/app", "app/directives/tip", "app/directives/date", "app/services/time", "app/controllers/day/dayCtrl", "app/controllers/day/planCtrl"}
	if needApp, ok := this.Data["NeedApp"].(bool); ok && needApp {
		js = append(js, "app/controllers/day/appCtrl")
	}
	admin := this.IsAdmin()
	if admin {
		js = append(js, "app/controllers/day/dayAdminCtrl")
	}
	this.Data["Script"] = js
	this.Data["Admin"] = admin
	this.TplNames = "day.html"
}

func (this *DayController) Add() {
	this.info = &map[string]string{
		"post" : "所有信息必须填写",
		"day" : "添加今天的工作内容失败（day）",
		"work" : "添加今天的工作内容失败（work）",
		"plan" : "没有这个工作计划",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	this.prefixInit()
	content := this.GetString("content")
	startTime, _ := this.GetInt("startTime")
	endTime, _ := this.GetInt("endTime")
	status, _ := this.GetInt8("status")
	planId, _ := this.GetInt("planId")
	appId, _ := this.GetInt("appId")
	if len(content) == 0 || startTime == 0 || endTime == 0 || status == 0 {
		this.error("post")
		return
	}
	if planId != 0 {
		planOption := models.NewPlanOption()
		plan, err := planOption.Read(planId)
		if err != nil {
			this.error("plan")
			return
		}
		plan.Status = status
		if plan.Status == 100 {
			plan.RealTime = this.Today.Unix()
		}
		planOption.Update(plan)
	}
	dayOption := models.NewDayOption()
	day, err := dayOption.ReadDay(this.user.Id, this.Y, this.M, this.D)
	if err != nil {
		day = models.NewDay(this.user.Id, this.Y, this.M, this.D)
		_, err = dayOption.Insert(day)
		if err != nil {
			this.error("day")
			return
		}
	}
	workid, err := dayOption.InsertWork(models.NewDayWork(content, startTime, endTime, status, day, planId))
	if err != nil {
		this.error("work")
		return
	}
	if appId > 0 {
		models.NewAppOption().Insert(models.NewApp(this.user.Id, appId, content))
	}
	this.success(map[string]interface{} {"id" : workid})
}


func (this *DayController) prefixCheck() bool {
	this.info = &map[string]string{
		"auth" : "没有权限更新",
		"post" : "所有信息必须填写",
		"id" : "没有该工作记录",
		"day" : "更新时添加工作内容失败（day）",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	if !this.IsAdmin() {
		this.error("auth")
		return false
	}
	return true
}

func (this *DayController) Search() {
	if !this.prefixCheck() {
		return
	}
	id, _ := this.GetInt("id")
	dayWork, err := models.NewDayOption().ReadDayWork(id)
	if err != nil {
		this.error("id")
		return
	}
	this.success(dayWork)
}

func (this *DayController) Update() {
	if !this.prefixCheck() {
		return
	}
	content := this.GetString("content")
	startTime, _ := this.GetInt("startTime")
	endTime, _ := this.GetInt("endTime")
	status, _ := this.GetInt8("status")
	if len(content) == 0 || startTime == 0 || endTime == 0 || status == 0 {
		this.error("post")
		return
	}

	date, _ := this.GetInt64("date")
	timer := time.Unix(date, 0)
	y := timer.Year()
	m := int(timer.Month())
	d := timer.Day()

	dayOption := models.NewDayOption()
	id, _ := this.GetInt("id")
	dayWork, err := dayOption.ReadDayWork(id)
	if err != nil {
		this.error("id")
		return
	}
	planId := dayWork.PlanId
	if planId != 0 {
		planOption := models.NewPlanOption()
		plan, err := planOption.Read(planId)
		if err == nil && plan != nil {
			plan.Status = status
			if plan.Status == 100 {
				plan.RealTime = this.Today.Unix()
			}
			planOption.Update(plan)
		}
	}
	oldDay, _ := dayOption.ReadDayById(dayWork.Day.Id, false)
	day, err := dayOption.ReadDay(oldDay.Uid, y, m, d)
	if err != nil {
		day = models.NewDay(oldDay.Uid, y, m, d)
		_, err = dayOption.Insert(day)
		if err != nil {
			this.error("day")
			return
		}
	}
	dayWork.Content = content
	dayWork.StartTime = startTime
	dayWork.EndTime = endTime
	dayWork.Status = status
	dayWork.Day = day
	dayWork.Time = uint64(date)
	_, err = dayOption.UpdateDayWork(dayWork, "")
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}

func (this *DayController) Delete() {
	if !this.prefixCheck() {
		return
	}
	id, _ := this.GetInt("id")
	err := models.NewDayOption().DeleteDayWork(id)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}
