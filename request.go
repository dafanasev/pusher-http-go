package pusher

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/rs/zerolog"
)

type BeforeRequestHandler func(req *http.Request) *http.Request
type AfterRequestHandler func(req *http.Request, resp *http.Response, err error)

// change timeout to time.Duration
func request(client *http.Client, method, url string, body []byte, logger *zerolog.Logger, before BeforeRequestHandler, after AfterRequestHandler) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if before != nil {
		req = before(req)
	}
	resp, err := client.Do(req)
	if after != nil {
		after(req, resp, err)
	}
	if err != nil {
		if logger != nil {
			logger.Error().Err(err).Interface("request", req).Str("url", url).Bytes("body", body).Msg("Cannot do http request")
		}
		return nil, err
	}
	defer resp.Body.Close()
	return processResponse(resp, logger)
}

func processResponse(response *http.Response, logger *zerolog.Logger) ([]byte, error) {
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		if logger != nil {
			logger.Error().Err(err).Interface("response", response).Msg("cannot read response body")
		}
		return nil, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return responseBody, nil
	}
	message := fmt.Sprintf("Status Code: %s - %s", strconv.Itoa(response.StatusCode), string(responseBody))
	err = errors.New(message)
	if logger != nil {
		logger.Error().Err(err).Msg(message)
	}
	return nil, err
}
