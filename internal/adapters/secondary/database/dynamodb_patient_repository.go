package database

import (
	"context"
	"fmt"
	"go-dentist-api/internal/core/domain"
	"go-dentist-api/internal/core/ports"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type DynamoDBPatientRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBPatientRepository(client *dynamodb.Client, tableName string) ports.PatientRepository {
	return &DynamoDBPatientRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *DynamoDBPatientRepository) Create(ctx context.Context, patient *domain.Patient) error {
	item, err := attributevalue.MarshalMap(patient)
	if err != nil {
		return fmt.Errorf("failed to marshal patient: %w", err)
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

func (r *DynamoDBPatientRepository) GetByID(ctx context.Context, id string) (*domain.Patient, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if result.Item == nil {
		return nil, nil
	}

	var patient domain.Patient
	err = attributevalue.UnmarshalMap(result.Item, &patient)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal patient: %w", err)
	}

	return &patient, nil
}

func (r *DynamoDBPatientRepository) GetByCPF(ctx context.Context, cpf string) (*domain.Patient, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("cpf-index"),
		KeyConditionExpression: aws.String("CPF = :cpf"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":cpf": &types.AttributeValueMemberS{Value: cpf},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by CPF: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var patient domain.Patient
	err = attributevalue.UnmarshalMap(result.Items[0], &patient)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal patient: %w", err)
	}

	return &patient, nil
}

func (r *DynamoDBPatientRepository) Update(ctx context.Context, patient *domain.Patient) error {
	patient.UpdatedAt = time.Now()

	item, err := attributevalue.MarshalMap(patient)
	if err != nil {
		return fmt.Errorf("failed to marshal patient: %w", err)
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

func (r *DynamoDBPatientRepository) Delete(ctx context.Context, id string) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}

func (r *DynamoDBPatientRepository) GetByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("clinic-id-index"),
		KeyConditionExpression: aws.String("ClinicID = :clinic_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":clinic_id": &types.AttributeValueMemberS{Value: clinicID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query by clinic ID: %w", err)
	}

	var patients []*domain.Patient
	for _, item := range result.Items {
		var patient domain.Patient
		err = attributevalue.UnmarshalMap(item, &patient)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal patient: %w", err)
		}
		patients = append(patients, &patient)
	}

	return patients, nil
}

func (r *DynamoDBPatientRepository) GetActiveByClinicID(ctx context.Context, clinicID string) ([]*domain.Patient, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("clinic-id-index"),
		KeyConditionExpression: aws.String("ClinicID = :clinic_id"),
		FilterExpression:       aws.String("Active = :active"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":clinic_id": &types.AttributeValueMemberS{Value: clinicID},
			":active":    &types.AttributeValueMemberBOOL{Value: true},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query active patients by clinic ID: %w", err)
	}

	var patients []*domain.Patient
	for _, item := range result.Items {
		var patient domain.Patient
		err = attributevalue.UnmarshalMap(item, &patient)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal patient: %w", err)
		}
		patients = append(patients, &patient)
	}

	return patients, nil
}

func (r *DynamoDBPatientRepository) SearchByName(ctx context.Context, clinicID, name string) ([]*domain.Patient, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("clinic-id-index"),
		KeyConditionExpression: aws.String("ClinicID = :clinic_id"),
		FilterExpression:       aws.String("contains(#name, :name)"),
		ExpressionAttributeNames: map[string]string{
			"#name": "Name",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":clinic_id": &types.AttributeValueMemberS{Value: clinicID},
			":name":      &types.AttributeValueMemberS{Value: name},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search patients by name: %w", err)
	}

	var patients []*domain.Patient
	for _, item := range result.Items {
		var patient domain.Patient
		err = attributevalue.UnmarshalMap(item, &patient)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal patient: %w", err)
		}
		patients = append(patients, &patient)
	}

	return patients, nil
}

func (r *DynamoDBPatientRepository) AddProcedure(ctx context.Context, patientID string, procedure *domain.PatientProcedure) error {
	// Primeiro, buscar o paciente
	patient, err := r.GetByID(ctx, patientID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return fmt.Errorf("patient not found")
	}

	// Adicionar o procedimento
	procedure.ID = uuid.New().String()
	procedure.PatientID = patientID
	procedure.CreatedAt = time.Now()
	
	patient.Procedures = append(patient.Procedures, *procedure)
	patient.UpdatedAt = time.Now()

	// Atualizar o paciente
	return r.Update(ctx, patient)
}

func (r *DynamoDBPatientRepository) GetProceduresByPatientID(ctx context.Context, patientID string) ([]*domain.PatientProcedure, error) {
	patient, err := r.GetByID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return nil, fmt.Errorf("patient not found")
	}

	var procedures []*domain.PatientProcedure
	for i := range patient.Procedures {
		procedures = append(procedures, &patient.Procedures[i])
	}

	return procedures, nil
}

func (r *DynamoDBPatientRepository) GetProceduresByDentistID(ctx context.Context, dentistID string) ([]*domain.PatientProcedure, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("contains(procedures, :dentist_id)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":dentist_id": &types.AttributeValueMemberS{Value: dentistID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan procedures by dentist ID: %w", err)
	}

	var allProcedures []*domain.PatientProcedure
	for _, item := range result.Items {
		var patient domain.Patient
		err = attributevalue.UnmarshalMap(item, &patient)
		if err != nil {
			continue
		}

		for i := range patient.Procedures {
			if patient.Procedures[i].DentistID == dentistID {
				allProcedures = append(allProcedures, &patient.Procedures[i])
			}
		}
	}

	return allProcedures, nil
}

func (r *DynamoDBPatientRepository) UpdateProcedure(ctx context.Context, procedure *domain.PatientProcedure) error {
	patient, err := r.GetByID(ctx, procedure.PatientID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return fmt.Errorf("patient not found")
	}

	// Encontrar e atualizar o procedimento
	for i := range patient.Procedures {
		if patient.Procedures[i].ID == procedure.ID {
			patient.Procedures[i] = *procedure
			patient.UpdatedAt = time.Now()
			return r.Update(ctx, patient)
		}
	}

	return fmt.Errorf("procedure not found")
}

func (r *DynamoDBPatientRepository) AddExam(ctx context.Context, patientID string, exam *domain.Exam) error {
	patient, err := r.GetByID(ctx, patientID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return fmt.Errorf("patient not found")
	}

	exam.ID = uuid.New().String()
	exam.PatientID = patientID
	exam.CreatedAt = time.Now()
	exam.UpdatedAt = time.Now()
	
	patient.Exams = append(patient.Exams, *exam)
	patient.UpdatedAt = time.Now()

	return r.Update(ctx, patient)
}

func (r *DynamoDBPatientRepository) GetExamsByPatientID(ctx context.Context, patientID string) ([]*domain.Exam, error) {
	patient, err := r.GetByID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return nil, fmt.Errorf("patient not found")
	}

	var exams []*domain.Exam
	for i := range patient.Exams {
		exams = append(exams, &patient.Exams[i])
	}

	return exams, nil
}

func (r *DynamoDBPatientRepository) UpdateExam(ctx context.Context, exam *domain.Exam) error {
	patient, err := r.GetByID(ctx, exam.PatientID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return fmt.Errorf("patient not found")
	}

	// Encontrar e atualizar o exame
	for i := range patient.Exams {
		if patient.Exams[i].ID == exam.ID {
			exam.UpdatedAt = time.Now()
			patient.Exams[i] = *exam
			patient.UpdatedAt = time.Now()
			return r.Update(ctx, patient)
		}
	}

	return fmt.Errorf("exam not found")
}

func (r *DynamoDBPatientRepository) GetExamByID(ctx context.Context, examID string) (*domain.Exam, error) {
	// Como os exames estão dentro dos pacientes, precisamos fazer um scan
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String(r.tableName),
		FilterExpression: aws.String("contains(exams, :exam_id)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":exam_id": &types.AttributeValueMemberS{Value: examID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to scan for exam: %w", err)
	}

	for _, item := range result.Items {
		var patient domain.Patient
		err = attributevalue.UnmarshalMap(item, &patient)
		if err != nil {
			continue
		}

		for i := range patient.Exams {
			if patient.Exams[i].ID == examID {
				return &patient.Exams[i], nil
			}
		}
	}

	return nil, nil
}

func (r *DynamoDBPatientRepository) UpdateMedicalHistory(ctx context.Context, patientID string, medicalHistory *domain.MedicalHistory) error {
	patient, err := r.GetByID(ctx, patientID)
	if err != nil {
		return fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return fmt.Errorf("patient not found")
	}

	medicalHistory.LastUpdate = time.Now()
	patient.MedicalHistory = *medicalHistory
	patient.UpdatedAt = time.Now()

	return r.Update(ctx, patient)
}

func (r *DynamoDBPatientRepository) GetMedicalHistory(ctx context.Context, patientID string) (*domain.MedicalHistory, error) {
	patient, err := r.GetByID(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}
	if patient == nil {
		return nil, fmt.Errorf("patient not found")
	}

	return &patient.MedicalHistory, nil
}