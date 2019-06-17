#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

version=`cat VERSION`
commit=$(git rev-parse HEAD)
bucket=gs://projectriff/riff-cnab

staged_bundle="${bucket}/builds/riff-bundle-${version}-${commit}.json"

if echo "${version}" | grep -iq -e '-snapshot$'; then
  # release if the version is not a '-snapshot'
  gsutil cp -a public-read -n ${staged_bundle} ${bucket}/releases/riff-bundle-${version}.json
fi
gsutil cp -a public-read -n ${staged_bundle} ${bucket}/snapshots/riff-bundle-${version}-${commit}.json
gsutil cp -a public-read    ${staged_bundle} ${bucket}/snapshots/riff-bundle-${version}.json
gsutil cp -a public-read    ${staged_bundle} ${bucket}/snapshots/riff-bundle-latest.json
