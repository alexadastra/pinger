package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"url-microservice/url_service"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := url_service.NewUrlServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter date: ")
	date, _ := reader.ReadString('\n')
	date = strings.Trim(date, "\n")
	fmt.Print("Enter n: ")
	n, _ := reader.ReadString('\n')
	n = strings.Trim(n, "\n")
	limitInt, err := strconv.ParseInt(n, 10, 64)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	data, err := c.GetUrls(ctx, &url_service.UrlGetRequest{
		Date: date,
		N:    int32(limitInt),
	})
	if err != nil {
		log.Fatalf("Could not make request :%v", err)
	}
	if data != nil {
		for i := 0; i < len(data.Urls); i++ {
			fmt.Println(data.Urls[i])
		}
	}
}
