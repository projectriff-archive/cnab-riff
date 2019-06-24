[![Build Status](https://dev.azure.com/projectriff/projectriff/_apis/build/status/projectriff.cnab-riff?branchName=master)](https://dev.azure.com/projectriff/projectriff/_build/latest?definitionId=16&branchName=master)

# cnab-riff
A CNAB bundle for installing [riff](https://projectriff.io/)

## Getting started
To install this cnab bundle, you can use any available cnab runtime. The following are installation instructions using [duffle](https://duffle.sh/).

### Steps
1. Download the latest duffle release for your operating system from duffle's [release page](https://github.com/deislabs/duffle/releases),
 put it somewhere on your path, and make it executable using, for example, `chmod +x /usr/local/bin/duffle`.
1. Start a kubernetes cluster that does not need resources other than kubeconfig for authentication (like minikube OR docker-for-desktop) by following the instructions in [getting started](https://projectriff.io/docs/getting-started/minikube/).
 
1. Download the latest snapshot of the bundle file:
    ```
    curl -O https://storage.googleapis.com/projectriff/riff-cnab/snapshots/riff-bundle-latest.json
    ```
1. Create a service account for installing riff
    ```
    export SERVICE_ACCOUNT=duffle-runtime
    export KUBE_NAMESPACE=kube-system
    kubectl create serviceaccount "${SERVICE_ACCOUNT}" -n "${KUBE_NAMESPACE}"
    kubectl create clusterrolebinding "${SERVICE_ACCOUNT}-cluster-admin" --clusterrole cluster-admin --serviceaccount "${KUBE_NAMESPACE}:${SERVICE_ACCOUNT}"
    ```
1. Install riff
    ```
    duffle install myriff riff-bundle-latest.json --bundle-is-file -s node_port=true -d k8s
    ```
    where `node_port=true` parameter changes all service types to NodePort from LoadBalancer 
    and `-d k8s` uses the duffle kubernetes driver to run the installer image in kubernetes cluster
1. You should now be able to see riff components installed on your kubernetes cluster:
    ```
    kubectl get pods --all-namespaces
    ```

## Uninstall
To uninstall, set the SERVICE_ACCOUNT and KUBE_NAMESPACE environment variables as above and use the below command:
```
duffle uninstall myriff -d k8s
```
