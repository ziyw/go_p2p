package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

// server
type server struct {
	Name string
	Addr string
}

// service
func (s *server) SayHello(ctx context.Context, request *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Hello From " + s.Name}, nil
}

func (s *server) Start() {
	// listen to address
	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// register new server
	_s := grpc.NewServer()
	RegisterHelloServiceServer(_s, s)
	if err := _s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

func (s *server) SendRequest(other server) {

	log.Printf("Current Server: %s", s.Name)

	log.Printf("Request to: %s", other.Name)
	conn, err := grpc.Dial(other.Addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect %v", err)
	}
	defer conn.Close()

	client := NewHelloServiceClient(conn)
	req := HelloRequest{Name: s.Name}
	r, err := client.SayHello(context.Background(), &req)
	if err != nil {
		log.Fatalf("Service failed %v", err)
	}
	log.Printf("Reply: %s", r.Message)
}

func main() {

	server1 := server{Name: "NodeOne", Addr: "localhost:60001"}
	server2 := server{Name: "NodeTwo", Addr: "localhost:60002"}
	server3 := server{Name: "NodeThree", Addr: "localhost:60003"}

	go server1.Start()
	go server2.Start()
	go server3.Start()

	server1.SendRequest(server2)
	server1.SendRequest(server3)
	server3.SendRequest(server2)
}
