#!/usr/bin/env bash
terraform init
terraform plan -var-file=./variables/${USER}.tfvars -out terraform.plan
terraform apply terraform.plan