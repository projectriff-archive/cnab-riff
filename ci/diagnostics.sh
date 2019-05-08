
#!/bin/bash

kubectl get deployments,services,pods --all-namespaces || true
echo "FAILING PODS:"
kubectl get pods --all-namespaces --field-selector=status.phase!=Running \
| tail -n +2 | awk '{print "-n", $1, $2}' | xargs -L 1 kubectl describe pod || true
echo "NODE:"
kubectl describe node || true