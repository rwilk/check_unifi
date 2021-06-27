package main

// code partially based on: github.com/paultyng/go-unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"
)

type APIUnifi struct {
	sync.Mutex

	csrf       string
	url        string
	httpclient *http.Client
}

func (api *APIUnifi) Login(ctx context.Context, URL string, user string, password string) (err error) {

	if api.httpclient == nil {
		api.httpclient = &http.Client{}

		cj, _ := cookiejar.New(nil)
		api.httpclient.Jar = cj

		api.url = strings.TrimRight(URL, "/")
	}

	err = api.request(ctx, "POST", api.url+"/api/login", &struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: user,
		Password: password,
	}, nil)

	if err != nil {
		return err
	}

	return nil

}

func (api *APIUnifi) GetDeviceBasic(ctx context.Context, site string) (err error, deviceBasic *DeviceBasicResponse) {

	var respBody DeviceBasicResponse

	url := fmt.Sprintf("%s/api/s/%s/stat/device-basic", api.url, site)
	err = api.request(ctx, "GET", url, nil, &respBody)

	if err != nil {
		return err, nil
	}

	return nil, &respBody

}

func (api *APIUnifi) request(ctx context.Context, method, url string, reqBody interface{}, respBody interface{}) error {
	api.Lock()
	defer api.Unlock()

	var (
		reqReader io.Reader
		err       error
		reqBytes  []byte
	)

	if reqBody != nil {
		reqBytes, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("unable to marshal JSON: %s %s %w", method, url, err)
		}
		reqReader = bytes.NewReader(reqBytes)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqReader)
	if err != nil {
		return fmt.Errorf("unable to create request: %s %s %w", method, url, err)
	}

	req.Header.Set("User-Agent", "check-unifi/0.1")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")

	if api.csrf != "" {
		req.Header.Set("X-CSRF-Token", api.csrf)
	}

	resp, err := api.httpclient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to perform request: %s %s %w", method, url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &NotFoundError{}
	}

	if csrf := resp.Header.Get("x-csrf-token"); csrf != "" {
		api.csrf = resp.Header.Get("x-csrf-token")
	}

	if resp.StatusCode != 200 {
		errBody := struct {
			Meta meta `json:"meta"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&errBody)
		return fmt.Errorf("%w (%s) for %s %s", errBody.Meta.error(), resp.Status, method, url)
	}

	if respBody == nil || resp.ContentLength == 0 {
		return nil
	}

	err = json.NewDecoder(resp.Body).Decode(respBody)
	if err != nil {
		return fmt.Errorf("unable to decode body: %s %s %w", method, url, err)
	}

	return nil

}

type NotFoundError struct{}

func (err *NotFoundError) Error() string {
	return "not found"
}

type APIError struct {
	RC      string
	Message string
}

func (err *APIError) Error() string {
	return err.Message
}

type meta struct {
	RC      string `json:"rc"`
	Message string `json:"msg"`
}

func (m *meta) error() error {
	if m.RC != "ok" {
		return &APIError{
			RC:      m.RC,
			Message: m.Message,
		}
	}

	return nil
}
