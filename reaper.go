package bat

import (
	"fmt"
	"net/http"
)

func Reap(req *http.Request) {
	fmt.Println(req.RemoteAddr)
	fmt.Println(req.Referer())
	fmt.Println(req.Cookies())
	fmt.Println(req.UserAgent())
}
