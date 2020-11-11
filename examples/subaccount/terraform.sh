#!/usr/bin/env bash
export TF_LOG=DEBUG
terraform init
TF_VAR_api_token=${LOGZIO_API_TOKEN} TF_VAR_account_id=${LOGZIO_ACCOUNT_ID} terraform plan -out terraform.plan
terraform apply terraform.plan