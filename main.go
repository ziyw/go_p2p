package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type node struct {
	Name string
	Addr string

	Clients map[string]HelloServiceClient
}

// Server side service implementation
func (n *node) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello From" + n.Name}, nil
}

// start server
func (n *node) StartServer() {
	lis, err := net.Listen("tcp", n.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	_n := grpc.NewServer()
	RegisterHelloServiceServer(_n, n)

	if err := _n.Serve(lis); err != nil {
		log.Fatalf("failed to serve %v", err)
	}

}

func (n *node) Start() {
	n.Clients = make(map[string]HelloServiceClient)
	go n.StartServer()
}

func (n *node) SendRequest(name string, addr string) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %w", err)
	}

	defer conn.Close()

	n.Clients[name] = NewHelloServiceClient(conn)

	req := HelloRequest{Name: n.Name}
	r, err := n.Clients[name].SayHello(context.Background(), &req)
	if err != nil {
		log.Fatalf("could not greet %v", err)
	}

	log.Printf("Reply from server %s", r.Message)
}

func main() {

	node1 := node{Name: "Node One", Addr: "localhost:60001", Clients: nil}
	node2 := node{Name: "Node Two", Addr: "localhost:60002", Clients: nil}

	node1.Start()
	node2.Start()

	node1.SendRequest("Node Two", "localhost:60002")
}
