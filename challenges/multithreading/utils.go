package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func doGet[T any](url string, ctx context.Context) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rDto T
	err = json.Unmarshal(body, &rDto)
	if err != nil {
		return nil, err
	}

	return &rDto, nil
}
