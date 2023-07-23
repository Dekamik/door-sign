package integrations

import (
	"door-sign/configuration"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func get[T any](url string) (*T, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return call[T](req)
}

func call[T any](req *http.Request) (*T, error) {
	client := &http.Client{}
	userAgent := fmt.Sprintf("door-sign/%s", configuration.ReadVersion())
	req.Header.Set("User-Agent", userAgent)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("%s %s returned the following error:\n%s\n", req.Method, req.URL.String(), string(bodyBytes))
	}
	log.Printf("%s %s - %s\n", req.Method, req.URL.String(), res.Status)

	var response T
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
