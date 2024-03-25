#!/usr/bin/env bash

go mod tidy

./scripts/build.sh
./build/server
