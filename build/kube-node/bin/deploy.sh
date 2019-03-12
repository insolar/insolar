#!/bin/bash
set -e
./bin/prepare-templates.sh

set +e
kubectl delete -f .out/ve-sts.yaml 2>/dev/null

set -e
kubectl apply -f .out/ve-sts.yaml

# set -x
# kubectl get pods
echo "kubectl hints:"
echo ""
echo "kubectl describe pod insolar-ve-0"
echo "kubectl logs -f insolar-ve-0 -c insolard"
echo 'kubectl logs insolar-ve-0 -c genesis"
