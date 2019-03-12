#!/bin/bash

# ./bin/prepare-templates.sh

kubectl delete -f manifests/ve-test.yaml
kubectl apply -f manifests/ve-test.yaml
