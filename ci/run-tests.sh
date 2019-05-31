#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

fats_repo="projectriff/fats"
source $FATS_DIR/.configure.sh
source $FATS_DIR/functions/helpers.sh
$FATS_DIR/install.sh kail

echo "Checking for ready ingress"
wait_for_ingress_ready 'istio-ingressgateway' 'istio-system'

# setup namespace
kubectl create namespace $NAMESPACE
fats_create_push_credentials $NAMESPACE

for test in command; do
  path="$FATS_DIR/functions/uppercase/${test}"
  function_name=fats-cluster-uppercase-${test}
  image=$(fats_image_repo ${function_name})
  create_args="--git-repo https://github.com/${fats_repo}.git --git-revision ${FATS_REFSPEC} --sub-path functions/uppercase/${test}"
  input_data=cnab
  expected_data=CNAB

  run_function $path $function_name $image "$create_args" $input_data $expected_data
done
