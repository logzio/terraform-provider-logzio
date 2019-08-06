#!/usr/bin/env bash -xe
GO111MODULE=on go get -v -t ./...
GO111MODULE=on go build -o ./build/terraform-provider-logzio
cp ./build/terraform-provider-logzio ~/.terraform.d/plugins/
go run utils/template.go
mkdir -p ~./terraform.d/metadata-repo/terraform/model/providers/
echo "logzio" >> ~./terraform.d/metadata-repo/terraform/model/providers.list
cp ./build/logzio.json ~./terraform.d/metadata-repo/terraform/model/providers/