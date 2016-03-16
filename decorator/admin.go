package decorator

import (
    "net/http"

    "github.com/gorilla/context"
)

type Admin struct {
	Callback func(interface{}) bool
}

func NewAdmin(callback func(interface{}) bool) *Admin {
    return &Admin{Callback: callback}
}

func (a Admin) Do(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        u := context.Get(r, "user")
        if u == nil {
            http.Redirect(w, r, "/login", http.StatusFound)
            return
        }

		if a.Callback(u) {
			h(w,r)
			return
		}

        http.Redirect(w, r, "/login", http.StatusFound)
    }
}
