package sidecar

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

type CoreClient struct {
	Http   *resty.Client
	BaseUrl string
}

func NewCoreClient(port int, token string) *CoreClient {
	return &CoreClient{
		Http: resty.New().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
			SetBaseURL(fmt.Sprintf("http://localhost:%d", port)),
	}
}

func (c *CoreClient) GetStatus() (map[string]interface{}, error) {
	var result map[string]interface{}
	resp, err := c.Http.R().
		SetResult(&result).
		Get("/health")
	
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("core error: %s", resp.Status())
	}
	return result, nil
}
