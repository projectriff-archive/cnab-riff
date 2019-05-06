#!/bin/bash

set -x
set -o errexit
set -o nounset
set -o pipefail

fats_repo="projectriff/fats"
source $FATS_DIR/.configure.sh
source $FATS_DIR/functions/helpers.sh
$FATS_DIR/install.sh kail

echo "Checking for ready pods"
wait_pod_selector_ready 'app=controller' 'knative-serving'
wait_pod_selector_ready 'app=webhook' 'knative-serving'
wait_pod_selector_ready 'app=build-controller' 'knative-build'
wait_pod_selector_ready 'app=build-webhook' 'knative-build'
echo "Checking for ready ingress"
wait_for_ingress_ready 'istio-ingressgateway' 'istio-system'

# setup namespace
kubectl create namespace $NAMESPACE
fats_create_push_credentials $NAMESPACE
riff namespace init $NAMESPACE $NAMESPACE_INIT_FLAGS

# in cluster builds
#for test in java java-boot node npm command; do
for test in java; do
  path="$FATS_DIR/functions/uppercase/${test}"
  function_name=fats-cluster-uppercase-${test}
  image=$(fats_image_repo ${function_name})
  create_args="--git-repo https://github.com/${fats_repo}.git --git-revision ${FATS_REFSPEC} --sub-path functions/uppercase/${test}"
  input_data=cnab
  expected_data=CNAB

  run_function $path $function_name $image "$create_args" $input_data $expected_data
done

# local builds
# if [ "$machine" != "MinGw" ]; then
#   # TODO enable for windows once we have a linux docker daemon available
#   for test in java java-boot node npm command; do
#     path="$FATS_DIR/functions/uppercase/${test}"
#     function_name=fats-local-uppercase-${test}
#     image=$(fats_image_repo ${function_name})
#     create_args="--local-path ."
#     input_data=cnab
#     expected_data=CNAB

#     run_function $path $function_name $image "$create_args" $input_data $expected_data
#   done
# fi
