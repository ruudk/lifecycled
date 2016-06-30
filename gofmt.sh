#!/bin/bash
set -o errexit

cd $(dirname $0)

PACKAGES=$(go list ./... | grep -v /vendor/)

echo "running gofmt"
for package in $PACKAGES
do
    gofmt -w $GOPATH/src/$package/
done
