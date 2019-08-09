#!/usr/bin/env bash

docker build -t logzio_terraform_provider:build .
docker run -v ~/go/src/github.com/jonboydell/logzio_terraform_provider:/go/src/github.com/jonboydell/logzio_terraform_provider logzio_terraform_provider:build go get -v ./...