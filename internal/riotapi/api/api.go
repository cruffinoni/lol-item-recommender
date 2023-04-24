package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"LoLItemRecommender/internal/printer"
)

type Client struct {
	client *http.Client
	mx     *sync.RWMutex
	limit  *Rate
}

var ErrInvalidStatusCode = errors.New("invalid status code returned, expected 200 got %d")

func NewClient() *Client {
	return &Client{
		client: &http.Client{},
		mx:     &sync.RWMutex{},
		limit:  NewRate(),
	}
}

func generateErrorInvalidStatusCode(status int) error {
	return fmt.Errorf(ErrInvalidStatusCode.Error(), status)
}

func (c *Client) Get(url string) ([]byte, error) {
	canConsume, t := c.limit.CanConsumeTokens()
	if !canConsume {
		for !canConsume {
			printer.Debug("Rate limit exceeded, sleeping for %v | %d tokens consumed", t, c.limit.GetTotalUsage())
			time.Sleep(t)
			printer.Debug("Waking up and try to consume tokens")
			canConsume, t = c.limit.CanConsumeTokens()
		}
	}
	c.mx.Lock()
	defer c.mx.Unlock()
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	var b []byte
	if b, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return b, nil
	case http.StatusTooManyRequests:
		return nil, generateErrorInvalidStatusCode(resp.StatusCode)
	default:
		printer.Error("route: %s resp: '%s' w/ %d ", url, b, resp.StatusCode)
		return nil, generateErrorInvalidStatusCode(resp.StatusCode)
	}
}
