Keeperd
===============

Keeperd is a daemon that can check if network is available to process new requests. 

Keeperd has its own api with `/check` handler.
Response format:

    {
        "available": true,
    }


Usage
----------
#### Build

    make keeperd
   
#### Run

    keeperd --config=<path_to_config_file>
    
### Config example

    log:
      level: info
    keeper:
      listenaddress: ':12012'
      faketrue: false
      pollperiod: 5s
      queryurl: 'https://prometheus.k8s-dev.insolar.io/api/v1/query?query='
      queries:
        - 'go_memstats_heap_inuse_bytes{installation="dev-alpha"} <bool 6000000000'
        - 'insolar_requests_opened{installation="dev-alpha"} - insolar_requests_closed{installation="dev-alpha"} <bool 5000'
        - 'rate(insolar_filament_length{installation="dev-alpha"}[1m]) <bool 300000'
      maxmetriclag: 2m
