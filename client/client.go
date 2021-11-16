package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"

	proto "github.com/anguud/DS_Mandatory_miniproject_3/proto"
	"google.golang.org/grpc"
)

var amount int64
var clients [3]proto.ProjectBidClient

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

	clients := [client1, client2, client3];

	ctx := context.Background()
	defer disconnect(ctx, client)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		request := proto.Request{}
		request.IpAddress = GetOutboundIP().String()
		request.Message = scanner.Text()

		if hasjoined {
			if inSecret && request.Message == "join" {
				log.Println("Can't join whil in the secret section")
			} else {
				response, err := client.RequestAccess(ctx, &request)

				if err != nil {
					log.Fatal(err)
				}
				if response.Response == ("WELCOME TO CIA SECRET AREA" + string(request.IpAddress)) {
					inSecret = true
				} else if response.Response == ("LEFT THE CIA SECRET AREA" + string(request.IpAddress)) {
					inSecret = false
				}
				log.Println(response)
			}
		} else {
			log.Println("you have to try to join first in order to leave")
		}
	}

}

func disconnect(ctx context.Context, client proto.MutualExclusionClient) {
	request := proto.Request{}
	request.IpAddress = GetOutboundIP().String()
	request.Message = "leave"

	if inSecret {
		_, err := client.RequestAccess(ctx, &request)

		if err != nil {
			log.Println(err)
		}
	}
	log.Println("DISCONNECTED!" + request.IpAddress)
}

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
