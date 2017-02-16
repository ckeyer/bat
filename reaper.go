package bat

import (
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Reaper interface {
	Reap(*http.Request) error
}

type Record struct {
	Name       string `json:"name" bson:"name"`
	Path       string `json:"path" bson:"path"`
	RemoteAddr string `json:"remote_addr" bson:"remote_addr"`
	Referer    string `json:"referer" bson:"referer"`
	Agent      string `json:"agent" bson:"agent"`
	// for url's Params
	Labels map[string]interface{} `json:"labels" bson:"labels"`
	// for more headers
	Supplement map[string]interface{} `json:"supplement" bson:"supplement"`
	Cookies    []*http.Cookie         `json:"cookies" bson:"cookies"`

	RecordAt time.Time `json:"record_at" bson:"record_at"`
}

func NewRecord(name string, req *http.Request) *Record {
	r := &Record{
		Name:       name,
		Path:       req.URL.Path,
		RemoteAddr: req.RemoteAddr,
		Referer:    req.Referer(),
		Agent:      req.UserAgent(),
		Cookies:    req.Cookies(),

		Labels:     map[string]interface{}{},
		Supplement: map[string]interface{}{},

		RecordAt: time.Now(),
	}

	for k, v := range req.URL.Query() {
		log.Debugf("url params key: %s, value: %v", k, v)
		if len(v) == 1 {
			r.Labels[k] = v[0]
		} else {
			r.Labels[k] = v
		}
	}

	for k, v := range req.Header {
		switch strings.ToLower(k) {
		case "host",
			"user-agent",
			"cookie",
			"cache-control",
			"accept-encoding",
			"connection",
			"pragma",
			"accept",
			"dnt":
			// ignore headers
		default:
			if len(v) == 1 {
				r.Supplement[k] = v[0]
			} else {
				r.Supplement[k] = v
			}
		}
	}

	return r
}

func Reap(req *http.Request) error {
	log.Debugf("Path: %s", req.URL.Path)
	log.Debug("RemoteAddr:", req.RemoteAddr)
	log.Debug(req.Referer())
	log.Debug("cookies", req.Cookies())
	log.Debug(req.UserAgent())
	r := NewRecord("", req)
	log.Infof("%+v", r)
	return nil
}
