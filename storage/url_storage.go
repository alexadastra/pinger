package storage

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"url-microservice/config"
	"url-microservice/url_service"
)

type UrlStorage interface {
	Save(url *url_service.Url) error
	View() ([]url_service.Url, error)
	ViewIdByUrl (url string) (int, error)
}

type DataBaseUrlStorage struct {
	mutex sync.RWMutex
}

func NewDataBaseUrlStorage() *DataBaseUrlStorage {
	return &DataBaseUrlStorage{}
}

func (storage *DataBaseUrlStorage) Save(url *url_service.Url) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	sqlStr := "INSERT INTO urls(url_string, url_method, time_interval, unix_time_added) VALUES ($1, $2, $3, $4)"

	stmt, err := config.DB.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(url.Url, url.Method, url.TimeInterval, url.TimeCreated)
	if err != nil {
		return fmt.Errorf("could not save new blog post in memory")
	}
	return nil
}

func (storage *DataBaseUrlStorage) View() ([]url_service.Url, error){
	rows, err := config.DB.Query("SELECT * FROM urls")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]url_service.Url, 0)
	for rows.Next() {
		url := url_service.Url{}
		err := rows.Scan(&url.Id, &url.Url, &url.Method, &url.TimeInterval, &url.TimeCreated)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (storage *DataBaseUrlStorage) ViewIdByUrl(url string) (int, error){
	rows, err := config.DB.Query("SELECT * FROM urls WHERE url_string = '" + url + "'")

	if err != nil {
		return 0, err
	}
	defer rows.Close()

	urls := make([]url_service.Url, 0)
	for rows.Next() {
		url := url_service.Url{}
		err := rows.Scan(&url.Id, &url.Url, &url.Method, &url.TimeInterval, &url.TimeCreated)
		if err != nil {
			return 0, err
		}
		urls = append(urls, url)
	}
	if err = rows.Err(); err != nil {
		return 0, err
	}
	id, err := strconv.ParseInt(urls[0].Id, 10, 64)
	return int(id), nil
}
