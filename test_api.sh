#!/bin/bash

echo "Testing Dental API Endpoints"
echo "============================="

BASE_URL="http://localhost:8080/api/v1"

echo "1. Testing GET /dentists (List all dentists)"
curl -X GET "$BASE_URL/dentists" -H "Content-Type: application/json"
echo -e "\n"

echo "2. Testing POST /dentists (Create dentist)"
curl -X POST "$BASE_URL/dentists" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Dr. João Silva",
    "cro": "CRO-SP-12345",
    "email": "joao.silva@email.com",
    "phone": "+55 11 99999-9999",
    "specialization": "Ortodontia",
    "clinic_id": "clinic-123"
  }'
echo -e "\n"

echo "3. Testing GET /dentists (List all dentists after creation)"
curl -X GET "$BASE_URL/dentists" -H "Content-Type: application/json"
echo -e "\n"

echo "4. Testing GET /dentists/search/cro/{cro} (Search by CRO)"
curl -X GET "$BASE_URL/dentists/search/cro/CRO-SP-12345" -H "Content-Type: application/json"
echo -e "\n"

echo "5. Testing GET /dentists/search/name/{name} (Search by name)"
curl -X GET "$BASE_URL/dentists/search/name/João" -H "Content-Type: application/json"
echo -e "\n"

echo "API Testing Complete!"