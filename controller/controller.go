package controller

import (
	log "github.com/Sirupsen/logrus"

    "net/http"
    "html/template"

	"github.com/gorilla/context"
)

type Controller struct {
    Template   string
    Layout     string
    Controller string
    Title      string
    Profile    interface{}
    TplVars    interface{}
}

func (tpl Controller) Render(w http.ResponseWriter, r *http.Request, tplVars interface{}) {

    funcMap := template.FuncMap{
        "add": func(x int, y int) int {
            return x + y
        },
    }

    t := template.Must(template.New("").Funcs(funcMap).ParseFiles("src/templates/" + tpl.Layout, "src/templates/" + tpl.Template))

    tpl.Title      = "Go web!"
    tpl.TplVars    = tplVars
    tpl.Controller = r.URL.Path

    if profile := context.Get(r, "user"); profile != nil {
        tpl.Profile = profile
    }

    err := t.ExecuteTemplate(w, "layout", tpl)
    if err != nil {
        log.Error(err)
    }
}