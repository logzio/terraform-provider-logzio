#!/usr/bin/env bash -xe
go build -o ./build/terraform-provider-logzio
cp ./build/terraform-provider-logzio ~/.terraform.d/plugins