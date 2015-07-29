package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
	"crypto/md5"
	"fmt"
	"math/rand"
	"io"
	"time"
	"os"
	"io/ioutil"
	"encoding/json"
)

var cacheOrm orm.Ormer

func init() {
	orm.Debug = false
	orm.RegisterDataBase("default", "sqlite3", "data.db")
	orm.RegisterModel(new(Person))
	orm.RegisterModel(new(Day))
	orm.RegisterModel(new(DayWork))
	orm.RegisterModel(new(Plan))
	orm.RegisterModel(new(Mjx))
	orm.RegisterModel(new(Tb))
	orm.RegisterModel(new(Gg))
	orm.RegisterModel(new(App))
	orm.RunSyncdb("default", false, true)
	cacheOrm = orm.NewOrm()
}

func MD5(str, salt string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum([]byte(salt)))
}

func RandStr(length int, replaceWords string) string {
	rand.Seed(time.Now().Unix())
	wordStr := "abcdefghijklmnopqrstuvwxyz"
	if len(replaceWords) != 0 {
		wordStr = replaceWords
	}
	words := []byte(wordStr)
	wLen := len(words)
	rs := make([]byte, 0)
	for i := 0; i < length; i++ {
		w := words[rand.Intn(wLen)]
		rs = append(rs, w)
	}
	return string(rs)
}

func RandStr2(length int) string {
	return RandStr(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

var GroupConfig map[string]string
func init() {
	file, _ := os.Open("./conf/group.conf")
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	json.Unmarshal(content, &GroupConfig)
}
