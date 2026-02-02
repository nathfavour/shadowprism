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



func (c *CoreClient) GetHistory() ([]map[string]interface{}, error) {



	var result []map[string]interface{}



	resp, err := c.Http.R().



		SetResult(&result).



		Get("/v1/history")







	if err != nil {



		return nil, err



	}



	if resp.IsError() {



		return nil, fmt.Errorf("core error: %s", resp.Status())



	}



	return result, nil



}







func (c *CoreClient) Shield(amount uint64, dest string, strategy string, force bool) (map[string]interface{}, error) {







	var result map[string]interface{}







	payload := map[string]interface{}{







		"amount_lamports":  amount,







		"destination_addr": dest,







		"strategy":         strategy,







		"force":            force,







	}















	resp, err := c.Http.R().











		SetBody(payload).



		SetResult(&result).



		Post("/v1/shield")







	if err != nil {



		return nil, err



	}



	if resp.IsError() {



		return nil, fmt.Errorf("shield failed: %s", resp.String())



	}



	return result, nil



}







func (c *CoreClient) Swap(amount uint64, from, to string) (map[string]interface{}, error) {



	var result map[string]interface{}



	payload := map[string]interface{}{



		"amount_lamports": amount,



		"from_token":      from,



		"to_token":        to,



	}







	resp, err := c.Http.R().



		SetBody(payload).



		SetResult(&result).



		Post("/v1/swap")







	if err != nil {



		return nil, err



	}



	if resp.IsError() {



		return nil, fmt.Errorf("swap failed: %s", resp.String())



	}



	return result, nil



}







func (c *CoreClient) Pay(amount uint64, merchant string) (map[string]interface{}, error) {



	var result map[string]interface{}



	payload := map[string]interface{}{



		"amount_lamports": amount,



		"merchant_id":     merchant,



	}







	resp, err := c.Http.R().



		SetBody(payload).



		SetResult(&result).



		Post("/v1/pay")







	if err != nil {



		return nil, err



	}



	if resp.IsError() {



		return nil, fmt.Errorf("pay failed: %s", resp.String())



	}



	return result, nil



}







func (c *CoreClient) GetMarket() (map[string]interface{}, error) {



	var result map[string]interface{}



	resp, err := c.Http.R().



		SetResult(&result).



		Get("/v1/market")







	if err != nil {



		return nil, err



	}



	if resp.IsError() {



		return nil, fmt.Errorf("market error: %s", resp.Status())



	}



	return result, nil



}




