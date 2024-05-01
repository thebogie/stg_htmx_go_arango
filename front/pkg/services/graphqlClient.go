package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GraphQLClientService struct {
	endpoint string
}

func InitGraphQL(endpoint string) *GraphQLClientService {
	return &GraphQLClientService{
		endpoint: endpoint,
	}
}

func (cs *GraphQLClientService) Query(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error {
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cs.endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GraphQL request failed with status code %d: %s", resp.StatusCode, body)
	}

	var respData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return err
	}

	if errors, ok := respData["errors"]; ok {
		return fmt.Errorf("GraphQL request failed: %v", errors)
	}

	if data, ok := respData["data"]; ok {
		if err := json.NewEncoder(bytes.NewBuffer([]byte(fmt.Sprintf("%v", data)))).Encode(result); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("GraphQL response did not contain data")
	}

	return nil
}
