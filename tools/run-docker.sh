#!/usr/bin/env bash
docker build -t pl -f control-panel/Dockerfile .
docker kill platformlab-dev
docker rm platformlab-dev
docker run --name platformlab-dev --env-file control-panel/demo.env --network=host -v "$PWD"/data:/data pl