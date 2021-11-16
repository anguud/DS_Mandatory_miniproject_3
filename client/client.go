package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	proto "github.com/anguud/DS_Mandatory_miniproject_3/proto"
	"google.golang.org/grpc"
)

var amount int64

var clients []proto.ProjectBidClient
var ctx context.Context

func main() {
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	conn2, err := grpc.Dial(":9081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	conn3, err := grpc.Dial(":9082", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	client1 := proto.NewProjectBidClient(conn)
	client2 := proto.NewProjectBidClient(conn2)
	client3 := proto.NewProjectBidClient(conn3)

	clients := make([]proto.ProjectBidClient, 3, 3)

	clients = append(clients, client1, client2, client3)

	ctx = context.Background()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "bid" {
			fmt.Println("amount to bid:")
			scanner.Scan()
			amountString := scanner.Text()
			amountInteger, _ := strconv.Atoi(amountString)
			placeBid(amountInteger)
		} else if text == "result" {
			result()
		}

	}
}

func placeBid(amount int) {
	log.Println("placing bid")
	response := ""
	for _, client := range clients {
		log.Println("client loop")
		bid := proto.Amount{}
		bid.Amount = int64(amount)
		bid.ClientId = "id"

		ack, err := client.Bid(ctx, &bid)
		if err != nil {
			log.Println(err)
		}
		response = ack.Response
	}
	log.Println(response)
}

func result() {

	highestbid := 0
	auctionIsOver := false
	message := proto.Message{}
	for _, client := range clients {
		log.Println(client)
		res, err := client.Result(ctx, &message)
		if err != nil {
			log.Println(err)
		}
		if int(res.HighestBid) > highestbid {
			highestbid = int(res.HighestBid)
			auctionIsOver = res.IsAuctionOver
		}
	}
	if auctionIsOver {
		log.Printf("Auction is over! Highest bid was %d", highestbid)
	} else {
		log.Printf("Auction is NOT over! Highest bid is currently %d", highestbid)
	}
}
