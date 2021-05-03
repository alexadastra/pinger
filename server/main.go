package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
	"url-microservice/scheduler"
	"url-microservice/storage"
	"url-microservice/url_service"
)

const (
	port = ":50051"
	defaultRequestMethod = "get"
	defaultTimeInterval = 24 * 60 * 60
	defaultLimitForChecksView = 5
)

type UrlServer struct {
	url_service.UnimplementedUrlServiceServer
	urlStorage storage.UrlStorage
	checkStorage storage.CheckStorage
	checkScheduler scheduler.CheckScheduler
}

func (s *UrlServer) PostUrl(ctx context.Context, req *url_service.UrlPostRequest) (*url_service.UrlPostResponse, error) {
	// pull url string
	urlString := req.GetUrl()
	if urlString == "" {
		return nil, errors.New("url cannot be empty")
	}
	// pull time interval to check;
	// if not given, take number of seconds in 24 hours
	// if not a number, return error
	timeInterval := req.GetTimeInterval()
	if timeInterval == "" {
		timeInterval = strconv.Itoa(defaultTimeInterval)
	}
	i, err := strconv.ParseInt(timeInterval, 10, 64)
	if err != nil {
		return nil, errors.New("406. Not Acceptable. Time interval must be an integer")
	}
	// pull request method; if not given, take "get"
	requestMethod := req.GetMethod()
	if requestMethod == ""{
		requestMethod = defaultRequestMethod
	}
	// form Url struct
	url := url_service.Url{}
	url.Url = urlString
	url.TimeCreated = strconv.FormatInt(time.Now().Unix(), 10)
	url.Method = requestMethod
	url.TimeInterval = timeInterval
	// save url in database
	err = s.urlStorage.Save(&url)
	if err != nil {
		return nil, fmt.Errorf("url could not be added: %w",err)
	}
	// add url to checking scheduler
	err = s.checkScheduler.AddCheck(url.Url, url.Method, int(i))
	// form response
	res := &url_service.UrlPostResponse{
		Url: &url,
	}
	log.Printf("Url has been added")
	return res, nil
}

func (s *UrlServer) GetChecks(ctx context.Context, req *url_service.CheckGetRequest) (*url_service.CheckGetResponse, error) {
	url := req.GetUrl()
	if url == ""{
		return nil, errors.New("url cannot be empty")
	}
	limit := req.GetLimit()
	if limit == 0{
		limit = defaultLimitForChecksView
	}
	checks, err := s.checkStorage.ViewByUrl(url, int(limit))
	if err != nil{
		return nil, errors.New("error while requesting db")
	}
	res := &url_service.CheckGetResponse{
		Checks: checks,
	}
	return res, nil
}

func (s *UrlServer) DeleteUrl(ctx context.Context, req *url_service.CheckGetRequest) (*url_service.CheckGetRequest, error) {
	url := req.GetUrl()
	s.checkScheduler.RemoveCheck(url)
	res := &url_service.CheckGetRequest{
		Url: url,
	}
	return res, nil
}

func main(){
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()

	dataBaseUrlStorage := storage.NewDataBaseUrlStorage()
	dataBaseCheckStorage := storage.NewDataBaseCheckStorage()

	cronCheckScheduler := scheduler.NewCronCheckScheduler()
	urls, _ := dataBaseUrlStorage.View()
	_ = cronCheckScheduler.Init(urls)

	url_service.RegisterUrlServiceServer(server,&UrlServer{urlStorage: dataBaseUrlStorage,
			checkStorage:dataBaseCheckStorage, checkScheduler: cronCheckScheduler})

	log.Printf("Server started at port %s", port)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
