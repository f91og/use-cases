package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ApiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// NewClient -
func NewClient(host, username, password string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		HostURL:    host,
		Auth: AuthStruct{
			Username: username,
			Password: password,
		},
	}

	err := c.SignIn()
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("x-token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	apiResponse := ApiResponse{}

	if err = json.Unmarshal(body, &apiResponse); err != nil {
		return nil, err
	}

	if int(apiResponse.Code) != 20000 {
		if apiResponse.Message != "nil" {
			return nil, errors.New(apiResponse.Message)
		} else {
			return nil, errors.New("operation failed for unknown reason")
		}
	}

	return apiResponse.Data, nil
}
