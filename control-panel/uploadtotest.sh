#!/bin/bash

set -xe

docker build -t synt .
docker tag pl smith:5000/synt
docker push smith:5000/synt

ssh smith 'docker pull localhost:5000/synt && docker tag localhost:5000/synt pl && docker kill syntdemo && docker rm syntdemo && cd ~/platformlab && docker run --name syntdemo --env-file demo.env --network=host -v "$PWD"/data:/data synt'
