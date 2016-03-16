package gcore

import (
	log "github.com/Sirupsen/logrus"

	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Gcore struct {
	address    string
	decorators []httpHandlerDecorator
	router     *mux.Router
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

type httpHandlerDecorator interface {
	Do(http.HandlerFunc) http.HandlerFunc
}

func New(address string) *Gcore {

	router := mux.NewRouter().StrictSlash(true)

	return &Gcore{
		address: address,
		router:  router,
	}
}

func (g *Gcore) RegisterDecorator(f httpHandlerDecorator) {
	g.decorators = append(g.decorators, f)
}

func (g *Gcore) AddStatic(path string, uri string) {
	matcherFunc := func(r *http.Request, rm *mux.RouteMatch) bool {
		return strings.Contains(r.RequestURI, ".")
	}

	g.router.HandleFunc(uri + "{file}", serveStatics(path + uri)).MatcherFunc(matcherFunc)
}

func (g *Gcore) AddRoute(path string, f http.HandlerFunc, d ...httpHandlerDecorator) {
	var decorators []httpHandlerDecorator

	decorators = append(decorators, d...)
	decorators = append(decorators, g.decorators...)

	g.router.HandleFunc(path, useHandlers(f, decorators))
}

func (g *Gcore) Boot() {
	srv := &http.Server{
		Handler:        g.router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	ln, err := net.Listen("tcp", g.address)
	if err != nil {
		log.Fatal(err)
	}

	log.Info("Listening on... ", g.address)

	err = srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
	if err != nil {
		log.Fatal(err)
	}
}

func serveStatics(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		file := vars["file"]

		http.ServeFile(w, r, path + file)
	}
}

func useHandlers(h http.HandlerFunc, decorators []httpHandlerDecorator) http.HandlerFunc {
	n        := len(decorators)
	handlers := make(map[int]httpHandlerDecorator, n)

	var userKeys []int
	for k := range decorators {
		userKeys = append(userKeys, k)
	}

	sort.Ints(userKeys)

	for _, k := range userKeys {
		handlers[k] = decorators[k]
	}

	var keys []int
	for k := range handlers {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		h = handlers[k].Do(h)
	}

	return h
}