package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Plan struct {
	Id int `json:"id"`
	Uid int `json:"uid"`
	Content string `orm:"type(text)" json:"content"`
	StartTime int64 `json:"startTime"`
	EndTime int64 `json:"endTime"`
	RealTime int64 `json:"realTime"`
	Status int8 `json:"status"`
	Time int64 `json:"time"`
}

func NewPlan(uid int, content string, startTime, endTime int64) *Plan {
	return &Plan{
		Uid : uid,
		Content : content,
		StartTime : startTime,
		EndTime : endTime,
		Status : 0,
		Time : time.Now().Unix(),
	}
}

type PlanOption struct {
	Orm orm.Ormer
}

func NewPlanOption() *PlanOption {
	return &PlanOption{
		Orm : cacheOrm,
	}
}

func (this *PlanOption) Insert(plan *Plan) (int64, error) {
	return this.Orm.Insert(plan)
}

func (this *PlanOption) Update(plan *Plan) (int64, error) {
	return this.Orm.Update(plan)
}

func (this *PlanOption) Read(id int) (*Plan, error) {
	plan := &Plan{Id : id}
	err := this.Orm.Read(plan)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

func (this *PlanOption) ReadByUid(uid int) []Plan {
	plans := make([]Plan, 0)
	this.Orm.QueryTable("plan").Filter("uid", uid).Filter("status__lt", 100).All(&plans)
	return plans
}

func (this *PlanOption) ReadByUidOfComplete(uid int, number int) []Plan {
	plans := make([]Plan, 0)
	this.Orm.QueryTable("plan").Filter("uid", uid).Filter("status", 100).OrderBy("-id").Limit(number).All(&plans)
	return plans
}

func (this *PlanOption) ReadByUidsInMonth(uids []int, startTime, endTime int64) []Plan {
	plans := make([]Plan, 0)
	cond := orm.NewCondition()
	condF := cond.And("uid__in", uids)
	condS := cond.And("startTime__gt", startTime).And("startTime__lt", endTime)
	condE := cond.And("endTime__gt", startTime).And("endTime__lt", endTime)
	condR := cond.And("realTime__gt", startTime).And("realTime__lt", endTime)
	this.Orm.QueryTable("plan").SetCond(cond.AndCond(condF).AndCond(cond.AndCond(condS).OrCond(condE).OrCond(condR))).All(&plans)
	return plans
}

func (this *PlanOption) ReadInMonth(startTime, endTime int64) []Plan {
	plans := make([]Plan, 0)
	cond := orm.NewCondition()
	condS := cond.And("startTime__gt", startTime).And("startTime__lt", endTime)
	condE := cond.And("endTime__gt", startTime).And("endTime__lt", endTime)
	condR := cond.And("realTime__gt", startTime).And("realTime__lt", endTime)
	this.Orm.QueryTable("plan").SetCond(cond.AndCond(condS).OrCond(condE).OrCond(condR)).All(&plans)
	return plans
}

func (this *PlanOption) Delete(id int) error {
	_, err := this.Orm.QueryTable("plan").Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return nil
}
