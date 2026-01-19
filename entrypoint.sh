#!/bin/sh

echo "Starting PR Merge Action..."

# Set environment variables from inputs
export REPO_OWNER="$INPUT_REPO_OWNER"
export REPO="$INPUT_REPO"
export GITHUB_ACCESS_TOKEN="$INPUT_GITHUB_ACCESS_TOKEN"

# Run the Go binary
/app/pr-merger
