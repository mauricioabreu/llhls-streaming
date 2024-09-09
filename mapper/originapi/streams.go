package originapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mapper/config"
	"mapper/store"
	"net/http"
	"time"
)

type StreamsResponse struct {
	Response []string `json:"response"`
}

type Client struct {
	client *http.Client
	cfg    *config.Config
	storer *store.RedisStore
}

func New(cfg *config.Config, storer *store.RedisStore) *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Second * 3,
		},
		storer: storer,
		cfg:    cfg,
	}
}

func (c *Client) ListStreams(hosts []string) (map[string][]string, error) {
	streams := make(map[string][]string)

	for _, host := range hosts {
		url := fmt.Sprintf("%s/v1/vhosts/default/apps/ll/streams", host)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for host %s: %w", host, err)
		}

		encodedToken := base64.StdEncoding.EncodeToString([]byte(c.cfg.APIToken))
		req.Header.Set("Authorization", "Basic "+encodedToken)

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to perform request for host %s: %w", host, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("failed to get streams from host %s: %s\n", host, resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body for host %s: %w", host, err)
		}

		var streamsResponse StreamsResponse
		if err := json.Unmarshal(body, &streamsResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body for host %s: %w", host, err)
		}

		streams[host] = streamsResponse.Response
	}

	return streams, nil
}
