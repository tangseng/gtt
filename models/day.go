package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Day struct {
	Id int `json:"id"`
	Uid int `json:"uid"`
	Y int `json:"y"`
	M int `json:"m"`
	D int `json:"d"`

	Works []*DayWork `orm:"reverse(many)" json:"works"`
}

type DayWork struct {
	Id int `json:"id"`
	Content string `orm:"type(text)" json:"content"`
	StartTime int `json:"startTime"`
	EndTime int `json:"endTime"`
	Status int8 `json:"status"`
	Time uint64 `json:"time"`
	Score int8 `json:"score"`

	PlanId int `json:"planId"`

	Day *Day `orm:"rel(fk)"`
}

func NewDay(uid, y, m, d int) *Day {
	return &Day{
		Uid : uid,
		Y : y,
		M : m,
		D : d,
	}
}

func NewDayWork(content string, startTime, endTime int, status int8, day *Day, planId int) *DayWork {
	return &DayWork{
		Content : content,
		StartTime : startTime,
		EndTime : endTime,
		Status : status,
		Time : uint64(time.Now().Unix()),

		PlanId : planId,

		Day : day,
	}
}

type DayOption struct {
	Orm orm.Ormer
}

func NewDayOption() *DayOption {
	return &DayOption{
		Orm : cacheOrm,
	}
}

func (this *DayOption) Insert(day *Day) (int64, error) {
	return this.Orm.Insert(day)
}

func (this *DayOption) Update(day *Day) (int64, error) {
	return this.Orm.Update(day)
}

func (this *DayOption) ReadDay(uid, y, m, d int) (*Day, error) {
	day := new(Day)
	err := this.Orm.QueryTable("day").Filter("uid", uid).Filter("y", y).Filter("m", m).Filter("d", d).One(day)
	if err != nil {
		return nil, err
	}
	this.Orm.LoadRelated(day, "works")
	return day, nil
}

func (this *DayOption) ReadDayByUids(uids []int, y, m, d int) []Day {
	days := make([]Day, 0)
	this.Orm.QueryTable("day").Filter("uid__in", uids).Filter("y", y).Filter("m", m).Filter("d", d).All(&days)
	if len(days) > 0 {
		for index, day := range days {
			dayP := &day
			this.Orm.LoadRelated(dayP, "works")
			days[index] = *dayP
		}
	}
	return days
}

func (this *DayOption) ReadDayById(id int, need bool) (*Day, error) {
	day := new(Day)
	err := this.Orm.QueryTable("day").Filter("id", id).One(day)
	if err != nil {
		return nil, err
	}
	if need {
		this.Orm.LoadRelated(day, "works")
	}
	return day, nil
}

func (this *DayOption) ReadByPlanId(planId int) []DayWork {
	dayWorks := make([]DayWork, 0)
	this.Orm.QueryTable("day_work").Filter("plan_id", planId).All(&dayWorks)
	return dayWorks
}

func (this *DayOption) ReadDayWork(id int) (*DayWork, error) {
	dayWork := &DayWork{Id : id}
	err := this.Orm.Read(dayWork)
	if err != nil {
		return nil, err
	}
	return dayWork, nil
}

func (this *DayOption) UpdateDayWork(dayWork *DayWork, column string) (int64, error) {
	if column == "" {
		return this.Orm.Update(dayWork)
	}
	return this.Orm.Update(dayWork, column)
}

func (this *DayOption) DeleteDayWork(id int) error {
	_, err := this.Orm.QueryTable("day_work").Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (this *DayOption) Delete(id int) error {
	_, err := this.Orm.QueryTable("day").Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return nil
}

func (this *DayOption) InsertWork(work *DayWork) (int64, error) {
	return this.Orm.Insert(work)
}
