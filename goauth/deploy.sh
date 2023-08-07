#!/bin/bash

# --------- remote
app='goauth'
remote="alphaboi@curtisnewbie.com"
remote_build_path="~/services/${app}/build/"
remote_config_path="~/services/${app}/config/"
# ---------

ssh  "alphaboi@curtisnewbie.com" "rm -rv ${remote_build_path}*"

rsync -av -e ssh  --exclude='.git' \
    --exclude='storage' \
    --exclude='trash' \
    --exclude='.vscode' \
    --exclude='schema' \
    --exclude='LICENSE' \
    --exclude='README.md' \
    ./* "${remote}:${remote_build_path}"
if [ ! $? -eq 0 ]; then
    exit -1
fi

rsync -av -e ssh ./app-conf-prod.yml "${remote}:${remote_config_path}"
if [ ! $? -eq 0 ]; then
    exit -1
fi

ssh  "alphaboi@curtisnewbie.com" "cd services; docker-compose up -d --build ${app}"