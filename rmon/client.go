package rmon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	login      string
	password   string
	userAgent  string
	token      string
}

func NewClient(baseURL, login, password, userAgent string) (*Client, error) {
	client := &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		login:      login,
		password:   password,
		userAgent:  userAgent,
	}

	if err := client.authenticate(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) authenticate() error {
	authURL := fmt.Sprintf("%s/api/v1.0/login", c.baseURL) // Проверьте, что этот URL корректен
	authData := map[string]string{
		"login":    c.login,
		"password": c.password,
	}

	reqBody, err := json.Marshal(authData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Логирование заголовков и тела ответа для диагностики
	log.Printf("Authentication response status: %s", resp.Status)
	log.Printf("Authentication response headers: %v", resp.Header)
	log.Printf("Authentication response body: %s", respBody)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, respBody)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return fmt.Errorf("unable to find token in response: %v", result)
	}

	c.token = token
	return nil
}

func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, respBody)
	}

	return respBody, nil
}
