#!/usr/bin/env bash

docker build -t logzio_terraform_provider_build .
docker run -v `pwd`:/go/src/github.com/logzio/logzio_terraform_provider: logzio_terraform_provider_build