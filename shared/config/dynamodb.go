package config

import (
	"context"
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

	// Initialize tables for all modules
	ensureDentalTablesExist()
	ensureFinancialTablesExist()
}

// ensureDentalTablesExist creates tables for the dental module
func ensureDentalTablesExist() {
	ensureDentistTableExists()
	ensurePatientTableExists()
	ensureProcedureTableExists()
	ensureAppointmentTableExists()
}

// ensureFinancialTablesExist creates tables for the financial module
func ensureFinancialTablesExist() {
	ensureExpenseTableExists()
	ensureRevenueTableExists()
	ensureInvoiceTableExists()
}

func ensureDentistTableExists() {
	tableName := "Dentists"
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Table %s does not exist, creating...", tableName)
		_, err = DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
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
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			log.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		log.Printf("Table %s created successfully", tableName)
	} else {
		log.Printf("Table %s already exists", tableName)
	}
}

func ensurePatientTableExists() {
	tableName := "Patients"
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Table %s does not exist, creating...", tableName)
		_, err = DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
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
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			log.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		log.Printf("Table %s created successfully", tableName)
	} else {
		log.Printf("Table %s already exists", tableName)
	}
}

func ensureProcedureTableExists() {
	tableName := "Procedures"
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Table %s does not exist, creating...", tableName)
		_, err = DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
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
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			log.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		log.Printf("Table %s created successfully", tableName)
	} else {
		log.Printf("Table %s already exists", tableName)
	}
}

func ensureAppointmentTableExists() {
	tableName := "Appointments"
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Table %s does not exist, creating...", tableName)
		_, err = DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
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
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			log.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		log.Printf("Table %s created successfully", tableName)
	} else {
		log.Printf("Table %s already exists", tableName)
	}
}

func ensureExpenseTableExists() {
	tableName := "Expenses"
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Table %s does not exist, creating...", tableName)
		_, err = DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
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
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			log.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		log.Printf("Table %s created successfully", tableName)
	} else {
		log.Printf("Table %s already exists", tableName)
	}
}

func ensureRevenueTableExists() {
	tableName := "Revenues"
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Table %s does not exist, creating...", tableName)
		_, err = DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
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
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			log.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		log.Printf("Table %s created successfully", tableName)
	} else {
		log.Printf("Table %s already exists", tableName)
	}
}

func ensureInvoiceTableExists() {
	tableName := "Invoices"
	_, err := DBClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		log.Printf("Table %s does not exist, creating...", tableName)
		_, err = DBClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
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
			BillingMode: types.BillingModePayPerRequest,
		})
		if err != nil {
			log.Fatalf("Failed to create table %s: %v", tableName, err)
		}
		log.Printf("Table %s created successfully", tableName)
	} else {
		log.Printf("Table %s already exists", tableName)
	}
}