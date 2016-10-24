package bat

import (
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
)

const (
	API = ""
)

var (
	DefaultImage = GetImage()
)

func Serve(addr string) {
	r := martini.NewRouter()
	mar := martini.New()
	mar.Use(martini.Recovery())
	mar.Use(render.Renderer())
	mar.MapTo(r, (*martini.Routes)(nil))
	mar.Action(r.Handle)
	m := &martini.ClassicMartini{Martini: mar, Router: r}

	m.Group(API, func(r martini.Router) {
		r.Group("/image", func(r martini.Router) {
			r.Get("/png", PNG)
			r.Get("/jpg", JPG)
			r.Get("/gif", GIF)
			r.Get("/:key", IMG)
		}, Reap)
	})
	log.Infof("listening on %s", addr)
	m.RunOnAddr(addr)
}

func Hello(rw http.ResponseWriter, req *http.Request) {
	rw.Write([]byte("hello"))
}

func IMG(rw http.ResponseWriter, req *http.Request) {
	for k, handler := range map[string]http.HandlerFunc{
		".jpg": JPG,
		".gif": GIF,
		".png": PNG,
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
