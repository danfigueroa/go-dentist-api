package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DBClient *dynamodb.Client

func InitDynamoDB() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolver(
		aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{URL: "http://localhost:8000"}, nil
		}),
	))
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	DBClient = dynamodb.NewFromConfig(cfg)
	log.Println("DynamoDB Local connected")
}
