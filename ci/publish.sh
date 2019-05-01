#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [ "$#" -ne 1 ]; then
  echo "expected to receive one of (stage/snapshot/release) as argument"
  exit 1
fi

if [ "$1" == "stage" ]; then
  version="`cat VERSION`-${BUILD_NUMBER}"
  bucket=gs://projectriff/riff-cnab/builds
elif [ "$1" == "snapshot" ]; then
  version="`cat VERSION`-${BUILD_NUMBER}"
  bucket=gs://projectriff/riff-cnab/snapshots
elif [ "$1" == "release" ]; then
  version=`cat VERSION`
  bucket=gs://projectriff/riff-cnab/releases
else
  echo "unknown publish argument"
  exit 1
fi


duffle export riff -t
tar -xvzf riff-*.tgz
mv bundle.* riff-bundle-${version}.json

gcloud auth activate-service-account --key-file <(echo $GCLOUD_CLIENT_SECRET | base64 --decode)

gsutil cp -a public-read -n riff-bundle-${version}.json ${bucket}/
