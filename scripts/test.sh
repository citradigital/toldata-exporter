#!/bin/sh

apk add --update curl git gcc musl-dev
su -l nobody
cd /src/
export MIGRATION_PATH=/migrations/test
env
go test
