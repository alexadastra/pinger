package scheduler

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"strconv"
	"time"
	"url-microservice/check"
	"url-microservice/storage"
	"url-microservice/url_service"
)

type CheckScheduler interface {
	AddCheck(url string, method string, time int) error
	RemoveCheck(url string)
	Init(urls []url_service.Url) error
}

type CronCheckScheduler struct{
	scheduler    *gocron.Scheduler
	requester    *check.HttpRequester
	urlStorage *storage.DataBaseUrlStorage
	checkStorage *storage.DataBaseCheckStorage
	jobs         map[string]*gocron.Job
}

func NewCronCheckScheduler() *CronCheckScheduler {
	c := new(CronCheckScheduler)
	c.jobs = make(map[string]*gocron.Job)
	c.scheduler = gocron.NewScheduler(time.UTC)
	c.requester = check.NewHttpRequester()
	c.checkStorage = storage.NewDataBaseCheckStorage()
	c.urlStorage = storage.NewDataBaseUrlStorage()
	return c
}

func (checkScheduler *CronCheckScheduler) Init(urls []url_service.Url) error{
	for i:=0; i < len(urls); i++{
		timeInterval, _ := strconv.ParseInt(urls[i].TimeInterval, 10, 64)
		err := checkScheduler.AddCheck(urls[i].Url, urls[i].Method, int(timeInterval))
		if err != nil {
			return err
		}
	}
	checkScheduler.scheduler.StartAsync()
	return nil
}

func (checkScheduler *CronCheckScheduler) makeUrlRequest(url string, method string) {
	code, err := checkScheduler.requester.DoRequest(url, method)
	if err != nil{
		fmt.Println(nil)
		return
	}
	newCheck := url_service.Check{}
	urlId, err := checkScheduler.urlStorage.ViewIdByUrl(url)
	if err != nil {
		return
	}
	newCheck.Url = strconv.Itoa(urlId)
	newCheck.StatusCode = int32(code)
	newCheck.TimeChecked = strconv.Itoa(int(time.Now().Unix()))
	err = checkScheduler.checkStorage.Save(&newCheck)
	if err != nil {
		return
	}
}

func (checkScheduler *CronCheckScheduler) AddCheck(url string, method string, time int) error{
	job, err := checkScheduler.scheduler.Every(time).Seconds().Do(checkScheduler.makeUrlRequest, url, method)
	if err != nil {
		return err
	}
	checkScheduler.jobs[url] = job
	return nil
}

func (checkScheduler *CronCheckScheduler) RemoveCheck(a string) {
	job := checkScheduler.jobs[a]
	checkScheduler.scheduler.RemoveByReference(job)
}

