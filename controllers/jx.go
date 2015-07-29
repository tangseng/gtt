package controllers

import (
	"gtt/models"
	"time"
	"strconv"
)

type JXController struct {
	BaseController
	Today int64
	Y int
	M int
	D int
}


func (this *JXController) prefixInit() {
	y, m, d := time.Now().Date()
	this.Today = time.Date(y, m, d, 0, 0, 0, 0, time.UTC).Unix()
	this.Y = y
	this.M = int(m)
	this.D = d
}

func (this *JXController) Get() {
	this.prefixInit()
	cy, cm, cd := this.Y, this.M, this.D
	pm, _ := this.GetInt("m")
	if pm > 0 {
		if pm > cm {
			cy -= 1
		}
		cm = pm
		cd = 1
	}
	this.Data["PageJx"] = true
	this.Data["Y"] = cy
	this.Data["M"] = cm
	this.Data["D"] = cd
	admin := this.IsAdmin()
	if admin {
		persons := models.NewPersonOption().ReadAll()
		outPersons := map[string]string{}
		for _, person := range persons {
			outPersons[strconv.Itoa(person.Id)] = person.Name
		}
		this.Data["Persons"] = outPersons
		this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/directives/date", "app/services/month", "app/services/time", "app/controllers/jx/jxAdminCtrl"}
		this.TplNames = "jx.admin.html"
	} else {
		this.Data["Script"] = []string{"app/app", "app/directives/tip", "app/services/month", "app/services/time", "app/controllers/jx/jxCtrl"}
		this.TplNames = "jx.html"
	}
}

func (this *JXController) GetDay() {
	this.info = &map[string]string{
		"post" : "所有信息必须填写",
		"end" : "结束时间不能早于开始时间，并且不能早于今天",
	}
	y, _ := this.GetInt("y")
	m, _ := this.GetInt("m")
	d, _ := this.GetInt("d")
	uid := this.user.Id
	admin := this.IsAdmin()
	if admin {
		uid, _ = this.GetInt("uid")
	}
	day, err := models.NewDayOption().ReadDay(uid, y, m, d)
	works := make([]*models.DayWork, 0)
	if err == nil && len(day.Works) > 0 {
		works = day.Works
	}
	this.success(works)
}

func (this *JXController) GetMjx() {
	this.info = &map[string]string{
		"post" : "所有信息必须填写",
		"end" : "结束时间不能早于开始时间，并且不能早于今天",
	}
	uid := this.user.Id
	admin := this.IsAdmin()
	if admin {
		uid, _ = this.GetInt("uid")
	}
	startTime, _ := this.GetInt64("qian")
	endTime, _ := this.GetInt64("hou")
	if endTime < startTime {
		this.error("end")
		return
	}
	mjx := models.NewMjxOption().ReadByUidInMonth(uid, startTime, endTime)
	this.success(mjx)
}

func (this *JXController) AddMjx() {
	this.info = &map[string]string{
		"admin" : "没有权限添加月常规绩效",
	}
	if !this.IsAdmin() {
		this.error("admin")
		return
	}
	content := this.GetString("content")
	date, _ := this.GetInt64("date")
	score, _ := this.GetInt8("score")
	uid, _ := this.GetInt("uid")
	id, _ := models.NewMjxOption().Insert(models.NewMjx(uid, content, score, date))
	this.success(map[string]interface{}{"id" : id})
}

func (this *JXController) UpdateMjx() {
	this.info = &map[string]string{
		"admin" : "没有权限修改月常规绩效",
		"no" : "不存在该月常规绩效",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	if !this.IsAdmin() {
		this.error("admin")
		return
	}
	content := this.GetString("content")
	date, _ := this.GetInt64("date")
	score, _ := this.GetInt8("score")
	uid, _ := this.GetInt("uid")
	id, _ := this.GetInt("id")
	mjxOption := models.NewMjxOption()
	mjx, err := mjxOption.Read(id)
	if err != nil {
		this.error("no")
		return
	}
	mjx.Content = content
	mjx.Date = date
	mjx.Uid = uid
	mjx.Score = score
	mjx.Time = time.Now().Unix()
	_, err = mjxOption.Update(mjx)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}

func (this *JXController) DeleteMjx() {
	this.info = &map[string]string{
		"admin" : "没有权限删除月常规绩效",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	if !this.IsAdmin() {
		this.error("admin")
		return
	}
	id, _ := this.GetInt("id")
	err := models.NewMjxOption().Delete(id)
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}

func (this *JXController) Score() {
	this.info = &map[string]string{
		"admin" : "没有权限设置绩效分数",
		"err" : "操作失败",
		"ok" : "操作成功",
	}
	if !this.IsAdmin() {
		this.error("admin")
		return
	}
	id, _ := this.GetInt("id")
	score, _ := this.GetInt8("score")
	dayOption := models.NewDayOption()
	dayWork, err := dayOption.ReadDayWork(id)
	if err != nil {
		this.error("err")
		return
	}
	dayWork.Score = score
	_, err = dayOption.UpdateDayWork(dayWork, "score")
	if err != nil {
		this.error("err")
		return
	}
	this.success("ok")
}


func (this *JXController) XX() {
	this.prefixInit()
	this.Data["PageJx"] = true
	persons, _ := this.allPersons(false)
	this.Data["Persons"] = persons
	this.Data["Y"] = this.Y
	this.Data["M"] = this.M
	this.Data["D"] = this.D
	this.Data["Script"] = []string{"app/app", "app/directives/date", "app/services/time", "app/controllers/jx/jxaCtrl"}
	this.TplNames = "jxa.html"
}

func (this *JXController) XXs() {
	this.prefixInit()
	cy, cm, cd := this.Y, this.M, this.D
	if y, _ := this.GetInt("y"); y > 0 {
		cy = y
	}
	if m, _ := this.GetInt("m"); m > 0 {
		cm = m
	}
	if d, _ := this.GetInt("d"); d > 0 {
		cd = d
	}
	_, ids := this.allPersons(true)
	days := models.NewDayOption().ReadDayByUids(ids, cy, cm, cd)
	this.success(days)
}

func (this *JXController) allPersons(only bool) (map[string]map[string]string, []int) {
	ids := make([]int, 0)
	persons := models.NewPersonOption().ReadAll()
	if only {
		for _, person := range persons {
			if person.Group != "zj" {
				ids = append(ids, person.Id)
			}
		}
		return nil, ids
	} else {
		outPersons := map[string]map[string]string{}
		for _, person := range persons {
			if person.Group != "zj" {
				outPersons[strconv.Itoa(person.Id)] = map[string]string{
					"name" : person.Name,
					"color" : person.Color,
				}
				ids = append(ids, person.Id)
			}
		}
		return outPersons, ids
	}
}
