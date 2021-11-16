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
	var errorcount = 0
	response := ""
	for _, client := range clients {
		bid := proto.Amount{}
		bid.Amount = int64(amount)
		bid.ClientId = "id"

		ack, err := client.Bid(ctx, &bid)
		if err != nil {
			errorcount = errorcount + 1
			// log.Println(err)
			continue
		}
		response = ack.Response
	}
	if errorcount < len(clients) {
		log.Println(response)
	} else {
		log.Panicln("all replicas down")
	}
}

func result() {
	var errorcount = 0
	highestbid := 0
	auctionIsOver := false
	message := proto.Message{}
	for _, client := range clients {
		res, err := client.Result(ctx, &message)
		if err != nil {
			errorcount = errorcount + 1
			// log.Println(err)
			continue
		}
		if int(res.HighestBid) > highestbid {
			highestbid = int(res.HighestBid)
			auctionIsOver = res.IsAuctionOver
		}
	}
	if errorcount < len(clients) {
		if auctionIsOver {
			log.Printf("Auction is over! Highest bid was %d", highestbid)
		} else {
			log.Printf("Auction is NOT over! Highest bid is currently %d", highestbid)
		}
	} else {
		log.Panicln("all replicas down")
	}
}
