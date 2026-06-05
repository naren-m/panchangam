#!/bin/sh
# Verifies the live HTTP gateway path used by the phone, watch app, and complication.

set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
GRPC_PORT=${GRPC_PORT:-55051}
HTTP_PORT=${HTTP_PORT:-58080}
AT=${AT:-2026-06-02T12:00:00Z}
LAT=${LAT:-37.3382}
LNG=${LNG:--121.8863}
TZ=${TZ:-America/Los_Angeles}

TMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/panchangam-api-smoke.XXXXXX")
GRPC_PID=""
GATEWAY_PID=""

cleanup() {
	if [ -n "$GATEWAY_PID" ] && kill -0 "$GATEWAY_PID" 2>/dev/null; then
		kill "$GATEWAY_PID" 2>/dev/null || true
	fi
	if [ -n "$GRPC_PID" ] && kill -0 "$GRPC_PID" 2>/dev/null; then
		kill "$GRPC_PID" 2>/dev/null || true
	fi
	rm -rf "$TMP_DIR"
}
trap cleanup EXIT INT TERM

wait_for_url() {
	url=$1
	name=$2
	attempt=1
	while [ "$attempt" -le 30 ]; do
		if curl -fsS "$url" >/dev/null 2>&1; then
			return 0
		fi
		sleep 1
		attempt=$((attempt + 1))
	done

	echo "$name did not become ready at $url" >&2
	return 1
}

require_response_text() {
	text=$1
	if ! grep -Fq "$text" "$TMP_DIR/response.json"; then
		echo "Response missing expected text: $text" >&2
		echo "Response body:" >&2
		cat "$TMP_DIR/response.json" >&2
		return 1
	fi
}

require_header_text() {
	text=$1
	if ! grep -Fq "$text" "$TMP_DIR/headers.txt"; then
		echo "Response headers missing expected text: $text" >&2
		echo "Response headers:" >&2
		cat "$TMP_DIR/headers.txt" >&2
		return 1
	fi
}

validate_current_response() {
	response_file=$1
	before_epoch=$2
	after_epoch=$3

	python3 - "$response_file" "$before_epoch" "$after_epoch" <<'PY'
import datetime
import json
import sys

path, before_text, after_text = sys.argv[1:]
with open(path, "r", encoding="utf-8") as handle:
    response = json.load(handle)


def parse_rfc3339(value):
    if value.endswith("Z"):
        value = value[:-1] + "+00:00"
    return datetime.datetime.fromisoformat(value).timestamp()


generated_at = parse_rfc3339(response["generated_at"])
start_time = parse_rfc3339(response["tithi"]["start_time"])
end_time = parse_rfc3339(response["tithi"]["end_time"])
next_refresh_at = parse_rfc3339(response["next_refresh_at"])
before = float(before_text)
after = float(after_text)

if generated_at < before - 2 or generated_at > after + 2:
    raise SystemExit("generated_at should be within the current request window")
if generated_at < start_time or generated_at >= end_time:
    raise SystemExit("generated_at should be within the returned tithi window")
if next_refresh_at <= generated_at:
    raise SystemExit("next_refresh_at should be after generated_at")
PY
}

cd "$ROOT_DIR"

go build -o "$TMP_DIR/panchangam-server" ./cmd/server
go build -o "$TMP_DIR/panchangam-gateway" ./cmd/gateway

"$TMP_DIR/panchangam-server" --grpc-port="$GRPC_PORT" --log-level=warn >"$TMP_DIR/grpc.log" 2>&1 &
GRPC_PID=$!

sleep 1
if ! kill -0 "$GRPC_PID" 2>/dev/null; then
	echo "gRPC server exited during startup" >&2
	cat "$TMP_DIR/grpc.log" >&2
	exit 1
fi

ENABLE_CACHE=false "$TMP_DIR/panchangam-gateway" --grpc-endpoint="localhost:$GRPC_PORT" --http-port="$HTTP_PORT" --log-level=warn >"$TMP_DIR/gateway.log" 2>&1 &
GATEWAY_PID=$!

sleep 1
if ! kill -0 "$GATEWAY_PID" 2>/dev/null; then
	echo "HTTP gateway exited during startup" >&2
	cat "$TMP_DIR/gateway.log" >&2
	exit 1
fi

wait_for_url "http://localhost:$HTTP_PORT/api/v1/health" "HTTP gateway"

curl -fsS \
	-D "$TMP_DIR/headers.txt" \
	"http://localhost:$HTTP_PORT/api/v1/tithi/current?at=$AT&lat=$LAT&lng=$LNG&tz=$TZ&region=California&method=Drik&locale=en&calendar_system=Purnimanta" \
	-o "$TMP_DIR/response.json"

require_response_text '"date":"2026-06-02"'
require_response_text '"tithi":'
require_response_text '"pancha_anga":'
require_response_text '"day":'
require_response_text '"calculation":'
require_response_text '"timezone":"America/Los_Angeles"'
require_response_text "\"generated_at\":\"$AT\""
require_response_text '"next_refresh_at":'
require_header_text "Cache-Control: public, max-age=300"

CURRENT_BEFORE=$(python3 -c 'import time; print(time.time())')
curl -fsS \
	"http://localhost:$HTTP_PORT/api/v1/tithi/current?lat=$LAT&lng=$LNG&tz=$TZ&region=California&method=Drik&locale=en&calendar_system=Purnimanta" \
	-o "$TMP_DIR/current-response.json"
CURRENT_AFTER=$(python3 -c 'import time; print(time.time())')

grep -Fq '"generated_at":' "$TMP_DIR/current-response.json" || {
	echo "Current response missing generated_at" >&2
	cat "$TMP_DIR/current-response.json" >&2
	exit 1
}
validate_current_response "$TMP_DIR/current-response.json" "$CURRENT_BEFORE" "$CURRENT_AFTER"

echo "current tithi API smoke passed"
