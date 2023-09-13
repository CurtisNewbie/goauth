#!/bin/bash

# --------- remote
remote="alphaboi@curtisnewbie.com"
service="goauth"
remote_build_path="~/services/${service}/build/"
remote_config_path="~/services/${service}/config/"
os="linux"
arch="amd64"
build="${service}_build"
dockerfile="./Dockerfile_local"
# ---------

echo "Building Go app for platform $os/$arch to binary '$build'"
CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -o "$build"

ssh  "alphaboi@curtisnewbie.com" "rm -rv ${remote_build_path}*"

rsync -av -e ssh ./app-conf-prod.yml "${remote}:${remote_config_path}"
if [ ! $? -eq 0 ]; then
    exit -1
fi

rsync -av -e ssh "./$build" "${remote}:${remote_build_path}"
if [ ! $? -eq 0 ]; then
    exit -1
fi

rsync -av -e ssh "./$dockerfile" "${remote}:${remote_build_path}/Dockerfile"
if [ ! $? -eq 0 ]; then
    exit -1
fi

ssh  "alphaboi@curtisnewbie.com" "cd services; docker-compose up -d --build ${service}"

# cleanup
if [ -f "$build" ]; then
    rm "$build"
fi
