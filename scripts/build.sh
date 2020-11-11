#!/usr/bin/env bash
TERRAFORM_HOME=~/.terraform.d

#GO111MODULE=on go get -v -t ./...
GO111MODULE=on go build -o ./build/terraform-provider-logzio
cp ./build/terraform-provider-logzio ${TERRAFORM_HOME}/plugins/