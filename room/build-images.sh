#!/bin/bash

version=${1:-"dev"}

set -e

docker build -t "flmnchll/room-service:${version}" --file Dockerfile.roomservice .
