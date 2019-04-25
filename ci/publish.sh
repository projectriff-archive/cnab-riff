#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

version=`cat VERSION`
commit=$(git rev-parse HEAD)

gcloud auth activate-service-account --key-file <(echo $GCLOUD_CLIENT_SECRET | base64 --decode)

bucket=gs://projectriff/riff-cnab/releases

gsutil cp -a public-read -n riff-*{.tgz,.zip} ${bucket}/v${version}/
