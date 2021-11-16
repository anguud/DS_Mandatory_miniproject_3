package main

import (
	"context"
	"log"
	"net"

	proto "github.com/emredogan/ds_mandatory_exercise_2/proto"
	"google.golang.org/grpc"
)

type server struct {
	proto.UnimplementedMutualExclusionServer
	queue            chan string
	critical_section chan string
	activeConnection string
}

func main() {
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	server := server{}
	server.queue = make(chan string, 50)
	server.critical_section = make(chan string, 1)
	server.activeConnection = ""

	grpc := grpc.NewServer()
	proto.RegisterMutualExclusionServer(grpc, &server)

	if err := grpc.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}

}

func (s *server) RequestAccess(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	response := ""
	if in.Message == "join" {
		s.queue <- in.IpAddress
		log.Println("Client joined the queue " + in.IpAddress)
		s.critical_section <- (<-s.queue)
		s.activeConnection = in.IpAddress
		log.Println("Connected new client ", in.IpAddress)
		response = "WELCOME TO CIA SECRET AREA" + in.IpAddress
	} else if in.Message == "leave" {
		if in.IpAddress == s.activeConnection {
			leftClient := <-s.critical_section
			log.Println("Connected client LEFT ", leftClient)
			response = "LEFT THE CIA SECRET AREA " + in.IpAddress

		} else {
			log.Println("clien trying to leave not in the secret area ", in.IpAddress)
			response = "Not currently in the secret area " + in.IpAddress
		}
	}
	return &proto.Response{Response: response}, nil
}
