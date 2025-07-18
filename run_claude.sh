#!/bin/bash

# Usage: run_claude.sh [--risky] "<your prompt here>"

RISKY=""
PROMPT=""

# Parse arguments
if [[ "$1" == "--risky" ]]; then
  RISKY="--dangerously-skip-permissions"
  shift
fi

PROMPT="$*"

if [[ -z "$PROMPT" ]]; then
  echo "Usage: $0 [--risky] \"<your prompt>\""
  exit 1
fi

# Run Claude
claude --print $RISKY "$PROMPT" >> ~/claude_output.log 2>&1

