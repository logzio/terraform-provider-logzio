#!/usr/bin/env bash
PROVIDER_VERSION=v1.1.4
DOWNLOAD_URL=https://github.com/logzio/logzio_terraform_provider/releases/download/${PROVIDER_VERSION}

rm -f ~/.terraform.d/plugins/terraform-provider-logzio
mkdir -p  ~/.terraform.d/plugins/
wget ${DOWNLOAD_URL}/terraform-provider-logzio_${PROVIDER_VERSION}_darwin_amd64 -O ~/.terraform.d/plugins/terraform-provider-logzio
chmod +x ~/.terraform.d/plugins/terraform-provider-logzio