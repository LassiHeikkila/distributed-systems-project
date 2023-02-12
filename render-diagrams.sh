#!/bin/bash

inputs=("flmnchll-architecture.drawio")

for input in "${inputs[@]}"; do
    drawio --export --format png "${input}"
done
