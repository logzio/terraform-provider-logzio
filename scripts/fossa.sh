#!/bin/bash
if [ -z ${1} ]; then echo "provide fossa key as arg" && exit 1; fi
FOSSA_API_KEY=${1} fossa