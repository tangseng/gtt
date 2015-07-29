package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Mjx struct {
	Id int `json:"id"`
	Uid int `json:"uid"`
	Content string `orm:"type(text)" json:"content"`
	Score int8 `json:"score"`
	Date int64 `json:"date"`
	Time int64 `json:"time"`
}

func NewMjx(uid int, content string, score int8, date int64) *Mjx {
	return &Mjx{
		Uid : uid,
		Content : content,
		Score : score,
		Date : date,
		Time : time.Now().Unix(),
	}
}

type MjxOption struct {
	Orm orm.Ormer
}

func NewMjxOption() *MjxOption {
	return &MjxOption{
		Orm : cacheOrm,
	}
}

func (this *MjxOption) Insert(Mjx *Mjx) (int64, error) {
	return this.Orm.Insert(Mjx)
}

func (this *MjxOption) Update(Mjx *Mjx) (int64, error) {
	return this.Orm.Update(Mjx)
}

func (this *MjxOption) Read(id int) (*Mjx, error) {
	Mjx := &Mjx{Id : id}
	err := this.Orm.Read(Mjx)
	if err != nil {
		return nil, err
	}
	return Mjx, nil
}

func (this *MjxOption) ReadByUid(uid int) []Mjx {
	Mjxs := make([]Mjx, 0)
	this.Orm.QueryTable("mjx").Filter("uid", uid).Filter("status__lt", 100).All(&Mjxs)
	return Mjxs
}

func (this *MjxOption) ReadByUidInMonth(uid int, startTime, endTime int64) []Mjx {
	Mjxs := make([]Mjx, 0)
	this.Orm.QueryTable("mjx").Filter("uid", uid).Filter("date__gt", startTime).Filter("date__lt", endTime).OrderBy("-id").All(&Mjxs)
	return Mjxs
}

func (this *MjxOption) Delete(id int) error {
	_, err := this.Orm.QueryTable("mjx").Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return nil
}
