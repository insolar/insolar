#!/bin/bash
# set -x
mkdir -p .out

SECRETS_SOURCE=${SECRETS_SOURCE:-"./.secrets"}
kubectl delete configmap node-secrets
kubectl create configmap node-secrets --from-file=${SECRETS_SOURCE}
kubectl get configmaps node-secrets -o yaml > .out/secrets.yaml

# kubectl get configmaps node-secrets -o yaml
echo "you could check secrets:"
echo "kubectl get configmaps node-secrets -o yaml"
