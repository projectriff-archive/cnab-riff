[![Build Status](https://dev.azure.com/projectriff/projectriff/_apis/build/status/projectriff.cnab-riff?branchName=master)](https://dev.azure.com/projectriff/projectriff/_build/latest?definitionId=16&branchName=master)

# cnab-riff
A CNAB bundle for installing [riff](https://projectriff.io/)

## Getting started
To install this cnab bundle, you can use any available cnab runtime. The following are installation instructions using [duffle](https://duffle.sh/).

### Steps
The high level steps to install riff are:
1. Start a kubernetes cluster
1. Install a CNAB runtime. We will use [duffle](https://duffle.sh/).
1. Create service account in order to grant additional permissions like creating CRDs over the default sercice account.
1. Use the duffle k8s driver to install riff. This will run riff bundle's invocationImage in k8s cluster using the service account created in the previous step. We will also set a couple of environment variables required by the k8s driver.

### Detailed Instructions
Follow the steps below for your choice of kubernetes cluster
#### Minikube
1. [Install minikube](https://kubernetes.io/docs/tasks/tools/install-minikube/)
1. Start the minikube cluster
    ```
    minikube start --memory=4096 --cpus=4 \
    --kubernetes-version=v1.14.0 \
    --vm-driver=hyperkit \
    --bootstrapper=kubeadm \
    --extra-config=apiserver.enable-admission-plugins="LimitRanger,NamespaceExists,NamespaceLifecycle,ResourceQuota,ServiceAccount,DefaultStorageClass,MutatingAdmissionWebhook"
    ```
1. Confirm that your kubectl context is pointing to the new cluster
    ```
    kubectl config current-context
    ```
1. Install duffle by download the latest duffle release for your operating system from duffle's [release page](https://github.com/deislabs/duffle/releases), put it somewhere on your path, and make it executable, for example, `chmod +x /usr/local/bin/duffle`.
1. Set environment variables required by duffle k8s driver and create the the service account
    ```
    export SERVICE_ACCOUNT=duffle-runtime
    export KUBE_NAMESPACE=duffle
    kubectl create namespace $KUBE_NAMESPACE
    kubectl create serviceaccount "${SERVICE_ACCOUNT}" -n "${KUBE_NAMESPACE}"
    kubectl create clusterrolebinding "${SERVICE_ACCOUNT}-cluster-admin" --clusterrole cluster-admin --serviceaccount "${KUBE_NAMESPACE}:${SERVICE_ACCOUNT}"
    ```
1. install riff
    ```
    duffle install myriff https://storage.googleapis.com/projectriff/riff-cnab/snapshots/riff-bundle-latest.json --bundle-is-file -s node_port=true -d k8s
    ```
    where `node_port=true` parameter changes all service types to NodePort from LoadBalancer
    and `-d k8s` uses the duffle kubernetes driver to run the installer image in kubernetes cluster
1. You should now be able to see riff components installed on your kubernetes cluster:
    ```
    kubectl get pods --all-namespaces
    ```

#### Docker for Desktop
1. Install the latest docker-for-desktop edge release for your platform. [mac](https://hub.docker.com/editions/community/docker-ce-desktop-mac) OR [windows](https://hub.docker.com/editions/community/docker-ce-desktop-windows)
1. Start the [kubernetes cluser](https://docs.docker.com/docker-for-mac/#kubernetes)
1. Confirm that your kubectl context is pointing to the new cluster
    ```
    kubectl config current-context
    ```
1. Install duffle by [building from source](https://github.com/deislabs/duffle/blob/master/docs/developing.md)

1. Set environment variables required by duffle k8s driver and create the the service account
    ```
    export SERVICE_ACCOUNT=duffle-runtime
    export KUBE_NAMESPACE=duffle
    kubectl create namespace $KUBE_NAMESPACE
    kubectl create serviceaccount "${SERVICE_ACCOUNT}" -n "${KUBE_NAMESPACE}"
    kubectl create clusterrolebinding "${SERVICE_ACCOUNT}-cluster-admin" --clusterrole cluster-admin --serviceaccount "${KUBE_NAMESPACE}:${SERVICE_ACCOUNT}"
    ```
1. install riff
    ```
    duffle install myriff https://storage.googleapis.com/projectriff/riff-cnab/snapshots/riff-bundle-latest.json --bundle-is-file -s node_port=true -d k8s
    ```
    where `node_port=true` parameter changes all service types to NodePort from LoadBalancer 
    and `-d k8s` uses the duffle kubernetes driver to run the installer image in kubernetes cluster
1. You should now be able to see riff components installed on your kubernetes cluster:
    ```
    kubectl get pods --all-namespaces
    ```

#### GKE
1. Initialize gcloud cli by running `gcloud init`
1. Enable the necessary APIs for gcloud.
    ```
    gcloud services enable \
      cloudapis.googleapis.com \
      container.googleapis.com \
      containerregistry.googleapis.com
    ```
1. Create a GKE cluster
    ```
    gcloud container clusters create my-riff-cluster \
      --cluster-version=latest \
      --machine-type=n1-standard-2 \
      --enable-autoscaling --min-nodes=1 --max-nodes=3 \
      --enable-autorepair \
      --scopes=service-control,service-management,compute-rw,storage-ro,cloud-platform,logging-write,monitoring-write,pubsub,datastore \
      --num-nodes=3
    ```
1. Confirm that your kubectl context is pointing to the new cluster
    ```
    kubectl config current-context
    ```
1. Install duffle by [building from source](https://github.com/deislabs/duffle/blob/master/docs/developing.md)

1. Set environment variables required by duffle k8s driver and create the the service account
    ```
    export SERVICE_ACCOUNT=duffle-runtime
    export KUBE_NAMESPACE=duffle
    kubectl create namespace $KUBE_NAMESPACE
    kubectl create serviceaccount "${SERVICE_ACCOUNT}" -n "${KUBE_NAMESPACE}"
    kubectl create clusterrolebinding "${SERVICE_ACCOUNT}-cluster-admin" --clusterrole cluster-admin --serviceaccount "${KUBE_NAMESPACE}:${SERVICE_ACCOUNT}"
    ```
1. install riff
    ```
    duffle install myriff https://storage.googleapis.com/projectriff/riff-cnab/snapshots/riff-bundle-latest.json --bundle-is-file -d k8s
    ```
    and `-d k8s` uses the duffle kubernetes driver to run the installer image in kubernetes cluster

1. You should now be able to see riff components installed on your kubernetes cluster:
    ```
    kubectl get pods --all-namespaces
    ```

## Uninstall
To uninstall, set the SERVICE_ACCOUNT and KUBE_NAMESPACE environment variables as above and use the command:
```
duffle uninstall myriff -d k8s
```
