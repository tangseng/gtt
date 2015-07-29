package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type Gg struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Desc string `orm:"type(text)" json:"desc"`
	Date int64 `json:"date"`
	Class string `json:"class"`
	Time int64 `json:"time"`
}

func NewGg(title, desc, class string, date int64) *Gg {
	return &Gg{
		Title : title,
		Desc : desc,
		Date : date,
		Class : class,
		Time : time.Now().Unix(),
	}
}

type GgOption struct {
	Orm orm.Ormer
}

func NewGgOption() *GgOption {
	return &GgOption{
		Orm : cacheOrm,
	}
}

func (this *GgOption) Insert(Gg *Gg) (int64, error) {
	return this.Orm.Insert(Gg)
}

func (this *GgOption) Update(Gg *Gg) (int64, error) {
	return this.Orm.Update(Gg)
}

func (this *GgOption) Read(id int) (*Gg, error) {
	Gg := &Gg{Id : id}
	err := this.Orm.Read(Gg)
	if err != nil {
		return nil, err
	}
	return Gg, nil
}

func (this *GgOption) ReadMore(offset, limit int) []Gg {
	Ggs := make([]Gg, 0)
	this.Orm.QueryTable("Gg").OrderBy("-id").Offset(offset).Limit(limit).All(&Ggs)
	return Ggs
}

func (this *GgOption) Delete(id int) error {
	_, err := this.Orm.QueryTable("Gg").Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return nil
}
