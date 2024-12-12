package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	sgrpc "go.unistack.org/micro-server-grpc/v3"
	"go.unistack.org/micro/v3"
	"go.unistack.org/micro/v3/server"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	switch os.Args[1] {
	case "server":
		TestServer()
	case "client":
		TestClient()
	}
}

type serv struct {
	UnimplementedTestServiceServer
}

func (s *serv) DoWork(ctx context.Context, in *emptypb.Empty) (*WorkResponse, error) {
	fmt.Println("Starting long-running operation")
	time.Sleep(4 * time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("Operation interrupted")
		return nil, ctx.Err()
	default:
		fmt.Println("Operation completed")
		return &WorkResponse{Message: "Work done"}, nil
	}
}

func startServer(ctx context.Context) {
	s := sgrpc.NewServer(server.Name("Service"), server.Address("localhost:1234"))
	svc := micro.NewService(
		micro.Context(ctx),
		micro.Server(s),
	)
	svc.Init()
	RegisterTestServiceServer(s.GRPCServer(), &serv{})

	go func() {
		fmt.Printf("wait for ctx.Done\n")
		<-ctx.Done()
		fmt.Printf("wait for Stop\n")
		svc.Stop()
		fmt.Printf("Stopped\n")
	}()

	fmt.Printf("svc Run\n")
	svc.Run()
	fmt.Printf("svc End\n")
}

func TestClient() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("create grpc conn\n")
	conn, _ := grpc.NewClient("localhost:1234", grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()
	cli := NewTestServiceClient(conn)

	replyCh := make(chan string)
	go func() {
		resp, err := cli.DoWork(ctx, &emptypb.Empty{})
		if err != nil {
			fmt.Println("Client call failed:", err)
			replyCh <- ""
		} else {
			replyCh <- resp.Message
		}
	}()

	p, _ := os.FindProcess(os.Getpid())
	_ = p
	//_ = p.Signal(syscall.SIGTERM)

	select {
	case reply := <-replyCh:
		if reply != "Work done" {
			log.Printf("Expected reply 'Work done', got '%s'\n", reply)
		} else {
			log.Printf("all fine\n")
		}
	case <-ctx.Done():
		log.Printf("Request was not completed\n")
	}
}

func TestServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sigReceived := <-sig
		fmt.Printf("handle signal %v, exiting\n", sigReceived)
		cancel()
	}()

	log.Printf("run server\n")
	startServer(ctx)
}
