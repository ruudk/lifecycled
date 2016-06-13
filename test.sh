#!/bin/bash
set -o errexit

cd $(dirname $0)

PACKAGES=$(go list ./... | grep -v /vendor/)

go vet $PACKAGES
if ! which fgt > /dev/null ; then
    echo "Please install fgt from https://github.com/GeertJohan/fgt."
    exit 1
fi

if ! which golint > /dev/null ; then
    echo "Please install golint from github.com/golang/lint/golint."
    exit 1
fi

echo "running golint"
for package in $PACKAGES
do
    fgt golint ${package}
done

# check go fmt
echo "running gofmt"
for package in $PACKAGES
do
    test -z "$(gofmt -s -l $GOPATH/src/$package/ | grep -v /vendor/ | tee /dev/stderr)"
done

echo "running coverage"
COVPROFILES=""
for package in $(go list -f '{{if len .TestGoFiles}}{{.ImportPath}}{{end}}' $PACKAGES)
do
    profile="$GOPATH/src/$package/.coverprofile"
    go test -race --coverprofile=$profile $package
    [ -s $profile ] && COVPROFILES="$COVPROFILES $profile"
done
cat $COVPROFILES > coverprofile.txt