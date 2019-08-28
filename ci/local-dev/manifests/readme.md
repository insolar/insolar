1. All work goes through
    ``` 
    ./ci/local-dev/launcher
    ```
2. Run this script only from root of out repository
3. You need installed docker desktop with switched on [k8s](https://kubernetes.io)
4. All possible commands are available by launching **./ci/local-dev/launcher** without params:
    ```
    $ ./ci/local-dev/launcher
    Use cases:
    prepare - install ksonnet
    rebuild - rebuild source code only
    full-rebuild - prepare dependencies and rebuild source code
    start-nodes - deploy and start discovery nodes and pulsar
    start-all - like 'start-nodes' + start elk, jaeger, prometheus ( and grafana - not completely implemented )
    stop - stop all containers
    
    builded images
    pre - image with dependencies
    base - image with binaries
    ```

5. If you use it first time, it's required to do 
    ```
    ./ci/local-dev/launcher prepare
    ```
6. k8s helpers:
    ```
     - kubectl get pods  # to see running pods
     - kubectl exec -ti seed-0 -c insolard  -- bash. # login to first insolard node
     - kubectl logs seed-0 insolard  # logs of insolard on first node
    ```
7. Services:
     - [Jaeger]( http://localhost:30686 )
     - [Kibana]( http://localhost:30601 )
     - [Elasticsearch](http://localhost:30200 )
     - [Prometheus]( http://localhost:30090 )

8. Config files are in subdirectories of 
     ```
     ci/local-dev/manifests/ks/components/
     ```
     For example, config of **insolard** and ***pulsar***
     ```
     ci/local-dev/manifests/ks/components/insolar/insolard_config.libsonnet # insolard config
     ci/local-dev/manifests/ks/components/pulsar/pulsar_config.libsonnet    # pulsar config
     ```
9. Launch benchmark:
     ```
     kubectl exec -ti seed-0 -c insolard  -- bash  # login to first insolard node
     /go/bin/benchmark -c 2 -r 10 -k=/opt/insolar/config/ -a "http://seed-0:19001/admin-api/rpc -p "http://seed-0:19101/api/rpc". # start benchmark
     ```
