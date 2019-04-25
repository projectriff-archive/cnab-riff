#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

version=`cat VERSION`

duffle export riff -t
ls -hal
tar -xvzf riff-*.tgz
mv bundle.* riff-bundle-${version}.json

echo "foo is: $FOO"

gcloud auth activate-service-account --key-file <(echo $GCLOUD_CLIENT_SECRET | base64 --decode)

bucket=gs://projectriff/riff-cnab/releases

gsutil cp -a public-read -n riff-bundle-${version}.json ${bucket}/
