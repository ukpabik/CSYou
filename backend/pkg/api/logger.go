package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ukpabik/CSYou/pkg/api/model"
)

func LogSender(log model.Log, addr string, port int) error {

	logBytes, err := json.Marshal(log)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(logBytes)
	url := fmt.Sprintf("http://%s:%d/log", addr, port)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send log, status code: %d", resp.StatusCode)
	}
	return nil
}
