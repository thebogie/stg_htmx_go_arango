package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	player := session.Values["currentPlayer"].(types.Player)
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

	if resp.StatusCode == 401 {
		//id wrong user or id.. might need to login again:
		return &types.STGError{
			Msg:  "relogin",
			Code: 401,
		}
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

func (cs *GraphQLClientService) CheckLoginRefresh(w http.ResponseWriter, r *http.Request, err error) {

	session := r.Context().Value("session").(*sessions.Session)
	var stgError *types.STGError
	if errors.As(err, &stgError) {
		switch stgError.Code {
		case 401:
			session.Values["currentPlayer"] = types.Player{
				Firstname:   "",
				Email:       "",
				Password:    "",
				AccessToken: "",
			}
			err = session.Save(r, w)
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)

			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}

type graphqlClientKey struct{}

func WithGraphQLClient(ctx context.Context, client *GraphQLClientService) context.Context {
	return context.WithValue(ctx, graphqlClientKey{}, client)
}
func GraphqlClientFromContext(ctx context.Context) (*GraphQLClientService, bool) {
	client, ok := ctx.Value(graphqlClientKey{}).(*GraphQLClientService)
	return client, ok
}
