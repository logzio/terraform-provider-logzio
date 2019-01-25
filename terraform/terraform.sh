#!/usr/bin/env bash
export TF_LOG=DEBUG
terraform init
terraform plan -var-file=./${USER}.tfvars -out terraform.plan
terraform apply terraform.plan