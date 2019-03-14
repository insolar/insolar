# How to run virtual executor node in Kubernetes

This doc describes how to run VE-node in Kubernetes outside of Insolar network.

## Required software

* kubectl https://kubernetes.io/docs/tasks/tools/install-kubectl/
* envsubst (part of gnu gettext)

## Digital Ocean

1. Setup DO Kubernetes with single node (droplet). Configure your local environment – make sure your kubectl works properly with DO Kubernetes cluster.

2. Configure float IP in DO, and point it on node's droplet.

3. Delete firewall from node's droplet.

4. Wait until node will be reacheble over internet (i.e. ping from local host, until it works)

5. Prepare values

    export EXTERNAL_IP=174.138.105.150
    ./bin/prepare-templates.sh

6. run service

    kubectl delete service insolar-ve
    kubectl apply -f .out/service.yaml
    kubectl describe svc insolar-ve

7. Test is port forwarding works (requires netcat)

    kubectl delete -f manifests/ve-test.yaml
    kubectl apply -f manifests/ve-test.yaml

    kubectl describe pod insolar-ve-0
    kubectl logs insolar-ve-0

    echo 'Hi, TCP!' | nc $EXTERNAL_IP 30000
    echo 'Hi, UDP!' | nc -u $EXTERNAL_IP 30001

8. Prepare secrets

    mkdir -p .secrets
    cat your-cert-file.json > .secrets/cert.json
    cat your-keys-file.json > .secrets/keys.json
    ./bin/prepare-secrets.sh

check is all ok, output should contain cert and keys data:

    kubectl get configmaps node-secrets -o yaml

9. Run node

    kubectl delete -f manifests/ve-test.yaml
    ./bin/deploy.sh

check pod state

    kubectl get pods
    kubectl describe pod insolar-ve-0
    kubectl logs insolar-ve-0 -c insolard

10. Shut down and remove all created components

    ./bin/shutdown-all.sh

## GCP

TODO

## AWS

TODO

## Private

TODO

