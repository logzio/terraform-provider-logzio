#!/usr/bin/env bash -xe
go build -o ./build/terraform-provider-logzio
cp ./build/terraform-provider-logzio ~/.terraform.d/plugins/
go run utils/template.go
mkdir -p ~./terraform.d/metadata-repo/terraform/model/providers/
echo "logzio" >> ~./terraform.d/metadata-repo/terraform/model/providers.list
cp ./build/logzio.json ~./terraform.d/metadata-repo/terraform/model/providers/