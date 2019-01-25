#!/usr/bin/env bash
if [ -z ${1} ]; then echo "set ID as arg 1" && exit 1; fi
if [ -z ${LOGZIO_API_TOKEN} ]; then echo "set LOGZIO_API_TOKEN as env var" && exit 1; fi
ID=${1}
curl -sS -X "DELETE" -H "X-API-TOKEN: ${LOGZIO_API_TOKEN}"  -H "Content-Type: application/json" "https://api.logz.io/v1/endpoints/${ID}" | jq .