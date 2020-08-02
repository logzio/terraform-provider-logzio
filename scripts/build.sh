#!/usr/bin/env bash -xe
TERRAFORM_HOME=~/.terraform.d

#GO111MODULE=on go get -v -t ./...
GO111MODULE=on go build -o ./build/terraform-provider-logzio
cp ./build/terraform-provider-logzio ${TERRAFORM_HOME}/plugins/

go run utils/template.go
mkdir -p ${TERRAFORM_HOME}/metadata-repo/terraform/model/providers/
echo "logzio" >> ${TERRAFORM_HOME}/metadata-repo/terraform/model/providers.list
cp ./build/logzio.json ${TERRAFORM_HOME}/metadata-repo/terraform/model/providers/