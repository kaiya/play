package httprouter

import (
	"fmt"
	"net/http"
	"testing"

	"gitlab.momoso.com/cm/kit/third_party/lg"
)

func TestHttpRouter(t *testing.T) {
	hr := New()
	hr.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("opps, resource not found"))
	})
	hr.GET("/hello", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte("hello"))
	})
	hr.GET("/hello/:name", func(w http.ResponseWriter, r *http.Request, p Params) {
		w.Write([]byte(fmt.Sprintf("hello, %s", p.ByName("name"))))
	})
	lg.PanicError(http.ListenAndServe(":8888", hr))
}
