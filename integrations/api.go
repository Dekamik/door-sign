package integrations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func call[T any](req *http.Request) (*T, error) {
	client := &http.Client{}
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
		return nil, fmt.Errorf("Get returned the following error:\n%s\n", string(bodyBytes))
	}

	var response T
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
