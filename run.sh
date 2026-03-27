#!/usr/bin/env bash

set -e
set -u
set -o pipefail

load_env_file() {
    local env_file="$1"
    if [[ -f "$env_file" ]]; then
        set -a
        source "$env_file"
        set +a
    else
        echo "Error: $env_file not found" >&2
        return 1
    fi
}


case "${1:-}" in
    "test")
        cd function
        go clean -testcache
        go test ./...
        ;;
    "run")
        load_env_file ".env"
        cd function
        go run .
        ;;
    "trigger")
        load_env_file ".env"
        curl -H "X-API-KEY: $GREENMO_API_KEY" \
            "$GREENMO_API_URL?lon1=12.517685&lat1=55.739892&lon2=12.526059&lat2=55.734577&chargers=true&cars=true&desiredFuelLevel=60"
        ;;
    *) 
        echo "Usage: $0 {test|run|trigger}"
        exit 1 
        ;;
esac