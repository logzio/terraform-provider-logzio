#!/usr/bin/env bash
PROVIDER_VERSION=v1.1.3
DOWNLOAD_URL=https://github.com/jonboydell/logzio_terraform_provider/releases/download/${PROVIDER_VERSION}

rm ~/.terraform.d/plugins/terraform-provider-logzio
wget ${DOWNLOAD_URL}/terraform-provider-logzio_${PROVIDER_VERSION}_darwin_amd64 -O $HOME/.terraform.d/plugins/terraform-provider-logzio
chmod +x ~/.terraform.d/plugins/terraform-provider-logzio