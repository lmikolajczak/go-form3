#!/bin/sh
# run-tests.sh

set -e

healthcheck="$FORM3_API_BASE_URL/v1/health"
runtests="go test ./... -cover -tags=integration -count=1 -v"

until $(curl --output /dev/null --silent --fail $healthcheck); do
    echo "Waiting for $healthcheck..."
    sleep 1
done

echo "Account API is up - executing: $runtests"
exec $runtests
