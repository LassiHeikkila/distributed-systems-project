#!/bin/bash

version=${1:-"dev"}

set -e 

echo "Building Docker images in content/..."
cd content && ./build-images.sh "${version}" && cd -

echo "Building Docker images in account/..."
cd account && ./build-images.sh "${version}" && cd -

echo "Building Docker images in room/..."
cd room && ./build-images.sh "${version}" && cd -

# echo "Building frontend"
# cd frontend && npm run build && cd -