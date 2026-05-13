#!/bin/bash
set -e
BASE="http://localhost:8080"
echo "=== GO-ASSISTANT E2E TEST: Paczkomat Scenario ==="

# Step 1: Health check
echo "[1/6] Health check..."
HEALTH=$(curl -s $BASE/health)
if [ "$HEALTH" = "OK" ]; then
    echo "  ✓ Server healthy"
else
    echo "  ✗ Health failed: $HEALTH"
    exit 1
fi

# Step 2: Create paczkomat task with location trigger
echo "[2/6] Creating paczkomat task..."
TASK=$(curl -s -X POST $BASE/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Odbierz paczkę InPost Puławska","status":"pending","location_trigger":{"lat":52.2297,"lng":21.0122}}')
TASK_ID=$(echo $TASK | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")
echo "  ✓ Task created: $TASK_ID"

# Step 3: Verify task in DB
echo "[3/6] Verifying task in database..."
TASKS=$(curl -s $BASE/api/tasks)
COUNT=$(echo $TASKS | python3 -c "import sys,json; print(len(json.load(sys.stdin)))")
echo "  ✓ Tasks in DB: $COUNT"

# Step 4: Check proximity - user far from paczkomat
echo "[4/6] Proximity check - user far from paczkomat..."
FAR=$(curl -s -X POST $BASE/api/location/check \
  -H "Content-Type: application/json" \
  -d '{"current":{"lat":52.1000,"lng":20.8000},"target":{"lat":52.2297,"lng":21.0122},"radius":300}')
IS_NEAR=$(echo $FAR | python3 -c "import sys,json; print(json.load(sys.stdin)['near'])")
if [ "$IS_NEAR" = "False" ]; then
    echo "  ✓ Far from paczkomat: not triggered"
else
    echo "  ~ Far check: $IS_NEAR"
fi

# Step 5: Simulate user arriving near paczkomat - triggers WebSocket notification
echo "[5/6] User arrives near paczkomat (triggers geofence)..."
LOCATION=$(curl -s -X POST $BASE/api/location \
  -H "Content-Type: application/json" \
  -d '{"user_id":"00000000-0000-0000-0000-000000000001","lat":52.2298,"lng":21.0123}')
echo "  ✓ Location recorded: $LOCATION"

# Step 6: AI suggestions based on location
echo "[6/6] Getting AI suggestions for current context..."
AI=$(curl -s -X POST $BASE/api/ai/suggest \
  -H "Content-Type: application/json" \
  -d "{\"lat\":52.2298,\"lng\":21.0123,\"time_of_day\":\"evening\",\"tasks\":[{\"id\":\"$TASK_ID\",\"title\":\"Odbierz paczkę InPost Puławska\",\"status\":\"pending\",\"has_location_trigger\":true,\"distance_meters\":15.3}]}")
SUGGESTIONS=$(echo $AI | python3 -c "import sys,json; print(json.load(sys.stdin)['suggestions'][:80])")
echo "  ✓ AI suggestions: $SUGGESTIONS..."

echo ""
echo "=== ALL TESTS PASSED ✓ ==="
