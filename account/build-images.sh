#!/bin/bash

version=${1:-"dev"}

set -e

docker build -t "flmnchll/account-service:${version}" --file Dockerfile.accountservice .
