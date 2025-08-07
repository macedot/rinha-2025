package services

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type Client struct {
	client *fasthttp.Client
}

var client Client

var headerContentTypeJSON = []byte("application/json")

func HttpClientInstance() *Client {
	return &client
}

func (c *Client) Init() {
	c.client = &fasthttp.Client{
		ReadTimeout:                   5 * time.Second,
		WriteTimeout:                  5 * time.Second,
		MaxIdleConnDuration:           1 * time.Hour,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
}

func (c *Client) Get(url string) (int, []byte) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()
	err := c.client.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		log.Fatalf("HTTP GET request failed: %v", err)
	}
	statusCode := resp.StatusCode()
	if statusCode != fasthttp.StatusOK {
		log.Fatalf("HTTP GET request failed: %v", err)
	}
	return statusCode, resp.Body()

	// resp, err := c.http.Get(url)
	// if err != nil {
	// 	log.Fatalf("HTTP GET request failed: %v", err)
	// }
	// defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// if err == nil {
	// 	return resp.StatusCode, body
	// }
	// return resp.StatusCode, nil
}

func (c *Client) Post(url string, payload []byte) error {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJSON)
	req.SetBodyRaw(payload)
	resp := fasthttp.AcquireResponse()
	err := c.client.DoTimeout(req, resp, 1*time.Second)
	fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return fmt.Errorf("HTTP POST request failed: %v", err)
	}
	statusCode := resp.StatusCode()
	if statusCode != http.StatusOK {
		return fmt.Errorf("invalid HTTP response code: %d", statusCode)
	}
	return nil

	// resp, err := c.client.Post(url, "application/json", bytes.NewReader(payload))
	// if err != nil {
	// 	return fmt.Errorf("HTTP POST request failed: %v", err)
	// }
	// defer resp.Body.Close()
	// return nil
}
