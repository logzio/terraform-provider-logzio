#!/usr/bin/env bash
go build -o build/terraform-provider-logzio
pushd build
terraform init
terraform plan -out terraform.plan
terraform apply terraform.plan
popd