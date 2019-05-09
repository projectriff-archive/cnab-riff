#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

version="`cat VERSION`-${BUILD_NUMBER}"
bucket=gs://projectriff/riff-cnab/builds

duffle export riff -t
tar -xvzf riff-*.tgz
mv bundle.* riff-bundle-${version}.json

gcloud auth activate-service-account --key-file <(echo $GCLOUD_CLIENT_SECRET | base64 --decode)

gsutil cp -a public-read -n riff-bundle-${version}.json ${bucket}/
