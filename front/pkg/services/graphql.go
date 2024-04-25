package services

import (
	"context"
	"encoding/json"
	"fmt"
	"front/pkg/utils"
	"github.com/machinebox/graphql"
)

func Graphql_(query string) {

	// Load the GraphQL mutation from file

	// Create a client for the GraphQL server
	client := graphql.NewClient(utils.GetConfig().GetString("STG_GRAPHQL_API"))

	// Create a GraphQL request with the loaded mutation string
	req := graphql.NewRequest(string(query))

	// Define the variables map
	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"username": "mitch@gmail.com",
			"password": "letmein",
		},
	}

	// Set the variables on the request object
	req.Var("input", variables["input"])

	// Create a context
	ctx := context.Background()

	var respData map[string]interface{}
	if err := client.Run(ctx, req, &respData); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Access the raw data
	rawData, err := json.Marshal(respData)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Print the raw data
	fmt.Println("Raw Data:", string(rawData))
	fmt.Printf(" HEY")
}
