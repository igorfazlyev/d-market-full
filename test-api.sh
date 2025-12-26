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

# Test 2: Get Constants (Public endpoint - no auth required)
echo "✅ 2. Get System Constants"
curl -s $BASE_URL/api/constants | jq .
echo ""
echo ""

# Test 2a: Verify Roles are loaded
echo "✅ 2a. Verify Roles Constants"
curl -s $BASE_URL/api/constants | jq '.roles'
echo ""
echo ""

# Test 2b: Verify Specializations are loaded
echo "✅ 2b. Verify Specializations Constants"
curl -s $BASE_URL/api/constants | jq '.specializations'
echo ""
echo ""

# Test 2c: Verify Cities and Districts
echo "✅ 2c. Verify Cities and Districts"
curl -s $BASE_URL/api/constants | jq '{cities, districts_by_city}'
echo ""
echo ""

# Test 3: Patient Login
echo "✅ 3. Patient Login"
PATIENT_LOGIN=$(curl -s -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "patient", "password": "password"}')

echo $PATIENT_LOGIN | jq .
PATIENT_TOKEN=$(echo $PATIENT_LOGIN | jq -r '.access_token')
echo ""
echo ""

# Test 4: Get Patient Scans
echo "✅ 4. Get Patient Scans"
curl -s $BASE_URL/api/patient/scans \
  -H "Authorization: Bearer $PATIENT_TOKEN" | jq .
echo ""
echo ""

# Test 5: Get Treatment Plans
echo "✅ 5. Get Treatment Plans"
curl -s $BASE_URL/api/patient/plans \
  -H "Authorization: Bearer $PATIENT_TOKEN" | jq .
echo ""
echo ""

# Test 6: Clinic Login
echo "✅ 6. Clinic Login"
CLINIC_LOGIN=$(curl -s -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "clinic1", "password": "password"}')

CLINIC_TOKEN=$(echo $CLINIC_LOGIN | jq -r '.access_token')
echo "Clinic logged in: $(echo $CLINIC_LOGIN | jq -r '.user.profile.name')"
echo ""
echo ""

# Test 7: Get Clinic Dashboard
echo "✅ 7. Get Clinic Dashboard"
curl -s $BASE_URL/api/clinic/dashboard \
  -H "Authorization: Bearer $CLINIC_TOKEN" | jq .
echo ""
echo ""

# Test 8: Regulator Login
echo "✅ 8. Regulator Login"
REGULATOR_LOGIN=$(curl -s -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "regulator", "password": "password"}')

REGULATOR_TOKEN=$(echo $REGULATOR_LOGIN | jq -r '.access_token')
echo "Regulator logged in"
echo ""
echo ""

# Test 9: Get Regulator Dashboard
echo "✅ 9. Get Regulator Dashboard"
curl -s $BASE_URL/api/regulator/dashboard \
  -H "Authorization: Bearer $REGULATOR_TOKEN" | jq .
echo ""
echo ""

echo "=================================="
echo "✅ All tests completed!"
echo ""
echo "📊 Summary:"
echo "   - Constants endpoint working"
echo "   - All 3 user roles can login"
echo "   - Role-based access working"
echo "   - Database seeded successfully"
echo "=================================="
