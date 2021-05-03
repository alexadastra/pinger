package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"url-microservice/url_service"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := url_service.NewUrlServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Url: ")
	url, _ := reader.ReadString('\n')
	url = strings.Trim(url, "\n")
	fmt.Print("Enter time interval: ")
	interval, _ := reader.ReadString('\n')
	interval = strings.Trim(interval, "\n")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.PostUrl(ctx, &url_service.UrlPostRequest{
		Url: url,
		TimeInterval: interval,
	})
	if err != nil {
		log.Fatalf("Could not create Blog Post :%v", err)
	}
	log.Printf("Post Successfully Created")
}

