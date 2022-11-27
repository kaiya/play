package rpcproxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"vendor/github.com/golang/protobuf/jsonpb"
	"vendor/github.com/jhump/protoreflect/grpcreflect"
	"vendor/google.golang.org/grpc"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/proto"

	grv1 "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"github.com/gorilla/mux"
)

type Proxy struct {
	http.Handler
	backends   map[string]*grpc.ClientConn
	descSource map[string]grpcurl.DescriptorSource
	mu         sync.RWMutex
}

func New() *Proxy {
	p := &Proxy{}
	router := mux.NewRouter()
	router.Path("{backend}/{package}/{service}/{method}").Methods("POST").HandlerFunc(p.handle)
	p.Handler = router
	return p
}

func (p *Proxy) SetBackend(name string, cc *grpc.ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.backends[name] = cc
	p.descSource[name] = getDescSource(context.Background(), cc)
}

func getDescSource(ctx context.Context, cc *grpc.ClientConn) grpcurl.DescriptorSource {
	cli := grv1.NewServerReflectionClient(cc)

	return grpcurl.DescriptorSourceFromServer(ctx, grpcreflect.NewClient(ctx, cli))
}

func newRequestSupplier(json []byte) grpcurl.RequestSupplier {
	return func(m proto.Message) error {
		return jsonpb.Unmarshal(bytes.NewReader(json), m)
	}
}

func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	be := vars["backend"]
	pkg := vars["package"]
	srv := vars["service"]
	method := vars["method"]
	backend, ok := p.backends[be]
	if !ok {
		http.Error(w, fmt.Sprintf("backend:%s", be), http.StatusNotFound)
		return
	}
	ctx := r.Context()
	desc := p.descSource[be]
	headers := r.URL.Query()["metadata"]
	var buf bytes.Buffer
	// get metadata
	// read request in
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("read req body error:%s", err), http.StatusInternalServerError)
		return
	}
	// for receiving resp
	handler := grpcurl.NewDefaultEventHandler(&buf, desc, grpcurl.NewJSONFormatter(false, grpcurl.AnyResolverFromDescriptorSource(desc)), true)
	err = grpcurl.InvokeRPC(ctx, getDescSource(ctx, backend), backend, fmt.Sprintf("%s.%s.%s", pkg, srv, method), headers, handler, newRequestSupplier(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("invoke grpc error:%s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}
