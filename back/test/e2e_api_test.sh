#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:3434}"
PASS=0
FAIL=0
TOTAL=0

GREEN='\033[0;32m'
RED='\033[0;31m'
CYAN='\033[0;36m'
NC='\033[0m'

assert_status() {
  local test_name="$1" expected="$2" actual="$3" body="${4:-}"
  TOTAL=$((TOTAL + 1))
  if [ "$actual" -eq "$expected" ]; then
    PASS=$((PASS + 1))
    printf "${GREEN}PASS${NC} %s (HTTP %s)\n" "$test_name" "$actual"
  else
    FAIL=$((FAIL + 1))
    printf "${RED}FAIL${NC} %s — expected %s, got %s\n" "$test_name" "$expected" "$actual"
    [ -n "$body" ] && printf "      body: %s\n" "$body"
  fi
}

jq_field() {
  python3 -c "import sys,json; print(json.load(sys.stdin)$1)" <<< "$2"
}

# curl wrapper — sets BODY and CODE
api() {
  local resp
  resp=$(curl -s -w "\n%{http_code}" "$@")
  CODE=$(echo "$resp" | tail -1)
  BODY=$(echo "$resp" | sed '$d')
}

printf "\n${CYAN}=== NeoHome API E2E Tests ===${NC}\n"
printf "Target: %s\n\n" "$BASE_URL"

# ── Health ──
api "$BASE_URL/health"
assert_status "GET /health" 200 "$CODE"

# ── Register ──
TS=$(date +%s)
EMAIL="e2e-${TS}@test.com"
LOGIN="e2e-${TS}"

api -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"testtest1\",\"login\":\"${LOGIN}\",\"phone\":\"12345678\"}"
assert_status "POST /auth/register (new user)" 201 "$CODE" "$BODY"

# ── Duplicate register ──
api -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"testtest1\",\"login\":\"${LOGIN}dup\",\"phone\":\"12345678\"}"
assert_status "POST /auth/register (duplicate → 409)" 409 "$CODE"

# ── Login ──
api -X POST "$BASE_URL/api/v1/auth/login/email" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"testtest1\"}"
assert_status "POST /auth/login/email" 200 "$CODE" "$BODY"
TOKEN=$(jq_field "['accessToken']" "$BODY")

# ── Get current user ──
api "$BASE_URL/api/v1/users/me" -H "Authorization: Bearer $TOKEN"
assert_status "GET /users/me" 200 "$CODE" "$BODY"

# ── Unauthorized ──
api "$BASE_URL/api/v1/users/me"
assert_status "GET /users/me (no token → 401)" 401 "$CODE"

# ── Create device ──
api -X POST "$BASE_URL/api/v1/devices" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"deviceName":"E2E Sensor","deviceType":"sensor","roomName":"Lab","locationId":1,"locationName":"Office","status":"online"}'
assert_status "POST /devices (create)" 201 "$CODE" "$BODY"
DEVICE_ID=$(jq_field "['device']['deviceId']" "$BODY")
printf "      deviceId = %s\n" "$DEVICE_ID"

# ── List devices ──
api "$BASE_URL/api/v1/devices" -H "Authorization: Bearer $TOKEN"
assert_status "GET /devices (list)" 200 "$CODE" "$BODY"

# ── Ingest telemetry via REST ──
NOW_MS=$(python3 -c "import time; print(int(time.time()*1000))")
api -X POST "$BASE_URL/api/v1/telemetry" \
  -H "Content-Type: application/json" \
  -d "{\"deviceId\":${DEVICE_ID},\"recordedAt\":${NOW_MS},\"metricType\":\"temperature\",\"metricValue\":23.5,\"unit\":\"C\",\"roomName\":\"Lab\",\"locationName\":\"Office\",\"batteryLevel\":95,\"signalStrength\":-45}"
assert_status "POST /telemetry (REST ingest)" 201 "$CODE" "$BODY"

# ── Ingest telemetry via MQTT proxy ──
NOW_MS2=$(python3 -c "import time; print(int(time.time()*1000))")
api -X POST "$BASE_URL/api/v1/telemetry/mqtt" \
  -H "Content-Type: application/json" \
  -d "{\"topic\":\"neohome/devices/${DEVICE_ID}/telemetry\",\"payload\":{\"deviceId\":${DEVICE_ID},\"recordedAt\":${NOW_MS2},\"metricType\":\"humidity\",\"metricValue\":55.2,\"unit\":\"%\",\"roomName\":\"Lab\",\"locationName\":\"Office\",\"batteryLevel\":90,\"signalStrength\":-50}}"
assert_status "POST /telemetry/mqtt (MQTT proxy)" 201 "$CODE" "$BODY"

# ── Device telemetry history ──
api "$BASE_URL/api/v1/devices/${DEVICE_ID}/telemetry" -H "Authorization: Bearer $TOKEN"
assert_status "GET /devices/{id}/telemetry" 200 "$CODE" "$BODY"

# ── Latest telemetry ──
api "$BASE_URL/api/v1/devices/${DEVICE_ID}/latest" -H "Authorization: Bearer $TOKEN"
assert_status "GET /devices/{id}/latest" 200 "$CODE" "$BODY"

# ── Set thresholds ──
api -X PUT "$BASE_URL/api/v1/devices/${DEVICE_ID}/thresholds" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"thresholds":[{"metricType":"temperature","minValue":15,"maxValue":30,"severity":"critical"}]}'
assert_status "PUT /devices/{id}/thresholds" 200 "$CODE" "$BODY"

# ── Get thresholds ──
api "$BASE_URL/api/v1/devices/${DEVICE_ID}/thresholds" -H "Authorization: Bearer $TOKEN"
assert_status "GET /devices/{id}/thresholds" 200 "$CODE" "$BODY"

# ── Ingest above threshold → alert ──
NOW_MS3=$(python3 -c "import time; print(int(time.time()*1000))")
api -X POST "$BASE_URL/api/v1/telemetry" \
  -H "Content-Type: application/json" \
  -d "{\"deviceId\":${DEVICE_ID},\"recordedAt\":${NOW_MS3},\"metricType\":\"temperature\",\"metricValue\":42.0,\"unit\":\"C\",\"roomName\":\"Lab\",\"locationName\":\"Office\",\"batteryLevel\":85,\"signalStrength\":-40}"
assert_status "POST /telemetry (above threshold)" 201 "$CODE" "$BODY"
ALERT_COUNT=$(jq_field "['alerts'].__len__()" "$BODY")
if [ "$ALERT_COUNT" -ge 1 ]; then
  assert_status "  → generated alert" 200 200
else
  assert_status "  → generated alert" 200 500
fi

# ── List alerts ──
api "$BASE_URL/api/v1/alerts" -H "Authorization: Bearer $TOKEN"
assert_status "GET /alerts" 200 "$CODE" "$BODY"

# ── Resolve alert ──
ALERT_ID=$(jq_field "['alerts'][0]['alertId']" "$BODY" 2>/dev/null || echo "")
if [ -n "$ALERT_ID" ] && [ "$ALERT_ID" != "None" ]; then
  api -X PUT "$BASE_URL/api/v1/alerts/${ALERT_ID}/resolve" \
    -H "Authorization: Bearer $TOKEN"
  assert_status "PUT /alerts/{id}/resolve" 200 "$CODE" "$BODY"
else
  printf "${CYAN}SKIP${NC} PUT /alerts/{id}/resolve — no alerts\n"
fi

# ── Update device ──
api -X PATCH "$BASE_URL/api/v1/devices/${DEVICE_ID}" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"deviceName":"E2E Sensor Updated"}'
assert_status "PATCH /devices/{id} (update)" 200 "$CODE" "$BODY"

# ── Update profile ──
api -X PUT "$BASE_URL/api/v1/users/me" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"email\":\"${EMAIL}\",\"password\":\"\",\"login\":\"${LOGIN}-upd\",\"phone\":\"87654321\"}"
assert_status "PUT /users/me (update profile)" 200 "$CODE" "$BODY"

# ── Summary ──
printf "\n${CYAN}=== Results ===${NC}\n"
printf "Total: %d | ${GREEN}Pass: %d${NC} | ${RED}Fail: %d${NC}\n\n" "$TOTAL" "$PASS" "$FAIL"
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
