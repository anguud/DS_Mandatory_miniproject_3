package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"

	proto "github.com/emredogan/ds_mandatory_exercise_2/proto"
	"google.golang.org/grpc"
)

var hasjoined bool
var inSecret bool

func main() {
	conn, err := grpc.Dial(":9080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}

	client := proto.NewMutualExclusionClient(conn)
	ctx := context.Background()
	defer disconnect(ctx, client)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		request := proto.Request{}
		request.IpAddress = GetOutboundIP().String()
		request.Message = scanner.Text()
		if request.Message == "join" {
			hasjoined = true
		}

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
