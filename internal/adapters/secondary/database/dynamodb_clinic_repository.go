package database

import (
	"context"
	"fmt"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports/repositories"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClinicRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBClinicRepository(client *dynamodb.Client, tableName string) repositories.ClinicRepository {
	return &DynamoDBClinicRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBClinicRepository) Create(ctx context.Context, clinic *domain.Clinic) error {
	item, err := attributevalue.MarshalMap(clinic)
	if err != nil {
		return fmt.Errorf("failed to marshal clinic: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func (r *DynamoDBClinicRepository) GetByID(ctx context.Context, id string) (*domain.Clinic, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var clinic domain.Clinic
	err = attributevalue.UnmarshalMap(result.Item, &clinic)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal clinic: %w", err)
	}

	return &clinic, nil
}

func (r *DynamoDBClinicRepository) GetByCNPJ(ctx context.Context, cnpj string) (*domain.Clinic, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("cnpj-index"),
		KeyConditionExpression: aws.String("cnpj = :cnpj"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cnpj": &types.AttributeValueMemberS{Value: cnpj},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by CNPJ: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var clinic domain.Clinic
	err = attributevalue.UnmarshalMap(result.Items[0], &clinic)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal clinic: %w", err)
	}

	return &clinic, nil
}

func (r *DynamoDBClinicRepository) Update(ctx context.Context, clinic *domain.Clinic) error {
	clinic.UpdatedAt = time.Now()

	item, err := attributevalue.MarshalMap(clinic)
	if err != nil {
		return fmt.Errorf("failed to marshal clinic: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

func (r *DynamoDBClinicRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

func (r *DynamoDBClinicRepository) List(ctx context.Context, filters map[string]interface{}) ([]*domain.Clinic, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	}

	if len(filters) > 0 {
		filterExpression := ""
		expressionAttributeValues := make(map[string]types.AttributeValue)
		expressionAttributeNames := make(map[string]string)

		for key, value := range filters {
			if filterExpression != "" {
				filterExpression += " AND "
			}
			filterExpression += fmt.Sprintf("#%s = :%s", key, key)
			expressionAttributeNames[fmt.Sprintf("#%s", key)] = key

			switch v := value.(type) {
			case string:
				expressionAttributeValues[fmt.Sprintf(":%s", key)] = &types.AttributeValueMemberS{Value: v}
			}
		}

		input.FilterExpression = aws.String(filterExpression)
		input.ExpressionAttributeValues = expressionAttributeValues
		input.ExpressionAttributeNames = expressionAttributeNames
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan clinics: %w", err)
	}

	var clinics []*domain.Clinic
	for _, item := range result.Items {
		var clinic domain.Clinic
		err = attributevalue.UnmarshalMap(item, &clinic)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal clinic: %w", err)
		}
		clinics = append(clinics, &clinic)
	}

	return clinics, nil
}
