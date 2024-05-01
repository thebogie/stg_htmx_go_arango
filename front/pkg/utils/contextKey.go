package utils

import (
	"context"
	"front/pkg/services"
)

type redisClientKey struct{}
type graphqlClientKey struct{}

func WithRedisClient(ctx context.Context, client *services.RedisClient) context.Context {
	return context.WithValue(ctx, redisClientKey{}, client)
}

func RedisClientFromContext(ctx context.Context) (*services.RedisClient, bool) {
	client, ok := ctx.Value(redisClientKey{}).(*services.RedisClient)
	return client, ok
}

func WithGraphQLClient(ctx context.Context, client *services.GraphQLClientService) context.Context {
	return context.WithValue(ctx, graphqlClientKey{}, client)
}

func GraphqlClientFromContext(ctx context.Context) (*services.GraphQLClientService, bool) {
	client, ok := ctx.Value(graphqlClientKey{}).(*services.GraphQLClientService)
	return client, ok
}
