package main

import (
	_ "gtt/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
	"net/http"
	"html/template"
)

func main() {
	beego.Errorhandler("404", page404)
	beego.InsertFilter("*", beego.BeforeExec, (beego.FilterFunc)(filter))
	beego.Run()
}

func filter(ctx *context.Context) {
	url := ctx.Request.URL.Path
	if strings.HasPrefix(url, "/static/") || strings.HasPrefix(url, "/login/md5") {
		return
	}

	sess, _ := beego.GlobalSessions.SessionStart(ctx.ResponseWriter, ctx.Request)
	user := sess.Get("user")
	if user != nil {

	} else {
		if url != "/login/" && url != "/login" && url != "/login/login" {
			ctx.Redirect(302, "/login")
			return
		}
	}
}

func page404(rw http.ResponseWriter, r *http.Request){
	t,_:= template.New("404.html").ParseFiles(beego.ViewsPath + "/404.html")
	t.Execute(rw, nil)
}
