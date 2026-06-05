#!/bin/sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
APP_DIR="$ROOT_DIR/app"
PROJECT_FILE="$APP_DIR/PanchangamMobile.xcodeproj"
IOS_DESTINATION=${IOS_DESTINATION:-generic/platform=iOS Simulator}
WATCH_DESTINATION=${WATCH_DESTINATION:-generic/platform=watchOS Simulator}

fail() {
  printf 'build-simulators: %s\n' "$1" >&2
  exit 1
}

"$APP_DIR/scripts/check-native-tooling.sh"

cd "$APP_DIR"
xcodegen generate --spec project.yml

xcodebuild \
  -project "$PROJECT_FILE" \
  -scheme Panchangam \
  -configuration Debug \
  -destination "$IOS_DESTINATION" \
  CODE_SIGNING_ALLOWED=NO \
  CODE_SIGNING_REQUIRED=NO \
  build

xcodebuild \
  -project "$PROJECT_FILE" \
  -scheme PanchangamWatch \
  -configuration Debug \
  -destination "$WATCH_DESTINATION" \
  CODE_SIGNING_ALLOWED=NO \
  CODE_SIGNING_REQUIRED=NO \
  build
