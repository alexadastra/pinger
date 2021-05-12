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
	fmt.Print("Enter Url: ")
	url, _ := reader.ReadString('\n')
	url = strings.Trim(url, "\n")
	fmt.Print("Enter limit: ")
	limit, _ := reader.ReadString('\n')
	limit = strings.Trim(limit, "\n")
	limitInt, err := strconv.ParseInt(limit, 10, 64)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	data, err := c.GetChecks(ctx, &url_service.CheckGetRequest{
		Url:   url,
		Limit: int32(limitInt),
	})
	if err != nil {
		log.Fatalf("Could not make request :%v", err)
	}
	if data != nil {
		for i := 0; i < len(data.Checks); i++ {
			fmt.Println(strconv.Itoa(int(data.Checks[i].StatusCode)) + " " + data.Checks[i].TimeChecked)
		}
	} else {
		fmt.Println("nothing received!")
	}
}
