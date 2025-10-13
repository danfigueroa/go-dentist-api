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

type DynamoDBDentistRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBDentistRepository(client *dynamodb.Client, tableName string) repositories.DentistRepository {
	return &DynamoDBDentistRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBDentistRepository) Create(ctx context.Context, dentist *domain.Dentist) error {
	item, err := attributevalue.MarshalMap(dentist)
	if err != nil {
		return fmt.Errorf("failed to marshal dentist: %w", err)
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

func (r *DynamoDBDentistRepository) GetByID(ctx context.Context, id string) (*domain.Dentist, error) {
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

	var dentist domain.Dentist
	err = attributevalue.UnmarshalMap(result.Item, &dentist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dentist: %w", err)
	}

	return &dentist, nil
}

func (r *DynamoDBDentistRepository) GetByCRO(ctx context.Context, cro string) (*domain.Dentist, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("cro-index"),
		KeyConditionExpression: aws.String("cro = :cro"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cro": &types.AttributeValueMemberS{Value: cro},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by CRO: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var dentist domain.Dentist
	err = attributevalue.UnmarshalMap(result.Items[0], &dentist)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dentist: %w", err)
	}

	return &dentist, nil
}

func (r *DynamoDBDentistRepository) GetByName(ctx context.Context, name string) ([]*domain.Dentist, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("contains(#name, :name)"),
		ExpressionAttributeNames: map[string]string{
			"#name": "name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":name": &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan by name: %w", err)
	}

	var dentists []*domain.Dentist
	for _, item := range result.Items {
		var dentist domain.Dentist
		err = attributevalue.UnmarshalMap(item, &dentist)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal dentist: %w", err)
		}
		dentists = append(dentists, &dentist)
	}

	return dentists, nil
}

func (r *DynamoDBDentistRepository) GetByClinicID(ctx context.Context, clinicID string) ([]*domain.Dentist, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("clinic-id-index"),
		KeyConditionExpression: aws.String("clinic_id = :clinic_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":clinic_id": &types.AttributeValueMemberS{Value: clinicID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by clinic ID: %w", err)
	}

	var dentists []*domain.Dentist
	for _, item := range result.Items {
		var dentist domain.Dentist
		err = attributevalue.UnmarshalMap(item, &dentist)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal dentist: %w", err)
		}
		dentists = append(dentists, &dentist)
	}

	return dentists, nil
}

func (r *DynamoDBDentistRepository) Update(ctx context.Context, dentist *domain.Dentist) error {
	dentist.UpdatedAt = time.Now()

	item, err := attributevalue.MarshalMap(dentist)
	if err != nil {
		return fmt.Errorf("failed to marshal dentist: %w", err)
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

func (r *DynamoDBDentistRepository) Delete(ctx context.Context, id string) error {
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

func (r *DynamoDBDentistRepository) List(ctx context.Context, filters map[string]interface{}) ([]*domain.Dentist, error) {
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
			case bool:
				expressionAttributeValues[fmt.Sprintf(":%s", key)] = &types.AttributeValueMemberBOOL{Value: v}
			}
		}

		input.FilterExpression = aws.String(filterExpression)
		input.ExpressionAttributeValues = expressionAttributeValues
		input.ExpressionAttributeNames = expressionAttributeNames
	}

	result, err := r.client.Scan(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan dentists: %w", err)
	}

	var dentists []*domain.Dentist
	for _, item := range result.Items {
		var dentist domain.Dentist
		err = attributevalue.UnmarshalMap(item, &dentist)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal dentist: %w", err)
		}
		dentists = append(dentists, &dentist)
	}

	return dentists, nil
}
