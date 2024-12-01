#!/bin/bash

set -xe

docker build -t pl .
docker tag pl smith:5000/pl
docker push smith:5000/pl

ssh smith 'docker pull localhost:5000/pl && docker tag localhost:5000/pl pl && docker kill pldemo && docker rm pldemo && cd ~/platformlab && docker run --name pldemo --env-file demo.env --network=host -v "$PWD"/data:/data pl'