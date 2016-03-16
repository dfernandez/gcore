package decorator

import (
    "net/http"

    "github.com/gorilla/context"
    "github.com/gorilla/sessions"
)

type Context struct {
	HashKey     []byte
	SessionName string
}

func NewContext(hashKey []byte, sessionName string) *Context {
    return &Context{
		HashKey: hashKey,
		SessionName: sessionName,
	}
}

func (c Context) Do(h http.HandlerFunc) http.HandlerFunc {
    store := sessions.NewCookieStore(c.HashKey)

    return func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, c.SessionName)
        context.Set(r, "user", session.Values["user"])

        h(w, r)
    }
}