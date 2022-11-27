package rpcproxy

import (
	"net/http"
	"testing"

	"gitlab.momoso.com/mms2/utils/lg"
	"gitlab.momoso.com/mms2/utils/service"
)

func TestProxy(t *testing.T) {
	conn, err := service.DialGRPC("product-search")
	lg.PanicError(err)
	p := NewProxy(WithBackend("kaiya", conn))
	lg.PanicError(http.ListenAndServe(":8888", p))
}
