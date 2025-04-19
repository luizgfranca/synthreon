#!/usr/bin/env bash
docker build -t pl -f control-panel/Dockerfile .
docker kill synthreon-dev
docker rm synthreon-dev
docker run --name synthreon-dev --env-file control-panel/demo.env --network=host -v "$PWD"/data:/data synthreon
