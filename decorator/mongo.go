package decorator

import (
    "gopkg.in/mgo.v2"
	"net/http"

	"github.com/gorilla/context"
)

type Mongo struct {
	MongoServer string
}

func NewMongo(mongoServer string) *Mongo {
    return &Mongo{mongoServer}
}

func (m Mongo) Do(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        session, err := mgo.Dial(m.MongoServer)
        if err != nil {
            panic(err)
        }

        session.SetMode(mgo.Monotonic, true)
        context.Set(r, "mongoDB", session)

        h(w, r)

        session.Close()
    }
}
