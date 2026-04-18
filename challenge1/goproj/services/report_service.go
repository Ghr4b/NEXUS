package services

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func VisitLocalUri(uri string) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("http://localhost/%s", uri)

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
