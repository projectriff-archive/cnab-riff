#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

version="`cat VERSION`-${BUILD_NUMBER}"
bucket=gs://projectriff/riff-cnab/builds

duffle init
duffle build .
duffle export riff -t
tar -xvzf riff-*.tgz
mv bundle.* riff-bundle-${version}.json

echo "$DOCKER_PASSWORD" | docker login --username $DOCKER_USERNAME --password-stdin
gcloud auth activate-service-account --key-file <(echo $GCLOUD_CLIENT_SECRET | base64 --decode)

docker push `duffle show riff | jq -r .invocationImages[0].image`
gsutil cp -a public-read -n riff-bundle-${version}.json ${bucket}/
