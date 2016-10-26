package bat

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ckeyer/bat/store"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

const (
	API_PREFIX = "/api"
	IMG_PREFIX = "/static/image"
)

var (
	DefaultImage = GetImage()
)

func Serve(addr string, db store.Store) {
	r := martini.NewRouter()
	mar := martini.New()
	mar.Use(martini.Recovery())
	mar.Use(render.Renderer())
	mar.MapTo(r, (*martini.Routes)(nil))
	mar.Action(r.Handle)
	m := &martini.ClassicMartini{Martini: mar, Router: r}

	if db == nil {
		log.Fatal("db is nil")
	}
	m.Map(db)

	m.Group(IMG_PREFIX, func(r martini.Router) {
		r.Get("/:name", IMG)
	}, Reap)

	m.Group(API_PREFIX, func(r martini.Router) {
		r.Get("/history/:name", ListHistory)
		// r.Post()
	})
	log.Infof("listening on %s", addr)
	m.RunOnAddr(addr)
}

func Hello(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("hello"))
}

func IMG(rw http.ResponseWriter, req *http.Request) {
	for k, handler := range map[string]http.HandlerFunc{
		"jpg": JPG,
		"gif": GIF,
		"png": PNG,
	} {
		if strings.HasSuffix(req.URL.Path, k) {
			handler(rw, req)
			return
		}
	}

	PNG(rw, req)
}

func PNG(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "image/png")
	png.Encode(rw, DefaultImage)
}

func JPG(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "image/jpeg")
	jpeg.Encode(rw, DefaultImage, &jpeg.Options{1})
}

func GIF(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "image/gif")
	gif.Encode(rw, DefaultImage, &gif.Options{})
}

func Reap(req *http.Request, params martini.Params, db store.Store) {
	name := params["name"]
	for _, v := range []string{
		"",
		"png",
		"jpg",
		"gif",
		"public.jpg",
		"public.png",
	} {
		// ignore
		if name == v {
			name = "public"
			break
		}
	}

	r := NewRecord(name, req)
	if err := db.LogHistory(r); err != nil {
		log.Errorf("log history failed. error: %s", err)
	}

	log.Debugf("%+v", r)
}

func ListHistory(wr http.ResponseWriter, req *http.Request, params martini.Params, db store.Store) {
	name := params["name"]

	data, err := db.ListHistory(name)
	if err != nil {
		return
	}
	// log.Debugf("wr retrun: %+v", ret)
	// rnd.JSON(200, nil)
	wr.WriteHeader(200)
	wr.Header().Set("Content-Type", "application/json")
	wr.Header().Add("Content-Type", "charset=UTF-8")
	wr.Write(data)
}
