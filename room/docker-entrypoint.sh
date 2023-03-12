#!/bin/bash

# check that all needed env vars are set
: "${PEERJS_ADDR:?"PEERJS_ADDR is unset"}"
: "${ACCOUNT_SERVICE_ADDR:?"ACCOUNT_SERVICE_ADDR is unset"}"

/app/room-service \
    -peerjs "${PEERJS_ADDR}" \
    -accountServiceAddr "${ACCOUNT_SERVICE_ADDR}" 
