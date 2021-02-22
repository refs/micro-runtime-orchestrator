package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/asim/go-micro/plugins/client/grpc/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/google/uuid"
	pb "github.com/refs/whatever/helloworld/proto"
	"github.com/thejerf/suture"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	rsp.Greeting = "Hello " + req.Name
	return nil
}

// cancellations represents a slice of service name + contexts to cancel spawned services
// eventually move this to the runtime type.
var cancellations = make(map[string]withCancel)

//var r = registry.NewMemoryRegistry()
var r = registry.NewRegistry()

type withCancel struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var services = make(map[string]suture.ServiceToken)

func main() {
	supervisor := suture.NewSimple("universe")
	ctx, cancel := context.WithCancel(context.Background())
	// we can override go-micro defaults to correctly propagate our functionality.
	registry.DefaultRegistry = r

	go func() {
		watcher, err := r.Watch()
		if err != nil {
			panic(err)
		}

		for {
			if v, err := watcher.Next(); err == nil {
				d, _ := json.Marshal(v)
				fmt.Printf("event: %s\n", string(d))
			}
		}
	}()

	halt := make(chan os.Signal, 1)
	signal.Notify(halt, os.Interrupt)

	sMain := mainService{
		supervisor: supervisor,
		ctx:        ctx,
		cancel:     cancel,
	}

	s1ctx, s1cancel := context.WithCancel(ctx)
	s1addr := "0.0.0.0:9200"
	s1 := service{
		name:    "hello-world-1",
		address: &s1addr,
		r:       r,
		ctx:     s1ctx,
		cancel:  s1cancel,
	}

	s2ctx, s2cancel := context.WithCancel(ctx)
	s2addr := "0.0.0.0:9201"
	s2 := service{
		name:    "hello-world-2",
		address: &s2addr,
		r:       r,
		ctx:     s2ctx,
		cancel:  s2cancel,
	}

	services["main"] = supervisor.Add(sMain)
	services["s1"] = supervisor.Add(s1)
	services["s2"] = supervisor.Add(s2)

	go supervisor.Serve()

	for {
		select {
		case <-ctx.Done():
			supervisor.Stop()
			return
		case <-halt:
			cancel()
			supervisor.Stop()
			return
		}
	}
}

func newGrpcClient() client.Client {
	c := grpc.NewClient(
		client.RequestTimeout(10 * time.Second),
	)
	return c
}

// service implements the suture.Service interface.
type service struct {
	name    string
	address *string
	r       registry.Registry
	ctx     context.Context
	cancel  context.CancelFunc
}

func (s service) Serve() {
	microServer := server.NewServer(
		server.Name(s.name),
		server.Id(uuid.New().String()),
	)

	service := micro.NewService(
		micro.Client(newGrpcClient()),
		micro.Name(s.name),
		micro.Registry(r),
		micro.Server(microServer),
		micro.Context(s.ctx),
	)

	// we DON'T need to call init as we are not using env vars or cli flags.
	//service.Init()

	pb.RegisterGreeterHandler(service.Server(), new(Greeter))

	if err := service.Run(); err != nil {
		// [...] if this function either returns or panics, the Supervisor will call it again.
		// https://pkg.go.dev/github.com/thejerf/suture#hdr-Serve_Method
		panic(err)
	}
}

func (s service) Stop() {
	s.cancel() // calling cancel will cause a go-micro service to end and unregister it from the service registry.
}

// mainService provides control to a supervisor tree using simple http endpoints
type mainService struct {
	supervisor *suture.Supervisor
	ctx        context.Context
	cancel     context.CancelFunc
}

func (m mainService) Serve() {
	mux := http.ServeMux{}
	mux.Handle("/create-random", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		sname := req.Header.Get("sname")
		sx, sxcancel := context.WithCancel(m.ctx)
		serv := service{
			name:   sname,
			r:      r,
			ctx:    sx,
			cancel: sxcancel,
		}

		w.Write([]byte(fmt.Sprintf("created service %s", sname)))
		m.supervisor.Add(serv)
	}))

	mux.Handle("/cancel1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("closing :9200"))
		m.supervisor.Remove(services["s1"])
	}))

	mux.Handle("/cancel2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("closing :9201"))
		m.supervisor.Remove(services["s2"])
	}))

	mux.Handle("/cancelglobal", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("closing all services"))
		m.supervisor.Remove(services["main"])
	}))

	s := http.Server{
		Addr:    ":9300",
		Handler: &mux,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-m.ctx.Done()
}

func (m mainService) Stop() {
	m.cancel()
}
