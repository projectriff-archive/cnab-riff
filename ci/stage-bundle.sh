#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

version=`cat VERSION`
commit=$(git rev-parse HEAD)
bucket=gs://projectriff/riff-cnab

duffle export riff -t
tar -xvzf riff-*.tgz

echo "$DOCKER_PASSWORD" | docker login --username $DOCKER_USERNAME --password-stdin
gcloud auth activate-service-account --key-file <(echo $GCLOUD_CLIENT_SECRET | base64 --decode)

docker push `duffle show riff | jq -r .invocationImages[0].image`
gsutil cp -a public-read -n bundle.* ${bucket}/builds/riff-bundle-${version}-${commit}.json
