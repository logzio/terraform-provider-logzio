#!/usr/bin/env bash -xe
go build -o terraform/terraform-provider-logzio
pushd terraform
export TF_LOG=DEBUG
terraform init
terraform plan -var-file=./variables/${USER}.tfvars -out terraform.plan
terraform apply terraform.plan
popd