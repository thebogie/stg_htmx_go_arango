package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"front/pkg/types"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
)

type GraphQLClientService struct {
	endpoint string
}

type Request struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type Response struct {
	Data   interface{} `json:"data"`
	Errors []Error     `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
}

func InitGraphQL(endpoint string) *GraphQLClientService {
	return &GraphQLClientService{
		endpoint: endpoint,
	}
}

func (cs *GraphQLClientService) Query(ctx context.Context, query string, variables map[string]interface{}, result *[]byte) error {
	session := ctx.Value("session").(*sessions.Session)
	currentPlayer := session.Values["currentPlayer"]
	player, ok := currentPlayer.(*types.Player)
	if !ok {
		fmt.Println("Error: Unexpected value for currentPlayer")
		return nil
	}
	authKey := player.AccessToken

	req := Request{
		Query:     query,
		Variables: variables,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("GraphQL errors: %v", err)
	}

	httpReq, err := http.NewRequest("POST", cs.endpoint, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("GraphQL errors: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", authKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("GraphQL errors: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("GraphQL errors: %v", err)
	}

	var graphqlResp Response
	if err := json.Unmarshal(respBody, &graphqlResp); err != nil {
		return fmt.Errorf("GraphQL errors: %v", err)
	}

	if len(graphqlResp.Errors) > 0 {
		return fmt.Errorf("GraphQL errors: %v", graphqlResp.Errors)
	}

	*result, err = json.Marshal(graphqlResp.Data)
	if err != nil {
		return fmt.Errorf("GraphQL errors: %v", err)
	}

	//jsonString := string(result)
	//fmt.Println("JSON data:", jsonString)

	//for _, item := range dataMap {
	//	dataItem := item.(map[string]interface{})
	//	fmt.Printf("DATA: $w", dataItem)
	//}

	//if err := json.Unmarshal([]byte(fmt.Sprintf("%v", graphqlResp.Data)), result); err != nil {
	//	return err
	//}

	return nil
}

type graphqlClientKey struct{}

func WithGraphQLClient(ctx context.Context, client *GraphQLClientService) context.Context {
	return context.WithValue(ctx, graphqlClientKey{}, client)
}
func GraphqlClientFromContext(ctx context.Context) (*GraphQLClientService, bool) {
	client, ok := ctx.Value(graphqlClientKey{}).(*GraphQLClientService)
	return client, ok
}
