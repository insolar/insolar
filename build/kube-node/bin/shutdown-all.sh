#!/bin/bash
shopt -s nullglob
for each in .out/*.yaml ; do
    kubectl delete -f $each
done
# kubectl delete -f .out/*.yaml
for each in manifests/*.yaml ; do
    kubectl delete -f $each
done
