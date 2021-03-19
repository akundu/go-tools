package url_downloader

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/akundu/utilities/jobs"
	"github.com/mreiferson/go-httpclient"
)

var client *http.Client = nil
var transport *httpclient.Transport = nil

type GoTaskTest struct {
	url string
}

func (this *GoTaskTest) Run() *jobs.GoTaskResult {
	if transport == nil {
		transport = &httpclient.Transport{
			ConnectTimeout:        1 * time.Second,
			RequestTimeout:        10 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
		}
	}
	if client == nil {
		client = &http.Client{Transport: transport}
	}

	req, _ := http.NewRequest("GET", this.url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return &jobs.GoTaskResult{nil, err}
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	return &jobs.GoTaskResult{b, err}
}
