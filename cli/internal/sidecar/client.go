package sidecar

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

type CoreClient struct {
	Http   *resty.Client
	Socket string
}

func NewCoreClient(socketPath string, token string) *CoreClient {
	client := resty.New()

	// Configure UDS Dialer
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}

	client.SetTransport(transport).
		SetBaseURL("http://localhost"). // Hostname is ignored by UDS dialer
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", token)).
		SetTimeout(5 * time.Second)

	return &CoreClient{
		Http:   client,
		Socket: socketPath,
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