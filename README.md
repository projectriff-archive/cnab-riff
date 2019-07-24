[![Build Status](https://dev.azure.com/projectriff/projectriff/_apis/build/status/projectriff.cnab-riff?branchName=master)](https://dev.azure.com/projectriff/projectriff/_build/latest?definitionId=16&branchName=master)

# cnab-riff
A CNAB bundle for installing [riff](https://projectriff.io/)

## Getting started
A CNAB runtime is required to install this cnab bundle. Currently, [duffle](https://duffle.sh/) is the only supported runtime.

### Steps
The high level steps to install riff are:
1. Start a kubernetes cluster
1. Install [duffle](https://duffle.sh/) which is a CNAB runtime
1. Download and install the riff CNAB bundle into the cluster

### Parameters
There are two parameters that can be provided for this bundle
1. `node_port`: set this parameter to true if your kubernetes cluster does not have a LoadBalancer configured. This will ensure that kubernetes services can be accessed via NodePort rather than LoadBalancer.
    - default value: `false`
1. `log_level`: set the log level to one of `panic|fatal|error|warn|info|debug|trace`
    - default value: `info`

### Detailed Instructions
Once you have a running kubernetes cluster follow the steps below

1. Install duffle by download the latest duffle release for your operating system from duffle's [release page](https://github.com/deislabs/duffle/releases), put it somewhere on your path, and make it executable, for example, `chmod +x /usr/local/bin/duffle`.
1. Set environment variables required by duffle k8s driver and create the service account
    ```
    export SERVICE_ACCOUNT=duffle-runtime
    export KUBE_NAMESPACE=duffle
    kubectl create namespace $KUBE_NAMESPACE
    kubectl create serviceaccount "${SERVICE_ACCOUNT}" -n "${KUBE_NAMESPACE}"
    kubectl create clusterrolebinding "${SERVICE_ACCOUNT}-cluster-admin" --clusterrole cluster-admin --serviceaccount "${KUBE_NAMESPACE}:${SERVICE_ACCOUNT}"
    ```
1. install riff

    Append `-s node_port=true` as shown below for clusters that do not support LoadBalancer services, like Minikube.
    ```
    duffle install myriff https://storage.googleapis.com/projectriff/riff-cnab/snapshots/riff-bundle-latest.json --bundle-is-file -s node_port=true -d k8s
    ```
    where `-d k8s` uses the duffle kubernetes driver to run the installer image in kubernetes cluster
1. You should now be able to see riff components installed on your kubernetes cluster:
    ```
    kubectl get pods --all-namespaces
    ```

## Uninstall
To uninstall, set the SERVICE_ACCOUNT and KUBE_NAMESPACE environment variables as above and use the command:
```
duffle uninstall myriff -d k8s
```

## Developing
Modify the `kab-manifest.yaml` file to update riff components. Then run `make bundle` which will:
- generate `duffle.json`
- generate `cnab/app/kab/manifest.yaml`
- build the cnab bundle using `duffle`
