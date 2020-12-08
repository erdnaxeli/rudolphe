#!/usr/bin/env sh

apk add sqlite-static sqlite-dev
crystal build --static --release src/cli.cr
