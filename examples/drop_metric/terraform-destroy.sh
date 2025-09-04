#!/usr/bin/env bash
export TF_LOG=DEBUG
TF_VAR_api_token=${LOGZIO_API_TOKEN} TF_VAR_account_id=${LOGZIO_ACCOUNT_ID} terraform destroy 
