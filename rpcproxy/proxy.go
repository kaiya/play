package rpcproxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"gitlab.momoso.com/mms2/utils/lg"
	"google.golang.org/grpc"

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

func NewProxy(opts ...Option) *Proxy {
	p := &Proxy{
		backends:   make(map[string]*grpc.ClientConn),
		descSource: make(map[string]grpcurl.DescriptorSource),
	}
	for _, opt := range opts {
		opt(p)
	}
	router := mux.NewRouter()
	router.Path("/{backend}/{package}/{service}/{method}").Methods("POST").HandlerFunc(p.handle)
	p.Handler = router
	return p
}

func (p *Proxy) SetBackend(name string, cc *grpc.ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.backends[name] = cc
	p.descSource[name] = getDescSource(context.Background(), cc)
}

type Option func(*Proxy)

func WithBackend(name string, cc *grpc.ClientConn) Option {
	return func(p *Proxy) {
		p.backends[name] = cc
		p.descSource[name] = getDescSource(context.Background(), cc)
	}
}

func getDescSource(ctx context.Context, cc *grpc.ClientConn) grpcurl.DescriptorSource {
	cli := grv1.NewServerReflectionClient(cc)
	return grpcurl.DescriptorSourceFromServer(ctx, grpcreflect.NewClient(ctx, cli))
}

func newRequestSupplier(json []byte) grpcurl.RequestSupplier {
	first := true
	return func(m proto.Message) error {
		if first {
			first = false
			return jsonpb.Unmarshal(bytes.NewReader(json), m)
		}
		return io.EOF
	}
}

func (p *Proxy) handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	be := vars["backend"]
	pkg := vars["package"]
	srv := vars["service"]
	method := vars["method"]
	p.mu.RLock()
	backend, ok := p.backends[be]
	desc := p.descSource[be]
	p.mu.RUnlock()
	if !ok {
		http.Error(w, fmt.Sprintf("backend:%s", be), http.StatusNotFound)
		return
	}
	ctx := r.Context()
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
	handler := grpcurl.NewDefaultEventHandler(&buf, desc, grpcurl.NewJSONFormatter(false, grpcurl.AnyResolverFromDescriptorSource(desc)), false)

	start := time.Now()
	err = grpcurl.InvokeRPC(ctx, desc, backend, fmt.Sprintf("%s.%s.%s", pkg, srv, method), headers, handler, newRequestSupplier(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("invoke grpc error:%s", err), http.StatusInternalServerError)
		return
	}
	lg.Info(time.Since(start).String())
	if handler.Status.Err() != nil {
		http.Error(w, fmt.Sprintf("event handler error:%s", handler.Status.Message()), http.StatusInternalServerError)
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// desc source impl cached wrapper
type CacheDescSource struct {
	source   grpcurl.DescriptorSource
	services []string
	descFDs  []*desc.FieldDescriptor
	lock     sync.RWMutex
	cache    map[string]interface{}
}

// ListServices returns a list of fully-qualified service names. It will be all services in a set of
// descriptor files or the set of all services exposed by a gRPC server.
func (cds *CacheDescSource) ListServices() ([]string, error) {
	if len(cds.services) == 0 {
		srvs, err := cds.source.ListServices()
		if err != nil {
			return cds.services, err
		}
		cds.lock.Lock()
		cds.services = srvs
		cds.lock.Unlock()
	}
	return cds.services, nil
}

// FindSymbol returns a descriptor for the given fully-qualified symbol name.
func (cds *CacheDescSource) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	d, ok := cds.cache["descriptor"]
	if !ok {
		de, err := cds.FindSymbol(fullyQualifiedName)
		if err != nil {
			return nil, err
		}
		cds.cache["descriptor"] = de
	}
	return d.(desc.Descriptor), nil
}

// AllExtensionsForType returns all known extension fields that extend the given message type name.
func (cds *CacheDescSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {
	if len(cds.descFDs) == 0 {
		fds, err := cds.source.AllExtensionsForType(typeName)
		if err != nil {
			return nil, err
		}
		cds.descFDs = fds
	}
	return cds.descFDs, nil
}
