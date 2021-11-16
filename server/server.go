package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strconv"
	"time"

	proto "github.com/anguud/DS_Mandatory_miniproject_3/proto"
	"google.golang.org/grpc"
)

type server struct {
	highestbid    int64
	isAuctionOver bool
	proto.UnimplementedProjectBidServer
}

var port = flag.String("port", "9080", "Port for server")

func main() {

	flag.Parse()
	portnumber := ":" + *port
	list, err := net.Listen("tcp", portnumber)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", portnumber, err)
	}
	server := server{}
	server.highestbid = 0
	server.isAuctionOver = false

	grpc := grpc.NewServer()
	proto.RegisterProjectBidServer(grpc, &server)

	if err := grpc.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
	timer := time.NewTimer(30 * time.Second)
	timer1 := time.NewTimer(2 * time.Second)

	go func() {
		for {
			<-timer.C
			server.isAuctionOver = true
			<-timer1.C
			server.isAuctionOver = false
		}
	}()

}

func (s *server) Bid(ctx context.Context, in *proto.Amount) (*proto.Ack, error) {
	response := ""
	if in.Amount <= s.highestbid {
		response = "oh no bid: " + strconv.Itoa(int(in.Amount)) + " is not high enough"
		log.Println("Bidder bid a amount smaller that highest bid.")
	} else {
		s.highestbid = in.Amount
		response = "You know have the higste bid with your bid " + strconv.Itoa(int(in.Amount))
		log.Println("new highet bid.")
	}
	return &proto.Ack{Response: response}, nil
}

func (s *server) Result(ctx context.Context, in *proto.Message) (*proto.Outcome, error) {
	log.Println("result request")
	return &proto.Outcome{HighestBid: s.highestbid, IsAuctionOver: s.isAuctionOver}, nil
}
