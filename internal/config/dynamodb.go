package config

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DBClient *dynamodb.Client

func InitDynamoDB() {
	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID && region == "us-west-2" {
				return aws.Endpoint{
					URL:           "http://localhost:8000",
					SigningRegion: "us-west-2",
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		},
	)

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-west-2"),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
				SessionToken:    "",
			},
		}),
	)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	DBClient = dynamodb.NewFromConfig(cfg)
	log.Println("DynamoDB Local connected")

	ensureDentistTableExists()
	ensureProcedureTableExists()
}

func ensureDentistTableExists() {
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String("Dentists"),
	})

	if err != nil {
		var resourceNotFoundError *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundError) {
			log.Println("Table 'Dentists' does not exist. Creating table...")
			_, err := DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
				TableName: aws.String("Dentists"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("ID"),
						KeyType:       types.KeyTypeHash,
					},
				},
				AttributeDefinitions: []types.AttributeDefinition{
					{
						AttributeName: aws.String("ID"),
						AttributeType: types.ScalarAttributeTypeS,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			})
			if err != nil {
				log.Fatalf("Failed to create table: %v", err)
			}
			log.Println("Table 'Dentists' created successfully.")
		} else {
			log.Fatalf("Failed to describe table: %v", err)
		}
	} else {
		log.Println("Table 'Dentists' already exists, skipping creation.")
	}
}

func ensureProcedureTableExists() {
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String("Procedures"),
	})

	if err != nil {
		var resourceNotFoundError *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundError) {
			log.Println("Table 'Procedures' does not exist. Creating table...")
			_, err := DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
				TableName: aws.String("Procedures"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("ID"),
						KeyType:       types.KeyTypeHash,
					},
				},
				AttributeDefinitions: []types.AttributeDefinition{
					{
						AttributeName: aws.String("ID"),
						AttributeType: types.ScalarAttributeTypeS,
					},
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				},
			})
			if err != nil {
				log.Fatalf("Failed to create Procedures table: %v", err)
			}
			log.Println("Table 'Procedures' created successfully.")
		} else {
			log.Fatalf("Failed to describe Procedures table: %v", err)
		}
	} else {
		log.Println("Table 'Procedures' already exists, skipping creation.")
	}
}
