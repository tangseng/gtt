package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

type Person struct {
	Id int `json:"id"`
	Name string `json:"name"`
	LoginName string `json:"loginName"`
	LoginPass string `json:"loginPass"`
	Color string `json:"color"`
	Group string `json:"group"`
	Time uint64 `json:"time"`
}

func NewPerson(name, loginName, loginPass, color, group string) *Person {
	return &Person{
		Name : name,
		LoginName : loginName,
		LoginPass : loginPass,
		Color : color,
		Group : group,
		Time : uint64(time.Now().Unix()),
	}
}

type PersonOption struct {
	Orm orm.Ormer
}

func NewPersonOption() *PersonOption {
	return &PersonOption{
		Orm : cacheOrm,
	}
}

func (this *PersonOption) Insert(person *Person) (int64, error) {
	return this.Orm.Insert(person)
}

func (this *PersonOption) Update(person *Person) (int64, error) {
	return this.Orm.Update(person)
}

func (this *PersonOption) Read(loginName, loginPass string) (*Person, error) {
	person := new(Person)
	err := this.Orm.QueryTable("person").Filter("loginName", loginName).Filter("loginPass", loginPass).One(person)
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (this *PersonOption) ReadById(id int) (*Person, error) {
	person := new(Person)
	err := this.Orm.QueryTable("person").Filter("id", id).One(person)
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (this *PersonOption) ReadByIds(ids []int) []Person {
	persons := make([]Person, 0)
	this.Orm.QueryTable("person").Filter("id__in", ids).All(&persons)
	return persons
}

func (this *PersonOption) ReadByName(name string) (*Person, error) {
	person := new(Person)
	err := this.Orm.QueryTable("person").Filter("name", name).One(person)
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (this *PersonOption) ReadByLoginName(loginName string) (*Person, error) {
	person := new(Person)
	err := this.Orm.QueryTable("person").Filter("loginName", loginName).One(person)
	if err != nil {
		return nil, err
	}
	return person, nil
}

func (this *PersonOption) ReadByGroup(group string) []Person {
	persons := make([]Person, 0)
	this.Orm.QueryTable("person").Filter("group", group).All(&persons)
	return persons
}

func (this *PersonOption) ReadAll() []Person {
	persons := make([]Person, 0)
	this.Orm.QueryTable("person").All(&persons)
	return persons
}

func (this *PersonOption) Delete(id int) error {
	_, err := this.Orm.QueryTable("person").Filter("id", id).Delete()
	if err != nil {
		return err
	}
	return nil
}
