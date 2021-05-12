package storage

import (
	"fmt"
	"log"
	"strconv"
	"sync"
	"url-microservice/config"
	"url-microservice/url_service"
)

type CheckStorage interface {
	Save(check *url_service.Check) error                           // save new check in db
	View() ([]url_service.Check, error)                            // get all checks
	ViewByUrl(url string, limit int) ([]*url_service.Check, error) // get 'limit' check of the 'url'
}

type DataBaseCheckStorage struct {
	mutex sync.RWMutex
}

func NewDataBaseCheckStorage() *DataBaseCheckStorage {
	return &DataBaseCheckStorage{}
}

func (storage *DataBaseCheckStorage) Save(check *url_service.Check) error {
	storage.mutex.Lock()
	defer storage.mutex.Unlock()

	sqlStr := "INSERT INTO checks(url_id, status_code, unix_time_added) VALUES ($1, $2, $3)"

	stmt, err := config.DB.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(check.Url, check.StatusCode, check.TimeChecked)
	if err != nil {
		return fmt.Errorf("could not save new blog post in memory")
	}
	return nil
}

func (storage *DataBaseCheckStorage) View() ([]url_service.Check, error) {
	rows, err := config.DB.Query("SELECT * FROM checks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checks := make([]url_service.Check, 0)
	for rows.Next() {
		check := url_service.Check{}
		err := rows.Scan(&check.Id, &check.Url, &check.StatusCode, &check.TimeChecked)
		if err != nil {
			return nil, err
		}
		checks = append(checks, check)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return checks, nil
}

func (storage *DataBaseCheckStorage) ViewByUrl(url string, limit int) ([]*url_service.Check, error) {
	rows, err := config.DB.Query("SELECT c.status_code, c.unix_time_added FROM checks c JOIN urls u ON c.url_id = u.url_id WHERE u.url_string = '" + url + "' ORDER BY c.unix_time_added DESC LIMIT " + strconv.Itoa(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checks := make([]*url_service.Check, 0)
	for rows.Next() {
		check := url_service.Check{}
		err := rows.Scan(&check.StatusCode, &check.TimeChecked)
		if err != nil {
			return nil, err
		}
		checks = append(checks, &check)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return checks, nil
}
