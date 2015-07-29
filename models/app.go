package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type App struct {
	Id int `json:"id"`
	Uid int `json:"uid"`
	Appid int `orm:"null" json:"appid"`
	Desc string `orm:"type(text)" json:"desc"`
	Time int64 `json:"time"`
}

func NewApp(uid, appid int, desc string) *App {
	return &App{
		Uid : uid,
		Appid : appid,
		Desc : desc,
		Time : time.Now().Unix(),
	}
}

type AppOption struct {
	Orm orm.Ormer
}

func NewAppOption() *AppOption {
	return &AppOption{
		Orm : cacheOrm,
	}
}

func (this *AppOption) Insert(app *App) (int64, error) {
	return this.Orm.Insert(app)
}

func (this *AppOption) Update(app *App) (int64, error) {
	return this.Orm.Update(app)
}

func (this *AppOption) Read(id int) (*App, error) {
	app := &App{Id : id}
	err := this.Orm.Read(app)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (this *AppOption) ReadMore(appid, offset, limit int) []App {
	apps := make([]App, 0)
	this.Orm.QueryTable("app").Filter("appid", appid).OrderBy("-id").Offset(offset).Limit(limit).All(&apps)
	return apps
}

func (this *AppOption) Delete(id int) error {
	_, err := this.Orm.QueryTable("app").Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return nil
}
