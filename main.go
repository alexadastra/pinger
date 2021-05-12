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
	port                      = ":50051"
	defaultRequestMethod      = "get"
	defaultTimeInterval       = 24 * 60 * 60
	defaultLimitForChecksView = 5
)

type UrlServer struct {
	url_service.UnimplementedUrlServiceServer
	urlStorage     storage.UrlStorage       // url repository
	checkStorage   storage.CheckStorage     // check repository
	checkScheduler scheduler.CheckScheduler // scheduler module, that automatically performs url requests
}

func (s *UrlServer) PostUrl(ctx context.Context, req *url_service.UrlPostRequest) (*url_service.UrlPostResponse, error) {
	// get dto from request
	urlDto := req.GetUrl()
	// pull url string
	urlString := urlDto.GetUrl()
	if urlString == "" {
		return nil, errors.New("url cannot be empty")
	}
	// pull time interval to check
	timeInterval := urlDto.GetTimeInterval()
	// if not given, take number of seconds in 24 hours
	if timeInterval == "" {
		timeInterval = strconv.Itoa(defaultTimeInterval)
	}
	i, err := strconv.ParseInt(timeInterval, 10, 64)
	// if not a number, return error
	if err != nil {
		return nil, errors.New("406. Not Acceptable. Time interval must be an integer")
	}
	// pull request method; if not given, take "get"
	requestMethod := urlDto.GetMethod()
	if requestMethod == "" {
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
		return nil, fmt.Errorf("url could not be added: %w", err)
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
	if url == "" {
		return nil, errors.New("url cannot be empty")
	}
	limit := req.GetLimit()
	if limit == 0 {
		limit = defaultLimitForChecksView
	}
	checks, err := s.checkStorage.ViewByUrl(url, int(limit))
	if err != nil {
		return nil, errors.New("error while requesting db")
	}
	// create []CheckDto from []Check by passing status code and formatting unix time to YYYY-MM-DD HH:MM:SS
	checkDtos := make([]*url_service.CheckDto, 0)
	for i := 0; i < len(checks); i++ {
		checkDto := url_service.CheckDto{}
		checkDto.StatusCode = checks[i].StatusCode
		timeInt, _ := strconv.ParseInt(checks[i].TimeChecked, 10, 64)
		checkDto.TimeChecked = time.Unix(timeInt, 0).Format("2006-01-02 15:04:05")
		checkDtos = append(checkDtos, &checkDto)
	}
	res := &url_service.CheckGetResponse{
		Checks: checkDtos,
	}
	return res, nil
}

func (s *UrlServer) DeleteUrl(ctx context.Context, req *url_service.UrlDeleteRequest) (*url_service.UrlDeleteResponse, error) {
	url := req.GetUrl()
	s.checkScheduler.RemoveCheck(url)
	res := &url_service.UrlDeleteResponse{}
	return res, nil
}

func (s *UrlServer) GetUrls(ctx context.Context, req *url_service.UrlGetRequest) (*url_service.UrlGetResponse, error) {
	// parse date as YYYY-MMM-DD (time is set to 00:00 UTC)
	date, err := time.Parse("2006-Jan-02", req.GetDate())
	if err != nil {
		return nil, err
	}
	// convert date to unix time format for querying the database
	dateInt := int(date.Unix())
	n := int(req.GetN())
	urls, err := s.urlStorage.ViewUrlByDateAndN(dateInt, n)
	if err != nil {
		return nil, err
	}
	res := &url_service.UrlGetResponse{
		Urls: urls,
	}
	return res, nil
}

func main() {
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

	url_service.RegisterUrlServiceServer(server, &UrlServer{urlStorage: dataBaseUrlStorage,
		checkStorage: dataBaseCheckStorage, checkScheduler: cronCheckScheduler})

	log.Printf("Server started at port %s", port)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
