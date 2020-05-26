#!/bin/bash
# usage: ./deploy.sh
# Pushing and releasing are one and the same thing because SLA is low enough that I don't mind occasional mistakes in exchange for simplicity

set -eu

if [ -z "$(git status --porcelain)" ]; then
    echo "Tagging and deploying..."
else
    echo "Commit any changes before deploying"
    exit 1
fi

tag="$(date +"%Y-%m-%d-%H%M")-$(git rev-parse --short HEAD)"
git tag -a "${tag}" --message "Tagging ${tag} for deployment"
git push
terraform apply -var "tag=${tag}"