#!/bin/bash

echo "🧪 Testing Dental Marketplace API"
echo "=================================="
echo ""

BASE_URL="http://localhost:8080"

# Test 1: Health Check
echo "✅ 1. Health Check"
curl -s $BASE_URL/health | jq .
echo ""
echo ""

# Test 2: Patient Login
echo "✅ 2. Patient Login"
PATIENT_LOGIN=$(curl -s -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "patient", "password": "password"}')

echo $PATIENT_LOGIN | jq .
PATIENT_TOKEN=$(echo $PATIENT_LOGIN | jq -r '.access_token')
echo ""
echo ""

# Test 3: Get Patient Scans
echo "✅ 3. Get Patient Scans"
curl -s $BASE_URL/api/patient/scans \
  -H "Authorization: Bearer $PATIENT_TOKEN" | jq .
echo ""
echo ""

# Test 4: Get Treatment Plans
echo "✅ 4. Get Treatment Plans"
curl -s $BASE_URL/api/patient/plans \
  -H "Authorization: Bearer $PATIENT_TOKEN" | jq .
echo ""
echo ""

# Test 5: Clinic Login
echo "✅ 5. Clinic Login"
CLINIC_LOGIN=$(curl -s -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "clinic1", "password": "password"}')

CLINIC_TOKEN=$(echo $CLINIC_LOGIN | jq -r '.access_token')
echo "Clinic logged in: $(echo $CLINIC_LOGIN | jq -r '.user.profile.name')"
echo ""
echo ""

# Test 6: Get Clinic Dashboard
echo "✅ 6. Get Clinic Dashboard"
curl -s $BASE_URL/api/clinic/dashboard \
  -H "Authorization: Bearer $CLINIC_TOKEN" | jq .
echo ""
echo ""

# Test 7: Regulator Login
echo "✅ 7. Regulator Login"
REGULATOR_LOGIN=$(curl -s -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "regulator", "password": "password"}')

REGULATOR_TOKEN=$(echo $REGULATOR_LOGIN | jq -r '.access_token')
echo "Regulator logged in"
echo ""
echo ""

# Test 8: Get Regulator Dashboard
echo "✅ 8. Get Regulator Dashboard"
curl -s $BASE_URL/api/regulator/dashboard \
  -H "Authorization: Bearer $REGULATOR_TOKEN" | jq .
echo ""
echo ""

echo "=================================="
echo "✅ All tests completed!"
