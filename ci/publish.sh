#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

if [ "$#" -ne 1 ]; then
  echo "expected to receive one of (stage/snapshot/release) as argument"
  exit 1
fi

build_bucket=gs://projectriff/riff-cnab/builds
build_file="riff-bundle-`cat VERSION`-${BUILD_NUMBER}.json"

if [ "$1" == "snapshot" ]; then
  version="`cat VERSION`-${BUILD_NUMBER}"
  bucket=gs://projectriff/riff-cnab/snapshots
elif [ "$1" == "release" ]; then
  version=`cat VERSION`
  bucket=gs://projectriff/riff-cnab/releases
else
  echo "unknown publish argument"
  exit 1
fi

dest_file="riff-bundle-${version}.json"

gsutil cp -a public-read -n ${build_bucket}/${build_file} ${bucket}/${dest_file}
gsutil cp -a public-read ${build_bucket}/${build_file} gs://projectriff/riff-cnab/snapshots/riff-bundle-`cat VERSION`.json
gsutil cp -a public-read ${build_bucket}/${build_file} gs://projectriff/riff-cnab/snapshots/riff-bundle-latest.json
