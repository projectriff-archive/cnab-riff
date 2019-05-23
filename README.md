[![Build Status](https://dev.azure.com/projectriff/projectriff/_apis/build/status/projectriff.cnab-riff?branchName=master)](https://dev.azure.com/projectriff/projectriff/_build/latest?definitionId=16&branchName=master)

# cnab-riff
A CNAB bundle for installing riff (experimental)

## Getting started
To install this cnab bundle, you can use any available cnab runtime. The following are installation instructions using [duffle](https://duffle.sh/).

### Steps
1. Get the latest duffle release for your operating system from duffle's [release page](https://github.com/deislabs/duffle/releases).
1. Start a kubernetes cluster and ensure that your kubeconfig file is configured to talk to that cluster. If using `minikube`, please run the below command before starting `minikube`:
    ```
    minikube config set embed-certs true
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
1. install riff
    ```
    duffle install myriff riff-bundle-latest.json --bundle-is-file --credentials myk8s --insecure
    ```

You should now be able to see riff components installed on your kubernetes cluster.
