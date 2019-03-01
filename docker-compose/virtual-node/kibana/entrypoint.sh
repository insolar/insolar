#!/bin/bash

if [ "$1" == "inside" ]; then
    sleep 30
    curl -XPUT --header 'Content-Type: application/json' 'elasticsearch:9200/_snapshot/original_backup' --data '{ "type": "fs", "settings": { "location": "/tmp/backup" } }'
    # curl -XPUT --header 'Content-Type: application/json' 'elasticsearch:9200/_snapshot/original_backup/snapshot?wait_for_completion=true' --data '{ "indices": ".kibana_1", "ignore_unavailable": true, "include_global_state": false }'
    curl -XPOST --header 'Content-Type: application/json' 'elasticsearch:9200/.kibana_1/_close'
    curl -XPOST --header 'Content-Type: application/json' 'elasticsearch:9200/_snapshot/original_backup/snapshot/_restore?wait_for_completion=true'
    curl -XDELETE --header 'Content-Type: application/json' 'elasticsearch:9200/_snapshot/original_backup'
else
    /usr/local/bin/entrypoint.sh inside &
    /usr/local/bin/kibana-docker
fi
