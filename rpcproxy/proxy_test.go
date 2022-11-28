package rpcproxy

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/grpc"

	"gitlab.momoso.com/cm/inventory/product-search/productsearchpb"
	"gitlab.momoso.com/mms2/utils/lg"
	"gitlab.momoso.com/mms2/utils/service"
)

func TestProxy(t *testing.T) {
	conn, err := service.DialGRPC("product-search")
	lg.PanicError(err)
	p := NewProxy(WithBackend("kaiya", conn))
	lg.PanicError(http.ListenAndServe(":8888", p))
}

func TestInvoke(t *testing.T) {
	addr := service.GetAddress("product-search")
	cc, err := grpc.Dial(addr, grpc.WithInsecure())
	lg.PanicError(err)
	defer cc.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cs, err := cc.NewStream(ctx, &grpc.StreamDesc{}, "/productsearchpb.SearchMaster/CalcSimpleRule")
	lg.PanicError(err)
	req := &productsearchpb.CalcSimpleRuleReq{
		Rule: &productsearchpb.SimpleSearchOrRule{
			AndRules: []*productsearchpb.SimpleSearchAndRule{{
				Items: []*productsearchpb.SearchItem{
					{
						Type:       "tag",
						Val:        "1",
						IsNegative: false,
					},
				},
			}},
		},
	}
	lg.PanicError(cs.SendMsg(req))
	resp := new(productsearchpb.CalcSimpleRuleResp)
	lg.PanicError(cs.RecvMsg(resp))
	lg.PrintJSON(resp)
}
