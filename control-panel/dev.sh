#!/usr/bin/env bash

cd ../packages/js-core
yarn build
npm link
cd ../../control-panel

cd web
npm link platformlab-core
yarn build
cd .. # web

export ACCESS_TOKEN_SECRET_KEY=supersecret
export ROOT_EMAIL=test@test.com
export ROOT_PASSWORD=password
export DATABASE=test.db
export STATIC_FILES_DIR=web/dist

go run .