package decorator

import (
    log "github.com/Sirupsen/logrus"

	"errors"
    "net/http"
)

type Recover struct {
}

func NewRecover() *Recover {
    return &Recover{}
}

func (e Recover) Do(h http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rcv := recover(); rcv != nil {
                var err error
                switch t := rcv.(type) {
                    case string:
                        err = errors.New(t)
                    case error:
                        err = t
                    default:
                        err = errors.New("Unknown error")
                }
                log.Error(err.Error())
            }
        }()

        h(w, r)
    }
}