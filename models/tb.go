package models

import (
	"github.com/astaxie/beego/orm"
	"time"
	"fmt"
)

type Tb struct {
	Id int `json:"id"`
	Uid int `json:"uid"`
	App string `orm:"null" json:"app"`
	Title string `json:"title"`
	Desc string `orm:"type(text)" json:"desc"`
	Custom string `json:"custom"`
	Date int64 `json:"date"`
	Why string `orm:"type(text)" json:"why"`
	Fix string `orm:"type(text)" json:"fix"`
	Time int64 `json:"time"`
}

func NewTb(uid int,app, title, desc, custom, why, fix string, date int64) *Tb {
	return &Tb{
		Uid : uid,
		App : app,
		Title : title,
		Desc : desc,
		Custom : custom,
		Date : date,
		Why : why,
		Fix : fix,
		Time : time.Now().Unix(),
	}
}

type TbOption struct {
	Orm orm.Ormer
}

func NewTbOption() *TbOption {
	return &TbOption{
		Orm : cacheOrm,
	}
}

func (this *TbOption) Insert(tb *Tb) (int64, error) {
	return this.Orm.Insert(tb)
}

func (this *TbOption) Update(tb *Tb) (int64, error) {
	return this.Orm.Update(tb)
}

func (this *TbOption) Read(id int) (*Tb, error) {
	tb := &Tb{Id : id}
	err := this.Orm.Read(tb)
	if err != nil {
		return nil, err
	}
	return tb, nil
}

func (this *TbOption) Search(key string, columns map[string]string, offset, limit int) []Tb {
	tbs := make([]Tb, 0)
	cond := orm.NewCondition()
	if len(columns) > 0 {
		condColumn := orm.NewCondition()
		for column, columnVal := range columns {
			condColumn = condColumn.And(fmt.Sprintf("%s__icontains", column), columnVal)
		}
		cond = cond.AndCond(condColumn)
	}
	if len(key) > 0 {
		cond = cond.AndCond(cond.Or("title__icontains", key).Or("desc__icontains", key).Or("why__icontains", key).Or("fix__icontains", key))
	}
	this.Orm.QueryTable("Tb").SetCond(cond).OrderBy("-id").Offset(offset).Limit(limit).All(&tbs)
	return tbs
}

func (this *TbOption) ReadMore(offset, limit int) []Tb {
	tbs := make([]Tb, 0)
	this.Orm.QueryTable("Tb").OrderBy("-id").Offset(offset).Limit(limit).All(&tbs)
	return tbs
}

func (this *TbOption) Delete(id int) error {
	_, err := this.Orm.QueryTable("tb").Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return nil
}
