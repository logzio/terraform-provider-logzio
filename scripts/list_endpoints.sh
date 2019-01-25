#!/usr/bin/env bash
if [ -z ${LOGZIO_API_TOKEN} ]; then echo "set LOGZIO_API_TOKEN as env var" && exit 1; fi
curl -sS -X "GET" -H "X-API-TOKEN: ${LOGZIO_API_TOKEN}"  -H "Content-Type: application/json" "https://api.logz.io/v1/endpoints" | jq .