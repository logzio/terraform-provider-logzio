#!/usr/bin/env bash -xe
go build -o terraform/terraform-provider-logzio
pushd terraform
export TF_LOG=DEBUG
terraform init
terraform destroy -var-file=./variables/${USER}.tfvars
popd