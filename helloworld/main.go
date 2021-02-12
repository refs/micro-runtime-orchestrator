package main

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/micro/examples/helloworld/proto"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"log"
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
var r = registry.NewRegistry()

type withCancel struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func main() {
	halt := make(chan os.Signal)
	defer close(halt)
	signal.Notify(halt, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	cancellations["global"] = withCancel{
		ctx:    ctx,
		cancel: cancel,
	}

	go run("helloworld", "-1", "0.0.0.0:9200", r)
	go run("helloworld", "-2", "0.0.0.0:9201", r)

	mux := http.ServeMux{}
	mux.Handle("/cancel1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("closing :9200"))
		cancellations["helloworld-1"].cancel()
	}))

	mux.Handle("/cancel2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("closing :9201"))
		cancellations["helloworld-2"].cancel()
	}))

	mux.Handle("/cancelglobal", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("closing all services"))
		cancellations["global"].cancel()
	}))

	s := http.Server{
		Addr:    ":9300",
		Handler: &mux,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Print(err)
		}
	}()

	go func() {
		<-halt
		cancellations["global"].cancel()
	}()

	<-cancellations["global"].ctx.Done()
}

func newGrpcClient() client.Client {
	c := grpc.NewClient(
		client.RequestTimeout(10 * time.Second),
	)
	return c
}

func run(name, instance, address string, r registry.Registry) {
	ctx, cancel := context.WithCancel(cancellations["global"].ctx)
	cancellations[name+instance] = withCancel{
		ctx:    ctx,
		cancel: cancel,
	}

	server := server.NewServer(
		server.Name(name+instance),
		server.Id(uuid.New().String()),
	)

	service := micro.NewService(
		micro.Client(newGrpcClient()),
		micro.Name(name+instance),
		micro.Registry(r),
		micro.Context(ctx),
		micro.Server(server),
		micro.Address(address),
	)

	service.Init()

	pb.RegisterGreeterHandler(service.Server(), new(Greeter))

	go func() {
		if err := service.Run(); err != nil {
			log.Fatal(err)
		}
	}()
}

func initGlobalContext() context.Context {
	global, cancel := context.WithCancel(context.Background())
	cancellations["global"] = withCancel{
		ctx:    global,
		cancel: cancel,
	}
	return global
}
