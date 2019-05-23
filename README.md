[![Build Status](https://dev.azure.com/projectriff/projectriff/_apis/build/status/projectriff.cnab-riff?branchName=master)](https://dev.azure.com/projectriff/projectriff/_build/latest?definitionId=16&branchName=master)

# cnab-riff
A CNAB bundle for installing riff (experimental)

## Getting started
To install this cnab bundle, you can use any available cnab runtime. The following are installation instructions using [duffle](https://duffle.sh/).

### Steps
1. Download the latest duffle release for your operating system from duffle's [release page](https://github.com/deislabs/duffle/releases),
 put it somewhere on your path, and make it executable using, for example, `chmod +x /usr/local/bin/duffle`.
1. Start a kubernetes cluster as described in the riff [getting started](https://projectriff.io/docs/getting-started/) documentation,
but stop once you've created a cluster and, if appropriate, given yourself cluster-admin permissions.
 If using `minikube`, please run the following command before starting `minikube`:
    ```
    minikube config set embed-certs true
    ```
1. Check that your kubeconfig file is configured to talk to that cluster:
    ```
    kubectl config current-context
    ```
 
1. Add a duffle credential to allow this bundle to talk to your kubernetes cluster.
    * Create a file named `myk8s.yaml` with the following content:
     ```
     name: myk8s
     credentials:
     - name: kubeconfig
       source:
         path: $HOME/.kube/config
     ```
    * Create a duffle credential using the above file:
     ```
     duffle credentials add myk8s.yaml
     ```
1. Download the latest snapshot of the bundle file:
    ```
    curl -O https://storage.googleapis.com/projectriff/riff-cnab/snapshots/riff-bundle-latest.json
    ```
1. Install riff
    ```
    duffle install myriff riff-bundle-latest.json --bundle-is-file --credentials myk8s --insecure
    ```
1. You should now be able to see riff components installed on your kubernetes cluster:
    ```
    kubectl get pods --all-namespaces
    ```
