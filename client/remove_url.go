package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"time"
	"url-microservice/url_service"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := url_service.NewUrlServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Url: ")
	url, _ := reader.ReadString('\n')
	url = strings.Trim(url, "\n")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.DeleteUrl(ctx, &url_service.CheckGetRequest{
		Url: url,
	})
	if err != nil {
		log.Fatalf("Could not make request :%v", err)
	}
	log.Printf("Url successfully deleted")
}
