package controllers

import (
	"gtt/models"
	"time"
)

type PlanController struct {
	BaseController
	Today int64
	Content string
	StartTime int64
	EndTime int64
}

func (this *PlanController) prefixInit() {
	y, m, d := time.Now().Date()
	this.Today = time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
}

func (this *PlanController) Get() {
	plans := models.NewPlanOption().ReadByUid(this.user.Id)
	this.success(plans)
}

func (this *PlanController) GetComplete() {
	plans := models.NewPlanOption().ReadByUidOfComplete(this.user.Id, 10)
	this.success(plans)
}

func (this *PlanController) GetPlan() {
	planId, _ := this.GetInt("id")
	dayWorks := models.NewDayOption().ReadByPlanId(planId)
	this.success(dayWorks)
}

func (this *PlanController) checkForm(isadd bool) bool {
	this.info = &map[string]string{
		"id" : "不存在该计划",
		"auth" : "该计划你没有权限修改",
		"post" : "所有信息必须填写",
		"end" : "结束时间不能早于开始时间，并且不能早于今天",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	content := this.GetString("content")
	startTime, _ := this.GetInt64("startTime")
	endTime, _ := this.GetInt64("endTime")
	if len(content) == 0 || startTime == 0 || endTime == 0 {
		this.error("post")
		return false
	}
	if (isadd && startTime < this.Today) || endTime <= startTime {
		this.error("end")
		return false
	}
	this.Content = content
	this.StartTime = startTime
	this.EndTime = endTime
	return true
}

func (this *PlanController) Add() {
	this.prefixInit()
	if !this.checkForm(true) {
		return
	}
	planOption := models.NewPlanOption()
	id, err := planOption.Insert(models.NewPlan(this.user.Id, this.Content, this.StartTime, this.EndTime))
	if err != nil {
		this.error("err")
		return
	}
	this.success(map[string]interface{} {"id" : id})
}

func (this *PlanController) Update() {
	this.prefixInit()
	if !this.checkForm(false) {
		return
	}
	id, _ := this.GetInt("id")
	if id == 0 {
		this.error("id")
		return
	}
	planOption := models.NewPlanOption()
	plan, err := planOption.Read(id)
	if err != nil {
		this.error("id")
		return
	}
	if plan.Uid != this.user.Id {
		this.error("auth")
		return
	}
	plan.Content = this.Content
	plan.StartTime = this.StartTime
	plan.EndTime = this.EndTime
	_, err = planOption.Update(plan)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}

func (this *PlanController) Delete() {
	this.info = &map[string]string{
		"auth" : "该计划你没有权限删除",
		"id" : "不存在该计划",
		"err" : "操作失败",
		"ok" : "操作成功",
	}


	id, _ := this.GetInt("id")
	if id == 0 {
		this.error("id")
		return
	}
	planOption := models.NewPlanOption()
	plan, err := planOption.Read(id)
	if err != nil {
		this.error("id")
		return
	}
	if plan.Uid != this.user.Id {
		this.error("auth")
		return
	}
	err = planOption.Delete(id)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}
