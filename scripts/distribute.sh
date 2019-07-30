#!/usr/bin/env bash -xe
if [ -z ${1} ]; then echo "provide OS as arg 1" && exit 1; fi
if [ -z ${2} ]; then echo "provide ARCH as arg 2" && exit 1; fi
if [ -z ${3} ]; then echo "provide VERSION as arg 3" && exit 1; fi
GO111MODULE=on go get -v -t ./...
GO111MODULE=on GOOS=${1} GOARCH=${2} go build -o ./build/${1}_${2}/terraform-provider-logzio_${3}
cp ./build/${1}_${2}/terraform-provider-logzio_${3} ./build/terraform-provider-logzio_${3}_${1}_${2}