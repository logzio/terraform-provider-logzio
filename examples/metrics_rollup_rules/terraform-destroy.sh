#!/bin/bash

# Example script to destroy metrics rollup rules Terraform configuration

# Set required environment variables
export TF_VAR_api_token="your-logzio-api-token"
export TF_VAR_region="us"  # or your region
export TF_VAR_account_id="123456"  # your account ID

# Destroy resources
terraform destroy -auto-approve 