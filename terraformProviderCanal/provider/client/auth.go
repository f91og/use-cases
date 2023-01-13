package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

// SignIn - Get a new token for user
func (c *Client) SignIn() error {
	if c.Auth.Username == "" || c.Auth.Password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}
	jsonEncodedBody, err := json.Marshal(c.Auth)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/user/login", c.HostURL), strings.NewReader(string(jsonEncodedBody)))
	if err != nil {
		return err
	}

	data, err := c.doRequest(req)
	if err != nil {
		return err
	}

	signInResponse := SignInResponse{}
	err = json.Unmarshal(data, &signInResponse)
	if err != nil {
		return err
	}

	c.Token = signInResponse.Token
	return nil
}

// SignOut - Revoke the token for a user
func (c *Client) SignOut() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/user/logout", c.HostURL), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}
	c.Token = ""
	return nil
}
