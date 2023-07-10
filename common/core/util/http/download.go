package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type downloader struct {
	data      map[string]string
	maxInTime int
}

func NewDownloader(data map[string]string, maxInTime int) *downloader {
	return &downloader{data: data, maxInTime: maxInTime}
}

func (d *downloader) Download(url, fileName string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("url:%v, err:%v", url, err.Error())
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("url:%v, err:%v", url, err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("url:%v, err:%v", url, err.Error())
		}
		return os.WriteFile(fileName, b, os.ModePerm)
	}

	return fmt.Errorf("url:%v, code:%v", url, resp.StatusCode)
}

func (d *downloader) DownloadAll() {
	if len(d.data) == 0 {
		return
	}
	wg := &sync.WaitGroup{}
	limitChan := make(chan struct{}, d.maxInTime)
	for url, fileName := range d.data {
		limitChan <- struct{}{}
		wg.Add(1)
		go func(url, fileName string, wg *sync.WaitGroup) {
			err := d.Download(url, fileName)
			if err != nil {
				log.Errorf("[DownloadAll] download err:%v", err)
			}

			<-limitChan
			wg.Done()
		}(url, fileName, wg)
	}
	wg.Wait()
	close(limitChan)

}
