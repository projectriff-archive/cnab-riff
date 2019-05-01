#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

source $FATS_DIR/.configure.sh
source $FATS_DIR/functions/helpers.sh

# in cluster builds
for test in java java-boot node npm command; do
  path="$FATS_DIR/functions/uppercase/${test}"
  function_name=fats-cluster-uppercase-${test}
  image=$(fats_image_repo ${function_name})
  create_args="--git-repo $(git remote get-url origin) --git-revision $(git rev-parse HEAD) --sub-path functions/uppercase/${test}"
  input_data=cnab
  expected_data=CNAB

  run_function $path $function_name $image "$create_args" $input_data $expected_data
done

# local builds
if [ "$machine" != "MinGw" ]; then
  # TODO enable for windows once we have a linux docker daemon available
  for test in java java-boot node npm command; do
    path="$FATS_DIR/functions/uppercase/${test}"
    function_name=fats-local-uppercase-${test}
    image=$(fats_image_repo ${function_name})
    create_args="--local-path ."
    input_data=cnab
    expected_data=CNAB

    run_function $path $function_name $image "$create_args" $input_data $expected_data
  done
fi
