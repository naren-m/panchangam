#!/bin/sh
set -eu

fail() {
  printf 'check-native-tooling: %s\n' "$1" >&2
  exit 1
}

active_developer_dir=$(xcode-select -p 2>/dev/null || true)
xcode_developer_dir="/Applications/Xcode.app/Contents/Developer"

command -v xcodegen >/dev/null 2>&1 || fail "xcodegen is required"
command -v xcodebuild >/dev/null 2>&1 || fail "xcodebuild is required"

if ! xcodebuild -version >/dev/null 2>&1; then
  if [ ! -d "$xcode_developer_dir" ]; then
    fail "full Xcode is required; install Xcode, then run: sudo xcode-select -s $xcode_developer_dir"
  fi
  if [ -n "$active_developer_dir" ]; then
    fail "full Xcode is required; active developer directory is $active_developer_dir; run: sudo xcode-select -s $xcode_developer_dir"
  fi
  fail "full Xcode is required; no active developer directory is selected; run: sudo xcode-select -s $xcode_developer_dir"
fi

xcrun --find simctl >/dev/null 2>&1 || fail "simctl is required from full Xcode"
xcrun --show-sdk-path --sdk iphonesimulator >/dev/null 2>&1 || fail "iOS simulator SDK is required"
xcrun --show-sdk-path --sdk watchsimulator >/dev/null 2>&1 || fail "watchOS simulator SDK is required"

printf 'check-native-tooling: native simulator tooling available\n'
