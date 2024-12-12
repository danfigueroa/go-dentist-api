package config

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var DBClient *dynamodb.Client

func InitDynamoDB() {
	// Verificar se existe uma variável de ambiente para o endpoint do DynamoDB
	dynamodbEndpoint := "http://localhost:8000"
	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		dynamodbEndpoint = endpoint
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == dynamodb.ServiceID && region == "us-west-2" {
				return aws.Endpoint{
					URL:           dynamodbEndpoint,
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
	ensurePatientTableExists()
	ensureProcedureTableExists()
	ensureAppointmentTableExists()
}

// ensureDentistTableExists verifica se a tabela "Dentists" já existe.
// Caso não exista, cria a tabela. Se já existir, não faz nada.
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

// ensurePatientTableExists verifica se a tabela "Patients" já existe.
// Caso não exista, cria a tabela. Se já existir, não faz nada.
func ensurePatientTableExists() {
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String("Patients"),
	})

	if err != nil {
		var resourceNotFoundError *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundError) {
			log.Println("Table 'Patients' does not exist. Creating table...")
			_, err := DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
				TableName: aws.String("Patients"),
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
			log.Println("Table 'Patients' created successfully.")
		} else {
			log.Fatalf("Failed to describe table: %v", err)
		}
	} else {
		log.Println("Table 'Patients' already exists, skipping creation.")
	}
}

// ensureProcedureTableExists verifica se a tabela "Procedures" já existe.
// Caso não exista, cria a tabela. Se já existir, não faz nada.
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
				log.Fatalf("Failed to create table: %v", err)
			}
			log.Println("Table 'Procedures' created successfully.")
		} else {
			log.Fatalf("Failed to describe table: %v", err)
		}
	} else {
		log.Println("Table 'Procedures' already exists, skipping creation.")
	}
}

// ensureAppointmentTableExists verifica se a tabela "Appointments" já existe.
// Caso não exista, cria a tabela. Se já existir, não faz nada.
func ensureAppointmentTableExists() {
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String("Appointments"),
	})

	if err != nil {
		var resourceNotFoundError *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundError) {
			log.Println("Table 'Appointments' does not exist. Creating table...")
			_, err := DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
				TableName: aws.String("Appointments"),
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
			log.Println("Table 'Appointments' created successfully.")
		} else {
			log.Fatalf("Failed to describe table: %v", err)
		}
	} else {
		log.Println("Table 'Appointments' already exists, skipping creation.")
	}
}
